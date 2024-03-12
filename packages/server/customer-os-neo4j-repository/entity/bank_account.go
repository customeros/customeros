package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type BankAccountEntity struct {
	Id                  string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Source              DataSource
	SourceOfTruth       DataSource
	AppSource           string
	BankName            string
	BankTransferEnabled bool
	AllowInternational  bool
	Currency            enum.Currency
	Iban                string
	Bic                 string
	SortCode            string
	AccountNumber       string
	RoutingNumber       string
}

type BankAccountEntities []BankAccountEntity
