package domain

type Shop struct{}

type NewCustomerEmail string

func (n NewCustomerEmail) Valid() error {
	return nil
}

func (n NewCustomerEmail) String() string {
	return n.String()
}

type NewCustomerPassword string

func (n NewCustomerPassword) Valid() error {
	return nil
}

func (n NewCustomerPassword) String() string {
	return n.String()
}

type RegisteredCustomerEmail string

func (r RegisteredCustomerEmail) Valid() error {
	return nil
}

func (r RegisteredCustomerEmail) String() string {
	return r.String()
}

type RegisteredCustomerPassword string

func (r RegisteredCustomerPassword) Valid() error {
	return nil
}

func (r RegisteredCustomerPassword) String() string {
	return r.String()
}
