package glove

import (
	"container/ring"
	"github.com/7phs/wordReEmbed/console/options"
	"github.com/7phs/wordReEmbed/lib/storages"
)

type Cooccur struct {
	vocab VocabCount

	isSymetric bool

	history *ring.Ring
	weights []float64

	bigramStorage storages.BigramStorage
}

func NewCooccur(vocab VocabCount, bigramStorage storages.BigramStorage, op *options.GloveOptions) *Cooccur {
	return (&Cooccur{
		vocab: vocab,

		isSymetric: op.IsSymetric,

		history: ring.New(int(op.WindowSize)),

		bigramStorage: bigramStorage,
	}).init(op)
}

func (o *Cooccur) init(op *options.GloveOptions) *Cooccur {
	weights := make([]float64, op.WindowSize)

	for i := range weights {
		weights[i] = 1.

		if op.IsDistanceWeighting {
			weights[i] /= float64(op.WindowSize - i)
		}
	}

	o.weights = weights

	return o
}

func (o *Cooccur) NewLine() {
	for i, ln := 0, o.history.Len(); i < ln; i++ {
		o.history.Value = nil
		o.history.Next()
	}
}

func (o *Cooccur) Push(word string) {
	wordCount, ok := o.vocab.Get(word)
	if !ok {
		return
	}

	index := 0
	weights := o.weights[len(o.weights)-o.history.Len():]
	wIndex2 := wordCount.Index

	o.history.Do(func(v interface{}) {
		if v == nil {
			return
		}
		wIndex := v.(int)

		o.bigramStorage.WeightInc(wIndex, wIndex2, weights[index])
		if o.isSymetric {
			o.bigramStorage.WeightInc(wIndex2, wIndex, weights[index])
		}

		index++
	})

	o.history.Value = wIndex2
	o.history = o.history.Next()
}
