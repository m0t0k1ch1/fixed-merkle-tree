package merkle

import (
	"fmt"
	"hash"
	"math"
)

const (
	DepthMin    = 1
	DepthMax    = 32
	HashSizeMin = 1  // bytes
	HashSizeMax = 64 // bytes
)

var (
	ErrTooSmallDepth    = fmt.Errorf("depth must be %d or more", DepthMin)
	ErrTooLargeDepth    = fmt.Errorf("depth must be %d or less", DepthMax)
	ErrTooSmallHashSize = fmt.Errorf("leaf size must be %d or more", HashSizeMin)
	ErrTooLargeHashSize = fmt.Errorf("leaf size must be %d or less", HashSizeMax)
)

type Config struct {
	hasher       hash.Hash
	depth        uint64
	hashSize     uint64
	allLeavesNum uint64
	allNodesNum  uint64
}

func NewConfig(hasher hash.Hash, depth, hashSize uint64) (*Config, error) {
	if depth < DepthMin {
		return nil, ErrTooSmallDepth
	}
	if depth > DepthMax {
		return nil, ErrTooLargeDepth
	}
	if hashSize < HashSizeMin {
		return nil, ErrTooSmallHashSize
	}
	if hashSize > HashSizeMax {
		return nil, ErrTooLargeHashSize
	}

	allLeavesNum := uint64(math.Exp2(float64(depth)))

	allNodesNum := uint64(0)
	for i := allLeavesNum; i >= 1; i /= 2 {
		allNodesNum += i
	}

	return &Config{
		hasher:       hasher,
		depth:        depth,
		hashSize:     hashSize,
		allLeavesNum: allLeavesNum,
		allNodesNum:  allNodesNum,
	}, nil
}
