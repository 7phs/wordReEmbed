package options

import (
	"math"
)

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB

	COOREC_SIZE     int64 = 8 + 8
	MEM_LIMIT_COEFF       = 0.1544313298
	EPS                   = 1e-3
)

var (
	CorpusFileName      = "./bin/text8"
	VocabFileName       = "./bin/vocab.txt"
	CooccurenceFileName = "./bin/cooccurrence.bin"

	Op = GloveOptions{
		VocabMinCount: 5,
		VocabMaxCount: 0,

		WindowSize:          15,
		IsDistanceWeighting: true,
		IsSymetric:          true,

		MemoryLimit: 4 * GB,
	}
)

type GloveOptions struct {
	VocabMinCount int
	VocabMaxCount int

	WindowSize          int
	IsDistanceWeighting bool
	IsSymetric          bool

	MemoryLimit int64
}

func (o *GloveOptions) GetMemoryLimit() (int, int) {
	/* The memory_limit determines a limit on the number of elements in bigram_table and the overflow buffer */
	/* Estimate the maximum value that max_product can take so that this limit is still satisfied */
	var (
		overflowLimit = 0.85 * float64(o.MemoryLimit) / float64(COOREC_SIZE)
		maxProduct    = 1e5
	)

	for {
		v := math.Log(maxProduct) + MEM_LIMIT_COEFF
		if math.Abs(overflowLimit-maxProduct*v) <= EPS {
			break
		}

		maxProduct = overflowLimit / v
	}

	return int(maxProduct), int(overflowLimit) / 6
}
