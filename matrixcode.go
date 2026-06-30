package f5core

// CodeWordLengths maps the F5 matrix-encoding parameter k to the number of
// coefficients n in a single (1, n, k) codeword, where n = 2^k - 1.
//
// The F5 (1, n, k) matrix code embeds k message bits into n coefficients by
// changing at most one coefficient. Index 0 is reserved (k=0 is not a valid
// parameter); the valid range is [MinKParameter, MaxKParameter].
//
//	k=1 ->   1 coefficient  per k=1 message bit   (no matrix encoding)
//	k=2 ->   3 coefficients per k=2 message bits
//	k=3 ->   7 coefficients per k=3 message bits
//	k=4 ->  15 coefficients per k=4 message bits
//	k=5 ->  31 coefficients per k=5 message bits
//	k=6 ->  63 coefficients per k=6 message bits
//	k=7 -> 127 coefficients per k=7 message bits
//	k=8 -> 255 coefficients per k=8 message bits
//
// These values match Westfeld's original F5 specification (2001) and the
// reference Java implementation used by PixelKnot and F5.jar.
var CodeWordLengths = [MaxKParameter + 1]int{0, 1, 3, 7, 15, 31, 63, 127, 255}

// CodeWordLength returns the number of coefficients in a single (1, n, k)
// F5 codeword for the given k parameter, i.e. n = 2^k - 1.
//
// Parameters:
//   - k: matrix-encoding parameter in [MinKParameter, MaxKParameter]
//
// Returns:
//   - The codeword length n = 2^k - 1 for valid k
//
// Panics:
//
//	If k is outside [MinKParameter, MaxKParameter]. The package has no
//	error-return convention, and an invalid k always indicates a caller
//	bug (header validation should reject invalid k before reaching this
//	function).
//
// Example:
//
//	CodeWordLength(1) // returns 1
//	CodeWordLength(3) // returns 7
//	CodeWordLength(8) // returns 255
func CodeWordLength(k int) int {
	if k < MinKParameter || k > MaxKParameter {
		panic("f5core: CodeWordLength called with k outside [MinKParameter, MaxKParameter]")
	}
	return CodeWordLengths[k]
}

// RequiredCapacityBits returns the number of DCT coefficients required to
// embed messageBits bits under the F5 (1, n, k) matrix code.
//
// With matrix encoding each codeword carries k message bits and consumes
// n = 2^k - 1 coefficients. The caller must provide enough usable
// coefficients (non-DC, non-zero after permutation) to cover this capacity;
// shrinkage during embedding may require additional coefficients in
// practice, so callers doing sizing estimates should include a margin.
//
// Formula:
//
//	codewords      = ceil(messageBits / k)
//	requiredCoeffs = codewords * (2^k - 1)
//
// Parameters:
//   - messageBits: total number of message bits to embed (must be >= 0)
//   - k: matrix-encoding parameter in [MinKParameter, MaxKParameter]
//
// Returns:
//   - The number of coefficients required. Returns 0 when messageBits == 0.
//
// Panics:
//
//	If k is outside [MinKParameter, MaxKParameter] or messageBits is
//	negative. Both indicate caller bugs.
//
// Example:
//
//	// k=1: 24 message bits -> 24 codewords of length 1 -> 24 coefficients.
//	RequiredCapacityBits(24, 1) // returns 24
//
//	// k=3: 21 message bits -> 7 codewords of length 7 -> 49 coefficients.
//	RequiredCapacityBits(21, 3) // returns 49
func RequiredCapacityBits(messageBits, k int) int {
	if messageBits < 0 {
		panic("f5core: RequiredCapacityBits called with negative messageBits")
	}
	if k < MinKParameter || k > MaxKParameter {
		panic("f5core: RequiredCapacityBits called with k outside [MinKParameter, MaxKParameter]")
	}
	if messageBits == 0 {
		return 0
	}
	codewords := (messageBits + k - 1) / k // ceil(messageBits / k)
	return codewords * CodeWordLengths[k]
}
