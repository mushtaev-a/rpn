# Expression Calculator

A mathematical expression calculator implementation in Go that supports standard notation with parentheses.

## Overview

This calculator implements basic arithmetic operations using standard mathematical notation. It supports both simple expressions and complex expressions with parentheses.

## Getting Started

### Prerequisites

- Go 1.16 or later
- Git

### Installation

```bash
# Clone the repository
git clone https://github.com/mushtaev-a/rpn.git
cd rpn

# Build the application
go build -o calculator cmd/main.go
```

### Running the Server

You can run the server in two ways:

1. Using the built binary:

```bash
./calculator
```

2. Using go run:

```bash
go run cmd/main.go
```

The server will start on port 8080 by default. You can change the port by setting the PORT environment variable:

```bash
PORT=3000 ./calculator
```

### Health Check

Once the server is running, you can verify it's working:

```bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"expression": "2+2"}'
```

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

## API Usage

The calculator is available as an HTTP service. By default, it runs on port 8080 (configurable via PORT environment variable).

### HTTP Endpoint

```
POST http://localhost:8080/api/v1/calculate
Content-Type: application/json
```

### Request Format

```json
{
  "expression": "your_expression_here"
}
```

### Example Curl Commands

Calculate a simple expression:

```bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"expression": "2+2*3"}'

# Response: {"result":8}
```

Calculate with parentheses:

```bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"expression": "(2+2)*3"}'

# Response: {"result":12}
```

### Error Responses

The API returns the following status codes:

- 400 Bad Request: Invalid JSON or empty expression.
- 405 Method Not Allowed: When using non-POST methods
- 422 Unprocessable Entity: Invalid expressions, division by zero, etc.
- 500 Internal Server Error: Unexpected errors

Error response examples:

```bash
# Empty expression
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"expression": ""}'

# Response: Empty expression

# Invalid operation
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"expression": "2++2"}'

# Response: some of operations are dublicated
```
