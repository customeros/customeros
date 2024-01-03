package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Repositories struct {
	ActionReadRepository              ActionReadRepository
	ActionWriteRepository             ActionWriteRepository
	CommentWriteRepository            CommentWriteRepository
	CommonReadRepository              CommonReadRepository
	ContactWriteRepository            ContactWriteRepository
	ContractReadRepository            ContractReadRepository
	ContractWriteRepository           ContractWriteRepository
	CountryReadRepository             CountryReadRepository
	CustomFieldWriteRepository        CustomFieldWriteRepository
	EmailReadRepository               EmailReadRepository
	InteractionSessionWriteRepository InteractionSessionWriteRepository
	IssueWriteRepository              IssueWriteRepository
	LogEntryWriteRepository           LogEntryWriteRepository
	MasterPlanWriteRepository         MasterPlanWriteRepository
	PhoneNumberReadRepository         PhoneNumberReadRepository
	PlayerWriteRepository             PlayerWriteRepository
	SocialWriteRepository             SocialWriteRepository
	TagWriteRepository                TagWriteRepository
	UserReadRepository                UserReadRepository
	UserWriteRepository               UserWriteRepository
}

func InitNeo4jRepositories(driver *neo4j.DriverWithContext, neo4jDatabase string) *Repositories {
	repositories := Repositories{
		ActionReadRepository:              NewActionReadRepository(driver, neo4jDatabase),
		ActionWriteRepository:             NewActionWriteRepository(driver, neo4jDatabase),
		CommentWriteRepository:            NewCommentWriteRepository(driver, neo4jDatabase),
		CommonReadRepository:              NewCommonReadRepository(driver, neo4jDatabase),
		ContactWriteRepository:            NewContactWriteRepository(driver, neo4jDatabase),
		ContractReadRepository:            NewContractReadRepository(driver, neo4jDatabase),
		ContractWriteRepository:           NewContractWriteRepository(driver, neo4jDatabase),
		CountryReadRepository:             NewCountryReadRepository(driver, neo4jDatabase),
		CustomFieldWriteRepository:        NewCustomFieldWriteRepository(driver, neo4jDatabase),
		EmailReadRepository:               NewEmailReadRepository(driver, neo4jDatabase),
		InteractionSessionWriteRepository: NewInteractionSessionWriteRepository(driver, neo4jDatabase),
		IssueWriteRepository:              NewIssueWriteRepository(driver, neo4jDatabase),
		LogEntryWriteRepository:           NewLogEntryWriteRepository(driver, neo4jDatabase),
		MasterPlanWriteRepository:         NewMasterPlanWriteRepository(driver, neo4jDatabase),
		PhoneNumberReadRepository:         NewPhoneNumberReadRepository(driver, neo4jDatabase),
		PlayerWriteRepository:             NewPlayerWriteRepository(driver, neo4jDatabase),
		SocialWriteRepository:             NewSocialWriteRepository(driver, neo4jDatabase),
		TagWriteRepository:                NewTagWriteRepository(driver, neo4jDatabase),
		UserReadRepository:                NewUserReadRepository(driver, neo4jDatabase),
		UserWriteRepository:               NewUserWriteRepository(driver, neo4jDatabase),
	}
	return &repositories
}
