// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package test_test

import (
	"testing"

	"backend.brokedaear.com/internal/common/tests/assert"
	"backend.brokedaear.com/internal/common/tests/test"
)

func TestNewCaseBase(t *testing.T) {
	tests := []struct {
		name          string
		want          any
		wantErr       bool
		expectName    string
		expectWant    any
		expectWantErr bool
	}{
		{
			name:          "simple case",
			want:          "expected value",
			wantErr:       false,
			expectName:    "simple case",
			expectWant:    "expected value",
			expectWantErr: false,
		},
		{
			name:          "error case",
			want:          nil,
			wantErr:       true,
			expectName:    "error case",
			expectWant:    nil,
			expectWantErr: true,
		},
		{
			name:          "numeric case",
			want:          42,
			wantErr:       false,
			expectName:    "numeric case",
			expectWant:    42,
			expectWantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				caseBase := test.NewCaseBase(tt.name, tt.want, tt.wantErr)

				assert.Equal(t, caseBase.Name, tt.expectName)
				if tt.expectWant != nil {
					assert.Equal(t, caseBase.Want, tt.expectWant)
				} else {
					assert.Equal(t, caseBase.Want, nil)
				}
				assert.Equal(t, caseBase.WantErr, tt.expectWantErr)
			},
		)
	}
}

func TestCaseBaseFields(t *testing.T) {
	t.Run("basic fields", func(t *testing.T) {
		name := "test case name"
			want := "expected result"

			caseBase := test.NewCaseBase(name, want, true)

			// Test that all fields are properly set
			assert.Equal(t, caseBase.Name, name)
			assert.Equal(t, caseBase.Want.(string), want)
			assert.Equal(t, caseBase.WantErr, true)
		},
	)

	t.Run(
		"embedded fields", func(t *testing.T) {
			// Test that CaseBase can be embedded in other structs
			type testCase struct {
				test.CaseBase
				input string
				other int
			}

			tc := testCase{
				CaseBase: test.NewCaseBase("embedded case", "result", false),
				input:    "test input",
				other:    123,
			}

			assert.Equal(t, tc.Name, "embedded case")
			assert.Equal(t, tc.Want.(string), "result")
			assert.Equal(t, tc.WantErr, false)
			assert.Equal(t, tc.input, "test input")
			assert.Equal(t, tc.other, 123)
		},
	)
}
