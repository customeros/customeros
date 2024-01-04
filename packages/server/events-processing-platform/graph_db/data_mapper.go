package graph_db

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"golang.org/x/exp/slices"
)

// Deprecated
func MapDbNodeToOrganizationEntity(node dbtype.Node) *entity.OrganizationEntity {
	props := utils.GetPropsFromNode(node)

	output := entity.OrganizationEntity{
		ID:                 utils.GetStringPropOrEmpty(props, "id"),
		CustomerOsId:       utils.GetStringPropOrEmpty(props, "customerOsId"),
		Name:               utils.GetStringPropOrEmpty(props, "name"),
		Description:        utils.GetStringPropOrEmpty(props, "description"),
		Website:            utils.GetStringPropOrEmpty(props, "website"),
		Industry:           utils.GetStringPropOrEmpty(props, "industry"),
		IndustryGroup:      utils.GetStringPropOrEmpty(props, "industryGroup"),
		SubIndustry:        utils.GetStringPropOrEmpty(props, "subIndustry"),
		TargetAudience:     utils.GetStringPropOrEmpty(props, "targetAudience"),
		ValueProposition:   utils.GetStringPropOrEmpty(props, "valueProposition"),
		LastFundingRound:   utils.GetStringPropOrEmpty(props, "lastFundingRound"),
		LastFundingAmount:  utils.GetStringPropOrEmpty(props, "lastFundingAmount"),
		ReferenceId:        utils.GetStringPropOrEmpty(props, "referenceId"),
		Note:               utils.GetStringPropOrEmpty(props, "note"),
		IsPublic:           utils.GetBoolPropOrFalse(props, "isPublic"),
		IsCustomer:         utils.GetBoolPropOrFalse(props, "isCustomer"),
		Hide:               utils.GetBoolPropOrFalse(props, "hide"),
		Employees:          utils.GetInt64PropOrZero(props, "employees"),
		Market:             utils.GetStringPropOrEmpty(props, "market"),
		CreatedAt:          utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:          utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:             neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:      neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:          utils.GetStringPropOrEmpty(props, "appSource"),
		LastTouchpointAt:   utils.GetTimePropOrNil(props, "lastTouchpointAt"),
		LastTouchpointId:   utils.GetStringPropOrNil(props, "lastTouchpointId"),
		YearFounded:        utils.GetInt64PropOrNil(props, "yearFounded"),
		Headquarters:       utils.GetStringPropOrEmpty(props, "headquarters"),
		EmployeeGrowthRate: utils.GetStringPropOrEmpty(props, "employeeGrowthRate"),
		LogoUrl:            utils.GetStringPropOrEmpty(props, "logoUrl"),
		RenewalSummary: entity.RenewalSummary{
			ArrForecast:            utils.GetFloatPropOrNil(props, "renewalForecastArr"),
			MaxArrForecast:         utils.GetFloatPropOrNil(props, "renewalForecastMaxArr"),
			RenewalLikelihood:      utils.GetStringPropOrEmpty(props, "derivedRenewalLikelihood"),
			RenewalLikelihoodOrder: utils.GetInt64PropOrNil(props, "derivedRenewalLikelihoodOrder"),
			NextRenewalAt:          utils.GetTimePropOrNil(props, "derivedNextRenewalAt"),
		},
		WebScrapeDetails: entity.WebScrapeDetails{
			WebScrapedUrl:             utils.GetStringPropOrEmpty(props, "webScrapedUrl"),
			WebScrapedAt:              utils.GetTimePropOrNil(props, "webScrapedAt"),
			WebScrapeLastRequestedAt:  utils.GetTimePropOrNil(props, "webScrapeLastRequestedAt"),
			WebScrapeLastRequestedUrl: utils.GetStringPropOrEmpty(props, "webScrapeLastRequestedUrl"),
			WebScrapeAttempts:         utils.GetInt64PropOrZero(props, "webScrapeAttempts"),
		},
		OnboardingDetails: entity.OnboardingDetails{
			Status:       utils.GetStringPropOrEmpty(props, "onboardingStatus"),
			SortingOrder: utils.GetInt64PropOrNil(props, "onboardingStatusOrder"),
			UpdatedAt:    utils.GetTimePropOrNil(props, "onboardingUpdatedAt"),
			Comments:     utils.GetStringPropOrEmpty(props, "onboardingComments"),
		},
	}
	return &output
}

// Deprecated
func MapDbNodeToUserEntity(node dbtype.Node) *entity.UserEntity {
	props := utils.GetPropsFromNode(node)
	return &entity.UserEntity{
		Id:              utils.GetStringPropOrEmpty(props, "id"),
		FirstName:       utils.GetStringPropOrEmpty(props, "firstName"),
		LastName:        utils.GetStringPropOrEmpty(props, "lastName"),
		Name:            utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:       utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:       utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:          neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:   neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:       utils.GetStringPropOrEmpty(props, "appSource"),
		Roles:           utils.GetListStringPropOrEmpty(props, "roles"),
		Internal:        utils.GetBoolPropOrFalse(props, "internal"),
		Bot:             utils.GetBoolPropOrFalse(props, "bot"),
		ProfilePhotoUrl: utils.GetStringPropOrEmpty(props, "profilePhotoUrl"),
		Timezone:        utils.GetStringPropOrEmpty(props, "timezone"),
	}
}

// Deprecated
func MapDbNodeToActionEntity(node dbtype.Node) *entity.ActionEntity {
	props := utils.GetPropsFromNode(node)
	action := entity.ActionEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Type:          neo4jentity.GetActionType(utils.GetStringPropOrEmpty(props, "type")),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		Metadata:      utils.GetStringPropOrEmpty(props, "metadata"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &action
}

// Deprecated
func MapDbNodeToLogEntryEntity(node dbtype.Node) *entity.LogEntryEntity {
	props := utils.GetPropsFromNode(node)
	logEntry := entity.LogEntryEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		ContentType:   utils.GetStringPropOrEmpty(props, "contentType"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		StartedAt:     utils.GetTimePropOrEpochStart(props, "startedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &logEntry
}

// Deprecated
func MapDbNodeToSocialEntity(node dbtype.Node) *entity.SocialEntity {
	props := utils.GetPropsFromNode(node)
	social := entity.SocialEntity{
		Id:           utils.GetStringPropOrEmpty(props, "id"),
		PlatformName: utils.GetStringPropOrEmpty(props, "platformName"),
		Url:          utils.GetStringPropOrEmpty(props, "url"),
		CreatedAt:    utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:    utils.GetTimePropOrEpochStart(props, "updatedAt"),
		SourceFields: entity.SourceFields{
			Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
			AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
			SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		},
	}
	return &social
}

// Deprecated
func MapDbNodeToInteractionEventEntity(node dbtype.Node) *entity.InteractionEventEntity {
	props := utils.GetPropsFromNode(node)
	issue := entity.InteractionEventEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		ContentType:   utils.GetStringPropOrEmpty(props, "contentType"),
		Channel:       utils.GetStringPropOrEmpty(props, "channel"),
		ChannelData:   utils.GetStringPropOrEmpty(props, "channelData"),
		Identifier:    utils.GetStringPropOrEmpty(props, "identifier"),
		EventType:     utils.GetStringPropOrEmpty(props, "eventType"),
		Hide:          utils.GetBoolPropOrFalse(props, "hide"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &issue
}

// Deprecated
func MapDbNodeToInteractionSessionEntity(node dbtype.Node) *entity.InteractionSessionEntity {
	props := utils.GetPropsFromNode(node)
	issue := entity.InteractionSessionEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Channel:       utils.GetStringPropOrEmpty(props, "channel"),
		ChannelData:   utils.GetStringPropOrEmpty(props, "channelData"),
		Identifier:    utils.GetStringPropOrEmpty(props, "identifier"),
		Type:          utils.GetStringPropOrEmpty(props, "type"),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		Status:        utils.GetStringPropOrEmpty(props, "status"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &issue
}

// Deprecated
func MapDbNodeToIssueEntity(node dbtype.Node) *entity.IssueEntity {
	props := utils.GetPropsFromNode(node)
	issue := entity.IssueEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Subject:       utils.GetStringPropOrEmpty(props, "subject"),
		Description:   utils.GetStringPropOrEmpty(props, "description"),
		Status:        utils.GetStringPropOrEmpty(props, "status"),
		Priority:      utils.GetStringPropOrEmpty(props, "priority"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &issue
}

// Deprecated
func MapDbNodeToCommentEntity(node dbtype.Node) *entity.CommentEntity {
	props := utils.GetPropsFromNode(node)
	comment := entity.CommentEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		ContentType:   utils.GetStringPropOrEmpty(props, "contentType"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &comment
}

// Deprecated
func MapDbNodeToOpportunityEntity(node *dbtype.Node) *entity.OpportunityEntity {
	if node == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*node)
	opportunity := entity.OpportunityEntity{
		Id:                utils.GetStringPropOrEmpty(props, "id"),
		Name:              utils.GetStringPropOrEmpty(props, "name"),
		Amount:            utils.GetFloatPropOrZero(props, "amount"),
		MaxAmount:         utils.GetFloatPropOrZero(props, "maxAmount"),
		InternalType:      utils.GetStringPropOrEmpty(props, "internalType"),
		ExternalType:      utils.GetStringPropOrEmpty(props, "externalType"),
		InternalStage:     utils.GetStringPropOrEmpty(props, "internalStage"),
		ExternalStage:     utils.GetStringPropOrEmpty(props, "externalStage"),
		EstimatedClosedAt: utils.GetTimePropOrNil(props, "estimatedClosedAt"),
		ClosedAt:          utils.GetTimePropOrNil(props, "closedAt"),
		GeneralNotes:      utils.GetStringPropOrEmpty(props, "generalNotes"),
		NextSteps:         utils.GetStringPropOrEmpty(props, "nextSteps"),
		Comments:          utils.GetStringPropOrEmpty(props, "comments"),
		CreatedAt:         utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:         utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource:         utils.GetStringPropOrEmpty(props, "appSource"),
		Source:            neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:     neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		OwnerUserId:       utils.GetStringPropOrEmpty(props, "ownerUserId"),
		RenewalDetails: entity.RenewalDetails{
			RenewedAt:              utils.GetTimePropOrNil(props, "renewedAt"),
			RenewalLikelihood:      utils.GetStringPropOrEmpty(props, "renewalLikelihood"),
			RenewalUpdatedByUserId: utils.GetStringPropOrEmpty(props, "renewalUpdatedByUserId"),
			RenewalUpdatedByUserAt: utils.GetTimePropOrNil(props, "renewalUpdatedByUserAt"),
		},
	}
	return &opportunity
}

// Deprecated
func MapDbNodeToContractEntity(node *dbtype.Node) *entity.ContractEntity {
	if node == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*node)
	contract := entity.ContractEntity{
		Id:                              utils.GetStringPropOrEmpty(props, "id"),
		Name:                            utils.GetStringPropOrEmpty(props, "name"),
		ContractUrl:                     utils.GetStringPropOrEmpty(props, "contractUrl"),
		CreatedAt:                       utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:                       utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource:                       utils.GetStringPropOrEmpty(props, "appSource"),
		Source:                          neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:                   neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		ServiceStartedAt:                utils.GetTimePropOrNil(props, "serviceStartedAt"),
		SignedAt:                        utils.GetTimePropOrNil(props, "signedAt"),
		EndedAt:                         utils.GetTimePropOrNil(props, "endedAt"),
		RenewalCycle:                    utils.GetStringPropOrEmpty(props, "renewalCycle"),
		RenewalPeriods:                  utils.GetInt64PropOrNil(props, "renewalPeriods"),
		Status:                          utils.GetStringPropOrEmpty(props, "status"),
		TriggeredOnboardingStatusChange: utils.GetBoolPropOrFalse(props, "triggeredOnboardingStatusChange"),
	}
	return &contract
}

// Deprecated
func MapDbNodeToPageView(node dbtype.Node) *entity.PageViewEntity {
	props := utils.GetPropsFromNode(node)
	pageViewAction := entity.PageViewEntity{
		Id:             utils.GetStringPropOrEmpty(props, "id"),
		Application:    utils.GetStringPropOrEmpty(props, "application"),
		TrackerName:    utils.GetStringPropOrEmpty(props, "trackerName"),
		SessionId:      utils.GetStringPropOrEmpty(props, "sessionId"),
		PageUrl:        utils.GetStringPropOrEmpty(props, "pageUrl"),
		PageTitle:      utils.GetStringPropOrEmpty(props, "pageTitle"),
		OrderInSession: utils.GetInt64PropOrZero(props, "orderInSession"),
		EngagedTime:    utils.GetInt64PropOrZero(props, "engagedTime"),
		StartedAt:      utils.GetTimePropOrNow(props, "startedAt"),
		EndedAt:        utils.GetTimePropOrNow(props, "endedAt"),
		Source:         neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:  neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:      utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &pageViewAction
}

// Deprecated
func MapDbNodeToNoteEntity(node dbtype.Node) *entity.NoteEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.NoteEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		ContentType:   utils.GetStringPropOrEmpty(props, "contentType"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &result
}

// Deprecated
func MapDbNodeToMeetingEntity(node dbtype.Node) *entity.MeetingEntity {
	props := utils.GetPropsFromNode(node)
	status := entity.GetMeetingStatus(utils.GetStringPropOrEmpty(props, "status"))
	meetingEntity := entity.MeetingEntity{
		Id:                 utils.GetStringPropOrEmpty(props, "id"),
		Name:               utils.GetStringPropOrNil(props, "name"),
		ConferenceUrl:      utils.GetStringPropOrNil(props, "conferenceUrl"),
		MeetingExternalUrl: utils.GetStringPropOrNil(props, "meetingExternalUrl"),
		Agenda:             utils.GetStringPropOrNil(props, "agenda"),
		AgendaContentType:  utils.GetStringPropOrNil(props, "agendaContentType"),
		CreatedAt:          MigrateStartedAt(props),
		UpdatedAt:          utils.GetTimePropOrNow(props, "updatedAt"),
		StartedAt:          utils.GetTimePropOrNil(props, "startedAt"),
		EndedAt:            utils.GetTimePropOrNil(props, "endedAt"),
		Recording:          utils.GetStringPropOrNil(props, "recording"),
		AppSource:          utils.GetStringPropOrEmpty(props, "appSource"),
		Source:             neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:      neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		Status:             &status,
	}

	return &meetingEntity
}

// Deprecated
func MapDbNodeToAnalysisEntity(node dbtype.Node) *entity.AnalysisEntity {
	props := utils.GetPropsFromNode(node)
	createdAt := utils.GetTimePropOrEpochStart(props, "createdAt")
	analysisEntity := entity.AnalysisEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:     &createdAt,
		AnalysisType:  utils.GetStringPropOrEmpty(props, "analysisType"),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		ContentType:   utils.GetStringPropOrEmpty(props, "contentType"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &analysisEntity
}

func MigrateStartedAt(props map[string]any) time.Time {
	if props["createdAt"] != nil {
		return utils.GetTimePropOrNow(props, "createdAt")
	}
	if props["startedAt"] != nil {
		return utils.GetTimePropOrNow(props, "startedAt")
	}
	return time.Now()
}

// Deprecated
func MapDbNodeToTimelineEvent(dbNode *dbtype.Node) entity.TimelineEvent {
	if slices.Contains(dbNode.Labels, entity.NodeLabel_PageView) {
		return MapDbNodeToPageView(*dbNode)
	} else if slices.Contains(dbNode.Labels, entity.NodeLabel_InteractionSession) {
		return MapDbNodeToInteractionSessionEntity(*dbNode)
	} else if slices.Contains(dbNode.Labels, entity.NodeLabel_Issue) {
		return MapDbNodeToIssueEntity(*dbNode)
	} else if slices.Contains(dbNode.Labels, entity.NodeLabel_Note) {
		return MapDbNodeToNoteEntity(*dbNode)
	} else if slices.Contains(dbNode.Labels, entity.NodeLabel_InteractionEvent) {
		return MapDbNodeToInteractionEventEntity(*dbNode)
	} else if slices.Contains(dbNode.Labels, entity.NodeLabel_Analysis) {
		return MapDbNodeToAnalysisEntity(*dbNode)
	} else if slices.Contains(dbNode.Labels, entity.NodeLabel_Meeting) {
		return MapDbNodeToMeetingEntity(*dbNode)
	} else if slices.Contains(dbNode.Labels, entity.NodeLabel_Action) {
		return MapDbNodeToActionEntity(*dbNode)
	} else if slices.Contains(dbNode.Labels, entity.NodeLabel_LogEntry) {
		return MapDbNodeToLogEntryEntity(*dbNode)
	}
	return nil
}

// Deprecated
func MapDbNodeToServiceLineItemEntity(node dbtype.Node) *entity.ServiceLineItemEntity {
	props := utils.GetPropsFromNode(node)
	serviceLineItem := entity.ServiceLineItemEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		StartedAt:     utils.GetTimePropOrEpochStart(props, "startedAt"),
		EndedAt:       utils.GetTimePropOrNil(props, "endedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		Billed:        utils.GetStringPropOrEmpty(props, "billed"),
		Price:         utils.GetFloatPropOrZero(props, "price"),
		Quantity:      utils.GetInt64PropOrZero(props, "quantity"),
		Comments:      utils.GetStringPropOrEmpty(props, "comments"),
		ParentId:      utils.GetStringPropOrEmpty(props, "parentId"),
		IsCanceled:    utils.GetBoolPropOrFalse(props, "isCanceled"),
	}
	return &serviceLineItem
}

// Deprecated
func MapDbNodeToEmailEntity(node dbtype.Node) *entity.EmailEntity {
	props := utils.GetPropsFromNode(node)
	return &entity.EmailEntity{
		Id:             utils.GetStringPropOrEmpty(props, "id"),
		Email:          utils.GetStringPropOrEmpty(props, "email"),
		RawEmail:       utils.GetStringPropOrEmpty(props, "rawEmail"),
		CreatedAt:      utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:      utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Primary:        utils.GetBoolPropOrFalse(props, "primary"),
		Source:         neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:  neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:      utils.GetStringPropOrEmpty(props, "appSource"),
		Label:          utils.GetStringPropOrEmpty(props, "label"),
		Validated:      utils.GetBoolPropOrNil(props, "validated"),
		IsReachable:    utils.GetStringPropOrNil(props, "isReachable"),
		IsValidSyntax:  utils.GetBoolPropOrNil(props, "isValidSyntax"),
		CanConnectSMTP: utils.GetBoolPropOrNil(props, "canConnectSMTP"),
		AcceptsMail:    utils.GetBoolPropOrNil(props, "acceptsMail"),
		HasFullInbox:   utils.GetBoolPropOrNil(props, "hasFullInbox"),
		IsCatchAll:     utils.GetBoolPropOrNil(props, "isCatchAll"),
		IsDeliverable:  utils.GetBoolPropOrNil(props, "isDeliverable"),
		IsDisabled:     utils.GetBoolPropOrNil(props, "isDisabled"),
		Error:          utils.GetStringPropOrNil(props, "error"),
	}
}
