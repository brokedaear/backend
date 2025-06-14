// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package server_test

import (
	"testing"

	"backend.brokedaear.com/internal/core/server"
	"backend.brokedaear.com/internal/common/tests/assert"
	"backend.brokedaear.com/internal/common/tests/test"
)

func TestPortConfig(t *testing.T) {
	t.Run(
		"validation", func(t *testing.T) {
			tests := []struct {
				test.CaseBase
				p server.Port
			}{
				{
					CaseBase: test.NewCaseBase(
						"valid Port lower bound",
						nil,
						false,
					),
					p: 1024,
				},
				{
					CaseBase: test.NewCaseBase(
						"valid Port upper bound",
						nil,
						false,
					),
					p: 65533,
				},
				{
					CaseBase: test.NewCaseBase(
						"invalid Port below range",
						nil,
						true,
					),
					p: 0,
				},
				{
					CaseBase: test.NewCaseBase(
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
						a server.Address
					}{
						{
							CaseBase: test.NewCaseBase(
								"empty Address",
								server.ErrInvalidAddressLength,
								true,
							),
							a: "",
						},
						{
							CaseBase: test.NewCaseBase(
								"Address with colon",
								server.ErrInvalidAddressColon,
								true,
							),
							a: "127.0.0.1:8080",
						},
						{
							CaseBase: test.NewCaseBase(
								"Address with path",
								server.ErrInvalidAddressWithPath,
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
						a server.Address
					}{
						{
							CaseBase: test.NewCaseBase(
								"just hostname",
								nil,
								false,
							),
							a: "localhost",
						},
						{
							CaseBase: test.NewCaseBase(
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
				v server.Version
			}{
				{
					CaseBase: test.NewCaseBase(
						"valid version",
						"",
						false,
					),
					v: "1.2.3",
				},
				{
					CaseBase: test.NewCaseBase(
						"too few elements",
						server.ErrInvalidVersionFormat,
						true,
					),
					v: "1.2",
				},
				{
					CaseBase: test.NewCaseBase(
						"too many elements",
						server.ErrInvalidVersionFormat,
						true,
					),
					v: "1.2.3.4",
				},
				{
					CaseBase: test.NewCaseBase(
						"non-numeric element",
						server.ErrInvalidVersionChars,
						true,
					),
					v: "1.2.alpha",
				},
				{
					CaseBase: test.NewCaseBase(
						"non-numeric element with numeric",
						server.ErrInvalidVersionChars,
						true,
					),
					v: "1.2.7ae",
				},
				{
					CaseBase: test.NewCaseBase(
						"negative number",
						server.ErrInvalidVersionSign,
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
