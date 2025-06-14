// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package assert_test

import (
	"testing"

	"backend.brokedaear.com/internal/common/tests/assert"
)

func TestAssertFunctions(t *testing.T) {
	t.Run(
		"on integers", func(t *testing.T) {
			assert.Equal(t, 1, 1)
			assert.NotEqual(t, 1, 2)
		},
	)

	t.Run(
		"on strings", func(t *testing.T) {
			assert.Equal(t, "hello", "hello")
			assert.NotEqual(t, "hello", "Grace")
		},
	)

	t.Run(
		"on booleans", func(t *testing.T) {
			assert.True(t, true)
			assert.False(t, false)
		},
	)

	t.Run(
		"on errors", func(t *testing.T) {
			err := errFakeNewError
			assert.Error(t, err, errFakeNewError)
		},
	)

	t.Run(
		"no error", func(t *testing.T) {
			assert.NoError(t, nil)
		},
	)
}

func TestErrorAndWant(t *testing.T) {
	t.Run(
		"want error and got error", func(t *testing.T) {
			assert.ErrorAndWant(t, errFakeNewError, true)
		},
	)

	t.Run(
		"want no error and got nil", func(t *testing.T) {
			assert.ErrorAndWant(t, nil, false)
		},
	)
}

type fakeError string

func (r fakeError) Error() string {
	return string(r)
}

var errFakeNewError fakeError = "fake failure"
