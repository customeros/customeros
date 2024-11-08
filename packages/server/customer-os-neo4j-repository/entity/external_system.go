package entity

import (
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

const PropertyExternalSystemStripePaymentMethodTypes = "stripePaymentMethodTypes"

type ExternalSystemEntity struct {
	DataLoaderKey
	ExternalSystemId neo4jenum.ExternalSystemId
	Name             string
	Relationship     struct {
		ExternalId     string
		SyncDate       *time.Time
		ExternalUrl    *string
		ExternalSource *string
		Primary        bool
	}
	Stripe struct {
		PaymentMethodTypes []string
	}
}

type ExternalSystemEntities []ExternalSystemEntity
