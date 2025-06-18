// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"testing"

	"backend.brokedaear.com"
	"backend.brokedaear.com/internal/common/tests/assert"
	"backend.brokedaear.com/internal/common/tests/test"
)

func TestConfigNewURIAddress(t *testing.T) {
	tests := []struct {
		test.CaseBase
		config Config
		want   string
	}{
		{
			CaseBase: test.NewCaseBase(
				"valid config creates URI address",
				"localhost:8080",
				false,
			),
			config: Config{
				Addr:      Address("localhost"),
				Port:      Port(8080),
				Env:       backend.EnvDevelopment,
				Version:   Version("1.0.0"),
				Telemetry: true,
			},
			want: "localhost:8080",
		},
		{
			CaseBase: test.NewCaseBase(
				"valid config with IP address",
				"127.0.0.1:9090",
				false,
			),
			config: Config{
				Addr:      Address("127.0.0.1"),
				Port:      Port(9090),
				Env:       backend.EnvProduction,
				Version:   Version("2.1.0"),
				Telemetry: false,
			},
			want: "127.0.0.1:9090",
		},
		{
			CaseBase: test.NewCaseBase(
				"invalid address returns error",
				"",
				true,
			),
			config: Config{
				Addr:      Address(""),
				Port:      Port(8080),
				Env:       backend.EnvDevelopment,
				Version:   Version("1.0.0"),
				Telemetry: true,
			},
			want: "",
		},
		{
			CaseBase: test.NewCaseBase(
				"invalid port returns error",
				"",
				true,
			),
			config: Config{
				Addr:      Address("localhost"),
				Port:      Port(0),
				Env:       backend.EnvDevelopment,
				Version:   Version("1.0.0"),
				Telemetry: true,
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.Name, func(t *testing.T) {
				got, err := tt.config.newURIAddress()
				assert.ErrorOrNoError(t, err, tt.WantErr)

				if !tt.WantErr {
					assert.Equal(t, got, tt.want)
				}
			},
		)
	}
}
