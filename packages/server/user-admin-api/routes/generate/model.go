package generate

import (
	"time"
)

type SourceData struct {
	Users []struct {
		FirstName       string  `json:"firstName"`
		LastName        string  `json:"lastName"`
		Email           string  `json:"email"`
		ProfilePhotoURL *string `json:"profilePhotoUrl,omitempty"`
	} `json:"users"`
	Contacts []struct {
		FirstName       string  `json:"firstName"`
		LastName        string  `json:"lastName"`
		Email           string  `json:"email"`
		ProfilePhotoURL *string `json:"profilePhotoUrl,omitempty"`
	} `json:"contacts"`
	TenantBillingProfiles []struct {
		LegalName                     string `json:"legalName"`
		Email                         string `json:"email"`
		AddressLine1                  string `json:"addressLine1"`
		Locality                      string `json:"locality"`
		Country                       string `json:"country"`
		Zip                           string `json:"zip"`
		DomesticPaymentsBankInfo      string `json:"domesticPaymentsBankInfo"`
		InternationalPaymentsBankInfo string `json:"internationalPaymentsBankInfo"`
		VatNumber                     string `json:"vatNumber"`
		SendInvoicesFrom              string `json:"sendInvoicesFrom"`
		CanPayWithCard                bool   `json:"canPayWithCard"`
		CanPayWithDirectDebitSEPA     bool   `json:"canPayWithDirectDebitSEPA"`
		CanPayWithDirectDebitACH      bool   `json:"canPayWithDirectDebitACH"`
		CanPayWithDirectDebitBacs     bool   `json:"canPayWithDirectDebitBacs"`
		CanPayWithPigeon              bool   `json:"canPayWithPigeon"`
	} `json:"tenantBillingProfiles"`
	Organizations []struct {
		Id                    string `json:"id"`
		Name                  string `json:"name"`
		Domain                string `json:"domain"`
		OnboardingStatusInput []struct {
			Status   string `json:"status"`
			Comments string `json:"comments"`
		} `json:"onboardingStatusInput"`
		Contracts []struct {
			Name                  string     `json:"name"`
			RenewalCycle          string     `json:"renewalCycle"`
			RenewalPeriods        int64      `json:"renewalPeriods"`
			ContractUrl           string     `json:"contractUrl"`
			ServiceStartedAt      time.Time  `json:"serviceStartedAt"`
			SignedAt              time.Time  `json:"signedAt"`
			InvoicingStartDate    *time.Time `json:"invoicingStartDate"`
			BillingCycle          string     `json:"billingCycle"`
			Currency              string     `json:"currency"`
			AddressLine1          string     `json:"addressLine1"`
			AddressLine2          string     `json:"addressLine2"`
			Zip                   string     `json:"zip"`
			Locality              string     `json:"locality"`
			Country               string     `json:"country"`
			OrganizationLegalName string     `json:"organizationLegalName"`
			InvoiceEmail          string     `json:"invoiceEmail"`
			InvoiceNote           string     `json:"invoiceNote"`
			ServiceLines          []struct {
				Name      string    `json:"description"`
				Billed    string    `json:"billingCycle"`
				Price     int       `json:"price"`
				Quantity  int       `json:"quantity"`
				StartedAt time.Time `json:"serviceStarted"`
				EndedAt   time.Time `json:"serviceEnded,omitempty"`
			} `json:"serviceLines"`
		} `json:"contracts,omitempty"`
		People []struct {
			Email       string `json:"email"`
			JobRole     string `json:"jobRole"`
			Description string `json:"description"`
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
		Intercom [][]struct {
			CreatedBy string    `json:"createdBy"`
			CreatedAt time.Time `json:"createdAt"`
			Message   string    `json:"message"`
		} `json:"intercom"`
	} `json:"organizations"`
	MasterPlans []struct {
		Name       string `json:"name"`
		Milestones []struct {
			Name          string   `json:"name"`
			Order         int64    `json:"order"`
			DurationHours int64    `json:"durationHours"`
			Optional      bool     `json:"optional"`
			Items         []string `json:"items"`
		} `json:"milestones"`
	} `json:"masterPlans"`
}
