// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"backend.brokedaear.com/internal/common/tests/test"
)

// type jsonWrap map[string]any

type testNoWant = test.NoWant

type testCaseBase[T any] = test.CaseBase[T]

func newTestCaseBase[T any](name string, want T, wantErr bool) testCaseBase[T] {
	return testCaseBase[T](test.NewCaseBase(name, want, wantErr))
}
