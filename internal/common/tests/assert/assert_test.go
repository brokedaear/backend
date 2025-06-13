// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package assert

import (
	"testing"

	"backend.brokedaear.com/internal/common/tests/test"
)

func TestAssertFunctions(t *testing.T) {
	t.Run(
		"on integers", func(t *testing.T) {
			Equal(t, 1, 1)
			NotEqual(t, 1, 2)
		},
	)

	t.Run(
		"on strings", func(t *testing.T) {
			Equal(t, "hello", "hello")
			NotEqual(t, "hello", "Grace")
		},
	)

	t.Run(
		"on booleans", func(t *testing.T) {
			True(t, true)
			False(t, false)
		},
	)

	t.Run(
		"on errors", func(t *testing.T) {
			err := fakeErrNewError
			Error(t, err, fakeErrNewError)
		},
	)

	t.Run(
		"no error", func(t *testing.T) {
			NoError(t, nil)
		},
	)
}

func TestErrorAndWant(t *testing.T) {
	t.Run(
		"want error and got error", func(t *testing.T) {
			ErrorAndWant(t, fakeErrNewError, true)
		},
	)

	t.Run(
		"want no error and got nil", func(t *testing.T) {
			ErrorAndWant(t, nil, false)
		},
	)
}

func Test_errorAndWant(t *testing.T) {
	t.Run(
		"got error and want error", func(t *testing.T) {
			err := errorAndWant(fakeErrNewError, true)
			if err != nil {
				t.Errorf(test.ErrStringFormat, err, true)
			}
		},
	)
	t.Run(
		"got nil and want error", func(t *testing.T) {
			err := errorAndWant(nil, true)
			if err == nil {
				t.Errorf(test.ErrStringFormat, err, true)
			}
		},
	)
}

func Test_errorAndNoWant(t *testing.T) {
	var want bool
	t.Run(
		"got error but don't want error", func(t *testing.T) {
			err := errorAndNoWant(fakeErrNewError, want)
			if err != nil {
				t.Errorf(test.ErrStringFormat, err, want)
			}
		},
	)
	t.Run(
		"got nil and don't want error", func(t *testing.T) {
			err := errorAndNoWant(nil, want)
			if err == nil {
				t.Errorf(test.ErrStringFormat, err, want)
			}
		},
	)
}

func Test_noErrorAndNoWant(t *testing.T) {
	var want bool
	t.Run(
		"got nil and don't want error", func(t *testing.T) {
			err := noErrorAndNoWant(nil, want)
			if err != nil {
				t.Errorf(test.ErrStringFormat, err, want)
			}
		},
	)
	t.Run(
		"got error but don't want error", func(t *testing.T) {
			err := noErrorAndNoWant(nil, want)
			if err != nil {
				t.Errorf(test.ErrStringFormat, err, want)
			}
		},
	)
}

type fakeError string

func (r fakeError) Error() string {
	return string(r)
}

var fakeErrNewError fakeError = "fake failure"
