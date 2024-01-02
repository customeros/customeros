package subscriptions

import (
	"context"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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
	defaultSettings := esdb.SubscriptionSettingsDefault()
	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.GraphSubscription.GroupName,
		nil,
		&defaultSettings,
		false,
		false,
		esdb.Start{},
	); err != nil {
		return err
	}

	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.EmailValidationSubscription.GroupName,
		&esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: []string{s.cfg.Subscriptions.EmailValidationSubscription.Prefix}},
		&defaultSettings,
		false,
		false,
		esdb.Start{},
	); err != nil {
		return err
	}

	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.PhoneNumberValidationSubscription.GroupName,
		&esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: []string{s.cfg.Subscriptions.PhoneNumberValidationSubscription.Prefix}},
		&defaultSettings,
		false,
		false,
		esdb.Start{},
	); err != nil {
		return err
	}

	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.LocationValidationSubscription.GroupName,
		&esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: []string{s.cfg.Subscriptions.LocationValidationSubscription.Prefix}},
		&defaultSettings,
		false,
		false,
		esdb.Start{},
	); err != nil {
		return err
	}

	organizationSubSettings := esdb.SubscriptionSettingsDefault()
	organizationSubSettings.MessageTimeout = s.cfg.Subscriptions.OrganizationSubscription.MessageTimeoutSec * 1000
	organizationSubSettings.CheckpointLowerBound = s.cfg.Subscriptions.OrganizationSubscription.CheckpointLowerBound
	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.OrganizationSubscription.GroupName,
		&esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: []string{s.cfg.Subscriptions.OrganizationSubscription.Prefix}},
		&organizationSubSettings,
		false,
		false,
		esdb.Start{},
	); err != nil {
		return err
	}

	organizationWebscrapeSubSettings := esdb.SubscriptionSettingsDefault()
	organizationWebscrapeSubSettings.MessageTimeout = s.cfg.Subscriptions.OrganizationWebscrapeSubscription.MessageTimeoutSec * 1000
	organizationWebscrapeSubSettings.CheckpointLowerBound = s.cfg.Subscriptions.OrganizationWebscrapeSubscription.CheckpointLowerBound
	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.OrganizationWebscrapeSubscription.GroupName,
		&esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: []string{s.cfg.Subscriptions.OrganizationWebscrapeSubscription.Prefix}},
		&organizationWebscrapeSubSettings,
		false,
		false,
		esdb.Position{
			Commit:  s.cfg.Subscriptions.OrganizationWebscrapeSubscription.StartPosition,
			Prepare: s.cfg.Subscriptions.OrganizationWebscrapeSubscription.StartPosition,
		},
	); err != nil {
		return err
	}

	interactionEventSubSettings := esdb.SubscriptionSettingsDefault()
	interactionEventSubSettings.MessageTimeout = s.cfg.Subscriptions.InteractionEventSubscription.MessageTimeoutSec * 1000
	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.InteractionEventSubscription.GroupName,
		&esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: []string{s.cfg.Subscriptions.InteractionEventSubscription.Prefix}},
		&interactionEventSubSettings,
		false,
		false,
		esdb.Start{},
	); err != nil {
		return err
	}

	notificationEventSubSettings := esdb.SubscriptionSettingsDefault()
	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.NotificationsSubscription.GroupName,
		nil,
		&notificationEventSubSettings,
		false,
		false,
		esdb.Start{},
	); err != nil {
		return err
	}

	return nil
}

func (s *Subscriptions) subscribeToAll(ctx context.Context, groupName string, filter *esdb.SubscriptionFilter, settings *esdb.PersistentSubscriptionSettings,
	deletePersistentSubscription, updatePersistentSubscription bool, startPosition esdb.AllPosition) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Subscriptions.subscribeToAll")
	defer span.Finish()
	span.LogFields(log.String("groupName", groupName))
	s.log.Infof("creating persistent subscription to $all: {%v}", groupName)

	// USE WITH EXTRA CARE, DELETES PERSISTENT SUBSCRIPTION AND RECREATES IT
	//if deletePersistentSubscription {
	//	err := s.db.DeletePersistentSubscriptionToAll(ctx, groupName, esdb.DeletePersistentSubscriptionOptions{})
	//	if err != nil {
	//		tracing.TraceErr(span, err)
	//		s.log.Errorf("error while deleting persistent subscription: %v", err.Error())
	//	} else {
	//		s.log.Infof("persistent subscription deleted: %v", groupName)
	//	}
	//}
	options := esdb.PersistentAllSubscriptionOptions{
		Settings:  settings,
		Filter:    filter,
		StartFrom: startPosition,
	}
	//if s.cfg.EventStoreConfig.AdminUsername != "" && s.cfg.EventStoreConfig.AdminPassword != "" {
	//	options.Authenticated = &esdb.Credentials{Login: s.cfg.EventStoreConfig.AdminUsername, Password: s.cfg.EventStoreConfig.AdminPassword}
	//}
	err := s.db.CreatePersistentSubscriptionToAll(ctx, groupName, options)
	if err != nil {
		esdbErr, _ := esdb.FromError(err)
		if !eventstore.IsEventStoreErrorCodeResourceAlreadyExists(err) {
			tracing.TraceErr(span, esdbErr)
			s.log.Fatalf("(Subscriptions.CreatePersistentSubscriptionToAll) err code: {%v}", esdbErr.Code())
		} else {
			s.log.Warnf("err code: %v, error: %v", esdbErr.Code(), esdbErr.Error())
			// UPDATING PERSISTENT SUBSCRIPTION IS NOT WORKING AS EXPECTED, FILTERS ARE REMOVED AFTER UPDATE
			//if updatePersistentSubscription {
			//	err = s.db.UpdatePersistentSubscriptionToAll(ctx, groupName, options)
			//	if err != nil {
			//		tracing.TraceErr(span, esdbErr)
			//		s.log.Fatalf("err code: %v, err: %v", esdbErr.Code(), esdbErr.Error())
			//	}
			//}
		}
	}
	return nil
}
