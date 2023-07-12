package service

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	localEntity "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/zendesk_support/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/zendesk_support/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	common_utils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type zendeskSupportDataService struct {
	airbyteStoreDb       *config.AirbyteStoreDB
	tenant               string
	instance             string
	users                map[string]localEntity.User
	organizations        map[string]localEntity.Organization
	usersAsOrganizations map[string]localEntity.UserAsOrganization
	tickets              map[string]localEntity.Ticket
	ticketComments       map[string]localEntity.TicketComment
}

func NewZendeskSupportDataService(airbyteStoreDb *config.AirbyteStoreDB, tenant string) common.SourceDataService {
	return &zendeskSupportDataService{
		airbyteStoreDb:       airbyteStoreDb,
		tenant:               tenant,
		users:                map[string]localEntity.User{},
		usersAsOrganizations: map[string]localEntity.UserAsOrganization{},
		organizations:        map[string]localEntity.Organization{},
		tickets:              map[string]localEntity.Ticket{},
		ticketComments:       map[string]localEntity.TicketComment{},
	}
}

func (s *zendeskSupportDataService) Refresh() {
	err := s.getDb().AutoMigrate(&localEntity.SyncStatusUser{})
	if err != nil {
		logrus.Error(err)
	}
	err = s.getDb().AutoMigrate(&localEntity.SyncStatusOrganization{})
	if err != nil {
		logrus.Error(err)
	}
	err = s.getDb().AutoMigrate(&localEntity.SyncStatusUserAsOrganization{})
	if err != nil {
		logrus.Error(err)
	}
	err = s.getDb().AutoMigrate(&localEntity.SyncStatusTicket{})
	if err != nil {
		logrus.Error(err)
	}
	err = s.getDb().AutoMigrate(&localEntity.SyncStatusTicketComment{})
	if err != nil {
		logrus.Error(err)
	}
}

func (s *zendeskSupportDataService) getDb() *gorm.DB {
	schemaName := s.SourceId()

	if len(s.instance) > 0 {
		schemaName = schemaName + "_" + s.instance
	}
	schemaName = schemaName + "_" + s.tenant
	return s.airbyteStoreDb.GetDBHandler(&config.Context{
		Schema: schemaName,
	})
}

func (s *zendeskSupportDataService) Close() {
	s.users = make(map[string]localEntity.User)
	s.usersAsOrganizations = make(map[string]localEntity.UserAsOrganization)
	s.organizations = make(map[string]localEntity.Organization)
	s.tickets = make(map[string]localEntity.Ticket)
	s.ticketComments = make(map[string]localEntity.TicketComment)
}

func (s *zendeskSupportDataService) SourceId() string {
	return string(entity.AirbyteSourceZendeskSupport)
}

func (s *zendeskSupportDataService) GetContactsForSync(batchSize int, runId string) []entity.ContactData {
	return nil
}

func (s *zendeskSupportDataService) GetOrganizationsForSync(batchSize int, runId string) []entity.OrganizationData {
	zendeskOrganizations, err := repository.GetOrganizations(s.getDb(), batchSize, runId)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	var zendeskUsersOrganizations localEntity.UsersAsOrganizations
	if len(zendeskOrganizations) < batchSize {
		zendeskUsersOrganizations, err = repository.GetUsersAsOrganizations(s.getDb(), batchSize, runId)
		if err != nil {
			logrus.Error(err)
			return nil
		}
	}

	customerOsOrganizations := make([]entity.OrganizationData, 0, len(zendeskOrganizations)+len(zendeskUsersOrganizations))
	for _, v := range zendeskOrganizations {
		organizationData := entity.OrganizationData{
			ExternalId:          strconv.FormatInt(v.Id, 10),
			ExternalSyncId:      strconv.FormatInt(v.Id, 10),
			ExternalUrl:         v.Url,
			ExternalSystem:      s.SourceId(),
			CreatedAt:           common_utils.TimePtr(v.CreateDate.UTC()),
			UpdatedAt:           common_utils.TimePtr(v.UpdatedDate.UTC()),
			Name:                v.Name,
			ExternalSourceTable: common_utils.StringPtr("organizations"),
		}
		if len(v.Details) > 0 {
			organizationData.Notes = append(organizationData.Notes, entity.OrganizationNote{
				Note:        v.Details,
				FieldSource: "details",
			})
		}
		organizationData.Domains = utils.GetUniqueElements(utils.ConvertJsonbToStringSlice(v.DomainNames))

		customerOsOrganizations = append(customerOsOrganizations, organizationData)
		s.organizations[organizationData.ExternalSyncId] = v
	}

	for _, v := range zendeskUsersOrganizations {
		organizationData := entity.OrganizationData{
			ExternalId:          strconv.FormatInt(v.Id, 10),
			ExternalSyncId:      strconv.FormatInt(v.Id, 10),
			ExternalUrl:         v.Url,
			ExternalSystem:      s.SourceId(),
			ExternalSourceTable: common_utils.StringPtr("users"),
			CreatedAt:           common_utils.TimePtr(v.CreateDate.UTC()),
			UpdatedAt:           common_utils.TimePtr(v.UpdatedDate.UTC()),
			PhoneNumber:         v.Phone,
			Name:                v.Name,
		}
		if len(v.Email) > 0 && !strings.HasSuffix(v.Email, "@without-email.com") {
			organizationData.Email = v.Email
		}
		if len(v.Notes) > 0 {
			organizationData.Notes = append(organizationData.Notes, entity.OrganizationNote{
				Note:        v.Notes,
				FieldSource: "notes",
			})
		}
		if len(v.Details) > 0 {
			organizationData.Notes = append(organizationData.Notes, entity.OrganizationNote{
				Note:        v.Details,
				FieldSource: "details",
			})
		}
		if v.ParentOrganizationId > 0 {
			organizationData.ParentOrganization = &entity.ParentOrganization{
				ExternalId:           strconv.FormatInt(v.ParentOrganizationId, 10),
				OrganizationRelation: entity.Subsidiary,
				Type:                 "store",
			}
		}

		customerOsOrganizations = append(customerOsOrganizations, organizationData)
		s.usersAsOrganizations[organizationData.ExternalSyncId] = v
	}
	return customerOsOrganizations
}

func (s *zendeskSupportDataService) GetUsersForSync(batchSize int, runId string) []entity.UserData {
	zendeskUsers, err := repository.GetUsers(s.getDb(), batchSize, runId)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	customerOsUsers := make([]entity.UserData, 0, len(zendeskUsers))
	for _, v := range zendeskUsers {
		userData := entity.UserData{
			ExternalId:     strconv.FormatInt(v.Id, 10),
			ExternalSystem: s.SourceId(),
			Name:           v.Name,
			Email:          v.Email,
			PhoneNumber:    v.Phone,
			CreatedAt:      common_utils.TimePtr(v.CreateDate.UTC()),
			UpdatedAt:      common_utils.TimePtr(v.UpdatedDate.UTC()),
			ExternalSyncId: strconv.FormatInt(v.Id, 10),
		}
		customerOsUsers = append(customerOsUsers, userData)

		s.users[userData.ExternalSyncId] = v
	}
	return customerOsUsers
}

func (s *zendeskSupportDataService) GetNotesForSync(batchSize int, runId string) []entity.NoteData {
	zendeskInternalTicketComments, err := repository.GetInternalTicketComments(s.getDb(), batchSize, runId)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	notesToReturn := make([]entity.NoteData, 0, len(zendeskInternalTicketComments))

	for _, v := range zendeskInternalTicketComments {
		noteData := entity.NoteData{
			ExternalId:     strconv.FormatInt(v.Id, 10),
			ExternalSyncId: strconv.FormatInt(v.Id, 10),
			ExternalSystem: s.SourceId(),
			CreatedAt:      common_utils.TimePtr(v.CreateDate.UTC()),
			Html:           v.HtmlBody,
			Text:           v.Body,
		}
		if v.TicketId > 0 {
			ticket, err := repository.GetTicket(s.getDb(), v.TicketId)
			if err == nil {
				noteData.MentionedTags = append(noteData.MentionedTags, ticket.Subject+" - "+strconv.FormatInt(v.TicketId, 10))
				noteData.NotedOrganizationsExternalIds = append(noteData.NotedOrganizationsExternalIds, strconv.FormatInt(ticket.RequesterId, 10))
			}
		}
		if v.AuthorId > 0 {
			noteData.CreatorExternalId = strconv.FormatInt(v.AuthorId, 10)
		}

		notesToReturn = append(notesToReturn, noteData)
		s.ticketComments[noteData.ExternalSyncId] = v
	}
	return notesToReturn
}

func (s *zendeskSupportDataService) GetEmailMessagesForSync(batchSize int, runId string) []entity.EmailMessageData {
	return nil
}

func (s *zendeskSupportDataService) MarkContactProcessed(externalSyncId, runId string, synced bool) error {
	return nil
}

func (s *zendeskSupportDataService) MarkOrganizationProcessed(externalSyncId, runId string, synced bool) error {
	organization, ok := s.organizations[externalSyncId]
	if ok {
		err := repository.MarkOrganizationProcessed(s.getDb(), organization, synced, runId)
		if err != nil {
			logrus.Errorf("error while marking organization with external reference %s as synced for zendesk support", externalSyncId)
		}
		return err
	} else if userAsOrganizations, ok := s.usersAsOrganizations[externalSyncId]; ok {
		err := repository.MarkUserAsOrganizationProcessed(s.getDb(), userAsOrganizations, synced, runId)
		if err != nil {
			logrus.Errorf("error while marking organization with external reference %s as synced for zendesk support", externalSyncId)
		}
		return err
	}
	return nil
}

func (s *zendeskSupportDataService) MarkUserProcessed(externalSyncId, runId string, synced bool) error {
	user, ok := s.users[externalSyncId]
	if ok {
		err := repository.MarkUserProcessed(s.getDb(), user, synced, runId)
		if err != nil {
			logrus.Errorf("error while marking owner with external reference %s as synced for zendesk support", externalSyncId)
		}
		return err
	}
	return nil
}

func (s *zendeskSupportDataService) MarkNoteProcessed(externalSyncId, runId string, synced bool) error {
	ticketComment, ok := s.ticketComments[externalSyncId]
	if ok {
		err := repository.MarkTicketCommentProcessed(s.getDb(), ticketComment, synced, runId)
		if err != nil {
			logrus.Errorf("error while marking ticket comment with external reference %s as synced for zendesk support", externalSyncId)
		}
		return err
	}
	return nil
}

func (s *zendeskSupportDataService) MarkEmailMessageProcessed(externalSyncId, runId string, synced bool) error {
	//TODO implement me
	return nil
}

func (s *zendeskSupportDataService) GetIssuesForSync(batchSize int, runId string) []entity.IssueData {
	zendeskTickets, err := repository.GetTickets(s.getDb(), batchSize, runId)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	ticketsToReturn := make([]entity.IssueData, 0, len(zendeskTickets))

	for _, v := range zendeskTickets {
		ticketData := entity.IssueData{
			ExternalId:     strconv.FormatInt(v.Id, 10),
			ExternalSyncId: strconv.FormatInt(v.Id, 10),
			ExternalSystem: s.SourceId(),
			ExternalUrl:    v.Url,
			CreatedAt:      v.CreateDate.UTC(),
			UpdatedAt:      v.UpdatedDate.UTC(),
			Subject:        v.Subject,
			Status:         v.Status,
			Priority:       v.Priority,
			Description:    v.Description,
		}
		ticketData.CollaboratorUserExternalIds = utils.GetUniqueElements(utils.ConvertJsonbToStringSlice(v.CollaboratorIds))
		ticketData.FollowerUserExternalIds = utils.GetUniqueElements(utils.ConvertJsonbToStringSlice(v.FollowerIds))
		if v.RequesterId > 0 {
			ticketData.ReporterOrganizationExternalId = strconv.FormatInt(v.RequesterId, 10)
		}
		if v.AssigneeId > 0 {
			ticketData.AssigneeUserExternalId = strconv.FormatInt(v.AssigneeId, 10)
		}
		if len(v.Type) > 0 {
			ticketData.Tags = append(ticketData.Tags, "type:"+v.Type)
		}
		ticketData.Tags = append(ticketData.Tags, utils.GetUniqueElements(utils.ConvertJsonbToStringSlice(v.Tags))...)

		ticketsToReturn = append(ticketsToReturn, ticketData)
		s.tickets[ticketData.ExternalSyncId] = v
	}
	return ticketsToReturn
}

func (s *zendeskSupportDataService) MarkIssueProcessed(externalSyncId, runId string, synced bool) error {
	ticket, ok := s.tickets[externalSyncId]
	if ok {
		err := repository.MarkTicketProcessed(s.getDb(), ticket, synced, runId)
		if err != nil {
			logrus.Errorf("error while marking ticket with external reference %s as synced for zendesk support", externalSyncId)
		}
		return err
	}
	return nil
}

func (s *zendeskSupportDataService) GetMeetingsForSync(batchSize int, runId string) []entity.MeetingData {
	return nil
}

func (s *zendeskSupportDataService) MarkMeetingProcessed(externalSyncId, runId string, synced bool) error {
	return nil
}

func (s *zendeskSupportDataService) GetInteractionEventsForSync(batchSize int, runId string) []entity.InteractionEventData {
	zendeskPublicTicketComments, err := repository.GetPublicTicketComments(s.getDb(), batchSize, runId)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	interactionEventsToReturn := make([]entity.InteractionEventData, 0, len(zendeskPublicTicketComments))

	for _, v := range zendeskPublicTicketComments {
		interactionEventData := entity.InteractionEventData{
			ExternalId:     strconv.FormatInt(v.Id, 10),
			ExternalSyncId: strconv.FormatInt(v.Id, 10),
			ExternalSystem: s.SourceId(),
			CreatedAt:      v.CreateDate.UTC(),
		}
		interactionEventData.Type = "ISSUE"
		if len(v.HtmlBody) > 0 {
			interactionEventData.Content = v.HtmlBody
			interactionEventData.ContentType = "text/html"
		} else if len(v.PlainBody) > 0 {
			interactionEventData.Content = v.PlainBody
			interactionEventData.ContentType = "text/plain"
		}
		if v.AuthorId > 0 {
			user, err := repository.GetZendeskUser(s.getDb(), v.AuthorId)
			if err == nil {
				participant := entity.InteractionEventParticipant{
					ExternalId: strconv.FormatInt(user.Id, 10),
				}
				if user.IsEndUser() {
					participant.ParticipantType = entity.ORGANIZATION
				} else {
					participant.ParticipantType = entity.USER
				}
				interactionEventData.SentBy = participant
			}
		}
		var ticket localEntity.Ticket
		if v.TicketId > 0 {
			ticket, err = repository.GetTicket(s.getDb(), v.TicketId)
			if err == nil {
				interactionEventData.PartOfExternalId = strconv.FormatInt(v.TicketId, 10)
			}
		}

		if interactionEventData.HasSender() {
			if interactionEventData.SentBy.ParticipantType == entity.USER {
				participant := entity.InteractionEventParticipant{
					ExternalId:      strconv.FormatInt(ticket.RequesterId, 10),
					RelationType:    "TO",
					ParticipantType: entity.ORGANIZATION,
				}
				interactionEventData.SentTo = make(map[string]entity.InteractionEventParticipant)
				interactionEventData.SentTo[participant.ExternalId] = participant
			} else {
				participant := entity.InteractionEventParticipant{
					ExternalId:      strconv.FormatInt(ticket.AssigneeId, 10),
					RelationType:    "TO",
					ParticipantType: entity.USER,
				}
				interactionEventData.SentTo = make(map[string]entity.InteractionEventParticipant)
				interactionEventData.SentTo[participant.ExternalId] = participant

				collaboratorUserExternalIds := utils.GetUniqueElements(utils.ConvertJsonbToStringSlice(ticket.CollaboratorIds))
				for _, collaboratorUserExternalId := range collaboratorUserExternalIds {
					collaboratorParticipant := entity.InteractionEventParticipant{
						ExternalId:      collaboratorUserExternalId,
						RelationType:    "COLLABORATOR",
						ParticipantType: entity.USER,
					}
					_, isPresent := interactionEventData.SentTo[collaboratorParticipant.ExternalId]
					if !isPresent {
						interactionEventData.SentTo[collaboratorParticipant.ExternalId] = collaboratorParticipant
					}
				}

				followerUserExternalIds := utils.GetUniqueElements(utils.ConvertJsonbToStringSlice(ticket.FollowerIds))
				for _, followerUserExternalId := range followerUserExternalIds {
					followerParticipant := entity.InteractionEventParticipant{
						ExternalId:      followerUserExternalId,
						RelationType:    "FOLLOWER",
						ParticipantType: entity.USER,
					}
					_, isPresent := interactionEventData.SentTo[followerParticipant.ExternalId]
					if !isPresent {
						interactionEventData.SentTo[followerParticipant.ExternalId] = followerParticipant
					}
				}
			}
		}

		interactionEventsToReturn = append(interactionEventsToReturn, interactionEventData)
		s.ticketComments[interactionEventData.ExternalSyncId] = v
	}
	return interactionEventsToReturn
}

func (s *zendeskSupportDataService) MarkInteractionEventProcessed(externalSyncId, runId string, synced bool) error {
	ticketComment, ok := s.ticketComments[externalSyncId]
	if ok {
		err := repository.MarkTicketCommentProcessed(s.getDb(), ticketComment, synced, runId)
		if err != nil {
			logrus.Errorf("error while marking ticket comment with external reference %s as synced for zendesk support", externalSyncId)
		}
		return err
	}
	return nil
}
