package storages

import "sort"

type Chunk struct {
	chunk      []*BigramRec
	chunkIndex int
	limit      int
}

func NewChunk(capacity, limit int) *Chunk {
	return &Chunk{
		chunk: make([]*BigramRec, 0, capacity),
		limit: capacity - limit,
	}
}

func (o *Chunk) Reset(chunkIndex int) {
	o.chunk = o.chunk[:0]
	o.chunkIndex = chunkIndex
}

func (o *Chunk) ChunkIndex() int {
	return o.chunkIndex
}

func (o *Chunk) Len() int {
	return len(o.chunk)
}

func (o *Chunk) IsOverflow() bool {
	// magic value 16, need maybe more accurate
	return len(o.chunk) >= o.limit
}

func (o *Chunk) Add(wordIndex, wordIndex2 int, weight float64) {
	o.chunk = append(o.chunk, &BigramRec{
		wordIndexS: int64(wordIndex)<<32 | int64(wordIndex2),
		weight:     weight,
	})
}

func (o *Chunk) Sort() {
	sort.Slice(o.chunk, func(i, j int) bool {
		return o.chunk[i].wordIndexS < o.chunk[j].wordIndexS
	})
}

func (o *Chunk) Records() []*BigramRec {
	return o.chunk
}
