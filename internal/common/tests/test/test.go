// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package test

// CaseBase is the base for table driven test cases.
type CaseBase struct {
	// Name is the name of the test.
	Name string

	// Want is what is wanted from the output.
	Want any

	// WantErr is whether an error should be expected.
	WantErr bool
}

// NewCaseBase creates a new TestCaseBase.
func NewCaseBase(name string, want any, wantErr bool) CaseBase {
	return CaseBase{
		Name:    name,
		Want:    want,
		WantErr: wantErr,
	}
}
