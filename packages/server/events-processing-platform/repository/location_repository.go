package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type LocationRepository interface {
	CreateLocation(ctx context.Context, locationId string, event events.LocationCreateEvent) error
	UpdateLocation(ctx context.Context, locationId string, event events.LocationUpdateEvent) error
	FailLocationValidation(ctx context.Context, locationId string, event events.LocationFailedValidationEvent) error
	LocationValidated(ctx context.Context, locationId string, event events.LocationValidatedEvent) error
	LinkWithOrganization(ctx context.Context, tenant, organizationId, locationId string, updatedAt time.Time) error
	LinkWithContact(ctx context.Context, tenant, contactId, locationId string, updatedAt time.Time) error
}

type locationRepository struct {
	driver *neo4j.DriverWithContext
}

func NewLocationRepository(driver *neo4j.DriverWithContext) LocationRepository {
	return &locationRepository{
		driver: driver,
	}
}

func (r *locationRepository) CreateLocation(ctx context.Context, locationId string, event events.LocationCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationRepository.CreateLocation")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("locationId", locationId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant}) 
		 MERGE (t)<-[:LOCATION_BELONGS_TO_TENANT]-(l:Location:Location_%s {id:$id}) 
		 ON CREATE SET l.rawAddress = $rawAddress,
						l.name = $name,
						l.country = $country,
						l.region = $region,
						l.district = $district,
						l.locality = $locality,
						l.street = $street,	
						l.address = $address,
						l.address2 = $address2,
						l.zip = $zip,
						l.addressType = $addressType,
						l.houseNumber = $houseNumber,
						l.postalCode = $postalCode,
						l.plusFour = $plusFour,
						l.commercial = $commercial,
						l.predirection = $predirection,
						l.validated = null,
						l.latitude = $latitude,
						l.longitude = $longitude,
						l.timeZone = $timeZone,
						l.utcOffset = $utcOffset,
						l.source = $source,
						l.sourceOfTruth = $sourceOfTruth,
						l.appSource = $appSource,
						l.createdAt = $createdAt,
						l.updatedAt = $updatedAt,
						l.syncedWithEventStore = true 
		 ON MATCH SET 	l.syncedWithEventStore = true
`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":            locationId,
				"tenant":        event.Tenant,
				"source":        helper.GetSource(utils.StringFirstNonEmpty(event.SourceFields.Source, event.Source)),
				"sourceOfTruth": helper.GetSourceOfTruth(utils.StringFirstNonEmpty(event.SourceFields.SourceOfTruth, event.SourceOfTruth)),
				"appSource":     helper.GetAppSource(utils.StringFirstNonEmpty(event.SourceFields.AppSource, event.AppSource)),
				"createdAt":     event.CreatedAt,
				"updatedAt":     event.UpdatedAt,
				"rawAddress":    event.RawAddress,
				"name":          event.Name,
				"latitude":      event.LocationAddress.Latitude,
				"longitude":     event.LocationAddress.Longitude,
				"country":       event.LocationAddress.Country,
				"region":        event.LocationAddress.Region,
				"district":      event.LocationAddress.District,
				"locality":      event.LocationAddress.Locality,
				"street":        event.LocationAddress.Street,
				"address":       event.LocationAddress.Address1,
				"address2":      event.LocationAddress.Address2,
				"zip":           event.LocationAddress.Zip,
				"addressType":   event.LocationAddress.AddressType,
				"houseNumber":   event.LocationAddress.HouseNumber,
				"postalCode":    event.LocationAddress.PostalCode,
				"plusFour":      event.LocationAddress.PlusFour,
				"commercial":    event.LocationAddress.Commercial,
				"predirection":  event.LocationAddress.Predirection,
				"timeZone":      event.LocationAddress.TimeZone,
				"utcOffset":     event.LocationAddress.UtcOffset,
			})
		return nil, err
	})
	return err
}

func (r *locationRepository) UpdateLocation(ctx context.Context, locationId string, event events.LocationUpdateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationRepository.UpdateLocation")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("locationId", locationId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:LOCATION_BELONGS_TO_TENANT]-(l:Location:Location_%s {id:$id})
		 SET l.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE l.sourceOfTruth END,,
			l.updatedAt = $updatedAt,
			l.syncedWithEventStore = true,
			l.rawAddress = $rawAddress,
			l.name = $name,
			l.country = $country,
			l.region = $region,
			l.district = $district,
			l.locality = $locality,
			l.street = $street,	
			l.address = $address,
			l.address2 = $address2,
			l.zip = $zip,
			l.addressType = $addressType,
			l.houseNumber = $houseNumber,
			l.postalCode = $postalCode,
			l.plusFour = $plusFour,
			l.commercial = $commercial,
			l.predirection = $predirection,
			l.latitude = $latitude,
			l.longitude = $longitude,
			l.timeZone = $timeZone,
			l.utcOffset = $utcOffset`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":            locationId,
				"tenant":        event.Tenant,
				"sourceOfTruth": event.Source,
				"updatedAt":     event.UpdatedAt,
				"rawAddress":    event.RawAddress,
				"name":          event.Name,
				"latitude":      event.LocationAddress.Latitude,
				"longitude":     event.LocationAddress.Longitude,
				"country":       event.LocationAddress.Country,
				"region":        event.LocationAddress.Region,
				"district":      event.LocationAddress.District,
				"locality":      event.LocationAddress.Locality,
				"street":        event.LocationAddress.Street,
				"address":       event.LocationAddress.Address1,
				"address2":      event.LocationAddress.Address2,
				"zip":           event.LocationAddress.Zip,
				"addressType":   event.LocationAddress.AddressType,
				"houseNumber":   event.LocationAddress.HouseNumber,
				"postalCode":    event.LocationAddress.PostalCode,
				"plusFour":      event.LocationAddress.PlusFour,
				"commercial":    event.LocationAddress.Commercial,
				"predirection":  event.LocationAddress.Predirection,
				"timeZone":      event.LocationAddress.TimeZone,
				"utcOffset":     event.LocationAddress.UtcOffset,
				"overwrite":     event.Source == constants.SourceOpenline,
			})
		return nil, err
	})
	return err
}

func (r *locationRepository) FailLocationValidation(ctx context.Context, locationId string, event events.LocationFailedValidationEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationRepository.FailLocationValidation")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("locationId", locationId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:LOCATION_BELONGS_TO_TENANT]-(l:Location:Location_%s {id:$id})
		 		SET l.validationError = $validationError,
		     		l.validated = false,
					l.updatedAt = $validatedAt`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":              locationId,
				"tenant":          event.Tenant,
				"validationError": event.ValidationError,
				"validatedAt":     event.ValidatedAt,
			})
		return nil, err
	})
	return err
}

func (r *locationRepository) LocationValidated(ctx context.Context, locationId string, event events.LocationValidatedEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationRepository.LocationValidated")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("locationId", locationId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:LOCATION_BELONGS_TO_TENANT]-(l:Location:Location_%s {id:$id})
		 		SET l.validationError = $validationError,
		     		l.validated = true,
					l.updatedAt = $validatedAt,
					l.commercial = $commercial,
					l.country = CASE WHEN $country <> '' or l.country is null or l.country = '' THEN $country ELSE l.subject END, 
					l.region = CASE WHEN $region <> '' or l.region is null or l.region = '' THEN $region ELSE l.region END, 
					l.district = CASE WHEN $district <> '' or l.district is null or l.district = '' THEN $district ELSE l.district END, 
					l.locality = CASE WHEN $locality <> '' or l.locality is null or l.locality = '' THEN $locality ELSE l.locality END, 
					l.street = CASE WHEN $street <> '' or l.street is null or l.street = '' THEN $street ELSE l.street END, 
					l.address = CASE WHEN $address <> '' or l.address is null or l.address = '' THEN $address ELSE l.address END, 
					l.address2 = CASE WHEN $address2 <> '' or l.address2 is null or l.address2 = '' THEN $address2 ELSE l.address2 END, 
					l.zip = CASE WHEN $zip <> '' or l.zip is null or l.zip = '' THEN $zip ELSE l.zip END, 
					l.addressType = CASE WHEN $addressType <> '' or l.addressType is null or l.addressType = '' THEN $addressType ELSE l.addressType END, 
					l.houseNumber = CASE WHEN $houseNumber <> '' or l.houseNumber is null or l.houseNumber = '' THEN $houseNumber ELSE l.houseNumber END, 
					l.postalCode = CASE WHEN $postalCode <> '' or l.postalCode is null or l.postalCode = '' THEN $postalCode ELSE l.postalCode END, 
					l.plusFour = CASE WHEN $plusFour <> '' or l.plusFour is null or l.plusFour = '' THEN $plusFour ELSE l.plusFour END, 
					l.predirection = CASE WHEN $predirection <> '' or l.predirection is null or l.predirection = '' THEN $predirection ELSE l.predirection END,
					l.latitude = CASE WHEN $latitude is not null or l.latitude is null THEN $latitude ELSE l.latitude END,
					l.longitude = CASE WHEN $longitude is not null or l.longitude is null THEN $longitude ELSE l.longitude END,
					l.timeZone = CASE WHEN $timeZone <> '' or l.timeZone is null or l.timeZone = '' THEN $timeZone ELSE l.timeZone END,
					l.utcOffset = CASE WHEN $utcOffset <> '' or l.utcOffset is null THEN $utcOffset ELSE l.utcOffset END`
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":              locationId,
				"tenant":          event.Tenant,
				"validationError": "",
				"validatedAt":     event.ValidatedAt,
				"latitude":        event.LocationAddress.Latitude,
				"longitude":       event.LocationAddress.Longitude,
				"country":         event.LocationAddress.Country,
				"region":          event.LocationAddress.Region,
				"district":        event.LocationAddress.District,
				"locality":        event.LocationAddress.Locality,
				"street":          event.LocationAddress.Street,
				"address":         event.LocationAddress.Address1,
				"address2":        event.LocationAddress.Address2,
				"zip":             event.LocationAddress.Zip,
				"addressType":     event.LocationAddress.AddressType,
				"houseNumber":     event.LocationAddress.HouseNumber,
				"postalCode":      event.LocationAddress.PostalCode,
				"plusFour":        event.LocationAddress.PlusFour,
				"commercial":      event.LocationAddress.Commercial,
				"predirection":    event.LocationAddress.Predirection,
				"timeZone":        event.LocationAddress.TimeZone,
				"utcOffset":       event.LocationAddress.UtcOffset,
			})
		return nil, err
	})
	return err
}

func (r *locationRepository) LinkWithOrganization(ctx context.Context, tenant, organizationId, locationId string, updatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationRepository.LinkWithOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("locationId", locationId), log.String("organizationId", organizationId))

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$organizationId}),
				(t)<-[:LOCATION_BELONGS_TO_TENANT]-(l:Location {id:$locationId})
		MERGE (o)-[:ASSOCIATED_WITH]->(l)
		SET	o.updatedAt = $updatedAt`

	return utils.ExecuteQuery(ctx, *r.driver, query, map[string]any{
		"tenant":         tenant,
		"locationId":     locationId,
		"organizationId": organizationId,
		"updatedAt":      updatedAt,
	})
}

func (r *locationRepository) LinkWithContact(ctx context.Context, tenant, contactId, locationId string, updatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationRepository.LinkWithContact")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("locationId", locationId), log.String("contactId", contactId))

	query := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId}),
				(t)<-[:LOCATION_BELONGS_TO_TENANT]-(l:Location {id:$locationId})
		MERGE (c)-[:ASSOCIATED_WITH]->(l)
		SET	c.updatedAt = $updatedAt`

	return utils.ExecuteQuery(ctx, *r.driver, query, map[string]any{
		"tenant":     tenant,
		"locationId": locationId,
		"contactId":  contactId,
		"updatedAt":  updatedAt,
	})
}
