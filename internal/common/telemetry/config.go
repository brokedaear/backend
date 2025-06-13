// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package telemetry

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Config is the configuration for telemetry.
type Config struct {
	ServiceName    string
	ServiceVersion string
	ServiceId      string
	ExporterConfig ExporterConfig
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.ServiceName) == "" {
		return ErrNoServiceName
	}

	if strings.TrimSpace(c.ServiceVersion) == "" {
		return ErrNoServiceVersion
	}

	if strings.TrimSpace(c.ServiceId) == "" {
		return ErrNoServiceId
	}

	err := validateServiceName(c.ServiceName)
	if err != nil {
		return errors.Wrap(err, ErrInvalidServiceName.Error())
	}

	err = validateServiceVersion(c.ServiceVersion)
	if err != nil {
		return errors.Wrap(err, ErrInvalidServiceVersion.Error())
	}

	err = c.ExporterConfig.Validate()
	if err != nil {
		return errors.Wrap(err, ErrInvalidExporterConfig.Error())
	}

	return nil
}

func (c Config) Value() any {
	return c
}

// ExporterConfig holds configuration for an OTEL exporter.
type ExporterConfig struct {
	Type     ExporterType
	Endpoint string
	Insecure bool
	Headers  map[string]string
}

func (e ExporterConfig) Validate() error {
	err := e.Type.Validate()
	if err != nil {
		return err
	}

	// Validate endpoint based on exporter type
	switch e.Type {
	case ExporterTypeGRPC, ExporterTypeHTTP:
		err := e.validateNetworkEndpoint()
		if err != nil {
			return err
		}
	case ExporterTypeStdout:
		err := e.validateStdoutEndpoint()
		if err != nil {
			return err
		}
	}

	err = e.validateHeaders()
	if err != nil {
		return err
	}

	return nil
}

func (e ExporterConfig) validateNetworkEndpoint() error {
	if strings.TrimSpace(e.Endpoint) == "" {
		return ErrEndpointRequired
	}

	parsedURL, err := url.Parse(e.Endpoint)
	if err != nil {
		return ErrInvalidEndpointURL
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return ErrInvalidEndpointScheme
	}

	if parsedURL.Host == "" {
		return ErrEndpointMissingHost
	}

	// Warn about insecure HTTPS endpoints
	// if parsedURL.Scheme == "https" && e.Insecure {
	// This is a warning case, not an error - log this in real implementation
	// For validation purposes, we'll allow it but could add a warning mechanism
	// }

	return nil
}

func (e ExporterConfig) validateStdoutEndpoint() error {
	// For stdout exporter, endpoint is optional (represents file path)
	if e.Endpoint != "" {
		if strings.TrimSpace(e.Endpoint) == "" {
			return ErrInvalidFilePath
		}

		invalidChars := []string{"\x00", "\n", "\r"}
		for _, char := range invalidChars {
			if strings.Contains(e.Endpoint, char) {
				return ErrFilePathInvalidChar
			}
		}
	}

	return nil
}

func (e ExporterConfig) validateHeaders() error {
	for key, value := range e.Headers {
		if strings.TrimSpace(key) == "" {
			return ErrHeaderKeyEmpty
		}

		if strings.ContainsAny(key, " \t\n\r") {
			return ConfigError(fmt.Sprintf("header key '%s' contains invalid characters", key))
		}

		if strings.ContainsAny(value, "\n\r") {
			return ConfigError(
				fmt.Sprintf(
					"header value for key '%s' contains invalid characters",
					key,
				),
			)
		}
	}

	return nil
}

func (e ExporterConfig) Value() any {
	return e
}

// ExporterType defines the type of OTLP exporter to use.
type ExporterType uint8

func (e ExporterType) Validate() error {
	if e > 2 {
		return ErrInvalidExporterType
	}

	return nil
}

func (e ExporterType) Value() any {
	return e
}

func (e ExporterType) String() string {
	switch e {
	case ExporterTypeGRPC:
		return "grpc"
	case ExporterTypeHTTP:
		return "http"
	case ExporterTypeStdout:
		return "stdout"
	default:
		return "INVALID"
	}
}

const (
	ExporterTypeGRPC ExporterType = iota
	ExporterTypeHTTP
	ExporterTypeStdout
)

const serviceNameLimit = 255
const serviceNameMinimum = 1

func validateServiceName(name string) error {
	name = strings.TrimSpace(name)

	// Check minimum length
	if len(name) < serviceNameMinimum {
		return ErrServiceNameEmpty
	}

	// Check maximum length (reasonable limit)
	if len(name) > serviceNameLimit {
		return ErrServiceNameTooLong
	}

	// OpenTelemetry recommends using dot notation for service names
	// Allow alphanumeric, dots, hyphens, and underscores
	for _, char := range name {
		if !isValidServiceNameChar(char) {
			return errors.Wrapf(ErrServiceNameInvalidChar, "char %s", string(char))
		}
	}

	// Service name should not start or end with dot
	if strings.HasPrefix(name, ".") || strings.HasSuffix(name, ".") {
		return ErrServiceNameStartEndDot
	}

	// Should not have consecutive dots
	if strings.Contains(name, "..") {
		return ErrServiceNameConsecutiveDots
	}

	return nil
}

const serviceVersionLimit = 128

func validateServiceVersion(version string) error {
	version = strings.TrimSpace(version)

	// Check minimum length
	if len(version) < 1 {
		return ErrServiceVersionEmpty
	}

	// Check maximum length
	if len(version) > serviceVersionLimit {
		return ErrServiceVersionTooLong
	}

	// Basic semantic version pattern check (flexible)
	// Allow versions like: 1.0.0, v1.0.0, 1.0.0-beta, 1.0.0+build.1, etc.
	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-+_"
	for _, char := range version {
		if !strings.ContainsRune(validChars, char) {
			return errors.Wrapf(ErrServiceVersionInvalidChar, "char %s", string(char))
		}
	}

	return nil
}

func isValidServiceNameChar(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9') ||
		char == '.' ||
		char == '-' ||
		char == '_'
}

type ConfigError string

func (e ConfigError) Error() string {
	return string(e)
}

var (
	ErrNoServiceName              ConfigError = "no service name provided"
	ErrNoServiceId                ConfigError = "no service id provided"
	ErrNoServiceVersion           ConfigError = "no service version provided"
	ErrInvalidServiceName         ConfigError = "invalid service name"
	ErrInvalidServiceVersion      ConfigError = "invalid service version"
	ErrInvalidExporterConfig      ConfigError = "invalid exporter config"
	ErrServiceNameConsecutiveDots ConfigError = "service name contains consecutive dots"
	ErrServiceNameStartEndDot     ConfigError = "service name starts or ends with dot"
	ErrServiceNameInvalidChar     ConfigError = "service name contains invalid character"
	ErrServiceNameTooLong                     = ConfigError("service name chars greater than " + strconv.Itoa(serviceNameLimit))
	ErrServiceNameEmpty           ConfigError = "service name is empty"
	ErrServiceVersionEmpty        ConfigError = "service version is empty"
	ErrServiceVersionTooLong                  = ConfigError("service version chars greater than " + strconv.Itoa(serviceVersionLimit))
	ErrServiceVersionInvalidChar  ConfigError = "service version contains invalid character"
	ErrEndpointRequired           ConfigError = "endpoint is required for exporter"
	ErrInvalidEndpointURL         ConfigError = "invalid endpoint URL"
	ErrInvalidEndpointScheme      ConfigError = "endpoint must use http or https scheme"
	ErrEndpointMissingHost        ConfigError = "endpoint must include a host"
	ErrInvalidFilePath            ConfigError = "if specified, endpoint must be a valid file path"
	ErrFilePathInvalidChar        ConfigError = "endpoint contains invalid character for file path"
	ErrHeaderKeyEmpty             ConfigError = "header key cannot be empty"
	ErrInvalidExporterType        ConfigError = "invalid exporter type"
)
