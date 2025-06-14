// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package validator_test

import (
	"testing"

	"backend.brokedaear.com/internal/common/tests/assert"
	"backend.brokedaear.com/internal/common/tests/test"
	"backend.brokedaear.com/internal/common/validator"
)

func TestCheck(t *testing.T) {
	tests := []struct {
		test.CaseBase
		args []validator.Verifiable
	}{
		{
			CaseBase: test.CaseBase{
				Name:    "no types",
				Want:    nil,
				WantErr: true,
			},
			args: []validator.Verifiable{},
		},
		{
			CaseBase: test.CaseBase{
				Name:    "valid types",
				Want:    nil,
				WantErr: false,
			},
			args: []validator.Verifiable{
				fakeValidType{"a"},
				fakeValidType{"b"},
				fakeValidType{"c"},
			},
		},
		{
			CaseBase: test.CaseBase{
				Name:    `invalid type"`,
				Want:    nil,
				WantErr: true,
			},
			args: []validator.Verifiable{
				fakeValidType{"a"},
				fakeInvalidType{"b"},
				fakeValidType{"c"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.Name, func(t *testing.T) {
				err := validator.Check(tt.args...)
				assert.ErrorOrNoError(t, err, tt.WantErr)
			},
		)
	}

	t.Run(
		"error at b not d", func(t *testing.T) {
			fake := fakeInvalidType{"b"}
			err := validator.Check(
				fakeValidType{"a"},
				fake,
				fakeValidType{"c"},
				fakeInvalidType{"d"},
			)
			assert.Error(t, err, errFakeInvalidTypeError)
			assert.Equal(t, fake.Value(), "b")
		},
	)
}

type fakeValidType struct {
	name string
}

func (f fakeValidType) Validate() error {
	return nil
}

func (f fakeValidType) Value() any {
	return f.name
}

type fakeInvalidType struct {
	name string
}

func (f fakeInvalidType) Validate() error {
	return errFakeInvalidTypeError
}

func (f fakeInvalidType) Value() any {
	return f.name
}

type fakeInvalidTypeError string

func (f fakeInvalidTypeError) Error() string {
	return string(f)
}

const errFakeInvalidTypeError fakeInvalidTypeError = "invalid type"
