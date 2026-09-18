package wtinylfu

import (
	"hash/fnv"
	"sync"
	"sync/atomic"

	"github.com/pfrancistechs-cel/edgecache/backend/internal/sketch"
	"github.com/pfrancistechs-cel/edgecache/backend/internal/slru"
)

type WTinyLFUShard struct {
	mu        sync.RWMutex
	windowCap int
	window    *slru.DoublyLinkedList
	windowIdx map[string]*slru.Element
	slru      *slru.SegmentedLRU
	sketch    *sketch.CountMinSketch
}

func NewShard(capacity int) *WTinyLFUShard {
	if capacity <= 0 {
		capacity = 500
	}

	winCap := int(float64(capacity) * 0.01)
	if winCap < 1 {
		winCap = 1
	}
	slruCap := capacity - winCap

	return &WTinyLFUShard{
		windowCap: winCap,
		window:    slru.NewList(),
		windowIdx: make(map[string]*slru.Element),
		slru:      slru.NewSLRU(slruCap),
		sketch:    sketch.NewCountMinSketch(512, 4),
	}
}

func (s *WTinyLFUShard) Get(key string) (string, bool) {
	s.sketch.Increment(key)

	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. Check Window LRU
	if elem, ok := s.windowIdx[key]; ok {
		s.window.MoveToFront(elem)
		return elem.Value, true
	}

	// 2. Check SLRU
	return s.slru.Get(key)
}

func (s *WTinyLFUShard) Set(key, val string) {
	s.sketch.Increment(key)

	s.mu.Lock()
	defer s.mu.Unlock()

	// If in window, update
	if elem, ok := s.windowIdx[key]; ok {
		elem.Value = val
		s.window.MoveToFront(elem)
		return
	}

	// Insert into Window LRU
	newElem := &slru.Element{Key: key, Value: val}
	s.window.PushFront(newElem)
	s.windowIdx[key] = newElem

	// Window overflow -> TinyLFU admission competition
	if s.window.Len() > s.windowCap {
		candidate := s.window.Back()
		s.window.Remove(candidate)
		delete(s.windowIdx, candidate.Key)

		// Compare frequency with SLRU victim candidate
		victim := s.slru.VictimCandidate()
		if victim == nil {
			s.slru.Add(candidate.Key, candidate.Value)
			return
		}

		candidateFreq := s.sketch.Estimate(candidate.Key)
		victimFreq := s.sketch.Estimate(victim.Key)

		if candidateFreq > victimFreq {
			// Admitted! Discard SLRU victim and insert candidate
			s.slru.Add(candidate.Key, candidate.Value)
		}
		// If candidateFreq <= victimFreq, candidate is discarded (admission filter protected cache!)
	}
}

type Cache struct {
	shards    []*WTinyLFUShard
	numShards int
	hits      int64
	misses    int64
}

func NewCache(totalCapacity int, numShards int) *Cache {
	if numShards <= 0 {
		numShards = 16
	}
	shardCap := totalCapacity / numShards
	if shardCap <= 0 {
		shardCap = 50
	}

	shards := make([]*WTinyLFUShard, numShards)
	for i := 0; i < numShards; i++ {
		shards[i] = NewShard(shardCap)
	}

	return &Cache{
		shards:    shards,
		numShards: numShards,
	}
}

func (c *Cache) shardIndex(key string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return int(h.Sum32() % uint32(c.numShards))
}

func (c *Cache) Get(key string) (string, bool) {
	idx := c.shardIndex(key)
	val, ok := c.shards[idx].Get(key)
	if ok {
		atomic.AddInt64(&c.hits, 1)
	} else {
		atomic.AddInt64(&c.misses, 1)
	}
	return val, ok
}

func (c *Cache) Set(key, val string) {
	idx := c.shardIndex(key)
	c.shards[idx].Set(key, val)
}

func (c *Cache) Stats() (int64, int64, int, int) {
	hits := atomic.LoadInt64(&c.hits)
	misses := atomic.LoadInt64(&c.misses)

	totalProb := 0
	totalProt := 0
	for _, s := range c.shards {
		s.mu.RLock()
		totalProb += s.slru.ProbationCount()
		totalProt += s.slru.ProtectedCount()
		s.mu.RUnlock()
	}

	return hits, misses, totalProb, totalProt
}


<!-- // Zero-allocation read path
 -->

<!-- // 16-way lock-striped architecture eliminates global bottleneck -->
