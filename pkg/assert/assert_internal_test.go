// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package assert

import (
	"testing"

	"backend.brokedaear.com/pkg/test"
)

func Test_errorAndWant(t *testing.T) {
	t.Run(
		"got error and want error", func(t *testing.T) {
			err := errorAndWant(errFakeNewError, true)
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
			err := errorAndNoWant(errFakeNewError, want)
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

var errFakeNewError fakeError = "fake failure"
