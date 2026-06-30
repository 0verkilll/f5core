package f5core

import (
	"testing"
)

func FuzzApplyDeZigZag(f *testing.F) {
	// Add seed corpus
	f.Add(0)
	f.Add(1)
	f.Add(63)
	f.Add(64)
	f.Add(65)
	f.Add(127)
	f.Add(128)
	f.Add(1000)
	f.Add(10000)

	f.Fuzz(func(t *testing.T, shuffled int) {
		// Skip negative inputs as they don't make sense for indices
		if shuffled < 0 {
			return
		}

		result := ApplyDeZigZag(shuffled)

		// Result should be non-negative
		if result < 0 {
			t.Errorf("ApplyDeZigZag(%d) = %d, want non-negative", shuffled, result)
		}

		// Result should be in the same block
		inputBlock := shuffled / 64
		resultBlock := result / 64
		if resultBlock != inputBlock {
			t.Errorf("ApplyDeZigZag(%d) = %d, crosses block boundary", shuffled, result)
		}

		// Result should be within valid range for the block
		blockStart := inputBlock * 64
		blockEnd := blockStart + 63
		if result < blockStart || result > blockEnd {
			t.Errorf("ApplyDeZigZag(%d) = %d, not in block range [%d, %d]",
				shuffled, result, blockStart, blockEnd)
		}
	})
}

func FuzzApplyDeZigZagDeterministic(f *testing.F) {
	// Verify that ApplyDeZigZag is deterministic
	f.Add(0)
	f.Add(100)
	f.Add(1000)

	f.Fuzz(func(t *testing.T, shuffled int) {
		if shuffled < 0 {
			return
		}

		result1 := ApplyDeZigZag(shuffled)
		result2 := ApplyDeZigZag(shuffled)

		if result1 != result2 {
			t.Errorf("ApplyDeZigZag(%d) not deterministic: %d vs %d",
				shuffled, result1, result2)
		}
	})
}

func FuzzDeZigZagTableAccess(f *testing.F) {
	// Fuzz test table access with modular arithmetic
	f.Add(0)
	f.Add(32)
	f.Add(63)

	f.Fuzz(func(t *testing.T, pos int) {
		// Normalize to valid table index
		if pos < 0 {
			return
		}
		idx := pos % 64

		value := DeZigZag[idx]

		// Value should be in valid range
		if value < 0 || value > 63 {
			t.Errorf("DeZigZag[%d] = %d, not in range [0, 63]", idx, value)
		}
	})
}

func FuzzHeaderConstruction(f *testing.F) {
	// Fuzz test header field construction and extraction
	f.Add(1, 100)
	f.Add(4, 1000)
	f.Add(8, 8388607)

	f.Fuzz(func(t *testing.T, k int, size int) {
		// Constrain inputs to valid ranges
		if k < MinKParameter || k > MaxKParameter {
			return
		}
		if size < 0 || size > MaxMessageSize {
			return
		}

		// Construct header
		header := uint32(k<<KParameterShift) | uint32(size&FileSizeMask)

		// Extract fields
		extractedK := int((header >> KParameterShift) & KParameterMask)
		extractedSize := int(header & FileSizeMask)

		// Verify roundtrip
		if extractedK != k {
			t.Errorf("k roundtrip failed: input %d, got %d", k, extractedK)
		}
		if extractedSize != size {
			t.Errorf("size roundtrip failed: input %d, got %d", size, extractedSize)
		}
	})
}
