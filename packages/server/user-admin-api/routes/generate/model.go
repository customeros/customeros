package generate

import "time"

type SourceData struct {
	Users []struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
	} `json:"users"`
	Contacts []struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
	} `json:"contacts"`
	Organizations []struct {
		Name   string `json:"name"`
		Domain string `json:"domain"`
		People []struct {
			Email string `json:"email"`
		} `json:"people"`
		Emails []struct {
			From        string    `json:"from"`
			To          []string  `json:"to"`
			Cc          []string  `json:"cc"`
			Bcc         []string  `json:"bcc"`
			Subject     string    `json:"subject"`
			Body        string    `json:"body"`
			ContentType string    `json:"contentType"`
			Date        time.Time `json:"date"`
		} `json:"emails"`
		Meetings []struct {
			CreatedBy string    `json:"createdBy"`
			Attendees []string  `json:"attendees"`
			Subject   string    `json:"subject"`
			Agenda    string    `json:"agenda"`
			StartedAt time.Time `json:"startedAt"`
			EndedAt   time.Time `json:"endedAt"`
		} `json:"meetings"`
		LogEntries []struct {
			CreatedBy   string    `json:"createdBy"`
			Content     string    `json:"content"`
			ContentType string    `json:"contentType"`
			Date        time.Time `json:"date"`
		} `json:"logEntries"`
		Issues []struct {
			CreatedBy   string    `json:"createdBy"`
			CreatedAt   time.Time `json:"createdAt"`
			Subject     string    `json:"subject"`
			Status      string    `json:"status"`
			Priority    string    `json:"priority"`
			Description string    `json:"description"`
		} `json:"issues"`
		Slack [][]struct {
			CreatedBy string    `json:"createdBy"`
			CreatedAt time.Time `json:"createdAt"`
			Message   string    `json:"message"`
		} `json:"slack"`
	} `json:"organizations"`
}
