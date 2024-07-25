import { ColumnViewType } from '@graphql/types';

type InvoicesColumnType =
  | ColumnViewType.InvoicesAmount
  | ColumnViewType.InvoicesBillingCycle
  | ColumnViewType.InvoicesContract
  | ColumnViewType.InvoicesDueDate
  | ColumnViewType.InvoicesIssueDatePast
  | ColumnViewType.InvoicesInvoicePreview
  | ColumnViewType.InvoicesIssueDate
  | ColumnViewType.InvoicesInvoiceStatus
  | ColumnViewType.InvoicesPaymentStatus
  | ColumnViewType.InvoicesInvoiceNumber;

export const invoicesOptionsMap: Record<InvoicesColumnType | string, string> = {
  [ColumnViewType.InvoicesAmount]: 'Amount',
  [ColumnViewType.InvoicesBillingCycle]: 'Billing Cycle',
  [ColumnViewType.InvoicesContract]: 'Contract',
  [ColumnViewType.InvoicesDueDate]: 'Due Date',
  [ColumnViewType.InvoicesInvoicePreview]: 'Invoice Preview',
  [ColumnViewType.InvoicesInvoiceNumber]: 'Invoice',
  [ColumnViewType.InvoicesIssueDate]: 'Issue Date',
  [ColumnViewType.InvoicesIssueDatePast]: 'Created At',
  [ColumnViewType.InvoicesInvoiceStatus]: 'Invoice Status',
  [ColumnViewType.InvoicesPaymentStatus]: 'Payment Status',
};

export const contactsOptionsMap: Record<InvoicesColumnType | string, string> = {
  [ColumnViewType.ContactsOrganization]: 'Organization',
  [ColumnViewType.ContactsName]: 'Name',
  [ColumnViewType.ContactsLinkedin]: 'LinkedIn',
  [ColumnViewType.ContactsCity]: 'City',
  [ColumnViewType.ContactsPersona]: 'Persona',
  [ColumnViewType.ContactsLastInteraction]: 'Last Interaction',
  [ColumnViewType.ContactsEmails]: 'Email',
  [ColumnViewType.ContactsPhoneNumbers]: 'Phone Number',
  [ColumnViewType.ContactsAvatar]: 'Avatar',
  [ColumnViewType.ContactsLanguages]: 'Languages',
  [ColumnViewType.ContactsTags]: 'Tags',
  [ColumnViewType.ContactsExperience]: 'Experience',
  [ColumnViewType.ContactsSchools]: 'Schools',
  [ColumnViewType.ContactsTimeInCurrentRole]: 'Time In Current Role',
  [ColumnViewType.ContactsJobTitle]: 'Job Title',
  [ColumnViewType.ContactsCountry]: 'Country',
  [ColumnViewType.ContactsSkills]: 'Skills',
  [ColumnViewType.ContactsLinkedinFollowerCount]: 'Linkedin Followers',
  [ColumnViewType.ContactsConnections]: 'Connected To',
  [ColumnViewType.ContactsRegion]: 'Region',
};

export const invoicesHelperTextMap: Record<
  InvoicesColumnType | string,
  string
> = {
  [ColumnViewType.InvoicesAmount]: 'E.g. $6,450',
  [ColumnViewType.InvoicesBillingCycle]: 'E.g. Monthly',
  [ColumnViewType.InvoicesContract]: 'E.g. Pile Contract',
  [ColumnViewType.InvoicesDueDate]: 'E.g. 15 Aug 2019',
  [ColumnViewType.InvoicesInvoicePreview]: 'E.g. RKD-04025',
  [ColumnViewType.InvoicesInvoiceNumber]: 'E.g. RKD-04025',
  [ColumnViewType.InvoicesIssueDate]: 'E.g. 15 Aug 2019',
  [ColumnViewType.InvoicesIssueDatePast]: 'E.g. 15 Jun 2019',
  [ColumnViewType.InvoicesInvoiceStatus]: 'E.g. Scheduled',
  [ColumnViewType.InvoicesPaymentStatus]: 'E.g. Paid',
};

type RenewalsColumnType =
  | ColumnViewType.RenewalsAvatar
  | ColumnViewType.RenewalsForecastArr
  | ColumnViewType.RenewalsLastTouchpoint
  | ColumnViewType.RenewalsName
  | ColumnViewType.RenewalsOwner
  | ColumnViewType.RenewalsRenewalDate
  | ColumnViewType.RenewalsRenewalLikelihood;

export const renewalsOptionsMap: Record<RenewalsColumnType | string, string> = {
  [ColumnViewType.RenewalsAvatar]: 'Logo',
  [ColumnViewType.RenewalsForecastArr]: 'ARR Forecast',
  [ColumnViewType.RenewalsLastTouchpoint]: 'Last Touchpoint',
  [ColumnViewType.RenewalsName]: 'Name',
  [ColumnViewType.RenewalsOwner]: 'Owner',
  [ColumnViewType.RenewalsRenewalDate]: 'Next Renewal',
  [ColumnViewType.RenewalsRenewalLikelihood]: 'Health',
};

export const renewalsHelperTextMap: Record<
  RenewalsColumnType | string,
  string
> = {
  [ColumnViewType.RenewalsAvatar]: 'E.g. Logo',
  [ColumnViewType.RenewalsForecastArr]: 'E.g. $6,450',
  [ColumnViewType.RenewalsLastTouchpoint]: 'E.g. Issue updated',
  [ColumnViewType.RenewalsName]: 'E.g. Pile Contract',
  [ColumnViewType.RenewalsOwner]: 'E.g. Howard Hu',
  [ColumnViewType.RenewalsRenewalDate]: 'E.g. 1 month',
  [ColumnViewType.RenewalsRenewalLikelihood]: 'E.g. High',
};

type OrganizationsColumnType =
  | ColumnViewType.OrganizationsAvatar
  | ColumnViewType.OrganizationsForecastArr
  | ColumnViewType.OrganizationsLastTouchpoint
  | ColumnViewType.OrganizationsName
  | ColumnViewType.OrganizationsOwner
  | ColumnViewType.OrganizationsOnboardingStatus
  | ColumnViewType.OrganizationsRelationship
  | ColumnViewType.OrganizationsRenewalLikelihood
  | ColumnViewType.OrganizationsRenewalDate
  | ColumnViewType.OrganizationsWebsite
  | ColumnViewType.OrganizationsChurnDate;

export const organizationsOptionsMap: Record<
  OrganizationsColumnType | string,
  string
> = {
  [ColumnViewType.OrganizationsAvatar]: 'Logo',
  [ColumnViewType.OrganizationsForecastArr]: 'ARR Forecast',
  [ColumnViewType.OrganizationsLastTouchpoint]: 'Last Touchpoint',
  [ColumnViewType.OrganizationsName]: 'Organization',
  [ColumnViewType.OrganizationsOwner]: 'Owner',
  [ColumnViewType.OrganizationsOnboardingStatus]: 'Onboarding',
  [ColumnViewType.OrganizationsRelationship]: 'Relationship',
  [ColumnViewType.OrganizationsRenewalLikelihood]: 'Health',
  [ColumnViewType.OrganizationsRenewalDate]: 'Next Renewal',
  [ColumnViewType.OrganizationsWebsite]: 'Website',
  [ColumnViewType.OrganizationsLeadSource]: 'Source',
  [ColumnViewType.OrganizationsSocials]: 'LinkedIn',
  [ColumnViewType.OrganizationsCreatedDate]: 'Created Date',
  [ColumnViewType.OrganizationsEmployeeCount]: 'Employee Count',
  [ColumnViewType.OrganizationsYearFounded]: 'Year Founded',
  [ColumnViewType.OrganizationsLastTouchpointDate]: 'Last Interacted',
  [ColumnViewType.OrganizationsChurnDate]: 'Churn Date',
  [ColumnViewType.OrganizationsLtv]: 'LTV',
  [ColumnViewType.OrganizationsIndustry]: 'Industry',
  [ColumnViewType.OrganizationsTags]: 'Tags',
  [ColumnViewType.OrganizationsLinkedinFollowerCount]: 'Linkedin Followers',
  [ColumnViewType.OrganizationsStage]: 'Stage',
  [ColumnViewType.OrganizationsCity]: 'Headquarters',
  [ColumnViewType.OrganizationsIsPublic]: 'Ownership Type',
  [ColumnViewType.OrganizationsContactCount]: 'Contacts',
};

export const organizationsHelperTextMap: Record<
  OrganizationsColumnType | string,
  string
> = {
  [ColumnViewType.OrganizationsAvatar]: 'E.g. Logo',
  [ColumnViewType.OrganizationsForecastArr]: 'E.g. $6,450',
  [ColumnViewType.OrganizationsLastTouchpoint]: 'E.g. Issue updated',
  [ColumnViewType.OrganizationsName]: 'E.g. Pile Contract',
  [ColumnViewType.OrganizationsOwner]: 'E.g. Howard Hu',
  [ColumnViewType.OrganizationsOnboardingStatus]: 'E.g. Onboarding',
  [ColumnViewType.OrganizationsRelationship]: 'E.g. Customer',
  [ColumnViewType.OrganizationsRenewalLikelihood]: 'E.g. High',
  [ColumnViewType.OrganizationsRenewalDate]: 'E.g. 1 month',
  [ColumnViewType.OrganizationsWebsite]: 'E.g. www.pile.com',
  [ColumnViewType.OrganizationsLeadSource]: 'E.g. Newsletter',
  [ColumnViewType.OrganizationsSocials]: 'E.g. /acmecorp',
  [ColumnViewType.OrganizationsCreatedDate]: 'E.g. 28 Mar 2019',
  [ColumnViewType.OrganizationsEmployeeCount]: 'E.g. 192',
  [ColumnViewType.OrganizationsYearFounded]: 'E.g. 2017',
  [ColumnViewType.OrganizationsLastTouchpointDate]: 'E.g. 16 Sep 2025',
  [ColumnViewType.OrganizationsCity]: 'E.g. Cape Town',
  [ColumnViewType.OrganizationsIsPublic]: 'E.g. Private',
  [ColumnViewType.OrganizationsLinkedinFollowerCount]: 'E.g. 15,930',
  [ColumnViewType.OrganizationsTags]: 'E.g. Solo RevOps',
  [ColumnViewType.OrganizationsContactCount]: 'E.g. 5',
  [ColumnViewType.OrganizationsIndustry]: 'E.g. Software',
  [ColumnViewType.OrganizationsStage]: 'E.g. Lead',
  [ColumnViewType.OrganizationsChurnDate]: 'E.g. 15 Aug 2024',
  [ColumnViewType.OrganizationsLtv]: 'E.g. $109,280',
};

export const contactsHelperTextMap: Record<string, string> = {
  [ColumnViewType.ContactsOrganization]: 'E.g. CustomerOs',
  [ColumnViewType.ContactsName]: 'E.g. Jane Doe',
  [ColumnViewType.ContactsLinkedin]: 'E.g. /janedoe',
  [ColumnViewType.ContactsCity]: 'E.g. Cape Town',
  [ColumnViewType.ContactsPersona]: 'E.g. Champion',
  [ColumnViewType.ContactsLastInteraction]: 'E.g. 16 Sep 2025',
  [ColumnViewType.ContactsEmails]: 'E.g. john.doe@acme.com',
  [ColumnViewType.ContactsPhoneNumbers]: 'E.g. (907) 834-2765',
  [ColumnViewType.ContactsLanguages]: 'E.g. English',
  [ColumnViewType.ContactsTimeInCurrentRole]: 'E.g. 2 years',
  [ColumnViewType.ContactsJobTitle]: 'E.g. CTO',
  [ColumnViewType.ContactsCountry]: 'E.g. South Africa',
  [ColumnViewType.ContactsLinkedinFollowerCount]: 'E.g. 15,930',
  [ColumnViewType.ContactsConnections]: 'E.g. Amy Liu',
  [ColumnViewType.ContactsSkills]: 'E.g. Data privacy',
  [ColumnViewType.ContactsSchools]: 'E.g. Stanford University',
  [ColumnViewType.ContactsExperience]: 'E.g. 4 yrs',
  [ColumnViewType.ContactsRegion]: 'E.g. California',
};
