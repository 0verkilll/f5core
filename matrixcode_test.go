package f5core

import (
	"testing"
)

// TestCodeWordLength_SpecVectors verifies that CodeWordLength returns the
// canonical F5 spec values 2^k-1 for k=1..8. These values must never drift:
// they are encoded in every F5 message header and any change would break
// bit-compatibility with PixelKnot/F5.jar.
func TestCodeWordLength_SpecVectors(t *testing.T) {
	tests := []struct {
		k    int
		want int
	}{
		{1, 1},
		{2, 3},
		{3, 7},
		{4, 15},
		{5, 31},
		{6, 63},
		{7, 127},
		{8, 255},
	}

	for _, tt := range tests {
		got := CodeWordLength(tt.k)
		if got != tt.want {
			t.Errorf("CodeWordLength(%d) = %d, want %d (2^%d - 1)", tt.k, got, tt.want, tt.k)
		}
		// Table should match function
		if CodeWordLengths[tt.k] != tt.want {
			t.Errorf("CodeWordLengths[%d] = %d, want %d", tt.k, CodeWordLengths[tt.k], tt.want)
		}
	}
}

func TestCodeWordLengths_IndexZeroReserved(t *testing.T) {
	if CodeWordLengths[0] != 0 {
		t.Errorf("CodeWordLengths[0] = %d, want 0 (reserved sentinel)", CodeWordLengths[0])
	}
}

func TestCodeWordLength_PanicsOnInvalidK(t *testing.T) {
	cases := []int{-1, 0, MaxKParameter + 1, 100}
	for _, k := range cases {
		func(k int) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("CodeWordLength(%d) did not panic", k)
				}
			}()
			_ = CodeWordLength(k)
		}(k)
	}
}

// TestRequiredCapacityBits verifies capacity sizing for representative k
// values. Vectors are chosen to exercise exact-fit, ceiling rounding, and
// zero-bit boundary cases.
func TestRequiredCapacityBits(t *testing.T) {
	tests := []struct {
		name        string
		messageBits int
		k           int
		want        int
	}{
		// k=1: 1 codeword per bit, 1 coefficient each.
		{"k=1 24 bits exact", 24, 1, 24},
		{"k=1 1 bit", 1, 1, 1},

		// k=3: 7 coefficients per 3 message bits. 21/3=7 codewords -> 49 coeffs.
		{"k=3 21 bits exact (7 codewords)", 21, 3, 49},
		{"k=3 1 bit rounds up to 1 codeword", 1, 3, 7},
		{"k=3 4 bits rounds up to 2 codewords", 4, 3, 14},

		// k=4: 15 coefficients per 4 bits. 32/4=8 codewords -> 120.
		{"k=4 32 bits exact", 32, 4, 120},
		{"k=4 33 bits -> 9 codewords", 33, 4, 135},

		// k=8: 255 coefficients per 8 bits. 16/8=2 codewords -> 510.
		{"k=8 16 bits exact", 16, 8, 510},
		{"k=8 1 bit rounds up", 1, 8, 255},

		// Zero message: no coefficients needed regardless of k.
		{"zero bits k=1", 0, 1, 0},
		{"zero bits k=8", 0, 8, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RequiredCapacityBits(tt.messageBits, tt.k)
			if got != tt.want {
				t.Errorf("RequiredCapacityBits(%d, %d) = %d, want %d",
					tt.messageBits, tt.k, got, tt.want)
			}
		})
	}
}

func TestRequiredCapacityBits_PanicsOnInvalidInputs(t *testing.T) {
	type args struct {
		messageBits int
		k           int
	}
	cases := []args{
		{-1, 1},                    // negative bits
		{100, 0},                   // k below min
		{100, MaxKParameter + 1},   // k above max
	}
	for _, c := range cases {
		func(c args) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("RequiredCapacityBits(%d, %d) did not panic", c.messageBits, c.k)
				}
			}()
			_ = RequiredCapacityBits(c.messageBits, c.k)
		}(c)
	}
}

// TestApplyDeZigZag_NegativeInput exercises the new guard: ApplyDeZigZag
// panics when given a negative shuffled index, rather than silently
// indexing DeZigZag[-N] (Go's `%` preserves sign, causing an out-of-range
// panic deep inside the lookup table with no diagnostic context).
func TestApplyDeZigZag_NegativeInput(t *testing.T) {
	cases := []int{-1, -64, -65, -1000}
	for _, shuffled := range cases {
		func(shuffled int) {
			defer func() {
				r := recover()
				if r == nil {
					t.Errorf("ApplyDeZigZag(%d) did not panic", shuffled)
					return
				}
				// Surface the message in test output for debuggability.
				t.Logf("ApplyDeZigZag(%d) panicked as expected: %v", shuffled, r)
			}()
			_ = ApplyDeZigZag(shuffled)
		}(shuffled)
	}
}
