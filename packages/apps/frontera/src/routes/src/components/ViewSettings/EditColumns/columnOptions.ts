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
  | ColumnViewType.InvoicesInvoiceNumber;

export const invoicesOptionsMap: Record<InvoicesColumnType | string, string> = {
  [ColumnViewType.InvoicesAmount]: 'Amount',
  [ColumnViewType.InvoicesBillingCycle]: 'Billing Cycle',
  [ColumnViewType.InvoicesContract]: 'Contract',
  [ColumnViewType.InvoicesDueDate]: 'Due Date',
  [ColumnViewType.InvoicesInvoicePreview]: 'Upcoming Invoices',
  [ColumnViewType.InvoicesInvoiceNumber]: 'Invoice',
  [ColumnViewType.InvoicesIssueDate]: 'Issue Date',
  [ColumnViewType.InvoicesIssueDatePast]: 'Created At',
  [ColumnViewType.InvoicesInvoiceStatus]: 'Invoice Status',
  [ColumnViewType.InvoicesOrganization]: 'Organization Name',
};

export const contactsOptionsMap: Record<InvoicesColumnType | string, string> = {
  [ColumnViewType.ContactsOrganization]: 'Organization',
  [ColumnViewType.ContactsName]: 'Name',
  [ColumnViewType.ContactsLinkedin]: 'LinkedIn',
  [ColumnViewType.ContactsCity]: 'City',
  [ColumnViewType.ContactsPersona]: 'Persona',
  [ColumnViewType.ContactsLastInteraction]: 'Last Interaction',
  [ColumnViewType.ContactsPhoneNumbers]: 'Mobile Number',
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
  [ColumnViewType.ContactsConnections]: 'LinkedIn Connections',
  [ColumnViewType.ContactsRegion]: 'Region',
  [ColumnViewType.ContactsFlows]: 'Current Flows',
  [ColumnViewType.ContactsPrimaryEmail]: 'Primary Email',
  [ColumnViewType.ContactsFlowStatus]: 'Status in Flow',
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
  | ColumnViewType.OrganizationsChurnDate
  | ColumnViewType.OrganizationsParentOrganization;

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
  [ColumnViewType.OrganizationsRenewalDate]: 'Renewal Date',
  [ColumnViewType.OrganizationsWebsite]: 'Website',
  [ColumnViewType.OrganizationsLeadSource]: 'Source',
  [ColumnViewType.OrganizationsSocials]: 'LinkedIn',
  [ColumnViewType.OrganizationsCreatedDate]: 'Created Date',
  [ColumnViewType.OrganizationsEmployeeCount]: 'Employees',
  [ColumnViewType.OrganizationsYearFounded]: 'Founded',
  [ColumnViewType.OrganizationsLastTouchpointDate]: 'Last Interacted',
  [ColumnViewType.OrganizationsChurnDate]: 'Churn Date',
  [ColumnViewType.OrganizationsLtv]: 'LTV',
  [ColumnViewType.OrganizationsIndustry]: 'Industry',
  [ColumnViewType.OrganizationsTags]: 'Tags',
  [ColumnViewType.OrganizationsLinkedinFollowerCount]: 'Linkedin Followers',
  [ColumnViewType.OrganizationsStage]: 'Stage',
  [ColumnViewType.OrganizationsCity]: 'City',
  [ColumnViewType.OrganizationsIsPublic]: 'Ownership Type',
  [ColumnViewType.OrganizationsContactCount]: 'Contacts',
  [ColumnViewType.OrganizationsHeadquarters]: 'Country',
  [ColumnViewType.OrganizationsParentOrganization]: 'Parent Org',
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
  [ColumnViewType.OrganizationsRenewalDate]: 'E.g. 3 Aug 2027',
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
  [ColumnViewType.OrganizationsHeadquarters]: 'E.g. Germany',
  [ColumnViewType.OrganizationsParentOrganization]: 'E.g. Alphabet',
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
  [ColumnViewType.ContactsFlows]: 'E.g. Education',
  [ColumnViewType.ContactsPersonalEmails]: 'E.g. steph@convoy.com',
  [ColumnViewType.ContactsFlowStatus]: 'E.g. Completed',
};

export const contractsMap: Record<string, string> = {
  [ColumnViewType.ContractsName]: 'Contract Name',
  [ColumnViewType.ContractsPeriod]: 'Period',
  [ColumnViewType.ContractsEnded]: 'Ended',
  [ColumnViewType.ContractsCurrency]: 'Currency',
  [ColumnViewType.ContractsStatus]: 'Status',
  [ColumnViewType.ContractsRenewal]: 'Renewal',
  [ColumnViewType.ContractsLtv]: 'LTV',
};

export const contractsHelperTextMap: Record<string, string> = {
  [ColumnViewType.ContractsName]: 'E.g. CustomerOs contract',
  [ColumnViewType.ContractsPeriod]: 'E.g. Monthly',
  [ColumnViewType.ContractsEnded]: 'E.g. 19 Jun 2022',
  [ColumnViewType.ContractsCurrency]: 'E.g. USD',
  [ColumnViewType.ContractsStatus]: 'E.g. Live',
  [ColumnViewType.ContractsRenewal]: 'E.g. Auto-renewing',
  [ColumnViewType.ContractsLtv]: 'E.g. $730,800',
  [ColumnViewType.ContractsOwner]: 'E.g. Sam Douglas',
  [ColumnViewType.ContractsHealth]: 'E.g. High',
  [ColumnViewType.ContractsRenewalDate]: 'E.g 12 Oct 2026',
  [ColumnViewType.ContractsForecastArr]: 'E.g. $120,930',
};

export const opportunitiesMap: Record<string, string> = {
  [ColumnViewType.OpportunitiesName]: 'Name',
  [ColumnViewType.OpportunitiesOrganization]: 'Organization',
  [ColumnViewType.OpportunitiesStage]: 'Stage',
  [ColumnViewType.OpportunitiesTimeInStage]: 'Time in Stage',
  [ColumnViewType.OpportunitiesEstimatedArr]: 'Estimated ARR',
  [ColumnViewType.OpportunitiesCreatedDate]: 'Created',
  [ColumnViewType.OpportunitiesNextStep]: 'Next Step',
};
export const opportunitiesHelperTextMap: Record<string, string> = {
  [ColumnViewType.OpportunitiesName]: 'E.g. CustomerOs opportunity',
  [ColumnViewType.OpportunitiesOrganization]: 'E.g. CustomerOs',
  [ColumnViewType.OpportunitiesStage]: 'E.g. Identified',
  [ColumnViewType.OpportunitiesTimeInStage]: 'E.g. 6 days',
  [ColumnViewType.OpportunitiesEstimatedArr]: 'E.g. $30,000',
  [ColumnViewType.OpportunitiesCreatedDate]: 'E.g. 12 Oct 2023',
  [ColumnViewType.OpportunitiesNextStep]: 'E.g. Qualified',
  [ColumnViewType.OpportunitiesOwner]: 'E.g. Howard Hu',
};

export const flowsMap: Record<string, string> = {
  [ColumnViewType.FlowName]: 'Flow',
  [ColumnViewType.FlowActionName]: 'Status', // This is the actual status of the flow - Wrong naming of enum -> to be refactored
  [ColumnViewType.FlowTotalCount]: 'Total',
  [ColumnViewType.FlowOnHoldCount]: 'Blocked',
  [ColumnViewType.FlowReadyCount]: 'Ready',
  [ColumnViewType.FlowScheduledCount]: 'Scheduled',
  [ColumnViewType.FlowInProgressCount]: 'In progress',
  [ColumnViewType.FlowCompletedCount]: 'Completed',
  [ColumnViewType.FlowGoalAchievedCount]: 'Goal Achieved',
};
export const flowsHelperTextMap: Record<string, string> = {
  [ColumnViewType.FlowName]: 'E.g. Aerospace CTO',
  [ColumnViewType.FlowActionName]: 'E.g. Live',
  [ColumnViewType.FlowTotalCount]: 'E.g. 125',
  [ColumnViewType.FlowOnHoldCount]: 'E.g. 21',
  [ColumnViewType.FlowReadyCount]: 'E.g. 12',
  [ColumnViewType.FlowScheduledCount]: 'E.g. 23',
  [ColumnViewType.FlowInProgressCount]: 'E.g. 34',
  [ColumnViewType.FlowCompletedCount]: 'E.g. 78',
  [ColumnViewType.FlowGoalAchievedCount]: 'E.g. 47',
};
