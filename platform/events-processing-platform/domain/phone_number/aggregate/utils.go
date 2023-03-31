package aggregate

import (
	"strings"

	"github.com/AleksK1NG/es-microservice/pkg/es"
)

// GetPhoneNumberAggregateID get phone_number aggregate id for eventstoredb
func GetPhoneNumberAggregateID(eventAggregateID string, tenant string) string {
	return strings.ReplaceAll(eventAggregateID, "phone_number-"+tenant+"-", "")
}

func IsAggregateNotFound(aggregate es.Aggregate) bool {
	return aggregate.GetVersion() == 0
}
