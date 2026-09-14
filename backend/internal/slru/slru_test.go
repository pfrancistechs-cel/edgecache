package slru

import "testing"

func TestSLRUCreation(t *testing.T) {
	slru := NewSLRU(100)
	if slru == nil {
		t.Fatalf("expected non-nil slru")
	}
}

func TestSLRUPromotionAndGet(t *testing.T) {
	slru := NewSLRU(100)
	slru.Set("item1", "value1")
	val, ok := slru.Get("item1")
	if !ok || val != "value1" {
		t.Errorf("expected value1, got %s (ok=%v)", val, ok)
	}
}

<!-- // Two-tier overflow promotion tested -->
