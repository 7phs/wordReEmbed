package glove

import (
	"sort"
	"strings"
)

type WordCount struct {
	Word  string
	Count int
	Index int
}

type VocabCount map[string]*WordCount

func NewVocabCount() VocabCount {
	return VocabCount(make(map[string]*WordCount))
}

func (o VocabCount) Len() int {
	return len(o)
}

func (o VocabCount) Get(word string) (*WordCount, bool) {
	wc, ok := o[word]
	return wc, ok
}

func (o *VocabCount) Add(word string) {
	if wordCount, ok := (*o)[word]; !ok {
		(*o)[word] = &WordCount{
			Word:  word,
			Count: 1,
		}
	} else {
		wordCount.Count++
	}
}

func (o *VocabCount) AddCount(wc *WordCount) {
	(*o)[wc.Word] = wc
}

func (o VocabCount) Sort(min, max int) (result []*WordCount) {
	result = make([]*WordCount, 0, len(o))

	for _, wc := range o {
		if min > 0 && wc.Count < min {
			continue
		}

		result = append(result, wc)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Count > result[j].Count {
			return true
		}

		if result[i].Count == result[j].Count {
			return strings.Compare(result[i].Word, result[j].Word) < 0
		}

		return false
	})

	if max > 0 && len(result) > max {
		result = result[:max]
	}

	return
}
