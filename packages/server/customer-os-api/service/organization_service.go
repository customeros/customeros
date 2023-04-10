package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	events_processing_organization "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/proto/organization"
	"github.com/sirupsen/logrus"
	"reflect"
)

type OrganizationService interface {
	Create(ctx context.Context, input *OrganizationCreateData) (*entity.OrganizationEntity, error)
	Update(ctx context.Context, input *OrganizationUpdateData) (*entity.OrganizationEntity, error)
	GetOrganizationForJobRole(ctx context.Context, roleId string) (*entity.OrganizationEntity, error)
	GetOrganizationById(ctx context.Context, organizationId string) (*entity.OrganizationEntity, error)
	FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	GetOrganizationsForContact(ctx context.Context, contactId string, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	PermanentDelete(ctx context.Context, organizationId string) (bool, error)
	Merge(ctx context.Context, primaryOrganizationId, mergedOrganizationId string) error
	GetOrganizationsForEmails(ctx context.Context, emailIds []string) (*entity.OrganizationEntities, error)
	GetOrganizationsForPhoneNumbers(ctx context.Context, phoneNumberIds []string) (*entity.OrganizationEntities, error)
	GetSubsidiaries(ctx context.Context, parentOrganizationId string) (*entity.OrganizationEntities, error)
	GetSubsidiaryOf(ctx context.Context, organizationId string) (*entity.OrganizationEntities, error)
	AddSubsidiary(ctx context.Context, organizationId, subsidiaryId, subsidiaryType string) error
	RemoveSubsidiary(ctx context.Context, organizationId, subsidiaryId string) error

	mapDbNodeToOrganizationEntity(node dbtype.Node) *entity.OrganizationEntity

	UpsertInEventStore(ctx context.Context, size int) (int, int, error)
	UpsertPhoneNumberRelationInEventStore(ctx context.Context, size int) (int, int, error)
	UpsertEmailRelationInEventStore(ctx context.Context, size int) (int, int, error)
}

type OrganizationCreateData struct {
	OrganizationEntity *entity.OrganizationEntity
	OrganizationTypeID *string
	Domains            []string
}

type OrganizationUpdateData struct {
	OrganizationEntity *entity.OrganizationEntity
	OrganizationTypeID *string
	Domains            []string
}

type organizationService struct {
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewOrganizationService(repositories *repository.Repositories, grpcClients *grpc_client.Clients) OrganizationService {
	return &organizationService{
		repositories: repositories,
		grpcClients:  grpcClients,
	}
}

func (s *organizationService) Create(ctx context.Context, input *OrganizationCreateData) (*entity.OrganizationEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	organizationDbNodePtr, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		tenant := common.GetTenantFromContext(ctx)

		for _, domain := range input.Domains {
			_, err := s.repositories.DomainRepository.Merge(ctx, entity.DomainEntity{
				Domain:    domain,
				Source:    input.OrganizationEntity.Source,
				AppSource: input.OrganizationEntity.AppSource,
			})
			if err != nil {
				return nil, err
			}
		}

		organizationDbNodePtr, err := s.repositories.OrganizationRepository.Create(ctx, tx, tenant, *input.OrganizationEntity)
		if err != nil {
			return nil, err
		}
		var organizationId = utils.GetPropsFromNode(*organizationDbNodePtr)["id"].(string)

		err = s.repositories.OrganizationRepository.LinkWithDomainsInTx(ctx, tx, tenant, organizationId, input.Domains)
		if err != nil {
			return nil, err
		}

		if input.OrganizationTypeID != nil {
			err = s.repositories.OrganizationRepository.LinkWithOrganizationTypeInTx(ctx, tx, tenant, organizationId, *input.OrganizationTypeID)
			if err != nil {
				return nil, err
			}
		}

		return organizationDbNodePtr, nil
	})
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*organizationDbNodePtr.(*dbtype.Node)), nil
}

func (s *organizationService) Update(ctx context.Context, input *OrganizationUpdateData) (*entity.OrganizationEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	organizationDbNodePtr, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		tenant := common.GetTenantFromContext(ctx)

		for _, domain := range input.Domains {
			_, err := s.repositories.DomainRepository.Merge(ctx, entity.DomainEntity{
				Domain:    domain,
				Source:    input.OrganizationEntity.Source,
				AppSource: input.OrganizationEntity.AppSource,
			})
			if err != nil {
				return nil, err
			}
		}

		organizationDbNodePtr, err := s.repositories.OrganizationRepository.Update(ctx, tx, tenant, *input.OrganizationEntity)
		if err != nil {
			return nil, err
		}
		var organizationId = utils.GetPropsFromNode(*organizationDbNodePtr)["id"].(string)

		err = s.repositories.OrganizationRepository.LinkWithDomainsInTx(ctx, tx, tenant, organizationId, input.Domains)
		if err != nil {
			return nil, err
		}

		err = s.repositories.OrganizationRepository.UnlinkFromDomainsNotInListInTx(ctx, tx, tenant, organizationId, input.Domains)
		if err != nil {
			return nil, err
		}

		err = s.repositories.OrganizationRepository.UnlinkFromOrganizationTypesInTx(ctx, tx, tenant, organizationId)
		if err != nil {
			return nil, err
		}
		if input.OrganizationTypeID != nil {
			err := s.repositories.OrganizationRepository.LinkWithOrganizationTypeInTx(ctx, tx, tenant, organizationId, *input.OrganizationTypeID)
			if err != nil {
				return nil, err
			}
		}

		return organizationDbNodePtr, nil
	})
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*organizationDbNodePtr.(*dbtype.Node)), nil
}

func (s *organizationService) FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error) {
	session := utils.NewNeo4jReadSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	cypherSort, err := buildSort(sortBy, reflect.TypeOf(entity.OrganizationEntity{}))
	if err != nil {
		return nil, err
	}
	cypherFilter, err := buildFilter(filter, reflect.TypeOf(entity.OrganizationEntity{}))
	if err != nil {
		return nil, err
	}

	dbNodesWithTotalCount, err := s.repositories.OrganizationRepository.GetPaginatedOrganizations(
		ctx,
		session,
		common.GetContext(ctx).Tenant,
		paginatedResult.GetSkip(),
		paginatedResult.GetLimit(),
		cypherFilter,
		cypherSort)
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodesWithTotalCount.Count)

	organizationEntities := entity.OrganizationEntities{}

	for _, v := range dbNodesWithTotalCount.Nodes {
		organizationEntities = append(organizationEntities, *s.mapDbNodeToOrganizationEntity(*v))
	}
	paginatedResult.SetRows(&organizationEntities)
	return &paginatedResult, nil
}

func (s *organizationService) GetOrganizationsForContact(ctx context.Context, contactId string, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error) {
	session := utils.NewNeo4jReadSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	cypherSort, err := buildSort(sortBy, reflect.TypeOf(entity.OrganizationEntity{}))
	if err != nil {
		return nil, err
	}
	cypherFilter, err := buildFilter(filter, reflect.TypeOf(entity.OrganizationEntity{}))
	if err != nil {
		return nil, err
	}

	dbNodesWithTotalCount, err := s.repositories.OrganizationRepository.GetPaginatedOrganizationsForContact(
		ctx,
		session,
		common.GetTenantFromContext(ctx),
		contactId,
		paginatedResult.GetSkip(),
		paginatedResult.GetLimit(),
		cypherFilter,
		cypherSort)
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodesWithTotalCount.Count)

	organizationEntities := entity.OrganizationEntities{}

	for _, v := range dbNodesWithTotalCount.Nodes {
		organizationEntities = append(organizationEntities, *s.mapDbNodeToOrganizationEntity(*v))
	}
	paginatedResult.SetRows(&organizationEntities)
	return &paginatedResult, nil
}

func (s *organizationService) GetOrganizationForJobRole(ctx context.Context, roleId string) (*entity.OrganizationEntity, error) {
	session := utils.NewNeo4jReadSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	dbNode, err := s.repositories.OrganizationRepository.GetOrganizationForJobRole(ctx, session, common.GetContext(ctx).Tenant, roleId)
	if dbNode == nil || err != nil {
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*dbNode), nil
}

func (s *organizationService) GetOrganizationById(ctx context.Context, organizationId string) (*entity.OrganizationEntity, error) {
	dbNode, err := s.repositories.OrganizationRepository.GetOrganizationById(ctx, common.GetTenantFromContext(ctx), organizationId)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToOrganizationEntity(*dbNode), nil
}

func (s *organizationService) PermanentDelete(ctx context.Context, organizationId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	err := s.repositories.OrganizationRepository.Delete(ctx, session, common.GetContext(ctx).Tenant, organizationId)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *organizationService) Merge(ctx context.Context, primaryOrganizationId, mergedOrganizationId string) error {
	session := utils.NewNeo4jWriteSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	_, err := s.GetOrganizationById(ctx, primaryOrganizationId)
	if err != nil {
		logrus.Errorf("Primary organization with id %s not found: %v", primaryOrganizationId, err)
		return err
	}
	_, err = s.GetOrganizationById(ctx, mergedOrganizationId)
	if err != nil {
		logrus.Errorf("Organization to merge with id %s not found: %v", mergedOrganizationId, err)
		return err
	}

	tenant := common.GetContext(ctx).Tenant
	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		err = s.repositories.OrganizationRepository.MergeOrganizationPropertiesInTx(ctx, tx, tenant, primaryOrganizationId, mergedOrganizationId, entity.DataSourceOpenline)
		if err != nil {
			return nil, err
		}

		err = s.repositories.OrganizationRepository.MergeOrganizationRelationsInTx(ctx, tx, tenant, primaryOrganizationId, mergedOrganizationId)
		if err != nil {
			return nil, err
		}

		err = s.repositories.OrganizationRepository.UpdateMergedOrganizationLabelsInTx(ctx, tx, tenant, mergedOrganizationId)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	return err
}

func (s *organizationService) GetOrganizationsForEmails(ctx context.Context, emailIds []string) (*entity.OrganizationEntities, error) {
	organizations, err := s.repositories.OrganizationRepository.GetAllForEmails(ctx, common.GetTenantFromContext(ctx), emailIds)
	if err != nil {
		return nil, err
	}
	organizationEntities := entity.OrganizationEntities{}
	for _, v := range organizations {
		organizationEntity := s.mapDbNodeToOrganizationEntity(*v.Node)
		organizationEntity.DataloaderKey = v.LinkedNodeId
		organizationEntities = append(organizationEntities, *organizationEntity)
	}
	return &organizationEntities, nil
}

func (s *organizationService) GetOrganizationsForPhoneNumbers(ctx context.Context, phoneNumberIds []string) (*entity.OrganizationEntities, error) {
	organizations, err := s.repositories.OrganizationRepository.GetAllForPhoneNumbers(ctx, common.GetTenantFromContext(ctx), phoneNumberIds)
	if err != nil {
		return nil, err
	}
	organizationEntities := entity.OrganizationEntities{}
	for _, v := range organizations {
		organizationEntity := s.mapDbNodeToOrganizationEntity(*v.Node)
		organizationEntity.DataloaderKey = v.LinkedNodeId
		organizationEntities = append(organizationEntities, *organizationEntity)
	}
	return &organizationEntities, nil
}

func (s *organizationService) GetSubsidiaries(ctx context.Context, parentOrganizationId string) (*entity.OrganizationEntities, error) {
	dbEntries, err := s.repositories.OrganizationRepository.GetLinkedSubOrganizations(ctx, common.GetTenantFromContext(ctx), parentOrganizationId, repository.Relationship_Subsidiary)
	if err != nil {
		return nil, err
	}
	organizationEntities := entity.OrganizationEntities{}
	for _, v := range dbEntries {
		organizationEntity := s.mapDbNodeToOrganizationEntity(*v.Node)
		s.addOrganizationRelationshipToOrganizationEntity(*v.Relationship, organizationEntity)
		organizationEntities = append(organizationEntities, *organizationEntity)
	}
	return &organizationEntities, nil
}

func (s *organizationService) AddSubsidiary(ctx context.Context, organizationId, subsidiaryId, subsidiaryType string) error {
	err := s.repositories.OrganizationRepository.LinkSubOrganization(ctx, common.GetTenantFromContext(ctx), organizationId, subsidiaryId, subsidiaryType, repository.Relationship_Subsidiary)
	return err
}

func (s *organizationService) RemoveSubsidiary(ctx context.Context, organizationId, subsidiaryId string) error {
	err := s.repositories.OrganizationRepository.UnlinkSubOrganization(ctx, common.GetTenantFromContext(ctx), organizationId, subsidiaryId, repository.Relationship_Subsidiary)
	return err
}

func (s *organizationService) GetSubsidiaryOf(ctx context.Context, organizationId string) (*entity.OrganizationEntities, error) {
	dbEntries, err := s.repositories.OrganizationRepository.GetLinkedParentOrganizations(ctx, common.GetTenantFromContext(ctx), organizationId, repository.Relationship_Subsidiary)
	if err != nil {
		return nil, err
	}
	organizationEntities := entity.OrganizationEntities{}
	for _, v := range dbEntries {
		organizationEntity := s.mapDbNodeToOrganizationEntity(*v.Node)
		s.addOrganizationRelationshipToOrganizationEntity(*v.Relationship, organizationEntity)
		organizationEntities = append(organizationEntities, *organizationEntity)
	}
	return &organizationEntities, nil
}

func (s *organizationService) UpsertInEventStore(ctx context.Context, size int) (int, int, error) {
	processedRecords := 0
	failedRecords := 0
	outputErr := error(nil)
	for size > 0 {
		batchSize := constants.Neo4jBatchSize
		if size < constants.Neo4jBatchSize {
			batchSize = size
		}
		records, err := s.repositories.OrganizationRepository.GetAllCrossTenants(ctx, batchSize)
		if err != nil {
			return 0, 0, err
		}
		for _, v := range records {
			_, err := s.grpcClients.OrganizationClient.UpsertOrganization(context.Background(), &events_processing_organization.UpsertOrganizationGrpcRequest{
				Id:            utils.GetStringPropOrEmpty(v.Node.Props, "id"),
				Tenant:        v.LinkedNodeId,
				Name:          utils.GetStringPropOrEmpty(v.Node.Props, "name"),
				Description:   utils.GetStringPropOrEmpty(v.Node.Props, "description"),
				Website:       utils.GetStringPropOrEmpty(v.Node.Props, "website"),
				Industry:      utils.GetStringPropOrEmpty(v.Node.Props, "industry"),
				IsPublic:      utils.GetBoolPropOrFalse(v.Node.Props, "isPublic"),
				Source:        utils.GetStringPropOrEmpty(v.Node.Props, "source"),
				SourceOfTruth: utils.GetStringPropOrEmpty(v.Node.Props, "sourceOfTruth"),
				AppSource:     utils.GetStringPropOrEmpty(v.Node.Props, "appSource"),
				CreatedAt:     utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(v.Node.Props, "createdAt")),
				UpdatedAt:     utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(v.Node.Props, "updatedAt")),
			})
			if err != nil {
				failedRecords++
				if outputErr != nil {
					outputErr = err
				}
				logrus.Errorf("Failed to call method: %v", err)
			} else {
				processedRecords++
			}
		}

		size -= batchSize
	}

	return processedRecords, failedRecords, outputErr
}

func (s *organizationService) UpsertPhoneNumberRelationInEventStore(ctx context.Context, size int) (int, int, error) {
	processedRecords := 0
	failedRecords := 0
	outputErr := error(nil)
	for size > 0 {
		batchSize := constants.Neo4jBatchSize
		if size < constants.Neo4jBatchSize {
			batchSize = size
		}
		records, err := s.repositories.OrganizationRepository.GetAllOrganizationPhoneNumberRelationships(ctx, batchSize)
		if err != nil {
			return 0, 0, err
		}
		for _, v := range records {
			_, err := s.grpcClients.OrganizationClient.LinkPhoneNumberToOrganization(context.Background(), &events_processing_organization.LinkPhoneNumberToOrganizationGrpcRequest{
				Primary:        utils.GetBoolPropOrFalse(v.Values[0].(neo4j.Relationship).Props, "primary"),
				Label:          utils.GetStringPropOrEmpty(v.Values[0].(neo4j.Relationship).Props, "label"),
				OrganizationId: v.Values[1].(string),
				PhoneNumberId:  v.Values[2].(string),
				Tenant:         v.Values[3].(string),
			})
			if err != nil {
				failedRecords++
				if outputErr != nil {
					outputErr = err
				}
				logrus.Errorf("Failed to call method: %v", err)
			} else {
				processedRecords++
			}
		}

		size -= batchSize
	}

	return processedRecords, failedRecords, outputErr
}

func (s *organizationService) UpsertEmailRelationInEventStore(ctx context.Context, size int) (int, int, error) {
	processedRecords := 0
	failedRecords := 0
	outputErr := error(nil)
	for size > 0 {
		batchSize := constants.Neo4jBatchSize
		if size < constants.Neo4jBatchSize {
			batchSize = size
		}
		records, err := s.repositories.OrganizationRepository.GetAllOrganizationEmailRelationships(ctx, batchSize)
		if err != nil {
			return 0, 0, err
		}
		for _, v := range records {
			_, err := s.grpcClients.OrganizationClient.LinkEmailToOrganization(context.Background(), &events_processing_organization.LinkEmailToOrganizationGrpcRequest{
				Primary:        utils.GetBoolPropOrFalse(v.Values[0].(neo4j.Relationship).Props, "primary"),
				Label:          utils.GetStringPropOrEmpty(v.Values[0].(neo4j.Relationship).Props, "label"),
				OrganizationId: v.Values[1].(string),
				EmailId:        v.Values[2].(string),
				Tenant:         v.Values[3].(string),
			})
			if err != nil {
				failedRecords++
				if outputErr != nil {
					outputErr = err
				}
				logrus.Errorf("Failed to call method: %v", err)
			} else {
				processedRecords++
			}
		}

		size -= batchSize
	}

	return processedRecords, failedRecords, outputErr
}

func (s *organizationService) mapDbNodeToOrganizationEntity(node dbtype.Node) *entity.OrganizationEntity {
	organizationEntityPtr := new(entity.OrganizationEntity)
	props := utils.GetPropsFromNode(node)
	organizationEntityPtr.ID = utils.GetStringPropOrEmpty(props, "id")
	organizationEntityPtr.Name = utils.GetStringPropOrEmpty(props, "name")
	organizationEntityPtr.Description = utils.GetStringPropOrEmpty(props, "description")
	organizationEntityPtr.Website = utils.GetStringPropOrEmpty(props, "website")
	organizationEntityPtr.Industry = utils.GetStringPropOrEmpty(props, "industry")
	organizationEntityPtr.IsPublic = utils.GetBoolPropOrFalse(props, "isPublic")
	organizationEntityPtr.CreatedAt = utils.GetTimePropOrEpochStart(props, "createdAt")
	organizationEntityPtr.UpdatedAt = utils.GetTimePropOrEpochStart(props, "updatedAt")
	organizationEntityPtr.Source = entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source"))
	organizationEntityPtr.SourceOfTruth = entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth"))
	organizationEntityPtr.AppSource = utils.GetStringPropOrEmpty(props, "appSource")
	return organizationEntityPtr
}

func (s *organizationService) addOrganizationRelationshipToOrganizationEntity(relationship dbtype.Relationship, organizationEntity *entity.OrganizationEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	organizationEntity.LinkedOrganizationType = utils.GetStringPropOrNil(props, "type")
}
