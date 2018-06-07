package storages

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/7phs/wordReEmbed/console/options"
	"github.com/7phs/wordReEmbed/lib/model/glove"
)

func StoreVocab(fileName string, vocab glove.VocabCount, op *options.GloveOptions) (int, error) {
	v := vocab.Sort(op.VocabMinCount, op.VocabMaxCount)

	vocabFile, err := os.Create(fileName)
	if err != nil {
		return len(v), errors.New("Error open file: '" + options.VocabFileName + "': " + err.Error())
	}
	defer vocabFile.Close()

	w := bufio.NewWriter(vocabFile)
	defer w.Flush()

	for _, word := range v {
		w.WriteString(fmt.Sprintf("%s %d\n", word.Word, word.Count))
	}

	return len(v), nil
}

func ReadVocab(fileName string) (glove.VocabCount, error) {
	vocab, err := OpenTextReader(fileName, '\n')
	if err != nil {
		return nil, errors.New("Error open file: '" + fileName + "': " + err.Error())
	}
	defer vocab.Close()

	var (
		vocabCount = glove.NewVocabCount()
		index      = 1
		count      = 0
	)

	err = nil

	vocab.ForEach(func(line []byte) bool {
		parts := bytes.Split(line, []byte{' '})
		if len(parts) != 2 {
			return true
		}

		count, err = strconv.Atoi(string(parts[1]))
		if err != nil {
			err = errors.New("failed to parse vocab files err: " + err.Error())
			return false
		}

		vocabCount.AddCount(&glove.WordCount{
			Word:  string(parts[0]),
			Count: count,
			Index: index,
		})

		index++

		return true
	})

	return vocabCount, nil
}
