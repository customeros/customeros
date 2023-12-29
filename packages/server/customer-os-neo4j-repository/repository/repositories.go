package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Repositories struct {
	CommentWriteRepository    CommentWriteRepository
	CommonReadRepository      CommonReadRepository
	ContractReadRepository    ContractReadRepository
	CountryReadRepository     CountryReadRepository
	EmailReadRepository       EmailReadRepository
	LogEntryWriteRepository   LogEntryWriteRepository
	PhoneNumberReadRepository PhoneNumberReadRepository
	PlayerWriteRepository     PlayerWriteRepository
	SocialWriteRepository     SocialWriteRepository
	TagWriteRepository        TagWriteRepository
	UserReadRepository        UserReadRepository
	UserWriteRepository       UserWriteRepository
}

func InitNeo4jRepositories(driver *neo4j.DriverWithContext, neo4jDatabase string) *Repositories {
	repositories := Repositories{
		CommentWriteRepository:    NewCommentWriteRepository(driver, neo4jDatabase),
		CommonReadRepository:      NewCommonReadRepository(driver, neo4jDatabase),
		ContractReadRepository:    NewContractReadRepository(driver, neo4jDatabase),
		CountryReadRepository:     NewCountryReadRepository(driver, neo4jDatabase),
		EmailReadRepository:       NewEmailReadRepository(driver, neo4jDatabase),
		LogEntryWriteRepository:   NewLogEntryWriteRepository(driver, neo4jDatabase),
		PhoneNumberReadRepository: NewPhoneNumberReadRepository(driver, neo4jDatabase),
		PlayerWriteRepository:     NewPlayerWriteRepository(driver, neo4jDatabase),
		SocialWriteRepository:     NewSocialWriteRepository(driver, neo4jDatabase),
		TagWriteRepository:        NewTagWriteRepository(driver, neo4jDatabase),
		UserReadRepository:        NewUserReadRepository(driver, neo4jDatabase),
		UserWriteRepository:       NewUserWriteRepository(driver, neo4jDatabase),
	}
	return &repositories
}
