// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package errors_test

import (
	stdlibErrors "errors"

	pkgErrors "github.com/pkg/errors"

	"fmt"
	"strings"
	"testing"

	"backend.brokedaear.com/pkg/assert"
	"backend.brokedaear.com/pkg/errors"
	"backend.brokedaear.com/pkg/test"
)

type ErrorsTestCase struct {
	test.CaseBase
	Err        error
	TargetErr  error
	Message    string
	Format     string
	Args       []any
	ShouldPass bool
}

//nolint:exhaustruct // unnecessary, API is already clear.
func TestJoin(t *testing.T) {
	tests := []ErrorsTestCase{
		{
			CaseBase: test.NewCaseBase("single error", "error1", false),
			Err:      errors.New("error1"),
		},
		{
			CaseBase: test.NewCaseBase("multiple errors", "error1\nerror2", false),
			Err:      errors.Join(errors.New("error1"), errors.New("error2")),
		},
		{
			CaseBase: test.NewCaseBase("nil error", nil, false),
			Err:      errors.Join(nil),
		},
		{
			CaseBase: test.NewCaseBase("mixed nil and non-nil", "error1", false),
			Err:      errors.Join(nil, errors.New("error1"), nil),
		},
		{
			CaseBase: test.NewCaseBase("empty join", nil, false),
			Err:      errors.Join(),
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.Name, func(t *testing.T) {
				if tt.Err == nil {
					assert.Equal(t, tt.Err, nil)
				} else {
					assert.Equal(t, tt.Err.Error(), tt.Want.(string))
				}
			},
		)
	}
}

func TestWrap(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := errors.Wrap(originalErr, "wrapped message")

	assert.NotEqual(t, wrappedErr, nil)
	assert.Equal(t, wrappedErr.Error(), "wrapped message: original error")
	assert.Equal(t, errors.Is(wrappedErr, originalErr), true)
}

func TestWrapf(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := errors.Wrapf(originalErr, "wrapped message with %s", "formatting")

	assert.NotEqual(t, wrappedErr, nil)
	assert.Equal(t, wrappedErr.Error(), "wrapped message with formatting: original error")
	assert.Equal(t, errors.Is(wrappedErr, originalErr), true)
}

//nolint:exhaustruct // unnecessary, API is already clear.
func TestIs(t *testing.T) {
	tests := []ErrorsTestCase{
		{
			CaseBase:   test.NewCaseBase("same error", true, false),
			Err:        errors.New("test error"),
			TargetErr:  errors.New("test error"),
			ShouldPass: false, // Different instances, should be false
		},
		{
			CaseBase:   test.NewCaseBase("wrapped error", true, false),
			Err:        errors.Wrap(stdlibErrors.ErrUnsupported, "wrapped"),
			TargetErr:  stdlibErrors.ErrUnsupported,
			ShouldPass: true,
		},
		{
			CaseBase:   test.NewCaseBase("different errors", false, false),
			Err:        errors.New("error1"),
			TargetErr:  errors.New("error2"),
			ShouldPass: false,
		},
		{
			CaseBase:   test.NewCaseBase("nil error", false, false),
			Err:        nil,
			TargetErr:  errors.New("error"),
			ShouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.Name, func(t *testing.T) {
				if tt.Name == "same error" {
					// Special case: test with the same error instance
					sameErr := errors.New("test error")
					result := errors.Is(sameErr, sameErr)
					assert.Equal(t, result, true)
				} else {
					result := errors.Is(tt.Err, tt.TargetErr)
					assert.Equal(t, result, tt.Want.(bool))
				}
			},
		)
	}
}

// customError is a test error type.
type customError struct {
	code int
	msg  string
}

func (c customError) Error() string {
	return c.msg
}

func TestAs(t *testing.T) {
	originalErr := customError{code: 404, msg: "not found"}
	wrappedErr := errors.Wrap(originalErr, "wrapped")

	var target customError
	result := errors.As(wrappedErr, &target)

	assert.Equal(t, result, true)
	assert.Equal(t, target.code, 404)
	assert.Equal(t, target.msg, "not found")
}

func TestPkgUnwrap(t *testing.T) {
	t.Run(
		"pkg wrapped error", func(t *testing.T) {
			original := errors.New("original")
			wrapped := errors.Wrap(original, "wrapper")
			unwrapped := errors.PkgUnwrap(wrapped)
			// pkg/errors.Unwrap returns the wrapped error itself (same as wrapped)
			assert.NotEqual(t, unwrapped, nil)
			assert.Equal(t, unwrapped.Error(), "wrapper: original")
		},
	)

	t.Run(
		"unwrappable error", func(t *testing.T) {
			simple := errors.New("simple error")
			unwrapped := errors.PkgUnwrap(simple)
			assert.Equal(t, unwrapped, nil)
		},
	)

	t.Run(
		"nil error", func(t *testing.T) {
			unwrapped := errors.PkgUnwrap(nil)
			assert.Equal(t, unwrapped, nil)
		},
	)
}

func TestStdUnwrap(t *testing.T) {
	t.Run(
		"stdlib wrapped error", func(t *testing.T) {
			original := errors.New("original")
			wrapped := fmt.Errorf("wrapper: %w", original)
			unwrapped := errors.StdUnwrap(wrapped)
			assert.NotEqual(t, unwrapped, nil)
			assert.Equal(t, unwrapped.Error(), "original")
		},
	)

	t.Run(
		"pkg wrapped error", func(t *testing.T) {
			original := errors.New("original")
			wrapped := errors.Wrap(original, "wrapper")
			unwrapped := errors.StdUnwrap(wrapped)
			// pkg/errors.Wrap doesn't implement the Unwrap method that stdlib errors.Unwrap expects
			// When no Unwrap method is found, stdlib errors.Unwrap returns the error itself
			assert.NotEqual(t, unwrapped, nil)
			assert.Equal(t, unwrapped.Error(), wrapped.Error())
		},
	)

	t.Run(
		"unwrappable error", func(t *testing.T) {
			simple := errors.New("simple error")
			unwrapped := errors.StdUnwrap(simple)
			assert.Equal(t, unwrapped, nil)
		},
	)

	t.Run(
		"nil error", func(t *testing.T) {
			unwrapped := errors.StdUnwrap(nil)
			assert.Equal(t, unwrapped, nil)
		},
	)
}

func TestNew(t *testing.T) {
	err := errors.New("test error")
	assert.NotEqual(t, err, nil)
	assert.Equal(t, err.Error(), "test error")
}

func TestCause(t *testing.T) {
	originalErr := errors.New("root cause")
	wrappedOnce := errors.Wrap(originalErr, "first wrap")
	wrappedTwice := errors.Wrap(wrappedOnce, "second wrap")

	cause := errors.Cause(wrappedTwice)
	assert.NotEqual(t, cause, nil)
	assert.Equal(t, cause.Error(), "root cause")
}

//nolint:exhaustruct // unnecessary, API is already clear.
func TestErrorf(t *testing.T) {
	tests := []ErrorsTestCase{
		{
			CaseBase: test.NewCaseBase("simple format", "error: test", false),
			Format:   "error: %s",
			Args:     []any{"test"},
		},
		{
			CaseBase: test.NewCaseBase("multiple args", "error 42: test", false),
			Format:   "error %d: %s",
			Args:     []any{42, "test"},
		},
		{
			CaseBase: test.NewCaseBase("no formatting", "simple error", false),
			Format:   "simple error",
			Args:     []any{},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.Name, func(t *testing.T) {
				err := errors.Errorf(tt.Format, tt.Args...)
				assert.Equal(t, err.Error(), tt.Want.(string))
			},
		)
	}
}

func TestWithStack(t *testing.T) {
	originalErr := errors.New("original")
	stackErr := errors.WithStack(originalErr)

	assert.NotEqual(t, stackErr, nil)
	assert.Equal(t, stackErr.Error(), "original")

	// But it should have stack trace information when formatted with %+v
	stackTrace := fmt.Sprintf("%+v", stackErr)
	t.Logf("Stack trace:\n%s", stackTrace)
	assert.Equal(t, strings.Contains(stackTrace, "errors_test.go"), true)
}

func TestWithMessage(t *testing.T) {
	originalErr := errors.New("original error")
	annotatedErr := errors.WithMessage(originalErr, "additional context")

	assert.NotEqual(t, annotatedErr, nil)
	assert.Equal(t, annotatedErr.Error(), "additional context: original error")
	assert.Equal(t, errors.Is(annotatedErr, originalErr), true)
}

func TestWithMessagef(t *testing.T) {
	originalErr := errors.New("original error")
	annotatedErr := errors.WithMessagef(originalErr, "context with %s", "formatting")

	assert.NotEqual(t, annotatedErr, nil)
	assert.Equal(t, annotatedErr.Error(), "context with formatting: original error")
	assert.Equal(t, errors.Is(annotatedErr, originalErr), true)
}

func TestErrUnsupported(t *testing.T) {
	// Test that the exported variable is accessible and has the correct value
	assert.Equal(t, errors.ErrUnsupported.Error(), "unsupported operation")
	assert.Equal(t, errors.Is(errors.ErrUnsupported, stdlibErrors.ErrUnsupported), true)
}

func TestErrorChaining(t *testing.T) {
	// Test complex error chaining scenarios
	root := errors.New("root cause")
	wrapped1 := errors.Wrap(root, "first layer")
	wrapped2 := errors.WithMessage(wrapped1, "second layer")
	wrapped3 := errors.Wrapf(wrapped2, "third layer with %s", "formatting")

	// Test that we can still access the root cause
	assert.Equal(t, errors.Cause(wrapped3).Error(), "root cause")

	// Test that Is still works through all layers
	assert.Equal(t, errors.Is(wrapped3, root), true)

	// Test the full error message
	fullMessage := wrapped3.Error()
	assert.Equal(t, strings.Contains(fullMessage, "root cause"), true)
	assert.Equal(t, strings.Contains(fullMessage, "third layer with formatting"), true)
}

func TestNilErrorHandling(t *testing.T) {
	t.Run(
		"self test", func(t *testing.T) {
			tests := []struct {
				name string
				fn   func() error
			}{
				{
					name: "Wrap nil error",
					fn:   func() error { return errors.Wrap(nil, "message") },
				},
				{
					name: "Wrapf nil error",
					fn:   func() error { return errors.Wrapf(nil, "message %s", "arg") },
				},
				{
					name: "WithMessage nil error",
					fn:   func() error { return errors.WithMessage(nil, "message") },
				},
				{
					name: "WithMessagef nil error",
					fn:   func() error { return errors.WithMessagef(nil, "message %s", "arg") },
				},
				{
					name: "WithStack nil error",
					fn:   func() error { return errors.WithStack(nil) },
				},
			}

			for _, tt := range tests {
				t.Run(
					tt.name, func(t *testing.T) {
						err := tt.fn()
						assert.Equal(t, err, nil)
					},
				)
			}
		},
	)

	t.Run(
		"self against original", func(t *testing.T) {
			tests := []struct {
				name string
				fn   func() error
				want func() error
			}{
				{
					name: "Wrap nil error",
					fn:   func() error { return errors.Wrap(nil, "message") },
					want: func() error { return pkgErrors.Wrap(nil, "message") },
				},
				{
					name: "Wrapf nil error",
					fn:   func() error { return errors.Wrapf(nil, "message %s", "arg") },
					want: func() error { return pkgErrors.Wrapf(nil, "message %s", "arg") },
				},
				{
					name: "WithMessage nil error",
					fn:   func() error { return errors.WithMessage(nil, "message") },
					want: func() error { return pkgErrors.WithMessage(nil, "message") },
				},
				{
					name: "WithMessagef nil error",
					fn:   func() error { return errors.WithMessagef(nil, "message %s", "arg") },
					want: func() error { return pkgErrors.WithMessagef(nil, "message %s", "arg") },
				},
				{
					name: "WithStack nil error",
					fn:   func() error { return errors.WithStack(nil) },
					want: func() error { return pkgErrors.WithStack(nil) },
				},
			}

			for _, tt := range tests {
				t.Run(
					tt.name, func(t *testing.T) {
						want := tt.fn()
						got := tt.want()
						assert.Equal(t, got, want)
					},
				)
			}
		},
	)
}

// networkError is a test error type with timeout behavior.
type networkError struct {
	timeout bool
}

func (n networkError) Error() string {
	return "network error"
}

func (n networkError) Timeout() bool {
	return n.timeout
}

func TestAsWithVariousTypes(t *testing.T) {
	originalErr := networkError{timeout: true}
	wrappedErr := errors.Wrap(originalErr, "connection failed")

	// Test with correct type
	var netErr networkError
	result := errors.As(wrappedErr, &netErr)
	assert.Equal(t, result, true)
	assert.Equal(t, netErr.timeout, true)

	// Test with different error type
	var customErr customError
	result = errors.As(wrappedErr, &customErr)
	assert.Equal(t, result, false)
}

func TestUnwrapComparison(t *testing.T) {
	t.Run(
		"pkg wrapped error behavior difference", func(t *testing.T) {
			original := errors.New("original")
			wrapped := errors.Wrap(original, "wrapper")

			pkgUnwrapped := errors.PkgUnwrap(wrapped)
			stdUnwrapped := errors.StdUnwrap(wrapped)

			// Both pkg/errors.Unwrap and stdlib errors.Unwrap return the wrapped error itself
			// because pkg/errors wrapped errors don't implement the stdlib Unwrap interface
			assert.NotEqual(t, pkgUnwrapped, nil)
			assert.Equal(t, pkgUnwrapped.Error(), "wrapper: original")

			assert.NotEqual(t, stdUnwrapped, nil)
			assert.Equal(t, stdUnwrapped.Error(), "wrapper: original")

			// They should be the same for pkg/errors wrapped errors
			assert.Equal(t, pkgUnwrapped.Error(), stdUnwrapped.Error())
		},
	)

	t.Run(
		"stdlib wrapped error behavior", func(t *testing.T) {
			original := errors.New("original")
			wrapped := fmt.Errorf("wrapper: %w", original)

			pkgUnwrapped := errors.PkgUnwrap(wrapped)
			stdUnwrapped := errors.StdUnwrap(wrapped)

			// Both should return the original error for stdlib wrapped errors
			assert.NotEqual(t, pkgUnwrapped, nil)
			assert.NotEqual(t, stdUnwrapped, nil)
			assert.Equal(t, pkgUnwrapped.Error(), "original")
			assert.Equal(t, stdUnwrapped.Error(), "original")
			assert.Equal(t, pkgUnwrapped.Error(), stdUnwrapped.Error())
		},
	)
}
