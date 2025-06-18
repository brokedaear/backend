// Package errors wraps the packages "github.com/pkg/errors" and "errors"
// from the Go standard library. This provides a nice interface to work
// with in code rather than having to import two distinct libraries
// simultaneously and to alias one of the package names to something like
// "errors2".
package errors

import (
	stdlibErrors "errors"
	"fmt"

	pkgErrors "github.com/pkg/errors"
)

// Join aggregates multiple errors into a single error, and does not replace
// Unwrap. This function is a facade for the standard library's errors.Join
// function.
func Join(errs ...error) error {
	return stdlibErrors.Join(errs...)
}

// Wrap wraps an error with a message. It is a facade for pkg/errors Wrap
// function.
func Wrap(err error, message string) error {
	return pkgErrors.Wrap(err, message)
}

// Wrapf wraps an error with a format and arguments, and does not replace
// Join. This function is a facade for pkg/errors Wrapf function.
func Wrapf(err error, format string, args ...any) error {
	return pkgErrors.Wrapf(err, format, args...)
}

// Is compares two errors and returns a bool representing their equality.
// This function wraps the standard library's errors.Is function.
func Is(err, target error) bool {
	return stdlibErrors.Is(err, target)
}

// As sets the value of the error tree in a target and sets the error to
// that target. It then returns true. This function is a facade for the
// standard library's errors.As function.
func As(err error, target any) bool {
	//goland:noinspection ALL
	return stdlibErrors.As(err, target)
}

// PkgUnwrap is a facade for pkg/errors Unwrap function. Another unwrap
// function exists under the standard library. They behave differently.
func PkgUnwrap(err error) error {
	return pkgErrors.Unwrap(err)
}

// StdUnwrap is a facade for the standard library's errors.Unwrap. Another
// Unwrap function exists under the pkg/errors library. They behave
// differently.
func StdUnwrap(err error) error {
	return stdlibErrors.Unwrap(err)
}

// New returns an error that formats as the given text. This function is a
// facade for the standard library's errors.New.
//
// The pkg/errors package provides a New function, however, it is 22x slower
// than both this facade and the standard library's New function. Therefore,
// this package does not provide an alternate implementation of the pkg/errors
// function.
func New(text string) error {
	return stdlibErrors.New(text)
}

// Cause returns the cause of an error, following the source of
// the error stack all the way to the root. The error must implement the
// Causer interface.
func Cause(err error) error {
	return pkgErrors.Cause(err)
}

// Errorf is for formatting an error like `Printf`. It is a facade for
// the standard library's Errorf.
//
// The pkg/errors package has its own implementation of Errorf, however
// synthetic benchmarks show it is ~5.5x slower than the standard library's
// implementation.
func Errorf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

// WithStack returns err with a stack trace. It is a facade for pkg/errors
// WithStack function.
func WithStack(err error) error {
	return pkgErrors.WithStack(err)
}

// WithMessage annotates an error string using a colon `:`. It is a facade
// for the pkg/errors WithMessage function.
func WithMessage(err error, message string) error {
	return pkgErrors.WithMessage(err, message)
}

// WithMessagef annotates an error string using a colon `:` and formats the
// message as well, if arguments are provided. It is a facade for pkg/errors
// WithMessagef function.
func WithMessagef(err error, format string, args ...any) error {
	return pkgErrors.WithMessagef(err, format, args...)
}

// ErrUnsupported is a facade for the standard library's ErrUnsupported error.
var ErrUnsupported = stdlibErrors.ErrUnsupported
