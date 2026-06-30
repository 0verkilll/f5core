# Examples

This directory contains examples demonstrating how to use the `f5core` package.

## Available Examples

| Example | Description |
|---------|-------------|
| [basic](./basic/) | Core functionality demonstration including de-zigzag transformation and constants |

## Running Examples

Each example can be run directly with `go run`:

```bash
cd examples/basic
go run main.go
```

Or run from the repository root:

```bash
go run ./examples/basic
```

## Example Output

Running the basic example produces output like:

```
=== F5 Core Examples ===

1. De-Zigzag Table (first 8 values):
   DeZigZag[0] = 0
   DeZigZag[1] = 1
   ...

2. ApplyDeZigZag Transformation:
   Index   0 (block 0, pos  0) ->   0
   Index   1 (block 0, pos  1) ->   1
   ...

3. F5 Algorithm Constants:
   MaxMessageSize:    8388607 bytes (8.00 MB)
   ...
```

## Creating New Examples

When adding new examples:

1. Create a new directory under `examples/`
2. Add a `main.go` file with a `main` function
3. Add a `go.mod` file with appropriate dependencies
4. Add a `README.md` describing the example
5. Update this README to include the new example in the table
