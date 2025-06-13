// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"testing"

	"backend.brokedaear.com/internal/common/tests/assert"
	"backend.brokedaear.com/internal/common/tests/test"
)

func TestEnvironmentConfig(t *testing.T) {
	t.Run(
		"string", func(t *testing.T) {
			tests := []struct {
				test.CaseBase
				e Environment
			}{
				{
					CaseBase: newTestCaseBase(
						"is development",
						"DEVELOPMENT",
						false,
					),
					e: EnvDevelopment,
				},
				{
					CaseBase: newTestCaseBase("is staging", "STAGING", false),
					e:        EnvStaging,
				},
				{
					CaseBase: newTestCaseBase(
						"is production",
						"PRODUCTION",
						false,
					),
					e: EnvProduction,
				},
				{
					CaseBase: newTestCaseBase("is ci", "CI", false),
					e:        EnvCI,
				},
				{
					CaseBase: newTestCaseBase("is invalid", "INVALID", false),
					e:        4,
				},
			}
			for _, tt := range tests {
				t.Run(
					tt.Name, func(t *testing.T) {
						assert.Equal(t, tt.e.String(), tt.Want.(string))
					},
				)
			}
		},
	)
	t.Run(
		"validation", func(t *testing.T) {
			tests := []struct {
				test.CaseBase
				e Environment
			}{
				{
					CaseBase: newTestCaseBase(
						"valid Environment",
						nil,
						false,
					),
					e: EnvDevelopment,
				},
				{
					CaseBase: newTestCaseBase(
						"invalid Environment",
						nil,
						true,
					),
					e: 100,
				},
			}
			for _, tt := range tests {
				t.Run(
					tt.Name, func(t *testing.T) {
						err := tt.e.Validate()
						assert.ErrorOrNoError(t, err, tt.WantErr)
					},
				)
			}
		},
	)
}

func TestPortConfig(t *testing.T) {
	t.Run(
		"validation", func(t *testing.T) {
			tests := []struct {
				test.CaseBase
				p Port
			}{
				{
					CaseBase: newTestCaseBase(
						"valid Port lower bound",
						nil,
						false,
					),
					p: 1024,
				},
				{
					CaseBase: newTestCaseBase(
						"valid Port upper bound",
						nil,
						false,
					),
					p: 65533,
				},
				{
					CaseBase: newTestCaseBase(
						"invalid Port below range",
						nil,
						true,
					),
					p: 0,
				},
				{
					CaseBase: newTestCaseBase(
						"invalid Port upper bound",
						nil,
						true,
					),
					p: 65535,
				},
			}
			for _, tt := range tests {
				t.Run(
					tt.Name, func(t *testing.T) {
						err := tt.p.Validate()
						assert.ErrorOrNoError(t, err, tt.WantErr)
					},
				)
			}
		},
	)
}

func TestAddressConfig(t *testing.T) {
	t.Run(
		"validation",
		func(t *testing.T) {
			t.Run(
				"should error", func(t *testing.T) {
					tests := []struct {
						test.CaseBase
						a Address
					}{
						{
							CaseBase: newTestCaseBase(
								"empty Address",
								ErrInvalidAddressLength,
								true,
							),
							a: "",
						},
						{
							CaseBase: newTestCaseBase(
								"Address with colon",
								ErrInvalidAddressColon,
								true,
							),
							a: "127.0.0.1:8080",
						},
						{
							CaseBase: newTestCaseBase(
								"Address with path",
								ErrInvalidAddressWithPath,
								true,
							),
							a: "dingdong.com/api/v1",
						},
					}
					for _, tt := range tests {
						t.Run(
							tt.Name, func(t *testing.T) {
								got := tt.a.Validate()
								assert.ErrorAndWant(t, got, tt.WantErr)
							},
						)
					}
				},
			)
			t.Run(
				"should pass", func(t *testing.T) {
					tests := []struct {
						test.CaseBase
						a Address
					}{
						{
							CaseBase: newTestCaseBase(
								"just hostname",
								nil,
								false,
							),
							a: "localhost",
						},
						{
							CaseBase: newTestCaseBase(
								"hostname with TLD",
								nil,
								false,
							),
							a: "shaboingboing.com",
						},
					}
					for _, tt := range tests {
						t.Run(
							tt.Name, func(t *testing.T) {
								got := tt.a.Validate()
								assert.NoErrorAndNoWant(t, got, tt.WantErr)
							},
						)
					}
				},
			)
		},
	)
}

func TestVersionConfig(t *testing.T) {
	t.Run(
		"validation", func(t *testing.T) {
			tests := []struct {
				test.CaseBase
				v Version
			}{
				{
					CaseBase: newTestCaseBase(
						"valid version",
						"",
						false,
					),
					v: "1.2.3",
				},
				{
					CaseBase: newTestCaseBase(
						"too few elements",
						ErrInvalidVersionFormat,
						true,
					),
					v: "1.2",
				},
				{
					CaseBase: newTestCaseBase(
						"too many elements",
						ErrInvalidVersionFormat,
						true,
					),
					v: "1.2.3.4",
				},
				{
					CaseBase: newTestCaseBase(
						"non-numeric element",
						ErrInvalidVersionChars,
						true,
					),
					v: "1.2.alpha",
				},
				{
					CaseBase: newTestCaseBase(
						"non-numeric element with numeric",
						ErrInvalidVersionChars,
						true,
					),
					v: "1.2.7ae",
				},
				{
					CaseBase: newTestCaseBase(
						"negative number",
						ErrInvalidVersionSign,
						true,
					),
					v: "1.2.-3",
				},
			}
			for _, tt := range tests {
				t.Run(
					tt.Name, func(t *testing.T) {
						err := tt.v.Validate()
						assert.ErrorOrNoError(t, err, tt.WantErr)
					},
				)
			}
		},
	)
}
