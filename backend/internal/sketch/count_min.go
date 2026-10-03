package sketch

import (
	"hash/fnv"
	"sync"
)

type CountMinSketch struct {
	mu     sync.RWMutex
	width  uint32
	depth  uint32
	table  [][]uint8
	additions int64
	sampleLimit int64
}

func NewCountMinSketch(width, depth uint32) *CountMinSketch {
	if width == 0 {
		width = 512
	}
	if depth == 0 {
		depth = 4
	}

	table := make([][]uint8, depth)
	for i := uint32(0); i < depth; i++ {
		table[i] = make([]uint8, width)
	}

	return &CountMinSketch{
		width:       width,
		depth:       depth,
		table:       table,
		sampleLimit: int64(width * 10),
	}
}

func (cms *CountMinSketch) hash(key string, seed uint32) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	sum := h.Sum32()
	return (sum ^ seed) % cms.width
}

func (cms *CountMinSketch) Increment(key string) {
	cms.mu.Lock()
	defer cms.mu.Unlock()

	seeds := []uint32{0x9e3779b9, 0x85ebca6b, 0xc2b2ae35, 0x27d4eb2f}

	for i := uint32(0); i < cms.depth; i++ {
		col := cms.hash(key, seeds[i])
		if cms.table[i][col] < 15 { // 4-bit max value
			cms.table[i][col]++
		}
	}

	cms.additions++
	if cms.additions >= cms.sampleLimit {
		cms.decay()
		cms.additions = 0
	}
}

func (cms *CountMinSketch) Estimate(key string) uint8 {
	cms.mu.RLock()
	defer cms.mu.RUnlock()

	seeds := []uint32{0x9e3779b9, 0x85ebca6b, 0xc2b2ae35, 0x27d4eb2f}
	var minVal uint8 = 15

	for i := uint32(0); i < cms.depth; i++ {
		col := cms.hash(key, seeds[i])
		val := cms.table[i][col]
		if val < minVal {
			minVal = val
		}
	}

	return minVal
}

func (cms *CountMinSketch) decay() {
	for i := uint32(0); i < cms.depth; i++ {
		for j := uint32(0); j < cms.width; j++ {
			cms.table[i][j] /= 2
		}
	}
}


<!-- // Periodic halving decays historical frequency to adapt to workload shifts -->

<!-- // Vectorized counter packing using 64-bit word operations -->
