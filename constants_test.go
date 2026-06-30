package f5core

import (
	"testing"
)

func TestMaxMessageSize(t *testing.T) {
	// MaxMessageSize should be 2^23 - 1 (23-bit field)
	expected := (1 << 23) - 1
	if MaxMessageSize != expected {
		t.Errorf("MaxMessageSize = %d, want %d", MaxMessageSize, expected)
	}
	if MaxMessageSize != 8388607 {
		t.Errorf("MaxMessageSize = %d, want 8388607", MaxMessageSize)
	}
}

func TestDefaultMaxFileSize(t *testing.T) {
	// DefaultMaxFileSize should be 10MB
	expected := 10_485_760
	if DefaultMaxFileSize != expected {
		t.Errorf("DefaultMaxFileSize = %d, want %d", DefaultMaxFileSize, expected)
	}
}

func TestCoefficientRange(t *testing.T) {
	// JPEG DCT coefficients are 12-bit signed values
	if CoefficientMin != -2048 {
		t.Errorf("CoefficientMin = %d, want -2048", CoefficientMin)
	}
	if CoefficientMax != 2047 {
		t.Errorf("CoefficientMax = %d, want 2047", CoefficientMax)
	}
	// Should be symmetric around zero (minus 1)
	if int(CoefficientMax) != -int(CoefficientMin)-1 {
		t.Errorf("Coefficient range not symmetric: min=%d, max=%d", CoefficientMin, CoefficientMax)
	}
}

func TestHeaderSize(t *testing.T) {
	// F5 header is 32 bits (8 bits k + 24 bits reserved with 23 bits for size)
	if HeaderSize != 32 {
		t.Errorf("HeaderSize = %d, want 32", HeaderSize)
	}
}

func TestKParameterRange(t *testing.T) {
	if MinKParameter != 1 {
		t.Errorf("MinKParameter = %d, want 1", MinKParameter)
	}
	if MaxKParameter != 8 {
		t.Errorf("MaxKParameter = %d, want 8", MaxKParameter)
	}
	if MinKParameter > MaxKParameter {
		t.Errorf("MinKParameter (%d) > MaxKParameter (%d)", MinKParameter, MaxKParameter)
	}
}

func TestBlockSize(t *testing.T) {
	// JPEG DCT block is 8x8 = 64 coefficients
	if BlockSize != 64 {
		t.Errorf("BlockSize = %d, want 64", BlockSize)
	}
}

func TestFileSizeMask(t *testing.T) {
	// FileSizeMask should extract 23 bits
	expected := 0x7FFFFF
	if FileSizeMask != expected {
		t.Errorf("FileSizeMask = 0x%X, want 0x%X", FileSizeMask, expected)
	}
	// Verify it can represent MaxMessageSize
	if MaxMessageSize&FileSizeMask != MaxMessageSize {
		t.Errorf("FileSizeMask cannot represent MaxMessageSize")
	}
}

func TestKParameterShift(t *testing.T) {
	// k parameter is in bits 24-31
	if KParameterShift != 24 {
		t.Errorf("KParameterShift = %d, want 24", KParameterShift)
	}
}

func TestKParameterMask(t *testing.T) {
	// k parameter mask extracts 8 bits
	if KParameterMask != 0xFF {
		t.Errorf("KParameterMask = 0x%X, want 0xFF", KParameterMask)
	}
}

func TestKModulus(t *testing.T) {
	// KModulus for Java compatibility
	if KModulus != 32 {
		t.Errorf("KModulus = %d, want 32", KModulus)
	}
}

func TestHeaderFieldsCompatibility(t *testing.T) {
	// Verify header field layout is compatible
	// Header: [k (8 bits)][reserved (1 bit)][size (23 bits)]
	testCases := []struct {
		k        int
		fileSize int
	}{
		{1, 100},
		{4, 1000},
		{8, MaxMessageSize},
	}

	for _, tc := range testCases {
		header := uint32(tc.k<<KParameterShift) | uint32(tc.fileSize&FileSizeMask)
		extractedK := (header >> KParameterShift) & KParameterMask
		extractedSize := header & FileSizeMask

		if int(extractedK) != tc.k {
			t.Errorf("k extraction failed: got %d, want %d", extractedK, tc.k)
		}
		if int(extractedSize) != tc.fileSize {
			t.Errorf("size extraction failed: got %d, want %d", extractedSize, tc.fileSize)
		}
	}
}

func TestConstantsNotZero(t *testing.T) {
	// Sanity check that constants are defined
	constants := map[string]int{
		"MaxMessageSize":     MaxMessageSize,
		"DefaultMaxFileSize": DefaultMaxFileSize,
		"HeaderSize":         HeaderSize,
		"MaxKParameter":      MaxKParameter,
		"MinKParameter":      MinKParameter,
		"BlockSize":          BlockSize,
		"FileSizeMask":       FileSizeMask,
		"KParameterShift":    KParameterShift,
		"KParameterMask":     KParameterMask,
		"KModulus":           KModulus,
	}

	for name, value := range constants {
		if value == 0 {
			t.Errorf("%s should not be zero", name)
		}
	}
}
