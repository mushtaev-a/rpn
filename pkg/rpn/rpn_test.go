package rpn_test

import (
	"testing"

	rpn "github.com/mushtaev-a/rpn/pkg/rpn"
)

func TestCalc(t *testing.T) {
	testCasesSuccess := []struct {
		name           string
		expression     string
		expectedResult float64
	}{
		{
			name:           "simple+",
			expression:     "1+1",
			expectedResult: 2,
		},
		{
			name:           "simpl-_",
			expression:     "1-1",
			expectedResult: 0,
		},
		{
			name:           "simple*",
			expression:     "1*1",
			expectedResult: 1,
		},
		{
			name:           "simple-/",
			expression:     "6/3",
			expectedResult: 2,
		},
		{
			name:           "multiple-equal+",
			expression:     "1+1+1",
			expectedResult: 3,
		},
		{
			name:           "multiple-equal-",
			expression:     "3-2-0",
			expectedResult: 1,
		},
		{
			name:           "multiple-equal*",
			expression:     "1*2*3",
			expectedResult: 6,
		},
		{
			name:           "multiple-equal/",
			expression:     "9/3/3",
			expectedResult: 1,
		},
		{
			name:           "priority",
			expression:     "2+2*3",
			expectedResult: 8,
		},
		{
			name:           "priority2",
			expression:     "2+1+6/2*3",
			expectedResult: 12,
		},
		{
			name:           "priority 3",
			expression:     "2+1-6/2*3",
			expectedResult: -6,
		},
		{
			name:           "priority braces",
			expression:     "(2+2)*2",
			expectedResult: 8,
		},
		{
			name:           "result > 10",
			expression:     "(2*3+4)*2",
			expectedResult: 20,
		},
		{
			name:           "priority braces 2",
			expression:     "(((2+1)-2)*6)*(4-1)",
			expectedResult: 18,
		},
		{
			name:           "priority braces 3",
			expression:     "3-(4*2)",
			expectedResult: -5,
		},
		{
			name:           "/",
			expression:     "1/2",
			expectedResult: 0.5,
		},
	}

	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := rpn.Calc(testCase.expression)
			if err != nil {
				t.Fatalf("successful case %s returns error", testCase.expression)
			}
			if val != testCase.expectedResult {
				t.Fatalf("%f should be equal %f", val, testCase.expectedResult)
			}
		})
	}

	testCasesFail := []struct {
		name        string
		expression  string
		expectedErr error
	}{
		{
			name:       "simple",
			expression: "1+1*",
		},
		{
			name:       "priority",
			expression: "2+2**2",
		},
		{
			name:       "priority",
			expression: "((2+2-*(2",
		},
		{
			name:       "empty",
			expression: "",
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := rpn.Calc(testCase.expression)
			if err == nil {
				t.Fatalf("expression %s is invalid but result  %f was obtained", testCase.expression, val)
			}
		})
	}
}
