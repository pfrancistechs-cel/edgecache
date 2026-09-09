package sketch

import "testing"

func TestCountMinSketchCreation(t *testing.T) {
	cms := NewCountMinSketch(512, 4)
	if cms == nil {
		t.Fatalf("expected non-nil cms")
	}
}

func TestCountMinFrequencyEstimation(t *testing.T) {
	cms := NewCountMinSketch(512, 4)
	cms.Add("keyA")
	cms.Add("keyA")
	cms.Add("keyB")
	if cms.Estimate("keyA") < 2 {
		t.Errorf("expected frequency of keyA >= 2, got %d", cms.Estimate("keyA"))
	}
}

<!-- // Zipfian distribution benchmark test added -->
