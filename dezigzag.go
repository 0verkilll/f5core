// Package f5core provides core F5 steganography primitives shared across
// f5messageextract, f5messageembed, and f5imagerecover packages.
package f5core

// DeZigZag is the JPEG de-zigzag transformation table.
//
// JPEG stores DCT coefficients in zigzag order to improve compression
// by grouping low-frequency coefficients (which tend to have larger values)
// at the beginning of the sequence. This table maps zigzag indices to
// natural 2D block positions (row-major order).
//
// The F5 algorithm applies this transformation to access coefficients
// in their original block positions during extraction.
//
// Transformation formula for coefficient access:
//
//	zigzag := shuffled - shuffled%64 + DeZigZag[shuffled%64]
//
// Where 'shuffled' is the permuted index and zigzag is the actual
// coefficient index in the original array.
//
// The table maps positions 0-63 within each 8x8 block:
//
//	Position 0  (DC) -> 0  (top-left corner, always skipped in F5)
//	Position 1       -> 1  (horizontal neighbor)
//	Position 2       -> 5  (diagonal)
//	...etc
//
// These values are copied exactly from the Java F5 reference implementation
// to ensure bit-perfect compatibility with PixelKnot and F5.jar.
var DeZigZag = [64]int{
	0, 1, 5, 6, 14, 15, 27, 28,
	2, 4, 7, 13, 16, 26, 29, 42,
	3, 8, 12, 17, 25, 30, 41, 43,
	9, 11, 18, 24, 31, 40, 44, 53,
	10, 19, 23, 32, 39, 45, 52, 54,
	20, 22, 33, 38, 46, 51, 55, 60,
	21, 34, 37, 47, 50, 56, 59, 61,
	35, 36, 48, 49, 57, 58, 62, 63,
}

// ApplyDeZigZag converts a shuffled coefficient index to its de-zigzagged position.
//
// The F5 algorithm applies permutation (Fisher-Yates shuffle) first to determine
// which coefficient to process, then applies de-zigzag transformation to locate
// the actual coefficient in the DCT coefficient array.
//
// The transformation preserves the 8x8 block structure:
//   - shuffled / 64 determines which block
//   - shuffled % 64 determines position within the block
//   - DeZigZag[shuffled % 64] gives the de-zigzagged position within the block
//
// Parameters:
//
//	shuffled - The permuted coefficient index from Fisher-Yates shuffle.
//	           Must be non-negative; negative values will panic to surface
//	           caller bugs before they silently corrupt memory (Go's `%`
//	           operator preserves sign, which would yield a negative
//	           positionInBlock and an out-of-bounds DeZigZag lookup).
//
// Returns:
//
//	The de-zigzagged index suitable for accessing the coefficient array
//
// Panics:
//
//	If shuffled is negative. The package has no error-return convention,
//	and negative indices always indicate a programmer error rather than
//	malformed user input.
//
// Example:
//
//	// Shuffled index 66 is in block 1 (64-127), position 2 within block
//	// DeZigZag[2] = 5, so de-zigzagged index is 64 + 5 = 69
//	zigzag := ApplyDeZigZag(66) // returns 69
func ApplyDeZigZag(shuffled int) int {
	if shuffled < 0 {
		panic("f5core: ApplyDeZigZag called with negative shuffled index")
	}
	blockStart := shuffled - shuffled%64
	positionInBlock := shuffled % 64
	return blockStart + DeZigZag[positionInBlock]
}
