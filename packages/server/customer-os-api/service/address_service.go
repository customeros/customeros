package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type AddressService interface {
	FindAllForContact(ctx context.Context, contactId string) (*entity.AddressEntities, error)
	FindAllForCompany(ctx context.Context, companyId string) (*entity.AddressEntities, error)
}

type addressService struct {
	repositories *repository.Repositories
}

func NewAddressService(repositories *repository.Repositories) AddressService {
	return &addressService{
		repositories: repositories,
	}
}

func (s *addressService) getDriver() neo4j.Driver {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *addressService) FindAllForContact(ctx context.Context, contactId string) (*entity.AddressEntities, error) {
	session := utils.NewNeo4jReadSession(s.getDriver())
	defer session.Close()

	dbNodes, err := s.repositories.AddressRepository.FindAllForContact(session, common.GetContext(ctx).Tenant, contactId)
	if err != nil {
		return nil, err
	}

	addressEntities := entity.AddressEntities{}
	for _, dbNode := range dbNodes {
		addressEntities = append(addressEntities, *s.mapDbNodeToAddressEntity(dbNode))
	}
	return &addressEntities, nil
}

func (s *addressService) FindAllForCompany(ctx context.Context, companyId string) (*entity.AddressEntities, error) {
	session := utils.NewNeo4jReadSession(s.getDriver())
	defer session.Close()

	dbNodes, err := s.repositories.AddressRepository.FindAllForCompany(session, common.GetContext(ctx).Tenant, companyId)
	if err != nil {
		return nil, err
	}

	addressEntities := entity.AddressEntities{}
	for _, dbNode := range dbNodes {
		addressEntities = append(addressEntities, *s.mapDbNodeToAddressEntity(dbNode))
	}
	return &addressEntities, nil
}

func (s *addressService) mapDbNodeToAddressEntity(node *dbtype.Node) *entity.AddressEntity {
	props := utils.GetPropsFromNode(*node)
	result := entity.AddressEntity{
		Id:       utils.GetStringPropOrEmpty(props, "id"),
		Source:   utils.GetStringPropOrEmpty(props, "source"),
		Country:  utils.GetStringPropOrEmpty(props, "country"),
		State:    utils.GetStringPropOrEmpty(props, "state"),
		City:     utils.GetStringPropOrEmpty(props, "city"),
		Address:  utils.GetStringPropOrEmpty(props, "address"),
		Address2: utils.GetStringPropOrEmpty(props, "address2"),
		Zip:      utils.GetStringPropOrEmpty(props, "zip"),
		Phone:    utils.GetStringPropOrEmpty(props, "phone"),
		Fax:      utils.GetStringPropOrEmpty(props, "fax"),
	}
	return &result
}
