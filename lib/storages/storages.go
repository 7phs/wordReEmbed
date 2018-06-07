package storages

type BigramStorage interface {
	WeightInc(wordIndex, wordIndex2 int, weight float64)
	Flush()
}

