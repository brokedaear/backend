// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"backend.brokedaear.com/internal/common/tests/test"
)

func newTestCaseBase(name string, want any, wantErr bool) test.CaseBase {
	return test.NewCaseBase(name, want, wantErr)
}
