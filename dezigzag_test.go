package f5core

import (
	"testing"
)

func TestDeZigZagTableSize(t *testing.T) {
	if len(DeZigZag) != 64 {
		t.Errorf("DeZigZag table size = %d, want 64", len(DeZigZag))
	}
}

func TestDeZigZagTableValues(t *testing.T) {
	// Each position 0-63 should appear exactly once in the table
	seen := make(map[int]bool)
	for i, v := range DeZigZag {
		if v < 0 || v > 63 {
			t.Errorf("DeZigZag[%d] = %d, want value in range [0, 63]", i, v)
		}
		if seen[v] {
			t.Errorf("DeZigZag contains duplicate value %d", v)
		}
		seen[v] = true
	}
}

func TestDeZigZagFirstRow(t *testing.T) {
	// Verify first row of the table (copied from Java F5 reference)
	expected := []int{0, 1, 5, 6, 14, 15, 27, 28}
	for i, want := range expected {
		if DeZigZag[i] != want {
			t.Errorf("DeZigZag[%d] = %d, want %d", i, DeZigZag[i], want)
		}
	}
}

func TestDeZigZagDCPosition(t *testing.T) {
	// DC coefficient (position 0) should map to 0
	if DeZigZag[0] != 0 {
		t.Errorf("DeZigZag[0] = %d, want 0 (DC coefficient)", DeZigZag[0])
	}
}

func TestApplyDeZigZag(t *testing.T) {
	tests := []struct {
		name     string
		shuffled int
		want     int
	}{
		{"block 0, position 0", 0, 0},
		{"block 0, position 1", 1, 1},
		{"block 0, position 2", 2, 5},
		{"block 0, position 63", 63, 63},
		{"block 1, position 0", 64, 64},
		{"block 1, position 2", 66, 69}, // 64 + DeZigZag[2] = 64 + 5 = 69
		{"block 2, position 0", 128, 128},
		{"block 2, position 3", 131, 134},     // 128 + DeZigZag[3] = 128 + 6 = 134
		{"large block", 640, 640},             // block 10, position 0
		{"large block with offset", 645, 655}, // block 10, position 5 -> 640 + DeZigZag[5] = 640 + 15 = 655
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApplyDeZigZag(tt.shuffled)
			if got != tt.want {
				t.Errorf("ApplyDeZigZag(%d) = %d, want %d", tt.shuffled, got, tt.want)
			}
		})
	}
}

func TestApplyDeZigZagBlockPreservation(t *testing.T) {
	// Verify that ApplyDeZigZag preserves block boundaries
	for block := 0; block < 10; block++ {
		blockStart := block * 64
		for pos := 0; pos < 64; pos++ {
			shuffled := blockStart + pos
			result := ApplyDeZigZag(shuffled)

			// Result should be in the same block
			resultBlock := result / 64
			if resultBlock != block {
				t.Errorf("ApplyDeZigZag(%d) = %d, crosses block boundary (block %d vs %d)",
					shuffled, result, block, resultBlock)
			}
		}
	}
}

func TestApplyDeZigZagBijection(t *testing.T) {
	// Within each block, ApplyDeZigZag should produce a bijection (all 64 unique values)
	for block := 0; block < 5; block++ {
		blockStart := block * 64
		seen := make(map[int]bool)

		for pos := 0; pos < 64; pos++ {
			shuffled := blockStart + pos
			result := ApplyDeZigZag(shuffled)

			if seen[result] {
				t.Errorf("Block %d: ApplyDeZigZag produced duplicate result %d", block, result)
			}
			seen[result] = true
		}

		if len(seen) != 64 {
			t.Errorf("Block %d: ApplyDeZigZag produced %d unique values, want 64", block, len(seen))
		}
	}
}

func TestApplyDeZigZagZeroInput(t *testing.T) {
	result := ApplyDeZigZag(0)
	if result != 0 {
		t.Errorf("ApplyDeZigZag(0) = %d, want 0", result)
	}
}

func BenchmarkApplyDeZigZag(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ApplyDeZigZag(i % 1000)
	}
}

func BenchmarkApplyDeZigZagParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ApplyDeZigZag(i % 1000)
			i++
		}
	})
}
