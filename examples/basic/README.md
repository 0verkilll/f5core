# Basic Example

This example demonstrates the core functionality of the `f5core` package.

## What It Shows

1. **De-Zigzag Table** - The lookup table used to reverse JPEG's zigzag encoding
2. **ApplyDeZigZag** - Transforming shuffled coefficient indices back to natural order
3. **F5 Algorithm Constants** - Configuration values for F5 steganography
4. **K Parameter Validation** - Checking valid encoding parameter values
5. **Code Word Lengths** - Matrix encoding dimensions for different k values

## Running

```bash
go run main.go
```

## Expected Output

```
=== F5 Core Examples ===

1. De-Zigzag Table (first 8 values):
   DeZigZag[0] = 0
   DeZigZag[1] = 1
   DeZigZag[2] = 5
   DeZigZag[3] = 6
   DeZigZag[4] = 14
   DeZigZag[5] = 15
   DeZigZag[6] = 27
   DeZigZag[7] = 28

2. ApplyDeZigZag Transformation:
   Index   0 (block 0, pos  0) ->   0
   Index   1 (block 0, pos  1) ->   1
   Index   2 (block 0, pos  2) ->   5
   Index  64 (block 1, pos  0) ->  64
   Index  65 (block 1, pos  1) ->  65
   Index  66 (block 1, pos  2) ->  69
   Index 128 (block 2, pos  0) -> 128

3. F5 Algorithm Constants:
   MaxMessageSize:    8388607 bytes (8.00 MB)
   DefaultMaxFileSize: 10485760 bytes (10.00 MB)
   HeaderSize:        32 bits
   K Parameter Range: 1 - 8
   BlockSize:         64 coefficients
   Coefficient Range: -2048 to 2047

4. K Parameter Validation:
   k=0: invalid
   k=1: valid
   ...
   k=8: valid
   k=9: invalid

5. Matrix Encoding Code Word Lengths:
   (n = 2^k - 1 coefficients needed per k bits)
   k=1: n=  1 coefficients
   k=2: n=  3 coefficients
   k=3: n=  7 coefficients
   k=4: n= 15 coefficients
   k=5: n= 31 coefficients
   k=6: n= 63 coefficients
   k=7: n=127 coefficients
   k=8: n=255 coefficients
```

## Key Concepts

### De-Zigzag Transformation

JPEG stores DCT coefficients in a zigzag order for efficient run-length encoding.
The `DeZigZag` table maps from zigzag order back to natural 8x8 block order.

### Block Structure

JPEG images are divided into 8x8 blocks (64 coefficients each). The first
coefficient in each block (position 0, 64, 128, etc.) is the DC coefficient,
which is not used for F5 steganography.

### K Parameter

The k parameter controls the efficiency of matrix encoding:
- Higher k = more bits embedded per change
- Higher k = requires more coefficients per code word
