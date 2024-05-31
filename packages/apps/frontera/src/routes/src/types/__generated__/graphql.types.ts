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
  Any: { input: any; output: any };
  Time: { input: any; output: any };
  Int64: { input: any; output: any };
  ID: { input: string; output: string };
  Int: { input: number; output: number };
  Float: { input: number; output: number };
  String: { input: string; output: string };
  Boolean: { input: boolean; output: boolean };
};

export type Action = {
  source: DataSource;
  __typename?: 'Action';
  actionType: ActionType;
  createdBy?: Maybe<User>;
  id: Scalars['ID']['output'];
  createdAt: Scalars['Time']['output'];
  appSource: Scalars['String']['output'];
  content?: Maybe<Scalars['String']['output']>;
  metadata?: Maybe<Scalars['String']['output']>;
};

export type ActionItem = {
  source: DataSource;
  __typename?: 'ActionItem';
  id: Scalars['ID']['output'];
  content: Scalars['String']['output'];
  createdAt: Scalars['Time']['output'];
  appSource: Scalars['String']['output'];
};

export enum ActionType {
  ContractRenewed = 'CONTRACT_RENEWED',
  ContractStatusUpdated = 'CONTRACT_STATUS_UPDATED',
  Created = 'CREATED',
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

export type Analysis = Node & {
  source: DataSource;
  __typename?: 'Analysis';
  sourceOfTruth: DataSource;
  id: Scalars['ID']['output'];
  describes: Array<DescriptionNode>;
  createdAt: Scalars['Time']['output'];
  appSource: Scalars['String']['output'];
  content?: Maybe<Scalars['String']['output']>;
  contentType?: Maybe<Scalars['String']['output']>;
  analysisType?: Maybe<Scalars['String']['output']>;
};

export type AnalysisDescriptionInput = {
  meetingId?: InputMaybe<Scalars['ID']['input']>;
  interactionEventId?: InputMaybe<Scalars['ID']['input']>;
  interactionSessionId?: InputMaybe<Scalars['ID']['input']>;
};

export type AnalysisInput = {
  appSource: Scalars['String']['input'];
  describes: Array<AnalysisDescriptionInput>;
  content?: InputMaybe<Scalars['String']['input']>;
  contentType?: InputMaybe<Scalars['String']['input']>;
  analysisType?: InputMaybe<Scalars['String']['input']>;
};

export type Attachment = Node & {
  source: DataSource;
  __typename?: 'Attachment';
  sourceOfTruth: DataSource;
  id: Scalars['ID']['output'];
  size: Scalars['Int64']['output'];
  cdnUrl: Scalars['String']['output'];
  createdAt: Scalars['Time']['output'];
  basePath: Scalars['String']['output'];
  fileName: Scalars['String']['output'];
  mimeType: Scalars['String']['output'];
  appSource: Scalars['String']['output'];
};

export type AttachmentInput = {
  size: Scalars['Int64']['input'];
  cdnUrl: Scalars['String']['input'];
  basePath: Scalars['String']['input'];
  fileName: Scalars['String']['input'];
  mimeType: Scalars['String']['input'];
  appSource: Scalars['String']['input'];
  id?: InputMaybe<Scalars['ID']['input']>;
  createdAt?: InputMaybe<Scalars['Time']['input']>;
};

export type BankAccount = MetadataInterface & {
  metadata: Metadata;
  __typename?: 'BankAccount';
  currency?: Maybe<Currency>;
  bic?: Maybe<Scalars['String']['output']>;
  iban?: Maybe<Scalars['String']['output']>;
  bankName?: Maybe<Scalars['String']['output']>;
  sortCode?: Maybe<Scalars['String']['output']>;
  allowInternational: Scalars['Boolean']['output'];
  bankTransferEnabled: Scalars['Boolean']['output'];
  otherDetails?: Maybe<Scalars['String']['output']>;
  accountNumber?: Maybe<Scalars['String']['output']>;
  routingNumber?: Maybe<Scalars['String']['output']>;
};

export type BankAccountCreateInput = {
  currency?: InputMaybe<Currency>;
  bic?: InputMaybe<Scalars['String']['input']>;
  iban?: InputMaybe<Scalars['String']['input']>;
  bankName?: InputMaybe<Scalars['String']['input']>;
  sortCode?: InputMaybe<Scalars['String']['input']>;
  otherDetails?: InputMaybe<Scalars['String']['input']>;
  accountNumber?: InputMaybe<Scalars['String']['input']>;
  routingNumber?: InputMaybe<Scalars['String']['input']>;
  allowInternational?: InputMaybe<Scalars['Boolean']['input']>;
  bankTransferEnabled?: InputMaybe<Scalars['Boolean']['input']>;
};

export type BankAccountUpdateInput = {
  id: Scalars['ID']['input'];
  currency?: InputMaybe<Currency>;
  bic?: InputMaybe<Scalars['String']['input']>;
  iban?: InputMaybe<Scalars['String']['input']>;
  bankName?: InputMaybe<Scalars['String']['input']>;
  sortCode?: InputMaybe<Scalars['String']['input']>;
  otherDetails?: InputMaybe<Scalars['String']['input']>;
  accountNumber?: InputMaybe<Scalars['String']['input']>;
  routingNumber?: InputMaybe<Scalars['String']['input']>;
  allowInternational?: InputMaybe<Scalars['Boolean']['input']>;
  bankTransferEnabled?: InputMaybe<Scalars['Boolean']['input']>;
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
  /** @deprecated Use billingCycleInMonths instead. */
  billingCycle?: Maybe<ContractBillingCycle>;
  check?: Maybe<Scalars['Boolean']['output']>;
  dueDays?: Maybe<Scalars['Int64']['output']>;
  region?: Maybe<Scalars['String']['output']>;
  country?: Maybe<Scalars['String']['output']>;
  locality?: Maybe<Scalars['String']['output']>;
  payOnline?: Maybe<Scalars['Boolean']['output']>;
  postalCode?: Maybe<Scalars['String']['output']>;
  invoiceNote?: Maybe<Scalars['String']['output']>;
  nextInvoicing?: Maybe<Scalars['Time']['output']>;
  addressLine1?: Maybe<Scalars['String']['output']>;
  addressLine2?: Maybe<Scalars['String']['output']>;
  billingEmail?: Maybe<Scalars['String']['output']>;
  invoicingStarted?: Maybe<Scalars['Time']['output']>;
  canPayWithCard?: Maybe<Scalars['Boolean']['output']>;
  payAutomatically?: Maybe<Scalars['Boolean']['output']>;
  billingCycleInMonths?: Maybe<Scalars['Int64']['output']>;
  billingEmailCC?: Maybe<Array<Scalars['String']['output']>>;
  organizationLegalName?: Maybe<Scalars['String']['output']>;
  billingEmailBCC?: Maybe<Array<Scalars['String']['output']>>;
  canPayWithDirectDebit?: Maybe<Scalars['Boolean']['output']>;
  canPayWithBankTransfer?: Maybe<Scalars['Boolean']['output']>;
};

export type BillingDetailsInput = {
  /** Deprecated, use billingCycleInMonths instead. */
  billingCycle?: InputMaybe<ContractBillingCycle>;
  check?: InputMaybe<Scalars['Boolean']['input']>;
  dueDays?: InputMaybe<Scalars['Int64']['input']>;
  region?: InputMaybe<Scalars['String']['input']>;
  country?: InputMaybe<Scalars['String']['input']>;
  locality?: InputMaybe<Scalars['String']['input']>;
  payOnline?: InputMaybe<Scalars['Boolean']['input']>;
  postalCode?: InputMaybe<Scalars['String']['input']>;
  invoiceNote?: InputMaybe<Scalars['String']['input']>;
  addressLine1?: InputMaybe<Scalars['String']['input']>;
  addressLine2?: InputMaybe<Scalars['String']['input']>;
  billingEmail?: InputMaybe<Scalars['String']['input']>;
  invoicingStarted?: InputMaybe<Scalars['Time']['input']>;
  canPayWithCard?: InputMaybe<Scalars['Boolean']['input']>;
  payAutomatically?: InputMaybe<Scalars['Boolean']['input']>;
  billingCycleInMonths?: InputMaybe<Scalars['Int64']['input']>;
  billingEmailCC?: InputMaybe<Array<Scalars['String']['input']>>;
  organizationLegalName?: InputMaybe<Scalars['String']['input']>;
  billingEmailBCC?: InputMaybe<Array<Scalars['String']['input']>>;
  canPayWithDirectDebit?: InputMaybe<Scalars['Boolean']['input']>;
  canPayWithBankTransfer?: InputMaybe<Scalars['Boolean']['input']>;
};

export type BillingProfile = Node &
  SourceFields & {
    source: DataSource;
    sourceOfTruth: DataSource;
    id: Scalars['ID']['output'];
    __typename?: 'BillingProfile';
    taxId: Scalars['String']['output'];
    createdAt: Scalars['Time']['output'];
    updatedAt: Scalars['Time']['output'];
    appSource: Scalars['String']['output'];
    legalName: Scalars['String']['output'];
  };

export type BillingProfileInput = {
  organizationId: Scalars['ID']['input'];
  taxId?: InputMaybe<Scalars['String']['input']>;
  createdAt?: InputMaybe<Scalars['Time']['input']>;
  legalName?: InputMaybe<Scalars['String']['input']>;
};

export type BillingProfileLinkEmailInput = {
  emailId: Scalars['ID']['input'];
  organizationId: Scalars['ID']['input'];
  billingProfileId: Scalars['ID']['input'];
  primary?: InputMaybe<Scalars['Boolean']['input']>;
};

export type BillingProfileLinkLocationInput = {
  locationId: Scalars['ID']['input'];
  organizationId: Scalars['ID']['input'];
  billingProfileId: Scalars['ID']['input'];
};

export type BillingProfileUpdateInput = {
  organizationId: Scalars['ID']['input'];
  billingProfileId: Scalars['ID']['input'];
  taxId?: InputMaybe<Scalars['String']['input']>;
  updatedAt?: InputMaybe<Scalars['Time']['input']>;
  legalName?: InputMaybe<Scalars['String']['input']>;
};

export enum CalculationType {
  RevenueShare = 'REVENUE_SHARE',
}

/**
 * Describes the relationship a Contact has with a Organization.
 * **A `return` object**
 */
export type Calendar = {
  source: DataSource;
  calType: CalendarType;
  __typename?: 'Calendar';
  sourceOfTruth: DataSource;
  id: Scalars['ID']['output'];
  createdAt: Scalars['Time']['output'];
  updatedAt: Scalars['Time']['output'];
  primary: Scalars['Boolean']['output'];
  appSource: Scalars['String']['output'];
  link?: Maybe<Scalars['String']['output']>;
};

export enum CalendarType {
  Calcom = 'CALCOM',
  Google = 'GOOGLE',
}

export enum ChargePeriod {
  Annually = 'ANNUALLY',
  Monthly = 'MONTHLY',
  Quarterly = 'QUARTERLY',
}

export type ColumnView = {
  __typename?: 'ColumnView';
  columnType: ColumnViewType;
  width: Scalars['Int']['output'];
  visible: Scalars['Boolean']['output'];
};

export type ColumnViewInput = {
  columnType: ColumnViewType;
  width: Scalars['Int']['input'];
  visible: Scalars['Boolean']['input'];
};

export enum ColumnViewType {
  InvoicesAmount = 'INVOICES_AMOUNT',
  InvoicesBillingCycle = 'INVOICES_BILLING_CYCLE',
  InvoicesContract = 'INVOICES_CONTRACT',
  InvoicesDueDate = 'INVOICES_DUE_DATE',
  InvoicesInvoiceNumber = 'INVOICES_INVOICE_NUMBER',
  InvoicesInvoicePreview = 'INVOICES_INVOICE_PREVIEW',
  InvoicesInvoiceStatus = 'INVOICES_INVOICE_STATUS',
  InvoicesIssueDate = 'INVOICES_ISSUE_DATE',
  InvoicesIssueDatePast = 'INVOICES_ISSUE_DATE_PAST',
  InvoicesPaymentStatus = 'INVOICES_PAYMENT_STATUS',
  OrganizationsAvatar = 'ORGANIZATIONS_AVATAR',
  OrganizationsContactCount = 'ORGANIZATIONS_CONTACT_COUNT',
  OrganizationsCreatedDate = 'ORGANIZATIONS_CREATED_DATE',
  OrganizationsEmployeeCount = 'ORGANIZATIONS_EMPLOYEE_COUNT',
  OrganizationsForecastArr = 'ORGANIZATIONS_FORECAST_ARR',
  OrganizationsLastTouchpoint = 'ORGANIZATIONS_LAST_TOUCHPOINT',
  OrganizationsLastTouchpointDate = 'ORGANIZATIONS_LAST_TOUCHPOINT_DATE',
  OrganizationsLeadSource = 'ORGANIZATIONS_LEAD_SOURCE',
  OrganizationsName = 'ORGANIZATIONS_NAME',
  OrganizationsOnboardingStatus = 'ORGANIZATIONS_ONBOARDING_STATUS',
  OrganizationsOwner = 'ORGANIZATIONS_OWNER',
  OrganizationsRelationship = 'ORGANIZATIONS_RELATIONSHIP',
  OrganizationsRenewalDate = 'ORGANIZATIONS_RENEWAL_DATE',
  OrganizationsRenewalLikelihood = 'ORGANIZATIONS_RENEWAL_LIKELIHOOD',
  OrganizationsSocials = 'ORGANIZATIONS_SOCIALS',
  OrganizationsStage = 'ORGANIZATIONS_STAGE',
  OrganizationsWebsite = 'ORGANIZATIONS_WEBSITE',
  OrganizationsYearFounded = 'ORGANIZATIONS_YEAR_FOUNDED',
  RenewalsAvatar = 'RENEWALS_AVATAR',
  RenewalsForecastArr = 'RENEWALS_FORECAST_ARR',
  RenewalsLastTouchpoint = 'RENEWALS_LAST_TOUCHPOINT',
  RenewalsName = 'RENEWALS_NAME',
  RenewalsOwner = 'RENEWALS_OWNER',
  RenewalsRenewalDate = 'RENEWALS_RENEWAL_DATE',
  RenewalsRenewalLikelihood = 'RENEWALS_RENEWAL_LIKELIHOOD',
}

export type Comment = {
  source: DataSource;
  __typename?: 'Comment';
  createdBy?: Maybe<User>;
  sourceOfTruth: DataSource;
  id: Scalars['ID']['output'];
  createdAt: Scalars['Time']['output'];
  externalLinks: Array<ExternalSystem>;
  updatedAt: Scalars['Time']['output'];
  appSource: Scalars['String']['output'];
  content?: Maybe<Scalars['String']['output']>;
  contentType?: Maybe<Scalars['String']['output']>;
};

export enum ComparisonOperator {
  Between = 'BETWEEN',
  Contains = 'CONTAINS',
  Eq = 'EQ',
  Gte = 'GTE',
  In = 'IN',
  IsEmpty = 'IS_EMPTY',
  IsNull = 'IS_NULL',
  Lte = 'LTE',
  StartsWith = 'STARTS_WITH',
}

export type Conditionals = {
  __typename?: 'Conditionals';
  minimumChargePeriod?: Maybe<ChargePeriod>;
  minimumChargeAmount: Scalars['Float']['output'];
};

/**
 * A contact represents an individual in customerOS.
 * **A `response` object.**
 */
export type Contact = ExtensibleEntity &
  Node & {
    /** Contact notes */
    notes: NotePage;
    source: DataSource;
    /** Contact owner (user) */
    owner?: Maybe<User>;
    /**
     * All email addresses associated with a contact in customerOS.
     * **Required.  If no values it returns an empty array.**
     */
    emails: Array<Email>;
    __typename?: 'Contact';
    socials: Array<Social>;
    /**
     * `organizationName` and `jobTitle` of the contact if it has been associated with an organization.
     * **Required.  If no values it returns an empty array.**
     */
    jobRoles: Array<JobRole>;
    notesByTime: Array<Note>;
    tags?: Maybe<Array<Tag>>;
    sourceOfTruth: DataSource;
    fieldSets: Array<FieldSet>;
    /**
     * All locations associated with a contact in customerOS.
     * **Required.  If no values it returns an empty array.**
     */
    locations: Array<Location>;
    /**
     * The unique ID associated with the contact in customerOS.
     * **Required**
     */
    id: Scalars['ID']['output'];
    organizations: OrganizationPage;
    /**
     * User defined metadata appended to the contact record in customerOS.
     * **Required.  If no values it returns an empty array.**
     */
    customFields: Array<CustomField>;
    /**
     * All phone numbers associated with a contact in customerOS.
     * **Required.  If no values it returns an empty array.**
     */
    phoneNumbers: Array<PhoneNumber>;
    /** Template of the contact in customerOS. */
    template?: Maybe<EntityTemplate>;
    /**
     * An ISO8601 timestamp recording when the contact was created in customerOS.
     * **Required**
     */
    createdAt: Scalars['Time']['output'];
    timelineEvents: Array<TimelineEvent>;
    updatedAt: Scalars['Time']['output'];
    /** The name of the contact in customerOS, alternative for firstName + lastName. */
    name?: Maybe<Scalars['String']['output']>;
    /**
     * Deprecated
     * @deprecated Use `tags` instead
     */
    label?: Maybe<Scalars['String']['output']>;
    /**
     * Deprecated
     * @deprecated Use `prefix` instead
     */
    title?: Maybe<Scalars['String']['output']>;
    prefix?: Maybe<Scalars['String']['output']>;
    /** The last name of the contact in customerOS. */
    lastName?: Maybe<Scalars['String']['output']>;
    timezone?: Maybe<Scalars['String']['output']>;
    appSource?: Maybe<Scalars['String']['output']>;
    /** The first name of the contact in customerOS. */
    firstName?: Maybe<Scalars['String']['output']>;
    description?: Maybe<Scalars['String']['output']>;
    profilePhotoUrl?: Maybe<Scalars['String']['output']>;
    timelineEventsTotalCount: Scalars['Int64']['output'];
  };

/**
 * A contact represents an individual in customerOS.
 * **A `response` object.**
 */
export type ContactNotesArgs = {
  pagination?: InputMaybe<Pagination>;
};

/**
 * A contact represents an individual in customerOS.
 * **A `response` object.**
 */
export type ContactNotesByTimeArgs = {
  pagination?: InputMaybe<TimeRange>;
};

/**
 * A contact represents an individual in customerOS.
 * **A `response` object.**
 */
export type ContactOrganizationsArgs = {
  where?: InputMaybe<Filter>;
  sort?: InputMaybe<Array<SortBy>>;
  pagination?: InputMaybe<Pagination>;
};

/**
 * A contact represents an individual in customerOS.
 * **A `response` object.**
 */
export type ContactTimelineEventsArgs = {
  size: Scalars['Int']['input'];
  from?: InputMaybe<Scalars['Time']['input']>;
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
  /** An email addresses associated with the contact. */
  email?: InputMaybe<EmailInput>;
  /** A phone number associated with the contact. */
  phoneNumber?: InputMaybe<PhoneNumberInput>;
  /** Deprecated */
  fieldSets?: InputMaybe<Array<FieldSetInput>>;
  /** Deprecated */
  ownerId?: InputMaybe<Scalars['ID']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  /** The prefix of the contact. */
  prefix?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  templateId?: InputMaybe<Scalars['ID']['input']>;
  /** An ISO8601 timestamp recording when the contact was created in customerOS. */
  createdAt?: InputMaybe<Scalars['Time']['input']>;
  /** The last name of the contact. */
  lastName?: InputMaybe<Scalars['String']['input']>;
  timezone?: InputMaybe<Scalars['String']['input']>;
  appSource?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  customFields?: InputMaybe<Array<CustomFieldInput>>;
  /** The first name of the contact. */
  firstName?: InputMaybe<Scalars['String']['input']>;
  description?: InputMaybe<Scalars['String']['input']>;
  profilePhotoUrl?: InputMaybe<Scalars['String']['input']>;
  externalReference?: InputMaybe<ExternalSystemReferenceInput>;
};

export type ContactOrganizationInput = {
  contactId: Scalars['ID']['input'];
  organizationId: Scalars['ID']['input'];
};

export type ContactParticipant = {
  contactParticipant: Contact;
  __typename?: 'ContactParticipant';
  type?: Maybe<Scalars['String']['output']>;
};

export type ContactTagInput = {
  tagId: Scalars['ID']['input'];
  contactId: Scalars['ID']['input'];
};

/**
 * Updates data fields associated with an existing customer record in customerOS.
 * **An `update` object.**
 */
export type ContactUpdateInput = {
  id: Scalars['ID']['input'];
  name?: InputMaybe<Scalars['String']['input']>;
  patch?: InputMaybe<Scalars['Boolean']['input']>;
  prefix?: InputMaybe<Scalars['String']['input']>;
  lastName?: InputMaybe<Scalars['String']['input']>;
  timezone?: InputMaybe<Scalars['String']['input']>;
  firstName?: InputMaybe<Scalars['String']['input']>;
  description?: InputMaybe<Scalars['String']['input']>;
  profilePhotoUrl?: InputMaybe<Scalars['String']['input']>;
};

/**
 * Specifies how many pages of contact information has been returned in the query response.
 * **A `response` object.**
 */
export type ContactsPage = Pages & {
  /**
   * A contact entity in customerOS.
   * **Required.  If no values it returns an empty array.**
   */
  content: Array<Contact>;
  __typename?: 'ContactsPage';
  /**
   * Total number of pages in the query response.
   * **Required.**
   */
  totalPages: Scalars['Int']['output'];
  /**
   * Total number of elements in the query response.
   * **Required.**
   */
  totalElements: Scalars['Int64']['output'];
};

export type Contract = MetadataInterface & {
  metadata: Metadata;
  /**
   * Deprecated, use metadata instead.
   * @deprecated Use metadata instead.
   */
  source: DataSource;
  owner?: Maybe<User>;
  /**
   * Deprecated, use contractStatus instead.
   * @deprecated Use contractStatus instead.
   */
  status: ContractStatus;
  __typename?: 'Contract';
  createdBy?: Maybe<User>;
  invoices: Array<Invoice>;
  /**
   * Deprecated, use metadata instead.
   * @deprecated Use metadata instead.
   */
  sourceOfTruth: DataSource;
  currency?: Maybe<Currency>;
  /**
   * Deprecated, use metadata instead.
   * @deprecated Use metadata instead.
   */
  id: Scalars['ID']['output'];
  contractStatus: ContractStatus;
  upcomingInvoices: Array<Invoice>;
  /**
   * Deprecated, use contractName instead.
   * @deprecated Use contractName instead.
   */
  name: Scalars['String']['output'];
  /**
   * Deprecated, use contractRenewalCycle instead.
   * @deprecated Use contractRenewalCycle instead.
   */
  renewalCycle: ContractRenewalCycle;
  /**
   * Deprecated, use metadata instead.
   * @deprecated Use metadata instead.
   */
  createdAt: Scalars['Time']['output'];
  externalLinks: Array<ExternalSystem>;
  /**
   * Deprecated, use metadata instead.
   * @deprecated Use metadata instead.
   */
  updatedAt: Scalars['Time']['output'];
  /**
   * Deprecated, use metadata instead.
   * @deprecated Use metadata instead.
   */
  appSource: Scalars['String']['output'];
  approved: Scalars['Boolean']['output'];
  attachments?: Maybe<Array<Attachment>>;
  billingDetails?: Maybe<BillingDetails>;
  autoRenew: Scalars['Boolean']['output'];
  /**
   * Deprecated, use billingDetails instead.
   * @deprecated Use billingDetails instead.
   */
  zip?: Maybe<Scalars['String']['output']>;
  contractName: Scalars['String']['output'];
  opportunities?: Maybe<Array<Opportunity>>;
  /**
   * Deprecated, use billingDetails instead.
   * @deprecated Use billingDetails instead.
   */
  billingCycle?: Maybe<ContractBillingCycle>;
  /**
   * Deprecated, use committedPeriodInMonths instead.
   * @deprecated Use committedPeriodInMonths instead.
   */
  contractRenewalCycle: ContractRenewalCycle;
  /**
   * Deprecated, use contractEnded instead.
   * @deprecated Use contractEnded instead.
   */
  endedAt?: Maybe<Scalars['Time']['output']>;
  /**
   * Deprecated, use contractSigned instead.
   * @deprecated Use contractSigned instead.
   */
  signedAt?: Maybe<Scalars['Time']['output']>;
  billingEnabled: Scalars['Boolean']['output'];
  /**
   * Deprecated, use billingDetails instead.
   * @deprecated Use billingDetails instead.
   */
  country?: Maybe<Scalars['String']['output']>;
  /**
   * Deprecated, use billingDetails instead.
   * @deprecated Use billingDetails instead.
   */
  locality?: Maybe<Scalars['String']['output']>;
  contractEnded?: Maybe<Scalars['Time']['output']>;
  contractUrl?: Maybe<Scalars['String']['output']>;
  /**
   * Deprecated, use billingDetails instead.
   * @deprecated Use billingDetails instead.
   */
  invoiceNote?: Maybe<Scalars['String']['output']>;
  /**
   * Deprecated, use contractLineItems instead.
   * @deprecated Use contractLineItems instead.
   */
  serviceLineItems?: Maybe<Array<ServiceLineItem>>;
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
  contractLineItems?: Maybe<Array<ServiceLineItem>>;
  contractSigned?: Maybe<Scalars['Time']['output']>;
  /**
   * Deprecated, use billingDetails instead.
   * @deprecated Use billingDetails instead.
   */
  invoiceEmail?: Maybe<Scalars['String']['output']>;
  serviceStarted?: Maybe<Scalars['Time']['output']>;
  /**
   * Deprecated, use committedPeriods instead.
   * @deprecated Use committedPeriods instead.
   */
  renewalPeriods?: Maybe<Scalars['Int64']['output']>;
  /**
   * Deprecated, use serviceStarted instead.
   * @deprecated Use serviceStarted instead.
   */
  serviceStartedAt?: Maybe<Scalars['Time']['output']>;
  /**
   * Deprecated, use committedPeriodInMonths instead.
   * @deprecated Use committedPeriodInMonths instead.
   */
  committedPeriods?: Maybe<Scalars['Int64']['output']>;
  /**
   * Deprecated, use billingDetails instead.
   * @deprecated Use billingDetails instead.
   */
  invoicingStartDate?: Maybe<Scalars['Time']['output']>;
  /**
   * Deprecated, use billingDetails instead.
   * @deprecated Use billingDetails instead.
   */
  organizationLegalName?: Maybe<Scalars['String']['output']>;
  committedPeriodInMonths?: Maybe<Scalars['Int64']['output']>;
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
  currency?: InputMaybe<Currency>;
  organizationId: Scalars['ID']['input'];
  /** Deprecated */
  name?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  billingCycle?: InputMaybe<ContractBillingCycle>;
  dueDays?: InputMaybe<Scalars['Int64']['input']>;
  /** Deprecated */
  renewalCycle?: InputMaybe<ContractRenewalCycle>;
  /** Deprecated */
  signedAt?: InputMaybe<Scalars['Time']['input']>;
  appSource?: InputMaybe<Scalars['String']['input']>;
  approved?: InputMaybe<Scalars['Boolean']['input']>;
  autoRenew?: InputMaybe<Scalars['Boolean']['input']>;
  contractUrl?: InputMaybe<Scalars['String']['input']>;
  contractName?: InputMaybe<Scalars['String']['input']>;
  contractSigned?: InputMaybe<Scalars['Time']['input']>;
  serviceStarted?: InputMaybe<Scalars['Time']['input']>;
  /** Deprecated */
  renewalPeriods?: InputMaybe<Scalars['Int64']['input']>;
  /** Deprecated, use committedPeriodInMonths instead. */
  contractRenewalCycle?: InputMaybe<ContractRenewalCycle>;
  /** Deprecated */
  serviceStartedAt?: InputMaybe<Scalars['Time']['input']>;
  billingEnabled?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated, use committedPeriodInMonths instead. */
  committedPeriods?: InputMaybe<Scalars['Int64']['input']>;
  /** Deprecated */
  invoicingStartDate?: InputMaybe<Scalars['Time']['input']>;
  /** Deprecated */
  externalReference?: InputMaybe<ExternalSystemReferenceInput>;
  committedPeriodInMonths?: InputMaybe<Scalars['Int64']['input']>;
};

export type ContractPage = Pages & {
  content: Array<Contract>;
  __typename?: 'ContractPage';
  totalPages: Scalars['Int']['output'];
  totalElements: Scalars['Int64']['output'];
  totalAvailable: Scalars['Int64']['output'];
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
  currency?: InputMaybe<Currency>;
  contractId: Scalars['ID']['input'];
  /** Deprecated */
  zip?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  name?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  endedAt?: InputMaybe<Scalars['Time']['input']>;
  /** Deprecated */
  billingCycle?: InputMaybe<ContractBillingCycle>;
  patch?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  renewalCycle?: InputMaybe<ContractRenewalCycle>;
  /** Deprecated */
  signedAt?: InputMaybe<Scalars['Time']['input']>;
  billingDetails?: InputMaybe<BillingDetailsInput>;
  /** Deprecated */
  country?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  locality?: InputMaybe<Scalars['String']['input']>;
  appSource?: InputMaybe<Scalars['String']['input']>;
  approved?: InputMaybe<Scalars['Boolean']['input']>;
  autoRenew?: InputMaybe<Scalars['Boolean']['input']>;
  contractEnded?: InputMaybe<Scalars['Time']['input']>;
  contractUrl?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  invoiceNote?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  addressLine1?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  addressLine2?: InputMaybe<Scalars['String']['input']>;
  contractName?: InputMaybe<Scalars['String']['input']>;
  contractSigned?: InputMaybe<Scalars['Time']['input']>;
  /** Deprecated */
  invoiceEmail?: InputMaybe<Scalars['String']['input']>;
  serviceStarted?: InputMaybe<Scalars['Time']['input']>;
  /** Deprecated */
  renewalPeriods?: InputMaybe<Scalars['Int64']['input']>;
  /** Deprecated, use committedPeriodInMonths instead. */
  contractRenewalCycle?: InputMaybe<ContractRenewalCycle>;
  /** Deprecated */
  serviceStartedAt?: InputMaybe<Scalars['Time']['input']>;
  billingEnabled?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  canPayWithCard?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated, use committedPeriodInMonths instead. */
  committedPeriods?: InputMaybe<Scalars['Int64']['input']>;
  /** Deprecated */
  invoicingStartDate?: InputMaybe<Scalars['Time']['input']>;
  /** Deprecated */
  organizationLegalName?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  canPayWithDirectDebit?: InputMaybe<Scalars['Boolean']['input']>;
  committedPeriodInMonths?: InputMaybe<Scalars['Int64']['input']>;
  /** Deprecated */
  canPayWithBankTransfer?: InputMaybe<Scalars['Boolean']['input']>;
};

export type Country = {
  __typename?: 'Country';
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  codeA2: Scalars['String']['output'];
  codeA3: Scalars['String']['output'];
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

/**
 * Describes a custom, user-defined field associated with a `Contact`.
 * **A `return` object.**
 */
export type CustomField = Node & {
  /** The source of the custom field value */
  source: DataSource;
  __typename?: 'CustomField';
  /**
   * The unique ID associated with the custom field.
   * **Required**
   */
  id: Scalars['ID']['output'];
  /**
   * Datatype of the custom field.
   * **Required**
   */
  datatype: CustomFieldDataType;
  /**
   * The value of the custom field.
   * **Required**
   */
  value: Scalars['Any']['output'];
  /**
   * The name of the custom field.
   * **Required**
   */
  name: Scalars['String']['output'];
  createdAt: Scalars['Time']['output'];
  updatedAt: Scalars['Time']['output'];
  template?: Maybe<CustomFieldTemplate>;
};

export enum CustomFieldDataType {
  Bool = 'BOOL',
  Datetime = 'DATETIME',
  Decimal = 'DECIMAL',
  Integer = 'INTEGER',
  Text = 'TEXT',
}

export type CustomFieldEntityType = {
  entityType: EntityType;
  id: Scalars['ID']['input'];
};

/**
 * Describes a custom, user-defined field associated with a `Contact` of type String.
 * **A `create` object.**
 */
export type CustomFieldInput = {
  /**
   * The value of the custom field.
   * **Required**
   */
  value: Scalars['Any']['input'];
  /** Deprecated */
  id?: InputMaybe<Scalars['ID']['input']>;
  /** Datatype of the custom field. */
  datatype?: InputMaybe<CustomFieldDataType>;
  /** The name of the custom field. */
  name?: InputMaybe<Scalars['String']['input']>;
  templateId?: InputMaybe<Scalars['ID']['input']>;
};

export type CustomFieldTemplate = Node & {
  id: Scalars['ID']['output'];
  type: CustomFieldTemplateType;
  order: Scalars['Int']['output'];
  name: Scalars['String']['output'];
  __typename?: 'CustomFieldTemplate';
  createdAt: Scalars['Time']['output'];
  updatedAt: Scalars['Time']['output'];
  max?: Maybe<Scalars['Int']['output']>;
  min?: Maybe<Scalars['Int']['output']>;
  mandatory: Scalars['Boolean']['output'];
  length?: Maybe<Scalars['Int']['output']>;
};

export type CustomFieldTemplateInput = {
  type: CustomFieldTemplateType;
  order: Scalars['Int']['input'];
  name: Scalars['String']['input'];
  max?: InputMaybe<Scalars['Int']['input']>;
  min?: InputMaybe<Scalars['Int']['input']>;
  length?: InputMaybe<Scalars['Int']['input']>;
  mandatory?: InputMaybe<Scalars['Boolean']['input']>;
};

export enum CustomFieldTemplateType {
  Link = 'LINK',
  Text = 'TEXT',
}

/**
 * Describes a custom, user-defined field associated with a `Contact`.
 * **An `update` object.**
 */
export type CustomFieldUpdateInput = {
  /**
   * The unique ID associated with the custom field.
   * **Required**
   */
  id: Scalars['ID']['input'];
  /**
   * Datatype of the custom field.
   * **Required**
   */
  datatype: CustomFieldDataType;
  /**
   * The value of the custom field.
   * **Required**
   */
  value: Scalars['Any']['input'];
  /**
   * The name of the custom field.
   * **Required**
   */
  name: Scalars['String']['input'];
};

export type CustomerContact = {
  email: CustomerEmail;
  id: Scalars['ID']['output'];
  __typename?: 'CustomerContact';
};

export type CustomerContactInput = {
  /** An email addresses associted with the contact. */
  email?: InputMaybe<EmailInput>;
  name?: InputMaybe<Scalars['String']['input']>;
  /** The prefix of the contact. */
  prefix?: InputMaybe<Scalars['String']['input']>;
  /** An ISO8601 timestamp recording when the contact was created in customerOS. */
  createdAt?: InputMaybe<Scalars['Time']['input']>;
  /** The last name of the contact. */
  lastName?: InputMaybe<Scalars['String']['input']>;
  timezone?: InputMaybe<Scalars['String']['input']>;
  appSource?: InputMaybe<Scalars['String']['input']>;
  /** The first name of the contact. */
  firstName?: InputMaybe<Scalars['String']['input']>;
  description?: InputMaybe<Scalars['String']['input']>;
};

export type CustomerEmail = {
  id: Scalars['ID']['output'];
  __typename?: 'CustomerEmail';
};

export type CustomerJobRole = {
  id: Scalars['ID']['output'];
  __typename?: 'CustomerJobRole';
};

export type CustomerUser = {
  jobRole: CustomerJobRole;
  __typename?: 'CustomerUser';
  id: Scalars['ID']['output'];
};

export type DashboardArrBreakdown = {
  __typename?: 'DashboardARRBreakdown';
  arrBreakdown: Scalars['Float']['output'];
  increasePercentage: Scalars['String']['output'];
  perMonth: Array<Maybe<DashboardArrBreakdownPerMonth>>;
};

export type DashboardArrBreakdownPerMonth = {
  year: Scalars['Int']['output'];
  month: Scalars['Int']['output'];
  churned: Scalars['Float']['output'];
  upsells: Scalars['Float']['output'];
  renewals: Scalars['Float']['output'];
  downgrades: Scalars['Float']['output'];
  cancellations: Scalars['Float']['output'];
  newlyContracted: Scalars['Float']['output'];
  __typename?: 'DashboardARRBreakdownPerMonth';
};

export type DashboardCustomerMap = {
  organization: Organization;
  arr: Scalars['Float']['output'];
  state: DashboardCustomerMapState;
  __typename?: 'DashboardCustomerMap';
  organizationId: Scalars['ID']['output'];
  contractSignedDate: Scalars['Time']['output'];
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
  /**
   * Deprecated
   * @deprecated Use increasePercentageValue instead
   */
  increasePercentage: Scalars['String']['output'];
  grossRevenueRetention: Scalars['Float']['output'];
  increasePercentageValue: Scalars['Float']['output'];
  perMonth: Array<Maybe<DashboardGrossRevenueRetentionPerMonth>>;
};

export type DashboardGrossRevenueRetentionPerMonth = {
  year: Scalars['Int']['output'];
  month: Scalars['Int']['output'];
  percentage: Scalars['Float']['output'];
  __typename?: 'DashboardGrossRevenueRetentionPerMonth';
};

export type DashboardMrrPerCustomer = {
  __typename?: 'DashboardMRRPerCustomer';
  mrrPerCustomer: Scalars['Float']['output'];
  increasePercentage: Scalars['String']['output'];
  perMonth: Array<Maybe<DashboardMrrPerCustomerPerMonth>>;
};

export type DashboardMrrPerCustomerPerMonth = {
  year: Scalars['Int']['output'];
  month: Scalars['Int']['output'];
  value: Scalars['Float']['output'];
  __typename?: 'DashboardMRRPerCustomerPerMonth';
};

export type DashboardNewCustomers = {
  __typename?: 'DashboardNewCustomers';
  thisMonthCount: Scalars['Int']['output'];
  perMonth: Array<Maybe<DashboardNewCustomersPerMonth>>;
  thisMonthIncreasePercentage: Scalars['String']['output'];
};

export type DashboardNewCustomersPerMonth = {
  year: Scalars['Int']['output'];
  count: Scalars['Int']['output'];
  month: Scalars['Int']['output'];
  __typename?: 'DashboardNewCustomersPerMonth';
};

export type DashboardOnboardingCompletion = {
  __typename?: 'DashboardOnboardingCompletion';
  increasePercentage: Scalars['Float']['output'];
  completionPercentage: Scalars['Float']['output'];
  perMonth: Array<DashboardOnboardingCompletionPerMonth>;
};

export type DashboardOnboardingCompletionPerMonth = {
  year: Scalars['Int']['output'];
  month: Scalars['Int']['output'];
  value: Scalars['Float']['output'];
  __typename?: 'DashboardOnboardingCompletionPerMonth';
};

export type DashboardPeriodInput = {
  end: Scalars['Time']['input'];
  start: Scalars['Time']['input'];
};

export type DashboardRetentionRate = {
  __typename?: 'DashboardRetentionRate';
  retentionRate: Scalars['Float']['output'];
  /**
   * Deprecated
   * @deprecated Use increasePercentageValue instead
   */
  increasePercentage: Scalars['String']['output'];
  increasePercentageValue: Scalars['Float']['output'];
  perMonth: Array<Maybe<DashboardRetentionRatePerMonth>>;
};

export type DashboardRetentionRatePerMonth = {
  year: Scalars['Int']['output'];
  month: Scalars['Int']['output'];
  churnCount: Scalars['Int']['output'];
  renewCount: Scalars['Int']['output'];
  __typename?: 'DashboardRetentionRatePerMonth';
};

export type DashboardRevenueAtRisk = {
  atRisk: Scalars['Float']['output'];
  __typename?: 'DashboardRevenueAtRisk';
  highConfidence: Scalars['Float']['output'];
};

export type DashboardTimeToOnboard = {
  __typename?: 'DashboardTimeToOnboard';
  perMonth: Array<DashboardTimeToOnboardPerMonth>;
  timeToOnboard?: Maybe<Scalars['Float']['output']>;
  increasePercentage?: Maybe<Scalars['Float']['output']>;
};

export type DashboardTimeToOnboardPerMonth = {
  year: Scalars['Int']['output'];
  month: Scalars['Int']['output'];
  value: Scalars['Float']['output'];
  __typename?: 'DashboardTimeToOnboardPerMonth';
};

export enum DataSource {
  Close = 'CLOSE',
  Hubspot = 'HUBSPOT',
  Intercom = 'INTERCOM',
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
  ZendeskSupport = 'ZENDESK_SUPPORT',
}

export type DeleteResponse = {
  __typename?: 'DeleteResponse';
  accepted: Scalars['Boolean']['output'];
  completed: Scalars['Boolean']['output'];
};

export type DescriptionNode = InteractionEvent | InteractionSession | Meeting;

/**
 * Describes an email address associated with a `Contact` in customerOS.
 * **A `return` object.**
 */
export type Email = {
  source: DataSource;
  users: Array<User>;
  __typename?: 'Email';
  contacts: Array<Contact>;
  /** Describes the type of email address (WORK, PERSONAL, etc). */
  label?: Maybe<EmailLabel>;
  sourceOfTruth: DataSource;
  /**
   * The unique ID associated with the contact in customerOS.
   * **Required**
   */
  id: Scalars['ID']['output'];
  organizations: Array<Organization>;
  createdAt: Scalars['Time']['output'];
  updatedAt: Scalars['Time']['output'];
  /**
   * Identifies whether the email address is primary or not.
   * **Required.**
   */
  primary: Scalars['Boolean']['output'];
  appSource: Scalars['String']['output'];
  /** An email address assocaited with the contact in customerOS. */
  email?: Maybe<Scalars['String']['output']>;
  rawEmail?: Maybe<Scalars['String']['output']>;
  emailValidationDetails: EmailValidationDetails;
};

/**
 * Describes an email address associated with a `Contact` in customerOS.
 * **A `create` object.**
 */
export type EmailInput = {
  /** Describes the type of email address (WORK, PERSONAL, etc). */
  label?: InputMaybe<EmailLabel>;
  /**
   * An email address associated with the contact in customerOS.
   * **Required.**
   */
  email: Scalars['String']['input'];
  /**
   * Identifies whether the email address is primary or not.
   * **Required.**
   */
  primary?: InputMaybe<Scalars['Boolean']['input']>;
  appSource?: InputMaybe<Scalars['String']['input']>;
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
  emailParticipant: Email;
  __typename?: 'EmailParticipant';
  type?: Maybe<Scalars['String']['output']>;
};

/**
 * Describes an email address associated with a `Contact` in customerOS.
 * **An `update` object.**
 */
export type EmailUpdateInput = {
  /**
   * An email address assocaited with the contact in customerOS.
   * **Required.**
   */
  id: Scalars['ID']['input'];
  /** Describes the type of email address (WORK, PERSONAL, etc). */
  label?: InputMaybe<EmailLabel>;
  email?: InputMaybe<Scalars['String']['input']>;
  /**
   * Identifies whether the email address is primary or not.
   * **Required.**
   */
  primary?: InputMaybe<Scalars['Boolean']['input']>;
};

export type EmailValidationDetails = {
  __typename?: 'EmailValidationDetails';
  error?: Maybe<Scalars['String']['output']>;
  validated?: Maybe<Scalars['Boolean']['output']>;
  isCatchAll?: Maybe<Scalars['Boolean']['output']>;
  isDisabled?: Maybe<Scalars['Boolean']['output']>;
  isReachable?: Maybe<Scalars['String']['output']>;
  acceptsMail?: Maybe<Scalars['Boolean']['output']>;
  hasFullInbox?: Maybe<Scalars['Boolean']['output']>;
  isDeliverable?: Maybe<Scalars['Boolean']['output']>;
  isValidSyntax?: Maybe<Scalars['Boolean']['output']>;
  canConnectSmtp?: Maybe<Scalars['Boolean']['output']>;
};

export type EntityTemplate = Node & {
  id: Scalars['ID']['output'];
  __typename?: 'EntityTemplate';
  name: Scalars['String']['output'];
  version: Scalars['Int']['output'];
  createdAt: Scalars['Time']['output'];
  updatedAt: Scalars['Time']['output'];
  extends?: Maybe<EntityTemplateExtension>;
  fieldSetTemplates: Array<FieldSetTemplate>;
  customFieldTemplates: Array<CustomFieldTemplate>;
};

export enum EntityTemplateExtension {
  Contact = 'CONTACT',
  Organization = 'ORGANIZATION',
}

export type EntityTemplateInput = {
  name: Scalars['String']['input'];
  extends?: InputMaybe<EntityTemplateExtension>;
  fieldSetTemplateInputs?: InputMaybe<Array<FieldSetTemplateInput>>;
  customFieldTemplateInputs?: InputMaybe<Array<CustomFieldTemplateInput>>;
};

export enum EntityType {
  Contact = 'Contact',
  Organization = 'Organization',
}

export type ExtensibleEntity = {
  id: Scalars['ID']['output'];
  template?: Maybe<EntityTemplate>;
};

export type ExternalSystem = {
  type: ExternalSystemType;
  __typename?: 'ExternalSystem';
  syncDate?: Maybe<Scalars['Time']['output']>;
  externalId?: Maybe<Scalars['String']['output']>;
  externalUrl?: Maybe<Scalars['String']['output']>;
  externalSource?: Maybe<Scalars['String']['output']>;
};

export type ExternalSystemInput = {
  name: Scalars['String']['input'];
};

export type ExternalSystemInstance = {
  type: ExternalSystemType;
  __typename?: 'ExternalSystemInstance';
  stripeDetails?: Maybe<ExternalSystemStripeDetails>;
};

export type ExternalSystemReferenceInput = {
  type: ExternalSystemType;
  externalId: Scalars['ID']['input'];
  syncDate?: InputMaybe<Scalars['Time']['input']>;
  externalUrl?: InputMaybe<Scalars['String']['input']>;
  externalSource?: InputMaybe<Scalars['String']['input']>;
};

export type ExternalSystemStripeDetails = {
  __typename?: 'ExternalSystemStripeDetails';
  paymentMethodTypes: Array<Scalars['String']['output']>;
};

export enum ExternalSystemType {
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
  ZendeskSupport = 'ZENDESK_SUPPORT',
}

export type FieldSet = {
  source: DataSource;
  __typename?: 'FieldSet';
  id: Scalars['ID']['output'];
  customFields: Array<CustomField>;
  name: Scalars['String']['output'];
  template?: Maybe<FieldSetTemplate>;
  createdAt: Scalars['Time']['output'];
  updatedAt: Scalars['Time']['output'];
};

export type FieldSetInput = {
  name: Scalars['String']['input'];
  id?: InputMaybe<Scalars['ID']['input']>;
  templateId?: InputMaybe<Scalars['ID']['input']>;
  customFields?: InputMaybe<Array<CustomFieldInput>>;
};

export type FieldSetTemplate = Node & {
  id: Scalars['ID']['output'];
  __typename?: 'FieldSetTemplate';
  order: Scalars['Int']['output'];
  name: Scalars['String']['output'];
  createdAt: Scalars['Time']['output'];
  updatedAt: Scalars['Time']['output'];
  customFieldTemplates: Array<CustomFieldTemplate>;
};

export type FieldSetTemplateInput = {
  order: Scalars['Int']['input'];
  name: Scalars['String']['input'];
  customFieldTemplateInputs?: InputMaybe<Array<CustomFieldTemplateInput>>;
};

export type FieldSetUpdateInput = {
  id: Scalars['ID']['input'];
  name: Scalars['String']['input'];
};

export type Filter = {
  NOT?: InputMaybe<Filter>;
  OR?: InputMaybe<Array<Filter>>;
  AND?: InputMaybe<Array<Filter>>;
  filter?: InputMaybe<FilterItem>;
};

export type FilterItem = {
  operation?: ComparisonOperator;
  value: Scalars['Any']['input'];
  property: Scalars['String']['input'];
  includeEmpty?: InputMaybe<Scalars['Boolean']['input']>;
  caseSensitive?: InputMaybe<Scalars['Boolean']['input']>;
};

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
  key: Scalars['String']['output'];
  value: Scalars['String']['output'];
  __typename?: 'GCliAttributeKeyValuePair';
  display?: Maybe<Scalars['String']['output']>;
};

export enum GCliCacheItemType {
  Contact = 'CONTACT',
  Organization = 'ORGANIZATION',
  State = 'STATE',
}

export type GCliItem = {
  __typename?: 'GCliItem';
  type: GCliSearchResultType;
  id: Scalars['ID']['output'];
  display: Scalars['String']['output'];
  data?: Maybe<Array<GCliAttributeKeyValuePair>>;
};

export enum GCliSearchResultType {
  Contact = 'CONTACT',
  Email = 'EMAIL',
  Organization = 'ORGANIZATION',
  OrganizationRelationship = 'ORGANIZATION_RELATIONSHIP',
  State = 'STATE',
}

export type GlobalCache = {
  user: User;
  __typename?: 'GlobalCache';
  gCliCache: Array<GCliItem>;
  isOwner: Scalars['Boolean']['output'];
  cdnLogoUrl: Scalars['String']['output'];
  contractsExist: Scalars['Boolean']['output'];
  isGoogleActive: Scalars['Boolean']['output'];
  maxARRForecastValue: Scalars['Float']['output'];
  minARRForecastValue: Scalars['Float']['output'];
  isGoogleTokenExpired: Scalars['Boolean']['output'];
};

export type InteractionEvent = Node & {
  source: DataSource;
  issue?: Maybe<Issue>;
  meeting?: Maybe<Meeting>;
  sourceOfTruth: DataSource;
  summary?: Maybe<Analysis>;
  id: Scalars['ID']['output'];
  includes: Array<Attachment>;
  actions?: Maybe<Array<Action>>;
  __typename?: 'InteractionEvent';
  repliesTo?: Maybe<InteractionEvent>;
  createdAt: Scalars['Time']['output'];
  externalLinks: Array<ExternalSystem>;
  actionItems?: Maybe<Array<ActionItem>>;
  appSource: Scalars['String']['output'];
  sentBy: Array<InteractionEventParticipant>;
  sentTo: Array<InteractionEventParticipant>;
  channel?: Maybe<Scalars['String']['output']>;
  content?: Maybe<Scalars['String']['output']>;
  eventType?: Maybe<Scalars['String']['output']>;
  interactionSession?: Maybe<InteractionSession>;
  channelData?: Maybe<Scalars['String']['output']>;
  contentType?: Maybe<Scalars['String']['output']>;
  eventIdentifier?: Maybe<Scalars['String']['output']>;
  customerOSInternalIdentifier?: Maybe<Scalars['String']['output']>;
};

export type InteractionEventInput = {
  appSource: Scalars['String']['input'];
  meetingId?: InputMaybe<Scalars['ID']['input']>;
  repliesTo?: InputMaybe<Scalars['ID']['input']>;
  sentBy: Array<InteractionEventParticipantInput>;
  sentTo: Array<InteractionEventParticipantInput>;
  channel?: InputMaybe<Scalars['String']['input']>;
  content?: InputMaybe<Scalars['String']['input']>;
  createdAt?: InputMaybe<Scalars['Time']['input']>;
  eventType?: InputMaybe<Scalars['String']['input']>;
  externalId?: InputMaybe<Scalars['String']['input']>;
  channelData?: InputMaybe<Scalars['String']['input']>;
  contentType?: InputMaybe<Scalars['String']['input']>;
  interactionSession?: InputMaybe<Scalars['ID']['input']>;
  eventIdentifier?: InputMaybe<Scalars['String']['input']>;
  externalSystemId?: InputMaybe<Scalars['String']['input']>;
  customerOSInternalIdentifier?: InputMaybe<Scalars['String']['input']>;
};

export type InteractionEventParticipant =
  | ContactParticipant
  | EmailParticipant
  | JobRoleParticipant
  | OrganizationParticipant
  | PhoneNumberParticipant
  | UserParticipant;

export type InteractionEventParticipantInput = {
  userID?: InputMaybe<Scalars['ID']['input']>;
  type?: InputMaybe<Scalars['String']['input']>;
  contactID?: InputMaybe<Scalars['ID']['input']>;
  email?: InputMaybe<Scalars['String']['input']>;
  phoneNumber?: InputMaybe<Scalars['String']['input']>;
};

export type InteractionSession = Node & {
  source: DataSource;
  sourceOfTruth: DataSource;
  id: Scalars['ID']['output'];
  includes: Array<Attachment>;
  describedBy: Array<Analysis>;
  events: Array<InteractionEvent>;
  __typename?: 'InteractionSession';
  name: Scalars['String']['output'];
  status: Scalars['String']['output'];
  createdAt: Scalars['Time']['output'];
  /**
   * Deprecated
   * @deprecated Use createdAt instead
   */
  startedAt: Scalars['Time']['output'];
  updatedAt: Scalars['Time']['output'];
  appSource: Scalars['String']['output'];
  type?: Maybe<Scalars['String']['output']>;
  /**
   * Deprecated
   * @deprecated Use updatedAt instead
   */
  endedAt?: Maybe<Scalars['Time']['output']>;
  channel?: Maybe<Scalars['String']['output']>;
  attendedBy: Array<InteractionSessionParticipant>;
  channelData?: Maybe<Scalars['String']['output']>;
  sessionIdentifier?: Maybe<Scalars['String']['output']>;
};

export type InteractionSessionInput = {
  name: Scalars['String']['input'];
  status: Scalars['String']['input'];
  appSource: Scalars['String']['input'];
  type?: InputMaybe<Scalars['String']['input']>;
  channel?: InputMaybe<Scalars['String']['input']>;
  channelData?: InputMaybe<Scalars['String']['input']>;
  sessionIdentifier?: InputMaybe<Scalars['String']['input']>;
  attendedBy?: InputMaybe<Array<InteractionSessionParticipantInput>>;
};

export type InteractionSessionParticipant =
  | ContactParticipant
  | EmailParticipant
  | PhoneNumberParticipant
  | UserParticipant;

export type InteractionSessionParticipantInput = {
  userID?: InputMaybe<Scalars['ID']['input']>;
  type?: InputMaybe<Scalars['String']['input']>;
  contactID?: InputMaybe<Scalars['ID']['input']>;
  email?: InputMaybe<Scalars['String']['input']>;
  phoneNumber?: InputMaybe<Scalars['String']['input']>;
};

export enum InternalStage {
  ClosedLost = 'CLOSED_LOST',
  ClosedWon = 'CLOSED_WON',
  Evaluating = 'EVALUATING',
  Open = 'OPEN',
}

export enum InternalType {
  CrossSell = 'CROSS_SELL',
  Nbo = 'NBO',
  Renewal = 'RENEWAL',
  Upsell = 'UPSELL',
}

export type Invoice = MetadataInterface & {
  contract: Contract;
  metadata: Metadata;
  __typename?: 'Invoice';
  customer: InvoiceCustomer;
  provider: InvoiceProvider;
  organization: Organization;
  status?: Maybe<InvoiceStatus>;
  due: Scalars['Time']['output'];
  issued: Scalars['Time']['output'];
  paid: Scalars['Boolean']['output'];
  taxDue: Scalars['Float']['output'];
  dryRun: Scalars['Boolean']['output'];
  invoiceLineItems: Array<InvoiceLine>;
  subtotal: Scalars['Float']['output'];
  amountDue: Scalars['Float']['output'];
  currency: Scalars['String']['output'];
  preview: Scalars['Boolean']['output'];
  amountPaid: Scalars['Float']['output'];
  offCycle: Scalars['Boolean']['output'];
  postpaid: Scalars['Boolean']['output'];
  invoiceUrl: Scalars['String']['output'];
  note?: Maybe<Scalars['String']['output']>;
  invoiceNumber: Scalars['String']['output'];
  amountRemaining: Scalars['Float']['output'];
  invoicePeriodEnd: Scalars['Time']['output'];
  invoicePeriodStart: Scalars['Time']['output'];
  repositoryFileId: Scalars['String']['output'];
  paymentLink?: Maybe<Scalars['String']['output']>;
  /**
   * Deprecated
   * @deprecated not used
   */
  domesticPaymentsBankInfo?: Maybe<Scalars['String']['output']>;
  /**
   * Deprecated
   * @deprecated not used
   */
  internationalPaymentsBankInfo?: Maybe<Scalars['String']['output']>;
};

export type InvoiceCustomer = {
  __typename?: 'InvoiceCustomer';
  name?: Maybe<Scalars['String']['output']>;
  email?: Maybe<Scalars['String']['output']>;
  addressZip?: Maybe<Scalars['String']['output']>;
  addressLine1?: Maybe<Scalars['String']['output']>;
  addressLine2?: Maybe<Scalars['String']['output']>;
  addressRegion?: Maybe<Scalars['String']['output']>;
  addressCountry?: Maybe<Scalars['String']['output']>;
  addressLocality?: Maybe<Scalars['String']['output']>;
};

export type InvoiceLine = MetadataInterface & {
  metadata: Metadata;
  __typename?: 'InvoiceLine';
  contractLineItem: ServiceLineItem;
  price: Scalars['Float']['output'];
  total: Scalars['Float']['output'];
  taxDue: Scalars['Float']['output'];
  quantity: Scalars['Int64']['output'];
  subtotal: Scalars['Float']['output'];
  description: Scalars['String']['output'];
};

export type InvoiceLineSimulate = {
  key: Scalars['String']['output'];
  price: Scalars['Float']['output'];
  total: Scalars['Float']['output'];
  __typename?: 'InvoiceLineSimulate';
  taxDue: Scalars['Float']['output'];
  quantity: Scalars['Int64']['output'];
  subtotal: Scalars['Float']['output'];
  description: Scalars['String']['output'];
};

export type InvoiceProvider = {
  __typename?: 'InvoiceProvider';
  name?: Maybe<Scalars['String']['output']>;
  logoUrl?: Maybe<Scalars['String']['output']>;
  addressZip?: Maybe<Scalars['String']['output']>;
  addressLine1?: Maybe<Scalars['String']['output']>;
  addressLine2?: Maybe<Scalars['String']['output']>;
  addressRegion?: Maybe<Scalars['String']['output']>;
  addressCountry?: Maybe<Scalars['String']['output']>;
  addressLocality?: Maybe<Scalars['String']['output']>;
  logoRepositoryFileId?: Maybe<Scalars['String']['output']>;
};

export type InvoiceSimulate = {
  customer: InvoiceCustomer;
  provider: InvoiceProvider;
  __typename?: 'InvoiceSimulate';
  due: Scalars['Time']['output'];
  issued: Scalars['Time']['output'];
  note: Scalars['String']['output'];
  total: Scalars['Float']['output'];
  amount: Scalars['Float']['output'];
  taxDue: Scalars['Float']['output'];
  subtotal: Scalars['Float']['output'];
  currency: Scalars['String']['output'];
  offCycle: Scalars['Boolean']['output'];
  postpaid: Scalars['Boolean']['output'];
  invoiceNumber: Scalars['String']['output'];
  invoicePeriodEnd: Scalars['Time']['output'];
  invoiceLineItems: Array<InvoiceLineSimulate>;
  invoicePeriodStart: Scalars['Time']['output'];
};

export type InvoiceSimulateInput = {
  contractId: Scalars['ID']['input'];
  serviceLines: Array<InvoiceSimulateServiceLineInput>;
};

export type InvoiceSimulateServiceLineInput = {
  billingCycle: BilledType;
  key: Scalars['String']['input'];
  price: Scalars['Float']['input'];
  quantity: Scalars['Int64']['input'];
  description: Scalars['String']['input'];
  serviceStarted: Scalars['Time']['input'];
  parentId?: InputMaybe<Scalars['ID']['input']>;
  taxRate?: InputMaybe<Scalars['Float']['input']>;
  closeVersion?: InputMaybe<Scalars['Boolean']['input']>;
  serviceLineItemId?: InputMaybe<Scalars['ID']['input']>;
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
  content: Array<Invoice>;
  __typename?: 'InvoicesPage';
  totalPages: Scalars['Int']['output'];
  totalElements: Scalars['Int64']['output'];
  totalAvailable: Scalars['Int64']['output'];
};

export type InvoicingCycle = Node &
  SourceFields & {
    source: DataSource;
    type: InvoicingCycleType;
    sourceOfTruth: DataSource;
    id: Scalars['ID']['output'];
    __typename?: 'InvoicingCycle';
    createdAt: Scalars['Time']['output'];
    updatedAt: Scalars['Time']['output'];
    appSource: Scalars['String']['output'];
  };

export type InvoicingCycleInput = {
  type: InvoicingCycleType;
};

export enum InvoicingCycleType {
  Anniversary = 'ANNIVERSARY',
  Date = 'DATE',
}

export type InvoicingCycleUpdateInput = {
  type: InvoicingCycleType;
  id: Scalars['ID']['input'];
};

export type Issue = Node &
  SourceFields & {
    source: DataSource;
    __typename?: 'Issue';
    comments: Array<Comment>;
    sourceOfTruth: DataSource;
    id: Scalars['ID']['output'];
    tags?: Maybe<Array<Maybe<Tag>>>;
    assignedTo: Array<IssueParticipant>;
    followedBy: Array<IssueParticipant>;
    status: Scalars['String']['output'];
    createdAt: Scalars['Time']['output'];
    externalLinks: Array<ExternalSystem>;
    reportedBy?: Maybe<IssueParticipant>;
    updatedAt: Scalars['Time']['output'];
    submittedBy?: Maybe<IssueParticipant>;
    appSource: Scalars['String']['output'];
    interactionEvents: Array<InteractionEvent>;
    subject?: Maybe<Scalars['String']['output']>;
    priority?: Maybe<Scalars['String']['output']>;
    description?: Maybe<Scalars['String']['output']>;
  };

export type IssueParticipant =
  | ContactParticipant
  | OrganizationParticipant
  | UserParticipant;

export type IssueSummaryByStatus = {
  count: Scalars['Int64']['output'];
  __typename?: 'IssueSummaryByStatus';
  status: Scalars['String']['output'];
};

/**
 * Describes the relationship a Contact has with a Organization.
 * **A `return` object**
 */
export type JobRole = {
  source: DataSource;
  __typename?: 'JobRole';
  contact?: Maybe<Contact>;
  sourceOfTruth: DataSource;
  id: Scalars['ID']['output'];
  /**
   * Organization associated with a Contact.
   * **Required.**
   */
  organization?: Maybe<Organization>;
  createdAt: Scalars['Time']['output'];
  updatedAt: Scalars['Time']['output'];
  primary: Scalars['Boolean']['output'];
  appSource: Scalars['String']['output'];
  endedAt?: Maybe<Scalars['Time']['output']>;
  company?: Maybe<Scalars['String']['output']>;
  startedAt?: Maybe<Scalars['Time']['output']>;
  /** The Contact's job title. */
  jobTitle?: Maybe<Scalars['String']['output']>;
  description?: Maybe<Scalars['String']['output']>;
};

/**
 * Describes the relationship a Contact has with an Organization.
 * **A `create` object**
 */
export type JobRoleInput = {
  endedAt?: InputMaybe<Scalars['Time']['input']>;
  company?: InputMaybe<Scalars['String']['input']>;
  startedAt?: InputMaybe<Scalars['Time']['input']>;
  jobTitle?: InputMaybe<Scalars['String']['input']>;
  primary?: InputMaybe<Scalars['Boolean']['input']>;
  appSource?: InputMaybe<Scalars['String']['input']>;
  organizationId?: InputMaybe<Scalars['ID']['input']>;
  description?: InputMaybe<Scalars['String']['input']>;
};

export type JobRoleParticipant = {
  jobRoleParticipant: JobRole;
  __typename?: 'JobRoleParticipant';
  type?: Maybe<Scalars['String']['output']>;
};

/**
 * Describes the relationship a Contact has with an Organization.
 * **A `create` object**
 */
export type JobRoleUpdateInput = {
  id: Scalars['ID']['input'];
  endedAt?: InputMaybe<Scalars['Time']['input']>;
  company?: InputMaybe<Scalars['String']['input']>;
  startedAt?: InputMaybe<Scalars['Time']['input']>;
  jobTitle?: InputMaybe<Scalars['String']['input']>;
  primary?: InputMaybe<Scalars['Boolean']['input']>;
  organizationId?: InputMaybe<Scalars['ID']['input']>;
  description?: InputMaybe<Scalars['String']['input']>;
};

export type LastTouchpoint = {
  __typename?: 'LastTouchpoint';
  lastTouchPointType?: Maybe<LastTouchpointType>;
  lastTouchPointTimelineEvent?: Maybe<TimelineEvent>;
  lastTouchPointAt?: Maybe<Scalars['Time']['output']>;
  lastTouchPointTimelineEventId?: Maybe<Scalars['ID']['output']>;
};

export enum LastTouchpointType {
  Action = 'ACTION',
  ActionCreated = 'ACTION_CREATED',
  Analysis = 'ANALYSIS',
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
  subsidiaryId: Scalars['ID']['input'];
  organizationId: Scalars['ID']['input'];
  type?: InputMaybe<Scalars['String']['input']>;
};

export type LinkedOrganization = {
  organization: Organization;
  __typename?: 'LinkedOrganization';
  type?: Maybe<Scalars['String']['output']>;
};

export type Location = Node &
  SourceFields & {
    source: DataSource;
    __typename?: 'Location';
    sourceOfTruth: DataSource;
    id: Scalars['ID']['output'];
    createdAt: Scalars['Time']['output'];
    updatedAt: Scalars['Time']['output'];
    appSource: Scalars['String']['output'];
    zip?: Maybe<Scalars['String']['output']>;
    name?: Maybe<Scalars['String']['output']>;
    region?: Maybe<Scalars['String']['output']>;
    street?: Maybe<Scalars['String']['output']>;
    address?: Maybe<Scalars['String']['output']>;
    country?: Maybe<Scalars['String']['output']>;
    latitude?: Maybe<Scalars['Float']['output']>;
    address2?: Maybe<Scalars['String']['output']>;
    district?: Maybe<Scalars['String']['output']>;
    locality?: Maybe<Scalars['String']['output']>;
    longitude?: Maybe<Scalars['Float']['output']>;
    plusFour?: Maybe<Scalars['String']['output']>;
    timeZone?: Maybe<Scalars['String']['output']>;
    utcOffset?: Maybe<Scalars['Int64']['output']>;
    postalCode?: Maybe<Scalars['String']['output']>;
    rawAddress?: Maybe<Scalars['String']['output']>;
    addressType?: Maybe<Scalars['String']['output']>;
    commercial?: Maybe<Scalars['Boolean']['output']>;
    houseNumber?: Maybe<Scalars['String']['output']>;
    predirection?: Maybe<Scalars['String']['output']>;
  };

export type LocationUpdateInput = {
  id: Scalars['ID']['input'];
  zip?: InputMaybe<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  region?: InputMaybe<Scalars['String']['input']>;
  street?: InputMaybe<Scalars['String']['input']>;
  address?: InputMaybe<Scalars['String']['input']>;
  country?: InputMaybe<Scalars['String']['input']>;
  latitude?: InputMaybe<Scalars['Float']['input']>;
  address2?: InputMaybe<Scalars['String']['input']>;
  district?: InputMaybe<Scalars['String']['input']>;
  locality?: InputMaybe<Scalars['String']['input']>;
  longitude?: InputMaybe<Scalars['Float']['input']>;
  plusFour?: InputMaybe<Scalars['String']['input']>;
  timeZone?: InputMaybe<Scalars['String']['input']>;
  utcOffset?: InputMaybe<Scalars['Int64']['input']>;
  postalCode?: InputMaybe<Scalars['String']['input']>;
  rawAddress?: InputMaybe<Scalars['String']['input']>;
  addressType?: InputMaybe<Scalars['String']['input']>;
  commercial?: InputMaybe<Scalars['Boolean']['input']>;
  houseNumber?: InputMaybe<Scalars['String']['input']>;
  predirection?: InputMaybe<Scalars['String']['input']>;
};

export type LogEntry = {
  tags: Array<Tag>;
  source: DataSource;
  __typename?: 'LogEntry';
  createdBy?: Maybe<User>;
  sourceOfTruth: DataSource;
  id: Scalars['ID']['output'];
  createdAt: Scalars['Time']['output'];
  externalLinks: Array<ExternalSystem>;
  startedAt: Scalars['Time']['output'];
  updatedAt: Scalars['Time']['output'];
  appSource: Scalars['String']['output'];
  content?: Maybe<Scalars['String']['output']>;
  contentType?: Maybe<Scalars['String']['output']>;
};

export type LogEntryInput = {
  tags?: InputMaybe<Array<TagIdOrNameInput>>;
  content?: InputMaybe<Scalars['String']['input']>;
  startedAt?: InputMaybe<Scalars['Time']['input']>;
  appSource?: InputMaybe<Scalars['String']['input']>;
  contentType?: InputMaybe<Scalars['String']['input']>;
};

export type LogEntryUpdateInput = {
  content?: InputMaybe<Scalars['String']['input']>;
  startedAt?: InputMaybe<Scalars['Time']['input']>;
  contentType?: InputMaybe<Scalars['String']['input']>;
};

export enum Market {
  B2B = 'B2B',
  B2C = 'B2C',
  Marketplace = 'MARKETPLACE',
}

export type MasterPlan = Node &
  SourceFields & {
    source: DataSource;
    __typename?: 'MasterPlan';
    sourceOfTruth: DataSource;
    id: Scalars['ID']['output'];
    name: Scalars['String']['output'];
    createdAt: Scalars['Time']['output'];
    updatedAt: Scalars['Time']['output'];
    retired: Scalars['Boolean']['output'];
    appSource: Scalars['String']['output'];
    milestones: Array<MasterPlanMilestone>;
    retiredMilestones: Array<MasterPlanMilestone>;
  };

export type MasterPlanInput = {
  name?: InputMaybe<Scalars['String']['input']>;
};

export type MasterPlanMilestone = Node &
  SourceFields & {
    source: DataSource;
    sourceOfTruth: DataSource;
    id: Scalars['ID']['output'];
    name: Scalars['String']['output'];
    order: Scalars['Int64']['output'];
    __typename?: 'MasterPlanMilestone';
    createdAt: Scalars['Time']['output'];
    updatedAt: Scalars['Time']['output'];
    retired: Scalars['Boolean']['output'];
    appSource: Scalars['String']['output'];
    optional: Scalars['Boolean']['output'];
    durationHours: Scalars['Int64']['output'];
    items: Array<Scalars['String']['output']>;
  };

export type MasterPlanMilestoneInput = {
  order: Scalars['Int64']['input'];
  masterPlanId: Scalars['ID']['input'];
  optional: Scalars['Boolean']['input'];
  durationHours: Scalars['Int64']['input'];
  items: Array<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
};

export type MasterPlanMilestoneReorderInput = {
  masterPlanId: Scalars['ID']['input'];
  orderedIds: Array<Scalars['ID']['input']>;
};

export type MasterPlanMilestoneUpdateInput = {
  id: Scalars['ID']['input'];
  masterPlanId: Scalars['ID']['input'];
  name?: InputMaybe<Scalars['String']['input']>;
  order?: InputMaybe<Scalars['Int64']['input']>;
  retired?: InputMaybe<Scalars['Boolean']['input']>;
  optional?: InputMaybe<Scalars['Boolean']['input']>;
  durationHours?: InputMaybe<Scalars['Int64']['input']>;
  items?: InputMaybe<Array<Scalars['String']['input']>>;
};

export type MasterPlanUpdateInput = {
  id: Scalars['ID']['input'];
  name?: InputMaybe<Scalars['String']['input']>;
  retired?: InputMaybe<Scalars['Boolean']['input']>;
};

export type Meeting = Node & {
  note: Array<Note>;
  source: DataSource;
  status: MeetingStatus;
  __typename?: 'Meeting';
  sourceOfTruth: DataSource;
  id: Scalars['ID']['output'];
  includes: Array<Attachment>;
  describedBy: Array<Analysis>;
  recording?: Maybe<Attachment>;
  events: Array<InteractionEvent>;
  createdAt: Scalars['Time']['output'];
  createdBy: Array<MeetingParticipant>;
  updatedAt: Scalars['Time']['output'];
  attendedBy: Array<MeetingParticipant>;
  externalSystem: Array<ExternalSystem>;
  appSource: Scalars['String']['output'];
  name?: Maybe<Scalars['String']['output']>;
  endedAt?: Maybe<Scalars['Time']['output']>;
  agenda?: Maybe<Scalars['String']['output']>;
  startedAt?: Maybe<Scalars['Time']['output']>;
  conferenceUrl?: Maybe<Scalars['String']['output']>;
  agendaContentType?: Maybe<Scalars['String']['output']>;
  meetingExternalUrl?: Maybe<Scalars['String']['output']>;
};

export type MeetingInput = {
  note?: InputMaybe<NoteInput>;
  status?: InputMaybe<MeetingStatus>;
  name?: InputMaybe<Scalars['String']['input']>;
  endedAt?: InputMaybe<Scalars['Time']['input']>;
  agenda?: InputMaybe<Scalars['String']['input']>;
  createdAt?: InputMaybe<Scalars['Time']['input']>;
  startedAt?: InputMaybe<Scalars['Time']['input']>;
  appSource?: InputMaybe<Scalars['String']['input']>;
  conferenceUrl?: InputMaybe<Scalars['String']['input']>;
  createdBy?: InputMaybe<Array<MeetingParticipantInput>>;
  attendedBy?: InputMaybe<Array<MeetingParticipantInput>>;
  externalSystem?: InputMaybe<ExternalSystemReferenceInput>;
  agendaContentType?: InputMaybe<Scalars['String']['input']>;
  meetingExternalUrl?: InputMaybe<Scalars['String']['input']>;
};

export type MeetingParticipant =
  | ContactParticipant
  | EmailParticipant
  | OrganizationParticipant
  | UserParticipant;

export type MeetingParticipantInput = {
  userId?: InputMaybe<Scalars['ID']['input']>;
  contactId?: InputMaybe<Scalars['ID']['input']>;
  organizationId?: InputMaybe<Scalars['ID']['input']>;
};

export enum MeetingStatus {
  Accepted = 'ACCEPTED',
  Canceled = 'CANCELED',
  Undefined = 'UNDEFINED',
}

export type MeetingUpdateInput = {
  note?: InputMaybe<NoteUpdateInput>;
  status?: InputMaybe<MeetingStatus>;
  name?: InputMaybe<Scalars['String']['input']>;
  endedAt?: InputMaybe<Scalars['Time']['input']>;
  agenda?: InputMaybe<Scalars['String']['input']>;
  startedAt?: InputMaybe<Scalars['Time']['input']>;
  appSource?: InputMaybe<Scalars['String']['input']>;
  conferenceUrl?: InputMaybe<Scalars['String']['input']>;
  externalSystem?: InputMaybe<ExternalSystemReferenceInput>;
  agendaContentType?: InputMaybe<Scalars['String']['input']>;
  meetingExternalUrl?: InputMaybe<Scalars['String']['input']>;
};

/**
 * Specifies how many pages of meeting information has been returned in the query response.
 * **A `response` object.**
 */
export type MeetingsPage = Pages & {
  /**
   * A contact entity in customerOS.
   * **Required.  If no values it returns an empty array.**
   */
  content: Array<Meeting>;
  __typename?: 'MeetingsPage';
  /**
   * Total number of pages in the query response.
   * **Required.**
   */
  totalPages: Scalars['Int']['output'];
  /**
   * Total number of elements in the query response.
   * **Required.**
   */
  totalElements: Scalars['Int64']['output'];
};

export type Metadata = Node &
  SourceFieldsInterface & {
    source: DataSource;
    __typename?: 'Metadata';
    sourceOfTruth: DataSource;
    id: Scalars['ID']['output'];
    created: Scalars['Time']['output'];
    appSource: Scalars['String']['output'];
    lastUpdated: Scalars['Time']['output'];
  };

export type MetadataInterface = {
  metadata: Metadata;
};

export type Mutation = {
  tag_Create: Tag;
  note_Update: Note;
  user_Create: User;
  user_Update: User;
  user_AddRole: User;
  emailDelete: Result;
  note_Delete: Result;
  user_Delete: Result;
  invoice_Pay: Invoice;
  player_Merge: Result;
  invoice_Void: Invoice;
  social_Remove: Result;
  social_Update: Social;
  user_RemoveRole: User;
  contact_Merge: Contact;
  jobRole_Delete: Result;
  __typename?: 'Mutation';
  contact_Archive: Result;
  contact_Create: Contact;
  contact_Update: Contact;
  emailMergeToUser: Email;
  invoice_Update: Invoice;
  jobRole_Create: JobRole;
  jobRole_Update: JobRole;
  meeting_Create: Meeting;
  meeting_Update: Meeting;
  tag_Update?: Maybe<Tag>;
  workspace_Merge: Result;
  contract_Renew: Contract;
  emailUpdateInUser: Email;
  meeting_AddNote: Meeting;
  analysis_Create: Analysis;
  contact_AddSocial: Social;
  contract_Create: Contract;
  contract_Update: Contract;
  location_Update: Location;
  note_LinkAttachment: Note;
  contact_HardDelete: Result;
  emailMergeToContact: Email;
  tag_Delete?: Maybe<Result>;
  user_AddRoleInTenant: User;
  contact_AddTagById: Contact;
  emailRemoveFromUser: Result;
  emailUpdateInContact: Email;
  note_CreateForContact: Note;
  note_UnlinkAttachment: Note;
  user_DeleteInTenant: Result;
  attachment_Create: Attachment;
  masterPlan_Create: MasterPlan;
  masterPlan_Update: MasterPlan;
  user_RemoveRoleInTenant: User;
  contact_RemoveTagById: Contact;
  emailRemoveFromContact: Result;
  meeting_LinkRecording: Meeting;
  opportunityUpdate: Opportunity;
  organization_AddSocial: Social;
  bankAccount_Create: BankAccount;
  bankAccount_Update: BankAccount;
  contact_RemoveLocation: Contact;
  contract_Delete: DeleteResponse;
  emailMergeToOrganization: Email;
  emailRemoveFromUserById: Result;
  meeting_LinkAttachment: Meeting;
  meeting_LinkAttendedBy: Meeting;
  workspace_MergeToTenant: Result;
  contact_AddNewLocation: Location;
  contract_AddAttachment: Contract;
  emailUpdateInOrganization: Email;
  masterPlan_Duplicate: MasterPlan;
  meeting_AddNewLocation: Location;
  meeting_UnlinkRecording: Meeting;
  note_CreateForOrganization: Note;
  organization_Merge: Organization;
  fieldSetDeleteFromContact: Result;
  meeting_UnlinkAttachment: Meeting;
  meeting_UnlinkAttendedBy: Meeting;
  organization_Create: Organization;
  organization_Update: Organization;
  tableViewDef_Create: TableViewDef;
  tableViewDef_Update: TableViewDef;
  bankAccount_Delete: DeleteResponse;
  contact_RestoreFromArchive: Result;
  emailRemoveFromContactById: Result;
  contract_RemoveAttachment: Contract;
  emailRemoveFromOrganization: Result;
  location_RemoveFromContact: Contact;
  organization_SetOwner: Organization;
  phoneNumberMergeToUser: PhoneNumber;
  contact_AddOrganizationById: Contact;
  entityTemplateCreate: EntityTemplate;
  masterPlan_CreateDefault: MasterPlan;
  organization_Archive?: Maybe<Result>;
  organization_HideAll?: Maybe<Result>;
  organization_ShowAll?: Maybe<Result>;
  phoneNumberUpdateInUser: PhoneNumber;
  invoicingCycle_Create: InvoicingCycle;
  invoicingCycle_Update: InvoicingCycle;
  opportunityRenewalUpdate: Opportunity;
  organization_AddNewLocation: Location;
  organization_UnsetOwner: Organization;
  phoneNumberRemoveFromUserById: Result;
  tenant_UpdateSettings: TenantSettings;
  contact_CreateForOrganization: Contact;
  customFieldMergeToContact: CustomField;
  customer_user_AddJobRole: CustomerUser;
  phoneNumberMergeToContact: PhoneNumber;
  serviceLineItem_Delete: DeleteResponse;
  contact_RemoveOrganizationById: Contact;
  customFieldMergeToFieldSet: CustomField;
  customFieldUpdateInContact: CustomField;
  emailRemoveFromOrganizationById: Result;
  organization_ArchiveAll?: Maybe<Result>;
  phoneNumberRemoveFromUserByE164: Result;
  phoneNumberUpdateInContact: PhoneNumber;
  contractLineItem_Create: ServiceLineItem;
  contractLineItem_Update: ServiceLineItem;
  customFieldDeleteFromContactById: Result;
  customFieldUpdateInFieldSet: CustomField;
  customer_contact_Create: CustomerContact;
  fieldSetMergeToContact?: Maybe<FieldSet>;
  invoice_Simulate: Array<InvoiceSimulate>;
  logEntry_AddTag: Scalars['ID']['output'];
  logEntry_Update: Scalars['ID']['output'];
  organization_AddSubsidiary: Organization;
  phoneNumberRemoveFromContactById: Result;
  customFieldDeleteFromFieldSetById: Result;
  fieldSetUpdateInContact?: Maybe<FieldSet>;
  interactionEvent_Create: InteractionEvent;
  organizationPlan_Create: OrganizationPlan;
  organizationPlan_Update: OrganizationPlan;
  tenant_Merge: Scalars['String']['output'];
  customFieldDeleteFromContactByName: Result;
  organization_Hide: Scalars['ID']['output'];
  organization_Show: Scalars['ID']['output'];
  phoneNumberRemoveFromContactByE164: Result;
  logEntry_RemoveTag: Scalars['ID']['output'];
  logEntry_ResetTags: Scalars['ID']['output'];
  organization_RemoveSubsidiary: Organization;
  organization_UnlinkAllDomains: Organization;
  phoneNumberMergeToOrganization: PhoneNumber;
  contractLineItem_NewVersion: ServiceLineItem;
  customFieldsMergeAndUpdateInContact: Contact;
  organizationPlan_Duplicate: OrganizationPlan;
  phoneNumberUpdateInOrganization: PhoneNumber;
  interactionSession_Create: InteractionSession;
  location_RemoveFromOrganization: Organization;
  phoneNumberRemoveFromOrganizationById: Result;
  billingProfile_Create: Scalars['ID']['output'];
  billingProfile_Update: Scalars['ID']['output'];
  contact_FindEmail: Scalars['String']['output'];
  externalSystem_Create: Scalars['ID']['output'];
  tenant_AddBillingProfile: TenantBillingProfile;
  contractLineItem_Close: Scalars['ID']['output'];
  customFieldTemplate_Create: CustomFieldTemplate;
  masterPlanMilestone_Create: MasterPlanMilestone;
  masterPlanMilestone_Update: MasterPlanMilestone;
  phoneNumberRemoveFromOrganizationByE164: Result;
  tenant_hardDelete: Scalars['Boolean']['output'];
  offering_Create?: Maybe<Scalars['ID']['output']>;
  offering_Update?: Maybe<Scalars['ID']['output']>;
  reminder_Create?: Maybe<Scalars['ID']['output']>;
  reminder_Update?: Maybe<Scalars['ID']['output']>;
  billingProfile_LinkEmail: Scalars['ID']['output'];
  interactionEvent_LinkAttachment: InteractionEvent;
  organization_UpdateOnboardingStatus: Organization;
  tenant_UpdateBillingProfile: TenantBillingProfile;
  masterPlanMilestone_Duplicate: MasterPlanMilestone;
  billingProfile_UnlinkEmail: Scalars['ID']['output'];
  billingProfile_LinkLocation: Scalars['ID']['output'];
  masterPlanMilestone_Reorder: Scalars['ID']['output'];
  interactionSession_LinkAttachment: InteractionSession;
  billingProfile_UnlinkLocation: Scalars['ID']['output'];
  invoice_NextDryRunForContract: Scalars['ID']['output'];
  logEntry_CreateForOrganization: Scalars['ID']['output'];
  opportunityRenewal_UpdateAllForOrganization: Organization;
  masterPlanMilestone_BulkUpdate: Array<MasterPlanMilestone>;
  organizationPlanMilestone_Reorder: Scalars['ID']['output'];
  serviceLineItem_BulkUpdate: Array<Scalars['ID']['output']>;
  organizationPlanMilestone_Create: OrganizationPlanMilestone;
  organizationPlanMilestone_Update: OrganizationPlanMilestone;
  organizationPlanMilestone_Duplicate: OrganizationPlanMilestone;
  organizationPlanMilestone_BulkUpdate: Array<OrganizationPlanMilestone>;
};

export type MutationAnalysis_CreateArgs = {
  analysis: AnalysisInput;
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
  input: SocialInput;
  contactId: Scalars['ID']['input'];
};

export type MutationContact_AddTagByIdArgs = {
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

export type MutationContact_FindEmailArgs = {
  contactId: Scalars['ID']['input'];
  organizationId: Scalars['ID']['input'];
};

export type MutationContact_HardDeleteArgs = {
  contactId: Scalars['ID']['input'];
};

export type MutationContact_MergeArgs = {
  primaryContactId: Scalars['ID']['input'];
  mergedContactIds: Array<Scalars['ID']['input']>;
};

export type MutationContact_RemoveLocationArgs = {
  contactId: Scalars['ID']['input'];
  locationId: Scalars['ID']['input'];
};

export type MutationContact_RemoveOrganizationByIdArgs = {
  input: ContactOrganizationInput;
};

export type MutationContact_RemoveTagByIdArgs = {
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

export type MutationContractLineItem_UpdateArgs = {
  input: ServiceLineItemUpdateInput;
};

export type MutationContract_AddAttachmentArgs = {
  contractId: Scalars['ID']['input'];
  attachmentId: Scalars['ID']['input'];
};

export type MutationContract_CreateArgs = {
  input: ContractInput;
};

export type MutationContract_DeleteArgs = {
  id: Scalars['ID']['input'];
};

export type MutationContract_RemoveAttachmentArgs = {
  contractId: Scalars['ID']['input'];
  attachmentId: Scalars['ID']['input'];
};

export type MutationContract_RenewArgs = {
  input: ContractRenewalInput;
};

export type MutationContract_UpdateArgs = {
  input: ContractUpdateInput;
};

export type MutationCustomFieldDeleteFromContactByIdArgs = {
  id: Scalars['ID']['input'];
  contactId: Scalars['ID']['input'];
};

export type MutationCustomFieldDeleteFromContactByNameArgs = {
  contactId: Scalars['ID']['input'];
  fieldName: Scalars['String']['input'];
};

export type MutationCustomFieldDeleteFromFieldSetByIdArgs = {
  id: Scalars['ID']['input'];
  contactId: Scalars['ID']['input'];
  fieldSetId: Scalars['ID']['input'];
};

export type MutationCustomFieldMergeToContactArgs = {
  input: CustomFieldInput;
  contactId: Scalars['ID']['input'];
};

export type MutationCustomFieldMergeToFieldSetArgs = {
  input: CustomFieldInput;
  contactId: Scalars['ID']['input'];
  fieldSetId: Scalars['ID']['input'];
};

export type MutationCustomFieldTemplate_CreateArgs = {
  input: CustomFieldTemplateInput;
};

export type MutationCustomFieldUpdateInContactArgs = {
  input: CustomFieldUpdateInput;
  contactId: Scalars['ID']['input'];
};

export type MutationCustomFieldUpdateInFieldSetArgs = {
  input: CustomFieldUpdateInput;
  contactId: Scalars['ID']['input'];
  fieldSetId: Scalars['ID']['input'];
};

export type MutationCustomFieldsMergeAndUpdateInContactArgs = {
  contactId: Scalars['ID']['input'];
  fieldSets?: InputMaybe<Array<FieldSetInput>>;
  customFields?: InputMaybe<Array<CustomFieldInput>>;
};

export type MutationCustomer_Contact_CreateArgs = {
  input: CustomerContactInput;
};

export type MutationCustomer_User_AddJobRoleArgs = {
  id: Scalars['ID']['input'];
  jobRoleInput: JobRoleInput;
};

export type MutationEmailDeleteArgs = {
  id: Scalars['ID']['input'];
};

export type MutationEmailMergeToContactArgs = {
  input: EmailInput;
  contactId: Scalars['ID']['input'];
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

export type MutationEmailRemoveFromContactByIdArgs = {
  id: Scalars['ID']['input'];
  contactId: Scalars['ID']['input'];
};

export type MutationEmailRemoveFromOrganizationArgs = {
  email: Scalars['String']['input'];
  organizationId: Scalars['ID']['input'];
};

export type MutationEmailRemoveFromOrganizationByIdArgs = {
  id: Scalars['ID']['input'];
  organizationId: Scalars['ID']['input'];
};

export type MutationEmailRemoveFromUserArgs = {
  userId: Scalars['ID']['input'];
  email: Scalars['String']['input'];
};

export type MutationEmailRemoveFromUserByIdArgs = {
  id: Scalars['ID']['input'];
  userId: Scalars['ID']['input'];
};

export type MutationEmailUpdateInContactArgs = {
  input: EmailUpdateInput;
  contactId: Scalars['ID']['input'];
};

export type MutationEmailUpdateInOrganizationArgs = {
  input: EmailUpdateInput;
  organizationId: Scalars['ID']['input'];
};

export type MutationEmailUpdateInUserArgs = {
  input: EmailUpdateInput;
  userId: Scalars['ID']['input'];
};

export type MutationEntityTemplateCreateArgs = {
  input: EntityTemplateInput;
};

export type MutationExternalSystem_CreateArgs = {
  input: ExternalSystemInput;
};

export type MutationFieldSetDeleteFromContactArgs = {
  id: Scalars['ID']['input'];
  contactId: Scalars['ID']['input'];
};

export type MutationFieldSetMergeToContactArgs = {
  input: FieldSetInput;
  contactId: Scalars['ID']['input'];
};

export type MutationFieldSetUpdateInContactArgs = {
  input: FieldSetUpdateInput;
  contactId: Scalars['ID']['input'];
};

export type MutationInteractionEvent_CreateArgs = {
  event: InteractionEventInput;
};

export type MutationInteractionEvent_LinkAttachmentArgs = {
  eventId: Scalars['ID']['input'];
  attachmentId: Scalars['ID']['input'];
};

export type MutationInteractionSession_CreateArgs = {
  session: InteractionSessionInput;
};

export type MutationInteractionSession_LinkAttachmentArgs = {
  sessionId: Scalars['ID']['input'];
  attachmentId: Scalars['ID']['input'];
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

export type MutationInvoicingCycle_CreateArgs = {
  input: InvoicingCycleInput;
};

export type MutationInvoicingCycle_UpdateArgs = {
  input: InvoicingCycleUpdateInput;
};

export type MutationJobRole_CreateArgs = {
  input: JobRoleInput;
  contactId: Scalars['ID']['input'];
};

export type MutationJobRole_DeleteArgs = {
  roleId: Scalars['ID']['input'];
  contactId: Scalars['ID']['input'];
};

export type MutationJobRole_UpdateArgs = {
  input: JobRoleUpdateInput;
  contactId: Scalars['ID']['input'];
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
  input: TagIdOrNameInput;
  id: Scalars['ID']['input'];
};

export type MutationLogEntry_CreateForOrganizationArgs = {
  input: LogEntryInput;
  organizationId: Scalars['ID']['input'];
};

export type MutationLogEntry_RemoveTagArgs = {
  input: TagIdOrNameInput;
  id: Scalars['ID']['input'];
};

export type MutationLogEntry_ResetTagsArgs = {
  id: Scalars['ID']['input'];
  input?: InputMaybe<Array<TagIdOrNameInput>>;
};

export type MutationLogEntry_UpdateArgs = {
  id: Scalars['ID']['input'];
  input: LogEntryUpdateInput;
};

export type MutationMasterPlanMilestone_BulkUpdateArgs = {
  input: Array<MasterPlanMilestoneUpdateInput>;
};

export type MutationMasterPlanMilestone_CreateArgs = {
  input: MasterPlanMilestoneInput;
};

export type MutationMasterPlanMilestone_DuplicateArgs = {
  id: Scalars['ID']['input'];
  masterPlanId: Scalars['ID']['input'];
};

export type MutationMasterPlanMilestone_ReorderArgs = {
  input: MasterPlanMilestoneReorderInput;
};

export type MutationMasterPlanMilestone_UpdateArgs = {
  input: MasterPlanMilestoneUpdateInput;
};

export type MutationMasterPlan_CreateArgs = {
  input: MasterPlanInput;
};

export type MutationMasterPlan_CreateDefaultArgs = {
  input: MasterPlanInput;
};

export type MutationMasterPlan_DuplicateArgs = {
  id: Scalars['ID']['input'];
};

export type MutationMasterPlan_UpdateArgs = {
  input: MasterPlanUpdateInput;
};

export type MutationMeeting_AddNewLocationArgs = {
  meetingId: Scalars['ID']['input'];
};

export type MutationMeeting_AddNoteArgs = {
  note?: InputMaybe<NoteInput>;
  meetingId: Scalars['ID']['input'];
};

export type MutationMeeting_CreateArgs = {
  meeting: MeetingInput;
};

export type MutationMeeting_LinkAttachmentArgs = {
  meetingId: Scalars['ID']['input'];
  attachmentId: Scalars['ID']['input'];
};

export type MutationMeeting_LinkAttendedByArgs = {
  meetingId: Scalars['ID']['input'];
  participant: MeetingParticipantInput;
};

export type MutationMeeting_LinkRecordingArgs = {
  meetingId: Scalars['ID']['input'];
  attachmentId: Scalars['ID']['input'];
};

export type MutationMeeting_UnlinkAttachmentArgs = {
  meetingId: Scalars['ID']['input'];
  attachmentId: Scalars['ID']['input'];
};

export type MutationMeeting_UnlinkAttendedByArgs = {
  meetingId: Scalars['ID']['input'];
  participant: MeetingParticipantInput;
};

export type MutationMeeting_UnlinkRecordingArgs = {
  meetingId: Scalars['ID']['input'];
  attachmentId: Scalars['ID']['input'];
};

export type MutationMeeting_UpdateArgs = {
  meeting: MeetingUpdateInput;
  meetingId: Scalars['ID']['input'];
};

export type MutationNote_CreateForContactArgs = {
  input: NoteInput;
  contactId: Scalars['ID']['input'];
};

export type MutationNote_CreateForOrganizationArgs = {
  input: NoteInput;
  organizationId: Scalars['ID']['input'];
};

export type MutationNote_DeleteArgs = {
  id: Scalars['ID']['input'];
};

export type MutationNote_LinkAttachmentArgs = {
  noteId: Scalars['ID']['input'];
  attachmentId: Scalars['ID']['input'];
};

export type MutationNote_UnlinkAttachmentArgs = {
  noteId: Scalars['ID']['input'];
  attachmentId: Scalars['ID']['input'];
};

export type MutationNote_UpdateArgs = {
  input: NoteUpdateInput;
};

export type MutationOffering_CreateArgs = {
  input?: InputMaybe<OfferingCreateInput>;
};

export type MutationOffering_UpdateArgs = {
  input?: InputMaybe<OfferingUpdateInput>;
};

export type MutationOpportunityRenewalUpdateArgs = {
  input: OpportunityRenewalUpdateInput;
  ownerUserId?: InputMaybe<Scalars['ID']['input']>;
};

export type MutationOpportunityRenewal_UpdateAllForOrganizationArgs = {
  input: OpportunityRenewalUpdateAllForOrganizationInput;
};

export type MutationOpportunityUpdateArgs = {
  input: OpportunityUpdateInput;
};

export type MutationOrganizationPlanMilestone_BulkUpdateArgs = {
  input: Array<OrganizationPlanMilestoneUpdateInput>;
};

export type MutationOrganizationPlanMilestone_CreateArgs = {
  input: OrganizationPlanMilestoneInput;
};

export type MutationOrganizationPlanMilestone_DuplicateArgs = {
  id: Scalars['ID']['input'];
  organizationId: Scalars['ID']['input'];
  organizationPlanId: Scalars['ID']['input'];
};

export type MutationOrganizationPlanMilestone_ReorderArgs = {
  input: OrganizationPlanMilestoneReorderInput;
};

export type MutationOrganizationPlanMilestone_UpdateArgs = {
  input: OrganizationPlanMilestoneUpdateInput;
};

export type MutationOrganizationPlan_CreateArgs = {
  input: OrganizationPlanInput;
};

export type MutationOrganizationPlan_DuplicateArgs = {
  id: Scalars['ID']['input'];
  organizationId: Scalars['ID']['input'];
};

export type MutationOrganizationPlan_UpdateArgs = {
  input: OrganizationPlanUpdateInput;
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

export type MutationOrganization_ArchiveArgs = {
  id: Scalars['ID']['input'];
};

export type MutationOrganization_ArchiveAllArgs = {
  ids: Array<Scalars['ID']['input']>;
};

export type MutationOrganization_CreateArgs = {
  input: OrganizationInput;
};

export type MutationOrganization_HideArgs = {
  id: Scalars['ID']['input'];
};

export type MutationOrganization_HideAllArgs = {
  ids: Array<Scalars['ID']['input']>;
};

export type MutationOrganization_MergeArgs = {
  primaryOrganizationId: Scalars['ID']['input'];
  mergedOrganizationIds: Array<Scalars['ID']['input']>;
};

export type MutationOrganization_RemoveSubsidiaryArgs = {
  subsidiaryId: Scalars['ID']['input'];
  organizationId: Scalars['ID']['input'];
};

export type MutationOrganization_SetOwnerArgs = {
  userId: Scalars['ID']['input'];
  organizationId: Scalars['ID']['input'];
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
  input: PhoneNumberInput;
  contactId: Scalars['ID']['input'];
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
  e164: Scalars['String']['input'];
  contactId: Scalars['ID']['input'];
};

export type MutationPhoneNumberRemoveFromContactByIdArgs = {
  id: Scalars['ID']['input'];
  contactId: Scalars['ID']['input'];
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
  userId: Scalars['ID']['input'];
  e164: Scalars['String']['input'];
};

export type MutationPhoneNumberRemoveFromUserByIdArgs = {
  id: Scalars['ID']['input'];
  userId: Scalars['ID']['input'];
};

export type MutationPhoneNumberUpdateInContactArgs = {
  input: PhoneNumberUpdateInput;
  contactId: Scalars['ID']['input'];
};

export type MutationPhoneNumberUpdateInOrganizationArgs = {
  input: PhoneNumberUpdateInput;
  organizationId: Scalars['ID']['input'];
};

export type MutationPhoneNumberUpdateInUserArgs = {
  input: PhoneNumberUpdateInput;
  userId: Scalars['ID']['input'];
};

export type MutationPlayer_MergeArgs = {
  input: PlayerInput;
  userId: Scalars['ID']['input'];
};

export type MutationReminder_CreateArgs = {
  input: ReminderInput;
};

export type MutationReminder_UpdateArgs = {
  input: ReminderUpdateInput;
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

export type MutationTableViewDef_CreateArgs = {
  input: TableViewDefCreateInput;
};

export type MutationTableViewDef_UpdateArgs = {
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

export type MutationTenant_MergeArgs = {
  tenant: TenantInput;
};

export type MutationTenant_UpdateBillingProfileArgs = {
  input: TenantBillingProfileUpdateInput;
};

export type MutationTenant_UpdateSettingsArgs = {
  input?: InputMaybe<TenantSettingsInput>;
};

export type MutationTenant_HardDeleteArgs = {
  tenant: Scalars['String']['input'];
  confirmTenant: Scalars['String']['input'];
};

export type MutationUser_AddRoleArgs = {
  role: Role;
  id: Scalars['ID']['input'];
};

export type MutationUser_AddRoleInTenantArgs = {
  role: Role;
  id: Scalars['ID']['input'];
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
  role: Role;
  id: Scalars['ID']['input'];
};

export type MutationUser_RemoveRoleInTenantArgs = {
  role: Role;
  id: Scalars['ID']['input'];
  tenant: Scalars['String']['input'];
};

export type MutationUser_UpdateArgs = {
  input: UserUpdateInput;
};

export type MutationWorkspace_MergeArgs = {
  workspace: WorkspaceInput;
};

export type MutationWorkspace_MergeToTenantArgs = {
  workspace: WorkspaceInput;
  tenant: Scalars['String']['input'];
};

export type Node = {
  id: Scalars['ID']['output'];
};

export type Note = {
  source: DataSource;
  __typename?: 'Note';
  createdBy?: Maybe<User>;
  noted: Array<NotedEntity>;
  sourceOfTruth: DataSource;
  id: Scalars['ID']['output'];
  includes: Array<Attachment>;
  createdAt: Scalars['Time']['output'];
  updatedAt: Scalars['Time']['output'];
  appSource: Scalars['String']['output'];
  content?: Maybe<Scalars['String']['output']>;
  contentType?: Maybe<Scalars['String']['output']>;
};

export type NoteInput = {
  content?: InputMaybe<Scalars['String']['input']>;
  appSource?: InputMaybe<Scalars['String']['input']>;
  contentType?: InputMaybe<Scalars['String']['input']>;
};

export type NotePage = Pages & {
  content: Array<Note>;
  __typename?: 'NotePage';
  totalPages: Scalars['Int']['output'];
  totalElements: Scalars['Int64']['output'];
};

export type NoteUpdateInput = {
  id: Scalars['ID']['input'];
  content?: InputMaybe<Scalars['String']['input']>;
  contentType?: InputMaybe<Scalars['String']['input']>;
};

export type NotedEntity = Contact | Organization;

export type Offering = MetadataInterface & {
  metadata: Metadata;
  __typename?: 'Offering';
  conditionals: Conditionals;
  currency?: Maybe<Currency>;
  type?: Maybe<OfferingType>;
  name: Scalars['String']['output'];
  price: Scalars['Float']['output'];
  priceCalculation: PriceCalculation;
  pricingModel?: Maybe<PricingModel>;
  active: Scalars['Boolean']['output'];
  externalLinks: Array<ExternalSystem>;
  taxable: Scalars['Boolean']['output'];
  conditional: Scalars['Boolean']['output'];
  priceCalculated: Scalars['Boolean']['output'];
  pricingPeriodInMonths: Scalars['Int64']['output'];
};

export type OfferingCreateInput = {
  currency?: InputMaybe<Currency>;
  type?: InputMaybe<OfferingType>;
  pricingModel?: InputMaybe<PricingModel>;
  name?: InputMaybe<Scalars['String']['input']>;
  price?: InputMaybe<Scalars['Float']['input']>;
  active?: InputMaybe<Scalars['Boolean']['input']>;
  taxable?: InputMaybe<Scalars['Boolean']['input']>;
  priceCalculationType?: InputMaybe<CalculationType>;
  conditional?: InputMaybe<Scalars['Boolean']['input']>;
  priceCalculated?: InputMaybe<Scalars['Boolean']['input']>;
  conditionalsMinimumChargePeriod?: InputMaybe<ChargePeriod>;
  pricingPeriodInMonths?: InputMaybe<Scalars['Int64']['input']>;
  conditionalsMinimumChargeAmount?: InputMaybe<Scalars['Float']['input']>;
  priceCalculationRevenueSharePercentage?: InputMaybe<
    Scalars['Float']['input']
  >;
};

export enum OfferingType {
  Product = 'PRODUCT',
  Service = 'SERVICE',
}

export type OfferingUpdateInput = {
  id: Scalars['ID']['input'];
  currency?: InputMaybe<Currency>;
  type?: InputMaybe<OfferingType>;
  pricingModel?: InputMaybe<PricingModel>;
  name?: InputMaybe<Scalars['String']['input']>;
  price?: InputMaybe<Scalars['Float']['input']>;
  active?: InputMaybe<Scalars['Boolean']['input']>;
  taxable?: InputMaybe<Scalars['Boolean']['input']>;
  priceCalculationType?: InputMaybe<CalculationType>;
  conditional?: InputMaybe<Scalars['Boolean']['input']>;
  priceCalculated?: InputMaybe<Scalars['Boolean']['input']>;
  conditionalsMinimumChargePeriod?: InputMaybe<ChargePeriod>;
  pricingPeriodInMonths?: InputMaybe<Scalars['Int64']['input']>;
  conditionalsMinimumChargeAmount?: InputMaybe<Scalars['Float']['input']>;
  priceCalculationRevenueSharePercentage?: InputMaybe<
    Scalars['Float']['input']
  >;
};

export type OnboardingDetails = {
  status: OnboardingStatus;
  __typename?: 'OnboardingDetails';
  updatedAt?: Maybe<Scalars['Time']['output']>;
  comments?: Maybe<Scalars['String']['output']>;
};

export enum OnboardingPlanMilestoneItemStatus {
  Done = 'DONE',
  DoneLate = 'DONE_LATE',
  NotDone = 'NOT_DONE',
  NotDoneLate = 'NOT_DONE_LATE',
  Skipped = 'SKIPPED',
  SkippedLate = 'SKIPPED_LATE',
}

export enum OnboardingPlanMilestoneStatus {
  Done = 'DONE',
  DoneLate = 'DONE_LATE',
  NotStarted = 'NOT_STARTED',
  NotStartedLate = 'NOT_STARTED_LATE',
  Started = 'STARTED',
  StartedLate = 'STARTED_LATE',
}

export enum OnboardingPlanStatus {
  Done = 'DONE',
  DoneLate = 'DONE_LATE',
  Late = 'LATE',
  NotStarted = 'NOT_STARTED',
  NotStartedLate = 'NOT_STARTED_LATE',
  OnTrack = 'ON_TRACK',
}

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
  status: OnboardingStatus;
  organizationId: Scalars['ID']['input'];
  comments?: InputMaybe<Scalars['String']['input']>;
};

export type Opportunity = Node & {
  source: DataSource;
  owner?: Maybe<User>;
  createdBy?: Maybe<User>;
  sourceOfTruth: DataSource;
  __typename?: 'Opportunity';
  internalType: InternalType;
  id: Scalars['ID']['output'];
  internalStage: InternalStage;
  name: Scalars['String']['output'];
  amount: Scalars['Float']['output'];
  createdAt: Scalars['Time']['output'];
  externalLinks: Array<ExternalSystem>;
  updatedAt: Scalars['Time']['output'];
  comments: Scalars['String']['output'];
  maxAmount: Scalars['Float']['output'];
  appSource: Scalars['String']['output'];
  nextSteps: Scalars['String']['output'];
  externalType: Scalars['String']['output'];
  generalNotes: Scalars['String']['output'];
  externalStage: Scalars['String']['output'];
  renewedAt?: Maybe<Scalars['Time']['output']>;
  renewalApproved: Scalars['Boolean']['output'];
  renewalAdjustedRate: Scalars['Int64']['output'];
  renewalLikelihood: OpportunityRenewalLikelihood;
  renewalUpdatedByUserId: Scalars['String']['output'];
  estimatedClosedAt?: Maybe<Scalars['Time']['output']>;
  renewalUpdatedByUserAt?: Maybe<Scalars['Time']['output']>;
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
  opportunityId: Scalars['ID']['input'];
  /** Deprecated */
  name?: InputMaybe<Scalars['String']['input']>;
  amount?: InputMaybe<Scalars['Float']['input']>;
  ownerUserId?: InputMaybe<Scalars['ID']['input']>;
  comments?: InputMaybe<Scalars['String']['input']>;
  appSource?: InputMaybe<Scalars['String']['input']>;
  renewalAdjustedRate?: InputMaybe<Scalars['Int64']['input']>;
  renewalLikelihood?: InputMaybe<OpportunityRenewalLikelihood>;
};

export type OpportunityUpdateInput = {
  opportunityId: Scalars['ID']['input'];
  name?: InputMaybe<Scalars['String']['input']>;
  amount?: InputMaybe<Scalars['Float']['input']>;
  appSource?: InputMaybe<Scalars['String']['input']>;
  nextSteps?: InputMaybe<Scalars['String']['input']>;
  externalType?: InputMaybe<Scalars['String']['input']>;
  generalNotes?: InputMaybe<Scalars['String']['input']>;
  externalStage?: InputMaybe<Scalars['String']['input']>;
  estimatedClosedDate?: InputMaybe<Scalars['Time']['input']>;
  externalReference?: InputMaybe<ExternalSystemReferenceInput>;
};

export type Order = {
  source: DataSource;
  __typename?: 'Order';
  sourceOfTruth: DataSource;
  id: Scalars['ID']['output'];
  createdAt: Scalars['Time']['output'];
  appSource: Scalars['String']['output'];
  paidAt?: Maybe<Scalars['Time']['output']>;
  cancelledAt?: Maybe<Scalars['Time']['output']>;
  confirmedAt?: Maybe<Scalars['Time']['output']>;
  fulfilledAt?: Maybe<Scalars['Time']['output']>;
};

export type OrgAccountDetails = {
  __typename?: 'OrgAccountDetails';
  onboarding?: Maybe<OnboardingDetails>;
  renewalSummary?: Maybe<RenewalSummary>;
};

export type Organization = MetadataInterface & {
  metadata: Metadata;
  /**
   * Deprecated
   * @deprecated Use metadata.source
   */
  source: DataSource;
  owner?: Maybe<User>;
  emails: Array<Email>;
  orders: Array<Order>;
  contacts: ContactsPage;
  market?: Maybe<Market>;
  /**
   * Deprecated
   * @deprecated Use socialMedia
   */
  socials: Array<Social>;
  jobRoles: Array<JobRole>;
  tags?: Maybe<Array<Tag>>;
  /**
   * Deprecated
   * @deprecated Use metadata.sourceOfTruth
   */
  sourceOfTruth: DataSource;
  fieldSets: Array<FieldSet>;
  locations: Array<Location>;
  socialMedia: Array<Social>;
  __typename?: 'Organization';
  /**
   * Deprecated
   * @deprecated Use metadata.id
   */
  id: Scalars['ID']['output'];
  customFields: Array<CustomField>;
  phoneNumbers: Array<PhoneNumber>;
  stage?: Maybe<OrganizationStage>;
  name: Scalars['String']['output'];
  contracts?: Maybe<Array<Contract>>;
  hide: Scalars['Boolean']['output'];
  /**
   * Deprecated
   * @deprecated Use metadata.created
   */
  createdAt: Scalars['Time']['output'];
  externalLinks: Array<ExternalSystem>;
  timelineEvents: Array<TimelineEvent>;
  /**
   * Deprecated
   * @deprecated Use metadata.lastUpdated
   */
  updatedAt: Scalars['Time']['output'];
  /**
   * Deprecated
   * @deprecated Use metadata.appSource
   */
  appSource: Scalars['String']['output'];
  entityTemplate?: Maybe<EntityTemplate>;
  lastFundingRound?: Maybe<FundingRound>;
  lastTouchpoint?: Maybe<LastTouchpoint>;
  subsidiaries: Array<LinkedOrganization>;
  /**
   * Deprecated
   * @deprecated Use parentCompany
   */
  subsidiaryOf: Array<LinkedOrganization>;
  contactCount: Scalars['Int64']['output'];
  accountDetails?: Maybe<OrgAccountDetails>;
  customerOsId: Scalars['String']['output'];
  icon?: Maybe<Scalars['String']['output']>;
  logo?: Maybe<Scalars['String']['output']>;
  /**
   * Deprecated
   * @deprecated Use notes
   */
  note?: Maybe<Scalars['String']['output']>;
  notes?: Maybe<Scalars['String']['output']>;
  parentCompanies: Array<LinkedOrganization>;
  domains: Array<Scalars['String']['output']>;
  /**
   * Deprecated
   * @deprecated Use logo
   */
  logoUrl?: Maybe<Scalars['String']['output']>;
  public?: Maybe<Scalars['Boolean']['output']>;
  website?: Maybe<Scalars['String']['output']>;
  customId?: Maybe<Scalars['String']['output']>;
  employees?: Maybe<Scalars['Int64']['output']>;
  inboundCommsCount: Scalars['Int64']['output'];
  industry?: Maybe<Scalars['String']['output']>;
  /**
   * Deprecated
   * @deprecated Use public
   */
  isPublic?: Maybe<Scalars['Boolean']['output']>;
  /**
   * Deprecated
   * @deprecated Use lastTouchpoint
   */
  lastTouchPointType?: Maybe<LastTouchpointType>;
  outboundCommsCount: Scalars['Int64']['output'];
  relationship?: Maybe<OrganizationRelationship>;
  leadSource?: Maybe<Scalars['String']['output']>;
  yearFounded?: Maybe<Scalars['Int64']['output']>;
  description?: Maybe<Scalars['String']['output']>;
  /**
   * Deprecated, use relationship instead
   * @deprecated Use relationship
   */
  isCustomer?: Maybe<Scalars['Boolean']['output']>;
  /**
   * Deprecated
   * @deprecated Use customId
   */
  referenceId?: Maybe<Scalars['String']['output']>;
  subIndustry?: Maybe<Scalars['String']['output']>;
  headquarters?: Maybe<Scalars['String']['output']>;
  issueSummaryByStatus: Array<IssueSummaryByStatus>;
  industryGroup?: Maybe<Scalars['String']['output']>;
  /**
   * Deprecated
   * @deprecated Use lastTouchpoint
   */
  lastTouchPointTimelineEvent?: Maybe<TimelineEvent>;
  /**
   * Deprecated
   * @deprecated Use lastTouchpoint
   */
  lastTouchPointAt?: Maybe<Scalars['Time']['output']>;
  slackChannelId?: Maybe<Scalars['String']['output']>;
  stageLastUpdated?: Maybe<Scalars['Time']['output']>;
  suggestedMergeTo: Array<SuggestedMergeOrganization>;
  targetAudience?: Maybe<Scalars['String']['output']>;
  timelineEventsTotalCount: Scalars['Int64']['output'];
  valueProposition?: Maybe<Scalars['String']['output']>;
  lastFundingAmount?: Maybe<Scalars['String']['output']>;
  employeeGrowthRate?: Maybe<Scalars['String']['output']>;
  /** Deprecated */
  lastTouchPointTimelineEventId?: Maybe<Scalars['ID']['output']>;
};

export type OrganizationContactsArgs = {
  where?: InputMaybe<Filter>;
  sort?: InputMaybe<Array<SortBy>>;
  pagination?: InputMaybe<Pagination>;
};

export type OrganizationTimelineEventsArgs = {
  size: Scalars['Int']['input'];
  from?: InputMaybe<Scalars['Time']['input']>;
  timelineEventTypes?: InputMaybe<Array<TimelineEventType>>;
};

export type OrganizationTimelineEventsTotalCountArgs = {
  timelineEventTypes?: InputMaybe<Array<TimelineEventType>>;
};

export type OrganizationInput = {
  market?: InputMaybe<Market>;
  stage?: InputMaybe<OrganizationStage>;
  /** Deprecated */
  fieldSets?: InputMaybe<Array<FieldSetInput>>;
  icon?: InputMaybe<Scalars['String']['input']>;
  logo?: InputMaybe<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  note?: InputMaybe<Scalars['String']['input']>;
  notes?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  templateId?: InputMaybe<Scalars['ID']['input']>;
  /** Deprecated */
  logoUrl?: InputMaybe<Scalars['String']['input']>;
  public?: InputMaybe<Scalars['Boolean']['input']>;
  website?: InputMaybe<Scalars['String']['input']>;
  /**
   * The name of the organization.
   * **Required.**
   */
  customId?: InputMaybe<Scalars['String']['input']>;
  employees?: InputMaybe<Scalars['Int64']['input']>;
  industry?: InputMaybe<Scalars['String']['input']>;
  appSource?: InputMaybe<Scalars['String']['input']>;
  customFields?: InputMaybe<Array<CustomFieldInput>>;
  /** Deprecated */
  isPublic?: InputMaybe<Scalars['Boolean']['input']>;
  leadSource?: InputMaybe<Scalars['String']['input']>;
  relationship?: InputMaybe<OrganizationRelationship>;
  yearFounded?: InputMaybe<Scalars['Int64']['input']>;
  description?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated, use relationship instead */
  isCustomer?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  referenceId?: InputMaybe<Scalars['String']['input']>;
  subIndustry?: InputMaybe<Scalars['String']['input']>;
  headquarters?: InputMaybe<Scalars['String']['input']>;
  industryGroup?: InputMaybe<Scalars['String']['input']>;
  domains?: InputMaybe<Array<Scalars['String']['input']>>;
  slackChannelId?: InputMaybe<Scalars['String']['input']>;
  employeeGrowthRate?: InputMaybe<Scalars['String']['input']>;
};

export type OrganizationPage = Pages & {
  content: Array<Organization>;
  __typename?: 'OrganizationPage';
  totalPages: Scalars['Int']['output'];
  totalElements: Scalars['Int64']['output'];
  totalAvailable: Scalars['Int64']['output'];
};

export type OrganizationParticipant = {
  organizationParticipant: Organization;
  __typename?: 'OrganizationParticipant';
  type?: Maybe<Scalars['String']['output']>;
};

export type OrganizationPlan = Node &
  SourceFields & {
    source: DataSource;
    sourceOfTruth: DataSource;
    id: Scalars['ID']['output'];
    __typename?: 'OrganizationPlan';
    name: Scalars['String']['output'];
    createdAt: Scalars['Time']['output'];
    updatedAt: Scalars['Time']['output'];
    masterPlanId: Scalars['ID']['output'];
    retired: Scalars['Boolean']['output'];
    appSource: Scalars['String']['output'];
    milestones: Array<OrganizationPlanMilestone>;
    statusDetails: OrganizationPlanStatusDetails;
    retiredMilestones: Array<OrganizationPlanMilestone>;
  };

export type OrganizationPlanInput = {
  organizationId: Scalars['ID']['input'];
  name?: InputMaybe<Scalars['String']['input']>;
  masterPlanId?: InputMaybe<Scalars['String']['input']>;
};

export type OrganizationPlanMilestone = Node &
  SourceFields & {
    source: DataSource;
    sourceOfTruth: DataSource;
    id: Scalars['ID']['output'];
    name: Scalars['String']['output'];
    order: Scalars['Int64']['output'];
    dueDate: Scalars['Time']['output'];
    adhoc: Scalars['Boolean']['output'];
    createdAt: Scalars['Time']['output'];
    updatedAt: Scalars['Time']['output'];
    retired: Scalars['Boolean']['output'];
    appSource: Scalars['String']['output'];
    optional: Scalars['Boolean']['output'];
    __typename?: 'OrganizationPlanMilestone';
    items: Array<OrganizationPlanMilestoneItem>;
    statusDetails: OrganizationPlanMilestoneStatusDetails;
  };

export type OrganizationPlanMilestoneInput = {
  order: Scalars['Int64']['input'];
  dueDate: Scalars['Time']['input'];
  adhoc: Scalars['Boolean']['input'];
  createdAt: Scalars['Time']['input'];
  optional: Scalars['Boolean']['input'];
  organizationId: Scalars['ID']['input'];
  items: Array<Scalars['String']['input']>;
  organizationPlanId: Scalars['ID']['input'];
  name?: InputMaybe<Scalars['String']['input']>;
};

export type OrganizationPlanMilestoneItem = {
  uuid: Scalars['ID']['output'];
  text: Scalars['String']['output'];
  updatedAt: Scalars['Time']['output'];
  status: OnboardingPlanMilestoneItemStatus;
  __typename?: 'OrganizationPlanMilestoneItem';
};

export type OrganizationPlanMilestoneItemInput = {
  text: Scalars['String']['input'];
  updatedAt: Scalars['Time']['input'];
  status: OnboardingPlanMilestoneItemStatus;
  uuid?: InputMaybe<Scalars['ID']['input']>;
};

export type OrganizationPlanMilestoneReorderInput = {
  organizationId: Scalars['ID']['input'];
  orderedIds: Array<Scalars['ID']['input']>;
  organizationPlanId: Scalars['ID']['input'];
};

export type OrganizationPlanMilestoneStatusDetails = {
  text: Scalars['String']['output'];
  updatedAt: Scalars['Time']['output'];
  status: OnboardingPlanMilestoneStatus;
  __typename?: 'OrganizationPlanMilestoneStatusDetails';
};

export type OrganizationPlanMilestoneStatusDetailsInput = {
  text: Scalars['String']['input'];
  updatedAt: Scalars['Time']['input'];
  status: OnboardingPlanMilestoneStatus;
};

export type OrganizationPlanMilestoneUpdateInput = {
  id: Scalars['ID']['input'];
  updatedAt: Scalars['Time']['input'];
  organizationId: Scalars['ID']['input'];
  organizationPlanId: Scalars['ID']['input'];
  name?: InputMaybe<Scalars['String']['input']>;
  order?: InputMaybe<Scalars['Int64']['input']>;
  dueDate?: InputMaybe<Scalars['Time']['input']>;
  adhoc?: InputMaybe<Scalars['Boolean']['input']>;
  retired?: InputMaybe<Scalars['Boolean']['input']>;
  optional?: InputMaybe<Scalars['Boolean']['input']>;
  statusDetails?: InputMaybe<OrganizationPlanMilestoneStatusDetailsInput>;
  items?: InputMaybe<Array<InputMaybe<OrganizationPlanMilestoneItemInput>>>;
};

export type OrganizationPlanStatusDetails = {
  status: OnboardingPlanStatus;
  text: Scalars['String']['output'];
  updatedAt: Scalars['Time']['output'];
  __typename?: 'OrganizationPlanStatusDetails';
};

export type OrganizationPlanStatusDetailsInput = {
  status: OnboardingPlanStatus;
  text: Scalars['String']['input'];
  updatedAt: Scalars['Time']['input'];
};

export type OrganizationPlanUpdateInput = {
  id: Scalars['ID']['input'];
  organizationId: Scalars['ID']['input'];
  name?: InputMaybe<Scalars['String']['input']>;
  retired?: InputMaybe<Scalars['Boolean']['input']>;
  statusDetails?: InputMaybe<OrganizationPlanStatusDetailsInput>;
};

export enum OrganizationRelationship {
  Customer = 'CUSTOMER',
  FormerCustomer = 'FORMER_CUSTOMER',
  NotAFit = 'NOT_A_FIT',
  Prospect = 'PROSPECT',
}

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
  Unqualified = 'UNQUALIFIED',
}

export type OrganizationUpdateInput = {
  id: Scalars['ID']['input'];
  market?: InputMaybe<Market>;
  stage?: InputMaybe<OrganizationStage>;
  lastFundingRound?: InputMaybe<FundingRound>;
  icon?: InputMaybe<Scalars['String']['input']>;
  logo?: InputMaybe<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  /** Deprecatedm, use notes instead */
  note?: InputMaybe<Scalars['String']['input']>;
  notes?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  patch?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated, use logo instead */
  logoUrl?: InputMaybe<Scalars['String']['input']>;
  public?: InputMaybe<Scalars['Boolean']['input']>;
  website?: InputMaybe<Scalars['String']['input']>;
  customId?: InputMaybe<Scalars['String']['input']>;
  employees?: InputMaybe<Scalars['Int64']['input']>;
  industry?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated, use public instead */
  isPublic?: InputMaybe<Scalars['Boolean']['input']>;
  relationship?: InputMaybe<OrganizationRelationship>;
  yearFounded?: InputMaybe<Scalars['Int64']['input']>;
  description?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated, use relationship instead */
  isCustomer?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated, use customId instead */
  referenceId?: InputMaybe<Scalars['String']['input']>;
  subIndustry?: InputMaybe<Scalars['String']['input']>;
  headquarters?: InputMaybe<Scalars['String']['input']>;
  industryGroup?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  domains?: InputMaybe<Array<Scalars['String']['input']>>;
  slackChannelId?: InputMaybe<Scalars['String']['input']>;
  targetAudience?: InputMaybe<Scalars['String']['input']>;
  valueProposition?: InputMaybe<Scalars['String']['input']>;
  lastFundingAmount?: InputMaybe<Scalars['String']['input']>;
  employeeGrowthRate?: InputMaybe<Scalars['String']['input']>;
};

export type PageView = Node &
  SourceFields & {
    source: DataSource;
    __typename?: 'PageView';
    sourceOfTruth: DataSource;
    id: Scalars['ID']['output'];
    endedAt: Scalars['Time']['output'];
    sessionId: Scalars['ID']['output'];
    pageUrl: Scalars['String']['output'];
    startedAt: Scalars['Time']['output'];
    appSource: Scalars['String']['output'];
    pageTitle: Scalars['String']['output'];
    engagedTime: Scalars['Int64']['output'];
    application: Scalars['String']['output'];
    orderInSession: Scalars['Int64']['output'];
  };

/**
 * Describes the number of pages and total elements included in a query response.
 * **A `response` object.**
 */
export type Pages = {
  /**
   * The total number of pages included in the query response.
   * **Required.**
   */
  totalPages: Scalars['Int']['output'];
  /**
   * The total number of elements included in the query response.
   * **Required.**
   */
  totalElements: Scalars['Int64']['output'];
};

/** If provided as part of the request, results will be filtered down to the `page` and `limit` specified. */
export type Pagination = {
  /**
   * The results page to return in the response.
   * **Required.**
   */
  page: Scalars['Int']['input'];
  /**
   * The maximum number of results in the response.
   * **Required.**
   */
  limit: Scalars['Int']['input'];
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
  source: DataSource;
  users: Array<User>;
  contacts: Array<Contact>;
  country?: Maybe<Country>;
  __typename?: 'PhoneNumber';
  /**
   * The unique ID associated with the phone number.
   * **Required**
   */
  id: Scalars['ID']['output'];
  /** Defines the type of phone number. */
  label?: Maybe<PhoneNumberLabel>;
  organizations: Array<Organization>;
  createdAt: Scalars['Time']['output'];
  updatedAt: Scalars['Time']['output'];
  /**
   * Determines if the phone number is primary or not.
   * **Required**
   */
  primary: Scalars['Boolean']['output'];
  /** The phone number in e164 format.  */
  e164?: Maybe<Scalars['String']['output']>;
  appSource?: Maybe<Scalars['String']['output']>;
  validated?: Maybe<Scalars['Boolean']['output']>;
  rawPhoneNumber?: Maybe<Scalars['String']['output']>;
};

/**
 * Describes a phone number associated with a `Contact` in customerOS.
 * **A `create` object.**
 */
export type PhoneNumberInput = {
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
  countryCodeA2?: InputMaybe<Scalars['String']['input']>;
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
  phoneNumberParticipant: PhoneNumber;
  __typename?: 'PhoneNumberParticipant';
  type?: Maybe<Scalars['String']['output']>;
};

/**
 * Describes a phone number associated with a `Contact` in customerOS.
 * **An `update` object.**
 */
export type PhoneNumberUpdateInput = {
  /**
   * The unique ID associated with the phone number.
   * **Required**
   */
  id: Scalars['ID']['input'];
  /** Defines the type of phone number. */
  label?: InputMaybe<PhoneNumberLabel>;
  /**
   * Determines if the phone number is primary or not.
   * **Required**
   */
  primary?: InputMaybe<Scalars['Boolean']['input']>;
  phoneNumber?: InputMaybe<Scalars['String']['input']>;
  countryCodeA2?: InputMaybe<Scalars['String']['input']>;
};

export type Player = {
  source: DataSource;
  __typename?: 'Player';
  users: Array<PlayerUser>;
  sourceOfTruth: DataSource;
  id: Scalars['ID']['output'];
  authId: Scalars['String']['output'];
  createdAt: Scalars['Time']['output'];
  updatedAt: Scalars['Time']['output'];
  provider: Scalars['String']['output'];
  appSource: Scalars['String']['output'];
  identityId?: Maybe<Scalars['String']['output']>;
};

export type PlayerInput = {
  authId: Scalars['String']['input'];
  provider: Scalars['String']['input'];
  appSource?: InputMaybe<Scalars['String']['input']>;
  identityId?: InputMaybe<Scalars['String']['input']>;
};

export type PlayerUpdate = {
  appSource?: InputMaybe<Scalars['String']['input']>;
  identityId?: InputMaybe<Scalars['String']['input']>;
};

export type PlayerUser = {
  user: User;
  __typename?: 'PlayerUser';
  tenant: Scalars['String']['output'];
  default: Scalars['Boolean']['output'];
};

export type PriceCalculation = {
  __typename?: 'PriceCalculation';
  calculationType?: Maybe<CalculationType>;
  revenueSharePercentage: Scalars['Float']['output'];
};

export enum PricingModel {
  OneTime = 'ONE_TIME',
  Subscription = 'SUBSCRIPTION',
  Usage = 'USAGE',
}

export type Query = {
  user: User;
  email: Email;
  issue: Issue;
  users: UserPage;
  invoice: Invoice;
  meeting: Meeting;
  tags: Array<Tag>;
  analysis: Analysis;
  contract: Contract;
  logEntry: LogEntry;
  reminder: Reminder;
  user_ByEmail: User;
  __typename?: 'Query';
  attachment: Attachment;
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
  invoices: InvoicesPage;
  masterPlan: MasterPlan;
  contracts: ContractPage;
  /** Fetch a single contact from customerOS by contact ID. */
  contact?: Maybe<Contact>;
  contact_ByEmail: Contact;
  contact_ByPhone: Contact;
  phoneNumber: PhoneNumber;
  global_Cache: GlobalCache;
  invoice_ByNumber: Invoice;
  offerings: Array<Offering>;
  gcli_Search: Array<GCliItem>;
  externalMeetings: MeetingsPage;
  invoicingCycle: InvoicingCycle;
  masterPlans: Array<MasterPlan>;
  tenantSettings: TenantSettings;
  organizations: OrganizationPage;
  player_ByAuthIdProvider: Player;
  bankAccounts: Array<BankAccount>;
  billableInfo: TenantBillableInfo;
  opportunity?: Maybe<Opportunity>;
  serviceLineItem: ServiceLineItem;
  slack_Channels: SlackChannelPage;
  interactionEvent: InteractionEvent;
  organization?: Maybe<Organization>;
  organizationPlan: OrganizationPlan;
  tableViewDefs: Array<TableViewDef>;
  tenant: Scalars['String']['output'];
  timelineEvents: Array<TimelineEvent>;
  entityTemplates: Array<EntityTemplate>;
  interactionSession: InteractionSession;
  organization_DistinctOwners: Array<User>;
  remindersForOrganization: Array<Reminder>;
  organizationPlans: Array<OrganizationPlan>;
  tenantBillingProfile: TenantBillingProfile;
  dashboardView_Renewals?: Maybe<RenewalsPage>;
  organization_ByCustomId?: Maybe<Organization>;
  organization_ByCustomerOsId?: Maybe<Organization>;
  tenantBillingProfiles: Array<TenantBillingProfile>;
  tenant_ByEmail?: Maybe<Scalars['String']['output']>;
  interactionEvent_ByEventIdentifier: InteractionEvent;
  /** sort.By available options: ORGANIZATION, IS_CUSTOMER, DOMAIN, LOCATION, OWNER, LAST_TOUCHPOINT, RENEWAL_LIKELIHOOD, FORECAST_ARR, RENEWAL_DATE, ONBOARDING_STATUS */
  dashboardView_Organizations?: Maybe<OrganizationPage>;
  dashboard_ARRBreakdown?: Maybe<DashboardArrBreakdown>;
  dashboard_NewCustomers?: Maybe<DashboardNewCustomers>;
  externalSystemInstances: Array<ExternalSystemInstance>;
  dashboard_RetentionRate?: Maybe<DashboardRetentionRate>;
  dashboard_RevenueAtRisk?: Maybe<DashboardRevenueAtRisk>;
  dashboard_TimeToOnboard?: Maybe<DashboardTimeToOnboard>;
  tenant_ByWorkspace?: Maybe<Scalars['String']['output']>;
  interactionSession_ByEventIdentifier: InteractionSession;
  dashboard_MRRPerCustomer?: Maybe<DashboardMrrPerCustomer>;
  organizationPlansForOrganization: Array<OrganizationPlan>;
  dashboard_CustomerMap?: Maybe<Array<DashboardCustomerMap>>;
  interactionSession_BySessionIdentifier: InteractionSession;
  dashboard_OnboardingCompletion?: Maybe<DashboardOnboardingCompletion>;
  dashboard_GrossRevenueRetention?: Maybe<DashboardGrossRevenueRetention>;
};

export type QueryAnalysisArgs = {
  id: Scalars['ID']['input'];
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
  where?: InputMaybe<Filter>;
  sort?: InputMaybe<Array<SortBy>>;
  pagination?: InputMaybe<Pagination>;
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

export type QueryEntityTemplatesArgs = {
  extends?: InputMaybe<EntityTemplateExtension>;
};

export type QueryExternalMeetingsArgs = {
  where?: InputMaybe<Filter>;
  sort?: InputMaybe<Array<SortBy>>;
  pagination?: InputMaybe<Pagination>;
  externalSystemId: Scalars['String']['input'];
  externalId?: InputMaybe<Scalars['ID']['input']>;
};

export type QueryGcli_SearchArgs = {
  keyword: Scalars['String']['input'];
  limit?: InputMaybe<Scalars['Int']['input']>;
};

export type QueryInteractionEventArgs = {
  id: Scalars['ID']['input'];
};

export type QueryInteractionEvent_ByEventIdentifierArgs = {
  eventIdentifier: Scalars['String']['input'];
};

export type QueryInteractionSessionArgs = {
  id: Scalars['ID']['input'];
};

export type QueryInteractionSession_ByEventIdentifierArgs = {
  eventIdentifier: Scalars['String']['input'];
};

export type QueryInteractionSession_BySessionIdentifierArgs = {
  sessionIdentifier: Scalars['String']['input'];
};

export type QueryInvoiceArgs = {
  id: Scalars['ID']['input'];
};

export type QueryInvoice_ByNumberArgs = {
  number: Scalars['String']['input'];
};

export type QueryInvoicesArgs = {
  where?: InputMaybe<Filter>;
  sort?: InputMaybe<Array<SortBy>>;
  pagination?: InputMaybe<Pagination>;
  organizationId?: InputMaybe<Scalars['ID']['input']>;
};

export type QueryIssueArgs = {
  id: Scalars['ID']['input'];
};

export type QueryLogEntryArgs = {
  id: Scalars['ID']['input'];
};

export type QueryMasterPlanArgs = {
  id: Scalars['ID']['input'];
};

export type QueryMasterPlansArgs = {
  retired?: InputMaybe<Scalars['Boolean']['input']>;
};

export type QueryMeetingArgs = {
  id: Scalars['ID']['input'];
};

export type QueryOpportunityArgs = {
  id: Scalars['ID']['input'];
};

export type QueryOrganizationArgs = {
  id: Scalars['ID']['input'];
};

export type QueryOrganizationPlanArgs = {
  id: Scalars['ID']['input'];
};

export type QueryOrganizationPlansArgs = {
  retired?: InputMaybe<Scalars['Boolean']['input']>;
};

export type QueryOrganizationPlansForOrganizationArgs = {
  organizationId: Scalars['ID']['input'];
};

export type QueryOrganization_ByCustomIdArgs = {
  customId: Scalars['String']['input'];
};

export type QueryOrganization_ByCustomerOsIdArgs = {
  customerOsId: Scalars['String']['input'];
};

export type QueryOrganizationsArgs = {
  where?: InputMaybe<Filter>;
  sort?: InputMaybe<Array<SortBy>>;
  pagination?: InputMaybe<Pagination>;
};

export type QueryPhoneNumberArgs = {
  id: Scalars['ID']['input'];
};

export type QueryPlayer_ByAuthIdProviderArgs = {
  authId: Scalars['String']['input'];
  provider: Scalars['String']['input'];
};

export type QueryReminderArgs = {
  id: Scalars['ID']['input'];
};

export type QueryRemindersForOrganizationArgs = {
  organizationId: Scalars['ID']['input'];
  dismissed?: InputMaybe<Scalars['Boolean']['input']>;
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

export type QueryTenant_ByEmailArgs = {
  email: Scalars['String']['input'];
};

export type QueryTenant_ByWorkspaceArgs = {
  workspace: WorkspaceInput;
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
  where?: InputMaybe<Filter>;
  sort?: InputMaybe<Array<SortBy>>;
  pagination?: InputMaybe<Pagination>;
};

export type Reminder = MetadataInterface & {
  metadata: Metadata;
  owner?: Maybe<User>;
  __typename?: 'Reminder';
  dueDate?: Maybe<Scalars['Time']['output']>;
  content?: Maybe<Scalars['String']['output']>;
  dismissed?: Maybe<Scalars['Boolean']['output']>;
};

export type ReminderInput = {
  userId: Scalars['ID']['input'];
  dueDate: Scalars['Time']['input'];
  content: Scalars['String']['input'];
  organizationId: Scalars['ID']['input'];
};

export type ReminderUpdateInput = {
  id: Scalars['ID']['input'];
  dueDate?: InputMaybe<Scalars['Time']['input']>;
  content?: InputMaybe<Scalars['String']['input']>;
  dismissed?: InputMaybe<Scalars['Boolean']['input']>;
};

export type RenewalRecord = {
  contract: Contract;
  organization: Organization;
  __typename?: 'RenewalRecord';
  opportunity?: Maybe<Opportunity>;
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
  totalPages: Scalars['Int']['output'];
  totalElements: Scalars['Int64']['output'];
  totalAvailable: Scalars['Int64']['output'];
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
  tax: Tax;
  metadata: Metadata;
  createdBy?: Maybe<User>;
  billingCycle: BilledType;
  __typename?: 'ServiceLineItem';
  parentId: Scalars['ID']['output'];
  price: Scalars['Float']['output'];
  closed: Scalars['Boolean']['output'];
  externalLinks: Array<ExternalSystem>;
  quantity: Scalars['Int64']['output'];
  comments: Scalars['String']['output'];
  description: Scalars['String']['output'];
  serviceStarted: Scalars['Time']['output'];
  serviceEnded?: Maybe<Scalars['Time']['output']>;
};

export type ServiceLineItemBulkUpdateInput = {
  contractId: Scalars['ID']['input'];
  invoiceNote?: InputMaybe<Scalars['String']['input']>;
  serviceLineItems: Array<InputMaybe<ServiceLineItemBulkUpdateItem>>;
};

export type ServiceLineItemBulkUpdateItem = {
  billed?: InputMaybe<BilledType>;
  name?: InputMaybe<Scalars['String']['input']>;
  price?: InputMaybe<Scalars['Float']['input']>;
  vatRate?: InputMaybe<Scalars['Float']['input']>;
  quantity?: InputMaybe<Scalars['Int64']['input']>;
  comments?: InputMaybe<Scalars['String']['input']>;
  newVersion?: InputMaybe<Scalars['Boolean']['input']>;
  serviceStarted?: InputMaybe<Scalars['Time']['input']>;
  closeVersion?: InputMaybe<Scalars['Boolean']['input']>;
  serviceLineItemId?: InputMaybe<Scalars['ID']['input']>;
  isRetroactiveCorrection?: InputMaybe<Scalars['Boolean']['input']>;
};

export type ServiceLineItemCloseInput = {
  id: Scalars['ID']['input'];
  endedAt?: InputMaybe<Scalars['Time']['input']>;
  serviceEnded?: InputMaybe<Scalars['Time']['input']>;
};

export type ServiceLineItemInput = {
  tax?: InputMaybe<TaxInput>;
  contractId: Scalars['ID']['input'];
  billingCycle?: InputMaybe<BilledType>;
  price?: InputMaybe<Scalars['Float']['input']>;
  quantity?: InputMaybe<Scalars['Int64']['input']>;
  appSource?: InputMaybe<Scalars['String']['input']>;
  serviceEnded?: InputMaybe<Scalars['Time']['input']>;
  description?: InputMaybe<Scalars['String']['input']>;
  serviceStarted?: InputMaybe<Scalars['Time']['input']>;
};

export type ServiceLineItemNewVersionInput = {
  tax?: InputMaybe<TaxInput>;
  id?: InputMaybe<Scalars['ID']['input']>;
  price?: InputMaybe<Scalars['Float']['input']>;
  quantity?: InputMaybe<Scalars['Int64']['input']>;
  comments?: InputMaybe<Scalars['String']['input']>;
  appSource?: InputMaybe<Scalars['String']['input']>;
  description?: InputMaybe<Scalars['String']['input']>;
  serviceStarted?: InputMaybe<Scalars['Time']['input']>;
};

export type ServiceLineItemUpdateInput = {
  tax?: InputMaybe<TaxInput>;
  /** Deprecated: billing cycle is not updatable. */
  billingCycle?: InputMaybe<BilledType>;
  id?: InputMaybe<Scalars['ID']['input']>;
  price?: InputMaybe<Scalars['Float']['input']>;
  quantity?: InputMaybe<Scalars['Int64']['input']>;
  comments?: InputMaybe<Scalars['String']['input']>;
  appSource?: InputMaybe<Scalars['String']['input']>;
  serviceEnded?: InputMaybe<Scalars['Time']['input']>;
  description?: InputMaybe<Scalars['String']['input']>;
  serviceStarted?: InputMaybe<Scalars['Time']['input']>;
  isRetroactiveCorrection?: InputMaybe<Scalars['Boolean']['input']>;
};

export type SlackChannel = {
  metadata: Metadata;
  __typename?: 'SlackChannel';
  organization?: Maybe<Organization>;
  channelId: Scalars['String']['output'];
  channelName: Scalars['String']['output'];
};

export type SlackChannelPage = Pages & {
  content: Array<SlackChannel>;
  __typename?: 'SlackChannelPage';
  totalPages: Scalars['Int']['output'];
  totalElements: Scalars['Int64']['output'];
  totalAvailable: Scalars['Int64']['output'];
};

export type Social = Node &
  SourceFields & {
    source: DataSource;
    __typename?: 'Social';
    sourceOfTruth: DataSource;
    id: Scalars['ID']['output'];
    url: Scalars['String']['output'];
    createdAt: Scalars['Time']['output'];
    updatedAt: Scalars['Time']['output'];
    appSource: Scalars['String']['output'];
  };

export type SocialInput = {
  url: Scalars['String']['input'];
  appSource?: InputMaybe<Scalars['String']['input']>;
};

export type SocialUpdateInput = {
  id: Scalars['ID']['input'];
  url: Scalars['String']['input'];
};

export type SortBy = {
  direction?: SortingDirection;
  by: Scalars['String']['input'];
  caseSensitive?: InputMaybe<Scalars['Boolean']['input']>;
};

export enum SortingDirection {
  Asc = 'ASC',
  Desc = 'DESC',
}

export type SourceFields = {
  source: DataSource;
  sourceOfTruth: DataSource;
  id: Scalars['ID']['output'];
  appSource: Scalars['String']['output'];
};

export type SourceFieldsInterface = {
  source: DataSource;
  sourceOfTruth: DataSource;
  appSource: Scalars['String']['output'];
};

export type State = {
  country: Country;
  __typename?: 'State';
  id: Scalars['ID']['output'];
  code: Scalars['String']['output'];
  name: Scalars['String']['output'];
};

export type SuggestedMergeOrganization = {
  organization: Organization;
  __typename?: 'SuggestedMergeOrganization';
  confidence?: Maybe<Scalars['Float']['output']>;
  suggestedAt?: Maybe<Scalars['Time']['output']>;
  suggestedBy?: Maybe<Scalars['String']['output']>;
};

export enum TableIdType {
  AnnualRenewals = 'ANNUAL_RENEWALS',
  Churn = 'CHURN',
  Customers = 'CUSTOMERS',
  Leads = 'LEADS',
  MonthlyRenewals = 'MONTHLY_RENEWALS',
  MyPortfolio = 'MY_PORTFOLIO',
  Nurture = 'NURTURE',
  Organizations = 'ORGANIZATIONS',
  PastInvoices = 'PAST_INVOICES',
  QuarterlyRenewals = 'QUARTERLY_RENEWALS',
  UpcomingInvoices = 'UPCOMING_INVOICES',
}

export type TableViewDef = Node & {
  tableId: TableIdType;
  tableType: TableViewType;
  columns: Array<ColumnView>;
  __typename?: 'TableViewDef';
  id: Scalars['ID']['output'];
  order: Scalars['Int']['output'];
  icon: Scalars['String']['output'];
  name: Scalars['String']['output'];
  createdAt: Scalars['Time']['output'];
  filters: Scalars['String']['output'];
  sorting: Scalars['String']['output'];
  updatedAt: Scalars['Time']['output'];
};

export type TableViewDefCreateInput = {
  tableId: TableIdType;
  tableType: TableViewType;
  order: Scalars['Int']['input'];
  columns: Array<ColumnViewInput>;
  icon: Scalars['String']['input'];
  name: Scalars['String']['input'];
  filters: Scalars['String']['input'];
  sorting: Scalars['String']['input'];
};

export type TableViewDefUpdateInput = {
  id: Scalars['ID']['input'];
  order: Scalars['Int']['input'];
  columns: Array<ColumnViewInput>;
  icon: Scalars['String']['input'];
  name: Scalars['String']['input'];
  filters: Scalars['String']['input'];
  sorting: Scalars['String']['input'];
};

export enum TableViewType {
  Invoices = 'INVOICES',
  Organizations = 'ORGANIZATIONS',
  Renewals = 'RENEWALS',
}

export type Tag = {
  __typename?: 'Tag';
  source: DataSource;
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  createdAt: Scalars['Time']['output'];
  updatedAt: Scalars['Time']['output'];
  appSource: Scalars['String']['output'];
};

export type TagIdOrNameInput = {
  id?: InputMaybe<Scalars['ID']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
};

export type TagInput = {
  name: Scalars['String']['input'];
  appSource?: InputMaybe<Scalars['String']['input']>;
};

export type TagUpdateInput = {
  id: Scalars['ID']['input'];
  name: Scalars['String']['input'];
};

export type Tax = {
  __typename?: 'Tax';
  vat: Scalars['Boolean']['output'];
  taxRate: Scalars['Float']['output'];
  salesTax: Scalars['Boolean']['output'];
};

export type TaxInput = {
  taxRate: Scalars['Float']['input'];
};

export type TenantBillableInfo = {
  __typename?: 'TenantBillableInfo';
  greylistedContacts: Scalars['Int64']['output'];
  whitelistedContacts: Scalars['Int64']['output'];
  greylistedOrganizations: Scalars['Int64']['output'];
  whitelistedOrganizations: Scalars['Int64']['output'];
};

export type TenantBillingProfile = Node &
  SourceFields & {
    source: DataSource;
    sourceOfTruth: DataSource;
    id: Scalars['ID']['output'];
    zip: Scalars['String']['output'];
    /**
     * Deprecated
     * @deprecated Use sendInvoicesFrom
     */
    email: Scalars['String']['output'];
    phone: Scalars['String']['output'];
    __typename?: 'TenantBillingProfile';
    check: Scalars['Boolean']['output'];
    region: Scalars['String']['output'];
    country: Scalars['String']['output'];
    createdAt: Scalars['Time']['output'];
    updatedAt: Scalars['Time']['output'];
    locality: Scalars['String']['output'];
    appSource: Scalars['String']['output'];
    legalName: Scalars['String']['output'];
    vatNumber: Scalars['String']['output'];
    addressLine1: Scalars['String']['output'];
    addressLine2: Scalars['String']['output'];
    addressLine3: Scalars['String']['output'];
    sendInvoicesBcc: Scalars['String']['output'];
    sendInvoicesFrom: Scalars['String']['output'];
    canPayWithPigeon: Scalars['Boolean']['output'];
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
    domesticPaymentsBankInfo?: Maybe<Scalars['String']['output']>;
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
    /**
     * Deprecated
     * @deprecated Not used
     */
    internationalPaymentsBankInfo?: Maybe<Scalars['String']['output']>;
  };

export type TenantBillingProfileInput = {
  check: Scalars['Boolean']['input'];
  vatNumber: Scalars['String']['input'];
  sendInvoicesFrom: Scalars['String']['input'];
  zip?: InputMaybe<Scalars['String']['input']>;
  canPayWithPigeon: Scalars['Boolean']['input'];
  /** Deprecated */
  email?: InputMaybe<Scalars['String']['input']>;
  phone?: InputMaybe<Scalars['String']['input']>;
  region?: InputMaybe<Scalars['String']['input']>;
  country?: InputMaybe<Scalars['String']['input']>;
  locality?: InputMaybe<Scalars['String']['input']>;
  legalName?: InputMaybe<Scalars['String']['input']>;
  canPayWithBankTransfer: Scalars['Boolean']['input'];
  addressLine1?: InputMaybe<Scalars['String']['input']>;
  addressLine2?: InputMaybe<Scalars['String']['input']>;
  addressLine3?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  canPayWithCard?: InputMaybe<Scalars['Boolean']['input']>;
  sendInvoicesBcc?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  domesticPaymentsBankInfo?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  canPayWithDirectDebitACH?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  canPayWithDirectDebitBacs?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  canPayWithDirectDebitSEPA?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  internationalPaymentsBankInfo?: InputMaybe<Scalars['String']['input']>;
};

export type TenantBillingProfileUpdateInput = {
  id: Scalars['ID']['input'];
  zip?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  email?: InputMaybe<Scalars['String']['input']>;
  phone?: InputMaybe<Scalars['String']['input']>;
  check?: InputMaybe<Scalars['Boolean']['input']>;
  patch?: InputMaybe<Scalars['Boolean']['input']>;
  region?: InputMaybe<Scalars['String']['input']>;
  country?: InputMaybe<Scalars['String']['input']>;
  locality?: InputMaybe<Scalars['String']['input']>;
  legalName?: InputMaybe<Scalars['String']['input']>;
  vatNumber?: InputMaybe<Scalars['String']['input']>;
  addressLine1?: InputMaybe<Scalars['String']['input']>;
  addressLine2?: InputMaybe<Scalars['String']['input']>;
  addressLine3?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  canPayWithCard?: InputMaybe<Scalars['Boolean']['input']>;
  sendInvoicesBcc?: InputMaybe<Scalars['String']['input']>;
  sendInvoicesFrom?: InputMaybe<Scalars['String']['input']>;
  canPayWithPigeon?: InputMaybe<Scalars['Boolean']['input']>;
  canPayWithBankTransfer?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  domesticPaymentsBankInfo?: InputMaybe<Scalars['String']['input']>;
  /** Deprecated */
  canPayWithDirectDebitACH?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  canPayWithDirectDebitBacs?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  canPayWithDirectDebitSEPA?: InputMaybe<Scalars['Boolean']['input']>;
  /** Deprecated */
  internationalPaymentsBankInfo?: InputMaybe<Scalars['String']['input']>;
};

export type TenantInput = {
  name: Scalars['String']['input'];
  appSource?: InputMaybe<Scalars['String']['input']>;
};

export type TenantSettings = {
  __typename?: 'TenantSettings';
  baseCurrency?: Maybe<Currency>;
  /**
   * Deprecated
   * @deprecated Use logoRepositoryFileId
   */
  logoUrl: Scalars['String']['output'];
  billingEnabled: Scalars['Boolean']['output'];
  logoRepositoryFileId?: Maybe<Scalars['String']['output']>;
};

export type TenantSettingsInput = {
  baseCurrency?: InputMaybe<Currency>;
  patch?: InputMaybe<Scalars['Boolean']['input']>;
  logoUrl?: InputMaybe<Scalars['String']['input']>;
  billingEnabled?: InputMaybe<Scalars['Boolean']['input']>;
  logoRepositoryFileId?: InputMaybe<Scalars['String']['input']>;
};

export type TimeRange = {
  /**
   * The end time of the time range.
   * **Required.**
   */
  to: Scalars['Time']['input'];
  /**
   * The start time of the time range.
   * **Required.**
   */
  from: Scalars['Time']['input'];
};

export type TimelineEvent =
  | Action
  | Analysis
  | InteractionEvent
  | InteractionSession
  | Issue
  | LogEntry
  | Meeting
  | Note
  | Order
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
  player: Player;
  roles: Array<Role>;
  source: DataSource;
  __typename?: 'User';
  jobRoles: Array<JobRole>;
  sourceOfTruth: DataSource;
  calendars: Array<Calendar>;
  /**
   * The unique ID associated with the customerOS user.
   * **Required**
   */
  id: Scalars['ID']['output'];
  /**
   * All email addresses associated with a user in customerOS.
   * **Required.  If no values it returns an empty array.**
   */
  emails?: Maybe<Array<Email>>;
  phoneNumbers: Array<PhoneNumber>;
  bot: Scalars['Boolean']['output'];
  /**
   * Timestamp of user creation.
   * **Required**
   */
  createdAt: Scalars['Time']['output'];
  updatedAt: Scalars['Time']['output'];
  /**
   * The last name of the customerOS user.
   * **Required**
   */
  lastName: Scalars['String']['output'];
  appSource: Scalars['String']['output'];
  /**
   * The first name of the customerOS user.
   * **Required**
   */
  firstName: Scalars['String']['output'];
  internal: Scalars['Boolean']['output'];
  name?: Maybe<Scalars['String']['output']>;
  timezone?: Maybe<Scalars['String']['output']>;
  profilePhotoUrl?: Maybe<Scalars['String']['output']>;
};

/**
 * Describes the User of customerOS.  A user is the person who logs into the Openline platform.
 * **A `create` object.**
 */
export type UserInput = {
  /**
   * The email address of the customerOS user.
   * **Required**
   */
  email: EmailInput;
  /**
   * Player to associate with the user with. If the person does not exist, it will be created.
   * **Required**
   */
  player: PlayerInput;
  /**
   * The last name of the customerOS user.
   * **Required**
   */
  lastName: Scalars['String']['input'];
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
  name?: InputMaybe<Scalars['String']['input']>;
  timezone?: InputMaybe<Scalars['String']['input']>;
  /**
   * The name of the app performing the create.
   * **Optional**
   */
  appSource?: InputMaybe<Scalars['String']['input']>;
  profilePhotoUrl?: InputMaybe<Scalars['String']['input']>;
};

/**
 * Specifies how many pages of `User` information has been returned in the query response.
 * **A `return` object.**
 */
export type UserPage = Pages & {
  /**
   * A `User` entity in customerOS.
   * **Required.  If no values it returns an empty array.**
   */
  content: Array<User>;
  __typename?: 'UserPage';
  /**
   * Total number of pages in the query response.
   * **Required.**
   */
  totalPages: Scalars['Int']['output'];
  /**
   * Total number of elements in the query response.
   * **Required.**
   */
  totalElements: Scalars['Int64']['output'];
};

export type UserParticipant = {
  userParticipant: User;
  __typename?: 'UserParticipant';
  type?: Maybe<Scalars['String']['output']>;
};

export type UserUpdateInput = {
  id: Scalars['ID']['input'];
  /**
   * The last name of the customerOS user.
   * **Required**
   */
  lastName: Scalars['String']['input'];
  /**
   * The first name of the customerOS user.
   * **Required**
   */
  firstName: Scalars['String']['input'];
  name?: InputMaybe<Scalars['String']['input']>;
  timezone?: InputMaybe<Scalars['String']['input']>;
  profilePhotoUrl?: InputMaybe<Scalars['String']['input']>;
};

export type Workspace = {
  source: DataSource;
  __typename?: 'Workspace';
  sourceOfTruth: DataSource;
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  createdAt: Scalars['Time']['output'];
  updatedAt: Scalars['Time']['output'];
  provider: Scalars['String']['output'];
  appSource: Scalars['String']['output'];
};

export type WorkspaceInput = {
  name: Scalars['String']['input'];
  provider: Scalars['String']['input'];
  appSource?: InputMaybe<Scalars['String']['input']>;
};
