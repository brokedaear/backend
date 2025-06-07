// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package test

// CaseBase is the base for table driven test cases.
type CaseBase[T any] struct {
	// Name is the name of the test.
	Name string

	// Want is what is wanted from the output.
	Want T

	// WantErr is whether an error should be expected.
	WantErr bool
}

// NewCaseBase creates a new TestCaseBase.
func NewCaseBase[T any](name string, want T, wantErr bool) CaseBase[T] {
	return CaseBase[T]{
		Name:    name,
		Want:    want,
		WantErr: wantErr,
	}
}

// NoWant is a no-op type for table-test cases that implement
// the TestCaseBase. Use it for generic type instantiation to specify that
// a test case will not use the `Want` property.
type NoWant struct{}
