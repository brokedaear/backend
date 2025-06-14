// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package backend

import "errors"

// Environment specifies the application runtime environment, like "staging"
// or "production".
type Environment struct {
	v string
}

func (e Environment) Validate() error {
	if e.v == "" {
		return errors.New("environment is nil")
	}

	if e.v == "invalid" {
		return errors.New("invalid environment")
	}
	return nil
}

func (e Environment) Value() any {
	return e.v
}

func (e Environment) String() string {
	return e.v
}

//nolint:gochecknoglobals // These simulate enums.
var (
	EnvDevelopment = Environment{"development"}
	EnvStaging     = Environment{"staging"}
	EnvProduction  = Environment{"production"}
)

// EnvFromString returns an environment given a string.
func EnvFromString(s string) (Environment, error) {
	switch s {
	case "development":
		return EnvDevelopment, nil
	case "staging":
		return EnvStaging, nil
	case "production":
		return EnvProduction, nil
	default:
		return Environment{"invalid"}, errors.New("invalid environment")
	}
}
