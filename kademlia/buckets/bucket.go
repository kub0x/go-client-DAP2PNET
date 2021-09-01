package buckets

import (
	"dap2pnet/client/models"
	"math"
	"math/big"
)

type Bucket struct {
	Triplets map[string]models.Triplet
	maxLen   int
	curLen   int
}

func NewBucket(bitOrder int) *Bucket {
	maxLen := 25
	if bitOrder < 5 {
		maxLen = int(math.Pow(2, float64(bitOrder)))
	}

	return &Bucket{
		maxLen:   maxLen,
		curLen:   0,
		Triplets: make(map[string]models.Triplet, maxLen),
	}
}

func (buck *Bucket) Add(triplet models.Triplet, distance big.Int) {
	if buck.curLen < buck.maxLen {
		triplet.Distance = distance
		buck.Triplets[triplet.ID] = triplet
		buck.curLen++
	} else {
		// TODO HANDLE PEER REMOVAL
		println("EXCEEDED MAX CAPACITY")
	}

}
