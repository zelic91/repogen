package spec_test

import (
	"testing"

	"github.com/zelic91/repogen/internal/code"
	"github.com/zelic91/repogen/internal/spec"
)

type ErrorTestCase struct {
	Name           string
	Error          error
	ExpectedString string
}

func TestError(t *testing.T) {
	testTable := []ErrorTestCase{
		{
			Name:           "UnknownOperationError",
			Error:          spec.NewUnknownOperationError("Search"),
			ExpectedString: "unknown operation 'Search'",
		},
		{
			Name:           "StructFieldNotFoundError",
			Error:          spec.NewStructFieldNotFoundError([]string{"Phone", "Number"}),
			ExpectedString: "struct field 'PhoneNumber' not found",
		},
		{
			Name:           "UnsupportedReturnError",
			Error:          spec.NewUnsupportedReturnError(code.SimpleType("User"), 0),
			ExpectedString: "return type 'User' at index 0 is not supported",
		},
		{
			Name:           "OperationReturnCountUnmatchedError",
			Error:          spec.NewOperationReturnCountUnmatchedError(2),
			ExpectedString: "operation requires return count of 2",
		},
		{
			Name:           "InvalidQueryError",
			Error:          spec.NewInvalidQueryError([]string{"And"}),
			ExpectedString: "invalid query 'And'",
		},
		{
			Name: "IncompatibleComparatorError",
			Error: spec.NewIncompatibleComparatorError(spec.ComparatorTrue, code.StructField{
				Name: "Age",
				Type: code.TypeInt,
			}),
			ExpectedString: "cannot use comparator EQUAL_TRUE with struct field 'Age' of type 'int'",
		},
		{
			Name:           "InvalidSortError",
			Error:          spec.NewInvalidSortError([]string{"Order", "By"}),
			ExpectedString: "invalid sort 'OrderBy'",
		},
		{
			Name:           "ArgumentTypeNotMatchedError",
			Error:          spec.NewArgumentTypeNotMatchedError("Age", code.TypeInt, code.TypeFloat64),
			ExpectedString: "field 'Age' requires an argument of type 'int' (got 'float64')",
		},
		{
			Name: "IncompatibleUpdateOperatorError",
			Error: spec.NewIncompatibleUpdateOperatorError(spec.UpdateOperatorInc, spec.FieldReference{
				code.StructField{
					Name: "City",
					Type: code.TypeString,
				},
			}),
			ExpectedString: "cannot use update operator INC with struct field 'City' of type 'string'",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			if testCase.Error.Error() != testCase.ExpectedString {
				t.Errorf("Expected = %+v\nReceived = %+v", testCase.ExpectedString, testCase.Error.Error())
			}
		})
	}
}
