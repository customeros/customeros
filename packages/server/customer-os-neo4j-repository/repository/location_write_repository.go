package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type AddressDetails struct {
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	Country      string   `json:"country"`
	Region       string   `json:"region"`
	District     string   `json:"district"`
	Locality     string   `json:"locality"`
	Street       string   `json:"street"`
	Address      string   `json:"address"`
	Address2     string   `json:"address2"`
	Zip          string   `json:"zip"`
	AddressType  string   `json:"addressType"`
	HouseNumber  string   `json:"houseNumber"`
	PostalCode   string   `json:"postalCode"`
	PlusFour     string   `json:"plusFour"`
	Commercial   bool     `json:"commercial"`
	Predirection string   `json:"predirection"`
	TimeZone     string   `json:"timeZone"`
	UtcOffset    int      `json:"utcOffset"`
}

type LocationCreateFields struct {
	SourceFields   model.Source   `json:"sourceFields"`
	CreatedAt      time.Time      `json:"createdAt"`
	RawAddress     string         `json:"rawAddress"`
	Name           string         `json:"name"`
	AddressDetails AddressDetails `json:"addressDetails"`
}

type LocationUpdateFields struct {
	AddressDetails AddressDetails `json:"addressDetails"`
	Source         string         `json:"source"`
	RawAddress     string         `json:"rawAddress"`
	Name           string         `json:"name"`
}

type LocationWriteRepository interface {
	CreateLocation(ctx context.Context, tenant, locationId string, data LocationCreateFields) error
	UpdateLocation(ctx context.Context, tenant, locationId string, data LocationUpdateFields) error
	FailLocationValidation(ctx context.Context, tenant, locationId, validationError string, validatedAt time.Time) error
	LocationValidated(ctx context.Context, tenant, locationId string, addressDetails AddressDetails, validatedAt time.Time) error
	LinkWithOrganization(ctx context.Context, tenant, organizationId, locationId string) error
	LinkWithContact(ctx context.Context, tenant, contactId, locationId string) error
}

type locationRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewLocationWriteRepository(driver *neo4j.DriverWithContext, database string) LocationWriteRepository {
	return &locationRepository{
		driver:   driver,
		database: database,
	}
}

func (r *locationRepository) CreateLocation(ctx context.Context, tenant, locationId string, data LocationCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationWriteRepository.CreateLocation")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, locationId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
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
						l.updatedAt = datetime(),
						l.syncedWithEventStore = true 
		 ON MATCH SET 	l.syncedWithEventStore = true`, tenant)
	params := map[string]any{
		"id":            locationId,
		"tenant":        tenant,
		"source":        data.SourceFields.Source,
		"sourceOfTruth": data.SourceFields.SourceOfTruth,
		"appSource":     data.SourceFields.AppSource,
		"createdAt":     data.CreatedAt,
		"rawAddress":    data.RawAddress,
		"name":          data.Name,
		"latitude":      data.AddressDetails.Latitude,
		"longitude":     data.AddressDetails.Longitude,
		"country":       data.AddressDetails.Country,
		"region":        data.AddressDetails.Region,
		"district":      data.AddressDetails.District,
		"locality":      data.AddressDetails.Locality,
		"street":        data.AddressDetails.Street,
		"address":       data.AddressDetails.Address,
		"address2":      data.AddressDetails.Address2,
		"zip":           data.AddressDetails.Zip,
		"addressType":   data.AddressDetails.AddressType,
		"houseNumber":   data.AddressDetails.HouseNumber,
		"postalCode":    data.AddressDetails.PostalCode,
		"plusFour":      data.AddressDetails.PlusFour,
		"commercial":    data.AddressDetails.Commercial,
		"predirection":  data.AddressDetails.Predirection,
		"timeZone":      data.AddressDetails.TimeZone,
		"utcOffset":     data.AddressDetails.UtcOffset,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *locationRepository) UpdateLocation(ctx context.Context, tenant, locationId string, data LocationUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationWriteRepository.UpdateLocation")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, locationId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:LOCATION_BELONGS_TO_TENANT]-(l:Location {id:$id})
			WHERE l:Location_%s
		 SET l.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE l.sourceOfTruth END,
			l.updatedAt = datetime(),
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
			l.utcOffset = $utcOffset`, tenant)
	params := map[string]any{
		"id":            locationId,
		"tenant":        tenant,
		"sourceOfTruth": data.Source,
		"rawAddress":    data.RawAddress,
		"name":          data.Name,
		"latitude":      data.AddressDetails.Latitude,
		"longitude":     data.AddressDetails.Longitude,
		"country":       data.AddressDetails.Country,
		"region":        data.AddressDetails.Region,
		"district":      data.AddressDetails.District,
		"locality":      data.AddressDetails.Locality,
		"street":        data.AddressDetails.Street,
		"address":       data.AddressDetails.Address,
		"address2":      data.AddressDetails.Address2,
		"zip":           data.AddressDetails.Zip,
		"addressType":   data.AddressDetails.AddressType,
		"houseNumber":   data.AddressDetails.HouseNumber,
		"postalCode":    data.AddressDetails.PostalCode,
		"plusFour":      data.AddressDetails.PlusFour,
		"commercial":    data.AddressDetails.Commercial,
		"predirection":  data.AddressDetails.Predirection,
		"timeZone":      data.AddressDetails.TimeZone,
		"utcOffset":     data.AddressDetails.UtcOffset,
		"overwrite":     data.Source == constants.SourceOpenline,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *locationRepository) FailLocationValidation(ctx context.Context, tenant, locationId, validationError string, validatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationWriteRepository.FailLocationValidation")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, locationId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:LOCATION_BELONGS_TO_TENANT]-(l:Location:Location_%s {id:$id})
		 		SET l.validationError = $validationError,
		     		l.validated = false,
					l.updatedAt = datetime()`, tenant)
	params := map[string]any{
		"id":              locationId,
		"tenant":          tenant,
		"validationError": validationError,
		"validatedAt":     validatedAt,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *locationRepository) LocationValidated(ctx context.Context, tenant, locationId string, data AddressDetails, validatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationWriteRepository.LocationValidated")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, locationId)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:LOCATION_BELONGS_TO_TENANT]-(l:Location:Location_%s {id:$id})
		 		SET l.validationError = $validationError,
		     		l.validated = true,
					l.updatedAt = datetime(),
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
					l.utcOffset = CASE WHEN $utcOffset <> '' or l.utcOffset is null THEN $utcOffset ELSE l.utcOffset END`, tenant)
	params := map[string]any{
		"id":              locationId,
		"tenant":          tenant,
		"validationError": "",
		"validatedAt":     validatedAt,
		"latitude":        data.Latitude,
		"longitude":       data.Longitude,
		"country":         data.Country,
		"region":          data.Region,
		"district":        data.District,
		"locality":        data.Locality,
		"street":          data.Street,
		"address":         data.Address,
		"address2":        data.Address2,
		"zip":             data.Zip,
		"addressType":     data.AddressType,
		"houseNumber":     data.HouseNumber,
		"postalCode":      data.PostalCode,
		"plusFour":        data.PlusFour,
		"commercial":      data.Commercial,
		"predirection":    data.Predirection,
		"timeZone":        data.TimeZone,
		"utcOffset":       data.UtcOffset,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *locationRepository) LinkWithOrganization(ctx context.Context, tenant, organizationId, locationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationWriteRepository.LinkWithOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, locationId)
	span.LogFields(log.String("organizationId", organizationId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$organizationId}),
				(t)<-[:LOCATION_BELONGS_TO_TENANT]-(l:Location {id:$locationId})
		MERGE (o)-[:ASSOCIATED_WITH]->(l)
		SET	o.updatedAt = datetime()`
	params := map[string]any{
		"tenant":         tenant,
		"locationId":     locationId,
		"organizationId": organizationId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *locationRepository) LinkWithContact(ctx context.Context, tenant, contactId, locationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationWriteRepository.LinkWithContact")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, locationId)
	span.LogFields(log.String("contactId", contactId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId}),
				(t)<-[:LOCATION_BELONGS_TO_TENANT]-(l:Location {id:$locationId})
		MERGE (c)-[:ASSOCIATED_WITH]->(l)
		SET	c.updatedAt = datetime()`
	params := map[string]any{
		"tenant":     tenant,
		"locationId": locationId,
		"contactId":  contactId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
