export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = {
  [K in keyof T]: T[K];
};
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]?: Maybe<T[SubKey]>;
};
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]: Maybe<T[SubKey]>;
};
export type MakeEmpty<
  T extends { [key: string]: unknown },
  K extends keyof T,
> = { [_ in K]?: never };
export type Incremental<T> =
  | T
  | {
      [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never;
    };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string };
  String: { input: string; output: string };
  Boolean: { input: boolean; output: boolean };
  Int: { input: number; output: number };
  Float: { input: number; output: number };
  Any: { input: any; output: any };
  Int64: { input: any; output: any };
  Time: { input: any; output: any };
};

export type Action = {
  __typename?: 'Action';
  actionType: ActionType;
  appSource: Scalars['String']['output'];
  content?: Maybe<Scalars['String']['output']>;
  createdAt: Scalars['Time']['output'];
  createdBy?: Maybe<User>;
  id: Scalars['ID']['output'];
  metadata?: Maybe<Scalars['String']['output']>;
  source: DataSource;
};

export type ActionItem = {
  __typename?: 'ActionItem';
  appSource: Scalars['String']['output'];
  content: Scalars['String']['output'];
  createdAt: Scalars['Time']['output'];
  id: Scalars['ID']['output'];
  source: DataSource;
};

export type ActionResponse = {
  __typename?: 'ActionResponse';
  accepted: Scalars['Boolean']['output'];
};

export enum ActionType {
  ContractRenewed = 'CONTRACT_RENEWED',
  ContractStatusUpdated = 'CONTRACT_STATUS_UPDATED',
  Created = 'CREATED',
  InteractionEventRead = 'INTERACTION_EVENT_READ',
  InvoiceIssued = 'INVOICE_ISSUED',
  InvoiceOverdue = 'INVOICE_OVERDUE',
  InvoicePaid = 'INVOICE_PAID',
  InvoiceSent = 'INVOICE_SENT',
  InvoiceVoided = 'INVOICE_VOIDED',
  OnboardingStatusChanged = 'ONBOARDING_STATUS_CHANGED',
  RenewalForecastUpdated = 'RENEWAL_FORECAST_UPDATED',
  RenewalLikelihoodUpdated = 'RENEWAL_LIKELIHOOD_UPDATED',
  ServiceLineItemBilledTypeOnceCreated = 'SERVICE_LINE_ITEM_BILLED_TYPE_ONCE_CREATED',
  ServiceLineItemBilledTypeRecurringCreated = 'SERVICE_LINE_ITEM_BILLED_TYPE_RECURRING_CREATED',
  /** Deprecated */
  ServiceLineItemBilledTypeUpdated = 'SERVICE_LINE_ITEM_BILLED_TYPE_UPDATED',
  ServiceLineItemBilledTypeUsageCreated = 'SERVICE_LINE_ITEM_BILLED_TYPE_USAGE_CREATED',
  ServiceLineItemPriceUpdated = 'SERVICE_LINE_ITEM_PRICE_UPDATED',
  ServiceLineItemQuantityUpdated = 'SERVICE_LINE_ITEM_QUANTITY_UPDATED',
  ServiceLineItemRemoved = 'SERVICE_LINE_ITEM_REMOVED',
}

export type AddTagInput = {
  entityId: Scalars['ID']['input'];
  entityType: EntityType;
  tag: TagIdOrNameInput;
};

export type Attachment = Node & {
  __typename?: 'Attachment';
  appSource: Scalars['String']['output'];
  basePath: Scalars['String']['output'];
  cdnUrl: Scalars['String']['output'];
  createdAt: Scalars['Time']['output'];
  fileName: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  mimeType: Scalars['String']['output'];
  size: Scalars['Int64']['output'];
  source: DataSource;
  sourceOfTruth: DataSource;
};

export type AttachmentInput = {
  appSource: Scalars['String']['input'];
  basePath: Scalars['String']['input'];
  cdnUrl: Scalars['String']['input'];
  createdAt?: InputMaybe<Scalars['Time']['input']>;
  fileName: Scalars['String']['input'];
  id?: InputMaybe<Scalars['ID']['input']>;
  mimeType: Scalars['String']['input'];
  size: Scalars['Int64']['input'];
};

export type BankAccount = MetadataInterface & {
  __typename?: 'BankAccount';
  accountNumber?: Maybe<Scalars['String']['output']>;
  allowInternational: Scalars['Boolean']['output'];
  bankName?: Maybe<Scalars['String']['output']>;
  bankTransferEnabled: Scalars['Boolean']['output'];
  bic?: Maybe<Scalars['String']['output']>;
  currency?: Maybe<Currency>;
  iban?: Maybe<Scalars['String']['output']>;
  metadata: Metadata;
  otherDetails?: Maybe<Scalars['String']['output']>;
  routingNumber?: Maybe<Scalars['String']['output']>;
  sortCode?: Maybe<Scalars['String']['output']>;
};

export type BankAccountCreateInput = {
  accountNumber?: InputMaybe<Scalars['String']['input']>;
  allowInternational?: InputMaybe<Scalars['Boolean']['input']>;
  bankName?: InputMaybe<Scalars['String']['input']>;
  bankTransferEnabled?: InputMaybe<Scalars['Boolean']['input']>;
  bic?: InputMaybe<Scalars['String']['input']>;
  currency?: InputMaybe<Currency>;
  iban?: InputMaybe<Scalars['String']['input']>;
  otherDetails?: InputMaybe<Scalars['String']['input']>;
  routingNumber?: InputMaybe<Scalars['String']['input']>;
  sortCode?: InputMaybe<Scalars['String']['input']>;
};

export type BankAccountUpdateInput = {
  accountNumber?: InputMaybe<Scalars['String']['input']>;
  allowInternational?: InputMaybe<Scalars['Boolean']['input']>;
  bankName?: InputMaybe<Scalars['String']['input']>;
  bankTransferEnabled?: InputMaybe<Scalars['Boolean']['input']>;
  bic?: InputMaybe<Scalars['String']['input']>;
  currency?: InputMaybe<Currency>;
  iban?: InputMaybe<Scalars['String']['input']>;
  id: Scalars['ID']['input'];
  otherDetails?: InputMaybe<Scalars['String']['input']>;
  routingNumber?: InputMaybe<Scalars['String']['input']>;
  sortCode?: InputMaybe<Scalars['String']['input']>;
};

export enum BilledType {
  Annually = 'ANNUALLY',
  Monthly = 'MONTHLY',
  /**
   * Deprecated
   * @deprecated MONTHLY will be used instead.
   */
  None = 'NONE',
  Once = 'ONCE',
  Quarterly = 'QUARTERLY',
  /**
   * Deprecated
   * @deprecated Not supported yet.
   */
  Usage = 'USAGE',
}

export type BillingDetails = {
  __typename?: 'BillingDetails';
  addressLine1?: Maybe<Scalars['String']['output']>;
  addressLine2?: Maybe<Scalars['String']['output']>;
  /** @deprecated Use billingCycleInMonths instead. */
  billingCycle?: Maybe<ContractBillingCycle>;
  billingCycleInMonths?: Maybe<Scalars['Int64']['output']>;
  billingEmail?: Maybe<Scalars['String']['output']>;
  billingEmailBCC?: Maybe<Array<Scalars['String']['output']>>;
  billingEmailCC?: Maybe<Array<Scalars['String']['output']>>;
  canPayWithBankTransfer?: Maybe<Scalars['Boolean']['output']>;
  canPayWithCard?: Maybe<Scalars['Boolean']['output']>;
  canPayWithDirectDebit?: Maybe<Scalars['Boolean']['output']>;
  check?: Maybe<Scalars['Boolean']['output']>;
  country?: Maybe<Scalars['String']['output']>;
  dueDays?: Maybe<Scalars['Int64']['output']>;
  invoiceNote?: Maybe<Scalars['String']['output']>;
  invoicingStarted?: Maybe<Scalars['Time']['output']>;
  locality?: Maybe<Scalars['String']['output']>;
  nextInvoicing?: Maybe<Scalars['Time']['output']>;
  organizationLegalName?: Maybe<Scalars['String']['output']>;
  payAutomatically?: Maybe<Scalars['Boolean']['output']>;
  payOnline?: Maybe<Scalars['Boolean']['output']>;
  postalCode?: Maybe<Scalars['String']['output']>;
  region?: Maybe<Scalars['String']['output']>;
};

export type BillingDetailsInput = {
  addressLine1?: InputMaybe<Scalars['String']['input']>;
  addressLine2?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated, use billingCycleInMonths instead. */
  billingCycle?: InputMaybe<ContractBillingCycle>;
  billingCycleInMonths?: InputMaybe<Scalars['Int64']['input']>;
  billingEmail?: InputMaybe<Scalars['String']['input']>;
  billingEmailBCC?: InputMaybe<Array<Scalars['String']['input']>>;
  billingEmailCC?: InputMaybe<Array<Scalars['String']['input']>>;
  canPayWithBankTransfer?: InputMaybe<Scalars['Boolean']['input']>;
  canPayWithCard?: InputMaybe<Scalars['Boolean']['input']>;
  canPayWithDirectDebit?: InputMaybe<Scalars['Boolean']['input']>;
  check?: InputMaybe<Scalars['Boolean']['input']>;
  country?: InputMaybe<Scalars['String']['input']>;
  dueDays?: InputMaybe<Scalars['Int64']['input']>;
  invoiceNote?: InputMaybe<Scalars['String']['input']>;
  invoicingStarted?: InputMaybe<Scalars['Time']['input']>;
  locality?: InputMaybe<Scalars['String']['input']>;
  organizationLegalName?: InputMaybe<Scalars['String']['input']>;
  payAutomatically?: InputMaybe<Scalars['Boolean']['input']>;
  payOnline?: InputMaybe<Scalars['Boolean']['input']>;
  postalCode?: InputMaybe<Scalars['String']['input']>;
  region?: InputMaybe<Scalars['String']['input']>;
};

export type BillingProfile = Node &
  SourceFields & {
    __typename?: 'BillingProfile';
    appSource: Scalars['String']['output'];
    createdAt: Scalars['Time']['output'];
    id: Scalars['ID']['output'];
    legalName: Scalars['String']['output'];
    source: DataSource;
    sourceOfTruth: DataSource;
    taxId: Scalars['String']['output'];
    updatedAt: Scalars['Time']['output'];
  };

export type BillingProfileInput = {
  createdAt?: InputMaybe<Scalars['Time']['input']>;
  legalName?: InputMaybe<Scalars['String']['input']>;
  organizationId: Scalars['ID']['input'];
  taxId?: InputMaybe<Scalars['String']['input']>;
};

export type BillingProfileLinkEmailInput = {
  billingProfileId: Scalars['ID']['input'];
  emailId: Scalars['ID']['input'];
  organizationId: Scalars['ID']['input'];
  primary?: InputMaybe<Scalars['Boolean']['input']>;
};

export type BillingProfileLinkLocationInput = {
  billingProfileId: Scalars['ID']['input'];
  locationId: Scalars['ID']['input'];
  organizationId: Scalars['ID']['input'];
};

export type BillingProfileUpdateInput = {
  billingProfileId: Scalars['ID']['input'];
  legalName?: InputMaybe<Scalars['String']['input']>;
  organizationId: Scalars['ID']['input'];
  taxId?: InputMaybe<Scalars['String']['input']>;
  updatedAt?: InputMaybe<Scalars['Time']['input']>;
};

/**
 * Describes the relationship a Contact has with a Organization.
 * **A `return` object**
 */
export type Calendar = {
  __typename?: 'Calendar';
  appSource: Scalars['String']['output'];
  calType: CalendarType;
  createdAt: Scalars['Time']['output'];
  id: Scalars['ID']['output'];
  link?: Maybe<Scalars['String']['output']>;
  primary: Scalars['Boolean']['output'];
  source: DataSource;
  sourceOfTruth: DataSource;
  updatedAt: Scalars['Time']['output'];
};

export enum CalendarType {
  Calcom = 'CALCOM',
  Google = 'GOOGLE',
}

export type ColumnView = {
  __typename?: 'ColumnView';
  columnId: Scalars['Int']['output'];
  columnType: ColumnViewType;
  filter: Scalars['String']['output'];
  name: Scalars['String']['output'];
  visible: Scalars['Boolean']['output'];
  width: Scalars['Int']['output'];
};

export type ColumnViewInput = {
  columnId: Scalars['Int']['input'];
  columnType: ColumnViewType;
  filter: Scalars['String']['input'];
  name: Scalars['String']['input'];
  visible: Scalars['Boolean']['input'];
  width: Scalars['Int']['input'];
};

export enum ColumnViewType {
  ContactsAvatar = 'CONTACTS_AVATAR',
  ContactsCity = 'CONTACTS_CITY',
  ContactsConnections = 'CONTACTS_CONNECTIONS',
  ContactsCountry = 'CONTACTS_COUNTRY',
  ContactsCreatedAt = 'CONTACTS_CREATED_AT',
  ContactsEmails = 'CONTACTS_EMAILS',
  ContactsExperience = 'CONTACTS_EXPERIENCE',
  ContactsFlows = 'CONTACTS_FLOWS',
  ContactsFlowStatus = 'CONTACTS_FLOW_STATUS',
  ContactsJobTitle = 'CONTACTS_JOB_TITLE',
  ContactsLanguages = 'CONTACTS_LANGUAGES',
  ContactsLastInteraction = 'CONTACTS_LAST_INTERACTION',
  ContactsLinkedin = 'CONTACTS_LINKEDIN',
  ContactsLinkedinFollowerCount = 'CONTACTS_LINKEDIN_FOLLOWER_COUNT',
  ContactsName = 'CONTACTS_NAME',
  ContactsOrganization = 'CONTACTS_ORGANIZATION',
  ContactsPersona = 'CONTACTS_PERSONA',
  ContactsPersonalEmails = 'CONTACTS_PERSONAL_EMAILS',
  ContactsPhoneNumbers = 'CONTACTS_PHONE_NUMBERS',
  ContactsPrimaryEmail = 'CONTACTS_PRIMARY_EMAIL',
  ContactsRegion = 'CONTACTS_REGION',
  ContactsSchools = 'CONTACTS_SCHOOLS',
  ContactsSkills = 'CONTACTS_SKILLS',
  ContactsTags = 'CONTACTS_TAGS',
  ContactsTimeInCurrentRole = 'CONTACTS_TIME_IN_CURRENT_ROLE',
  ContactsUpdatedAt = 'CONTACTS_UPDATED_AT',
  ContractsCurrency = 'CONTRACTS_CURRENCY',
  ContractsEnded = 'CONTRACTS_ENDED',
  ContractsForecastArr = 'CONTRACTS_FORECAST_ARR',
  ContractsHealth = 'CONTRACTS_HEALTH',
  ContractsLtv = 'CONTRACTS_LTV',
  ContractsName = 'CONTRACTS_NAME',
  ContractsOwner = 'CONTRACTS_OWNER',
  ContractsPeriod = 'CONTRACTS_PERIOD',
  ContractsRenewal = 'CONTRACTS_RENEWAL',
  ContractsRenewalDate = 'CONTRACTS_RENEWAL_DATE',
  ContractsStatus = 'CONTRACTS_STATUS',
  FlowActionName = 'FLOW_ACTION_NAME',
  FlowActionStatus = 'FLOW_ACTION_STATUS',
  FlowCompletedCount = 'FLOW_COMPLETED_COUNT',
  FlowGoalAchievedCount = 'FLOW_GOAL_ACHIEVED_COUNT',
  FlowInProgressCount = 'FLOW_IN_PROGRESS_COUNT',
  FlowName = 'FLOW_NAME',
  FlowOnHoldCount = 'FLOW_ON_HOLD_COUNT',
  FlowReadyCount = 'FLOW_READY_COUNT',
  FlowScheduledCount = 'FLOW_SCHEDULED_COUNT',
  FlowStatus = 'FLOW_STATUS',
  FlowTotalCount = 'FLOW_TOTAL_COUNT',
  InvoicesAmount = 'INVOICES_AMOUNT',
  InvoicesBillingCycle = 'INVOICES_BILLING_CYCLE',
  InvoicesContract = 'INVOICES_CONTRACT',
  InvoicesDueDate = 'INVOICES_DUE_DATE',
  InvoicesInvoiceNumber = 'INVOICES_INVOICE_NUMBER',
  InvoicesInvoicePreview = 'INVOICES_INVOICE_PREVIEW',
  InvoicesInvoiceStatus = 'INVOICES_INVOICE_STATUS',
  InvoicesIssueDate = 'INVOICES_ISSUE_DATE',
  InvoicesIssueDatePast = 'INVOICES_ISSUE_DATE_PAST',
  InvoicesOrganization = 'INVOICES_ORGANIZATION',
  OpportunitiesCommonColumn = 'OPPORTUNITIES_COMMON_COLUMN',
  OpportunitiesCreatedDate = 'OPPORTUNITIES_CREATED_DATE',
  OpportunitiesEstimatedArr = 'OPPORTUNITIES_ESTIMATED_ARR',
  OpportunitiesName = 'OPPORTUNITIES_NAME',
  OpportunitiesNextStep = 'OPPORTUNITIES_NEXT_STEP',
  OpportunitiesOrganization = 'OPPORTUNITIES_ORGANIZATION',
  OpportunitiesOwner = 'OPPORTUNITIES_OWNER',
  OpportunitiesStage = 'OPPORTUNITIES_STAGE',
  OpportunitiesTimeInStage = 'OPPORTUNITIES_TIME_IN_STAGE',
  OrganizationsAvatar = 'ORGANIZATIONS_AVATAR',
  OrganizationsChurnDate = 'ORGANIZATIONS_CHURN_DATE',
  OrganizationsCity = 'ORGANIZATIONS_CITY',
  OrganizationsContactCount = 'ORGANIZATIONS_CONTACT_COUNT',
  OrganizationsCreatedDate = 'ORGANIZATIONS_CREATED_DATE',
  OrganizationsEmployeeCount = 'ORGANIZATIONS_EMPLOYEE_COUNT',
  OrganizationsForecastArr = 'ORGANIZATIONS_FORECAST_ARR',
  OrganizationsHeadquarters = 'ORGANIZATIONS_HEADQUARTERS',
  OrganizationsIndustry = 'ORGANIZATIONS_INDUSTRY',
  OrganizationsIsPublic = 'ORGANIZATIONS_IS_PUBLIC',
  OrganizationsLastTouchpoint = 'ORGANIZATIONS_LAST_TOUCHPOINT',
  OrganizationsLastTouchpointDate = 'ORGANIZATIONS_LAST_TOUCHPOINT_DATE',
  OrganizationsLeadSource = 'ORGANIZATIONS_LEAD_SOURCE',
  OrganizationsLinkedinFollowerCount = 'ORGANIZATIONS_LINKEDIN_FOLLOWER_COUNT',
  OrganizationsLtv = 'ORGANIZATIONS_LTV',
  OrganizationsName = 'ORGANIZATIONS_NAME',
  OrganizationsOnboardingStatus = 'ORGANIZATIONS_ONBOARDING_STATUS',
  OrganizationsOwner = 'ORGANIZATIONS_OWNER',
  OrganizationsParentOrganization = 'ORGANIZATIONS_PARENT_ORGANIZATION',
  OrganizationsRelationship = 'ORGANIZATIONS_RELATIONSHIP',
  OrganizationsRenewalDate = 'ORGANIZATIONS_RENEWAL_DATE',
  OrganizationsRenewalLikelihood = 'ORGANIZATIONS_RENEWAL_LIKELIHOOD',
  OrganizationsSocials = 'ORGANIZATIONS_SOCIALS',
  OrganizationsStage = 'ORGANIZATIONS_STAGE',
  OrganizationsTags = 'ORGANIZATIONS_TAGS',
  OrganizationsWebsite = 'ORGANIZATIONS_WEBSITE',
  OrganizationsYearFounded = 'ORGANIZATIONS_YEAR_FOUNDED',
}

export type Comment = {
  __typename?: 'Comment';
  appSource: Scalars['String']['output'];
  content?: Maybe<Scalars['String']['output']>;
  contentType?: Maybe<Scalars['String']['output']>;
  createdAt: Scalars['Time']['output'];
  createdBy?: Maybe<User>;
  externalLinks: Array<ExternalSystem>;
  id: Scalars['ID']['output'];
  source: DataSource;
  sourceOfTruth: DataSource;
  updatedAt: Scalars['Time']['output'];
};

export enum ComparisonOperator {
  Between = 'BETWEEN',
  Contains = 'CONTAINS',
  Eq = 'EQ',
  Gt = 'GT',
  Gte = 'GTE',
  In = 'IN',
  IsEmpty = 'IS_EMPTY',
  IsNoneOf = 'IS_NONE_OF',
  IsNotEmpty = 'IS_NOT_EMPTY',
  IsNull = 'IS_NULL',
  Lt = 'LT',
  Lte = 'LTE',
  NotContains = 'NOT_CONTAINS',
  /** Not supported yet */
  NotEqual = 'NOT_EQUAL',
  StartsWith = 'STARTS_WITH',
}

/**
 * A contact represents an individual in customerOS.
 * **A `response` object.**
 */
export type Contact = MetadataInterface &
  Node & {
    __typename?: 'Contact';
    appSource?: Maybe<Scalars['String']['output']>;
    /** All users associated on linkedin to this contact */
    connectedUsers: Array<User>;
    /**
     * An ISO8601 timestamp recording when the contact was created in customerOS.
     * **Required**
     */
    createdAt: Scalars['Time']['output'];
    /**
     * User defined metadata appended to the contact record in customerOS.
     * **Required.  If no values it returns an empty array.**
     */
    customFields: Array<CustomField>;
    description?: Maybe<Scalars['String']['output']>;
    /**
     * All email addresses associated with a contact in customerOS.
     * **Required.  If no values it returns an empty array.**
     */
    emails: Array<Email>;
    enrichDetails: EnrichDetails;
    /** The first name of the contact in customerOS. */
    firstName?: Maybe<Scalars['String']['output']>;
    flows: Array<Flow>;
    hide?: Maybe<Scalars['Boolean']['output']>;
    /**
     * Deprecated, use metadata instead
     * @deprecated Use `metadata.id` instead
     */
    id: Scalars['ID']['output'];
    /**
     * `organizationName` and `jobTitle` of the contact if it has been associated with an organization.
     * **Required.  If no values it returns an empty array.**
     */
    jobRoles: Array<JobRole>;
    /**
     * Deprecated
     * @deprecated Use `tags` instead
     */
    label?: Maybe<Scalars['String']['output']>;
    /** The last name of the contact in customerOS. */
    lastName?: Maybe<Scalars['String']['output']>;
    latestOrganizationWithJobRole?: Maybe<OrganizationWithJobRole>;
    /**
     * All locations associated with a contact in customerOS.
     * **Required.  If no values it returns an empty array.**
     */
    locations: Array<Location>;
    metadata: Metadata;
    /** The name of the contact in customerOS, alternative for firstName + lastName. */
    name?: Maybe<Scalars['String']['output']>;
    organizations: OrganizationPage;
    /** Contact owner (user) */
    owner?: Maybe<User>;
    /**
     * All phone numbers associated with a contact in customerOS.
     * **Required.  If no values it returns an empty array.**
     */
    phoneNumbers: Array<PhoneNumber>;
    prefix?: Maybe<Scalars['String']['output']>;
    primaryEmail?: Maybe<Email>;
    profilePhotoUrl?: Maybe<Scalars['String']['output']>;
    socials: Array<Social>;
    source: DataSource;
    tags?: Maybe<Array<Tag>>;
    timelineEvents: Array<TimelineEvent>;
    timelineEventsTotalCount: Scalars['Int64']['output'];
    timezone?: Maybe<Scalars['String']['output']>;
    /**
     * Deprecated
     * @deprecated Use `prefix` instead
     */
    title?: Maybe<Scalars['String']['output']>;
    updatedAt: Scalars['Time']['output'];
    username?: Maybe<Scalars['String']['output']>;
  };

/**
 * A contact represents an individual in customerOS.
 * **A `response` object.**
 */
export type ContactOrganizationsArgs = {
  pagination?: InputMaybe<Pagination>;
  sort?: InputMaybe<Array<SortBy>>;
  where?: InputMaybe<Filter>;
};

/**
 * A contact represents an individual in customerOS.
 * **A `response` object.**
 */
export type ContactTimelineEventsArgs = {
  from?: InputMaybe<Scalars['Time']['input']>;
  size: Scalars['Int']['input'];
  timelineEventTypes?: InputMaybe<Array<TimelineEventType>>;
};

/**
 * A contact represents an individual in customerOS.
 * **A `response` object.**
 */
export type ContactTimelineEventsTotalCountArgs = {
  timelineEventTypes?: InputMaybe<Array<TimelineEventType>>;
};

/**
 * Create an individual in customerOS.
 * **A `create` object.**
 */
export type ContactInput = {
  appSource?: InputMaybe<Scalars['String']['input']>;
  /** An ISO8601 timestamp recording when the contact was created in customerOS. */
  createdAt?: InputMaybe<Scalars['Time']['input']>;
  /** Deprecated */
  customFields?: InputMaybe<Array<CustomFieldInput>>;
  description?: InputMaybe<Scalars['String']['input']>;
  /** An email addresses associated with the contact. */
  email?: InputMaybe<EmailInput>;
  externalReference?: InputMaybe<ExternalSystemReferenceInput>;
  /** The first name of the contact. */
  firstName?: InputMaybe<Scalars['String']['input']>;
  /** The last name of the contact. */
  lastName?: InputMaybe<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  ownerId?: InputMaybe<Scalars['ID']['input']>;
  /** A phone number associated with the contact. */
  phoneNumber?: InputMaybe<PhoneNumberInput>;
  /** The prefix of the contact. */
  prefix?: InputMaybe<Scalars['String']['input']>;
  profilePhotoUrl?: InputMaybe<Scalars['String']['input']>;
  socialUrl?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  templateId?: InputMaybe<Scalars['ID']['input']>;
  timezone?: InputMaybe<Scalars['String']['input']>;
  username?: InputMaybe<Scalars['String']['input']>;
};

export type ContactOrganizationInput = {
  contactId: Scalars['ID']['input'];
  organizationId: Scalars['ID']['input'];
};

export type ContactParticipant = {
  __typename?: 'ContactParticipant';
  contactParticipant: Contact;
  type?: Maybe<Scalars['String']['output']>;
};

export type ContactTagInput = {
  contactId: Scalars['ID']['input'];
  tag: TagIdOrNameInput;
};

/**
 * Updates data fields associated with an existing customer record in customerOS.
 * **An `update` object.**
 */
export type ContactUpdateInput = {
  description?: InputMaybe<Scalars['String']['input']>;
  firstName?: InputMaybe<Scalars['String']['input']>;
  id: Scalars['ID']['input'];
  lastName?: InputMaybe<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  patch?: InputMaybe<Scalars['Boolean']['input']>;
  prefix?: InputMaybe<Scalars['String']['input']>;
  profilePhotoUrl?: InputMaybe<Scalars['String']['input']>;
  timezone?: InputMaybe<Scalars['String']['input']>;
  username?: InputMaybe<Scalars['String']['input']>;
};

/**
 * Specifies how many pages of contact information has been returned in the query response.
 * **A `response` object.**
 */
export type ContactsPage = Pages & {
  __typename?: 'ContactsPage';
  /**
   * A contact entity in customerOS.
   * **Required.  If no values it returns an empty array.**
   */
  content: Array<Contact>;
  /**
   * Total number of elements in the query response.
   * **Required.**
   */
  totalElements: Scalars['Int64']['output'];
  /**
   * Total number of pages in the query response.
   * **Required.**
   */
  totalPages: Scalars['Int']['output'];
};

export type Contract = MetadataInterface & {
  __typename?: 'Contract';
  /**
   * Deprecated, use billingDetails instead.
   * @deprecated Use billingDetails instead.
   */
  addressLine1?: Maybe<Scalars['String']['output']>;
  /**
   * Deprecated, use billingDetails instead.
   * @deprecated Use billingDetails instead.
   */
  addressLine2?: Maybe<Scalars['String']['output']>;
  /**
   * Deprecated, use metadata instead.
   * @deprecated Use metadata instead.
   */
  appSource: Scalars['String']['output'];
  approved: Scalars['Boolean']['output'];
  attachments?: Maybe<Array<Attachment>>;
  autoRenew: Scalars['Boolean']['output'];
  /**
   * Deprecated, use billingDetails instead.
   * @deprecated Use billingDetails instead.
   */
  billingCycle?: Maybe<ContractBillingCycle>;
  billingDetails?: Maybe<BillingDetails>;
  billingEnabled: Scalars['Boolean']['output'];
  committedPeriodInMonths?: Maybe<Scalars['Int64']['output']>;
  /**
   * Deprecated, use committedPeriodInMonths instead.
   * @deprecated Use committedPeriodInMonths instead.
   */
  committedPeriods?: Maybe<Scalars['Int64']['output']>;
  contractEnded?: Maybe<Scalars['Time']['output']>;
  contractLineItems?: Maybe<Array<ServiceLineItem>>;
  contractName: Scalars['String']['output'];
  /**
   * Deprecated, use committedPeriodInMonths instead.
   * @deprecated Use committedPeriodInMonths instead.
   */
  contractRenewalCycle: ContractRenewalCycle;
  contractSigned?: Maybe<Scalars['Time']['output']>;
  contractStatus: ContractStatus;
  contractUrl?: Maybe<Scalars['String']['output']>;
  /**
   * Deprecated, use billingDetails instead.
   * @deprecated Use billingDetails instead.
   */
  country?: Maybe<Scalars['String']['output']>;
  /**
   * Deprecated, use metadata instead.
   * @deprecated Use metadata instead.
   */
  createdAt: Scalars['Time']['output'];
  createdBy?: Maybe<User>;
  currency?: Maybe<Currency>;
  /**
   * Deprecated, use contractEnded instead.
   * @deprecated Use contractEnded instead.
   */
  endedAt?: Maybe<Scalars['Time']['output']>;
  externalLinks: Array<ExternalSystem>;
  /**
   * Deprecated, use metadata instead.
   * @deprecated Use metadata instead.
   */
  id: Scalars['ID']['output'];
  /**
   * Deprecated, use billingDetails instead.
   * @deprecated Use billingDetails instead.
   */
  invoiceEmail?: Maybe<Scalars['String']['output']>;
  /**
   * Deprecated, use billingDetails instead.
   * @deprecated Use billingDetails instead.
   */
  invoiceNote?: Maybe<Scalars['String']['output']>;
  invoices: Array<Invoice>;
  /**
   * Deprecated, use billingDetails instead.
   * @deprecated Use billingDetails instead.
   */
  invoicingStartDate?: Maybe<Scalars['Time']['output']>;
  /**
   * Deprecated, use billingDetails instead.
   * @deprecated Use billingDetails instead.
   */
  locality?: Maybe<Scalars['String']['output']>;
  ltv: Scalars['Float']['output'];
  metadata: Metadata;
  /**
   * Deprecated, use contractName instead.
   * @deprecated Use contractName instead.
   */
  name: Scalars['String']['output'];
  opportunities?: Maybe<Array<Opportunity>>;
  /**
   * Deprecated, use billingDetails instead.
   * @deprecated Use billingDetails instead.
   */
  organizationLegalName?: Maybe<Scalars['String']['output']>;
  owner?: Maybe<User>;
  /**
   * Deprecated, use contractRenewalCycle instead.
   * @deprecated Use contractRenewalCycle instead.
   */
  renewalCycle: ContractRenewalCycle;
  /**
   * Deprecated, use committedPeriods instead.
   * @deprecated Use committedPeriods instead.
   */
  renewalPeriods?: Maybe<Scalars['Int64']['output']>;
  /**
   * Deprecated, use contractLineItems instead.
   * @deprecated Use contractLineItems instead.
   */
  serviceLineItems?: Maybe<Array<ServiceLineItem>>;
  serviceStarted?: Maybe<Scalars['Time']['output']>;
  /**
   * Deprecated, use serviceStarted instead.
   * @deprecated Use serviceStarted instead.
   */
  serviceStartedAt?: Maybe<Scalars['Time']['output']>;
  /**
   * Deprecated, use contractSigned instead.
   * @deprecated Use contractSigned instead.
   */
  signedAt?: Maybe<Scalars['Time']['output']>;
  /**
   * Deprecated, use metadata instead.
   * @deprecated Use metadata instead.
   */
  source: DataSource;
  /**
   * Deprecated, use metadata instead.
   * @deprecated Use metadata instead.
   */
  sourceOfTruth: DataSource;
  /**
   * Deprecated, use contractStatus instead.
   * @deprecated Use contractStatus instead.
   */
  status: ContractStatus;
  upcomingInvoices: Array<Invoice>;
  /**
   * Deprecated, use metadata instead.
   * @deprecated Use metadata instead.
   */
  updatedAt: Scalars['Time']['output'];
  /**
   * Deprecated, use billingDetails instead.
   * @deprecated Use billingDetails instead.
   */
  zip?: Maybe<Scalars['String']['output']>;
};

/** Deprecated */
export enum ContractBillingCycle {
  AnnualBilling = 'ANNUAL_BILLING',
  CustomBilling = 'CUSTOM_BILLING',
  MonthlyBilling = 'MONTHLY_BILLING',
  None = 'NONE',
  QuarterlyBilling = 'QUARTERLY_BILLING',
}

export type ContractInput = {
  appSource?: InputMaybe<Scalars['String']['input']>;
  approved?: InputMaybe<Scalars['Boolean']['input']>;
  autoRenew?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  billingCycle?: InputMaybe<ContractBillingCycle>;
  billingEnabled?: InputMaybe<Scalars['Boolean']['input']>;
  committedPeriodInMonths?: InputMaybe<Scalars['Int64']['input']>;
  /** Deprecated, use committedPeriodInMonths instead. */
  committedPeriods?: InputMaybe<Scalars['Int64']['input']>;
  contractName?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated, use committedPeriodInMonths instead. */
  contractRenewalCycle?: InputMaybe<ContractRenewalCycle>;
  contractSigned?: InputMaybe<Scalars['Time']['input']>;
  contractUrl?: InputMaybe<Scalars['String']['input']>;
  currency?: InputMaybe<Currency>;
  dueDays?: InputMaybe<Scalars['Int64']['input']>;
  /** Deprecated */
  externalReference?: InputMaybe<ExternalSystemReferenceInput>;
  /** Deprecated */
  invoicingStartDate?: InputMaybe<Scalars['Time']['input']>;
  /** Deprecated */
  name?: InputMaybe<Scalars['String']['input']>;
  organizationId: Scalars['ID']['input'];
  /** Deprecated */
  renewalCycle?: InputMaybe<ContractRenewalCycle>;
  /** Deprecated */
  renewalPeriods?: InputMaybe<Scalars['Int64']['input']>;
  serviceStarted?: InputMaybe<Scalars['Time']['input']>;
  /** Deprecated */
  serviceStartedAt?: InputMaybe<Scalars['Time']['input']>;
  /** Deprecated */
  signedAt?: InputMaybe<Scalars['Time']['input']>;
};

export type ContractPage = Pages & {
  __typename?: 'ContractPage';
  content: Array<Contract>;
  totalAvailable: Scalars['Int64']['output'];
  totalElements: Scalars['Int64']['output'];
  totalPages: Scalars['Int']['output'];
};

/** Deprecated */
export enum ContractRenewalCycle {
  AnnualRenewal = 'ANNUAL_RENEWAL',
  MonthlyRenewal = 'MONTHLY_RENEWAL',
  None = 'NONE',
  QuarterlyRenewal = 'QUARTERLY_RENEWAL',
}

export type ContractRenewalInput = {
  contractId: Scalars['ID']['input'];
  renewalDate?: InputMaybe<Scalars['Time']['input']>;
};

export enum ContractStatus {
  Draft = 'DRAFT',
  Ended = 'ENDED',
  Live = 'LIVE',
  OutOfContract = 'OUT_OF_CONTRACT',
  Scheduled = 'SCHEDULED',
  Undefined = 'UNDEFINED',
}

export type ContractUpdateInput = {
  /** Deprecated */
  addressLine1?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  addressLine2?: InputMaybe<Scalars['String']['input']>;
  appSource?: InputMaybe<Scalars['String']['input']>;
  approved?: InputMaybe<Scalars['Boolean']['input']>;
  autoRenew?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  billingCycle?: InputMaybe<ContractBillingCycle>;
  billingDetails?: InputMaybe<BillingDetailsInput>;
  billingEnabled?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  canPayWithBankTransfer?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  canPayWithCard?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  canPayWithDirectDebit?: InputMaybe<Scalars['Boolean']['input']>;
  committedPeriodInMonths?: InputMaybe<Scalars['Int64']['input']>;
  /** Deprecated, use committedPeriodInMonths instead. */
  committedPeriods?: InputMaybe<Scalars['Int64']['input']>;
  contractEnded?: InputMaybe<Scalars['Time']['input']>;
  contractId: Scalars['ID']['input'];
  contractName?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated, use committedPeriodInMonths instead. */
  contractRenewalCycle?: InputMaybe<ContractRenewalCycle>;
  contractSigned?: InputMaybe<Scalars['Time']['input']>;
  contractUrl?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  country?: InputMaybe<Scalars['String']['input']>;
  currency?: InputMaybe<Currency>;
  /** Deprecated */
  endedAt?: InputMaybe<Scalars['Time']['input']>;
  /** Deprecated */
  invoiceEmail?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  invoiceNote?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  invoicingStartDate?: InputMaybe<Scalars['Time']['input']>;
  /** Deprecated */
  locality?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  name?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  organizationLegalName?: InputMaybe<Scalars['String']['input']>;
  patch?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  renewalCycle?: InputMaybe<ContractRenewalCycle>;
  /** Deprecated */
  renewalPeriods?: InputMaybe<Scalars['Int64']['input']>;
  serviceStarted?: InputMaybe<Scalars['Time']['input']>;
  /** Deprecated */
  serviceStartedAt?: InputMaybe<Scalars['Time']['input']>;
  /** Deprecated */
  signedAt?: InputMaybe<Scalars['Time']['input']>;
  /** Deprecated */
  zip?: InputMaybe<Scalars['String']['input']>;
};

export type Country = {
  __typename?: 'Country';
  codeA2: Scalars['String']['output'];
  codeA3: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  phoneCode: Scalars['String']['output'];
};

export enum Currency {
  Aud = 'AUD',
  Brl = 'BRL',
  Cad = 'CAD',
  Chf = 'CHF',
  Cny = 'CNY',
  Eur = 'EUR',
  Gbp = 'GBP',
  Hkd = 'HKD',
  Inr = 'INR',
  Jpy = 'JPY',
  Krw = 'KRW',
  Mxn = 'MXN',
  Nok = 'NOK',
  Nzd = 'NZD',
  Ron = 'RON',
  Sek = 'SEK',
  Sgd = 'SGD',
  Try = 'TRY',
  Usd = 'USD',
  Zar = 'ZAR',
}

export enum CustomEntityType {
  Contact = 'Contact',
  Organization = 'Organization',
}

/**
 * Describes a custom, user-defined field associated with a `Contact`.
 * **A `return` object.**
 */
export type CustomField = Node & {
  __typename?: 'CustomField';
  createdAt: Scalars['Time']['output'];
  /**
   * Datatype of the custom field.
   * **Required**
   */
  datatype: CustomFieldDataType;
  /**
   * The unique ID associated with the custom field.
   * **Required**
   */
  id: Scalars['ID']['output'];
  /**
   * The name of the custom field.
   * **Required**
   */
  name: Scalars['String']['output'];
  /** The source of the custom field value */
  source: DataSource;
  template?: Maybe<CustomFieldTemplate>;
  updatedAt: Scalars['Time']['output'];
  /**
   * The value of the custom field.
   * **Required**
   */
  value: Scalars['Any']['output'];
};

export enum CustomFieldDataType {
  Bool = 'BOOL',
  Datetime = 'DATETIME',
  Decimal = 'DECIMAL',
  Integer = 'INTEGER',
  Text = 'TEXT',
}

export type CustomFieldEntityType = {
  entityType: CustomEntityType;
  id: Scalars['ID']['input'];
};

/**
 * Describes a custom, user-defined field associated with a `Contact` of type String.
 * **A `create` object.**
 */
export type CustomFieldInput = {
  /** Datatype of the custom field. */
  datatype?: InputMaybe<CustomFieldDataType>;
  /** Deprecated */
  id?: InputMaybe<Scalars['ID']['input']>;
  /** The name of the custom field. */
  name?: InputMaybe<Scalars['String']['input']>;
  templateId?: InputMaybe<Scalars['ID']['input']>;
  /**
   * The value of the custom field.
   * **Required**
   */
  value: Scalars['Any']['input'];
};

export type CustomFieldTemplate = Node & {
  __typename?: 'CustomFieldTemplate';
  createdAt: Scalars['Time']['output'];
  entityType: EntityType;
  id: Scalars['ID']['output'];
  length?: Maybe<Scalars['Int64']['output']>;
  max?: Maybe<Scalars['Int64']['output']>;
  min?: Maybe<Scalars['Int64']['output']>;
  name: Scalars['String']['output'];
  order?: Maybe<Scalars['Int64']['output']>;
  required?: Maybe<Scalars['Boolean']['output']>;
  type: CustomFieldTemplateType;
  updatedAt: Scalars['Time']['output'];
  validValues: Array<Scalars['String']['output']>;
};

export type CustomFieldTemplateInput = {
  entityType?: InputMaybe<EntityType>;
  id?: InputMaybe<Scalars['ID']['input']>;
  length?: InputMaybe<Scalars['Int64']['input']>;
  max?: InputMaybe<Scalars['Int64']['input']>;
  min?: InputMaybe<Scalars['Int64']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  order?: InputMaybe<Scalars['Int64']['input']>;
  required?: InputMaybe<Scalars['Boolean']['input']>;
  type?: InputMaybe<CustomFieldTemplateType>;
  validValues?: InputMaybe<Array<Scalars['String']['input']>>;
};

export enum CustomFieldTemplateType {
  FreeText = 'FREE_TEXT',
  Number = 'NUMBER',
  SingleSelect = 'SINGLE_SELECT',
}

/**
 * Describes a custom, user-defined field associated with a `Contact`.
 * **An `update` object.**
 */
export type CustomFieldUpdateInput = {
  /**
   * Datatype of the custom field.
   * **Required**
   */
  datatype: CustomFieldDataType;
  /**
   * The unique ID associated with the custom field.
   * **Required**
   */
  id: Scalars['ID']['input'];
  /**
   * The name of the custom field.
   * **Required**
   */
  name: Scalars['String']['input'];
  /**
   * The value of the custom field.
   * **Required**
   */
  value: Scalars['Any']['input'];
};

export type CustomerContact = {
  __typename?: 'CustomerContact';
  email: CustomerEmail;
  id: Scalars['ID']['output'];
};

export type CustomerContactInput = {
  appSource?: InputMaybe<Scalars['String']['input']>;
  /** An ISO8601 timestamp recording when the contact was created in customerOS. */
  createdAt?: InputMaybe<Scalars['Time']['input']>;
  description?: InputMaybe<Scalars['String']['input']>;
  /** An email addresses associted with the contact. */
  email?: InputMaybe<EmailInput>;
  /** The first name of the contact. */
  firstName?: InputMaybe<Scalars['String']['input']>;
  /** The last name of the contact. */
  lastName?: InputMaybe<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  /** The prefix of the contact. */
  prefix?: InputMaybe<Scalars['String']['input']>;
  timezone?: InputMaybe<Scalars['String']['input']>;
};

export type CustomerEmail = {
  __typename?: 'CustomerEmail';
  id: Scalars['ID']['output'];
};

export type CustomerJobRole = {
  __typename?: 'CustomerJobRole';
  id: Scalars['ID']['output'];
};

export type CustomerUser = {
  __typename?: 'CustomerUser';
  id: Scalars['ID']['output'];
  jobRole: CustomerJobRole;
};

export type DashboardArrBreakdown = {
  __typename?: 'DashboardARRBreakdown';
  arrBreakdown: Scalars['Float']['output'];
  increasePercentage: Scalars['String']['output'];
  perMonth: Array<Maybe<DashboardArrBreakdownPerMonth>>;
};

export type DashboardArrBreakdownPerMonth = {
  __typename?: 'DashboardARRBreakdownPerMonth';
  cancellations: Scalars['Float']['output'];
  churned: Scalars['Float']['output'];
  downgrades: Scalars['Float']['output'];
  month: Scalars['Int']['output'];
  newlyContracted: Scalars['Float']['output'];
  renewals: Scalars['Float']['output'];
  upsells: Scalars['Float']['output'];
  year: Scalars['Int']['output'];
};

export type DashboardCustomerMap = {
  __typename?: 'DashboardCustomerMap';
  arr: Scalars['Float']['output'];
  contractSignedDate: Scalars['Time']['output'];
  organization: Organization;
  organizationId: Scalars['ID']['output'];
  state: DashboardCustomerMapState;
};

export enum DashboardCustomerMapState {
  /**
   * Deprecated
   * @deprecated Use HIGH_RISK instead
   */
  AtRisk = 'AT_RISK',
  Churned = 'CHURNED',
  HighRisk = 'HIGH_RISK',
  MediumRisk = 'MEDIUM_RISK',
  Ok = 'OK',
}

export type DashboardGrossRevenueRetention = {
  __typename?: 'DashboardGrossRevenueRetention';
  grossRevenueRetention: Scalars['Float']['output'];
  /**
   * Deprecated
   * @deprecated Use increasePercentageValue instead
   */
  increasePercentage: Scalars['String']['output'];
  increasePercentageValue: Scalars['Float']['output'];
  perMonth: Array<Maybe<DashboardGrossRevenueRetentionPerMonth>>;
};

export type DashboardGrossRevenueRetentionPerMonth = {
  __typename?: 'DashboardGrossRevenueRetentionPerMonth';
  month: Scalars['Int']['output'];
  percentage: Scalars['Float']['output'];
  year: Scalars['Int']['output'];
};

export type DashboardMrrPerCustomer = {
  __typename?: 'DashboardMRRPerCustomer';
  increasePercentage: Scalars['String']['output'];
  mrrPerCustomer: Scalars['Float']['output'];
  perMonth: Array<Maybe<DashboardMrrPerCustomerPerMonth>>;
};

export type DashboardMrrPerCustomerPerMonth = {
  __typename?: 'DashboardMRRPerCustomerPerMonth';
  month: Scalars['Int']['output'];
  value: Scalars['Float']['output'];
  year: Scalars['Int']['output'];
};

export type DashboardNewCustomers = {
  __typename?: 'DashboardNewCustomers';
  perMonth: Array<Maybe<DashboardNewCustomersPerMonth>>;
  thisMonthCount: Scalars['Int']['output'];
  thisMonthIncreasePercentage: Scalars['String']['output'];
};

export type DashboardNewCustomersPerMonth = {
  __typename?: 'DashboardNewCustomersPerMonth';
  count: Scalars['Int']['output'];
  month: Scalars['Int']['output'];
  year: Scalars['Int']['output'];
};

export type DashboardOnboardingCompletion = {
  __typename?: 'DashboardOnboardingCompletion';
  completionPercentage: Scalars['Float']['output'];
  increasePercentage: Scalars['Float']['output'];
  perMonth: Array<DashboardOnboardingCompletionPerMonth>;
};

export type DashboardOnboardingCompletionPerMonth = {
  __typename?: 'DashboardOnboardingCompletionPerMonth';
  month: Scalars['Int']['output'];
  value: Scalars['Float']['output'];
  year: Scalars['Int']['output'];
};

export type DashboardPeriodInput = {
  end: Scalars['Time']['input'];
  start: Scalars['Time']['input'];
};

export type DashboardRetentionRate = {
  __typename?: 'DashboardRetentionRate';
  /**
   * Deprecated
   * @deprecated Use increasePercentageValue instead
   */
  increasePercentage: Scalars['String']['output'];
  increasePercentageValue: Scalars['Float']['output'];
  perMonth: Array<Maybe<DashboardRetentionRatePerMonth>>;
  retentionRate: Scalars['Float']['output'];
};

export type DashboardRetentionRatePerMonth = {
  __typename?: 'DashboardRetentionRatePerMonth';
  churnCount: Scalars['Int']['output'];
  month: Scalars['Int']['output'];
  renewCount: Scalars['Int']['output'];
  year: Scalars['Int']['output'];
};

export type DashboardRevenueAtRisk = {
  __typename?: 'DashboardRevenueAtRisk';
  atRisk: Scalars['Float']['output'];
  highConfidence: Scalars['Float']['output'];
};

export type DashboardTimeToOnboard = {
  __typename?: 'DashboardTimeToOnboard';
  increasePercentage?: Maybe<Scalars['Float']['output']>;
  perMonth: Array<DashboardTimeToOnboardPerMonth>;
  timeToOnboard?: Maybe<Scalars['Float']['output']>;
};

export type DashboardTimeToOnboardPerMonth = {
  __typename?: 'DashboardTimeToOnboardPerMonth';
  month: Scalars['Int']['output'];
  value: Scalars['Float']['output'];
  year: Scalars['Int']['output'];
};

export enum DataSource {
  Attio = 'ATTIO',
  Close = 'CLOSE',
  Hubspot = 'HUBSPOT',
  Intercom = 'INTERCOM',
  Mailstack = 'MAILSTACK',
  Mixpanel = 'MIXPANEL',
  Na = 'NA',
  Openline = 'OPENLINE',
  Outlook = 'OUTLOOK',
  Pipedrive = 'PIPEDRIVE',
  Salesforce = 'SALESFORCE',
  Shopify = 'SHOPIFY',
  Slack = 'SLACK',
  Stripe = 'STRIPE',
  Unthread = 'UNTHREAD',
  Webscrape = 'WEBSCRAPE',
  ZendeskSell = 'ZENDESK_SELL',
  ZendeskSupport = 'ZENDESK_SUPPORT',
}

export type DeleteResponse = {
  __typename?: 'DeleteResponse';
  accepted: Scalars['Boolean']['output'];
  completed: Scalars['Boolean']['output'];
};

/**
 * Describes an email address associated with a `Contact` in customerOS.
 * **A `return` object.**
 */
export type Email = {
  __typename?: 'Email';
  /** @deprecated No longer supported */
  appSource: Scalars['String']['output'];
  contacts: Array<Contact>;
  createdAt: Scalars['Time']['output'];
  /** An email address assocaited with the contact in customerOS. */
  email?: Maybe<Scalars['String']['output']>;
  emailValidationDetails: EmailValidationDetails;
  /**
   * The unique ID associated with the contact in customerOS.
   * **Required**
   */
  id: Scalars['ID']['output'];
  /**
   * Describes the type of email address (WORK, PERSONAL, etc).
   * @deprecated No longer supported
   */
  label?: Maybe<EmailLabel>;
  organizations: Array<Organization>;
  /**
   * Identifies whether the email address is primary or not.
   * **Required.**
   */
  primary: Scalars['Boolean']['output'];
  rawEmail?: Maybe<Scalars['String']['output']>;
  source: DataSource;
  updatedAt: Scalars['Time']['output'];
  users: Array<User>;
  work?: Maybe<Scalars['Boolean']['output']>;
};

export enum EmailDeliverable {
  Deliverable = 'DELIVERABLE',
  Undeliverable = 'UNDELIVERABLE',
  Unknown = 'UNKNOWN',
}

/**
 * Describes an email address associated with a `Contact` in customerOS.
 * **A `create` object.**
 */
export type EmailInput = {
  appSource?: InputMaybe<Scalars['String']['input']>;
  /**
   * An email address associated with the contact in customerOS.
   * **Required.**
   */
  email: Scalars['String']['input'];
  label?: InputMaybe<EmailLabel>;
  /**
   * Identifies whether the email address is primary or not.
   * **Required.**
   */
  primary?: InputMaybe<Scalars['Boolean']['input']>;
};

/**
 * Describes the type of email address (WORK, PERSONAL, etc).
 * **A `return` object.
 */
export enum EmailLabel {
  Main = 'MAIN',
  Other = 'OTHER',
  Personal = 'PERSONAL',
  Work = 'WORK',
}

export type EmailParticipant = {
  __typename?: 'EmailParticipant';
  emailParticipant: Email;
  type?: Maybe<Scalars['String']['output']>;
};

/**
 * Describes an email address associated with a `Contact` in customerOS.
 * **An `update` object.**
 */
export type EmailRelationUpdateInput = {
  /** Deprecated */
  email?: InputMaybe<Scalars['String']['input']>;
  /**
   * An email address assocaited with the contact in customerOS.
   * **Required.**
   */
  id: Scalars['ID']['input'];
  /** Describes the type of email address (WORK, PERSONAL, etc). */
  label?: InputMaybe<EmailLabel>;
  /**
   * Identifies whether the email address is primary or not.
   * **Required.**
   */
  primary?: InputMaybe<Scalars['Boolean']['input']>;
};

export type EmailUpdateAddressInput = {
  email: Scalars['String']['input'];
  id: Scalars['ID']['input'];
};

export type EmailValidationDetails = {
  __typename?: 'EmailValidationDetails';
  alternateEmail?: Maybe<Scalars['String']['output']>;
  canConnectSmtp?: Maybe<Scalars['Boolean']['output']>;
  deliverable?: Maybe<EmailDeliverable>;
  firewall?: Maybe<Scalars['String']['output']>;
  isCatchAll?: Maybe<Scalars['Boolean']['output']>;
  /** @deprecated No longer supported */
  isDeliverable?: Maybe<Scalars['Boolean']['output']>;
  isFirewalled?: Maybe<Scalars['Boolean']['output']>;
  isFreeAccount?: Maybe<Scalars['Boolean']['output']>;
  isMailboxFull?: Maybe<Scalars['Boolean']['output']>;
  isPrimaryDomain?: Maybe<Scalars['Boolean']['output']>;
  isRisky?: Maybe<Scalars['Boolean']['output']>;
  isRoleAccount?: Maybe<Scalars['Boolean']['output']>;
  isValidSyntax?: Maybe<Scalars['Boolean']['output']>;
  primaryDomain?: Maybe<Scalars['String']['output']>;
  provider?: Maybe<Scalars['String']['output']>;
  smtpSuccess?: Maybe<Scalars['Boolean']['output']>;
  verified: Scalars['Boolean']['output'];
  verifyingCheckAll: Scalars['Boolean']['output'];
};

export type EmailVariableEntity = {
  __typename?: 'EmailVariableEntity';
  type: EmailVariableEntityType;
  variables: Array<EmailVariableName>;
};

export enum EmailVariableEntityType {
  Contact = 'CONTACT',
}

export enum EmailVariableName {
  ContactEmail = 'CONTACT_EMAIL',
  ContactFirstName = 'CONTACT_FIRST_NAME',
  ContactFullName = 'CONTACT_FULL_NAME',
  ContactLastName = 'CONTACT_LAST_NAME',
  OrganizationName = 'ORGANIZATION_NAME',
  SenderFirstName = 'SENDER_FIRST_NAME',
  SenderLastName = 'SENDER_LAST_NAME',
}

export type EnrichDetails = {
  __typename?: 'EnrichDetails';
  emailEnrichedAt?: Maybe<Scalars['Time']['output']>;
  emailFound?: Maybe<Scalars['Boolean']['output']>;
  emailRequestedAt?: Maybe<Scalars['Time']['output']>;
  enrichedAt?: Maybe<Scalars['Time']['output']>;
  failedAt?: Maybe<Scalars['Time']['output']>;
  mobilePhoneEnrichedAt?: Maybe<Scalars['Time']['output']>;
  mobilePhoneFound?: Maybe<Scalars['Boolean']['output']>;
  mobilePhoneRequestedAt?: Maybe<Scalars['Time']['output']>;
  requestedAt?: Maybe<Scalars['Time']['output']>;
};

export enum EntityType {
  Contact = 'CONTACT',
  LogEntry = 'LOG_ENTRY',
  Opportunity = 'OPPORTUNITY',
  Organization = 'ORGANIZATION',
}

export type ExternalSystem = {
  __typename?: 'ExternalSystem';
  externalId?: Maybe<Scalars['String']['output']>;
  externalSource?: Maybe<Scalars['String']['output']>;
  externalUrl?: Maybe<Scalars['String']['output']>;
  syncDate?: Maybe<Scalars['Time']['output']>;
  type: ExternalSystemType;
};

export type ExternalSystemInput = {
  name: Scalars['String']['input'];
};

export type ExternalSystemInstance = {
  __typename?: 'ExternalSystemInstance';
  stripeDetails?: Maybe<ExternalSystemStripeDetails>;
  type: ExternalSystemType;
};

export type ExternalSystemReferenceInput = {
  externalId: Scalars['ID']['input'];
  externalSource?: InputMaybe<Scalars['String']['input']>;
  externalUrl?: InputMaybe<Scalars['String']['input']>;
  syncDate?: InputMaybe<Scalars['Time']['input']>;
  type: ExternalSystemType;
};

export type ExternalSystemStripeDetails = {
  __typename?: 'ExternalSystemStripeDetails';
  paymentMethodTypes: Array<Scalars['String']['output']>;
};

export enum ExternalSystemType {
  Attio = 'ATTIO',
  Calcom = 'CALCOM',
  Close = 'CLOSE',
  Hubspot = 'HUBSPOT',
  Intercom = 'INTERCOM',
  Mixpanel = 'MIXPANEL',
  Outlook = 'OUTLOOK',
  Pipedrive = 'PIPEDRIVE',
  Salesforce = 'SALESFORCE',
  Slack = 'SLACK',
  Stripe = 'STRIPE',
  Unthread = 'UNTHREAD',
  Weconnect = 'WECONNECT',
  ZendeskSell = 'ZENDESK_SELL',
  ZendeskSupport = 'ZENDESK_SUPPORT',
}

export type Filter = {
  AND?: InputMaybe<Array<Filter>>;
  NOT?: InputMaybe<Filter>;
  OR?: InputMaybe<Array<Filter>>;
  filter?: InputMaybe<FilterItem>;
};

export type FilterItem = {
  caseSensitive?: InputMaybe<Scalars['Boolean']['input']>;
  includeEmpty?: InputMaybe<Scalars['Boolean']['input']>;
  operation?: ComparisonOperator;
  property: Scalars['String']['input'];
  value: Scalars['Any']['input'];
};

export type Flow = MetadataInterface & {
  __typename?: 'Flow';
  contacts: Array<FlowContact>;
  description: Scalars['String']['output'];
  edges: Scalars['String']['output'];
  metadata: Metadata;
  name: Scalars['String']['output'];
  nodes: Scalars['String']['output'];
  senders: Array<FlowSender>;
  statistics: FlowStatistics;
  status: FlowStatus;
};

export type FlowActionInputData = {
  email_new?: InputMaybe<FlowActionInputDataEmail>;
  email_reply?: InputMaybe<FlowActionInputDataEmail>;
  linkedin_connection_request?: InputMaybe<FlowActionInputDataLinkedinConnectionRequest>;
  linkedin_message?: InputMaybe<FlowActionInputDataLinkedinMessage>;
  wait?: InputMaybe<FlowActionInputDataWait>;
};

export type FlowActionInputDataEmail = {
  bodyTemplate: Scalars['String']['input'];
  replyToId?: InputMaybe<Scalars['String']['input']>;
  subject: Scalars['String']['input'];
};

export type FlowActionInputDataLinkedinConnectionRequest = {
  messageTemplate: Scalars['String']['input'];
};

export type FlowActionInputDataLinkedinMessage = {
  messageTemplate: Scalars['String']['input'];
};

export type FlowActionInputDataWait = {
  minutes: Scalars['Int64']['input'];
};

export enum FlowActionStatus {
  Active = 'ACTIVE',
  Archived = 'ARCHIVED',
  Inactive = 'INACTIVE',
  Paused = 'PAUSED',
}

export type FlowContact = MetadataInterface & {
  __typename?: 'FlowContact';
  contact: Contact;
  metadata: Metadata;
  scheduledAction?: Maybe<Scalars['String']['output']>;
  scheduledAt?: Maybe<Scalars['Time']['output']>;
  status: FlowParticipantStatus;
};

export type FlowMergeInput = {
  edges: Scalars['String']['input'];
  id?: InputMaybe<Scalars['ID']['input']>;
  name: Scalars['String']['input'];
  nodes: Scalars['String']['input'];
};

export enum FlowParticipantStatus {
  Completed = 'COMPLETED',
  GoalAchieved = 'GOAL_ACHIEVED',
  InProgress = 'IN_PROGRESS',
  OnHold = 'ON_HOLD',
  Ready = 'READY',
  Scheduled = 'SCHEDULED',
}

export type FlowSender = MetadataInterface & {
  __typename?: 'FlowSender';
  flow?: Maybe<Flow>;
  metadata: Metadata;
  user?: Maybe<User>;
};

export type FlowSenderMergeInput = {
  id?: InputMaybe<Scalars['ID']['input']>;
  userId?: InputMaybe<Scalars['ID']['input']>;
};

export type FlowStatistics = {
  __typename?: 'FlowStatistics';
  completed: Scalars['Int64']['output'];
  goalAchieved: Scalars['Int64']['output'];
  inProgress: Scalars['Int64']['output'];
  onHold: Scalars['Int64']['output'];
  ready: Scalars['Int64']['output'];
  scheduled: Scalars['Int64']['output'];
  total: Scalars['Int64']['output'];
};

export enum FlowStatus {
  Active = 'ACTIVE',
  Archived = 'ARCHIVED',
  Inactive = 'INACTIVE',
  Paused = 'PAUSED',
  Scheduling = 'SCHEDULING',
}

export enum FundingRound {
  Angel = 'ANGEL',
  Bridge = 'BRIDGE',
  FriendsAndFamily = 'FRIENDS_AND_FAMILY',
  Ipo = 'IPO',
  PreSeed = 'PRE_SEED',
  Seed = 'SEED',
  SeriesA = 'SERIES_A',
  SeriesB = 'SERIES_B',
  SeriesC = 'SERIES_C',
  SeriesD = 'SERIES_D',
  SeriesE = 'SERIES_E',
  SeriesF = 'SERIES_F',
}

export type GCliAttributeKeyValuePair = {
  __typename?: 'GCliAttributeKeyValuePair';
  display?: Maybe<Scalars['String']['output']>;
  key: Scalars['String']['output'];
  value: Scalars['String']['output'];
};

export enum GCliCacheItemType {
  Contact = 'CONTACT',
  Organization = 'ORGANIZATION',
  State = 'STATE',
}

export type GCliItem = {
  __typename?: 'GCliItem';
  data?: Maybe<Array<GCliAttributeKeyValuePair>>;
  display: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  type: GCliSearchResultType;
};

export enum GCliSearchResultType {
  Contact = 'CONTACT',
  Email = 'EMAIL',
  Organization = 'ORGANIZATION',
  OrganizationRelationship = 'ORGANIZATION_RELATIONSHIP',
  State = 'STATE',
}

export type GlobalCache = {
  __typename?: 'GlobalCache';
  activeEmailTokens: Array<GlobalCacheEmailToken>;
  cdnLogoUrl: Scalars['String']['output'];
  contractsExist: Scalars['Boolean']['output'];
  gCliCache: Array<GCliItem>;
  inactiveEmailTokens: Array<GlobalCacheEmailToken>;
  isOwner: Scalars['Boolean']['output'];
  mailboxes: Array<Scalars['String']['output']>;
  maxARRForecastValue: Scalars['Float']['output'];
  minARRForecastValue: Scalars['Float']['output'];
  user: User;
};

export type GlobalCacheEmailToken = {
  __typename?: 'GlobalCacheEmailToken';
  email: Scalars['String']['output'];
  provider: Scalars['String']['output'];
};

export type InteractionEvent = Node & {
  __typename?: 'InteractionEvent';
  actionItems?: Maybe<Array<ActionItem>>;
  actions?: Maybe<Array<Action>>;
  appSource: Scalars['String']['output'];
  channel: Scalars['String']['output'];
  channelData?: Maybe<Scalars['String']['output']>;
  content?: Maybe<Scalars['String']['output']>;
  contentType?: Maybe<Scalars['String']['output']>;
  createdAt: Scalars['Time']['output'];
  customerOSInternalIdentifier?: Maybe<Scalars['String']['output']>;
  eventIdentifier?: Maybe<Scalars['String']['output']>;
  eventType?: Maybe<Scalars['String']['output']>;
  externalLinks: Array<ExternalSystem>;
  id: Scalars['ID']['output'];
  includes: Array<Attachment>;
  interactionSession?: Maybe<InteractionSession>;
  issue?: Maybe<Issue>;
  meeting?: Maybe<Meeting>;
  repliesTo?: Maybe<InteractionEvent>;
  sentBy: Array<InteractionEventParticipant>;
  sentTo: Array<InteractionEventParticipant>;
  source: DataSource;
  sourceOfTruth: DataSource;
};

export type InteractionEventParticipant =
  | ContactParticipant
  | EmailParticipant
  | JobRoleParticipant
  | OrganizationParticipant
  | PhoneNumberParticipant
  | UserParticipant;

export type InteractionSession = Node & {
  __typename?: 'InteractionSession';
  appSource: Scalars['String']['output'];
  attendedBy: Array<InteractionSessionParticipant>;
  channel?: Maybe<Scalars['String']['output']>;
  channelData?: Maybe<Scalars['String']['output']>;
  createdAt: Scalars['Time']['output'];
  events: Array<InteractionEvent>;
  id: Scalars['ID']['output'];
  identifier: Scalars['String']['output'];
  name: Scalars['String']['output'];
  source: DataSource;
  sourceOfTruth: DataSource;
  status: Scalars['String']['output'];
  type?: Maybe<Scalars['String']['output']>;
  updatedAt: Scalars['Time']['output'];
};

export type InteractionSessionParticipant =
  | ContactParticipant
  | EmailParticipant
  | PhoneNumberParticipant
  | UserParticipant;

export enum InternalStage {
  ClosedLost = 'CLOSED_LOST',
  ClosedWon = 'CLOSED_WON',
  Open = 'OPEN',
}

export enum InternalType {
  CrossSell = 'CROSS_SELL',
  Nbo = 'NBO',
  Renewal = 'RENEWAL',
  Upsell = 'UPSELL',
}

export type Invoice = MetadataInterface & {
  __typename?: 'Invoice';
  amountDue: Scalars['Float']['output'];
  amountPaid: Scalars['Float']['output'];
  amountRemaining: Scalars['Float']['output'];
  billingCycleInMonths: Scalars['Int64']['output'];
  contract: Contract;
  currency: Scalars['String']['output'];
  customer: InvoiceCustomer;
  /**
   * Deprecated
   * @deprecated not used
   */
  domesticPaymentsBankInfo?: Maybe<Scalars['String']['output']>;
  dryRun: Scalars['Boolean']['output'];
  due: Scalars['Time']['output'];
  /**
   * Deprecated
   * @deprecated not used
   */
  internationalPaymentsBankInfo?: Maybe<Scalars['String']['output']>;
  invoiceLineItems: Array<InvoiceLine>;
  invoiceNumber: Scalars['String']['output'];
  invoicePeriodEnd: Scalars['Time']['output'];
  invoicePeriodStart: Scalars['Time']['output'];
  invoiceUrl: Scalars['String']['output'];
  issued: Scalars['Time']['output'];
  metadata: Metadata;
  note?: Maybe<Scalars['String']['output']>;
  offCycle: Scalars['Boolean']['output'];
  organization: Organization;
  paid: Scalars['Boolean']['output'];
  paymentLink?: Maybe<Scalars['String']['output']>;
  postpaid: Scalars['Boolean']['output'];
  preview: Scalars['Boolean']['output'];
  provider: InvoiceProvider;
  repositoryFileId: Scalars['String']['output'];
  status?: Maybe<InvoiceStatus>;
  subtotal: Scalars['Float']['output'];
  taxDue: Scalars['Float']['output'];
};

export type InvoiceCustomer = {
  __typename?: 'InvoiceCustomer';
  addressCountry?: Maybe<Scalars['String']['output']>;
  addressLine1?: Maybe<Scalars['String']['output']>;
  addressLine2?: Maybe<Scalars['String']['output']>;
  addressLocality?: Maybe<Scalars['String']['output']>;
  addressRegion?: Maybe<Scalars['String']['output']>;
  addressZip?: Maybe<Scalars['String']['output']>;
  email?: Maybe<Scalars['String']['output']>;
  name?: Maybe<Scalars['String']['output']>;
};

export type InvoiceLine = MetadataInterface & {
  __typename?: 'InvoiceLine';
  contractLineItem: ServiceLineItem;
  description: Scalars['String']['output'];
  metadata: Metadata;
  price: Scalars['Float']['output'];
  quantity: Scalars['Int64']['output'];
  subtotal: Scalars['Float']['output'];
  taxDue: Scalars['Float']['output'];
  total: Scalars['Float']['output'];
};

export type InvoiceLineSimulate = {
  __typename?: 'InvoiceLineSimulate';
  description: Scalars['String']['output'];
  key: Scalars['String']['output'];
  price: Scalars['Float']['output'];
  quantity: Scalars['Int64']['output'];
  subtotal: Scalars['Float']['output'];
  taxDue: Scalars['Float']['output'];
  total: Scalars['Float']['output'];
};

export type InvoiceProvider = {
  __typename?: 'InvoiceProvider';
  addressCountry?: Maybe<Scalars['String']['output']>;
  addressLine1?: Maybe<Scalars['String']['output']>;
  addressLine2?: Maybe<Scalars['String']['output']>;
  addressLocality?: Maybe<Scalars['String']['output']>;
  addressRegion?: Maybe<Scalars['String']['output']>;
  addressZip?: Maybe<Scalars['String']['output']>;
  logoRepositoryFileId?: Maybe<Scalars['String']['output']>;
  logoUrl?: Maybe<Scalars['String']['output']>;
  name?: Maybe<Scalars['String']['output']>;
};

export type InvoiceSimulate = {
  __typename?: 'InvoiceSimulate';
  amount: Scalars['Float']['output'];
  currency: Scalars['String']['output'];
  customer: InvoiceCustomer;
  due: Scalars['Time']['output'];
  invoiceLineItems: Array<InvoiceLineSimulate>;
  invoiceNumber: Scalars['String']['output'];
  invoicePeriodEnd: Scalars['Time']['output'];
  invoicePeriodStart: Scalars['Time']['output'];
  issued: Scalars['Time']['output'];
  note: Scalars['String']['output'];
  offCycle: Scalars['Boolean']['output'];
  postpaid: Scalars['Boolean']['output'];
  provider: InvoiceProvider;
  subtotal: Scalars['Float']['output'];
  taxDue: Scalars['Float']['output'];
  total: Scalars['Float']['output'];
};

export type InvoiceSimulateInput = {
  contractId: Scalars['ID']['input'];
  serviceLines: Array<InvoiceSimulateServiceLineInput>;
};

export type InvoiceSimulateServiceLineInput = {
  billingCycle: BilledType;
  closeVersion?: InputMaybe<Scalars['Boolean']['input']>;
  description: Scalars['String']['input'];
  key: Scalars['String']['input'];
  parentId?: InputMaybe<Scalars['ID']['input']>;
  price: Scalars['Float']['input'];
  quantity: Scalars['Int64']['input'];
  serviceLineItemId?: InputMaybe<Scalars['ID']['input']>;
  serviceStarted: Scalars['Time']['input'];
  taxRate?: InputMaybe<Scalars['Float']['input']>;
};

export enum InvoiceStatus {
  /**
   * Deprecated, replaced by INITIALIZED
   * @deprecated use INITIALIZED instead
   */
  Draft = 'DRAFT',
  Due = 'DUE',
  Empty = 'EMPTY',
  Initialized = 'INITIALIZED',
  OnHold = 'ON_HOLD',
  Overdue = 'OVERDUE',
  Paid = 'PAID',
  Scheduled = 'SCHEDULED',
  Void = 'VOID',
}

export type InvoiceUpdateInput = {
  id: Scalars['ID']['input'];
  patch: Scalars['Boolean']['input'];
  status?: InputMaybe<InvoiceStatus>;
};

export type InvoicesPage = Pages & {
  __typename?: 'InvoicesPage';
  content: Array<Invoice>;
  totalAvailable: Scalars['Int64']['output'];
  totalElements: Scalars['Int64']['output'];
  totalPages: Scalars['Int']['output'];
};

export type Issue = Node &
  SourceFields & {
    __typename?: 'Issue';
    appSource: Scalars['String']['output'];
    assignedTo: Array<IssueParticipant>;
    comments: Array<Comment>;
    createdAt: Scalars['Time']['output'];
    description?: Maybe<Scalars['String']['output']>;
    externalLinks: Array<ExternalSystem>;
    followedBy: Array<IssueParticipant>;
    id: Scalars['ID']['output'];
    interactionEvents: Array<InteractionEvent>;
    issueStatus: Scalars['String']['output'];
    priority?: Maybe<Scalars['String']['output']>;
    reportedBy?: Maybe<IssueParticipant>;
    source: DataSource;
    sourceOfTruth: DataSource;
    /**
     * Deprecated: Use issueStatus field instead
     * @deprecated Use issueStatus field instead
     */
    status: Scalars['String']['output'];
    subject?: Maybe<Scalars['String']['output']>;
    submittedBy?: Maybe<IssueParticipant>;
    tags?: Maybe<Array<Maybe<Tag>>>;
    updatedAt: Scalars['Time']['output'];
  };

export type IssueParticipant =
  | ContactParticipant
  | OrganizationParticipant
  | UserParticipant;

export type IssueSummaryByStatus = {
  __typename?: 'IssueSummaryByStatus';
  count: Scalars['Int64']['output'];
  status: Scalars['String']['output'];
};

/**
 * Describes the relationship a Contact has with a Organization.
 * **A `return` object**
 */
export type JobRole = {
  __typename?: 'JobRole';
  appSource: Scalars['String']['output'];
  company?: Maybe<Scalars['String']['output']>;
  contact?: Maybe<Contact>;
  createdAt: Scalars['Time']['output'];
  description?: Maybe<Scalars['String']['output']>;
  endedAt?: Maybe<Scalars['Time']['output']>;
  id: Scalars['ID']['output'];
  /** The Contact's job title. */
  jobTitle?: Maybe<Scalars['String']['output']>;
  /**
   * Organization associated with a Contact.
   * **Required.**
   */
  organization?: Maybe<Organization>;
  primary: Scalars['Boolean']['output'];
  source: DataSource;
  startedAt?: Maybe<Scalars['Time']['output']>;
  updatedAt: Scalars['Time']['output'];
};

/**
 * Describes the relationship a Contact has with an Organization.
 * **A `create` object**
 */
export type JobRoleInput = {
  appSource?: InputMaybe<Scalars['String']['input']>;
  company?: InputMaybe<Scalars['String']['input']>;
  description?: InputMaybe<Scalars['String']['input']>;
  endedAt?: InputMaybe<Scalars['Time']['input']>;
  jobTitle?: InputMaybe<Scalars['String']['input']>;
  organizationId?: InputMaybe<Scalars['ID']['input']>;
  primary?: InputMaybe<Scalars['Boolean']['input']>;
  startedAt?: InputMaybe<Scalars['Time']['input']>;
};

export type JobRoleParticipant = {
  __typename?: 'JobRoleParticipant';
  jobRoleParticipant: JobRole;
  type?: Maybe<Scalars['String']['output']>;
};

/**
 * Describes the relationship a Contact has with an Organization.
 * **A `create` object**
 */
export type JobRoleUpdateInput = {
  company?: InputMaybe<Scalars['String']['input']>;
  description?: InputMaybe<Scalars['String']['input']>;
  endedAt?: InputMaybe<Scalars['Time']['input']>;
  id: Scalars['ID']['input'];
  jobTitle?: InputMaybe<Scalars['String']['input']>;
  organizationId?: InputMaybe<Scalars['ID']['input']>;
  primary?: InputMaybe<Scalars['Boolean']['input']>;
  startedAt?: InputMaybe<Scalars['Time']['input']>;
};

export type LastTouchpoint = {
  __typename?: 'LastTouchpoint';
  lastTouchPointAt?: Maybe<Scalars['Time']['output']>;
  lastTouchPointTimelineEvent?: Maybe<TimelineEvent>;
  lastTouchPointTimelineEventId?: Maybe<Scalars['ID']['output']>;
  lastTouchPointType?: Maybe<LastTouchpointType>;
};

export enum LastTouchpointType {
  Action = 'ACTION',
  ActionCreated = 'ACTION_CREATED',
  InteractionEventChat = 'INTERACTION_EVENT_CHAT',
  InteractionEventEmailReceived = 'INTERACTION_EVENT_EMAIL_RECEIVED',
  InteractionEventEmailSent = 'INTERACTION_EVENT_EMAIL_SENT',
  InteractionEventPhoneCall = 'INTERACTION_EVENT_PHONE_CALL',
  InteractionSession = 'INTERACTION_SESSION',
  IssueCreated = 'ISSUE_CREATED',
  IssueUpdated = 'ISSUE_UPDATED',
  LogEntry = 'LOG_ENTRY',
  Meeting = 'MEETING',
  Note = 'NOTE',
  PageView = 'PAGE_VIEW',
}

export type LinkOrganizationsInput = {
  organizationId: Scalars['ID']['input'];
  removeExisting?: InputMaybe<Scalars['Boolean']['input']>;
  subsidiaryId: Scalars['ID']['input'];
  type?: InputMaybe<Scalars['String']['input']>;
};

export type LinkedOrganization = {
  __typename?: 'LinkedOrganization';
  organization: Organization;
  type?: Maybe<Scalars['String']['output']>;
};

export type Location = Node &
  SourceFields & {
    __typename?: 'Location';
    address?: Maybe<Scalars['String']['output']>;
    address2?: Maybe<Scalars['String']['output']>;
    addressType?: Maybe<Scalars['String']['output']>;
    appSource: Scalars['String']['output'];
    commercial?: Maybe<Scalars['Boolean']['output']>;
    country?: Maybe<Scalars['String']['output']>;
    countryCodeA2?: Maybe<Scalars['String']['output']>;
    countryCodeA3?: Maybe<Scalars['String']['output']>;
    createdAt: Scalars['Time']['output'];
    district?: Maybe<Scalars['String']['output']>;
    houseNumber?: Maybe<Scalars['String']['output']>;
    id: Scalars['ID']['output'];
    latitude?: Maybe<Scalars['Float']['output']>;
    locality?: Maybe<Scalars['String']['output']>;
    longitude?: Maybe<Scalars['Float']['output']>;
    name?: Maybe<Scalars['String']['output']>;
    plusFour?: Maybe<Scalars['String']['output']>;
    postalCode?: Maybe<Scalars['String']['output']>;
    predirection?: Maybe<Scalars['String']['output']>;
    rawAddress?: Maybe<Scalars['String']['output']>;
    region?: Maybe<Scalars['String']['output']>;
    source: DataSource;
    sourceOfTruth: DataSource;
    street?: Maybe<Scalars['String']['output']>;
    timeZone?: Maybe<Scalars['String']['output']>;
    updatedAt: Scalars['Time']['output'];
    utcOffset?: Maybe<Scalars['Float']['output']>;
    zip?: Maybe<Scalars['String']['output']>;
  };

export type LocationUpdateInput = {
  address?: InputMaybe<Scalars['String']['input']>;
  address2?: InputMaybe<Scalars['String']['input']>;
  addressType?: InputMaybe<Scalars['String']['input']>;
  commercial?: InputMaybe<Scalars['Boolean']['input']>;
  country?: InputMaybe<Scalars['String']['input']>;
  district?: InputMaybe<Scalars['String']['input']>;
  houseNumber?: InputMaybe<Scalars['String']['input']>;
  id: Scalars['ID']['input'];
  latitude?: InputMaybe<Scalars['Float']['input']>;
  locality?: InputMaybe<Scalars['String']['input']>;
  longitude?: InputMaybe<Scalars['Float']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  plusFour?: InputMaybe<Scalars['String']['input']>;
  postalCode?: InputMaybe<Scalars['String']['input']>;
  predirection?: InputMaybe<Scalars['String']['input']>;
  rawAddress?: InputMaybe<Scalars['String']['input']>;
  region?: InputMaybe<Scalars['String']['input']>;
  street?: InputMaybe<Scalars['String']['input']>;
  timeZone?: InputMaybe<Scalars['String']['input']>;
  utcOffset?: InputMaybe<Scalars['Float']['input']>;
  zip?: InputMaybe<Scalars['String']['input']>;
};

export type LogEntry = {
  __typename?: 'LogEntry';
  appSource: Scalars['String']['output'];
  content?: Maybe<Scalars['String']['output']>;
  contentType?: Maybe<Scalars['String']['output']>;
  createdAt: Scalars['Time']['output'];
  createdBy?: Maybe<User>;
  externalLinks: Array<ExternalSystem>;
  id: Scalars['ID']['output'];
  source: DataSource;
  sourceOfTruth: DataSource;
  startedAt: Scalars['Time']['output'];
  tags: Array<Tag>;
  updatedAt: Scalars['Time']['output'];
};

export type LogEntryInput = {
  appSource?: InputMaybe<Scalars['String']['input']>;
  content?: InputMaybe<Scalars['String']['input']>;
  contentType?: InputMaybe<Scalars['String']['input']>;
  startedAt?: InputMaybe<Scalars['Time']['input']>;
  tags?: InputMaybe<Array<TagIdOrNameInput>>;
};

export type LogEntryUpdateInput = {
  content?: InputMaybe<Scalars['String']['input']>;
  contentType?: InputMaybe<Scalars['String']['input']>;
  startedAt?: InputMaybe<Scalars['Time']['input']>;
};

export enum Market {
  B2B = 'B2B',
  B2C = 'B2C',
  Marketplace = 'MARKETPLACE',
}

export type Meeting = Node & {
  __typename?: 'Meeting';
  agenda?: Maybe<Scalars['String']['output']>;
  agendaContentType?: Maybe<Scalars['String']['output']>;
  appSource: Scalars['String']['output'];
  attendedBy: Array<MeetingParticipant>;
  conferenceUrl?: Maybe<Scalars['String']['output']>;
  createdAt: Scalars['Time']['output'];
  createdBy: Array<MeetingParticipant>;
  endedAt?: Maybe<Scalars['Time']['output']>;
  events: Array<InteractionEvent>;
  externalSystem: Array<ExternalSystem>;
  id: Scalars['ID']['output'];
  includes: Array<Attachment>;
  meetingExternalUrl?: Maybe<Scalars['String']['output']>;
  name?: Maybe<Scalars['String']['output']>;
  note: Array<Note>;
  recording?: Maybe<Attachment>;
  source: DataSource;
  sourceOfTruth: DataSource;
  startedAt?: Maybe<Scalars['Time']['output']>;
  status: MeetingStatus;
  updatedAt: Scalars['Time']['output'];
};

export type MeetingInput = {
  agenda?: InputMaybe<Scalars['String']['input']>;
  agendaContentType?: InputMaybe<Scalars['String']['input']>;
  appSource?: InputMaybe<Scalars['String']['input']>;
  attendedBy?: InputMaybe<Array<MeetingParticipantInput>>;
  conferenceUrl?: InputMaybe<Scalars['String']['input']>;
  createdAt?: InputMaybe<Scalars['Time']['input']>;
  createdBy?: InputMaybe<Array<MeetingParticipantInput>>;
  endedAt?: InputMaybe<Scalars['Time']['input']>;
  externalSystem?: InputMaybe<ExternalSystemReferenceInput>;
  meetingExternalUrl?: InputMaybe<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  note?: InputMaybe<NoteInput>;
  startedAt?: InputMaybe<Scalars['Time']['input']>;
  status?: InputMaybe<MeetingStatus>;
};

export type MeetingParticipant =
  | ContactParticipant
  | EmailParticipant
  | OrganizationParticipant
  | UserParticipant;

export type MeetingParticipantInput = {
  contactId?: InputMaybe<Scalars['ID']['input']>;
  organizationId?: InputMaybe<Scalars['ID']['input']>;
  userId?: InputMaybe<Scalars['ID']['input']>;
};

export enum MeetingStatus {
  Accepted = 'ACCEPTED',
  Canceled = 'CANCELED',
  Undefined = 'UNDEFINED',
}

export type MeetingUpdateInput = {
  agenda?: InputMaybe<Scalars['String']['input']>;
  agendaContentType?: InputMaybe<Scalars['String']['input']>;
  appSource?: InputMaybe<Scalars['String']['input']>;
  conferenceUrl?: InputMaybe<Scalars['String']['input']>;
  endedAt?: InputMaybe<Scalars['Time']['input']>;
  externalSystem?: InputMaybe<ExternalSystemReferenceInput>;
  meetingExternalUrl?: InputMaybe<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  note?: InputMaybe<NoteUpdateInput>;
  startedAt?: InputMaybe<Scalars['Time']['input']>;
  status?: InputMaybe<MeetingStatus>;
};

/**
 * Specifies how many pages of meeting information has been returned in the query response.
 * **A `response` object.**
 */
export type MeetingsPage = Pages & {
  __typename?: 'MeetingsPage';
  /**
   * A contact entity in customerOS.
   * **Required.  If no values it returns an empty array.**
   */
  content: Array<Meeting>;
  /**
   * Total number of elements in the query response.
   * **Required.**
   */
  totalElements: Scalars['Int64']['output'];
  /**
   * Total number of pages in the query response.
   * **Required.**
   */
  totalPages: Scalars['Int']['output'];
};

export type Metadata = Node &
  SourceFieldsInterface & {
    __typename?: 'Metadata';
    appSource: Scalars['String']['output'];
    created: Scalars['Time']['output'];
    id: Scalars['ID']['output'];
    lastUpdated: Scalars['Time']['output'];
    source: DataSource;
    sourceOfTruth: DataSource;
    /** Aggregate version from event store db */
    version?: Maybe<Scalars['Int64']['output']>;
  };

export type MetadataInterface = {
  metadata: Metadata;
};

export type Mutation = {
  __typename?: 'Mutation';
  addTag: Scalars['ID']['output'];
  attachment_Create: Attachment;
  bankAccount_Create: BankAccount;
  bankAccount_Delete: DeleteResponse;
  bankAccount_Update: BankAccount;
  billingProfile_Create: Scalars['ID']['output'];
  billingProfile_LinkEmail: Scalars['ID']['output'];
  billingProfile_LinkLocation: Scalars['ID']['output'];
  billingProfile_UnlinkEmail: Scalars['ID']['output'];
  billingProfile_UnlinkLocation: Scalars['ID']['output'];
  billingProfile_Update: Scalars['ID']['output'];
  contact_AddNewLocation: Location;
  contact_AddOrganizationById: Contact;
  contact_AddSocial: Social;
  contact_AddTag: ActionResponse;
  /** @deprecated Use contact_Hide instead */
  contact_Archive: Result;
  contact_Create: Scalars['ID']['output'];
  contact_CreateForOrganization: Contact;
  contact_FindWorkEmail: ActionResponse;
  contact_HardDelete: Result;
  contact_Hide: ActionResponse;
  contact_Merge: Contact;
  contact_RemoveLocation: Contact;
  contact_RemoveOrganizationById: Contact;
  contact_RemoveSocial: ActionResponse;
  contact_RemoveTag: ActionResponse;
  contact_RestoreFromArchive: Result;
  contact_Update: Contact;
  contractLineItem_Close: Scalars['ID']['output'];
  contractLineItem_Create: ServiceLineItem;
  contractLineItem_NewVersion: ServiceLineItem;
  contractLineItem_Pause: ActionResponse;
  contractLineItem_Resume: ActionResponse;
  contractLineItem_Update: ServiceLineItem;
  contract_AddAttachment: Contract;
  contract_Create: Contract;
  contract_Delete: DeleteResponse;
  contract_RemoveAttachment: Contract;
  contract_Renew: Contract;
  contract_Update: Contract;
  customFieldDeleteFromContactById: Result;
  customFieldDeleteFromContactByName: Result;
  customFieldMergeToContact: CustomField;
  customFieldTemplate_Delete?: Maybe<Scalars['Boolean']['output']>;
  customFieldTemplate_Save: CustomFieldTemplate;
  customFieldUpdateInContact: CustomField;
  customFieldsMergeAndUpdateInContact: Contact;
  customer_contact_Create: CustomerContact;
  customer_user_AddJobRole: CustomerUser;
  emailMergeToContact: Email;
  emailMergeToOrganization: Email;
  emailMergeToUser: Email;
  emailRemoveFromContact: Result;
  emailRemoveFromOrganization: Result;
  emailRemoveFromUser: Result;
  emailReplaceForContact: Email;
  emailReplaceForOrganization: Email;
  emailReplaceForUser: Email;
  emailSetPrimaryForContact: Result;
  email_Validate: ActionResponse;
  externalSystem_Create: Scalars['ID']['output'];
  flowContact_Add: FlowContact;
  flowContact_AddBulk: Result;
  flowContact_Delete: Result;
  flowContact_DeleteBulk: Result;
  flowSender_Delete: Result;
  flowSender_Merge: FlowSender;
  flow_ChangeStatus: Flow;
  flow_Merge: Flow;
  interactionEvent_LinkAttachment: Result;
  invoice_NextDryRunForContract: Scalars['ID']['output'];
  invoice_Pay: Invoice;
  invoice_Simulate: Array<InvoiceSimulate>;
  invoice_Update: Invoice;
  invoice_Void: Invoice;
  jobRole_Create: JobRole;
  jobRole_Delete: Result;
  jobRole_Update: JobRole;
  location_RemoveFromContact: Contact;
  location_RemoveFromOrganization: Organization;
  /** @deprecated No longer supported */
  location_Update: Location;
  logEntry_AddTag: Scalars['ID']['output'];
  logEntry_CreateForOrganization: Scalars['ID']['output'];
  logEntry_RemoveTag: Scalars['ID']['output'];
  logEntry_ResetTags: Scalars['ID']['output'];
  logEntry_Update: Scalars['ID']['output'];
  meeting_AddNewLocation: Meeting;
  meeting_AddNote: Meeting;
  meeting_Create: Meeting;
  meeting_LinkAttachment: Meeting;
  meeting_LinkAttendedBy: Meeting;
  meeting_LinkRecording: Meeting;
  meeting_UnlinkAttachment: Meeting;
  meeting_UnlinkAttendedBy: Meeting;
  meeting_UnlinkRecording: Meeting;
  meeting_Update: Meeting;
  note_Delete: Result;
  note_LinkAttachment: Note;
  note_UnlinkAttachment: Note;
  note_Update: Note;
  opportunityRenewalUpdate: Opportunity;
  opportunityRenewal_UpdateAllForOrganization: Organization;
  opportunity_Archive: ActionResponse;
  opportunity_Create: Opportunity;
  opportunity_RemoveOwner: ActionResponse;
  opportunity_Save: Opportunity;
  opportunity_SetOwner: ActionResponse;
  opportunity_Update: Opportunity;
  organization_AddNewLocation: Location;
  organization_AddSocial: Social;
  organization_AddSubsidiary: Organization;
  /** @deprecated No longer supported */
  organization_AddTag: ActionResponse;
  organization_Archive?: Maybe<Result>;
  organization_ArchiveAll?: Maybe<Result>;
  organization_Hide: Scalars['ID']['output'];
  organization_HideAll?: Maybe<Result>;
  organization_Merge: Organization;
  organization_RemoveSocial: ActionResponse;
  organization_RemoveSubsidiary: Organization;
  /** @deprecated No longer supported */
  organization_RemoveTag: ActionResponse;
  organization_Save: Organization;
  /** @deprecated No longer supported */
  organization_SetOwner: Organization;
  organization_Show: Scalars['ID']['output'];
  organization_ShowAll?: Maybe<Result>;
  organization_UnlinkAllDomains: Organization;
  /** @deprecated No longer supported */
  organization_UnsetOwner: Organization;
  /** @deprecated No longer supported */
  organization_Update: Organization;
  organization_UpdateOnboardingStatus: Organization;
  phoneNumberMergeToContact: PhoneNumber;
  phoneNumberMergeToOrganization: PhoneNumber;
  phoneNumberMergeToUser: PhoneNumber;
  phoneNumberRemoveFromContactByE164: Result;
  phoneNumberRemoveFromContactById: Result;
  phoneNumberRemoveFromOrganizationByE164: Result;
  phoneNumberRemoveFromOrganizationById: Result;
  phoneNumberRemoveFromUserByE164: Result;
  phoneNumberRemoveFromUserById: Result;
  phoneNumberUpdateInContact: PhoneNumber;
  phoneNumberUpdateInOrganization: PhoneNumber;
  phoneNumberUpdateInUser: PhoneNumber;
  phoneNumber_Update: PhoneNumber;
  reminder_Create?: Maybe<Scalars['ID']['output']>;
  reminder_Update?: Maybe<Scalars['ID']['output']>;
  removeTag?: Maybe<Result>;
  serviceLineItem_BulkUpdate: Array<Scalars['ID']['output']>;
  serviceLineItem_Delete: DeleteResponse;
  social_Remove: Result;
  social_Update: Social;
  tableViewDef_Archive: ActionResponse;
  tableViewDef_Create: TableViewDef;
  tableViewDef_Update: TableViewDef;
  tableViewDef_UpdateShared: TableViewDef;
  tag_Create: Tag;
  tag_Delete?: Maybe<Result>;
  tag_Update?: Maybe<Tag>;
  tenant_AddBillingProfile: TenantBillingProfile;
  tenant_UpdateBillingProfile: TenantBillingProfile;
  tenant_UpdateSettings: TenantSettings;
  tenant_UpdateSettingsOpportunityStage: ActionResponse;
  tenant_hardDelete: Scalars['Boolean']['output'];
  user_AddRole: User;
  user_AddRoleInTenant: User;
  user_Create: User;
  user_Delete: Result;
  user_DeleteInTenant: Result;
  user_RemoveRole: User;
  user_RemoveRoleInTenant: User;
  user_Update: User;
  workflow_Create: Workflow;
  workflow_Update: ActionResponse;
};

export type MutationAddTagArgs = {
  input: AddTagInput;
};

export type MutationAttachment_CreateArgs = {
  input: AttachmentInput;
};

export type MutationBankAccount_CreateArgs = {
  input?: InputMaybe<BankAccountCreateInput>;
};

export type MutationBankAccount_DeleteArgs = {
  id: Scalars['ID']['input'];
};

export type MutationBankAccount_UpdateArgs = {
  input?: InputMaybe<BankAccountUpdateInput>;
};

export type MutationBillingProfile_CreateArgs = {
  input: BillingProfileInput;
};

export type MutationBillingProfile_LinkEmailArgs = {
  input: BillingProfileLinkEmailInput;
};

export type MutationBillingProfile_LinkLocationArgs = {
  input: BillingProfileLinkLocationInput;
};

export type MutationBillingProfile_UnlinkEmailArgs = {
  input: BillingProfileLinkEmailInput;
};

export type MutationBillingProfile_UnlinkLocationArgs = {
  input: BillingProfileLinkLocationInput;
};

export type MutationBillingProfile_UpdateArgs = {
  input: BillingProfileUpdateInput;
};

export type MutationContact_AddNewLocationArgs = {
  contactId: Scalars['ID']['input'];
};

export type MutationContact_AddOrganizationByIdArgs = {
  input: ContactOrganizationInput;
};

export type MutationContact_AddSocialArgs = {
  contactId: Scalars['ID']['input'];
  input: SocialInput;
};

export type MutationContact_AddTagArgs = {
  input: ContactTagInput;
};

export type MutationContact_ArchiveArgs = {
  contactId: Scalars['ID']['input'];
};

export type MutationContact_CreateArgs = {
  input: ContactInput;
};

export type MutationContact_CreateForOrganizationArgs = {
  input: ContactInput;
  organizationId: Scalars['ID']['input'];
};

export type MutationContact_FindWorkEmailArgs = {
  contactId: Scalars['ID']['input'];
  domain?: InputMaybe<Scalars['String']['input']>;
  findMobileNumber?: InputMaybe<Scalars['Boolean']['input']>;
  organizationId?: InputMaybe<Scalars['ID']['input']>;
};

export type MutationContact_HardDeleteArgs = {
  contactId: Scalars['ID']['input'];
};

export type MutationContact_HideArgs = {
  contactId: Scalars['ID']['input'];
};

export type MutationContact_MergeArgs = {
  mergedContactIds: Array<Scalars['ID']['input']>;
  primaryContactId: Scalars['ID']['input'];
};

export type MutationContact_RemoveLocationArgs = {
  contactId: Scalars['ID']['input'];
  locationId: Scalars['ID']['input'];
};

export type MutationContact_RemoveOrganizationByIdArgs = {
  input: ContactOrganizationInput;
};

export type MutationContact_RemoveSocialArgs = {
  contactId: Scalars['ID']['input'];
  socialId: Scalars['ID']['input'];
};

export type MutationContact_RemoveTagArgs = {
  input: ContactTagInput;
};

export type MutationContact_RestoreFromArchiveArgs = {
  contactId: Scalars['ID']['input'];
};

export type MutationContact_UpdateArgs = {
  input: ContactUpdateInput;
};

export type MutationContractLineItem_CloseArgs = {
  input: ServiceLineItemCloseInput;
};

export type MutationContractLineItem_CreateArgs = {
  input: ServiceLineItemInput;
};

export type MutationContractLineItem_NewVersionArgs = {
  input: ServiceLineItemNewVersionInput;
};

export type MutationContractLineItem_PauseArgs = {
  id: Scalars['ID']['input'];
};

export type MutationContractLineItem_ResumeArgs = {
  id: Scalars['ID']['input'];
};

export type MutationContractLineItem_UpdateArgs = {
  input: ServiceLineItemUpdateInput;
};

export type MutationContract_AddAttachmentArgs = {
  attachmentId: Scalars['ID']['input'];
  contractId: Scalars['ID']['input'];
};

export type MutationContract_CreateArgs = {
  input: ContractInput;
};

export type MutationContract_DeleteArgs = {
  id: Scalars['ID']['input'];
};

export type MutationContract_RemoveAttachmentArgs = {
  attachmentId: Scalars['ID']['input'];
  contractId: Scalars['ID']['input'];
};

export type MutationContract_RenewArgs = {
  input: ContractRenewalInput;
};

export type MutationContract_UpdateArgs = {
  input: ContractUpdateInput;
};

export type MutationCustomFieldDeleteFromContactByIdArgs = {
  contactId: Scalars['ID']['input'];
  id: Scalars['ID']['input'];
};

export type MutationCustomFieldDeleteFromContactByNameArgs = {
  contactId: Scalars['ID']['input'];
  fieldName: Scalars['String']['input'];
};

export type MutationCustomFieldMergeToContactArgs = {
  contactId: Scalars['ID']['input'];
  input: CustomFieldInput;
};

export type MutationCustomFieldTemplate_DeleteArgs = {
  id: Scalars['ID']['input'];
};

export type MutationCustomFieldTemplate_SaveArgs = {
  input: CustomFieldTemplateInput;
};

export type MutationCustomFieldUpdateInContactArgs = {
  contactId: Scalars['ID']['input'];
  input: CustomFieldUpdateInput;
};

export type MutationCustomFieldsMergeAndUpdateInContactArgs = {
  contactId: Scalars['ID']['input'];
  customFields?: InputMaybe<Array<CustomFieldInput>>;
};

export type MutationCustomer_Contact_CreateArgs = {
  input: CustomerContactInput;
};

export type MutationCustomer_User_AddJobRoleArgs = {
  id: Scalars['ID']['input'];
  jobRoleInput: JobRoleInput;
};

export type MutationEmailMergeToContactArgs = {
  contactId: Scalars['ID']['input'];
  input: EmailInput;
};

export type MutationEmailMergeToOrganizationArgs = {
  input: EmailInput;
  organizationId: Scalars['ID']['input'];
};

export type MutationEmailMergeToUserArgs = {
  input: EmailInput;
  userId: Scalars['ID']['input'];
};

export type MutationEmailRemoveFromContactArgs = {
  contactId: Scalars['ID']['input'];
  email: Scalars['String']['input'];
};

export type MutationEmailRemoveFromOrganizationArgs = {
  email: Scalars['String']['input'];
  organizationId: Scalars['ID']['input'];
};

export type MutationEmailRemoveFromUserArgs = {
  email: Scalars['String']['input'];
  userId: Scalars['ID']['input'];
};

export type MutationEmailReplaceForContactArgs = {
  contactId: Scalars['ID']['input'];
  input: EmailInput;
  previousEmail?: InputMaybe<Scalars['String']['input']>;
};

export type MutationEmailReplaceForOrganizationArgs = {
  input: EmailInput;
  organizationId: Scalars['ID']['input'];
  previousEmail?: InputMaybe<Scalars['String']['input']>;
};

export type MutationEmailReplaceForUserArgs = {
  input: EmailInput;
  previousEmail?: InputMaybe<Scalars['String']['input']>;
  userId: Scalars['ID']['input'];
};

export type MutationEmailSetPrimaryForContactArgs = {
  contactId: Scalars['ID']['input'];
  email: Scalars['String']['input'];
};

export type MutationEmail_ValidateArgs = {
  id: Scalars['ID']['input'];
};

export type MutationExternalSystem_CreateArgs = {
  input: ExternalSystemInput;
};

export type MutationFlowContact_AddArgs = {
  contactId: Scalars['ID']['input'];
  flowId: Scalars['ID']['input'];
};

export type MutationFlowContact_AddBulkArgs = {
  contactId: Array<Scalars['ID']['input']>;
  flowId: Scalars['ID']['input'];
};

export type MutationFlowContact_DeleteArgs = {
  id: Scalars['ID']['input'];
};

export type MutationFlowContact_DeleteBulkArgs = {
  id: Array<Scalars['ID']['input']>;
};

export type MutationFlowSender_DeleteArgs = {
  id: Scalars['ID']['input'];
};

export type MutationFlowSender_MergeArgs = {
  flowId: Scalars['ID']['input'];
  input: FlowSenderMergeInput;
};

export type MutationFlow_ChangeStatusArgs = {
  id: Scalars['ID']['input'];
  status: FlowStatus;
};

export type MutationFlow_MergeArgs = {
  input: FlowMergeInput;
};

export type MutationInteractionEvent_LinkAttachmentArgs = {
  attachmentId: Scalars['ID']['input'];
  eventId: Scalars['ID']['input'];
};

export type MutationInvoice_NextDryRunForContractArgs = {
  contractId: Scalars['ID']['input'];
};

export type MutationInvoice_PayArgs = {
  id: Scalars['ID']['input'];
};

export type MutationInvoice_SimulateArgs = {
  input: InvoiceSimulateInput;
};

export type MutationInvoice_UpdateArgs = {
  input: InvoiceUpdateInput;
};

export type MutationInvoice_VoidArgs = {
  id: Scalars['ID']['input'];
};

export type MutationJobRole_CreateArgs = {
  contactId: Scalars['ID']['input'];
  input: JobRoleInput;
};

export type MutationJobRole_DeleteArgs = {
  contactId: Scalars['ID']['input'];
  roleId: Scalars['ID']['input'];
};

export type MutationJobRole_UpdateArgs = {
  contactId: Scalars['ID']['input'];
  input: JobRoleUpdateInput;
};

export type MutationLocation_RemoveFromContactArgs = {
  contactId: Scalars['ID']['input'];
  locationId: Scalars['ID']['input'];
};

export type MutationLocation_RemoveFromOrganizationArgs = {
  locationId: Scalars['ID']['input'];
  organizationId: Scalars['ID']['input'];
};

export type MutationLocation_UpdateArgs = {
  input: LocationUpdateInput;
};

export type MutationLogEntry_AddTagArgs = {
  id: Scalars['ID']['input'];
  input: TagIdOrNameInput;
};

export type MutationLogEntry_CreateForOrganizationArgs = {
  input: LogEntryInput;
  organizationId: Scalars['ID']['input'];
};

export type MutationLogEntry_RemoveTagArgs = {
  id: Scalars['ID']['input'];
  input: TagIdOrNameInput;
};

export type MutationLogEntry_ResetTagsArgs = {
  id: Scalars['ID']['input'];
  input?: InputMaybe<Array<TagIdOrNameInput>>;
};

export type MutationLogEntry_UpdateArgs = {
  id: Scalars['ID']['input'];
  input: LogEntryUpdateInput;
};

export type MutationMeeting_AddNewLocationArgs = {
  meetingId: Scalars['ID']['input'];
};

export type MutationMeeting_AddNoteArgs = {
  meetingId: Scalars['ID']['input'];
  note?: InputMaybe<NoteInput>;
};

export type MutationMeeting_CreateArgs = {
  meeting: MeetingInput;
};

export type MutationMeeting_LinkAttachmentArgs = {
  attachmentId: Scalars['ID']['input'];
  meetingId: Scalars['ID']['input'];
};

export type MutationMeeting_LinkAttendedByArgs = {
  meetingId: Scalars['ID']['input'];
  participant: MeetingParticipantInput;
};

export type MutationMeeting_LinkRecordingArgs = {
  attachmentId: Scalars['ID']['input'];
  meetingId: Scalars['ID']['input'];
};

export type MutationMeeting_UnlinkAttachmentArgs = {
  attachmentId: Scalars['ID']['input'];
  meetingId: Scalars['ID']['input'];
};

export type MutationMeeting_UnlinkAttendedByArgs = {
  meetingId: Scalars['ID']['input'];
  participant: MeetingParticipantInput;
};

export type MutationMeeting_UnlinkRecordingArgs = {
  attachmentId: Scalars['ID']['input'];
  meetingId: Scalars['ID']['input'];
};

export type MutationMeeting_UpdateArgs = {
  meeting: MeetingUpdateInput;
  meetingId: Scalars['ID']['input'];
};

export type MutationNote_DeleteArgs = {
  id: Scalars['ID']['input'];
};

export type MutationNote_LinkAttachmentArgs = {
  attachmentId: Scalars['ID']['input'];
  noteId: Scalars['ID']['input'];
};

export type MutationNote_UnlinkAttachmentArgs = {
  attachmentId: Scalars['ID']['input'];
  noteId: Scalars['ID']['input'];
};

export type MutationNote_UpdateArgs = {
  input: NoteUpdateInput;
};

export type MutationOpportunityRenewalUpdateArgs = {
  input: OpportunityRenewalUpdateInput;
  ownerUserId?: InputMaybe<Scalars['ID']['input']>;
};

export type MutationOpportunityRenewal_UpdateAllForOrganizationArgs = {
  input: OpportunityRenewalUpdateAllForOrganizationInput;
};

export type MutationOpportunity_ArchiveArgs = {
  id: Scalars['ID']['input'];
};

export type MutationOpportunity_CreateArgs = {
  input: OpportunityCreateInput;
};

export type MutationOpportunity_RemoveOwnerArgs = {
  opportunityId: Scalars['ID']['input'];
};

export type MutationOpportunity_SaveArgs = {
  input: OpportunitySaveInput;
};

export type MutationOpportunity_SetOwnerArgs = {
  opportunityId: Scalars['ID']['input'];
  userId: Scalars['ID']['input'];
};

export type MutationOpportunity_UpdateArgs = {
  input: OpportunityUpdateInput;
};

export type MutationOrganization_AddNewLocationArgs = {
  organizationId: Scalars['ID']['input'];
};

export type MutationOrganization_AddSocialArgs = {
  input: SocialInput;
  organizationId: Scalars['ID']['input'];
};

export type MutationOrganization_AddSubsidiaryArgs = {
  input: LinkOrganizationsInput;
};

export type MutationOrganization_AddTagArgs = {
  input: OrganizationTagInput;
};

export type MutationOrganization_ArchiveArgs = {
  id: Scalars['ID']['input'];
};

export type MutationOrganization_ArchiveAllArgs = {
  ids: Array<Scalars['ID']['input']>;
};

export type MutationOrganization_HideArgs = {
  id: Scalars['ID']['input'];
};

export type MutationOrganization_HideAllArgs = {
  ids: Array<Scalars['ID']['input']>;
};

export type MutationOrganization_MergeArgs = {
  mergedOrganizationIds: Array<Scalars['ID']['input']>;
  primaryOrganizationId: Scalars['ID']['input'];
};

export type MutationOrganization_RemoveSocialArgs = {
  organizationId: Scalars['ID']['input'];
  socialId: Scalars['ID']['input'];
};

export type MutationOrganization_RemoveSubsidiaryArgs = {
  organizationId: Scalars['ID']['input'];
  subsidiaryId: Scalars['ID']['input'];
};

export type MutationOrganization_RemoveTagArgs = {
  input: OrganizationTagInput;
};

export type MutationOrganization_SaveArgs = {
  input: OrganizationSaveInput;
};

export type MutationOrganization_SetOwnerArgs = {
  organizationId: Scalars['ID']['input'];
  userId: Scalars['ID']['input'];
};

export type MutationOrganization_ShowArgs = {
  id: Scalars['ID']['input'];
};

export type MutationOrganization_ShowAllArgs = {
  ids: Array<Scalars['ID']['input']>;
};

export type MutationOrganization_UnlinkAllDomainsArgs = {
  organizationId: Scalars['ID']['input'];
};

export type MutationOrganization_UnsetOwnerArgs = {
  organizationId: Scalars['ID']['input'];
};

export type MutationOrganization_UpdateArgs = {
  input: OrganizationUpdateInput;
};

export type MutationOrganization_UpdateOnboardingStatusArgs = {
  input: OnboardingStatusInput;
};

export type MutationPhoneNumberMergeToContactArgs = {
  contactId: Scalars['ID']['input'];
  input: PhoneNumberInput;
};

export type MutationPhoneNumberMergeToOrganizationArgs = {
  input: PhoneNumberInput;
  organizationId: Scalars['ID']['input'];
};

export type MutationPhoneNumberMergeToUserArgs = {
  input: PhoneNumberInput;
  userId: Scalars['ID']['input'];
};

export type MutationPhoneNumberRemoveFromContactByE164Args = {
  contactId: Scalars['ID']['input'];
  e164: Scalars['String']['input'];
};

export type MutationPhoneNumberRemoveFromContactByIdArgs = {
  contactId: Scalars['ID']['input'];
  id: Scalars['ID']['input'];
};

export type MutationPhoneNumberRemoveFromOrganizationByE164Args = {
  e164: Scalars['String']['input'];
  organizationId: Scalars['ID']['input'];
};

export type MutationPhoneNumberRemoveFromOrganizationByIdArgs = {
  id: Scalars['ID']['input'];
  organizationId: Scalars['ID']['input'];
};

export type MutationPhoneNumberRemoveFromUserByE164Args = {
  e164: Scalars['String']['input'];
  userId: Scalars['ID']['input'];
};

export type MutationPhoneNumberRemoveFromUserByIdArgs = {
  id: Scalars['ID']['input'];
  userId: Scalars['ID']['input'];
};

export type MutationPhoneNumberUpdateInContactArgs = {
  contactId: Scalars['ID']['input'];
  input: PhoneNumberRelationUpdateInput;
};

export type MutationPhoneNumberUpdateInOrganizationArgs = {
  input: PhoneNumberRelationUpdateInput;
  organizationId: Scalars['ID']['input'];
};

export type MutationPhoneNumberUpdateInUserArgs = {
  input: PhoneNumberRelationUpdateInput;
  userId: Scalars['ID']['input'];
};

export type MutationPhoneNumber_UpdateArgs = {
  input: PhoneNumberUpdateInput;
};

export type MutationReminder_CreateArgs = {
  input: ReminderInput;
};

export type MutationReminder_UpdateArgs = {
  input: ReminderUpdateInput;
};

export type MutationRemoveTagArgs = {
  input: RemoveTagInput;
};

export type MutationServiceLineItem_BulkUpdateArgs = {
  input: ServiceLineItemBulkUpdateInput;
};

export type MutationServiceLineItem_DeleteArgs = {
  id: Scalars['ID']['input'];
};

export type MutationSocial_RemoveArgs = {
  socialId: Scalars['ID']['input'];
};

export type MutationSocial_UpdateArgs = {
  input: SocialUpdateInput;
};

export type MutationTableViewDef_ArchiveArgs = {
  id: Scalars['ID']['input'];
};

export type MutationTableViewDef_CreateArgs = {
  input: TableViewDefCreateInput;
};

export type MutationTableViewDef_UpdateArgs = {
  input: TableViewDefUpdateInput;
};

export type MutationTableViewDef_UpdateSharedArgs = {
  input: TableViewDefUpdateInput;
};

export type MutationTag_CreateArgs = {
  input: TagInput;
};

export type MutationTag_DeleteArgs = {
  id: Scalars['ID']['input'];
};

export type MutationTag_UpdateArgs = {
  input: TagUpdateInput;
};

export type MutationTenant_AddBillingProfileArgs = {
  input: TenantBillingProfileInput;
};

export type MutationTenant_UpdateBillingProfileArgs = {
  input: TenantBillingProfileUpdateInput;
};

export type MutationTenant_UpdateSettingsArgs = {
  input?: InputMaybe<TenantSettingsInput>;
};

export type MutationTenant_UpdateSettingsOpportunityStageArgs = {
  input: TenantSettingsOpportunityStageConfigurationInput;
};

export type MutationTenant_HardDeleteArgs = {
  confirmTenant: Scalars['String']['input'];
  tenant: Scalars['String']['input'];
};

export type MutationUser_AddRoleArgs = {
  id: Scalars['ID']['input'];
  role: Role;
};

export type MutationUser_AddRoleInTenantArgs = {
  id: Scalars['ID']['input'];
  role: Role;
  tenant: Scalars['String']['input'];
};

export type MutationUser_CreateArgs = {
  input: UserInput;
};

export type MutationUser_DeleteArgs = {
  id: Scalars['ID']['input'];
};

export type MutationUser_DeleteInTenantArgs = {
  id: Scalars['ID']['input'];
  tenant: Scalars['String']['input'];
};

export type MutationUser_RemoveRoleArgs = {
  id: Scalars['ID']['input'];
  role: Role;
};

export type MutationUser_RemoveRoleInTenantArgs = {
  id: Scalars['ID']['input'];
  role: Role;
  tenant: Scalars['String']['input'];
};

export type MutationUser_UpdateArgs = {
  input: UserUpdateInput;
};

export type MutationWorkflow_CreateArgs = {
  input: WorkflowCreateInput;
};

export type MutationWorkflow_UpdateArgs = {
  input: WorkflowUpdateInput;
};

export type Node = {
  id: Scalars['ID']['output'];
};

export type Note = {
  __typename?: 'Note';
  appSource: Scalars['String']['output'];
  content?: Maybe<Scalars['String']['output']>;
  contentType?: Maybe<Scalars['String']['output']>;
  createdAt: Scalars['Time']['output'];
  createdBy?: Maybe<User>;
  id: Scalars['ID']['output'];
  includes: Array<Attachment>;
  source: DataSource;
  sourceOfTruth: DataSource;
  updatedAt: Scalars['Time']['output'];
};

export type NoteInput = {
  appSource?: InputMaybe<Scalars['String']['input']>;
  content?: InputMaybe<Scalars['String']['input']>;
  contentType?: InputMaybe<Scalars['String']['input']>;
};

export type NotePage = Pages & {
  __typename?: 'NotePage';
  content: Array<Note>;
  totalElements: Scalars['Int64']['output'];
  totalPages: Scalars['Int']['output'];
};

export type NoteUpdateInput = {
  content?: InputMaybe<Scalars['String']['input']>;
  contentType?: InputMaybe<Scalars['String']['input']>;
  id: Scalars['ID']['input'];
};

export type OnboardingDetails = {
  __typename?: 'OnboardingDetails';
  comments?: Maybe<Scalars['String']['output']>;
  status: OnboardingStatus;
  updatedAt?: Maybe<Scalars['Time']['output']>;
};

export enum OnboardingStatus {
  Done = 'DONE',
  Late = 'LATE',
  NotApplicable = 'NOT_APPLICABLE',
  NotStarted = 'NOT_STARTED',
  OnTrack = 'ON_TRACK',
  Stuck = 'STUCK',
  Successful = 'SUCCESSFUL',
}

export type OnboardingStatusInput = {
  comments?: InputMaybe<Scalars['String']['input']>;
  organizationId: Scalars['ID']['input'];
  status: OnboardingStatus;
};

export type Opportunity = MetadataInterface & {
  __typename?: 'Opportunity';
  amount: Scalars['Float']['output'];
  /** Deprecated, use metadata */
  appSource?: Maybe<Scalars['String']['output']>;
  comments: Scalars['String']['output'];
  /** Deprecated, use metadata */
  createdAt?: Maybe<Scalars['Time']['output']>;
  createdBy?: Maybe<User>;
  currency?: Maybe<Currency>;
  estimatedClosedAt?: Maybe<Scalars['Time']['output']>;
  externalLinks: Array<ExternalSystem>;
  externalStage: Scalars['String']['output'];
  externalType: Scalars['String']['output'];
  generalNotes: Scalars['String']['output'];
  /** Deprecated, use metadata */
  id: Scalars['ID']['output'];
  internalStage: InternalStage;
  internalType: InternalType;
  likelihoodRate: Scalars['Int64']['output'];
  maxAmount: Scalars['Float']['output'];
  metadata: Metadata;
  name: Scalars['String']['output'];
  nextSteps: Scalars['String']['output'];
  organization?: Maybe<Organization>;
  owner?: Maybe<User>;
  renewalAdjustedRate: Scalars['Int64']['output'];
  renewalApproved: Scalars['Boolean']['output'];
  renewalLikelihood: OpportunityRenewalLikelihood;
  renewalUpdatedByUserAt?: Maybe<Scalars['Time']['output']>;
  renewalUpdatedByUserId: Scalars['String']['output'];
  renewedAt?: Maybe<Scalars['Time']['output']>;
  /** Deprecated, use metadata */
  source?: Maybe<DataSource>;
  /** Deprecated, use metadata */
  sourceOfTruth?: Maybe<DataSource>;
  stageLastUpdated?: Maybe<Scalars['Time']['output']>;
  /** Deprecated, use metadata */
  updatedAt?: Maybe<Scalars['Time']['output']>;
};

export type OpportunityCreateInput = {
  comments?: InputMaybe<Scalars['String']['input']>;
  currency?: InputMaybe<Currency>;
  estimatedClosedDate?: InputMaybe<Scalars['Time']['input']>;
  externalStage?: InputMaybe<Scalars['String']['input']>;
  externalType?: InputMaybe<Scalars['String']['input']>;
  generalNotes?: InputMaybe<Scalars['String']['input']>;
  internalType?: InputMaybe<InternalType>;
  likelihoodRate?: InputMaybe<Scalars['Int64']['input']>;
  maxAmount?: InputMaybe<Scalars['Float']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  nextSteps?: InputMaybe<Scalars['String']['input']>;
  organizationId: Scalars['ID']['input'];
};

export type OpportunityPage = Pages & {
  __typename?: 'OpportunityPage';
  content: Array<Opportunity>;
  totalAvailable: Scalars['Int64']['output'];
  totalElements: Scalars['Int64']['output'];
  totalPages: Scalars['Int']['output'];
};

export enum OpportunityRenewalLikelihood {
  HighRenewal = 'HIGH_RENEWAL',
  LowRenewal = 'LOW_RENEWAL',
  MediumRenewal = 'MEDIUM_RENEWAL',
  ZeroRenewal = 'ZERO_RENEWAL',
}

export type OpportunityRenewalUpdateAllForOrganizationInput = {
  organizationId: Scalars['ID']['input'];
  renewalAdjustedRate?: InputMaybe<Scalars['Int64']['input']>;
  renewalLikelihood?: InputMaybe<OpportunityRenewalLikelihood>;
};

export type OpportunityRenewalUpdateInput = {
  amount?: InputMaybe<Scalars['Float']['input']>;
  appSource?: InputMaybe<Scalars['String']['input']>;
  comments?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  name?: InputMaybe<Scalars['String']['input']>;
  opportunityId: Scalars['ID']['input'];
  ownerUserId?: InputMaybe<Scalars['ID']['input']>;
  renewalAdjustedRate?: InputMaybe<Scalars['Int64']['input']>;
  renewalLikelihood?: InputMaybe<OpportunityRenewalLikelihood>;
};

export type OpportunitySaveInput = {
  amount?: InputMaybe<Scalars['Float']['input']>;
  currency?: InputMaybe<Currency>;
  estimatedClosedDate?: InputMaybe<Scalars['Time']['input']>;
  externalStage?: InputMaybe<Scalars['String']['input']>;
  externalType?: InputMaybe<Scalars['String']['input']>;
  internalStage?: InputMaybe<InternalStage>;
  internalType?: InputMaybe<InternalType>;
  likelihoodRate?: InputMaybe<Scalars['Int64']['input']>;
  maxAmount?: InputMaybe<Scalars['Float']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  nextSteps?: InputMaybe<Scalars['String']['input']>;
  opportunityId?: InputMaybe<Scalars['ID']['input']>;
  organizationId?: InputMaybe<Scalars['ID']['input']>;
  ownerId?: InputMaybe<Scalars['ID']['input']>;
};

export type OpportunityUpdateInput = {
  amount?: InputMaybe<Scalars['Float']['input']>;
  currency?: InputMaybe<Currency>;
  estimatedClosedDate?: InputMaybe<Scalars['Time']['input']>;
  externalStage?: InputMaybe<Scalars['String']['input']>;
  externalType?: InputMaybe<Scalars['String']['input']>;
  internalStage?: InputMaybe<InternalStage>;
  likelihoodRate?: InputMaybe<Scalars['Int64']['input']>;
  maxAmount?: InputMaybe<Scalars['Float']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  nextSteps?: InputMaybe<Scalars['String']['input']>;
  opportunityId: Scalars['ID']['input'];
};

export type OrgAccountDetails = {
  __typename?: 'OrgAccountDetails';
  churned?: Maybe<Scalars['Time']['output']>;
  ltv?: Maybe<Scalars['Float']['output']>;
  ltvCurrency?: Maybe<Currency>;
  onboarding?: Maybe<OnboardingDetails>;
  renewalSummary?: Maybe<RenewalSummary>;
};

export type Organization = MetadataInterface & {
  __typename?: 'Organization';
  accountDetails?: Maybe<OrgAccountDetails>;
  /**
   * Deprecated
   * @deprecated Use metadata.appSource
   */
  appSource: Scalars['String']['output'];
  contactCount: Scalars['Int64']['output'];
  contacts: ContactsPage;
  contracts?: Maybe<Array<Contract>>;
  /**
   * Deprecated
   * @deprecated Use metadata.created
   */
  createdAt: Scalars['Time']['output'];
  customFields: Array<CustomField>;
  /**
   * Deprecated
   * @deprecated Use referenceId
   */
  customId?: Maybe<Scalars['String']['output']>;
  customerOsId: Scalars['String']['output'];
  description?: Maybe<Scalars['String']['output']>;
  domains: Array<Scalars['String']['output']>;
  emails: Array<Email>;
  employeeGrowthRate?: Maybe<Scalars['String']['output']>;
  employees?: Maybe<Scalars['Int64']['output']>;
  enrichDetails: EnrichDetails;
  externalLinks: Array<ExternalSystem>;
  headquarters?: Maybe<Scalars['String']['output']>;
  hide: Scalars['Boolean']['output'];
  /** @deprecated Use logo */
  icon?: Maybe<Scalars['String']['output']>;
  iconUrl?: Maybe<Scalars['String']['output']>;
  icpFit: Scalars['Boolean']['output'];
  /**
   * Deprecated
   * @deprecated Use metadata.id
   */
  id: Scalars['ID']['output'];
  inboundCommsCount: Scalars['Int64']['output'];
  industry?: Maybe<Scalars['String']['output']>;
  industryGroup?: Maybe<Scalars['String']['output']>;
  /**
   * Deprecated, use relationship instead
   * @deprecated Use relationship
   */
  isCustomer?: Maybe<Scalars['Boolean']['output']>;
  /**
   * Deprecated
   * @deprecated Use public
   */
  isPublic?: Maybe<Scalars['Boolean']['output']>;
  issueSummaryByStatus: Array<IssueSummaryByStatus>;
  jobRoles: Array<JobRole>;
  lastFundingAmount?: Maybe<Scalars['String']['output']>;
  lastFundingRound?: Maybe<FundingRound>;
  /**
   * Deprecated
   * @deprecated Use lastTouchpoint
   */
  lastTouchPointAt?: Maybe<Scalars['Time']['output']>;
  /**
   * Deprecated
   * @deprecated Use lastTouchpoint
   */
  lastTouchPointTimelineEvent?: Maybe<TimelineEvent>;
  /** Deprecated */
  lastTouchPointTimelineEventId?: Maybe<Scalars['ID']['output']>;
  /**
   * Deprecated
   * @deprecated Use lastTouchpoint
   */
  lastTouchPointType?: Maybe<LastTouchpointType>;
  lastTouchpoint?: Maybe<LastTouchpoint>;
  leadSource?: Maybe<Scalars['String']['output']>;
  locations: Array<Location>;
  /** @deprecated Use logo */
  logo?: Maybe<Scalars['String']['output']>;
  logoUrl?: Maybe<Scalars['String']['output']>;
  market?: Maybe<Market>;
  metadata: Metadata;
  name: Scalars['String']['output'];
  /**
   * Deprecated
   * @deprecated Use notes
   */
  note?: Maybe<Scalars['String']['output']>;
  notes?: Maybe<Scalars['String']['output']>;
  opportunities?: Maybe<Array<Opportunity>>;
  outboundCommsCount: Scalars['Int64']['output'];
  owner?: Maybe<User>;
  parentCompanies: Array<LinkedOrganization>;
  phoneNumbers: Array<PhoneNumber>;
  public?: Maybe<Scalars['Boolean']['output']>;
  referenceId?: Maybe<Scalars['String']['output']>;
  relationship?: Maybe<OrganizationRelationship>;
  slackChannelId?: Maybe<Scalars['String']['output']>;
  socialMedia: Array<Social>;
  /**
   * Deprecated
   * @deprecated Use socialMedia
   */
  socials: Array<Social>;
  /**
   * Deprecated
   * @deprecated Use metadata.source
   */
  source: DataSource;
  /**
   * Deprecated
   * @deprecated Use metadata.sourceOfTruth
   */
  sourceOfTruth: DataSource;
  stage?: Maybe<OrganizationStage>;
  stageLastUpdated?: Maybe<Scalars['Time']['output']>;
  subIndustry?: Maybe<Scalars['String']['output']>;
  subsidiaries: Array<LinkedOrganization>;
  /**
   * Deprecated
   * @deprecated Use parentCompany
   */
  subsidiaryOf: Array<LinkedOrganization>;
  suggestedMergeTo: Array<SuggestedMergeOrganization>;
  tags?: Maybe<Array<Tag>>;
  targetAudience?: Maybe<Scalars['String']['output']>;
  timelineEvents: Array<TimelineEvent>;
  timelineEventsTotalCount: Scalars['Int64']['output'];
  /**
   * Deprecated
   * @deprecated Use metadata.lastUpdated
   */
  updatedAt: Scalars['Time']['output'];
  valueProposition?: Maybe<Scalars['String']['output']>;
  website?: Maybe<Scalars['String']['output']>;
  yearFounded?: Maybe<Scalars['Int64']['output']>;
};

export type OrganizationContactsArgs = {
  pagination?: InputMaybe<Pagination>;
  sort?: InputMaybe<Array<SortBy>>;
  where?: InputMaybe<Filter>;
};

export type OrganizationTimelineEventsArgs = {
  from?: InputMaybe<Scalars['Time']['input']>;
  size: Scalars['Int']['input'];
  timelineEventTypes?: InputMaybe<Array<TimelineEventType>>;
};

export type OrganizationTimelineEventsTotalCountArgs = {
  timelineEventTypes?: InputMaybe<Array<TimelineEventType>>;
};

export type OrganizationInput = {
  appSource?: InputMaybe<Scalars['String']['input']>;
  customFields?: InputMaybe<Array<CustomFieldInput>>;
  /**
   * The name of the organization.
   * **Required.**
   */
  customId?: InputMaybe<Scalars['String']['input']>;
  description?: InputMaybe<Scalars['String']['input']>;
  domains?: InputMaybe<Array<Scalars['String']['input']>>;
  employeeGrowthRate?: InputMaybe<Scalars['String']['input']>;
  employees?: InputMaybe<Scalars['Int64']['input']>;
  headquarters?: InputMaybe<Scalars['String']['input']>;
  icon?: InputMaybe<Scalars['String']['input']>;
  industry?: InputMaybe<Scalars['String']['input']>;
  industryGroup?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated, use relationship instead */
  isCustomer?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  isPublic?: InputMaybe<Scalars['Boolean']['input']>;
  leadSource?: InputMaybe<Scalars['String']['input']>;
  logo?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  logoUrl?: InputMaybe<Scalars['String']['input']>;
  market?: InputMaybe<Market>;
  name?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  note?: InputMaybe<Scalars['String']['input']>;
  notes?: InputMaybe<Scalars['String']['input']>;
  public?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  referenceId?: InputMaybe<Scalars['String']['input']>;
  relationship?: InputMaybe<OrganizationRelationship>;
  slackChannelId?: InputMaybe<Scalars['String']['input']>;
  stage?: InputMaybe<OrganizationStage>;
  subIndustry?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  templateId?: InputMaybe<Scalars['ID']['input']>;
  website?: InputMaybe<Scalars['String']['input']>;
  yearFounded?: InputMaybe<Scalars['Int64']['input']>;
};

export type OrganizationPage = Pages & {
  __typename?: 'OrganizationPage';
  content: Array<Organization>;
  totalAvailable: Scalars['Int64']['output'];
  totalElements: Scalars['Int64']['output'];
  totalPages: Scalars['Int']['output'];
};

export type OrganizationParticipant = {
  __typename?: 'OrganizationParticipant';
  organizationParticipant: Organization;
  type?: Maybe<Scalars['String']['output']>;
};

export enum OrganizationRelationship {
  Customer = 'CUSTOMER',
  FormerCustomer = 'FORMER_CUSTOMER',
  NotAFit = 'NOT_A_FIT',
  Prospect = 'PROSPECT',
}

export type OrganizationSaveInput = {
  description?: InputMaybe<Scalars['String']['input']>;
  domains?: InputMaybe<Array<Scalars['String']['input']>>;
  employeeGrowthRate?: InputMaybe<Scalars['String']['input']>;
  employees?: InputMaybe<Scalars['Int64']['input']>;
  headquarters?: InputMaybe<Scalars['String']['input']>;
  iconUrl?: InputMaybe<Scalars['String']['input']>;
  icpFit?: InputMaybe<Scalars['Boolean']['input']>;
  id?: InputMaybe<Scalars['ID']['input']>;
  industry?: InputMaybe<Scalars['String']['input']>;
  industryGroup?: InputMaybe<Scalars['String']['input']>;
  lastFundingAmount?: InputMaybe<Scalars['String']['input']>;
  lastFundingRound?: InputMaybe<FundingRound>;
  leadSource?: InputMaybe<Scalars['String']['input']>;
  logoUrl?: InputMaybe<Scalars['String']['input']>;
  market?: InputMaybe<Market>;
  name?: InputMaybe<Scalars['String']['input']>;
  notes?: InputMaybe<Scalars['String']['input']>;
  ownerId?: InputMaybe<Scalars['String']['input']>;
  public?: InputMaybe<Scalars['Boolean']['input']>;
  referenceId?: InputMaybe<Scalars['String']['input']>;
  relationship?: InputMaybe<OrganizationRelationship>;
  slackChannelId?: InputMaybe<Scalars['String']['input']>;
  stage?: InputMaybe<OrganizationStage>;
  subIndustry?: InputMaybe<Scalars['String']['input']>;
  targetAudience?: InputMaybe<Scalars['String']['input']>;
  valueProposition?: InputMaybe<Scalars['String']['input']>;
  website?: InputMaybe<Scalars['String']['input']>;
  yearFounded?: InputMaybe<Scalars['Int64']['input']>;
};

export enum OrganizationStage {
  Engaged = 'ENGAGED',
  InitialValue = 'INITIAL_VALUE',
  Lead = 'LEAD',
  MaxValue = 'MAX_VALUE',
  Onboarding = 'ONBOARDING',
  PendingChurn = 'PENDING_CHURN',
  ReadyToBuy = 'READY_TO_BUY',
  RecurringValue = 'RECURRING_VALUE',
  Target = 'TARGET',
  Trial = 'TRIAL',
  Unqualified = 'UNQUALIFIED',
}

export type OrganizationTagInput = {
  organizationId: Scalars['ID']['input'];
  tag: TagIdOrNameInput;
};

export type OrganizationUpdateInput = {
  customId?: InputMaybe<Scalars['String']['input']>;
  description?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  domains?: InputMaybe<Array<Scalars['String']['input']>>;
  employeeGrowthRate?: InputMaybe<Scalars['String']['input']>;
  employees?: InputMaybe<Scalars['Int64']['input']>;
  headquarters?: InputMaybe<Scalars['String']['input']>;
  icon?: InputMaybe<Scalars['String']['input']>;
  icpFit?: InputMaybe<Scalars['Boolean']['input']>;
  id: Scalars['ID']['input'];
  industry?: InputMaybe<Scalars['String']['input']>;
  industryGroup?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated, use relationship instead */
  isCustomer?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated, use public instead */
  isPublic?: InputMaybe<Scalars['Boolean']['input']>;
  lastFundingAmount?: InputMaybe<Scalars['String']['input']>;
  lastFundingRound?: InputMaybe<FundingRound>;
  logo?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated, use logo instead */
  logoUrl?: InputMaybe<Scalars['String']['input']>;
  market?: InputMaybe<Market>;
  name?: InputMaybe<Scalars['String']['input']>;
  /** Deprecatedm, use notes instead */
  note?: InputMaybe<Scalars['String']['input']>;
  notes?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  patch?: InputMaybe<Scalars['Boolean']['input']>;
  public?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated, use customId instead */
  referenceId?: InputMaybe<Scalars['String']['input']>;
  relationship?: InputMaybe<OrganizationRelationship>;
  slackChannelId?: InputMaybe<Scalars['String']['input']>;
  stage?: InputMaybe<OrganizationStage>;
  subIndustry?: InputMaybe<Scalars['String']['input']>;
  targetAudience?: InputMaybe<Scalars['String']['input']>;
  valueProposition?: InputMaybe<Scalars['String']['input']>;
  website?: InputMaybe<Scalars['String']['input']>;
  yearFounded?: InputMaybe<Scalars['Int64']['input']>;
};

export type OrganizationWithJobRole = {
  __typename?: 'OrganizationWithJobRole';
  jobRole: JobRole;
  organization: Organization;
};

export type PageView = Node &
  SourceFields & {
    __typename?: 'PageView';
    appSource: Scalars['String']['output'];
    application: Scalars['String']['output'];
    endedAt: Scalars['Time']['output'];
    engagedTime: Scalars['Int64']['output'];
    id: Scalars['ID']['output'];
    orderInSession: Scalars['Int64']['output'];
    pageTitle: Scalars['String']['output'];
    pageUrl: Scalars['String']['output'];
    sessionId: Scalars['ID']['output'];
    source: DataSource;
    sourceOfTruth: DataSource;
    startedAt: Scalars['Time']['output'];
  };

/**
 * Describes the number of pages and total elements included in a query response.
 * **A `response` object.**
 */
export type Pages = {
  /**
   * The total number of elements included in the query response.
   * **Required.**
   */
  totalElements: Scalars['Int64']['output'];
  /**
   * The total number of pages included in the query response.
   * **Required.**
   */
  totalPages: Scalars['Int']['output'];
};

/** If provided as part of the request, results will be filtered down to the `page` and `limit` specified. */
export type Pagination = {
  /**
   * The maximum number of results in the response.
   * **Required.**
   */
  limit: Scalars['Int']['input'];
  /**
   * The results page to return in the response.
   * **Required.**
   */
  page: Scalars['Int']['input'];
};

/**
 * The honorific title of an individual.
 * **A `response` object.**
 */
export enum PersonTitle {
  /** For the holder of a doctoral degree. */
  Dr = 'DR',
  /** For girls, unmarried women, and married women who continue to use their maiden name. */
  Miss = 'MISS',
  /** For men, regardless of marital status. */
  Mr = 'MR',
  /** For married women. */
  Mrs = 'MRS',
  /** For women, regardless of marital status, or when marital status is unknown. */
  Ms = 'MS',
}

/**
 * Describes a phone number associated with a `Contact` in customerOS.
 * **A `return` object.**
 */
export type PhoneNumber = {
  __typename?: 'PhoneNumber';
  appSource?: Maybe<Scalars['String']['output']>;
  contacts: Array<Contact>;
  country?: Maybe<Country>;
  createdAt: Scalars['Time']['output'];
  /** The phone number in e164 format. */
  e164?: Maybe<Scalars['String']['output']>;
  /**
   * The unique ID associated with the phone number.
   * **Required**
   */
  id: Scalars['ID']['output'];
  /** Defines the type of phone number. */
  label?: Maybe<PhoneNumberLabel>;
  organizations: Array<Organization>;
  /**
   * Determines if the phone number is primary or not.
   * **Required**
   */
  primary: Scalars['Boolean']['output'];
  rawPhoneNumber?: Maybe<Scalars['String']['output']>;
  source: DataSource;
  updatedAt: Scalars['Time']['output'];
  users: Array<User>;
  validated?: Maybe<Scalars['Boolean']['output']>;
};

/**
 * Describes a phone number associated with a `Contact` in customerOS.
 * **A `create` object.**
 */
export type PhoneNumberInput = {
  countryCodeA2?: InputMaybe<Scalars['String']['input']>;
  /** Defines the type of phone number. */
  label?: InputMaybe<PhoneNumberLabel>;
  /**
   * The phone number in e164 format.
   * **Required**
   */
  phoneNumber: Scalars['String']['input'];
  /**
   * Determines if the phone number is primary or not.
   * **Required**
   */
  primary?: InputMaybe<Scalars['Boolean']['input']>;
};

/**
 * Defines the type of phone number.
 * **A `response` object. **
 */
export enum PhoneNumberLabel {
  Home = 'HOME',
  Main = 'MAIN',
  Mobile = 'MOBILE',
  Other = 'OTHER',
  Work = 'WORK',
}

export type PhoneNumberParticipant = {
  __typename?: 'PhoneNumberParticipant';
  phoneNumberParticipant: PhoneNumber;
  type?: Maybe<Scalars['String']['output']>;
};

/**
 * Describes a phone number associated with a `Contact` in customerOS.
 * **An `update` object.**
 */
export type PhoneNumberRelationUpdateInput = {
  /** Deprecated */
  countryCodeA2?: InputMaybe<Scalars['String']['input']>;
  /**
   * The unique ID associated with the phone number.
   * **Required**
   */
  id: Scalars['ID']['input'];
  /** Defines the type of phone number. */
  label?: InputMaybe<PhoneNumberLabel>;
  /** Deprecated */
  phoneNumber?: InputMaybe<Scalars['String']['input']>;
  /**
   * Determines if the phone number is primary or not.
   * **Required**
   */
  primary?: InputMaybe<Scalars['Boolean']['input']>;
};

/**
 * Describes a phone number associated with a `Contact` in customerOS.
 * **An `update` object.**
 */
export type PhoneNumberUpdateInput = {
  countryCodeA2?: InputMaybe<Scalars['String']['input']>;
  id: Scalars['ID']['input'];
  phoneNumber: Scalars['String']['input'];
};

export type Query = {
  __typename?: 'Query';
  attachment: Attachment;
  bankAccounts: Array<BankAccount>;
  billableInfo: TenantBillableInfo;
  /** Fetch a single contact from customerOS by contact ID. */
  contact?: Maybe<Contact>;
  contact_ByEmail: Contact;
  contact_ByPhone: Contact;
  /**
   * Fetch paginated list of contacts
   * Possible values for sort:
   * - PREFIX
   * - FIRST_NAME
   * - LAST_NAME
   * - NAME
   * - DESCRIPTION
   * - CREATED_AT
   */
  contacts: ContactsPage;
  contract: Contract;
  contracts: ContractPage;
  customFieldTemplate_List: Array<CustomFieldTemplate>;
  /** sort.By available options: ORGANIZATION, IS_CUSTOMER, DOMAIN, LOCATION, OWNER, LAST_TOUCHPOINT, RENEWAL_LIKELIHOOD, FORECAST_ARR, RENEWAL_DATE, ONBOARDING_STATUS */
  dashboardView_Organizations?: Maybe<OrganizationPage>;
  dashboardView_Renewals?: Maybe<RenewalsPage>;
  dashboard_ARRBreakdown?: Maybe<DashboardArrBreakdown>;
  dashboard_CustomerMap?: Maybe<Array<DashboardCustomerMap>>;
  dashboard_GrossRevenueRetention?: Maybe<DashboardGrossRevenueRetention>;
  dashboard_MRRPerCustomer?: Maybe<DashboardMrrPerCustomer>;
  dashboard_NewCustomers?: Maybe<DashboardNewCustomers>;
  dashboard_OnboardingCompletion?: Maybe<DashboardOnboardingCompletion>;
  dashboard_RetentionRate?: Maybe<DashboardRetentionRate>;
  dashboard_RevenueAtRisk?: Maybe<DashboardRevenueAtRisk>;
  dashboard_TimeToOnboard?: Maybe<DashboardTimeToOnboard>;
  email: Email;
  externalMeetings: MeetingsPage;
  externalSystemInstances: Array<ExternalSystemInstance>;
  flow: Flow;
  flowParticipant: FlowContact;
  flow_emailVariables: Array<EmailVariableEntity>;
  flows: Array<Flow>;
  gcli_Search: Array<GCliItem>;
  global_Cache: GlobalCache;
  interactionEvent: InteractionEvent;
  invoice: Invoice;
  invoice_ByNumber: Invoice;
  invoices: InvoicesPage;
  issue: Issue;
  logEntry: LogEntry;
  meeting: Meeting;
  opportunities_LinkedToOrganizations: OpportunityPage;
  opportunity?: Maybe<Opportunity>;
  organization?: Maybe<Organization>;
  organization_ByCustomId?: Maybe<Organization>;
  organization_ByCustomerOsId?: Maybe<Organization>;
  organization_CheckWebsite: WebsiteDetails;
  organization_DistinctOwners: Array<User>;
  organizations: OrganizationPage;
  organizations_HiddenAfter: Array<Scalars['String']['output']>;
  phoneNumber: PhoneNumber;
  reminder: Reminder;
  remindersForOrganization: Array<Reminder>;
  serviceLineItem: ServiceLineItem;
  slack_Channels: SlackChannelPage;
  tableViewDefs: Array<TableViewDef>;
  tags: Array<Tag>;
  tenant: Scalars['String']['output'];
  tenantBillingProfile: TenantBillingProfile;
  tenantBillingProfiles: Array<TenantBillingProfile>;
  tenantSettings: TenantSettings;
  timelineEvents: Array<TimelineEvent>;
  user: User;
  user_ByEmail: User;
  users: UserPage;
  workflow_ByType: Workflow;
  workflows: Array<Workflow>;
};

export type QueryAttachmentArgs = {
  id: Scalars['ID']['input'];
};

export type QueryContactArgs = {
  id: Scalars['ID']['input'];
};

export type QueryContact_ByEmailArgs = {
  email: Scalars['String']['input'];
};

export type QueryContact_ByPhoneArgs = {
  e164: Scalars['String']['input'];
};

export type QueryContactsArgs = {
  pagination?: InputMaybe<Pagination>;
  sort?: InputMaybe<Array<SortBy>>;
  where?: InputMaybe<Filter>;
};

export type QueryContractArgs = {
  id: Scalars['ID']['input'];
};

export type QueryContractsArgs = {
  pagination?: InputMaybe<Pagination>;
};

export type QueryDashboardView_OrganizationsArgs = {
  pagination: Pagination;
  sort?: InputMaybe<SortBy>;
  where?: InputMaybe<Filter>;
};

export type QueryDashboardView_RenewalsArgs = {
  pagination: Pagination;
  sort?: InputMaybe<SortBy>;
  where?: InputMaybe<Filter>;
};

export type QueryDashboard_ArrBreakdownArgs = {
  period?: InputMaybe<DashboardPeriodInput>;
};

export type QueryDashboard_GrossRevenueRetentionArgs = {
  period?: InputMaybe<DashboardPeriodInput>;
};

export type QueryDashboard_MrrPerCustomerArgs = {
  period?: InputMaybe<DashboardPeriodInput>;
};

export type QueryDashboard_NewCustomersArgs = {
  period?: InputMaybe<DashboardPeriodInput>;
};

export type QueryDashboard_OnboardingCompletionArgs = {
  period?: InputMaybe<DashboardPeriodInput>;
};

export type QueryDashboard_RetentionRateArgs = {
  period?: InputMaybe<DashboardPeriodInput>;
};

export type QueryDashboard_RevenueAtRiskArgs = {
  period?: InputMaybe<DashboardPeriodInput>;
};

export type QueryDashboard_TimeToOnboardArgs = {
  period?: InputMaybe<DashboardPeriodInput>;
};

export type QueryEmailArgs = {
  id: Scalars['ID']['input'];
};

export type QueryExternalMeetingsArgs = {
  externalId?: InputMaybe<Scalars['ID']['input']>;
  externalSystemId: Scalars['String']['input'];
  pagination?: InputMaybe<Pagination>;
  sort?: InputMaybe<Array<SortBy>>;
  where?: InputMaybe<Filter>;
};

export type QueryFlowArgs = {
  id: Scalars['ID']['input'];
};

export type QueryFlowParticipantArgs = {
  id: Scalars['ID']['input'];
};

export type QueryGcli_SearchArgs = {
  keyword: Scalars['String']['input'];
  limit?: InputMaybe<Scalars['Int']['input']>;
};

export type QueryInteractionEventArgs = {
  id: Scalars['ID']['input'];
};

export type QueryInvoiceArgs = {
  id: Scalars['ID']['input'];
};

export type QueryInvoice_ByNumberArgs = {
  number: Scalars['String']['input'];
};

export type QueryInvoicesArgs = {
  organizationId?: InputMaybe<Scalars['ID']['input']>;
  pagination?: InputMaybe<Pagination>;
  sort?: InputMaybe<Array<SortBy>>;
  where?: InputMaybe<Filter>;
};

export type QueryIssueArgs = {
  id: Scalars['ID']['input'];
};

export type QueryLogEntryArgs = {
  id: Scalars['ID']['input'];
};

export type QueryMeetingArgs = {
  id: Scalars['ID']['input'];
};

export type QueryOpportunities_LinkedToOrganizationsArgs = {
  pagination?: InputMaybe<Pagination>;
};

export type QueryOpportunityArgs = {
  id: Scalars['ID']['input'];
};

export type QueryOrganizationArgs = {
  id: Scalars['ID']['input'];
};

export type QueryOrganization_ByCustomIdArgs = {
  customId: Scalars['String']['input'];
};

export type QueryOrganization_ByCustomerOsIdArgs = {
  customerOsId: Scalars['String']['input'];
};

export type QueryOrganization_CheckWebsiteArgs = {
  website: Scalars['String']['input'];
};

export type QueryOrganizationsArgs = {
  pagination?: InputMaybe<Pagination>;
  sort?: InputMaybe<Array<SortBy>>;
  where?: InputMaybe<Filter>;
};

export type QueryOrganizations_HiddenAfterArgs = {
  date: Scalars['Time']['input'];
};

export type QueryPhoneNumberArgs = {
  id: Scalars['ID']['input'];
};

export type QueryReminderArgs = {
  id: Scalars['ID']['input'];
};

export type QueryRemindersForOrganizationArgs = {
  dismissed?: InputMaybe<Scalars['Boolean']['input']>;
  organizationId: Scalars['ID']['input'];
};

export type QueryServiceLineItemArgs = {
  id: Scalars['ID']['input'];
};

export type QuerySlack_ChannelsArgs = {
  pagination?: InputMaybe<Pagination>;
};

export type QueryTenantBillingProfileArgs = {
  id: Scalars['ID']['input'];
};

export type QueryTimelineEventsArgs = {
  ids: Array<Scalars['ID']['input']>;
};

export type QueryUserArgs = {
  id: Scalars['ID']['input'];
};

export type QueryUser_ByEmailArgs = {
  email: Scalars['String']['input'];
};

export type QueryUsersArgs = {
  pagination?: InputMaybe<Pagination>;
  sort?: InputMaybe<Array<SortBy>>;
  where?: InputMaybe<Filter>;
};

export type QueryWorkflow_ByTypeArgs = {
  workflowType: WorkflowType;
};

export type Reminder = MetadataInterface & {
  __typename?: 'Reminder';
  content?: Maybe<Scalars['String']['output']>;
  dismissed?: Maybe<Scalars['Boolean']['output']>;
  dueDate?: Maybe<Scalars['Time']['output']>;
  metadata: Metadata;
  owner?: Maybe<User>;
};

export type ReminderInput = {
  content: Scalars['String']['input'];
  dueDate: Scalars['Time']['input'];
  organizationId: Scalars['ID']['input'];
  userId: Scalars['ID']['input'];
};

export type ReminderUpdateInput = {
  content?: InputMaybe<Scalars['String']['input']>;
  dismissed?: InputMaybe<Scalars['Boolean']['input']>;
  dueDate?: InputMaybe<Scalars['Time']['input']>;
  id: Scalars['ID']['input'];
};

export type RemoveTagInput = {
  entityId: Scalars['ID']['input'];
  entityType: EntityType;
  tagId: Scalars['ID']['input'];
};

export type RenewalRecord = {
  __typename?: 'RenewalRecord';
  contract: Contract;
  opportunity?: Maybe<Opportunity>;
  organization: Organization;
};

export type RenewalSummary = {
  __typename?: 'RenewalSummary';
  arrForecast?: Maybe<Scalars['Float']['output']>;
  maxArrForecast?: Maybe<Scalars['Float']['output']>;
  nextRenewalDate?: Maybe<Scalars['Time']['output']>;
  renewalLikelihood?: Maybe<OpportunityRenewalLikelihood>;
};

export type RenewalsPage = Pages & {
  __typename?: 'RenewalsPage';
  content: Array<RenewalRecord>;
  totalAvailable: Scalars['Int64']['output'];
  totalElements: Scalars['Int64']['output'];
  totalPages: Scalars['Int']['output'];
};

/**
 * Describes the success or failure of the GraphQL call.
 * **A `return` object**
 */
export type Result = {
  __typename?: 'Result';
  /**
   * The result of the GraphQL call.
   * **Required.**
   */
  result: Scalars['Boolean']['output'];
};

export enum Role {
  Admin = 'ADMIN',
  Owner = 'OWNER',
  PlatformOwner = 'PLATFORM_OWNER',
  User = 'USER',
}

export type ServiceLineItem = MetadataInterface & {
  __typename?: 'ServiceLineItem';
  billingCycle: BilledType;
  closed: Scalars['Boolean']['output'];
  comments: Scalars['String']['output'];
  createdBy?: Maybe<User>;
  description: Scalars['String']['output'];
  externalLinks: Array<ExternalSystem>;
  metadata: Metadata;
  parentId: Scalars['ID']['output'];
  paused: Scalars['Boolean']['output'];
  price: Scalars['Float']['output'];
  quantity: Scalars['Int64']['output'];
  serviceEnded?: Maybe<Scalars['Time']['output']>;
  serviceStarted: Scalars['Time']['output'];
  tax: Tax;
};

export type ServiceLineItemBulkUpdateInput = {
  contractId: Scalars['ID']['input'];
  invoiceNote?: InputMaybe<Scalars['String']['input']>;
  serviceLineItems: Array<InputMaybe<ServiceLineItemBulkUpdateItem>>;
};

export type ServiceLineItemBulkUpdateItem = {
  billed?: InputMaybe<BilledType>;
  closeVersion?: InputMaybe<Scalars['Boolean']['input']>;
  comments?: InputMaybe<Scalars['String']['input']>;
  isRetroactiveCorrection?: InputMaybe<Scalars['Boolean']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  newVersion?: InputMaybe<Scalars['Boolean']['input']>;
  price?: InputMaybe<Scalars['Float']['input']>;
  quantity?: InputMaybe<Scalars['Int64']['input']>;
  serviceLineItemId?: InputMaybe<Scalars['ID']['input']>;
  serviceStarted?: InputMaybe<Scalars['Time']['input']>;
  vatRate?: InputMaybe<Scalars['Float']['input']>;
};

export type ServiceLineItemCloseInput = {
  endedAt?: InputMaybe<Scalars['Time']['input']>;
  id: Scalars['ID']['input'];
  serviceEnded?: InputMaybe<Scalars['Time']['input']>;
};

export type ServiceLineItemInput = {
  appSource?: InputMaybe<Scalars['String']['input']>;
  billingCycle?: InputMaybe<BilledType>;
  contractId: Scalars['ID']['input'];
  description?: InputMaybe<Scalars['String']['input']>;
  price?: InputMaybe<Scalars['Float']['input']>;
  quantity?: InputMaybe<Scalars['Int64']['input']>;
  serviceEnded?: InputMaybe<Scalars['Time']['input']>;
  serviceStarted?: InputMaybe<Scalars['Time']['input']>;
  tax?: InputMaybe<TaxInput>;
};

export type ServiceLineItemNewVersionInput = {
  appSource?: InputMaybe<Scalars['String']['input']>;
  comments?: InputMaybe<Scalars['String']['input']>;
  description?: InputMaybe<Scalars['String']['input']>;
  id?: InputMaybe<Scalars['ID']['input']>;
  price?: InputMaybe<Scalars['Float']['input']>;
  quantity?: InputMaybe<Scalars['Int64']['input']>;
  serviceStarted?: InputMaybe<Scalars['Time']['input']>;
  tax?: InputMaybe<TaxInput>;
};

export type ServiceLineItemUpdateInput = {
  appSource?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated: billing cycle is not updatable. */
  billingCycle?: InputMaybe<BilledType>;
  comments?: InputMaybe<Scalars['String']['input']>;
  description?: InputMaybe<Scalars['String']['input']>;
  id?: InputMaybe<Scalars['ID']['input']>;
  isRetroactiveCorrection?: InputMaybe<Scalars['Boolean']['input']>;
  price?: InputMaybe<Scalars['Float']['input']>;
  quantity?: InputMaybe<Scalars['Int64']['input']>;
  serviceEnded?: InputMaybe<Scalars['Time']['input']>;
  serviceStarted?: InputMaybe<Scalars['Time']['input']>;
  tax?: InputMaybe<TaxInput>;
};

export type SlackChannel = {
  __typename?: 'SlackChannel';
  channelId: Scalars['String']['output'];
  channelName: Scalars['String']['output'];
  metadata: Metadata;
  organization?: Maybe<Organization>;
};

export type SlackChannelPage = Pages & {
  __typename?: 'SlackChannelPage';
  content: Array<SlackChannel>;
  totalAvailable: Scalars['Int64']['output'];
  totalElements: Scalars['Int64']['output'];
  totalPages: Scalars['Int']['output'];
};

export type Social = Node &
  SourceFields & {
    __typename?: 'Social';
    alias: Scalars['String']['output'];
    appSource: Scalars['String']['output'];
    createdAt: Scalars['Time']['output'];
    externalId: Scalars['String']['output'];
    followersCount: Scalars['Int64']['output'];
    id: Scalars['ID']['output'];
    metadata: Metadata;
    source: DataSource;
    sourceOfTruth: DataSource;
    updatedAt: Scalars['Time']['output'];
    url: Scalars['String']['output'];
  };

export type SocialInput = {
  appSource?: InputMaybe<Scalars['String']['input']>;
  url: Scalars['String']['input'];
};

export type SocialUpdateInput = {
  id: Scalars['ID']['input'];
  url: Scalars['String']['input'];
};

export type SortBy = {
  by: Scalars['String']['input'];
  caseSensitive?: InputMaybe<Scalars['Boolean']['input']>;
  direction?: SortingDirection;
};

export enum SortingDirection {
  Asc = 'ASC',
  Desc = 'DESC',
}

export type SourceFields = {
  appSource: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  source: DataSource;
  sourceOfTruth: DataSource;
};

export type SourceFieldsInterface = {
  appSource: Scalars['String']['output'];
  source: DataSource;
  sourceOfTruth: DataSource;
};

export type State = {
  __typename?: 'State';
  code: Scalars['String']['output'];
  country: Country;
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
};

export type SuggestedMergeOrganization = {
  __typename?: 'SuggestedMergeOrganization';
  confidence?: Maybe<Scalars['Float']['output']>;
  organization: Organization;
  suggestedAt?: Maybe<Scalars['Time']['output']>;
  suggestedBy?: Maybe<Scalars['String']['output']>;
};

export enum TableIdType {
  Contacts = 'CONTACTS',
  ContactsForTargetOrganizations = 'CONTACTS_FOR_TARGET_ORGANIZATIONS',
  Contracts = 'CONTRACTS',
  Customers = 'CUSTOMERS',
  FlowActions = 'FLOW_ACTIONS',
  FlowContacts = 'FLOW_CONTACTS',
  Opportunities = 'OPPORTUNITIES',
  OpportunitiesRecords = 'OPPORTUNITIES_RECORDS',
  Organizations = 'ORGANIZATIONS',
  PastInvoices = 'PAST_INVOICES',
  Targets = 'TARGETS',
  UpcomingInvoices = 'UPCOMING_INVOICES',
}

export type TableViewDef = Node & {
  __typename?: 'TableViewDef';
  columns: Array<ColumnView>;
  createdAt: Scalars['Time']['output'];
  defaultFilters: Scalars['String']['output'];
  filters: Scalars['String']['output'];
  icon: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  isPreset: Scalars['Boolean']['output'];
  isShared: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  order: Scalars['Int']['output'];
  sorting: Scalars['String']['output'];
  tableId: TableIdType;
  tableType: TableViewType;
  updatedAt: Scalars['Time']['output'];
};

export type TableViewDefCreateInput = {
  columns: Array<ColumnViewInput>;
  filters: Scalars['String']['input'];
  icon: Scalars['String']['input'];
  isPreset: Scalars['Boolean']['input'];
  isShared: Scalars['Boolean']['input'];
  name: Scalars['String']['input'];
  order: Scalars['Int']['input'];
  sorting: Scalars['String']['input'];
  tableId: TableIdType;
  tableType: TableViewType;
};

export type TableViewDefUpdateInput = {
  columns: Array<ColumnViewInput>;
  filters: Scalars['String']['input'];
  icon: Scalars['String']['input'];
  id: Scalars['ID']['input'];
  name: Scalars['String']['input'];
  order: Scalars['Int']['input'];
  sorting: Scalars['String']['input'];
};

export enum TableViewType {
  Contacts = 'CONTACTS',
  Contracts = 'CONTRACTS',
  Flow = 'FLOW',
  Invoices = 'INVOICES',
  Opportunities = 'OPPORTUNITIES',
  Organizations = 'ORGANIZATIONS',
}

export type Tag = {
  __typename?: 'Tag';
  appSource: Scalars['String']['output'];
  createdAt: Scalars['Time']['output'];
  id: Scalars['ID']['output'];
  metadata: Metadata;
  name: Scalars['String']['output'];
  source: DataSource;
  updatedAt: Scalars['Time']['output'];
};

export type TagIdOrNameInput = {
  id?: InputMaybe<Scalars['ID']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
};

export type TagInput = {
  appSource?: InputMaybe<Scalars['String']['input']>;
  name: Scalars['String']['input'];
};

export type TagUpdateInput = {
  id: Scalars['ID']['input'];
  name: Scalars['String']['input'];
};

export type Tax = {
  __typename?: 'Tax';
  salesTax: Scalars['Boolean']['output'];
  taxRate: Scalars['Float']['output'];
  vat: Scalars['Boolean']['output'];
};

export type TaxInput = {
  taxRate: Scalars['Float']['input'];
};

export type TenantBillableInfo = {
  __typename?: 'TenantBillableInfo';
  greylistedContacts: Scalars['Int64']['output'];
  greylistedOrganizations: Scalars['Int64']['output'];
  whitelistedContacts: Scalars['Int64']['output'];
  whitelistedOrganizations: Scalars['Int64']['output'];
};

export type TenantBillingProfile = Node &
  SourceFields & {
    __typename?: 'TenantBillingProfile';
    addressLine1: Scalars['String']['output'];
    addressLine2: Scalars['String']['output'];
    addressLine3: Scalars['String']['output'];
    appSource: Scalars['String']['output'];
    canPayWithBankTransfer: Scalars['Boolean']['output'];
    /**
     * Deprecated
     * @deprecated Not used
     */
    canPayWithCard?: Maybe<Scalars['Boolean']['output']>;
    /**
     * Deprecated
     * @deprecated Not used
     */
    canPayWithDirectDebitACH?: Maybe<Scalars['Boolean']['output']>;
    /**
     * Deprecated
     * @deprecated Not used
     */
    canPayWithDirectDebitBacs?: Maybe<Scalars['Boolean']['output']>;
    /**
     * Deprecated
     * @deprecated Not used
     */
    canPayWithDirectDebitSEPA?: Maybe<Scalars['Boolean']['output']>;
    canPayWithPigeon: Scalars['Boolean']['output'];
    check: Scalars['Boolean']['output'];
    country: Scalars['String']['output'];
    createdAt: Scalars['Time']['output'];
    /**
     * Deprecated
     * @deprecated Not used
     */
    domesticPaymentsBankInfo?: Maybe<Scalars['String']['output']>;
    /**
     * Deprecated
     * @deprecated Use sendInvoicesFrom
     */
    email: Scalars['String']['output'];
    id: Scalars['ID']['output'];
    /**
     * Deprecated
     * @deprecated Not used
     */
    internationalPaymentsBankInfo?: Maybe<Scalars['String']['output']>;
    legalName: Scalars['String']['output'];
    locality: Scalars['String']['output'];
    phone: Scalars['String']['output'];
    region: Scalars['String']['output'];
    sendInvoicesBcc: Scalars['String']['output'];
    sendInvoicesFrom: Scalars['String']['output'];
    source: DataSource;
    sourceOfTruth: DataSource;
    updatedAt: Scalars['Time']['output'];
    vatNumber: Scalars['String']['output'];
    zip: Scalars['String']['output'];
  };

export type TenantBillingProfileInput = {
  addressLine1?: InputMaybe<Scalars['String']['input']>;
  addressLine2?: InputMaybe<Scalars['String']['input']>;
  addressLine3?: InputMaybe<Scalars['String']['input']>;
  canPayWithBankTransfer: Scalars['Boolean']['input'];
  /** Deprecated */
  canPayWithCard?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  canPayWithDirectDebitACH?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  canPayWithDirectDebitBacs?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  canPayWithDirectDebitSEPA?: InputMaybe<Scalars['Boolean']['input']>;
  canPayWithPigeon: Scalars['Boolean']['input'];
  check: Scalars['Boolean']['input'];
  country?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  domesticPaymentsBankInfo?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  email?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  internationalPaymentsBankInfo?: InputMaybe<Scalars['String']['input']>;
  legalName?: InputMaybe<Scalars['String']['input']>;
  locality?: InputMaybe<Scalars['String']['input']>;
  phone?: InputMaybe<Scalars['String']['input']>;
  region?: InputMaybe<Scalars['String']['input']>;
  sendInvoicesBcc?: InputMaybe<Scalars['String']['input']>;
  sendInvoicesFrom: Scalars['String']['input'];
  vatNumber: Scalars['String']['input'];
  zip?: InputMaybe<Scalars['String']['input']>;
};

export type TenantBillingProfileUpdateInput = {
  addressLine1?: InputMaybe<Scalars['String']['input']>;
  addressLine2?: InputMaybe<Scalars['String']['input']>;
  addressLine3?: InputMaybe<Scalars['String']['input']>;
  canPayWithBankTransfer?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  canPayWithCard?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  canPayWithDirectDebitACH?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  canPayWithDirectDebitBacs?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  canPayWithDirectDebitSEPA?: InputMaybe<Scalars['Boolean']['input']>;
  canPayWithPigeon?: InputMaybe<Scalars['Boolean']['input']>;
  check?: InputMaybe<Scalars['Boolean']['input']>;
  country?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  domesticPaymentsBankInfo?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  email?: InputMaybe<Scalars['String']['input']>;
  id: Scalars['ID']['input'];
  /** Deprecated */
  internationalPaymentsBankInfo?: InputMaybe<Scalars['String']['input']>;
  legalName?: InputMaybe<Scalars['String']['input']>;
  locality?: InputMaybe<Scalars['String']['input']>;
  patch?: InputMaybe<Scalars['Boolean']['input']>;
  phone?: InputMaybe<Scalars['String']['input']>;
  region?: InputMaybe<Scalars['String']['input']>;
  sendInvoicesBcc?: InputMaybe<Scalars['String']['input']>;
  sendInvoicesFrom?: InputMaybe<Scalars['String']['input']>;
  vatNumber?: InputMaybe<Scalars['String']['input']>;
  zip?: InputMaybe<Scalars['String']['input']>;
};

export type TenantInput = {
  appSource?: InputMaybe<Scalars['String']['input']>;
  name: Scalars['String']['input'];
};

export type TenantSettings = {
  __typename?: 'TenantSettings';
  baseCurrency?: Maybe<Currency>;
  billingEnabled: Scalars['Boolean']['output'];
  logoRepositoryFileId?: Maybe<Scalars['String']['output']>;
  /**
   * Deprecated
   * @deprecated Use logoRepositoryFileId
   */
  logoUrl: Scalars['String']['output'];
  opportunityStages: Array<TenantSettingsOpportunityStageConfiguration>;
  workspaceLogo?: Maybe<Scalars['String']['output']>;
  workspaceName?: Maybe<Scalars['String']['output']>;
};

export type TenantSettingsInput = {
  baseCurrency?: InputMaybe<Currency>;
  billingEnabled?: InputMaybe<Scalars['Boolean']['input']>;
  logoRepositoryFileId?: InputMaybe<Scalars['String']['input']>;
  logoUrl?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  patch?: InputMaybe<Scalars['Boolean']['input']>;
  workspaceLogo?: InputMaybe<Scalars['String']['input']>;
  workspaceName?: InputMaybe<Scalars['String']['input']>;
};

export type TenantSettingsOpportunityStageConfiguration = {
  __typename?: 'TenantSettingsOpportunityStageConfiguration';
  id: Scalars['ID']['output'];
  label: Scalars['String']['output'];
  likelihoodRate: Scalars['Int64']['output'];
  order: Scalars['Int']['output'];
  value: Scalars['String']['output'];
  visible: Scalars['Boolean']['output'];
};

export type TenantSettingsOpportunityStageConfigurationInput = {
  id: Scalars['ID']['input'];
  label?: InputMaybe<Scalars['String']['input']>;
  likelihoodRate?: InputMaybe<Scalars['Int64']['input']>;
  visible?: InputMaybe<Scalars['Boolean']['input']>;
};

export type TimeRange = {
  /**
   * The start time of the time range.
   * **Required.**
   */
  from: Scalars['Time']['input'];
  /**
   * The end time of the time range.
   * **Required.**
   */
  to: Scalars['Time']['input'];
};

export type TimelineEvent =
  | Action
  | InteractionEvent
  | InteractionSession
  | Issue
  | LogEntry
  | Meeting
  | Note
  | PageView;

export enum TimelineEventType {
  Action = 'ACTION',
  Analysis = 'ANALYSIS',
  InteractionEvent = 'INTERACTION_EVENT',
  InteractionSession = 'INTERACTION_SESSION',
  Issue = 'ISSUE',
  LogEntry = 'LOG_ENTRY',
  Meeting = 'MEETING',
  Note = 'NOTE',
  Order = 'ORDER',
  PageView = 'PAGE_VIEW',
}

/**
 * Describes the User of customerOS.  A user is the person who logs into the Openline platform.
 * **A `return` object**
 */
export type User = {
  __typename?: 'User';
  appSource: Scalars['String']['output'];
  bot: Scalars['Boolean']['output'];
  calendars: Array<Calendar>;
  /**
   * Timestamp of user creation.
   * **Required**
   */
  createdAt: Scalars['Time']['output'];
  /**
   * All email addresses associated with a user in customerOS.
   * **Required.  If no values it returns an empty array.**
   */
  emails?: Maybe<Array<Email>>;
  /**
   * The first name of the customerOS user.
   * **Required**
   */
  firstName: Scalars['String']['output'];
  hasLinkedInToken: Scalars['Boolean']['output'];
  /**
   * The unique ID associated with the customerOS user.
   * **Required**
   */
  id: Scalars['ID']['output'];
  internal: Scalars['Boolean']['output'];
  jobRoles: Array<JobRole>;
  /**
   * The last name of the customerOS user.
   * **Required**
   */
  lastName: Scalars['String']['output'];
  mailboxes: Array<Scalars['String']['output']>;
  name?: Maybe<Scalars['String']['output']>;
  phoneNumbers: Array<PhoneNumber>;
  profilePhotoUrl?: Maybe<Scalars['String']['output']>;
  roles: Array<Role>;
  source: DataSource;
  sourceOfTruth: DataSource;
  timezone?: Maybe<Scalars['String']['output']>;
  updatedAt: Scalars['Time']['output'];
};

/**
 * Describes the User of customerOS.  A user is the person who logs into the Openline platform.
 * **A `create` object.**
 */
export type UserInput = {
  /**
   * The name of the app performing the create.
   * **Optional**
   */
  appSource?: InputMaybe<Scalars['String']['input']>;
  /**
   * The email address of the customerOS user.
   * **Required**
   */
  email: EmailInput;
  /**
   * The first name of the customerOS user.
   * **Required**
   */
  firstName: Scalars['String']['input'];
  /**
   * The Job Roles of the user.
   * **Optional**
   */
  jobRoles?: InputMaybe<Array<JobRoleInput>>;
  /**
   * The last name of the customerOS user.
   * **Required**
   */
  lastName: Scalars['String']['input'];
  name?: InputMaybe<Scalars['String']['input']>;
  profilePhotoUrl?: InputMaybe<Scalars['String']['input']>;
  timezone?: InputMaybe<Scalars['String']['input']>;
};

/**
 * Specifies how many pages of `User` information has been returned in the query response.
 * **A `return` object.**
 */
export type UserPage = Pages & {
  __typename?: 'UserPage';
  /**
   * A `User` entity in customerOS.
   * **Required.  If no values it returns an empty array.**
   */
  content: Array<User>;
  /**
   * Total number of elements in the query response.
   * **Required.**
   */
  totalElements: Scalars['Int64']['output'];
  /**
   * Total number of pages in the query response.
   * **Required.**
   */
  totalPages: Scalars['Int']['output'];
};

export type UserParticipant = {
  __typename?: 'UserParticipant';
  type?: Maybe<Scalars['String']['output']>;
  userParticipant: User;
};

export type UserUpdateInput = {
  /**
   * The first name of the customerOS user.
   * **Required**
   */
  firstName: Scalars['String']['input'];
  id: Scalars['ID']['input'];
  /**
   * The last name of the customerOS user.
   * **Required**
   */
  lastName: Scalars['String']['input'];
  name?: InputMaybe<Scalars['String']['input']>;
  profilePhotoUrl?: InputMaybe<Scalars['String']['input']>;
  timezone?: InputMaybe<Scalars['String']['input']>;
};

export type WebsiteDetails = {
  __typename?: 'WebsiteDetails';
  domain: Scalars['String']['output'];
  primary: Scalars['Boolean']['output'];
  website: Scalars['String']['output'];
};

export type Workflow = Node & {
  __typename?: 'Workflow';
  actionParam1: Scalars['String']['output'];
  condition: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  live: Scalars['Boolean']['output'];
  name?: Maybe<Scalars['String']['output']>;
  type: WorkflowType;
};

export type WorkflowCreateInput = {
  actionParam1?: InputMaybe<Scalars['String']['input']>;
  condition?: InputMaybe<Scalars['String']['input']>;
  live?: InputMaybe<Scalars['Boolean']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  type: WorkflowType;
};

export enum WorkflowType {
  IdealContactPersona = 'IDEAL_CONTACT_PERSONA',
  IdealCustomerProfile = 'IDEAL_CUSTOMER_PROFILE',
}

export type WorkflowUpdateInput = {
  actionParam1?: InputMaybe<Scalars['String']['input']>;
  condition?: InputMaybe<Scalars['String']['input']>;
  id: Scalars['ID']['input'];
  live?: InputMaybe<Scalars['Boolean']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
};
