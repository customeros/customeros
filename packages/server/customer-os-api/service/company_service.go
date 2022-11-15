package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type CompanyService interface {
	MergeCompanyToContact(ctx context.Context, contactId string, input *entity.CompanyPositionEntity) (*entity.CompanyPositionEntity, error)
	UpdateCompanyPosition(ctx context.Context, contactId, companyPositionId, jobTitle string) (*entity.CompanyPositionEntity, error)
	DeleteCompanyPositionFromContact(ctx context.Context, contactId, companyPositionId string) (bool, error)

	GetCompanyPositionsForContact(ctx context.Context, contactId string) (*entity.CompanyPositionEntities, error)

	FindCompaniesByNameLike(ctx context.Context, page, limit int, companyName string) (*utils.Pagination, error)

	getDriver() neo4j.Driver
}

type companyService struct {
	repository *repository.RepositoryContainer
}

func NewCompanyService(repository *repository.RepositoryContainer) CompanyService {
	return &companyService{
		repository: repository,
	}
}

func (s *companyService) getDriver() neo4j.Driver {
	return *s.repository.Drivers.Neo4jDriver
}

func (s *companyService) MergeCompanyToContact(ctx context.Context, contactId string, input *entity.CompanyPositionEntity) (*entity.CompanyPositionEntity, error) {
	var err error
	var companyDbNode *dbtype.Node
	var positionDbRelationship *dbtype.Relationship

	if len(input.Company.Id) == 0 {
		companyDbNode, positionDbRelationship, err = s.repository.CompanyRepository.LinkNewCompanyToContact(common.GetContext(ctx).Tenant, contactId, input.Company.Name, input.JobTitle)
	} else {
		companyDbNode, positionDbRelationship, err = s.repository.CompanyRepository.LinkExistingCompanyToContact(common.GetContext(ctx).Tenant, contactId, input.Company.Id, input.JobTitle)
	}
	if err != nil {
		return nil, err
	}
	companyPositionEntity := s.mapCompanyPositionDbRelationshipToEntity(positionDbRelationship)
	companyPositionEntity.Company = *s.mapCompanyDbNodeToEntity(companyDbNode)
	return companyPositionEntity, nil
}

func (s *companyService) UpdateCompanyPosition(ctx context.Context, contactId, companyPositionId, jobTitle string) (*entity.CompanyPositionEntity, error) {
	companyDbNode, positionDbRelationship, err := s.repository.CompanyRepository.UpdateCompanyPosition(common.GetContext(ctx).Tenant, contactId, companyPositionId, jobTitle)
	if err != nil {
		return nil, err
	}
	companyPositionEntity := s.mapCompanyPositionDbRelationshipToEntity(positionDbRelationship)
	companyPositionEntity.Company = *s.mapCompanyDbNodeToEntity(companyDbNode)
	return companyPositionEntity, nil
}

func (s *companyService) DeleteCompanyPositionFromContact(ctx context.Context, contactId, companyPositionId string) (bool, error) {
	err := s.repository.CompanyRepository.DeleteCompanyPosition(common.GetContext(ctx).Tenant, contactId, companyPositionId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *companyService) GetCompanyPositionsForContact(ctx context.Context, contactId string) (*entity.CompanyPositionEntities, error) {

	companiesAndPositionsDbNodes, err := s.repository.CompanyRepository.GetCompanyPositionsForContact(common.GetContext(ctx).Tenant, contactId)
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

	dbNodesWithTotalCount, err := s.repository.CompanyRepository.GetPaginatedCompaniesWithNameLike(common.GetContext(ctx).Tenant, companyName, paginatedResult.GetSkip(), paginatedResult.GetLimit())
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
