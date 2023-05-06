package subscriptions

import (
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"golang.org/x/net/context"
)

type Subscriptions struct {
	log logger.Logger
	db  *esdb.Client
	cfg *config.Config
}

func NewSubscriptions(log logger.Logger, db *esdb.Client, cfg *config.Config) *Subscriptions {
	return &Subscriptions{
		log: log,
		db:  db,
		cfg: cfg,
	}
}

func (s *Subscriptions) RefreshSubscriptions(ctx context.Context) error {
	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.GraphSubscription.GroupName,
		nil,
	); err != nil {
		return err
	}

	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.EmailValidationSubscription.GroupName,
		&esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: []string{s.cfg.Subscriptions.EmailPrefix}},
	); err != nil {
		return err
	}

	return nil
}

func (s *Subscriptions) subscribeToAll(ctx context.Context, groupName string, filter *esdb.SubscriptionFilter) error {
	s.log.Infof("creating persistent subscription to $all: {%v}", groupName)

	//s.db.DeletePersistentSubscriptionToAll(ctx, groupName, esdb.DeletePersistentSubscriptionOptions{})
	settings := esdb.SubscriptionSettingsDefault()
	options := esdb.PersistentAllSubscriptionOptions{
		Settings:  &settings,
		Filter:    filter,
		StartFrom: esdb.Start{},
		//Authenticated: &esdb.Credentials{Login: "admin", Password: "changeit"},
	}
	err := s.db.CreatePersistentSubscriptionToAll(ctx, groupName, options)
	if err != nil {
		if !eventstore.IsEventStoreErrorCodeResourceAlreadyExists(err) {
			s.log.Fatalf("(EmailValidationConsumer.CreatePersistentSubscriptionToAll) err: {%v}", err.Error())
		} else {
			//err = s.db.UpdatePersistentSubscriptionToAll(ctx, groupName, options)
			//if err != nil {
			//	s.log.Fatalf("(EmailValidationConsumer.UpdatePersistentSubscriptionToAll) err: {%v}", err.Error())
			//}
		}
	}
	return nil
}
