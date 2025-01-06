package entity

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
)

const PropertyExternalSystemStripePaymentMethodTypes = "stripePaymentMethodTypes"

type ExternalSystemEntity struct {
	DataLoaderKey
	ExternalSystemId enum.Source
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
