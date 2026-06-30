package f5core

// F5 Algorithm Constants define fixed values for the F5 steganography algorithm.
const (
	// MaxMessageSize is the maximum message size in bytes that can be embedded.
	// This is limited by the 23-bit file size field in the F5 header (2^23 - 1).
	MaxMessageSize = (1 << 23) - 1 // 8,388,607 bytes

	// DefaultMaxFileSize is the default maximum allowed size for extracted data.
	// This prevents memory exhaustion attacks from malformed headers.
	// Default: 10MB (10,485,760 bytes)
	DefaultMaxFileSize = 10_485_760

	// CoefficientMin is the minimum valid JPEG quantized DCT coefficient value.
	CoefficientMin int16 = -2048

	// CoefficientMax is the maximum valid JPEG quantized DCT coefficient value.
	CoefficientMax int16 = 2047

	// HeaderSize is the number of bits in the F5 message header.
	// The header contains the k parameter (8 bits) and file size (23 bits).
	HeaderSize = 32

	// MaxKParameter is the maximum valid k parameter value for matrix encoding.
	// k=8 uses (1,255,8) codes: 255 coefficients to extract 8 bits.
	MaxKParameter = 8

	// MinKParameter is the minimum valid k parameter value for matrix encoding.
	// k=1 uses simple (1,1,1) codes: 1 coefficient per bit (no matrix encoding).
	MinKParameter = 1

	// BlockSize is the size of a JPEG DCT block (8x8 = 64 coefficients).
	BlockSize = 64

	// FileSizeMask extracts bits 0-22 for file size (23 bits = max ~8MB).
	FileSizeMask = 0x7FFFFF

	// KParameterShift is the bit position of the k parameter (bits 24-31).
	KParameterShift = 24

	// KParameterMask extracts the k parameter after shifting.
	KParameterMask = 0xFF

	// KModulus ensures k wraps within valid range (for Java compatibility).
	KModulus = 32
)
