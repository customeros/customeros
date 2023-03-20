package dataloader

import (
	"context"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"net/http"
)

type loadersString string

const loadersKey = loadersString("dataloaders")

type Loaders struct {
	TagsForOrganization                    *dataloader.Loader
	TagsForContact                         *dataloader.Loader
	TagsForTicket                          *dataloader.Loader
	EmailsForContact                       *dataloader.Loader
	EmailsForOrganization                  *dataloader.Loader
	LocationsForContact                    *dataloader.Loader
	LocationsForOrganization               *dataloader.Loader
	JobRolesForContact                     *dataloader.Loader
	DomainsForOrganization                 *dataloader.Loader
	NotesForTicket                         *dataloader.Loader
	InteractionEventsForInteractionSession *dataloader.Loader
	InteractionSessionForInteractionEvent  *dataloader.Loader
	SentByParticipantsForInteractionEvent  *dataloader.Loader
	SentToParticipantsForInteractionEvent  *dataloader.Loader
	PhoneNumbersForOrganization            *dataloader.Loader
	PhoneNumbersForUser                    *dataloader.Loader
	PhoneNumbersForContact                 *dataloader.Loader
}

type tagBatcher struct {
	tagService service.TagService
}
type emailBatcher struct {
	emailService service.EmailService
}
type locationBatcher struct {
	locationService service.LocationService
}
type jobRoleBatcher struct {
	jobRoleService service.JobRoleService
}
type domainBatcher struct {
	domainService service.DomainService
}
type noteBatcher struct {
	noteService service.NoteService
}
type interactionEventBatcher struct {
	interactionEventService service.InteractionEventService
}
type interactionSessionBatcher struct {
	interactionSessionService service.InteractionSessionService
}
type interactionEventParticipantBatcher struct {
	interactionEventService service.InteractionEventService
}
type phoneNumberBatcher struct {
	phoneNumberService service.PhoneNumberService
}

// NewDataLoader returns the instantiated Loaders struct for use in a request
func NewDataLoader(services *service.Services) *Loaders {
	tagBatcher := &tagBatcher{
		tagService: services.TagService,
	}
	emailBatcher := &emailBatcher{
		emailService: services.EmailService,
	}
	locationBatcher := &locationBatcher{
		locationService: services.LocationService,
	}
	jobRoleBatcher := &jobRoleBatcher{
		jobRoleService: services.JobRoleService,
	}
	domainBatcher := &domainBatcher{
		domainService: services.DomainService,
	}
	noteBatcher := &noteBatcher{
		noteService: services.NoteService,
	}
	interactionEventBatcher := &interactionEventBatcher{
		interactionEventService: services.InteractionEventService,
	}
	interactionSessionBatcher := &interactionSessionBatcher{
		interactionSessionService: services.InteractionSessionService,
	}
	interactionEventParticipantBatcher := &interactionEventParticipantBatcher{
		interactionEventService: services.InteractionEventService,
	}
	phoneNumberBatcher := &phoneNumberBatcher{
		phoneNumberService: services.PhoneNumberService,
	}
	return &Loaders{
		TagsForOrganization:                    dataloader.NewBatchedLoader(tagBatcher.getTagsForOrganizations, dataloader.WithClearCacheOnBatch()),
		TagsForContact:                         dataloader.NewBatchedLoader(tagBatcher.getTagsForContacts, dataloader.WithClearCacheOnBatch()),
		TagsForTicket:                          dataloader.NewBatchedLoader(tagBatcher.getTagsForTickets, dataloader.WithClearCacheOnBatch()),
		EmailsForContact:                       dataloader.NewBatchedLoader(emailBatcher.getEmailsForContacts, dataloader.WithClearCacheOnBatch()),
		EmailsForOrganization:                  dataloader.NewBatchedLoader(emailBatcher.getEmailsForOrganizations, dataloader.WithClearCacheOnBatch()),
		LocationsForContact:                    dataloader.NewBatchedLoader(locationBatcher.getLocationsForContacts, dataloader.WithClearCacheOnBatch()),
		LocationsForOrganization:               dataloader.NewBatchedLoader(locationBatcher.getLocationsForOrganizations, dataloader.WithClearCacheOnBatch()),
		JobRolesForContact:                     dataloader.NewBatchedLoader(jobRoleBatcher.getJobRolesForContacts, dataloader.WithClearCacheOnBatch()),
		DomainsForOrganization:                 dataloader.NewBatchedLoader(domainBatcher.getDomainsForOrganizations, dataloader.WithClearCacheOnBatch()),
		NotesForTicket:                         dataloader.NewBatchedLoader(noteBatcher.getNotesForTickets, dataloader.WithClearCacheOnBatch()),
		InteractionEventsForInteractionSession: dataloader.NewBatchedLoader(interactionEventBatcher.getInteractionEventsForInteractionSessions, dataloader.WithClearCacheOnBatch()),
		InteractionSessionForInteractionEvent:  dataloader.NewBatchedLoader(interactionSessionBatcher.getInteractionSessionsForInteractionEvents, dataloader.WithClearCacheOnBatch()),
		SentByParticipantsForInteractionEvent:  dataloader.NewBatchedLoader(interactionEventParticipantBatcher.getSentByParticipantsForInteractionEvents, dataloader.WithClearCacheOnBatch()),
		SentToParticipantsForInteractionEvent:  dataloader.NewBatchedLoader(interactionEventParticipantBatcher.getSentToParticipantsForInteractionEvents, dataloader.WithClearCacheOnBatch()),
		PhoneNumbersForOrganization:            dataloader.NewBatchedLoader(phoneNumberBatcher.getPhoneNumbersForOrganizations, dataloader.WithClearCacheOnBatch()),
		PhoneNumbersForUser:                    dataloader.NewBatchedLoader(phoneNumberBatcher.getPhoneNumbersForUsers, dataloader.WithClearCacheOnBatch()),
		PhoneNumbersForContact:                 dataloader.NewBatchedLoader(phoneNumberBatcher.getPhoneNumbersForContacts, dataloader.WithClearCacheOnBatch()),
	}
}

// Middleware injects a DataLoader into the request context, so it can be used later in the schema resolvers
func Middleware(loaders *Loaders, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCtx := context.WithValue(r.Context(), loadersKey, loaders)
		r = r.WithContext(nextCtx)
		next.ServeHTTP(w, r)
	})
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
