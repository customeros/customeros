package subscriptions

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
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

	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.PhoneNumberValidationSubscription.GroupName,
		&esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: []string{s.cfg.Subscriptions.PhoneNumberPrefix}},
	); err != nil {
		return err
	}

	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.LocationValidationSubscription.GroupName,
		&esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: []string{s.cfg.Subscriptions.LocationPrefix}},
	); err != nil {
		return err
	}

	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.OrganizationSubscription.GroupName,
		&esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: []string{s.cfg.Subscriptions.OrganizationPrefix}},
	); err != nil {
		return err
	}

	return nil
}

func (s *Subscriptions) subscribeToAll(ctx context.Context, groupName string, filter *esdb.SubscriptionFilter) error {
	s.log.Infof("creating persistent subscription to $all: {%v}", groupName)

	// DO NOT UNCOMMENT THIS LINE, IT WILL DELETE THE PERSISTENT SUBSCRIPTION
	//s.db.DeletePersistentSubscriptionToAll(ctx, groupName, esdb.DeletePersistentSubscriptionOptions{})
	settings := esdb.SubscriptionSettingsDefault()
	options := esdb.PersistentAllSubscriptionOptions{
		Settings:  &settings,
		Filter:    filter,
		StartFrom: esdb.Start{},
	}
	//if s.cfg.EventStoreConfig.AdminUsername != "" && s.cfg.EventStoreConfig.AdminPassword != "" {
	//	options.Authenticated = &esdb.Credentials{Login: s.cfg.EventStoreConfig.AdminUsername, Password: s.cfg.EventStoreConfig.AdminPassword}
	//}
	err := s.db.CreatePersistentSubscriptionToAll(ctx, groupName, options)
	if err != nil {
		esdbErr, _ := esdb.FromError(err)
		if !eventstore.IsEventStoreErrorCodeResourceAlreadyExists(err) {
			s.log.Fatalf("(Subscriptions.CreatePersistentSubscriptionToAll) err code: {%v}", esdbErr.Code())
		} else {
			s.log.Warnf("(Subscriptions.CreatePersistentSubscriptionToAll) err code: {%v}", esdbErr.Code())
			// UPDATING PERSISTENT SUBSCRIPTION IS NOT WORKING AS EXPECTED, FILTERS ARE REMOVED AFTER UPDATE
			//err = s.db.UpdatePersistentSubscriptionToAll(ctx, groupName, options)
			//if err != nil {
			//	s.log.Fatalf("(EmailValidationConsumer.UpdatePersistentSubscriptionToAll) err: {%v}", err.Error())
			//}
		}
	}
	return nil
}
