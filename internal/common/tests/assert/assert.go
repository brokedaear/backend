// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package assert

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"

	"backend.brokedaear.com/internal/common/tests/test"
)

// TODO make Equal and NotEqual deep equals!

// Equal is a test helper that compares `got` and `want`
// for equality. If they are inequal, it calls `t.Errorf`.
func Equal[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

// NotEqual is a test helper that compares `got` and `want`
// for inequality. If they are equal, it calls `t.Errorf`.
func NotEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got == want {
		t.Errorf("didn't want %v", got)
	}
}

// True is a test helper that verifies `got` is true. Otherwise,
// it calls `t.Errorf`.
func True(t *testing.T, got bool) {
	t.Helper()
	if !got {
		t.Errorf("got %v, want true", got)
	}
}

// False is a test helper that verifies `got` is false. Otherwise,
// it calls `t.Errorf`.
func False(t *testing.T, got bool) {
	t.Helper()
	if got {
		t.Errorf("got %v, want false", got)
	}
}

// Error compares two errors for equality. Otherwise, it calls `t.Errorf`.
// For two errors to be equal, their underlying type must be equal.
func Error(t *testing.T, got error, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Errorf(test.ErrStringFormat, got, want)
	}
}

func NoError(t *testing.T, got error) {
	t.Helper()
	err := noError(got)
	if err != nil {
		t.Error(err)
	}
}

func noError(got error) error {
	var want error
	if !errors.Is(got, want) {
		return fmt.Errorf(test.ErrStringFormat, got, want)
	}

	return nil
}

func ErrorAndWant(t *testing.T, got error, want bool) {
	t.Helper()
	err := errorAndWant(got, want)
	if err != nil {
		t.Errorf(test.ErrStringFormat, got, want)
	}
}

func errorAndWant(got error, want bool) error {
	if (got != nil) != want {
		return fmt.Errorf("error = %v, wantErr = %v", got, want)
	}

	return nil
}

func ErrorAndNoWant(t *testing.T, got error, want bool) {
	t.Helper()
	err := errorAndNoWant(got, want)
	if err == nil {
		t.Errorf(test.ErrStringFormat, got, want)
	}
}

func errorAndNoWant(got error, want bool) error {
	if (got == nil) != want {
		return fmt.Errorf(test.ErrStringFormat, got, want)
	}

	return nil
}

func NoErrorAndNoWant(t *testing.T, got error, want bool) {
	t.Helper()
	err := noErrorAndNoWant(got, want)
	if err != nil {
		t.Errorf(test.ErrStringFormat, got, want)
	}
}

func noErrorAndNoWant(got error, want bool) error {
	if (got != nil) != want {
		return fmt.Errorf(test.ErrStringFormat, got, want)
	}

	return nil
}

// ErrorOrNoError decides whether to assert an error if an error is wanted
// or to assert no error when an error is not wanted.
func ErrorOrNoError(t *testing.T, got error, want bool) {
	t.Helper()
	err := errorOrNoError(got, want)
	if err != nil {
		t.Error(assertionError(got, want))
	}
}

func errorOrNoError(got error, want bool) error {
	if want {
		return errorAndWant(got, true)
	} else {
		return noError(got)
	}
}

func assertionError(got, want any) error {
	return fmt.Errorf(test.ErrStringFormat, got, want)
}
