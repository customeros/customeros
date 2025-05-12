package organization

import (
	nats_common "github.com/customeros/customeros/packages/server/customer-os-common-module/nats"
	"reflect"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/nats-io/nats.go"
)

type organizationService struct {
	log             logger.Logger
	natsConn        *nats_common.NATSConnections
	postgres        *postgres_repository.Repositories
	neo4j           *neo4j_repository.Repositories
	events          *events.EventsService
	domain          interfaces.DomainService
	industry        interfaces.IndustryService
	user            interfaces.UserService
	social          interfaces.SocialService
	currencyService interfaces.CurrencyService
	contractService interfaces.ContractService
	emailService    interfaces.EmailService
	subscriptions   []*nats.Subscription
}

func NewOrganizationService(
	log logger.Logger,
	natsConn *nats_common.NATSConnections,
	postgres *postgres_repository.Repositories,
	neo4j *neo4j_repository.Repositories,
	events *events.EventsService,
	domain interfaces.DomainService,
	industry interfaces.IndustryService,
	social interfaces.SocialService,
	user interfaces.UserService,
	currencyService interfaces.CurrencyService,
	emailService interfaces.EmailService,
) interfaces.OrganizationService {
	return &organizationService{
		log:             log,
		natsConn:        natsConn,
		postgres:        postgres,
		neo4j:           neo4j,
		events:          events,
		domain:          domain,
		industry:        industry,
		user:            user,
		social:          social,
		currencyService: currencyService,
		emailService:    emailService,
		subscriptions:   make([]*nats.Subscription, 0),
	}
}

func (s *organizationService) SetSocialService(social interfaces.SocialService) {
	s.social = social
}

func (s *organizationService) SetContractService(contractService interfaces.ContractService) {
	s.contractService = contractService
}

func (s *organizationService) IsInitialized() bool {
	return utils.IsInitialized(s, reflect.TypeOf((*nats_common.NATSConnections)(nil)))
}
