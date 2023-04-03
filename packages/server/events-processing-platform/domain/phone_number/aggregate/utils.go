package aggregate

import (
	es "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"strings"
)

// GetPhoneNumberAggregateID get phone_number aggregate id for eventstoredb
func GetPhoneNumberAggregateID(eventAggregateID string, tenant string) string {
	return strings.ReplaceAll(eventAggregateID, "phone_number-"+tenant+"-", "")
}

func IsAggregateNotFound(aggregate es.Aggregate) bool {
	return aggregate.GetVersion() == 0
}
