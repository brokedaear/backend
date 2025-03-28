// data.go contains models that correspond to the data access layer, that is,
// the raw data representations of the data in the database.

package domain

import "time"

type Customer struct {
	ID             int
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// Product represents a product that we sell. This could be an audio plugin,
// physical merchandise, or anything of that nature.
type Product struct {
	ID int

	// Type of the product, 0 is plugin, 1 is merchandise.
	Type int

	// Name is the name of the product, for example, "Microwave", which is
	// our very first audio plugin.
	Name string

	// PriceId is found on stripe.
	PriceId string

	// ProductId is found on stripe.
	ProductId string
}
