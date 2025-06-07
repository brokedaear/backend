// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package domain

import (
	"fmt"
)

type WebShop interface {
	SignUp() error
	Login() error
	Buy() error
}

// Credential represents any type of credential, such as an email,
// name, session JWT, etc.
type Credential interface {
	Valid() error
	String() string
}

// A transactional setting is one in which a customer can potentially purchase
// a product.

// user represents a user in the overall application system.
type user struct {
	email    Credential
	password Credential
}

// possibleCustomer represents a possible new customer based on a sign up query.
// To become a registeredCustomer, verified must be true, email must be valid,
// and password must be valid.
type possibleCustomer struct {
	user
	verified bool
}

func NewPossibleCustomer(email, password Credential, sessionToken ...Credential) (*possibleCustomer, error) {
	var err error

	err = email.Valid()
	if err != nil {
		return nil, invalidEmailError{email.String()}
	}

	err = password.Valid()
	if err != nil {
		return nil, invalidPasswordError{password.String()}
	}

	c := &possibleCustomer{
		user: user{
			email:    email,
			password: password,
		},
	}

	return c, nil
}

// registeredCustomer represents a customer in any possible transactional
// setting, such as being in a checkout scenario, looking at the cart, or even
// just being logged in. registeredCustomer is different from possibleCustomer,
// because of the password validation methods.
type registeredCustomer struct {
	user
	sessionToken Credential
}

// NewRegisteredCustomer takes an email, password, and optionally a
// sessionToken to create a new customer instance.
func NewRegisteredCustomer(email, password Credential, sessionToken ...Credential) (*registeredCustomer, error) {
	var err error

	err = email.Valid()
	if err != nil {
		return nil, invalidEmailError{email.String()}
	}

	err = password.Valid()
	if err != nil {
		return nil, invalidPasswordError{password.String()}
	}

	c := &registeredCustomer{
		user: user{
			email:    email,
			password: password,
		},
	}

	if len(sessionToken) >= 1 {
		err = sessionToken[0].Valid()
		if err != nil {
			return nil, invalidSessionTokenError{sessionToken[0].String()}
		}

		c.sessionToken = sessionToken[0]
	}

	return c, nil
}

type invalidEmailError struct {
	email string
}

func (e invalidEmailError) Error() string {
	return fmt.Sprintf("invalid email %s", e.email)
}

type invalidPasswordError struct {
	password string
}

func (e invalidPasswordError) Error() string {
	return fmt.Sprintf("invalid password %s", e.password)
}

type invalidSessionTokenError struct {
	token string
}

func (e invalidSessionTokenError) Error() string {
	return fmt.Sprintf("invalid session token %s", e.token)
}
