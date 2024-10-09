package subscriptions

import (
	"context"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
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
	graphSubscriptionSettings := esdb.SubscriptionSettingsDefault()
	graphSubscriptionSettings.ExtraStatistics = true
	graphSubscriptionSettings.CheckpointLowerBound = s.cfg.Subscriptions.GraphSubscription.CheckpointLowerBound
	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.GraphSubscription.GroupName,
		nil,
		&graphSubscriptionSettings,
		false,
		false,
		esdb.End{},
	); err != nil {
		return err
	}

	graphLowPrioritySubSettings := esdb.SubscriptionSettingsDefault()
	graphLowPrioritySubSettings.ExtraStatistics = true
	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.GraphLowPrioritySubscription.GroupName,
		&esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: s.cfg.Subscriptions.GraphLowPrioritySubscription.Prefixes},
		&graphLowPrioritySubSettings,
		false,
		false,
		esdb.End{},
	); err != nil {
		return err
	}

	phoneNumberValidationSubscriptionSettings := esdb.SubscriptionSettingsDefault()
	phoneNumberValidationSubscriptionSettings.ExtraStatistics = true
	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.PhoneNumberValidationSubscription.GroupName,
		&esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: []string{s.cfg.Subscriptions.PhoneNumberValidationSubscription.Prefix}},
		&phoneNumberValidationSubscriptionSettings,
		false,
		false,
		esdb.End{},
	); err != nil {
		return err
	}

	locationValidationSubscriptionSettings := esdb.SubscriptionSettingsDefault()
	locationValidationSubscriptionSettings.ExtraStatistics = true
	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.LocationValidationSubscription.GroupName,
		&esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: []string{s.cfg.Subscriptions.LocationValidationSubscription.Prefix}},
		&locationValidationSubscriptionSettings,
		false,
		false,
		esdb.End{},
	); err != nil {
		return err
	}

	organizationSubscriptionSettings := esdb.SubscriptionSettingsDefault()
	organizationSubscriptionSettings.MessageTimeout = s.cfg.Subscriptions.OrganizationSubscription.MessageTimeoutSec * 1000
	organizationSubscriptionSettings.CheckpointLowerBound = s.cfg.Subscriptions.OrganizationSubscription.CheckpointLowerBound
	organizationSubscriptionSettings.ExtraStatistics = true
	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.OrganizationSubscription.GroupName,
		&esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: []string{s.cfg.Subscriptions.OrganizationSubscription.Prefix}},
		&organizationSubscriptionSettings,
		false,
		false,
		esdb.End{},
	); err != nil {
		return err
	}

	organizationWebscrapeSubscriptionSettings := esdb.SubscriptionSettingsDefault()
	organizationWebscrapeSubscriptionSettings.MessageTimeout = s.cfg.Subscriptions.OrganizationWebscrapeSubscription.MessageTimeoutSec * 1000
	organizationWebscrapeSubscriptionSettings.CheckpointLowerBound = s.cfg.Subscriptions.OrganizationWebscrapeSubscription.CheckpointLowerBound
	organizationWebscrapeSubscriptionSettings.ExtraStatistics = true
	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.OrganizationWebscrapeSubscription.GroupName,
		&esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: []string{s.cfg.Subscriptions.OrganizationWebscrapeSubscription.Prefix}},
		&organizationWebscrapeSubscriptionSettings,
		false,
		false,
		esdb.End{},
	); err != nil {
		return err
	}

	enrichSubscriptionSettings := esdb.SubscriptionSettingsDefault()
	enrichSubscriptionSettings.MessageTimeout = s.cfg.Subscriptions.EnrichSubscription.MessageTimeoutSec * 1000
	enrichSubscriptionSettings.CheckpointLowerBound = s.cfg.Subscriptions.EnrichSubscription.CheckpointLowerBound
	enrichSubscriptionSettings.ExtraStatistics = true
	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.EnrichSubscription.GroupName,
		nil,
		&enrichSubscriptionSettings,
		false,
		false,
		esdb.End{},
	); err != nil {
		return err
	}

	notificationEventSubscriptionSettings := esdb.SubscriptionSettingsDefault()
	notificationEventSubscriptionSettings.ExtraStatistics = true
	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.NotificationsSubscription.GroupName,
		nil,
		&notificationEventSubscriptionSettings,
		false,
		false,
		esdb.End{},
	); err != nil {
		return err
	}

	invoiceEventSubSettings := esdb.SubscriptionSettingsDefault()
	invoiceEventSubSettings.MessageTimeout = s.cfg.Subscriptions.InvoiceSubscription.MessageTimeoutSec * 1000
	invoiceEventSubSettings.ExtraStatistics = true
	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.InvoiceSubscription.GroupName,
		nil,
		&invoiceEventSubSettings,
		false,
		false,
		esdb.End{},
	); err != nil {
		return err
	}

	reminderEventSubscriptionSettings := esdb.SubscriptionSettingsDefault()
	reminderEventSubscriptionSettings.ExtraStatistics = true
	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.ReminderSubscription.GroupName,
		nil,
		&reminderEventSubscriptionSettings,
		false,
		false,
		esdb.End{},
	); err != nil {
		return err
	}

	notifyRealtimeEventSubscriptionSettings := esdb.SubscriptionSettingsDefault()
	notifyRealtimeEventSubscriptionSettings.ExtraStatistics = false
	if err := s.subscribeToAll(ctx,
		s.cfg.Subscriptions.NotifyRealtimeSubscription.GroupName,
		&esdb.SubscriptionFilter{Type: esdb.StreamFilterType, Prefixes: []string{s.cfg.Subscriptions.NotifyRealtimeSubscription.Prefix}},
		&notifyRealtimeEventSubscriptionSettings,
		false,
		false,
		esdb.End{},
	); err != nil {
		return err
	}

	return nil
}

func (s *Subscriptions) permanentlyDeletePersistentSubscription(ctx context.Context, groupName string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Subscriptions.permanentlyDeletePersistentSubscription")
	defer span.Finish()
	span.LogFields(log.String("groupName", groupName))

	err := s.db.DeletePersistentSubscriptionToAll(ctx, groupName, esdb.DeletePersistentSubscriptionOptions{})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error while deleting persistent subscription: %v", err.Error())
	} else {
		s.log.Infof("persistent subscription deleted: %v", groupName)
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
