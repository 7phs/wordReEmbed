package storages

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/7phs/wordReEmbed/console/options"
)

type BigramRec struct {
	wordIndexS int64
	weight     float64
}

type FileBigramTable struct {
	fileName     string
	fileBaseName string
	fileExt      string

	windowSize int
	maxProduct int
	bigram     []float64
	lookup     []int

	chunkPool  sync.Pool
	chunk      *Chunk
	chunkIndex int

	bufPool sync.Pool

	ch   chan *Chunk
	wait sync.WaitGroup
}

func NewFileBigramTable(fileName string, vocabSize int, op *options.GloveOptions) *FileBigramTable {
	maxProduct, overflowLength := op.GetMemoryLimit()

	fileExt := filepath.Ext(fileName)
	fileBaseName := strings.TrimSuffix(fileName, fileExt)

	fmt.Println("overflowLength=", overflowLength)

	return (&FileBigramTable{
		fileName:     fileName,
		fileBaseName: fileBaseName,
		fileExt:      fileExt,

		windowSize: op.WindowSize,
		maxProduct: maxProduct,
		lookup:     make([]int, vocabSize+1),

		chunkPool: sync.Pool{
			New: func() interface{} {
				return NewChunk(overflowLength, op.WindowSize)
			},
		},
		bufPool: sync.Pool{
			New: func() interface{} {
				buf := bytes.NewBuffer([]byte{})
				buf.Grow(10 * 1024 * 1024)
				return buf
			},
		},

		chunkIndex: 1,
		ch:         make(chan *Chunk),
	}).init(vocabSize)
}

func (o *FileBigramTable) init(vocabSize int) *FileBigramTable {
	o.chunk = o.chunkPool.Get().(*Chunk)

	o.lookup[0] = 1
	for i := 1; i < len(o.lookup); i++ {
		o.lookup[i] = o.lookup[i-1]

		v := o.maxProduct / i
		if v >= vocabSize {
			v = vocabSize
		}

		o.lookup[i] += v
	}

	o.bigram = make([]float64, o.LookupMax())

	o.chunk.Reset(o.chunkIndex)

	o.wait.Add(2)
	for i := 0; i < 2; i++ {
		go o.storeChunk()
	}

	return o
}

func (o *FileBigramTable) LookupMax() int {
	return o.lookup[len(o.lookup)-1]
}

func (o *FileBigramTable) WeightInc(wordIndex, wordIndex2 int, weight float64) {
	if wordIndex < o.maxProduct/wordIndex2 {
		o.bigram[o.lookup[wordIndex-1]+wordIndex2-2] += weight
	} else {
		o.chunk.Add(wordIndex, wordIndex2, weight)

		if o.chunk.IsOverflow() {
			o.flushChunk()
		}
	}
}

func (o *FileBigramTable) flushChunk() {
	if o.chunk.Len() == 0 {
		return
	}

	chunk := o.chunk

	o.chunkIndex++
	o.chunk = o.chunkPool.Get().(*Chunk)
	o.chunk.Reset(o.chunkIndex)

	o.ch <- chunk
}

func (o *FileBigramTable) storeChunk() {
	for chunk := range o.ch {
		func(chunk *Chunk) {
			defer o.chunkPool.Put(chunk)

			fileName := fmt.Sprintf("%s_%04d%s", o.fileBaseName, chunk.ChunkIndex(), o.fileExt)

			fmt.Println("flushChunk: save: ", chunk.Len())

			chunkFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				// TODO: raise an error
				fmt.Println("flushChunk, ", fileName, ":", err)
				return
			}
			defer chunkFile.Close()

			fmt.Println("flushChunk: sort start")
			chunk.Sort()
			fmt.Println("flushChunk: sort finish")

			buf := o.bufPool.Get().(*bytes.Buffer)
			defer o.bufPool.Put(buf)

			buf.Reset()

			var (
				accum BigramRec
				index = 1
			)

			for _, rec := range chunk.Records() {
				if accum.wordIndexS > 0 {
					if accum.wordIndexS == rec.wordIndexS {
						accum.weight += rec.weight
						continue
					} else {
						binary.Write(buf, binary.LittleEndian, accum.wordIndexS)
						binary.Write(buf, binary.LittleEndian, accum.weight)
					}
				}

				accum.wordIndexS = rec.wordIndexS
				accum.weight = rec.weight
				index++

				if index%600*1024 == 0 {
					chunkFile.Write(buf.Bytes())
					buf.Reset()
				}
			}

			binary.Write(buf, binary.LittleEndian, accum.wordIndexS)
			binary.Write(buf, binary.LittleEndian, accum.weight)
			chunkFile.Write(buf.Bytes())

			fmt.Println("flushChunk: finish")
		}(chunk)
	}

	o.wait.Done()
}

func (o *FileBigramTable) flushBigram() {
	o.wait.Add(1)

	go func(lookup []int, bigram []float64) {
		defer o.wait.Done()

		fileName := fmt.Sprintf("%s_%04d%s", o.fileBaseName, 0, o.fileExt)

		fmt.Println("flushBigram: save")

		chunkFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			// TODO: raise an error
			fmt.Println("flushBigram, ", fileName, ":", err)
			return
		}
		defer chunkFile.Close()

		buf := o.bufPool.Get().(*bytes.Buffer)
		defer o.bufPool.Put(buf)
		buf.Reset()

		var (
			x     int64
			ln    = int64(len(lookup))
			index = 0
		)
		for x = 1; x < ln; x++ {
			x64 := x << 32

			for y := 1; y <= lookup[x]-lookup[x-1]; y++ {
				if r := bigram[lookup[x-1]-2+y]; r != 0 {
					binary.Write(buf, binary.LittleEndian, x64|int64(y))
					binary.Write(buf, binary.LittleEndian, r)

					index++

					if index%600*1024 == 0 {
						chunkFile.Write(buf.Bytes())
						buf.Reset()
					}
				}
			}
		}

		if buf.Len() > 0 {
			chunkFile.Write(buf.Bytes())
		}

		fmt.Println("flushBigram: finish")
	}(o.lookup, o.bigram)
}

func (o *FileBigramTable) mergeChunk() {
	// TODO: implement it
}

func (o *FileBigramTable) Flush() {
	o.flushChunk()
	o.flushBigram()

	o.mergeChunk()

	time.Sleep(100 * time.Millisecond)
	close(o.ch)

	o.wait.Wait()
}
