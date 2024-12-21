package rpn

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrDividingByZero             = errors.New("dividing by zero")
	ErrDuplicateOpertaionsSigns   = errors.New("some of operations is dublicated")
	ErrOpertaionsSigns            = errors.New("expression couldn't starts/ends with operation sign")
	ErrExpressionStringEmpty      = errors.New("expression is empty")
	ErrExpressionStringParetheses = errors.New("expression has problem with parentheses")
)

func sum(a, b float64) (float64, error) {
	return a + b, nil
}

func sub(a, b float64) (float64, error) {
	return a - b, nil
}

func mul(a, b float64) (float64, error) {
	return a * b, nil
}

func div(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDividingByZero
	}

	return a / b, nil
}

func checkExpressionOperantions(expression string) (bool, error) {
	reOperations := regexp.MustCompile(`[+\-*/]{2,}`)
	matches := reOperations.FindAllString(expression, -1)

	if len(matches) != 0 {
		return false, ErrDuplicateOpertaionsSigns
	}

	reWrongOperands := regexp.MustCompile(`^[+\-*/]|[+\-*/]$`)
	hasMistakes := reWrongOperands.Match([]byte(expression))

	if hasMistakes {
		return false, ErrOpertaionsSigns
	}

	return !hasMistakes, nil
}

func findParenthesesSubstringsIndices(s string) ([][]int, bool) {
	var result [][]int
	var stack []int

	for i, char := range s {
		if char == '(' {
			stack = append(stack, i)
		} else if char == ')' {
			if len(stack) > 0 {
				start := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if len(stack) == 0 {
					result = append(result, []int{start, i + 1})
				}
			}
		}
	}

	return result, len(stack) == 0
}

func getFuncMap() map[string]func(a, b float64) (float64, error) {
	return map[string]func(a, b float64) (float64, error){
		"+": sum,
		"-": sub,
		"*": mul,
		"/": div,
	}
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

func Calc(expression string) (float64, error) {
	if len(expression) == 0 {
		return 0, ErrExpressionStringEmpty
	}

	isValid, err := checkExpressionOperantions(expression)

	if !isValid {
		return 0, err
	}

	withBraces, isValidParentheses := findParenthesesSubstringsIndices(expression)

	if !isValidParentheses {
		return 0, ErrExpressionStringParetheses
	}

	if len(withBraces) > 0 {
		originalExpression := expression

		for _, borders := range withBraces {
			start, end := borders[0], borders[1]

			expressionForReplace := originalExpression[start:end]
			subExpression := originalExpression[start+1 : end-1]

			res, err := Calc(subExpression)

			if err != nil {
				return 0, err
			}

			expression = strings.Replace(expression, expressionForReplace, strconv.FormatFloat(res, 'f', -1, 64), 1)
		}
	}

	// Stack to hold numbers and operators
	var numStack []float64
	var opStack []string

	// Helper function to apply an operator to the top two numbers in the stack
	applyOp := func() error {
		if len(numStack) < 2 {
			return fmt.Errorf("not enough operands")
		}
		b := numStack[len(numStack)-1]
		a := numStack[len(numStack)-2]
		numStack = numStack[:len(numStack)-2]

		op := opStack[len(opStack)-1]
		opStack = opStack[:len(opStack)-1]

		res, err := getFuncMap()[op](a, b)
		if err != nil {
			return err
		}

		numStack = append(numStack, res)
		return nil
	}

	// Parse the expression
	i := 0
	for i < len(expression) {
		char := expression[i]

		if unicode.IsDigit(rune(char)) || char == '.' {
			// Parse number
			j := i
			for j < len(expression) && (unicode.IsDigit(rune(expression[j])) || expression[j] == '.') {
				j++
			}
			num, _ := strconv.ParseFloat(expression[i:j], 64)
			numStack = append(numStack, num)
			i = j
		} else if char == '(' {
			opStack = append(opStack, string(char))
			i++
		} else if char == ')' {
			for len(opStack) > 0 && opStack[len(opStack)-1] != "(" {
				if err := applyOp(); err != nil {
					return 0, err
				}
			}
			opStack = opStack[:len(opStack)-1]
			i++
		} else if strings.ContainsRune("+-*/", rune(char)) {
			for len(opStack) > 0 && precedence(opStack[len(opStack)-1]) >= precedence(string(char)) {
				if err := applyOp(); err != nil {
					return 0, err
				}
			}
			opStack = append(opStack, string(char))
			i++
		} else {
			i++
		}
	}

	// Apply remaining operators
	for len(opStack) > 0 {
		if err := applyOp(); err != nil {
			return 0, err
		}
	}

	if len(numStack) != 1 {
		return 0, fmt.Errorf("invalid expression")
	}

	return numStack[0], nil
}
