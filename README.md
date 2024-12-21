# Expression Calculator

A mathematical expression calculator implementation in Go that supports standard notation with parentheses.

## Overview

This calculator implements basic arithmetic operations using standard mathematical notation. It supports both simple expressions and complex expressions with parentheses.

## Supported Operations

- Addition (+)
- Subtraction (-)
- Multiplication (\*)
- Division (/)
- Parentheses for operation precedence

## Usage

```go
result, err := rpn.Calc("2+2*3")      // Returns 8
result, err := rpn.Calc("(2+2)*3")    // Returns 12
result, err := rpn.Calc("2+1+6/2*3")  // Returns 12
```

## Error Handling

The calculator handles various error cases:

- Division by zero
- Duplicate operation signs (e.g., "2++3")
- Invalid expressions starting or ending with operations
- Empty expressions
- Unmatched parentheses
- Invalid expression format

## Example Expressions

```
Simple Operations:
"1+1" = 2
"6/3" = 2
"1*2*3" = 6

With Operator Precedence:
"2+2*3" = 8
"2+1+6/2*3" = 12
"3-(4*2)" = -5

With Parentheses:
"(2+2)*2" = 8
"(2*3+4)*2" = 20
"(((2+1)-2)*6)*(4-1)" = 18
```

## Testing

The package includes comprehensive test cases covering:

- Basic arithmetic operations
- Complex expressions with multiple operators
- Parentheses handling
- Error cases
