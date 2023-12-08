package model

import "time"

type WorkspaceInput struct {
	Name      string  `json:"name"`
	Provider  string  `json:"provider"`
	AppSource *string `json:"appSource"`
}

type EmailInput struct {
	Email     string  `json:"email"`
	Primary   bool    `json:"primary"`
	AppSource *string `json:"appSource"`
}

type PlayerInput struct {
	IdentityId string  `json:"identityId"`
	AuthId     string  `json:"authId"`
	Provider   string  `json:"provider"`
	AppSource  *string `json:"appSource"`
}

type UserInput struct {
	FirstName       string      `json:"firstName"`
	LastName        string      `json:"lastName"`
	Email           EmailInput  `json:"email"`
	Player          PlayerInput `json:"player"`
	AppSource       *string     `json:"appSource"`
	ProfilePhotoURL *string     `json:"profilePhotoUrl,omitempty"`
}

type TenantInput struct {
	Name      string  `json:"name"`
	AppSource *string `json:"appSource"`
}

type ContactInput struct {
	FirstName       *string     `json:"firstName,omitempty"`
	LastName        *string     `json:"lastName,omitempty"`
	Email           *EmailInput `json:"email,omitempty"`
	ProfilePhotoURL *string     `json:"profilePhotoUrl,omitempty"`
}

type ContractInput struct {
	OrganizationId   string    `json:"organizationId"`
	Name             string    `json:"name"`
	RenewalCycle     string    `json:"renewalCycle"`
	RenewalPeriods   int64     `json:"renewalPeriods"`
	ContractUrl      string    `json:"contractUrl"`
	ServiceStartedAt time.Time `json:"serviceStartedAt"`
	SignedAt         time.Time `json:"signedAt"`
}

type ServiceLineInput struct {
	ContractId string    `json:"contractId,omitempty"`
	Name       string    `json:"name"`
	Billed     string    `json:"billed"`
	Price      int       `json:"price"`
	Quantity   int       `json:"quantity"`
	StartedAt  time.Time `json:"startedAt"`
}

type ServiceLineEndedInput struct {
	ContractId string    `json:"contractId,omitempty"`
	Name       string    `json:"name"`
	Billed     string    `json:"billed"`
	Price      int       `json:"price"`
	Quantity   int       `json:"quantity"`
	StartedAt  time.Time `json:"startedAt"`
	EndedAt    time.Time `json:"endedAt"`
}

type InteractionSessionParticipantInput struct {
	Email       *string `json:"email,omitempty"`
	PhoneNumber *string `json:"phoneNumber,omitempty"`
	ContactID   *string `json:"contactID,omitempty"`
	UserID      *string `json:"userID,omitempty"`
	Type        *string `json:"type,omitempty"`
}

type InteractionEventParticipantInput struct {
	Email       *string `json:"email,omitempty"`
	PhoneNumber *string `json:"phoneNumber,omitempty"`
	ContactID   *string `json:"contactID,omitempty"`
	UserID      *string `json:"userID,omitempty"`
	Type        *string `json:"type,omitempty"`
}

type EmailChannelData struct {
	Subject   string   `json:"Subject"`
	InReplyTo []string `json:"InReplyTo"`
	Reference []string `json:"Reference"`
}

type MeetingParticipantInput struct {
	ContactID      *string `json:"contactId,omitempty"`
	UserID         *string `json:"userId,omitempty"`
	OrganizationID *string `json:"organizationId,omitempty"`
}

type NoteInput struct {
	Content     *string `json:"content,omitempty"`
	ContentType *string `json:"contentType,omitempty"`
	AppSource   *string `json:"appSource,omitempty"`
}

type MeetingStatus string

const (
	MeetingStatusUndefined MeetingStatus = "UNDEFINED"
	MeetingStatusAccepted  MeetingStatus = "ACCEPTED"
	MeetingStatusCanceled  MeetingStatus = "CANCELED"
)

type MeetingInput struct {
	Name               *string                    `json:"name,omitempty"`
	AttendedBy         []*MeetingParticipantInput `json:"attendedBy,omitempty"`
	CreatedBy          []*MeetingParticipantInput `json:"createdBy,omitempty"`
	CreatedAt          *time.Time                 `json:"createdAt,omitempty"`
	StartedAt          *time.Time                 `json:"startedAt,omitempty"`
	EndedAt            *time.Time                 `json:"endedAt,omitempty"`
	ConferenceURL      *string                    `json:"conferenceUrl,omitempty"`
	MeetingExternalURL *string                    `json:"meetingExternalUrl,omitempty"`
	Agenda             *string                    `json:"agenda,omitempty"`
	AgendaContentType  *string                    `json:"agendaContentType,omitempty"`
	Note               *NoteInput                 `json:"note,omitempty"`
	AppSource          *string                    `json:"appSource,omitempty"`
	Status             *MeetingStatus             `json:"status,omitempty"`
}
