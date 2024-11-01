package test

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

func CleanupAllData(ctx context.Context, driver *neo4j.DriverWithContext) {
	ExecuteWriteQuery(ctx, driver, `MATCH (n) DETACH DELETE n`, map[string]any{})
}

func CreateCountry(ctx context.Context, driver *neo4j.DriverWithContext, entity entity.CountryEntity) string {
	var countryId = entity.Id
	if countryId == "" {
		countryUuid, _ := uuid.NewRandom()
		countryId = countryUuid.String()
	}
	query := `MERGE (c:Country{id:$id}) 
				ON CREATE SET 
					c.phoneCode = $phoneCode,
					c.codeA2 = $codeA2,
					c.codeA3 = $codeA3,
					c.name = $name, 
					c.createdAt = $createdAt, 
					c.updatedAt = $updatedAt`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"id":        entity.Id,
		"codeA2":    entity.CodeA2,
		"codeA3":    entity.CodeA3,
		"phoneCode": entity.PhoneCode,
		"name":      entity.Name,
		"createdAt": utils.Now(),
		"updatedAt": utils.Now(),
	})
	return countryId
}

func CreateTenant(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) {
	query := `MERGE (t:Tenant {name:$tenant}) ON CREATE SET t.createdAt=$now`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant": tenant,
		"now":    utils.Now(),
	})
}

func CreateTenantSettings(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, settings entity.TenantSettingsEntity) string {
	settingsId := utils.NewUUIDIfEmpty(settings.Id)
	query := `MATCH (t:Tenant {name:$tenant}) 
				MERGE (t)-[:HAS_SETTINGS]->(s:TenantSettings {tenant:$tenant})
				ON CREATE SET
					s.id=$id,
					s.createdAt=$createdAt,
					s.updatedAt=$updatedAt,
					s.invoicingEnabled=$invoicingEnabled,
					s.invoicingPostpaid=$invoicingPostpaid,
					s.logoRepositoryFileId=$logoRepositoryFileId,
					s.baseCurrency=$baseCurrency,
					s.enrichContacts=$enrichContacts`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"id":                   settingsId,
		"tenant":               tenant,
		"invoicingEnabled":     settings.InvoicingEnabled,
		"invoicingPostpaid":    settings.InvoicingPostpaid,
		"createdAt":            settings.CreatedAt,
		"updatedAt":            settings.UpdatedAt,
		"baseCurrency":         settings.BaseCurrency,
		"logoRepositoryFileId": settings.LogoRepositoryFileId,
		"enrichContacts":       settings.EnrichContacts,
	})
	return settingsId
}

func CreateTenantBillingProfile(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, profile entity.TenantBillingProfileEntity) string {
	profileId := utils.NewUUIDIfEmpty(profile.Id)
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
				MERGE (t)-[:HAS_BILLING_PROFILE]->(tbp:TenantBillingProfile {id:$profileId})
				ON CREATE SET
					tbp:TenantBillingProfile_%s,
					tbp.createdAt=$createdAt,
					tbp.updatedAt=$updatedAt,
					tbp.phone=$phone,
					tbp.legalName=$legalName,
					tbp.addressLine1=$addressLine1,
					tbp.addressLine2=$addressLine2,
					tbp.addressLine3=$addressLine3,
					tbp.locality=$locality,
					tbp.country=$country,
					tbp.region=$region,
					tbp.zip=$zip,
					tbp.vatNumber=$vatNumber,
					tbp.sendInvoicesFrom=$sendInvoicesFrom,
					tbp.sendInvoicesBcc=$sendInvoicesBcc,
					tbp.canPayWithPigeon=$canPayWithPigeon,
					tbp.canPayWithBankTransfer=$canPayWithBankTransfer,
					tbp.check=$check
`, tenant)
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":                 tenant,
		"profileId":              profileId,
		"createdAt":              profile.CreatedAt,
		"updatedAt":              profile.UpdatedAt,
		"phone":                  profile.Phone,
		"addressLine1":           profile.AddressLine1,
		"addressLine2":           profile.AddressLine2,
		"addressLine3":           profile.AddressLine3,
		"locality":               profile.Locality,
		"country":                profile.Country,
		"region":                 profile.Region,
		"zip":                    profile.Zip,
		"legalName":              profile.LegalName,
		"vatNumber":              profile.VatNumber,
		"sendInvoicesFrom":       profile.SendInvoicesFrom,
		"sendInvoicesBcc":        profile.SendInvoicesBcc,
		"canPayWithPigeon":       profile.CanPayWithPigeon,
		"canPayWithBankTransfer": profile.CanPayWithBankTransfer,
		"check":                  profile.Check,
	})
	return profileId
}

func CreateWorkspace(ctx context.Context, driver *neo4j.DriverWithContext, workspace string, provider string, tenant string) {
	query := `MATCH (t:Tenant {name: $tenant})
			  MERGE (t)-[:HAS_WORKSPACE]->(w:Workspace {name:$workspace, provider:$provider})`

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":    tenant,
		"provider":  provider,
		"workspace": workspace,
	})
}

func CreateDefaultUser(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) string {
	return CreateUser(ctx, driver, tenant, entity.UserEntity{
		FirstName:     "first",
		LastName:      "last",
		Source:        "openline",
		SourceOfTruth: "openline",
	})
}

func CreateDefaultUserAlpha(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) string {
	return CreateUser(ctx, driver, tenant, entity.UserEntity{
		FirstName:     "alpha",
		LastName:      "alpha",
		Source:        "openline",
		SourceOfTruth: "openline",
	})
}

func CreateDefaultUserWithId(ctx context.Context, driver *neo4j.DriverWithContext, tenant, userId string) string {
	return CreateUser(ctx, driver, tenant, entity.UserEntity{
		Id:            userId,
		FirstName:     "first",
		LastName:      "last",
		Source:        "openline",
		SourceOfTruth: "openline",
	})
}

func CreateUser(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, user entity.UserEntity) string {
	now := utils.Now()
	createdAt := user.CreatedAt
	if createdAt.IsZero() {
		createdAt = now
	}
	updatedAt := user.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = now
	}

	userId := utils.NewUUIDIfEmpty(user.Id)
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
			MERGE (u:User {id: $userId})-[:USER_BELONGS_TO_TENANT]->(t)
			SET u:User_%s, 
				u.roles=$roles,
				u.internal=$internal,
				u.bot=$bot,
				u.firstName=$firstName,
				u.lastName=$lastName,
				u.profilePhotoUrl=$profilePhotoUrl,
				u.createdAt=$createdAt,
				u.updatedAt=$updatedAt,
				u.source=$source,
				u.sourceOfTruth=$sourceOfTruth,
				u.appSource=$appSource`, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":          tenant,
		"userId":          userId,
		"firstName":       user.FirstName,
		"lastName":        user.LastName,
		"source":          user.Source,
		"sourceOfTruth":   user.SourceOfTruth,
		"appSource":       user.AppSource,
		"roles":           user.Roles,
		"internal":        user.Internal,
		"bot":             user.Bot,
		"profilePhotoUrl": user.ProfilePhotoUrl,
		"createdAt":       createdAt,
		"updatedAt":       updatedAt,
	})
	return userId
}

func CreateUserWithId(ctx context.Context, driver *neo4j.DriverWithContext, tenant, userId string) {
	CreateUser(ctx, driver, tenant, entity.UserEntity{
		Id: userId,
	})
}

func CreateOrganization(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, organization entity.OrganizationEntity) string {
	orgId := utils.NewUUIDIfEmpty(organization.ID)
	now := time.Now().UTC()
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization:Organization_%s {id:$id})
			ON CREATE SET 	org.name=$name, 
							org.customerOsId=$customerOsId,
							org.referenceId=$referenceId,
							org.description=$description, 
							org.website=$website,
							org.industry=$industry, 
							org.subIndustry=$subIndustry,
							org.industryGroup=$industryGroup,
							org.targetAudience=$targetAudience,	
							org.valueProposition=$valueProposition,
							org.lastFundingRound=$lastFundingRound,
							org.lastFundingAmount=$lastFundingAmount,
							org.lastTouchpointAt=$lastTouchpointAt,
							org.lastTouchpointType=$lastTouchpointType,
							org.note=$note,
							org.logoUrl=$logoUrl,
							org.iconUrl=$iconUrl,
							org.yearFounded=$yearFounded,
							org.headquarters=$headquarters,
							org.employeeGrowthRate=$employeeGrowthRate,
							org.isPublic=$isPublic, 
							org.hide=$hide,
							org.icpFit=$icpFit,
							org.createdAt=$now,
							org.updatedAt=$now,
							org.renewalForecastArr=$renewalForecastArr,
							org.renewalForecastMaxArr=$renewalForecastMaxArr,
							org.derivedNextRenewalAt=$derivedNextRenewalAt,
							org.derivedRenewalLikelihood=$derivedRenewalLikelihood,
							org.derivedRenewalLikelihoodOrder=$derivedRenewalLikelihoodOrder,
							org.onboardingStatus=$onboardingStatus,
							org.onboardingStatusOrder=$onboardingStatusOrder,
							org.onboardingUpdatedAt=$onboardingUpdatedAt,
							org.onboardingComments=$onboardingComments,
							org.relationship=$relationship,
							org.stage=$stage,
							org.stageUpdatedAt=$stageUpdatedAt,
							org.leadSource=$leadSource,
							org.derivedChurnedAt=$derivedChurnedAt,
							org.derivedLtv=$derivedLtv,
							org.derivedLtvCurrency=$derivedLtvCurrency
							`, tenant)
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"id":                            orgId,
		"customerOsId":                  organization.CustomerOsId,
		"referenceId":                   organization.ReferenceId,
		"tenant":                        tenant,
		"name":                          organization.Name,
		"description":                   organization.Description,
		"website":                       organization.Website,
		"industry":                      organization.Industry,
		"isPublic":                      organization.IsPublic,
		"subIndustry":                   organization.SubIndustry,
		"industryGroup":                 organization.IndustryGroup,
		"targetAudience":                organization.TargetAudience,
		"valueProposition":              organization.ValueProposition,
		"hide":                          organization.Hide,
		"lastTouchpointAt":              utils.TimePtrAsAny(organization.LastTouchpointAt, &now),
		"lastTouchpointType":            organization.LastTouchpointType,
		"lastFundingRound":              organization.LastFundingRound,
		"lastFundingAmount":             organization.LastFundingAmount,
		"note":                          organization.Note,
		"logoUrl":                       organization.LogoUrl,
		"iconUrl":                       organization.IconUrl,
		"yearFounded":                   organization.YearFounded,
		"headquarters":                  organization.Headquarters,
		"employeeGrowthRate":            organization.EmployeeGrowthRate,
		"renewalForecastArr":            organization.RenewalSummary.ArrForecast,
		"renewalForecastMaxArr":         organization.RenewalSummary.MaxArrForecast,
		"derivedNextRenewalAt":          utils.TimePtrAsAny(organization.RenewalSummary.NextRenewalAt),
		"derivedRenewalLikelihood":      organization.RenewalSummary.RenewalLikelihood,
		"derivedRenewalLikelihoodOrder": organization.RenewalSummary.RenewalLikelihoodOrder,
		"onboardingStatus":              organization.OnboardingDetails.Status,
		"onboardingStatusOrder":         organization.OnboardingDetails.SortingOrder,
		"onboardingUpdatedAt":           utils.TimePtrAsAny(organization.OnboardingDetails.UpdatedAt),
		"onboardingComments":            organization.OnboardingDetails.Comments,
		"now":                           utils.Now(),
		"relationship":                  organization.Relationship.String(),
		"stage":                         organization.Stage.String(),
		"stageUpdatedAt":                utils.TimePtrAsAny(organization.StageUpdatedAt),
		"leadSource":                    organization.LeadSource,
		"derivedChurnedAt":              utils.TimePtrAsAny(organization.DerivedData.ChurnedAt),
		"derivedLtv":                    organization.DerivedData.Ltv,
		"derivedLtvCurrency":            organization.DerivedData.LtvCurrency.String(),
		"icpFit":                        organization.IcpFit,
	})
	return orgId
}

func CreateLogEntry(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, logEntry entity.LogEntryEntity) string {
	logEntryId := utils.NewUUIDIfEmpty(logEntry.Id)
	query := fmt.Sprintf(`
			  MERGE (l:LogEntry {id:$id})
				SET l:LogEntry_%s,
					l:TimelineEvent,
					l:TimelineEvent_%s,
					l.content=$content,
					l.contentType=$contentType,
					l.startedAt=$startedAt
				`, tenant, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":      tenant,
		"id":          logEntryId,
		"content":     logEntry.Content,
		"contentType": logEntry.ContentType,
		"startedAt":   logEntry.StartedAt,
	})
	return logEntryId
}

func CreateBillingProfileForOrganization(ctx context.Context, driver *neo4j.DriverWithContext, tenant, orgId string, billingProfile entity.BillingProfileEntity) string {
	billingProfileId := CreateBillingProfile(ctx, driver, tenant, billingProfile)
	LinkNodes(ctx, driver, orgId, billingProfileId, "HAS_BILLING_PROFILE")
	return billingProfileId
}

func CreateBillingProfile(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, billingProfile entity.BillingProfileEntity) string {
	billingProfileId := utils.NewUUIDIfEmpty(billingProfile.Id)
	query := fmt.Sprintf(`
			  MERGE (bp:BillingProfile {id:$id})
				SET bp:BillingProfile_%s
				`, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant": tenant,
		"id":     billingProfileId,
	})
	return billingProfileId
}

func CreateLogEntryForOrganization(ctx context.Context, driver *neo4j.DriverWithContext, tenant, orgId string, logEntry entity.LogEntryEntity) string {
	logEntryId := CreateLogEntry(ctx, driver, tenant, logEntry)
	LinkNodes(ctx, driver, orgId, logEntryId, "LOGGED")
	return logEntryId
}

func LinkNodes(ctx context.Context, driver *neo4j.DriverWithContext, fromId, toId string, relation string, properties ...map[string]any) {
	query := fmt.Sprintf(`
			  MATCH (from {id: $fromId})
			  MATCH (to {id: $toId})
			  MERGE (from)-[rel:%s]->(to)`, relation)
	params := map[string]any{
		"fromId": fromId,
		"toId":   toId,
	}
	if len(properties) > 0 {
		params["properties"] = properties[0]
		query += " SET rel += $properties"
	}

	ExecuteWriteQuery(ctx, driver, query, params)
}

func CreateEmail(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, entity entity.EmailEntity) string {
	emailId := utils.NewUUIDIfEmpty(entity.Id)
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
								MERGE (e:Email {id:$emailId})
								MERGE (e)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t)
								ON CREATE SET e:Email_%s,
									e.email=$email,
									e.rawEmail=$rawEmail,
									e.isCatchAll=$isCatchAll,
									e.deliverable=$deliverable,
									e.isValidSyntax=$isValidSyntax,
									e.isRoleAccount=$isRoleAccount,
									e.username=$username,
									e.isRisky=$isRisky,
									e.isFirewalled=$isFirewalled,
									e.provider=$provider,
									e.firewall=$firewall,
									e.isMailboxFull=$isMailboxFull,
									e.isFreeAccount=$isFreeAccount,
									e.smtpSuccess=$smtpSuccess,
									e.verifyResponseCode=$verifyResponseCode,
									e.verifyErrorCode=$verifyErrorCode,
									e.verifyDescription=$verifyDescription,
									e.isPrimaryDomain=$isPrimaryDomain,
									e.primaryDomain=$primaryDomain,
									e.alternateEmail=$alternateEmail,
									e.createdAt=$createdAt,
									e.updatedAt=$updatedAt,
									e.work=$work,
									e.retryValidation=$retryValidation
							`, tenant)
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":             tenant,
		"emailId":            emailId,
		"email":              entity.Email,
		"rawEmail":           entity.RawEmail,
		"createdAt":          entity.CreatedAt,
		"updatedAt":          entity.UpdatedAt,
		"isCatchAll":         entity.IsCatchAll,
		"deliverable":        entity.Deliverable,
		"isValidSyntax":      entity.IsValidSyntax,
		"isRoleAccount":      entity.IsRoleAccount,
		"username":           entity.Username,
		"isRisky":            entity.IsRisky,
		"isFirewalled":       entity.IsFirewalled,
		"provider":           entity.Provider,
		"firewall":           entity.Firewall,
		"isMailboxFull":      entity.IsMailboxFull,
		"isFreeAccount":      entity.IsFreeAccount,
		"smtpSuccess":        entity.SmtpSuccess,
		"verifyResponseCode": entity.ResponseCode,
		"verifyErrorCode":    entity.ErrorCode,
		"verifyDescription":  entity.Description,
		"isPrimaryDomain":    entity.IsPrimaryDomain,
		"primaryDomain":      entity.PrimaryDomain,
		"alternateEmail":     entity.AlternateEmail,
		"work":               entity.Work,
		"retryValidation":    entity.RetryValidation,
	})
	return emailId
}

func CreateEmailForEntity(ctx context.Context, driver *neo4j.DriverWithContext, tenant, entityId string, emailEntity entity.EmailEntity) string {
	emailId := CreateEmail(ctx, driver, tenant, emailEntity)
	LinkNodes(ctx, driver, entityId, emailId, "HAS")
	return emailId
}

func CreateLocation(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, location entity.LocationEntity) string {
	locationId := utils.NewUUIDIfEmpty(location.Id)
	query := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})
			  MERGE (t)<-[:LOCATION_BELONGS_TO_TENANT]-(i:Location {id:$id})
				SET i:Location_%s,
					i.name=$name,
					i.createdAt=$createdAt,
					i.updatedAt=$updatedAt,
					i.country=$country,
					i.region=$region,    
					i.locality=$locality,    
					i.address=$address,    
					i.address2=$address2,    
					i.zip=$zip,    
					i.addressType=$addressType,    
					i.houseNumber=$houseNumber,    
					i.postalCode=$postalCode,    
					i.plusFour=$plusFour,    
					i.commercial=$commercial,    
					i.predirection=$predirection,    
					i.district=$district,    
					i.street=$street,    
					i.rawAddress=$rawAddress,    
					i.latitude=$latitude,    
					i.longitude=$longitude,    
					i.timeZone=$timeZone,    
					i.utcOffset=$utcOffset,    
					i.sourceOfTruth=$sourceOfTruth,
					i.source=$source,
					i.appSource=$appSource`, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":        tenant,
		"id":            locationId,
		"name":          location.Name,
		"createdAt":     location.CreatedAt,
		"updatedAt":     location.UpdatedAt,
		"country":       location.Country,
		"region":        location.Region,
		"locality":      location.Locality,
		"address":       location.Address,
		"address2":      location.Address2,
		"zip":           location.Zip,
		"addressType":   location.AddressType,
		"houseNumber":   location.HouseNumber,
		"postalCode":    location.PostalCode,
		"plusFour":      location.PlusFour,
		"commercial":    location.Commercial,
		"predirection":  location.Predirection,
		"district":      location.District,
		"street":        location.Street,
		"rawAddress":    location.RawAddress,
		"latitude":      location.Latitude,
		"longitude":     location.Longitude,
		"timeZone":      location.TimeZone,
		"utcOffset":     location.UtcOffset,
		"sourceOfTruth": location.SourceOfTruth,
		"source":        location.Source,
		"appSource":     location.AppSource,
	})
	return locationId
}

func CreatePhoneNumber(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, phoneNumber entity.PhoneNumberEntity) string {
	phoneNumberId := utils.NewUUIDIfEmpty(phoneNumber.Id)
	query := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})
			  MERGE (t)<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(i:PhoneNumber {id:$id})
				SET i:PhoneNumber_%s,
					i.e164=$e164,
					i.validated=$validated,
					i.rawPhoneNumber=$rawPhoneNumber,
					i.source=$source,
					i.sourceOfTruth=$sourceOfTruth,
					i.appSource=$appSource,
					i.createdAt=$createdAt,
					i.updatedAt=$updatedAt`, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":         tenant,
		"id":             phoneNumberId,
		"e164":           phoneNumber.E164,
		"validated":      phoneNumber.Validated,
		"rawPhoneNumber": phoneNumber.RawPhoneNumber,
		"source":         phoneNumber.Source,
		"sourceOfTruth":  phoneNumber.SourceOfTruth,
		"appSource":      phoneNumber.AppSource,
		"createdAt":      phoneNumber.CreatedAt,
		"updatedAt":      phoneNumber.UpdatedAt,
	})
	return phoneNumberId
}

func CreateContractForOrganization(ctx context.Context, driver *neo4j.DriverWithContext, tenant, organizationId string, contract entity.ContractEntity) string {
	contractId := utils.NewUUIDIfEmpty(contract.Id)
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}), (o:Organization {id:$organizationId})
				MERGE (t)<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$id})<-[:HAS_CONTRACT]-(o)
				SET 
					c:Contract_%s,
					c.source=$source,
					c.sourceOfTruth=$sourceOfTruth,
					c.appSource=$appSource,
					c.createdAt=$createdAt,
					c.updatedAt=$updatedAt,
					c.name=$name,
					c.contractUrl=$contractUrl,
					c.status=$status,
					c.signedAt=$signedAt,
					c.serviceStartedAt=$serviceStartedAt,
					c.endedAt=$endedAt,
					c.currency=$currency,
					c.invoicingStartDate=$invoicingStartDate,
					c.nextInvoiceDate=$nextInvoiceDate,
					c.billingCycleInMonths=$billingCycleInMonths,
					c.addressLine1=$addressLine1,
					c.addressLine2=$addressLine2,
					c.zip=$zip,
					c.locality=$locality,
					c.country=$country,
					c.region=$region,
					c.organizationLegalName=$organizationLegalName,
					c.invoiceEmail=$invoiceEmail,
					c.invoiceNote=$invoiceNote,
					c.canPayWithCard=$canPayWithCard,
					c.canPayWithDirectDebit=$canPayWithDirectDebit,
					c.canPayWithBankTransfer=$canPayWithBankTransfer,
					c.invoicingEnabled=$invoicingEnabled,
					c.payOnline=$payOnline,
					c.payAutomatically=$payAutomatically,
					c.autoRenew=$autoRenew,
					c.check=$check,
					c.dueDays=$dueDays,
					c.lengthInMonths=$lengthInMonths,
					c.approved=$approved,
					c.ltv=$ltv
				`, tenant)

	params := map[string]any{
		"id":                     contractId,
		"organizationId":         organizationId,
		"tenant":                 tenant,
		"source":                 contract.Source,
		"sourceOfTruth":          contract.SourceOfTruth,
		"appSource":              contract.AppSource,
		"createdAt":              contract.CreatedAt,
		"updatedAt":              contract.UpdatedAt,
		"name":                   contract.Name,
		"contractUrl":            contract.ContractUrl,
		"status":                 contract.ContractStatus.String(),
		"signedAt":               utils.ToDateAsAny(contract.SignedAt),
		"serviceStartedAt":       utils.ToDateAsAny(contract.ServiceStartedAt),
		"endedAt":                utils.ToDateAsAny(contract.EndedAt),
		"currency":               contract.Currency.String(),
		"invoicingStartDate":     utils.ToNeo4jDateAsAny(contract.InvoicingStartDate),
		"nextInvoiceDate":        utils.ToNeo4jDateAsAny(contract.NextInvoiceDate),
		"billingCycleInMonths":   contract.BillingCycleInMonths,
		"addressLine1":           contract.AddressLine1,
		"addressLine2":           contract.AddressLine2,
		"zip":                    contract.Zip,
		"locality":               contract.Locality,
		"country":                contract.Country,
		"region":                 contract.Region,
		"organizationLegalName":  contract.OrganizationLegalName,
		"invoiceEmail":           contract.InvoiceEmail,
		"invoiceNote":            contract.InvoiceNote,
		"canPayWithCard":         contract.CanPayWithCard,
		"canPayWithDirectDebit":  contract.CanPayWithDirectDebit,
		"canPayWithBankTransfer": contract.CanPayWithBankTransfer,
		"invoicingEnabled":       contract.InvoicingEnabled,
		"payOnline":              contract.PayOnline,
		"payAutomatically":       contract.PayAutomatically,
		"autoRenew":              contract.AutoRenew,
		"check":                  contract.Check,
		"dueDays":                contract.DueDays,
		"lengthInMonths":         contract.LengthInMonths,
		"approved":               contract.Approved,
		"ltv":                    contract.Ltv,
	}
	ExecuteWriteQuery(ctx, driver, query, params)

	return contractId
}

func CreateOpportunityForContract(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contractId string, opportunity entity.OpportunityEntity) string {
	opportunityId := CreateOpportunity(ctx, driver, tenant, opportunity)
	if opportunity.InternalType == enum.OpportunityInternalTypeRenewal {
		LinkContractWithOpportunity(ctx, driver, contractId, opportunityId, true)
	} else {
		LinkContractWithOpportunity(ctx, driver, contractId, opportunityId, false)
	}
	return opportunityId
}

func CreateOpportunityForOrganization(ctx context.Context, driver *neo4j.DriverWithContext, tenant, organizationId string, opportunity entity.OpportunityEntity) string {
	opportunityId := CreateOpportunity(ctx, driver, tenant, opportunity)
	LinkNodes(ctx, driver, organizationId, opportunityId, "HAS_OPPORTUNITY")
	return opportunityId
}

func CreateOpportunity(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, opportunity entity.OpportunityEntity) string {
	opportunityId := utils.NewUUIDIfEmpty(opportunity.Id)
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
				MERGE (t)<-[:OPPORTUNITY_BELONGS_TO_TENANT]-(op:Opportunity {id:$id})
				SET 
                    op:Opportunity_%s,
					op.name=$name,
					op.source=$source,
					op.sourceOfTruth=$sourceOfTruth,
					op.appSource=$appSource,
					op.amount=$amount,
					op.maxAmount=$maxAmount,
                    op.internalType=$internalType,
					op.externalType=$externalType,
					op.internalStage=$internalStage,
					op.externalStage=$externalStage,
					op.estimatedClosedAt=$estimatedClosedAt,
					op.generalNotes=$generalNotes,
                    op.comments=$comments,
                    op.renewedAt=$renewedAt,
                    op.renewalLikelihood=$renewalLikelihood,
                    op.renewalUpdatedByUserId=$renewalUpdatedByUserId,
                    op.renewalUpdatedByUserAt=$renewalUpdatedByUserAt,
					op.renewalApproved=$renewalApproved,
					op.renewalAdjustedRate=$renewalAdjustedRate,
					op.nextSteps=$nextSteps,
					op.createdAt=$createdAt,
					op.updatedAt=$updatedAt,
					op.currency=$currency,
					op.likelihoodRate=$likelihoodRate
				`, tenant)
	if opportunity.InternalType == enum.OpportunityInternalTypeRenewal {
		query += ", op:RenewalOpportunity"
	}

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"id":                     opportunityId,
		"tenant":                 tenant,
		"name":                   opportunity.Name,
		"source":                 opportunity.Source,
		"sourceOfTruth":          opportunity.SourceOfTruth,
		"appSource":              opportunity.AppSource,
		"amount":                 opportunity.Amount,
		"maxAmount":              opportunity.MaxAmount,
		"internalType":           opportunity.InternalType,
		"externalType":           opportunity.ExternalType,
		"internalStage":          opportunity.InternalStage,
		"externalStage":          opportunity.ExternalStage,
		"estimatedClosedAt":      opportunity.EstimatedClosedAt,
		"generalNotes":           opportunity.GeneralNotes,
		"nextSteps":              opportunity.NextSteps,
		"comments":               opportunity.Comments,
		"renewedAt":              utils.ToDateAsAny(opportunity.RenewalDetails.RenewedAt),
		"renewalLikelihood":      opportunity.RenewalDetails.RenewalLikelihood,
		"renewalUpdatedByUserId": opportunity.RenewalDetails.RenewalUpdatedByUserId,
		"renewalUpdatedByUserAt": utils.TimePtrAsAny(opportunity.RenewalDetails.RenewalUpdatedByUserAt),
		"renewalApproved":        opportunity.RenewalDetails.RenewalApproved,
		"renewalAdjustedRate":    opportunity.RenewalDetails.RenewalAdjustedRate,
		"createdAt":              opportunity.CreatedAt,
		"updatedAt":              opportunity.UpdatedAt,
		"currency":               opportunity.Currency,
		"likelihoodRate":         opportunity.LikelihoodRate,
	})
	return opportunityId
}

func LinkContractWithOpportunity(ctx context.Context, driver *neo4j.DriverWithContext, contractId, opportunityId string, renewal bool) {
	query := `MATCH (c:Contract {id:$contractId}), (o:Opportunity {id:$opportunityId})
				MERGE (c)-[:HAS_OPPORTUNITY]->(o)`
	if renewal {
		query += ` MERGE (c)-[:ACTIVE_RENEWAL]->(o)`
	}
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"contractId":    contractId,
		"opportunityId": opportunityId,
	})
}

func ActiveRenewalOpportunityForContract(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contractId, opportunityId string) string {
	query := fmt.Sprintf(`
				MATCH (c:Contract_%s {id:$contractId}), (op:Opportunity_%s {id:$opportunityId})
				MERGE (c)-[:ACTIVE_RENEWAL]->(op)
				`, tenant, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"opportunityId": opportunityId,
		"contractId":    contractId,
	})
	return opportunityId
}

func CreateServiceLineItemForContract(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contractId string, serviceLineItem entity.ServiceLineItemEntity) string {
	serviceLineItemId := utils.NewUUIDIfEmpty(serviceLineItem.ID)
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})
				MERGE (c)-[:HAS_SERVICE]->(sli:ServiceLineItem {id:$id})
				SET 
					sli:ServiceLineItem_%s,
					sli.name=$name,
					sli.source=$source,
					sli.sourceOfTruth=$sourceOfTruth,
					sli.appSource=$appSource,
					sli.isCanceled=$isCanceled,	
					sli.billed=$billed,	
					sli.quantity=$quantity,	
					sli.price=$price,
					sli.previousBilled=$previousBilled,	
					sli.previousQuantity=$previousQuantity,	
					sli.previousPrice=$previousPrice,
                    sli.comments=$comments,
					sli.startedAt=$startedAt,
					sli.endedAt=$endedAt,
					sli.createdAt=$createdAt,
					sli.updatedAt=$updatedAt,
	                sli.parentId=$parentId,
					sli.vatRate=toFloat($vatRate)
				`, tenant)

	params := map[string]any{
		"id":               serviceLineItemId,
		"contractId":       contractId,
		"tenant":           tenant,
		"name":             serviceLineItem.Name,
		"source":           serviceLineItem.Source,
		"sourceOfTruth":    serviceLineItem.SourceOfTruth,
		"appSource":        serviceLineItem.AppSource,
		"isCanceled":       serviceLineItem.Canceled,
		"billed":           serviceLineItem.Billed,
		"quantity":         serviceLineItem.Quantity,
		"price":            serviceLineItem.Price,
		"previousBilled":   serviceLineItem.PreviousBilled,
		"previousQuantity": serviceLineItem.PreviousQuantity,
		"previousPrice":    serviceLineItem.PreviousPrice,
		"startedAt":        utils.ToDate(serviceLineItem.StartedAt),
		"comments":         serviceLineItem.Comments,
		"createdAt":        serviceLineItem.CreatedAt,
		"updatedAt":        serviceLineItem.UpdatedAt,
		"parentId":         serviceLineItem.ParentID,
		"vatRate":          serviceLineItem.VatRate,
		"endedAt":          utils.ToDateAsAny(serviceLineItem.EndedAt),
	}

	ExecuteWriteQuery(ctx, driver, query, params)
	return serviceLineItemId
}

func InsertContractWithOpportunity(ctx context.Context, driver *neo4j.DriverWithContext, tenant, organizationId string, contract entity.ContractEntity, opportunity entity.OpportunityEntity) string {
	contractId := CreateContractForOrganization(ctx, driver, tenant, organizationId, contract)
	CreateOpportunityForContract(ctx, driver, tenant, contractId, opportunity)
	return contractId
}

func InsertContractWithActiveRenewalOpportunity(ctx context.Context, driver *neo4j.DriverWithContext, tenant, organizationId string, contract entity.ContractEntity, opportunity entity.OpportunityEntity) string {
	contractId := CreateContractForOrganization(ctx, driver, tenant, organizationId, contract)
	opportunityId := CreateOpportunityForContract(ctx, driver, tenant, contractId, opportunity)
	ActiveRenewalOpportunityForContract(ctx, driver, tenant, contractId, opportunityId)
	return contractId
}

func InsertServiceLineItem(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contractId string, billedType enum.BilledType, price float64, quantity int64, startedAt time.Time) string {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		ID:        id,
		ParentID:  id,
		Billed:    billedType,
		Price:     price,
		Quantity:  quantity,
		StartedAt: startedAt,
	})
	return id
}

func InsertServiceLineItemEnded(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contractId string, billedType enum.BilledType, price float64, quantity int64, startedAt, endedAt time.Time) string {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		ID:        id,
		ParentID:  id,
		Billed:    billedType,
		Price:     price,
		Quantity:  quantity,
		StartedAt: startedAt,
		EndedAt:   &endedAt,
	})
	return id
}

func InsertServiceLineItemCanceled(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contractId string, billedType enum.BilledType, price float64, quantity int64, startedAt, endedAt time.Time) string {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		ID:        id,
		ParentID:  id,
		Billed:    billedType,
		Price:     price,
		Quantity:  quantity,
		Canceled:  true,
		StartedAt: startedAt,
		EndedAt:   &endedAt,
	})
	return id
}

func InsertServiceLineItemWithParent(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contractId string, billedType enum.BilledType, price float64, quantity int64, previousBilledType enum.BilledType, previousPrice float64, previousQuantity int64, startedAt time.Time, parentId string) {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		ID:               id,
		ParentID:         parentId,
		Billed:           billedType,
		Price:            price,
		Quantity:         quantity,
		PreviousBilled:   previousBilledType,
		PreviousPrice:    previousPrice,
		PreviousQuantity: previousQuantity,
		StartedAt:        startedAt,
	})
}

func InsertServiceLineItemEndedWithParent(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contractId string, billedType enum.BilledType, price float64, quantity int64, previousBilledType enum.BilledType, previousPrice float64, previousQuantity int64, startedAt, endedAt time.Time, parentId string) {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		ID:               id,
		ParentID:         parentId,
		Billed:           billedType,
		Price:            price,
		Quantity:         quantity,
		PreviousBilled:   previousBilledType,
		PreviousPrice:    previousPrice,
		PreviousQuantity: previousQuantity,
		StartedAt:        startedAt,
		EndedAt:          &endedAt,
	})
}

func InsertServiceLineItemCanceledWithParent(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contractId string, billedType enum.BilledType, price float64, quantity int64, previousBilledType enum.BilledType, previousPrice float64, previousQuantity int64, startedAt, endedAt time.Time, parentId string) {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		ID:               id,
		ParentID:         parentId,
		Billed:           billedType,
		Price:            price,
		Quantity:         quantity,
		PreviousBilled:   previousBilledType,
		PreviousPrice:    previousPrice,
		PreviousQuantity: previousQuantity,
		Canceled:         true,
		StartedAt:        startedAt,
		EndedAt:          &endedAt,
	})
}

func CreateInvoiceForContract(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contractId string, invoice entity.InvoiceEntity) string {
	invoiceId := utils.NewUUIDIfEmpty(invoice.Id)
	query := fmt.Sprintf(`
			MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})
			MERGE (t)<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$id}) 
			ON CREATE SET 
				i:Invoice_%s,
				i.source=$source,
				i.sourceOfTruth=$sourceOfTruth,
				i.appSource=$appSource,
				i.createdAt=$createdAt,
				i.updatedAt=$updatedAt,
				i.dryRun=$dryRun,
				i.preview=$preview,
				i.offCycle=$offCycle,
				i.postpaid=$postpaid,
				i.number=$number,
				i.periodStartDate=$periodStartDate,
				i.periodEndDate=$periodEndDate,
				i.dueDate=$dueDate,
				i.issuedDate=$issuedDate,
				i.currency=$currency,
				i.amount=$amount,
				i.vat=$vat,
				i.totalAmount=$totalAmount,
				i.repositoryFileId=$repositoryFileId,
				i.status=$status,
				i.note=$note,
				i.customerEmail=$customerEmail,
				i.paymentLink=$paymentLink,
				i.billingCycleInMonths=$billingCycleInMonths
			WITH c, i 
			MERGE (c)-[:HAS_INVOICE]->(i) 
				`, tenant)

	params := map[string]any{
		"id":                   invoiceId,
		"contractId":           contractId,
		"tenant":               tenant,
		"source":               invoice.Source,
		"sourceOfTruth":        invoice.SourceOfTruth,
		"appSource":            invoice.AppSource,
		"createdAt":            invoice.CreatedAt,
		"updatedAt":            invoice.UpdatedAt,
		"dryRun":               invoice.DryRun,
		"offCycle":             invoice.OffCycle,
		"postpaid":             invoice.Postpaid,
		"preview":              invoice.Preview,
		"number":               invoice.Number,
		"periodStartDate":      utils.ToNeo4jDateAsAny(&invoice.PeriodStartDate),
		"periodEndDate":        utils.ToNeo4jDateAsAny(&invoice.PeriodEndDate),
		"dueDate":              utils.ToNeo4jDateAsAny(&invoice.DueDate),
		"issuedDate":           utils.ToDate(invoice.IssuedDate),
		"currency":             invoice.Currency,
		"amount":               invoice.Amount,
		"vat":                  invoice.Vat,
		"totalAmount":          invoice.TotalAmount,
		"repositoryFileId":     invoice.RepositoryFileId,
		"status":               invoice.Status.String(),
		"note":                 invoice.Note,
		"customerEmail":        invoice.Customer.Email,
		"paymentLink":          invoice.PaymentDetails.PaymentLink,
		"billingCycleInMonths": invoice.BillingCycleInMonths,
	}

	ExecuteWriteQuery(ctx, driver, query, params)
	return invoiceId
}

func CreateInvoiceLine(ctx context.Context, driver *neo4j.DriverWithContext, tenant, invoiceId string, invoiceLine entity.InvoiceLineEntity) string {
	invoiceLineId := utils.NewUUIDIfEmpty(invoiceLine.Id)
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract)-[:HAS_INVOICE]->(i:Invoice {id:$invoiceId})
				MERGE (i)-[:HAS_INVOICE_LINE]->(il:InvoiceLine {id:$id})
				ON CREATE SET  
					il:InvoiceLine_%s,
					il.source=$source,
					il.sourceOfTruth=$sourceOfTruth,
					il.appSource=$appSource,
					il.createdAt=$createdAt,
					il.updatedAt=$updatedAt,
					il.name=$name,
					il.price=$price,
					il.quantity=$quantity,
					il.amount=$amount,
					il.vat=$vat,
					il.totalAmount=$totalAmount,
					il.billedType=$billedType
				`, tenant)

	params := map[string]any{
		"id":            invoiceLineId,
		"invoiceId":     invoiceId,
		"tenant":        tenant,
		"source":        invoiceLine.Source,
		"sourceOfTruth": invoiceLine.SourceOfTruth,
		"appSource":     invoiceLine.AppSource,
		"createdAt":     invoiceLine.CreatedAt,
		"updatedAt":     invoiceLine.UpdatedAt,
		"name":          invoiceLine.Name,
		"price":         invoiceLine.Price,
		"quantity":      invoiceLine.Quantity,
		"amount":        invoiceLine.Amount,
		"vat":           invoiceLine.Vat,
		"totalAmount":   invoiceLine.TotalAmount,
		"billedType":    invoiceLine.BilledType.String(),
	}

	ExecuteWriteQuery(ctx, driver, query, params)
	return invoiceLineId
}

func MarkInvoicingStarted(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contractId string, invoicingStartedAt time.Time) {
	query := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(ct:Contract_%s {id:$contractId})
				SET ct.techInvoicingStartedAt=$invoicingStartedAt
				`, tenant)

	params := map[string]any{
		"tenant":             tenant,
		"contractId":         contractId,
		"invoicingStartedAt": invoicingStartedAt,
	}

	ExecuteWriteQuery(ctx, driver, query, params)
}

func CreateTag(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, tagEntity entity.TagEntity) string {
	tagId := utils.NewUUIDIfEmpty(tagEntity.Id)
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$id})
			ON CREATE SET 
				tag:Tag_%s,
				tag.name=$name, 
				tag.source=$source, 
				tag.createdAt=$createdAt, 
				tag.updatedAt=$updatedAt`, tenant)
	params := map[string]any{
		"tenant":    tenant,
		"id":        tagId,
		"name":      tagEntity.Name,
		"createdAt": tagEntity.CreatedAt,
		"updatedAt": tagEntity.UpdatedAt,
		"source":    tagEntity.Source,
	}
	ExecuteWriteQuery(ctx, driver, query, params)
	return tagId
}

func CreateReminder(ctx context.Context, driver *neo4j.DriverWithContext, tenant, userId, orgId string, createdAt time.Time, reminderEntity entity.ReminderEntity) string {
	reminderId := utils.NewUUIDIfEmpty(reminderEntity.Id)
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
				MERGE (t)<-[:REMINDER_BELONGS_TO_TENANT]-(r:Reminder {id:$id})
				ON CREATE SET  
					r:Reminder_%s,
					r.createdAt=$createdAt,
					r.updatedAt=$createdAt,	
					r.source=$source,
					r.sourceOfTruth=$source,
					r.appSource=$appSource,
					r.content=$content,	
					r.dueDate=$dueDate,
					r.dismissed=$dismissed
					
				WITH t, r	
			
				MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId})
				MERGE (u)<-[:REMINDER_BELONGS_TO_USER]-(r)

				WITH t, r
				MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$organizationId})
				MERGE (o)<-[:REMINDER_BELONGS_TO_ORGANIZATION]-(r)`, tenant)
	params := map[string]interface{}{
		"tenant":         tenant,
		"id":             reminderId,
		"userId":         userId,
		"organizationId": orgId,
		"content":        reminderEntity.Content,
		"source":         reminderEntity.Source,
		"appSource":      reminderEntity.AppSource,
		"createdAt":      createdAt,
		"dueDate":        reminderEntity.DueDate,
		"dismissed":      reminderEntity.Dismissed,
	}
	ExecuteWriteQuery(ctx, driver, query, params)
	return reminderId
}

func CreateBankAccount(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, bankAccount entity.BankAccountEntity) string {
	accountId := utils.NewUUIDIfEmpty(bankAccount.Id)
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
				MERGE (t)-[:HAS_BANK_ACCOUNT]->(ba:BankAccount {id:$accountId})
				ON CREATE SET
					ba:BankAccount_%s,
					ba.createdAt=$createdAt,
					ba.updatedAt=$updatedAt,
					ba.bankName=$bankName,
					ba.currency=$currency,
					ba.bankTransferEnabled=$bankTransferEnabled,
					ba.allowInternational=$allowInternational,
					ba.accountNumber=$accountNumber,
					ba.iban=$iban,
					ba.bic=$bic,
					ba.sortCode=$sortCode,
					ba.routingNumber=$routingNumber,
					ba.otherDetails=$otherDetails
`, tenant)
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":              tenant,
		"accountId":           accountId,
		"createdAt":           bankAccount.CreatedAt,
		"updatedAt":           bankAccount.UpdatedAt,
		"bankName":            bankAccount.BankName,
		"currency":            bankAccount.Currency.String(),
		"bankTransferEnabled": bankAccount.BankTransferEnabled,
		"allowInternational":  bankAccount.AllowInternational,
		"accountNumber":       bankAccount.AccountNumber,
		"iban":                bankAccount.Iban,
		"bic":                 bankAccount.Bic,
		"sortCode":            bankAccount.SortCode,
		"routingNumber":       bankAccount.RoutingNumber,
		"otherDetails":        bankAccount.OtherDetails,
	})
	return accountId
}

func CreateExternalSystem(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, externalSystem entity.ExternalSystemEntity) {
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
				MERGE (t)<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(es:ExternalSystem {id:$externalSystemId})
				ON CREATE SET
					es:ExternalSystem_%s,
					es.stripePaymentMethodTypes=$stripePaymentMethodTypes
`, tenant)
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":                   tenant,
		"externalSystemId":         externalSystem.ExternalSystemId.String(),
		"stripePaymentMethodTypes": externalSystem.Stripe.PaymentMethodTypes,
	})
}

func CreateContact(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, contact entity.ContactEntity) string {
	contactId := utils.NewUUIDIfEmpty(contact.Id)
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
				MERGE (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})
				ON CREATE SET
					c:Contact_%s,
					c.createdAt=$createdAt,
					c.updatedAt=$updatedAt,
					c.source=$source,
					c.appSource=$appSource,
					c.firstName=$firstName,
					c.lastName=$lastName,
					c.prefix=$prefix,
					c.name=$name,
					c.hide=$hide
`, tenant)
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":    tenant,
		"contactId": contactId,
		"createdAt": contact.CreatedAt,
		"updatedAt": contact.UpdatedAt,
		"source":    contact.Source,
		"appSource": contact.AppSource,
		"firstName": contact.FirstName,
		"lastName":  contact.LastName,
		"prefix":    contact.Prefix,
		"name":      contact.Name,
		"hide":      contact.Hide,
	})
	return contactId
}

func CreateSocial(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, social entity.SocialEntity) string {
	socialId := utils.NewUUIDIfEmpty(social.Id)
	query := fmt.Sprintf(`MERGE (s:Social:Social_%s {id: $id})
				SET s.url=$url,
					s.alias=$alias,
					s.followersCount=$followersCount,
					s.externalId=$externalId
				`, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"id":             socialId,
		"url":            social.Url,
		"alias":          social.Alias,
		"followersCount": social.FollowersCount,
		"externalId":     social.ExternalId,
	})
	return socialId
}

func UserOwnsOrganization(ctx context.Context, driver *neo4j.DriverWithContext, userId, organizationId string) {
	query := `MATCH (o:Organization {id:$organizationId}),
			        (u:User {id:$userId})
			MERGE (u)-[:OWNS]->(o)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"organizationId": organizationId,
		"userId":         userId,
	})
}

func LinkDomainToOrganization(ctx context.Context, driver *neo4j.DriverWithContext, organizationId, domain string) {
	query := ` MERGE (d:Domain {domain:$domain})
			ON CREATE SET
				d.id=randomUUID(),
				d.source="test",
				d.appSource="test",
				d.createdAt=$now,
				d.updatedAt=$now
			WITH d
			MATCH (o:Organization {id:$organizationId})
			MERGE (o)-[:HAS_DOMAIN]->(d)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"organizationId": organizationId,
		"domain":         domain,
		"now":            utils.Now(),
	})
}

func CreateCommentForIssue(ctx context.Context, driver *neo4j.DriverWithContext, tenant, issueId string, comment entity.CommentEntity) string {
	commentId := CreateComment(ctx, driver, tenant, comment)
	LinkNodes(ctx, driver, commentId, issueId, "COMMENTED")
	return commentId
}

func CreateComment(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, comment entity.CommentEntity) string {
	commentId := utils.NewUUIDIfEmpty(comment.Id)
	query := fmt.Sprintf(`
			  MERGE (c:Comment {id:$id})
				ON CREATE SET c:Comment_%s,
					c.content=$content,
					c.contentType=$contentType,
					c.createdAt=$createdAt
				`, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"id":          commentId,
		"content":     comment.Content,
		"contentType": comment.ContentType,
		"createdAt":   comment.CreatedAt,
	})
	return commentId
}

func CreateInteractionEvent(ctx context.Context, driver *neo4j.DriverWithContext, tenant, identifier, content, contentType string, channel string, createdAt time.Time) string {
	return CreateInteractionEventFromEntity(ctx, driver, tenant, entity.InteractionEventEntity{
		Identifier:  identifier,
		Content:     content,
		ContentType: contentType,
		Channel:     channel,
		CreatedAt:   createdAt,
	})
}

func CreateInteractionEventFromEntity(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, ie entity.InteractionEventEntity) string {
	var interactionEventId, _ = uuid.NewRandom()

	query := "MERGE (ie:InteractionEvent {id:$id})" +
		" ON CREATE SET " +
		"	ie.content=$content, " +
		"	ie.createdAt=$createdAt, " +
		"	ie.channel=$channel, " +
		"	ie.contentType=$contentType, " +
		"	ie.source=$source, " +
		"   ie.hide=$hide, " +
		"	ie.sourceOfTruth=$sourceOfTruth, " +
		"	ie.appSource=$appSource," +
		"	ie:InteractionEvent_%s, ie:TimelineEvent, ie:TimelineEvent_%s," +
		"   ie.identifier=$identifier"
	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant, tenant), map[string]any{
		"id":            interactionEventId.String(),
		"content":       ie.Content,
		"contentType":   ie.ContentType,
		"channel":       ie.Channel,
		"createdAt":     ie.CreatedAt,
		"source":        "openline",
		"sourceOfTruth": "openline",
		"appSource":     "test",
		"identifier":    ie.Identifier,
		"hide":          ie.Hide,
	})
	return interactionEventId.String()
}

func CreateInteractionSession(ctx context.Context, driver *neo4j.DriverWithContext, tenant, identifier, name, sessionType, status, channel string, createdAt time.Time, inTimeline bool) string {
	var interactionSessionId, _ = uuid.NewRandom()

	query := "MERGE (is:InteractionSession {id:$id})" +
		" ON CREATE SET " +
		"	is.createdAt=$createdAt, " +
		"	is.updatedAt=$updatedAt, " +
		"	is.name=$name, " +
		"	is.type=$type, " +
		"	is.channel=$channel, " +
		"	is.status=$status, " +
		"	is.source=$source, " +
		"	is.sourceOfTruth=$sourceOfTruth, " +
		"	is.appSource=$appSource," +
		"   is.identifier=$identifier, " +
		"	is:InteractionSession_%s"

	resolvedQuery := ""
	if inTimeline {
		query += ", is:TimelineEvent, is:TimelineEvent_%s"

		resolvedQuery = fmt.Sprintf(query, tenant, tenant)
	} else {
		resolvedQuery = fmt.Sprintf(query, tenant)
	}
	ExecuteWriteQuery(ctx, driver, resolvedQuery, map[string]any{
		"id":            interactionSessionId.String(),
		"name":          name,
		"type":          sessionType,
		"channel":       channel,
		"status":        status,
		"createdAt":     createdAt,
		"updatedAt":     createdAt.Add(time.Duration(10) * time.Minute),
		"source":        "openline",
		"sourceOfTruth": "openline",
		"appSource":     "test",
		"identifier":    identifier,
	})
	return interactionSessionId.String()
}
