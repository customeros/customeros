package command

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type CreateContractCommand struct {
	eventstore.BaseCommand
	DataFields     model.ContractDataFields
	Source         commonmodel.Source
	ExternalSystem commonmodel.ExternalSystem
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}

// NewCreateContractCommand creates a new CreateContractCommand.
func NewCreateContractCommand(
	contractId, tenant, loggedInUserId string,
	dataFields model.ContractDataFields,
	source commonmodel.Source,
	externalSystem commonmodel.ExternalSystem,
	createdAt, updatedAt *time.Time) *CreateContractCommand {

	return &CreateContractCommand{
		BaseCommand:    eventstore.NewBaseCommand(contractId, tenant, loggedInUserId),
		DataFields:     dataFields,
		Source:         source,
		ExternalSystem: externalSystem,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

// UpdateContractCommand contains the data needed to update a contract.
type UpdateContractCommand struct {
	eventstore.BaseCommand
	DataFields     model.ContractDataFields
	ExternalSystem commonmodel.ExternalSystem
	UpdatedAt      *time.Time
	Source         commonmodel.Source
}

func NewUpdateContractCommand(
	contractId, tenant, loggedInUserId string,
	dataFields model.ContractDataFields,
	source commonmodel.Source,
	externalSystem commonmodel.ExternalSystem,
	updatedAt *time.Time) *UpdateContractCommand {

	return &UpdateContractCommand{
		BaseCommand:    eventstore.NewBaseCommand(contractId, tenant, loggedInUserId),
		DataFields:     dataFields,
		Source:         source,
		ExternalSystem: externalSystem,
		UpdatedAt:      updatedAt,
	}
}
