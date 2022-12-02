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

type CompanyService interface {
	MergeCompanyToContact(ctx context.Context, contactId string, input *entity.CompanyPositionEntity) (*entity.CompanyPositionEntity, error)
	UpdateCompanyPosition(ctx context.Context, contactId, companyPositionId string, input *entity.CompanyPositionEntity) (*entity.CompanyPositionEntity, error)
	DeleteCompanyPositionFromContact(ctx context.Context, contactId, companyPositionId string) (bool, error)
	GetCompanyPositionsForContact(ctx context.Context, contactId string) (*entity.CompanyPositionEntities, error)
	FindCompaniesByNameLike(ctx context.Context, page, limit int, companyName string) (*utils.Pagination, error)
}

type companyService struct {
	repositories *repository.Repositories
}

func NewCompanyService(repositories *repository.Repositories) CompanyService {
	return &companyService{
		repositories: repositories,
	}
}

func (s *companyService) MergeCompanyToContact(ctx context.Context, contactId string, input *entity.CompanyPositionEntity) (*entity.CompanyPositionEntity, error) {
	var err error
	var companyDbNode *dbtype.Node
	var positionDbRelationship *dbtype.Relationship

	session := utils.NewNeo4jWriteSession(*s.repositories.Drivers.Neo4jDriver)
	defer session.Close()

	newPosition, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		if len(input.Company.Id) == 0 {
			companyDbNode, positionDbRelationship, err = s.repositories.CompanyRepository.LinkNewCompanyToContactInTx(tx, common.GetContext(ctx).Tenant, contactId, input.Company.Name, input.JobTitle)
		} else {
			companyDbNode, positionDbRelationship, err = s.repositories.CompanyRepository.LinkExistingCompanyToContactInTx(tx, common.GetContext(ctx).Tenant, contactId, input.Company.Id, input.JobTitle)
		}
		if err != nil {
			return nil, err
		}
		companyPositionEntity := s.mapCompanyPositionDbRelationshipToEntity(positionDbRelationship)
		companyPositionEntity.Company = *s.mapCompanyDbNodeToEntity(companyDbNode)
		return companyPositionEntity, nil
	})

	return newPosition.(*entity.CompanyPositionEntity), nil
}

func (s *companyService) UpdateCompanyPosition(ctx context.Context, contactId, companyPositionId string, input *entity.CompanyPositionEntity) (*entity.CompanyPositionEntity, error) {
	var err error
	var companyDbNode *dbtype.Node
	var positionDbRelationship *dbtype.Relationship
	tenant := common.GetContext(ctx).Tenant

	session := utils.NewNeo4jWriteSession(*s.repositories.Drivers.Neo4jDriver)
	defer session.Close()

	currentPositionDtls, err := s.repositories.CompanyRepository.GetCompanyPositionForContact(session, tenant, contactId, companyPositionId)
	if err != nil {
		return nil, err
	}
	currentPositionId := utils.GetStringPropOrEmpty(utils.GetPropsFromRelationship(*currentPositionDtls.Position), "id")

	updatedPosition, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		if len(input.Company.Id) == 0 {
			err := s.repositories.CompanyRepository.DeleteCompanyPositionInTx(tx, tenant, contactId, currentPositionId)
			if err != nil {
				return nil, err
			}
			companyDbNode, positionDbRelationship, err = s.repositories.CompanyRepository.LinkNewCompanyToContactInTx(tx, tenant, contactId, input.Company.Name, input.JobTitle)
		} else if input.Company.Id == currentPositionId {
			companyDbNode, positionDbRelationship, err = s.repositories.CompanyRepository.UpdateCompanyPositionInTx(tx, common.GetContext(ctx).Tenant, contactId, companyPositionId, input.JobTitle)
		} else {
			err := s.repositories.CompanyRepository.DeleteCompanyPositionInTx(tx, tenant, contactId, currentPositionId)
			if err != nil {
				return nil, err
			}
			companyDbNode, positionDbRelationship, err = s.repositories.CompanyRepository.LinkExistingCompanyToContactInTx(tx, tenant, contactId, input.Company.Id, input.JobTitle)
		}
		companyPositionEntity := s.mapCompanyPositionDbRelationshipToEntity(positionDbRelationship)
		companyPositionEntity.Company = *s.mapCompanyDbNodeToEntity(companyDbNode)
		return companyPositionEntity, nil
	})

	return updatedPosition.(*entity.CompanyPositionEntity), err
}

func (s *companyService) DeleteCompanyPositionFromContact(ctx context.Context, contactId, companyPositionId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(*s.repositories.Drivers.Neo4jDriver)
	defer session.Close()
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		return nil, s.repositories.CompanyRepository.DeleteCompanyPositionInTx(tx, common.GetContext(ctx).Tenant, contactId, companyPositionId)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *companyService) GetCompanyPositionsForContact(ctx context.Context, contactId string) (*entity.CompanyPositionEntities, error) {
	session := utils.NewNeo4jReadSession(*s.repositories.Drivers.Neo4jDriver)
	defer session.Close()

	companiesAndPositionsDbNodes, err := s.repositories.CompanyRepository.GetCompanyPositionsForContact(session, common.GetContext(ctx).Tenant, contactId)
	if err != nil {
		return nil, err
	}

	companyPositionEntities := entity.CompanyPositionEntities{}
	for _, v := range companiesAndPositionsDbNodes {
		companyPositionEntity := s.mapCompanyPositionDbRelationshipToEntity(v.Position)
		companyPositionEntity.Company = *s.mapCompanyDbNodeToEntity(v.Company)
		companyPositionEntities = append(companyPositionEntities, *companyPositionEntity)
	}
	return &companyPositionEntities, nil
}

func (s *companyService) FindCompaniesByNameLike(ctx context.Context, page, limit int, companyName string) (*utils.Pagination, error) {
	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}

	dbNodesWithTotalCount, err := s.repositories.CompanyRepository.GetPaginatedCompaniesWithNameLike(common.GetContext(ctx).Tenant, companyName, paginatedResult.GetSkip(), paginatedResult.GetLimit())
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodesWithTotalCount.Count)

	companyEntities := entity.CompanyEntities{}

	for _, v := range dbNodesWithTotalCount.Nodes {
		companyEntities = append(companyEntities, *s.mapCompanyDbNodeToEntity(v))
	}
	paginatedResult.SetRows(&companyEntities)
	return &paginatedResult, nil
}

func (s *companyService) mapCompanyDbNodeToEntity(node *dbtype.Node) *entity.CompanyEntity {
	props := utils.GetPropsFromNode(*node)
	companyEntity := new(entity.CompanyEntity)
	companyEntity.Id = utils.GetStringPropOrEmpty(props, "id")
	companyEntity.Name = utils.GetStringPropOrEmpty(props, "name")
	return companyEntity
}

func (s *companyService) mapCompanyPositionDbRelationshipToEntity(relationship *dbtype.Relationship) *entity.CompanyPositionEntity {
	props := utils.GetPropsFromRelationship(*relationship)
	companyPositionEntity := new(entity.CompanyPositionEntity)
	companyPositionEntity.Id = utils.GetStringPropOrEmpty(props, "id")
	companyPositionEntity.JobTitle = utils.GetStringPropOrEmpty(props, "jobTitle")
	return companyPositionEntity
}
