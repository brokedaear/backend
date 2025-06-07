// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package domain

type Shop struct{}

type NewCustomerEmail string

func (n NewCustomerEmail) Valid() error {
	return nil
}

func (n NewCustomerEmail) String() string {
	return string(n)
}

type NewCustomerPassword string

func (n NewCustomerPassword) Valid() error {
	return nil
}

func (n NewCustomerPassword) String() string {
	return string(n)
}

type RegisteredCustomerEmail string

func (r RegisteredCustomerEmail) Valid() error {
	return nil
}

func (r RegisteredCustomerEmail) String() string {
	return string(r)
}

type RegisteredCustomerPassword string

func (r RegisteredCustomerPassword) Valid() error {
	return nil
}

func (r RegisteredCustomerPassword) String() string {
	return string(r)
}
