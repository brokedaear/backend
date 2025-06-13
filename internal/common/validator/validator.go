// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package validator

// Verifiable is a type that can be validated against constraints defined in
// the Validate method. A Value method also exists on Verifiable to return
// the underlying content that had been verified.
type Verifiable interface {
	Validate() error
	Value() any
}

// Check iterates through a list of Verifiable types, calling their Validate
// methods to check if the type respects their constraints. It errors on the
// first type that returns an error.
func Check(types ...Verifiable) error {
	if len(types) == 0 {
		return ErrNoTypesProvided
	}

	for _, v := range types {
		err := v.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

type VerifiableError string

func (e VerifiableError) Error() string {
	return string(e)
}

const (
	ErrNoTypesProvided VerifiableError = "no list of types provided"
)
