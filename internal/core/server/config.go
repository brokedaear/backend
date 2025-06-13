// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"strconv"
	"strings"

	"backend.brokedaear.com/internal/common/validator"
)

// Config defines a default server configuration.
type Config struct {
	// Addr is the Address on which to bind the application.
	Addr Address

	// Port number to bind to for the application.
	Port Port

	// Env is the runtime environment, such as "development" or "production"
	Env Environment

	// Version is the Version of the software.
	Version Version

	// Telemetry determines if telemetry tracking for internals are enabled.
	Telemetry bool
}

func (c Config) Validate() error {
	return validator.Check(c.Addr, c.Port, c.Env, c.Version)
}

func (c Config) Value() any {
	return c
}

func NewConfig(addr string, port uint16, env uint8, version string) (*Config, error) {
	a := Address(addr)
	p := Port(port)
	e := Environment(env)
	v := Version(version)

	err := validator.Check(a, p, e, v)
	if err != nil {
		return nil, err
	}

	return &Config{
		Addr:    a,
		Port:    p,
		Env:     e,
		Version: v,
	}, nil
}

// Port represents a layer 4 OSI Port.
type Port uint16

func (p Port) String() string {
	return strconv.Itoa(int(p))
}

func (p Port) Validate() error {
	if p < 1024 || p >= 65534 {
		return ErrInvalidPortRange
	}

	return nil
}

func (p Port) Value() any {
	return uint16(p)
}

// Address represents a layer 4 OSI Address. An address must only be an
// IP address OR a domain name followed by a TLD.
type Address string

func (a Address) String() string {
	return string(a)
}

func (a Address) Validate() error {
	var (
		colon        = ":"
		space        = " "
		forwardSlash = "/"
		addr         = a.String()
	)

	if len(addr) == 0 {
		return ErrInvalidAddressLength
	}

	if strings.Contains(addr, colon) {
		return ErrInvalidAddressColon
	}

	if strings.Contains(addr, space) {
		return ErrInvalidAddressSpace
	}

	if strings.Contains(addr, forwardSlash) {
		return ErrInvalidAddressWithPath
	}

	return nil
}

func (a Address) Value() any {
	return a.String()
}

// Environment specifies the application runtime Environment.
type Environment uint8

const (
	EnvDevelopment Environment = iota
	EnvStaging
	EnvProduction
	EnvCI
)

var environments = [...]string{"DEVELOPMENT", "STAGING", "PRODUCTION", "CI", "INVALID"}

func (e Environment) String() string {
	return environments[e]
}

func (e Environment) Validate() error {
	totalEnvs := Environment(len(environments))
	if e > totalEnvs {
		return ErrInvalidEnvironment
	}

	return nil
}

func (e Environment) Value() any {
	return uint8(e)
}

type Version string

func (v Version) String() string {
	return string(v)
}

func (v Version) Validate() error {
	elements := strings.Split(v.String(), ".")
	if len(elements) != 3 {
		return ErrInvalidVersionFormat
	}

	for _, element := range elements {
		n, err := strconv.Atoi(element)
		if err != nil {
			return ErrInvalidVersionChars
		}

		if n < 0 {
			return ErrInvalidVersionSign
		}
	}

	return nil
}

func (v Version) Value() any {
	return v.String()
}

type ConfigError string

func (c ConfigError) Error() string {
	return string(c)
}

const (
	ErrInvalidEnvironment     ConfigError = "Configured Environment invalid"
	ErrInvalidPortRange       ConfigError = "Configured Port range must be [1024, 65535)"
	ErrInvalidVersionFormat   ConfigError = "Configured Version must be of the format x.x.x"
	ErrInvalidVersionChars    ConfigError = "Configured Version must only be an unsigned integer"
	ErrInvalidVersionSign     ConfigError = "Configured Version must be >= 0"
	ErrInvalidAddressLength   ConfigError = "Configured Address length must be greater than 0"
	ErrInvalidAddressColon    ConfigError = "Configured Address must not contain a colon"
	ErrInvalidAddressSpace    ConfigError = "Configured Address must not contain a space"
	ErrInvalidAddressWithPath ConfigError = "Configured Address must not contain a path"
	ErrInvalidVersionAlpha    ConfigError = "Configured Version cannot contain alpha chars"
)
