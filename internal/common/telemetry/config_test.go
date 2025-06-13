// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package telemetry

import (
	"strings"
	"testing"

	"backend.brokedaear.com/internal/common/tests/assert"
	"backend.brokedaear.com/internal/common/tests/test"
)

type ConfigTestCase struct {
	test.CaseBase
	Config Config
}

func TestConfig_Validate(t *testing.T) {
	tests := []ConfigTestCase{
		{
			CaseBase: test.NewCaseBase("valid config", nil, false),
			Config: Config{
				ServiceName:    "my.service",
				ServiceVersion: "1.0.0",
				ServiceId:      "service-123",
				ExporterConfig: ExporterConfig{
					Type:     ExporterTypeGRPC,
					Endpoint: "http://localhost:4317",
					Insecure: true,
				},
			},
		},
		{
			CaseBase: test.NewCaseBase("empty service name", "service name is required", true),
			Config: Config{
				ServiceName:    "",
				ServiceVersion: "1.0.0",
				ServiceId:      "service-123",
				ExporterConfig: ExporterConfig{Type: ExporterTypeStdout},
			},
		},
		{
			CaseBase: test.NewCaseBase(
				"whitespace only service name",
				"service name is required",
				true,
			),
			Config: Config{
				ServiceName:    "   ",
				ServiceVersion: "1.0.0",
				ServiceId:      "service-123",
				ExporterConfig: ExporterConfig{Type: ExporterTypeStdout},
			},
		},
		{
			CaseBase: test.NewCaseBase(
				"empty service version",
				"service version is required",
				true,
			),
			Config: Config{
				ServiceName:    "my.service",
				ServiceVersion: "",
				ServiceId:      "service-123",
				ExporterConfig: ExporterConfig{Type: ExporterTypeStdout},
			},
		},
		{
			CaseBase: test.NewCaseBase("empty service id", "service ID is required", true),
			Config: Config{
				ServiceName:    "my.service",
				ServiceVersion: "1.0.0",
				ServiceId:      "",
				ExporterConfig: ExporterConfig{Type: ExporterTypeStdout},
			},
		},
		{
			CaseBase: test.NewCaseBase("invalid service name format", "invalid service name", true),
			Config: Config{
				ServiceName:    "my service!",
				ServiceVersion: "1.0.0",
				ServiceId:      "service-123",
				ExporterConfig: ExporterConfig{Type: ExporterTypeStdout},
			},
		},
		{
			CaseBase: test.NewCaseBase("invalid exporter config", "invalid exporter config", true),
			Config: Config{
				ServiceName:    "my.service",
				ServiceVersion: "1.0.0",
				ServiceId:      "service-123",
				ExporterConfig: ExporterConfig{
					Type:     ExporterTypeGRPC,
					Endpoint: "", // missing required endpoint
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.Name, func(t *testing.T) {
				got := tt.Config.Validate()
				assert.ErrorOrNoError(t, got, tt.WantErr)
			},
		)
	}
}

type ExporterConfigTestCase struct {
	test.CaseBase
	Config ExporterConfig
}

func TestExporterConfig_Validate(t *testing.T) {
	tests := []ExporterConfigTestCase{
		{
			CaseBase: test.NewCaseBase("valid grpc config", nil, false),
			Config: ExporterConfig{
				Type:     ExporterTypeGRPC,
				Endpoint: "http://localhost:4317",
				Insecure: true,
				Headers:  map[string]string{"api-key": "secret"},
			},
		},
		{
			CaseBase: test.NewCaseBase("valid http config", nil, false),
			Config: ExporterConfig{
				Type:     ExporterTypeHTTP,
				Endpoint: "https://api.example.com/v1/traces",
				Insecure: false,
			},
		},
		{
			CaseBase: test.NewCaseBase("valid stdout config", nil, false),
			Config: ExporterConfig{
				Type: ExporterTypeStdout,
			},
		},
		{
			CaseBase: test.NewCaseBase("valid stdout config with file", nil, false),
			Config: ExporterConfig{
				Type:     ExporterTypeStdout,
				Endpoint: "/tmp/otel-output.json",
			},
		},
		{
			CaseBase: test.NewCaseBase("invalid exporter type", "invalid exporter type", true),
			Config: ExporterConfig{
				Type: ExporterType(99),
			},
		},
		{
			CaseBase: test.NewCaseBase(
				"grpc missing endpoint",
				"endpoint is required for grpc exporter",
				true,
			),
			Config: ExporterConfig{
				Type: ExporterTypeGRPC,
			},
		},
		{
			CaseBase: test.NewCaseBase(
				"http missing endpoint",
				"endpoint is required for http exporter",
				true,
			),
			Config: ExporterConfig{
				Type: ExporterTypeHTTP,
			},
		},
		{
			CaseBase: test.NewCaseBase("invalid endpoint url", "invalid endpoint URL", true),
			Config: ExporterConfig{
				Type:     ExporterTypeGRPC,
				Endpoint: "not-a-url",
			},
		},
		{
			CaseBase: test.NewCaseBase(
				"invalid endpoint scheme",
				"endpoint must use http or https scheme",
				true,
			),
			Config: ExporterConfig{
				Type:     ExporterTypeGRPC,
				Endpoint: "ftp://localhost:4317",
			},
		},
		{
			CaseBase: test.NewCaseBase(
				"endpoint missing host",
				"endpoint must include a host",
				true,
			),
			Config: ExporterConfig{
				Type:     ExporterTypeGRPC,
				Endpoint: "http://",
			},
		},
		{
			CaseBase: test.NewCaseBase(
				"stdout with invalid file path",
				"endpoint contains invalid character for file path",
				true,
			),
			Config: ExporterConfig{
				Type:     ExporterTypeStdout,
				Endpoint: "/path/with\x00null",
			},
		},
		{
			CaseBase: test.NewCaseBase("invalid header key", "header key cannot be empty", true),
			Config: ExporterConfig{
				Type:     ExporterTypeGRPC,
				Endpoint: "http://localhost:4317",
				Headers:  map[string]string{"": "value"},
			},
		},
		{
			CaseBase: test.NewCaseBase(
				"header key with spaces",
				"header key 'api key' contains invalid characters",
				true,
			),
			Config: ExporterConfig{
				Type:     ExporterTypeGRPC,
				Endpoint: "http://localhost:4317",
				Headers:  map[string]string{"api key": "value"},
			},
		},
		{
			CaseBase: test.NewCaseBase(
				"header value with newline",
				"header value for key 'api-key' contains invalid characters",
				true,
			),
			Config: ExporterConfig{
				Type:     ExporterTypeGRPC,
				Endpoint: "http://localhost:4317",
				Headers:  map[string]string{"api-key": "value\nwith\nnewlines"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.Name, func(t *testing.T) {
				got := tt.Config.Validate()
				assert.ErrorOrNoError(t, got, tt.WantErr)
			},
		)
	}
}

type ExporterTypeTestCase struct {
	test.CaseBase
	ExporterType ExporterType
}

func TestExporterType_Validate(t *testing.T) {
	tests := []ExporterTypeTestCase{
		{
			CaseBase:     test.NewCaseBase("valid grpc type", nil, false),
			ExporterType: ExporterTypeGRPC,
		},
		{
			CaseBase:     test.NewCaseBase("valid http type", nil, false),
			ExporterType: ExporterTypeHTTP,
		},
		{
			CaseBase:     test.NewCaseBase("valid stdout type", nil, false),
			ExporterType: ExporterTypeStdout,
		},
		{
			CaseBase: test.NewCaseBase(
				"invalid type - too high",
				"invalid exporter type: 3",
				true,
			),
			ExporterType: ExporterType(3),
		},
		{
			CaseBase: test.NewCaseBase(
				"invalid type - way too high",
				"invalid exporter type: 99",
				true,
			),
			ExporterType: ExporterType(99),
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.Name, func(t *testing.T) {
				got := tt.ExporterType.Validate()
				assert.ErrorOrNoError(t, got, tt.WantErr)
			},
		)
	}
}

type ExporterTypeStringTestCase struct {
	test.CaseBase
	ExporterType ExporterType
}

func TestExporterType_String(t *testing.T) {
	tests := []ExporterTypeStringTestCase{
		{
			CaseBase:     test.NewCaseBase("grpc type", "grpc", false),
			ExporterType: ExporterTypeGRPC,
		},
		{
			CaseBase:     test.NewCaseBase("http type", "http", false),
			ExporterType: ExporterTypeHTTP,
		},
		{
			CaseBase:     test.NewCaseBase("stdout type", "stdout", false),
			ExporterType: ExporterTypeStdout,
		},
		{
			CaseBase:     test.NewCaseBase("unknown type", "INVALID", false),
			ExporterType: ExporterType(99),
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.Name, func(t *testing.T) {
				got := tt.ExporterType.String()
				assert.Equal(t, got, tt.Want.(string))
			},
		)
	}
}

type ServiceNameTestCase struct {
	test.CaseBase
	ServiceName string
}

func TestValidateServiceName(t *testing.T) {
	tests := []ServiceNameTestCase{
		{
			CaseBase:    test.NewCaseBase("valid simple name", nil, false),
			ServiceName: "myservice",
		},
		{
			CaseBase:    test.NewCaseBase("valid dotted name", nil, false),
			ServiceName: "my.service.name",
		},
		{
			CaseBase:    test.NewCaseBase("valid with hyphens and underscores", nil, false),
			ServiceName: "my-service_name",
		},
		{
			CaseBase:    test.NewCaseBase("valid with numbers", nil, false),
			ServiceName: "service123",
		},
		{
			CaseBase:    test.NewCaseBase("empty name", "service name cannot be empty", true),
			ServiceName: "",
		},
		{
			CaseBase: test.NewCaseBase(
				"name starting with dot",
				"service name cannot start or end with a dot",
				true,
			),
			ServiceName: ".service",
		},
		{
			CaseBase: test.NewCaseBase(
				"name ending with dot",
				"service name cannot start or end with a dot",
				true,
			),
			ServiceName: "service.",
		},
		{
			CaseBase: test.NewCaseBase(
				"consecutive dots",
				"service name cannot contain consecutive dots",
				true,
			),
			ServiceName: "my..service",
		},
		{
			CaseBase: test.NewCaseBase(
				"invalid characters",
				"service name contains invalid character",
				true,
			),
			ServiceName: "my service!",
		},
		{
			CaseBase:    test.NewCaseBase("too long name", "service name too long", true),
			ServiceName: strings.Repeat("a", 256),
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.Name, func(t *testing.T) {
				got := validateServiceName(tt.ServiceName)
				assert.ErrorOrNoError(t, got, tt.WantErr)
			},
		)
	}
}

type ServiceVersionTestCase struct {
	test.CaseBase
	ServiceVersion string
}

func TestValidateServiceVersion(t *testing.T) {
	tests := []ServiceVersionTestCase{
		{
			CaseBase:       test.NewCaseBase("valid semver", nil, false),
			ServiceVersion: "1.0.0",
		},
		{
			CaseBase:       test.NewCaseBase("valid semver with v prefix", nil, false),
			ServiceVersion: "v1.0.0",
		},
		{
			CaseBase:       test.NewCaseBase("valid with pre-release", nil, false),
			ServiceVersion: "1.0.0-alpha.1",
		},
		{
			CaseBase:       test.NewCaseBase("valid with build metadata", nil, false),
			ServiceVersion: "1.0.0+build.123",
		},
		{
			CaseBase:       test.NewCaseBase("valid complex version", nil, false),
			ServiceVersion: "v2.1.0-beta.1+build.456",
		},
		{
			CaseBase: test.NewCaseBase(
				"empty version",
				"service version cannot be empty",
				true,
			),
			ServiceVersion: "",
		},
		{
			CaseBase: test.NewCaseBase(
				"version with invalid characters",
				"service version contains invalid character",
				true,
			),
			ServiceVersion: "1.0.0@invalid",
		},
		{
			CaseBase:       test.NewCaseBase("too long version", "service version too long", true),
			ServiceVersion: strings.Repeat("a", 129),
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.Name, func(t *testing.T) {
				got := validateServiceVersion(tt.ServiceVersion)
				assert.ErrorOrNoError(t, got, tt.WantErr)
			},
		)
	}
}

type ServiceNameCharTestCase struct {
	test.CaseBase
	Char rune
}

func TestIsValidServiceNameChar(t *testing.T) {
	tests := []ServiceNameCharTestCase{
		{
			CaseBase: test.NewCaseBase("lowercase letter", true, false),
			Char:     'a',
		},
		{
			CaseBase: test.NewCaseBase("uppercase letter", true, false),
			Char:     'Z',
		},
		{
			CaseBase: test.NewCaseBase("digit", true, false),
			Char:     '5',
		},
		{
			CaseBase: test.NewCaseBase("dot", true, false),
			Char:     '.',
		},
		{
			CaseBase: test.NewCaseBase("hyphen", true, false),
			Char:     '-',
		},
		{
			CaseBase: test.NewCaseBase("underscore", true, false),
			Char:     '_',
		},
		{
			CaseBase: test.NewCaseBase("space", false, false),
			Char:     ' ',
		},
		{
			CaseBase: test.NewCaseBase("exclamation", false, false),
			Char:     '!',
		},
		{
			CaseBase: test.NewCaseBase("at symbol", false, false),
			Char:     '@',
		},
		{
			CaseBase: test.NewCaseBase("slash", false, false),
			Char:     '/',
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.Name, func(t *testing.T) {
				got := isValidServiceNameChar(tt.Char)
				assert.Equal(t, got, tt.Want.(bool))
			},
		)
	}
}

func BenchmarkConfig_Validate(b *testing.B) {
	config := Config{
		ServiceName:    "my.service",
		ServiceVersion: "1.0.0",
		ServiceId:      "service-123",
		ExporterConfig: ExporterConfig{
			Type:     ExporterTypeGRPC,
			Endpoint: "http://localhost:4317",
			Insecure: true,
			Headers:  map[string]string{"api-key": "secret"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.Validate()
	}
}

func BenchmarkExporterConfig_Validate(b *testing.B) {
	config := ExporterConfig{
		Type:     ExporterTypeGRPC,
		Endpoint: "http://localhost:4317",
		Insecure: true,
		Headers:  map[string]string{"api-key": "secret"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.Validate()
	}
}
