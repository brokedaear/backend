// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"testing"

	"backend.brokedaear.com/internal/common/tests/assert"
)

func TestEnvironmentConfig(t *testing.T) {
	t.Run(
		"string", func(t *testing.T) {
			tests := []struct {
				testCaseBase[string]
				e Environment
			}{
				{
					testCaseBase: newTestCaseBase(
						"is development",
						"DEVELOPMENT",
						false,
					),
					e: EnvDevelopment,
				},
				{
					testCaseBase: newTestCaseBase("is staging", "STAGING", false),
					e:            EnvStaging,
				},
				{
					testCaseBase: newTestCaseBase(
						"is production",
						"PRODUCTION",
						false,
					),
					e: EnvProduction,
				},
				{
					testCaseBase: newTestCaseBase("is ci", "CI", false),
					e:            EnvCI,
				},
				{
					testCaseBase: newTestCaseBase("is invalid", "INVALID", false),
					e:            4,
				},
			}
			for _, tt := range tests {
				t.Run(
					tt.Name, func(t *testing.T) {
						assert.Equal(t, tt.e.String(), tt.Want)
					},
				)
			}
		},
	)
	t.Run(
		"validation", func(t *testing.T) {
			tests := []struct {
				testCaseBase[testNoWant]
				e Environment
			}{
				{
					testCaseBase: newTestCaseBase(
						"valid Environment",
						testNoWant{},
						false,
					),
					e: EnvDevelopment,
				},
				{
					testCaseBase: newTestCaseBase(
						"invalid Environment",
						testNoWant{},
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
				testCaseBase[testNoWant]
				p Port
			}{
				{
					testCaseBase: newTestCaseBase(
						"valid Port lower bound",
						testNoWant{},
						false,
					),
					p: 1024,
				},
				{
					testCaseBase: newTestCaseBase(
						"valid Port upper bound",
						testNoWant{},
						false,
					),
					p: 65533,
				},
				{
					testCaseBase: newTestCaseBase(
						"invalid Port below range",
						testNoWant{},
						true,
					),
					p: 0,
				},
				{
					testCaseBase: newTestCaseBase(
						"invalid Port upper bound",
						testNoWant{},
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
						testCaseBase[ConfigError]
						a Address
					}{
						{
							testCaseBase: newTestCaseBase(
								"empty Address",
								ErrInvalidAddressLength,
								true,
							),
							a: "",
						},
						{
							testCaseBase: newTestCaseBase(
								"Address with colon",
								ErrInvalidAddressColon,
								true,
							),
							a: "127.0.0.1:8080",
						},
						{
							testCaseBase: newTestCaseBase(
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
						testCaseBase[testNoWant]
						a Address
					}{
						{
							testCaseBase: newTestCaseBase(
								"just hostname",
								testNoWant{},
								false,
							),
							a: "localhost",
						},
						{
							testCaseBase: newTestCaseBase(
								"hostname with TLD",
								testNoWant{},
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
				testCaseBase[ConfigError]
				v Version
			}{
				{
					testCaseBase: newTestCaseBase[ConfigError](
						"valid version",
						"",
						false,
					),
					v: "1.2.3",
				},
				{
					testCaseBase: newTestCaseBase(
						"too few elements",
						ErrInvalidVersionFormat,
						true,
					),
					v: "1.2",
				},
				{
					testCaseBase: newTestCaseBase(
						"too many elements",
						ErrInvalidVersionFormat,
						true,
					),
					v: "1.2.3.4",
				},
				{
					testCaseBase: newTestCaseBase(
						"non-numeric element",
						ErrInvalidVersionChars,
						true,
					),
					v: "1.2.alpha",
				},
				{
					testCaseBase: newTestCaseBase(
						"non-numeric element with numeric",
						ErrInvalidVersionChars,
						true,
					),
					v: "1.2.7ae",
				},
				{
					testCaseBase: newTestCaseBase(
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
