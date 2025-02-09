package models

import (
	"fmt"
)

// Credential represents any type of credential, such as an email,
// name, session JWT, etc.
type Credential interface {
	Valid() error
	fmt.Stringer
}

// A transactional setting is one in which a customer can potentially purchase
// a product.

// Customer represents a customer in any possible transactional setting, such as
// being in a checkout scenario, looking at the cart, or even just being logged
// in.
type customer struct {
	email        Credential
	password     Credential
	sessionToken Credential
}

// NewCustomer takes an email, password, and optionally a sessionToken to create
// a new customer instance.
func NewCustomer(email, password Credential, sessionToken ...Credential) (*customer, error) {

	var err error
	err = email.Valid()
	if err != nil {
		return nil, invalidEmailError{email.String()}
	}

	err = password.Valid()
	if err != nil {
		return nil, invalidPasswordError{password.String()}
	}

	c := &customer{
		email:    email,
		password: password,
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

// Product represents a product that we sell. This could be an audio plugin,
// physical merchandise, or anything of that nature.
type Product struct {
	// Name is the name of the product, for example, "Microwave", which is
	// our very first audio plugin.
	Name string

	// PriceId is found on stripe.
	PriceId string

	// ProductId is found on stripe.
	ProductId string
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
