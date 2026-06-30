// Example demonstrating f5core package usage.
package main

import (
	"fmt"

	"github.com/0verkilll/f5core"
)

func main() {
	fmt.Println("=== F5 Core Examples ===")
	fmt.Println()

	// Example 1: De-zigzag table values
	fmt.Println("1. De-Zigzag Table (first 8 values):")
	for i := 0; i < 8; i++ {
		fmt.Printf("   DeZigZag[%d] = %d\n", i, f5core.DeZigZag[i])
	}
	fmt.Println()

	// Example 2: Apply de-zigzag transformation
	fmt.Println("2. ApplyDeZigZag Transformation:")
	testIndices := []int{0, 1, 2, 64, 65, 66, 128}
	for _, idx := range testIndices {
		result := f5core.ApplyDeZigZag(idx)
		block := idx / 64
		posInBlock := idx % 64
		fmt.Printf("   Index %3d (block %d, pos %2d) -> %3d\n",
			idx, block, posInBlock, result)
	}
	fmt.Println()

	// Example 3: F5 Algorithm Constants
	fmt.Println("3. F5 Algorithm Constants:")
	fmt.Printf("   MaxMessageSize:    %d bytes (%.2f MB)\n",
		f5core.MaxMessageSize, float64(f5core.MaxMessageSize)/1024/1024)
	fmt.Printf("   DefaultMaxFileSize: %d bytes (%.2f MB)\n",
		f5core.DefaultMaxFileSize, float64(f5core.DefaultMaxFileSize)/1024/1024)
	fmt.Printf("   HeaderSize:        %d bits\n", f5core.HeaderSize)
	fmt.Printf("   K Parameter Range: %d - %d\n",
		f5core.MinKParameter, f5core.MaxKParameter)
	fmt.Printf("   BlockSize:         %d coefficients\n", f5core.BlockSize)
	fmt.Printf("   Coefficient Range: %d to %d\n",
		f5core.CoefficientMin, f5core.CoefficientMax)
	fmt.Println()

	// Example 4: Validate K parameter
	fmt.Println("4. K Parameter Validation:")
	for k := 0; k <= 9; k++ {
		valid := k >= f5core.MinKParameter && k <= f5core.MaxKParameter
		status := "invalid"
		if valid {
			status = "valid"
		}
		fmt.Printf("   k=%d: %s\n", k, status)
	}
	fmt.Println()

	// Example 5: Calculate code word length for each k
	fmt.Println("5. Matrix Encoding Code Word Lengths:")
	fmt.Println("   (n = 2^k - 1 coefficients needed per k bits)")
	for k := f5core.MinKParameter; k <= f5core.MaxKParameter; k++ {
		n := (1 << k) - 1
		fmt.Printf("   k=%d: n=%3d coefficients\n", k, n)
	}
}
