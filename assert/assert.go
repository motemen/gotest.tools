/*Package assert provides assertions for comparing expected values to actual
values in tests. When an assertion fails a helpful error message is printed.

Example usage

The example below shows assert used with some common types and the failure
messages it produces.


	import (
	    "testing"

	    "gotest.tools/v3/assert"
	    is "gotest.tools/v3/assert/cmp"
	)

	func TestEverything(t *testing.T) {
	    // booleans
	    assert.Assert(t, ok)
	    // assertion failed: ok is false
	    assert.Assert(t, !missing)
	    // assertion failed: missing is true

	    // primitives
	    assert.Equal(t, count, 1)
	    // assertion failed: 0 (count int) != 1 (int)
	    assert.Equal(t, msg, "the message")
	    // assertion failed: my message (msg string) != the message (string)
	    assert.Assert(t, total != 10) // use Assert for NotEqual
	    // assertion failed: total is 10
	    assert.Assert(t, count > 20, "count=%v", count)
	    // assertion failed: count is <= 20: count=1

	    // errors
	    assert.NilError(t, closer.Close())
	    // assertion failed: error is not nil: close /file: errno 11
	    assert.Error(t, err, "the exact error message")
	    // assertion failed: expected error "the exact error message", got "oops"
	    assert.ErrorContains(t, err, "includes this")
	    // assertion failed: expected error to contain "includes this", got "oops"
	    assert.ErrorIs(t, err, os.ErrNotExist)
	    // assertion failed: error is "oops" (err *errors.errorString), not "file does not exist" (os.ErrNotExist *errors.errorString)

	    // complex types
	    assert.DeepEqual(t, result, myStruct{Name: "title"})
	    assert.Assert(t, is.Len(items, 3))
	    // assertion failed: expected [] (length 0) to have length 3
	    assert.Assert(t, len(sequence) != 0) // use Assert for NotEmpty
	    // assertion failed: len(sequence) is 0
	    assert.Assert(t, is.Contains(mapping, "key"))
	    // assertion failed: map[other:1] does not contain key

	    // pointers and interface
	    assert.Assert(t, ref == nil)
	    // assertion failed: ref is not nil
	    assert.Assert(t, ref != nil) // use Assert for NotNil
	    // assertion failed: ref is nil
	}

Assert and Check

Assert() and Check() both accept a Comparison, and fail the test when the
comparison fails. The one difference is that Assert() will end the test execution
immediately (using t.FailNow()) whereas Check() will fail the test (using t.Fail()),
return the value of the comparison, then proceed with the rest of the test case.

Comparisons

Package http://pkg.go.dev/gotest.tools/v3/assert/cmp provides
many common comparisons. Additional comparisons can be written to compare
values in other ways. See the example Assert (CustomComparison).

Automated migration from testify

gty-migrate-from-testify is a command which translates Go source code from
testify assertions to the assertions provided by this package.

See http://pkg.go.dev/gotest.tools/v3/assert/cmd/gty-migrate-from-testify.


*/
package assert // import "gotest.tools/v3/assert"

import (
	gocmp "github.com/google/go-cmp/cmp"
	"gotest.tools/v3/assert/cmp"
	"gotest.tools/v3/internal/assert"
)

// BoolOrComparison can be a bool, or cmp.Comparison. See Assert() for usage.
type BoolOrComparison interface{}

// TestingT is the subset of testing.T used by the assert package.
type TestingT interface {
	FailNow()
	Fail()
	Log(args ...interface{})
}

type helperT interface {
	Helper()
}

// Assert performs a comparison. If the comparison fails, the test is marked as
// failed, a failure message is logged, and execution is stopped immediately.
//
// The comparison argument may be one of three types:
//   bool
// True is success. False is a failure.
// The failure message will contain the literal source code of the expression.
//   cmp.Comparison
// Uses cmp.Result.Success() to check for success of failure.
// The comparison is responsible for producing a helpful failure message.
// http://pkg.go.dev/gotest.tools/v3/assert/cmp provides many common comparisons.
//   error
// A nil value is considered success.
// A non-nil error is a failure, err.Error() is used as the failure message.
func Assert(t TestingT, comparison BoolOrComparison, msgAndArgs ...interface{}) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}
	if !assert.Eval(t, assert.ArgsFromComparisonCall, comparison, msgAndArgs...) {
		t.FailNow()
	}
}

// Check performs a comparison. If the comparison fails the test is marked as
// failed, a failure message is logged, and Check returns false. Otherwise returns
// true.
//
// See Assert for details about the comparison arg and failure messages.
func Check(t TestingT, comparison BoolOrComparison, msgAndArgs ...interface{}) bool {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}
	if !assert.Eval(t, assert.ArgsFromComparisonCall, comparison, msgAndArgs...) {
		t.Fail()
		return false
	}
	return true
}

// NilError fails the test immediately if err is not nil.
// This is equivalent to Assert(t, err)
func NilError(t TestingT, err error, msgAndArgs ...interface{}) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}
	if !assert.Eval(t, assert.ArgsAfterT, err, msgAndArgs...) {
		t.FailNow()
	}
}

// Equal uses the == operator to assert two values are equal and fails the test
// if they are not equal.
//
// If the comparison fails Equal will use the variable names for x and y as part
// of the failure message to identify the actual and expected values.
//
// If either x or y are a multi-line string the failure message will include a
// unified diff of the two values. If the values only differ by whitespace
// the unified diff will be augmented by replacing whitespace characters with
// visible characters to identify the whitespace difference.
//
// This is equivalent to Assert(t, cmp.Equal(x, y)).
func Equal(t TestingT, x, y interface{}, msgAndArgs ...interface{}) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}
	if !assert.Eval(t, assert.ArgsAfterT, cmp.Equal(x, y), msgAndArgs...) {
		t.FailNow()
	}
}

// DeepEqual uses google/go-cmp (https://godoc.org/github.com/google/go-cmp/cmp)
// to assert two values are equal and fails the test if they are not equal.
//
// Package http://pkg.go.dev/gotest.tools/v3/assert/opt provides some additional
// commonly used Options.
//
// This is equivalent to Assert(t, cmp.DeepEqual(x, y)).
func DeepEqual(t TestingT, x, y interface{}, opts ...gocmp.Option) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}
	if !assert.Eval(t, assert.ArgsAfterT, cmp.DeepEqual(x, y, opts...)) {
		t.FailNow()
	}
}

// Error fails the test if err is nil, or the error message is not the expected
// message.
// Equivalent to Assert(t, cmp.Error(err, message)).
func Error(t TestingT, err error, message string, msgAndArgs ...interface{}) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}
	if !assert.Eval(t, assert.ArgsAfterT, cmp.Error(err, message), msgAndArgs...) {
		t.FailNow()
	}
}

// ErrorContains fails the test if err is nil, or the error message does not
// contain the expected substring.
// Equivalent to Assert(t, cmp.ErrorContains(err, substring)).
func ErrorContains(t TestingT, err error, substring string, msgAndArgs ...interface{}) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}
	if !assert.Eval(t, assert.ArgsAfterT, cmp.ErrorContains(err, substring), msgAndArgs...) {
		t.FailNow()
	}
}

// ErrorType fails the test if err is nil, or err is not the expected type.
// Equivalent to Assert(t, cmp.ErrorType(err, expected)).
//
// Expected can be one of:
//   func(error) bool
// Function should return true if the error is the expected type.
//   type struct{}, type &struct{}
// A struct or a pointer to a struct.
// Fails if the error is not of the same type as expected.
//   type &interface{}
// A pointer to an interface type.
// Fails if err does not implement the interface.
//   reflect.Type
// Fails if err does not implement the reflect.Type
func ErrorType(t TestingT, err error, expected interface{}, msgAndArgs ...interface{}) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}
	if !assert.Eval(t, assert.ArgsAfterT, cmp.ErrorType(err, expected), msgAndArgs...) {
		t.FailNow()
	}
}

// ErrorIs fails the test if err is nil, or the error does not match expected
// when compared using errors.Is. See https://golang.org/pkg/errors/#Is for
// accepted argument values.
func ErrorIs(t TestingT, err error, expected error, msgAndArgs ...interface{}) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}
	if !assert.Eval(t, assert.ArgsAfterT, cmp.ErrorIs(err, expected), msgAndArgs...) {
		t.FailNow()
	}
}
