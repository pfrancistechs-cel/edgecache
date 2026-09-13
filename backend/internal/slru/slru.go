package slru

type SegmentedLRU struct {
	probationCap int
	protectedCap int
	probation    *DoublyLinkedList
	protected    *DoublyLinkedList
	index        map[string]*Element
	isProtected  map[string]bool
}

func NewSLRU(totalCap int) *SegmentedLRU {
	if totalCap <= 0 {
		totalCap = 1000
	}

	// 80% Protected, 20% Probation
	protCap := int(float64(totalCap) * 0.8)
	if protCap == 0 {
		protCap = 1
	}
	probCap := totalCap - protCap

	return &SegmentedLRU{
		probationCap: probCap,
		protectedCap: protCap,
		probation:    NewList(),
		protected:    NewList(),
		index:        make(map[string]*Element),
		isProtected:  make(map[string]bool),
	}
}

func (s *SegmentedLRU) Get(key string) (string, bool) {
	elem, exists := s.index[key]
	if !exists {
		return "", false
	}

	if s.isProtected[key] {
		s.protected.MoveToFront(elem)
	} else {
		// Promote from Probation to Protected
		s.probation.Remove(elem)
		s.protected.PushFront(elem)
		s.isProtected[key] = true

		// If protected overflows, demote protected tail back to probation
		if s.protected.Len() > s.protectedCap {
			demoted := s.protected.Back()
			s.protected.Remove(demoted)
			s.probation.PushFront(demoted)
			s.isProtected[demoted.Key] = false
		}
	}

	return elem.Value, true
}

func (s *SegmentedLRU) Add(key, val string) *Element {
	elem, exists := s.index[key]
	if exists {
		elem.Value = val
		_, _ = s.Get(key)
		return nil
	}

	newElem := &Element{Key: key, Value: val}
	s.probation.PushFront(newElem)
	s.index[key] = newElem
	s.isProtected[key] = false

	// If probation overflows, evict candidate
	if s.probation.Len() > s.probationCap {
		victim := s.probation.Back()
		s.probation.Remove(victim)
		delete(s.index, victim.Key)
		delete(s.isProtected, victim.Key)
		return victim
	}

	return nil
}

func (s *SegmentedLRU) VictimCandidate() *Element {
	return s.probation.Back()
}

func (s *SegmentedLRU) ProbationCount() int {
	return s.probation.Len()
}

func (s *SegmentedLRU) ProtectedCount() int {
	return s.protected.Len()
}


<!-- // Optimized demotion lock release
 -->
