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
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  Any: any;
  Time: any;
  ID: string;
  Int64: any;
  Int: number;
  Float: number;
  String: string;
  Boolean: boolean;
};

export type Action = {
  id: Scalars['ID'];
  source: DataSource;
  __typename?: 'Action';
  actionType: ActionType;
  createdBy?: Maybe<User>;
  createdAt: Scalars['Time'];
  appSource: Scalars['String'];
  content?: Maybe<Scalars['String']>;
  metadata?: Maybe<Scalars['String']>;
};

export type ActionItem = {
  id: Scalars['ID'];
  source: DataSource;
  __typename?: 'ActionItem';
  content: Scalars['String'];
  createdAt: Scalars['Time'];
  appSource: Scalars['String'];
};

export enum ActionType {
  Created = 'CREATED',
  RenewalForecastUpdated = 'RENEWAL_FORECAST_UPDATED',
  RenewalLikelihoodUpdated = 'RENEWAL_LIKELIHOOD_UPDATED',
}

export type Analysis = Node & {
  id: Scalars['ID'];
  source: DataSource;
  __typename?: 'Analysis';
  sourceOfTruth: DataSource;
  createdAt: Scalars['Time'];
  appSource: Scalars['String'];
  describes: Array<DescriptionNode>;
  content?: Maybe<Scalars['String']>;
  contentType?: Maybe<Scalars['String']>;
  analysisType?: Maybe<Scalars['String']>;
};

export type AnalysisDescriptionInput = {
  meetingId?: InputMaybe<Scalars['ID']>;
  interactionEventId?: InputMaybe<Scalars['ID']>;
  interactionSessionId?: InputMaybe<Scalars['ID']>;
};

export type AnalysisInput = {
  appSource: Scalars['String'];
  content?: InputMaybe<Scalars['String']>;
  describes: Array<AnalysisDescriptionInput>;
  contentType?: InputMaybe<Scalars['String']>;
  analysisType?: InputMaybe<Scalars['String']>;
};

export type Attachment = Node & {
  id: Scalars['ID'];
  source: DataSource;
  size: Scalars['Int64'];
  name: Scalars['String'];
  __typename?: 'Attachment';
  sourceOfTruth: DataSource;
  createdAt: Scalars['Time'];
  mimeType: Scalars['String'];
  appSource: Scalars['String'];
  extension: Scalars['String'];
};

export type AttachmentInput = {
  size: Scalars['Int64'];
  name: Scalars['String'];
  mimeType: Scalars['String'];
  appSource: Scalars['String'];
  extension: Scalars['String'];
};

export type BillingDetails = {
  __typename?: 'BillingDetails';
  frequency?: Maybe<RenewalCycle>;
  amount?: Maybe<Scalars['Float']>;
  renewalCycle?: Maybe<RenewalCycle>;
  renewalCycleNext?: Maybe<Scalars['Time']>;
  renewalCycleStart?: Maybe<Scalars['Time']>;
};

export type BillingDetailsInput = {
  id: Scalars['ID'];
  frequency?: InputMaybe<RenewalCycle>;
  amount?: InputMaybe<Scalars['Float']>;
  renewalCycle?: InputMaybe<RenewalCycle>;
  renewalCycleStart?: InputMaybe<Scalars['Time']>;
};

/**
 * Describes the relationship a Contact has with a Organization.
 * **A `return` object**
 */
export type Calendar = {
  id: Scalars['ID'];
  source: DataSource;
  calType: CalendarType;
  __typename?: 'Calendar';
  sourceOfTruth: DataSource;
  createdAt: Scalars['Time'];
  updatedAt: Scalars['Time'];
  primary: Scalars['Boolean'];
  appSource: Scalars['String'];
  link?: Maybe<Scalars['String']>;
};

export enum CalendarType {
  Calcom = 'CALCOM',
  Google = 'GOOGLE',
}

export type Comment = {
  id: Scalars['ID'];
  source: DataSource;
  __typename?: 'Comment';
  createdBy?: Maybe<User>;
  sourceOfTruth: DataSource;
  createdAt: Scalars['Time'];
  updatedAt: Scalars['Time'];
  appSource: Scalars['String'];
  content?: Maybe<Scalars['String']>;
  externalLinks: Array<ExternalSystem>;
  contentType?: Maybe<Scalars['String']>;
};

export enum ComparisonOperator {
  Contains = 'CONTAINS',
  Eq = 'EQ',
  StartsWith = 'STARTS_WITH',
}

/**
 * A contact represents an individual in customerOS.
 * **A `response` object.**
 */
export type Contact = ExtensibleEntity &
  Node & {
    /** Contact notes */
    notes: NotePage;
    /**
     * The unique ID associated with the contact in customerOS.
     * **Required**
     */
    id: Scalars['ID'];
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
    /**
     * An ISO8601 timestamp recording when the contact was created in customerOS.
     * **Required**
     */
    createdAt: Scalars['Time'];
    fieldSets: Array<FieldSet>;
    /**
     * All locations associated with a contact in customerOS.
     * **Required.  If no values it returns an empty array.**
     */
    locations: Array<Location>;
    updatedAt: Scalars['Time'];
    /** The name of the contact in customerOS, alternative for firstName + lastName. */
    name?: Maybe<Scalars['String']>;
    organizations: OrganizationPage;
    /**
     * User defined metadata appended to the contact record in customerOS.
     * **Required.  If no values it returns an empty array.**
     */
    customFields: Array<CustomField>;
    /** @deprecated Use `tags` instead */
    label?: Maybe<Scalars['String']>;
    /**
     * All phone numbers associated with a contact in customerOS.
     * **Required.  If no values it returns an empty array.**
     */
    phoneNumbers: Array<PhoneNumber>;
    /** Template of the contact in customerOS. */
    template?: Maybe<EntityTemplate>;
    /**
     * The title associate with the contact in customerOS.
     * @deprecated Use `prefix` instead
     */
    title?: Maybe<Scalars['String']>;
    prefix?: Maybe<Scalars['String']>;
    /** The last name of the contact in customerOS. */
    lastName?: Maybe<Scalars['String']>;
    timezone?: Maybe<Scalars['String']>;
    appSource?: Maybe<Scalars['String']>;
    /** The first name of the contact in customerOS. */
    firstName?: Maybe<Scalars['String']>;
    timelineEvents: Array<TimelineEvent>;
    description?: Maybe<Scalars['String']>;
    profilePhotoUrl?: Maybe<Scalars['String']>;
    timelineEventsTotalCount: Scalars['Int64'];
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
  size: Scalars['Int'];
  from?: InputMaybe<Scalars['Time']>;
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
  /** Id of the contact owner (user) */
  ownerId?: InputMaybe<Scalars['ID']>;
  name?: InputMaybe<Scalars['String']>;
  /** The prefix of the contact. */
  prefix?: InputMaybe<Scalars['String']>;
  /** The unique ID associated with the template of the contact in customerOS. */
  templateId?: InputMaybe<Scalars['ID']>;
  /** An ISO8601 timestamp recording when the contact was created in customerOS. */
  createdAt?: InputMaybe<Scalars['Time']>;
  /** The last name of the contact. */
  lastName?: InputMaybe<Scalars['String']>;
  timezone?: InputMaybe<Scalars['String']>;
  appSource?: InputMaybe<Scalars['String']>;
  /** The first name of the contact. */
  firstName?: InputMaybe<Scalars['String']>;
  /** A phone number associated with the contact. */
  phoneNumber?: InputMaybe<PhoneNumberInput>;
  description?: InputMaybe<Scalars['String']>;
  fieldSets?: InputMaybe<Array<FieldSetInput>>;
  /**
   * User defined metadata appended to contact.
   * **Required.**
   */
  customFields?: InputMaybe<Array<CustomFieldInput>>;
  externalReference?: InputMaybe<ExternalSystemReferenceInput>;
};

export type ContactOrganizationInput = {
  contactId: Scalars['ID'];
  organizationId: Scalars['ID'];
};

export type ContactParticipant = {
  contactParticipant: Contact;
  type?: Maybe<Scalars['String']>;
  __typename?: 'ContactParticipant';
};

export type ContactTagInput = {
  tagId: Scalars['ID'];
  contactId: Scalars['ID'];
};

/**
 * Updates data fields associated with an existing customer record in customerOS.
 * **An `update` object.**
 */
export type ContactUpdateInput = {
  /**
   * The unique ID associated with the contact in customerOS.
   * **Required.**
   */
  id: Scalars['ID'];
  /** Id of the contact owner (user) */
  ownerId?: InputMaybe<Scalars['ID']>;
  name?: InputMaybe<Scalars['String']>;
  label?: InputMaybe<Scalars['String']>;
  /** The prefix associate with the contact in customerOS. */
  prefix?: InputMaybe<Scalars['String']>;
  /** The last name of the contact in customerOS. */
  lastName?: InputMaybe<Scalars['String']>;
  timezone?: InputMaybe<Scalars['String']>;
  /** The first name of the contact in customerOS. */
  firstName?: InputMaybe<Scalars['String']>;
  description?: InputMaybe<Scalars['String']>;
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
  /**
   * Total number of pages in the query response.
   * **Required.**
   */
  totalPages: Scalars['Int'];
  __typename?: 'ContactsPage';
  /**
   * Total number of elements in the query response.
   * **Required.**
   */
  totalElements: Scalars['Int64'];
};

export type Country = {
  id: Scalars['ID'];
  __typename?: 'Country';
  name: Scalars['String'];
  codeA2: Scalars['String'];
  codeA3: Scalars['String'];
  phoneCode: Scalars['String'];
};

/**
 * Describes a custom, user-defined field associated with a `Contact`.
 * **A `return` object.**
 */
export type CustomField = Node & {
  /**
   * The unique ID associated with the custom field.
   * **Required**
   */
  id: Scalars['ID'];
  /** The source of the custom field value */
  source: DataSource;
  /**
   * The value of the custom field.
   * **Required**
   */
  value: Scalars['Any'];
  /**
   * The name of the custom field.
   * **Required**
   */
  name: Scalars['String'];
  __typename?: 'CustomField';
  createdAt: Scalars['Time'];
  updatedAt: Scalars['Time'];
  /**
   * Datatype of the custom field.
   * **Required**
   */
  datatype: CustomFieldDataType;
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
  id: Scalars['ID'];
  entityType: EntityType;
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
  value: Scalars['Any'];
  id?: InputMaybe<Scalars['ID']>;
  /** The name of the custom field. */
  name?: InputMaybe<Scalars['String']>;
  templateId?: InputMaybe<Scalars['ID']>;
  /** Datatype of the custom field. */
  datatype?: InputMaybe<CustomFieldDataType>;
};

export type CustomFieldTemplate = Node & {
  id: Scalars['ID'];
  order: Scalars['Int'];
  name: Scalars['String'];
  createdAt: Scalars['Time'];
  updatedAt: Scalars['Time'];
  max?: Maybe<Scalars['Int']>;
  min?: Maybe<Scalars['Int']>;
  mandatory: Scalars['Boolean'];
  type: CustomFieldTemplateType;
  length?: Maybe<Scalars['Int']>;
  __typename?: 'CustomFieldTemplate';
};

export type CustomFieldTemplateInput = {
  order: Scalars['Int'];
  name: Scalars['String'];
  type: CustomFieldTemplateType;
  max?: InputMaybe<Scalars['Int']>;
  min?: InputMaybe<Scalars['Int']>;
  length?: InputMaybe<Scalars['Int']>;
  mandatory?: InputMaybe<Scalars['Boolean']>;
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
  id: Scalars['ID'];
  /**
   * The value of the custom field.
   * **Required**
   */
  value: Scalars['Any'];
  /**
   * The name of the custom field.
   * **Required**
   */
  name: Scalars['String'];
  /**
   * Datatype of the custom field.
   * **Required**
   */
  datatype: CustomFieldDataType;
};

export type CustomerContact = {
  id: Scalars['ID'];
  email: CustomerEmail;
  __typename?: 'CustomerContact';
};

export type CustomerContactInput = {
  /** An email addresses associted with the contact. */
  email?: InputMaybe<EmailInput>;
  name?: InputMaybe<Scalars['String']>;
  /** The prefix of the contact. */
  prefix?: InputMaybe<Scalars['String']>;
  /** An ISO8601 timestamp recording when the contact was created in customerOS. */
  createdAt?: InputMaybe<Scalars['Time']>;
  /** The last name of the contact. */
  lastName?: InputMaybe<Scalars['String']>;
  timezone?: InputMaybe<Scalars['String']>;
  appSource?: InputMaybe<Scalars['String']>;
  /** The first name of the contact. */
  firstName?: InputMaybe<Scalars['String']>;
  description?: InputMaybe<Scalars['String']>;
};

export type CustomerEmail = {
  id: Scalars['ID'];
  __typename?: 'CustomerEmail';
};

export type CustomerJobRole = {
  id: Scalars['ID'];
  __typename?: 'CustomerJobRole';
};

export type CustomerUser = {
  id: Scalars['ID'];
  jobRole: CustomerJobRole;
  __typename?: 'CustomerUser';
};

export enum DataSource {
  Hubspot = 'HUBSPOT',
  Intercom = 'INTERCOM',
  Na = 'NA',
  Openline = 'OPENLINE',
  Pipedrive = 'PIPEDRIVE',
  Salesforce = 'SALESFORCE',
  Slack = 'SLACK',
  Webscrape = 'WEBSCRAPE',
  ZendeskSupport = 'ZENDESK_SUPPORT',
}

export type DescriptionNode = InteractionEvent | InteractionSession | Meeting;

/**
 * Describes an email address associated with a `Contact` in customerOS.
 * **A `return` object.**
 */
export type Email = {
  /**
   * The unique ID associated with the contact in customerOS.
   * **Required**
   */
  id: Scalars['ID'];
  source: DataSource;
  users: Array<User>;
  __typename?: 'Email';
  contacts: Array<Contact>;
  /** Describes the type of email address (WORK, PERSONAL, etc). */
  label?: Maybe<EmailLabel>;
  sourceOfTruth: DataSource;
  createdAt: Scalars['Time'];
  updatedAt: Scalars['Time'];
  /**
   * Identifies whether the email address is primary or not.
   * **Required.**
   */
  primary: Scalars['Boolean'];
  appSource: Scalars['String'];
  /** An email address assocaited with the contact in customerOS. */
  email?: Maybe<Scalars['String']>;
  organizations: Array<Organization>;
  rawEmail?: Maybe<Scalars['String']>;
  emailValidationDetails: EmailValidationDetails;
};

/**
 * Describes an email address associated with a `Contact` in customerOS.
 * **A `create` object.**
 */
export type EmailInput = {
  /**
   * An email address associated with the contact in customerOS.
   * **Required.**
   */
  email: Scalars['String'];
  /** Describes the type of email address (WORK, PERSONAL, etc). */
  label?: InputMaybe<EmailLabel>;
  /**
   * Identifies whether the email address is primary or not.
   * **Required.**
   */
  primary?: InputMaybe<Scalars['Boolean']>;
  appSource?: InputMaybe<Scalars['String']>;
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
  type?: Maybe<Scalars['String']>;
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
  id: Scalars['ID'];
  /** Describes the type of email address (WORK, PERSONAL, etc). */
  label?: InputMaybe<EmailLabel>;
  email?: InputMaybe<Scalars['String']>;
  /**
   * Identifies whether the email address is primary or not.
   * **Required.**
   */
  primary?: InputMaybe<Scalars['Boolean']>;
};

export type EmailValidationDetails = {
  error?: Maybe<Scalars['String']>;
  __typename?: 'EmailValidationDetails';
  validated?: Maybe<Scalars['Boolean']>;
  isCatchAll?: Maybe<Scalars['Boolean']>;
  isDisabled?: Maybe<Scalars['Boolean']>;
  isReachable?: Maybe<Scalars['String']>;
  acceptsMail?: Maybe<Scalars['Boolean']>;
  hasFullInbox?: Maybe<Scalars['Boolean']>;
  isDeliverable?: Maybe<Scalars['Boolean']>;
  isValidSyntax?: Maybe<Scalars['Boolean']>;
  canConnectSmtp?: Maybe<Scalars['Boolean']>;
};

export type EntityTemplate = Node & {
  id: Scalars['ID'];
  name: Scalars['String'];
  version: Scalars['Int'];
  createdAt: Scalars['Time'];
  updatedAt: Scalars['Time'];
  __typename?: 'EntityTemplate';
  extends?: Maybe<EntityTemplateExtension>;
  fieldSetTemplates: Array<FieldSetTemplate>;
  customFieldTemplates: Array<CustomFieldTemplate>;
};

export enum EntityTemplateExtension {
  Contact = 'CONTACT',
  Organization = 'ORGANIZATION',
}

export type EntityTemplateInput = {
  name: Scalars['String'];
  extends?: InputMaybe<EntityTemplateExtension>;
  fieldSetTemplateInputs?: InputMaybe<Array<FieldSetTemplateInput>>;
  customFieldTemplateInputs?: InputMaybe<Array<CustomFieldTemplateInput>>;
};

export enum EntityType {
  Contact = 'Contact',
  Organization = 'Organization',
}

export type ExtensibleEntity = {
  id: Scalars['ID'];
  template?: Maybe<EntityTemplate>;
};

export type ExternalSystem = {
  type: ExternalSystemType;
  __typename?: 'ExternalSystem';
  syncDate?: Maybe<Scalars['Time']>;
  externalId?: Maybe<Scalars['String']>;
  externalUrl?: Maybe<Scalars['String']>;
  externalSource?: Maybe<Scalars['String']>;
};

export type ExternalSystemReferenceInput = {
  type: ExternalSystemType;
  externalId: Scalars['ID'];
  syncDate?: InputMaybe<Scalars['Time']>;
  externalUrl?: InputMaybe<Scalars['String']>;
  externalSource?: InputMaybe<Scalars['String']>;
};

export enum ExternalSystemType {
  Calcom = 'CALCOM',
  Hubspot = 'HUBSPOT',
  Intercom = 'INTERCOM',
  Pipedrive = 'PIPEDRIVE',
  Salesforce = 'SALESFORCE',
  Slack = 'SLACK',
  ZendeskSupport = 'ZENDESK_SUPPORT',
}

export type FieldSet = {
  id: Scalars['ID'];
  source: DataSource;
  __typename?: 'FieldSet';
  name: Scalars['String'];
  createdAt: Scalars['Time'];
  updatedAt: Scalars['Time'];
  customFields: Array<CustomField>;
  template?: Maybe<FieldSetTemplate>;
};

export type FieldSetInput = {
  name: Scalars['String'];
  id?: InputMaybe<Scalars['ID']>;
  templateId?: InputMaybe<Scalars['ID']>;
  customFields?: InputMaybe<Array<CustomFieldInput>>;
};

export type FieldSetTemplate = Node & {
  id: Scalars['ID'];
  order: Scalars['Int'];
  name: Scalars['String'];
  createdAt: Scalars['Time'];
  updatedAt: Scalars['Time'];
  __typename?: 'FieldSetTemplate';
  customFieldTemplates: Array<CustomFieldTemplate>;
};

export type FieldSetTemplateInput = {
  order: Scalars['Int'];
  name: Scalars['String'];
  customFieldTemplateInputs?: InputMaybe<Array<CustomFieldTemplateInput>>;
};

export type FieldSetUpdateInput = {
  id: Scalars['ID'];
  name: Scalars['String'];
};

export type Filter = {
  NOT?: InputMaybe<Filter>;
  OR?: InputMaybe<Array<Filter>>;
  AND?: InputMaybe<Array<Filter>>;
  filter?: InputMaybe<FilterItem>;
};

export type FilterItem = {
  value: Scalars['Any'];
  property: Scalars['String'];
  operation?: ComparisonOperator;
  caseSensitive?: InputMaybe<Scalars['Boolean']>;
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
  key: Scalars['String'];
  value: Scalars['String'];
  display?: Maybe<Scalars['String']>;
  __typename?: 'GCliAttributeKeyValuePair';
};

export enum GCliCacheItemType {
  Contact = 'CONTACT',
  Organization = 'ORGANIZATION',
  State = 'STATE',
}

export type GCliItem = {
  id: Scalars['ID'];
  __typename?: 'GCliItem';
  display: Scalars['String'];
  type: GCliSearchResultType;
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
  isOwner: Scalars['Boolean'];
  isGoogleActive: Scalars['Boolean'];
  isGoogleTokenExpired: Scalars['Boolean'];
};

export type InteractionEvent = Node & {
  id: Scalars['ID'];
  source: DataSource;
  issue?: Maybe<Issue>;
  meeting?: Maybe<Meeting>;
  sourceOfTruth: DataSource;
  summary?: Maybe<Analysis>;
  createdAt: Scalars['Time'];
  includes: Array<Attachment>;
  appSource: Scalars['String'];
  __typename?: 'InteractionEvent';
  channel?: Maybe<Scalars['String']>;
  content?: Maybe<Scalars['String']>;
  repliesTo?: Maybe<InteractionEvent>;
  eventType?: Maybe<Scalars['String']>;
  externalLinks: Array<ExternalSystem>;
  actionItems?: Maybe<Array<ActionItem>>;
  channelData?: Maybe<Scalars['String']>;
  contentType?: Maybe<Scalars['String']>;
  eventIdentifier?: Maybe<Scalars['String']>;
  sentBy: Array<InteractionEventParticipant>;
  sentTo: Array<InteractionEventParticipant>;
  interactionSession?: Maybe<InteractionSession>;
};

export type InteractionEventInput = {
  appSource: Scalars['String'];
  meetingId?: InputMaybe<Scalars['ID']>;
  repliesTo?: InputMaybe<Scalars['ID']>;
  channel?: InputMaybe<Scalars['String']>;
  content?: InputMaybe<Scalars['String']>;
  createdAt?: InputMaybe<Scalars['Time']>;
  eventType?: InputMaybe<Scalars['String']>;
  externalId?: InputMaybe<Scalars['String']>;
  channelData?: InputMaybe<Scalars['String']>;
  contentType?: InputMaybe<Scalars['String']>;
  interactionSession?: InputMaybe<Scalars['ID']>;
  eventIdentifier?: InputMaybe<Scalars['String']>;
  sentBy: Array<InteractionEventParticipantInput>;
  sentTo: Array<InteractionEventParticipantInput>;
  externalSystemId?: InputMaybe<Scalars['String']>;
};

export type InteractionEventParticipant =
  | ContactParticipant
  | EmailParticipant
  | JobRoleParticipant
  | OrganizationParticipant
  | PhoneNumberParticipant
  | UserParticipant;

export type InteractionEventParticipantInput = {
  userID?: InputMaybe<Scalars['ID']>;
  type?: InputMaybe<Scalars['String']>;
  contactID?: InputMaybe<Scalars['ID']>;
  email?: InputMaybe<Scalars['String']>;
  phoneNumber?: InputMaybe<Scalars['String']>;
};

export type InteractionSession = Node & {
  id: Scalars['ID'];
  source: DataSource;
  name: Scalars['String'];
  sourceOfTruth: DataSource;
  status: Scalars['String'];
  createdAt: Scalars['Time'];
  /** @deprecated Use createdAt instead */
  startedAt: Scalars['Time'];
  updatedAt: Scalars['Time'];
  includes: Array<Attachment>;
  appSource: Scalars['String'];
  describedBy: Array<Analysis>;
  events: Array<InteractionEvent>;
  type?: Maybe<Scalars['String']>;
  /** @deprecated Use updatedAt instead */
  endedAt?: Maybe<Scalars['Time']>;
  __typename?: 'InteractionSession';
  channel?: Maybe<Scalars['String']>;
  channelData?: Maybe<Scalars['String']>;
  sessionIdentifier?: Maybe<Scalars['String']>;
  attendedBy: Array<InteractionSessionParticipant>;
};

export type InteractionSessionInput = {
  name: Scalars['String'];
  status: Scalars['String'];
  appSource: Scalars['String'];
  type?: InputMaybe<Scalars['String']>;
  channel?: InputMaybe<Scalars['String']>;
  channelData?: InputMaybe<Scalars['String']>;
  sessionIdentifier?: InputMaybe<Scalars['String']>;
  attendedBy?: InputMaybe<Array<InteractionSessionParticipantInput>>;
};

export type InteractionSessionParticipant =
  | ContactParticipant
  | EmailParticipant
  | PhoneNumberParticipant
  | UserParticipant;

export type InteractionSessionParticipantInput = {
  userID?: InputMaybe<Scalars['ID']>;
  type?: InputMaybe<Scalars['String']>;
  contactID?: InputMaybe<Scalars['ID']>;
  email?: InputMaybe<Scalars['String']>;
  phoneNumber?: InputMaybe<Scalars['String']>;
};

export type Issue = Node &
  SourceFields & {
    id: Scalars['ID'];
    source: DataSource;
    __typename?: 'Issue';
    comments: Array<Comment>;
    sourceOfTruth: DataSource;
    status: Scalars['String'];
    createdAt: Scalars['Time'];
    updatedAt: Scalars['Time'];
    appSource: Scalars['String'];
    tags?: Maybe<Array<Maybe<Tag>>>;
    subject?: Maybe<Scalars['String']>;
    assignedTo: Array<IssueParticipant>;
    followedBy: Array<IssueParticipant>;
    priority?: Maybe<Scalars['String']>;
    externalLinks: Array<ExternalSystem>;
    reportedBy?: Maybe<IssueParticipant>;
    submittedBy?: Maybe<IssueParticipant>;
    description?: Maybe<Scalars['String']>;
    interactionEvents: Array<InteractionEvent>;
  };

export type IssueParticipant =
  | ContactParticipant
  | OrganizationParticipant
  | UserParticipant;

export type IssueSummaryByStatus = {
  count: Scalars['Int64'];
  status: Scalars['String'];
  __typename?: 'IssueSummaryByStatus';
};

/**
 * Describes the relationship a Contact has with a Organization.
 * **A `return` object**
 */
export type JobRole = {
  id: Scalars['ID'];
  source: DataSource;
  __typename?: 'JobRole';
  contact?: Maybe<Contact>;
  sourceOfTruth: DataSource;
  createdAt: Scalars['Time'];
  updatedAt: Scalars['Time'];
  primary: Scalars['Boolean'];
  appSource: Scalars['String'];
  endedAt?: Maybe<Scalars['Time']>;
  company?: Maybe<Scalars['String']>;
  /**
   * Organization associated with a Contact.
   * **Required.**
   */
  organization?: Maybe<Organization>;
  startedAt?: Maybe<Scalars['Time']>;
  /** The Contact's job title. */
  jobTitle?: Maybe<Scalars['String']>;
  description?: Maybe<Scalars['String']>;
};

/**
 * Describes the relationship a Contact has with an Organization.
 * **A `create` object**
 */
export type JobRoleInput = {
  endedAt?: InputMaybe<Scalars['Time']>;
  company?: InputMaybe<Scalars['String']>;
  startedAt?: InputMaybe<Scalars['Time']>;
  jobTitle?: InputMaybe<Scalars['String']>;
  primary?: InputMaybe<Scalars['Boolean']>;
  appSource?: InputMaybe<Scalars['String']>;
  organizationId?: InputMaybe<Scalars['ID']>;
  description?: InputMaybe<Scalars['String']>;
};

export type JobRoleParticipant = {
  jobRoleParticipant: JobRole;
  type?: Maybe<Scalars['String']>;
  __typename?: 'JobRoleParticipant';
};

/**
 * Describes the relationship a Contact has with an Organization.
 * **A `create` object**
 */
export type JobRoleUpdateInput = {
  id: Scalars['ID'];
  endedAt?: InputMaybe<Scalars['Time']>;
  company?: InputMaybe<Scalars['String']>;
  startedAt?: InputMaybe<Scalars['Time']>;
  jobTitle?: InputMaybe<Scalars['String']>;
  primary?: InputMaybe<Scalars['Boolean']>;
  organizationId?: InputMaybe<Scalars['ID']>;
  description?: InputMaybe<Scalars['String']>;
};

export type LinkOrganizationsInput = {
  organizationId: Scalars['ID'];
  subOrganizationId: Scalars['ID'];
  type?: InputMaybe<Scalars['String']>;
};

export type LinkedOrganization = {
  organization: Organization;
  type?: Maybe<Scalars['String']>;
  __typename?: 'LinkedOrganization';
};

export type Location = Node &
  SourceFields & {
    id: Scalars['ID'];
    source: DataSource;
    __typename?: 'Location';
    sourceOfTruth: DataSource;
    createdAt: Scalars['Time'];
    updatedAt: Scalars['Time'];
    appSource: Scalars['String'];
    zip?: Maybe<Scalars['String']>;
    name?: Maybe<Scalars['String']>;
    region?: Maybe<Scalars['String']>;
    street?: Maybe<Scalars['String']>;
    address?: Maybe<Scalars['String']>;
    country?: Maybe<Scalars['String']>;
    latitude?: Maybe<Scalars['Float']>;
    address2?: Maybe<Scalars['String']>;
    district?: Maybe<Scalars['String']>;
    locality?: Maybe<Scalars['String']>;
    longitude?: Maybe<Scalars['Float']>;
    plusFour?: Maybe<Scalars['String']>;
    timeZone?: Maybe<Scalars['String']>;
    utcOffset?: Maybe<Scalars['Int64']>;
    postalCode?: Maybe<Scalars['String']>;
    rawAddress?: Maybe<Scalars['String']>;
    addressType?: Maybe<Scalars['String']>;
    commercial?: Maybe<Scalars['Boolean']>;
    houseNumber?: Maybe<Scalars['String']>;
    predirection?: Maybe<Scalars['String']>;
  };

export type LocationUpdateInput = {
  id: Scalars['ID'];
  zip?: InputMaybe<Scalars['String']>;
  name?: InputMaybe<Scalars['String']>;
  region?: InputMaybe<Scalars['String']>;
  street?: InputMaybe<Scalars['String']>;
  address?: InputMaybe<Scalars['String']>;
  country?: InputMaybe<Scalars['String']>;
  latitude?: InputMaybe<Scalars['Float']>;
  address2?: InputMaybe<Scalars['String']>;
  district?: InputMaybe<Scalars['String']>;
  locality?: InputMaybe<Scalars['String']>;
  longitude?: InputMaybe<Scalars['Float']>;
  plusFour?: InputMaybe<Scalars['String']>;
  timeZone?: InputMaybe<Scalars['String']>;
  utcOffset?: InputMaybe<Scalars['Int64']>;
  postalCode?: InputMaybe<Scalars['String']>;
  rawAddress?: InputMaybe<Scalars['String']>;
  addressType?: InputMaybe<Scalars['String']>;
  commercial?: InputMaybe<Scalars['Boolean']>;
  houseNumber?: InputMaybe<Scalars['String']>;
  predirection?: InputMaybe<Scalars['String']>;
};

export type LogEntry = {
  tags: Array<Tag>;
  id: Scalars['ID'];
  source: DataSource;
  __typename?: 'LogEntry';
  createdBy?: Maybe<User>;
  sourceOfTruth: DataSource;
  createdAt: Scalars['Time'];
  startedAt: Scalars['Time'];
  updatedAt: Scalars['Time'];
  appSource: Scalars['String'];
  content?: Maybe<Scalars['String']>;
  externalLinks: Array<ExternalSystem>;
  contentType?: Maybe<Scalars['String']>;
};

export type LogEntryInput = {
  content?: InputMaybe<Scalars['String']>;
  startedAt?: InputMaybe<Scalars['Time']>;
  appSource?: InputMaybe<Scalars['String']>;
  tags?: InputMaybe<Array<TagIdOrNameInput>>;
  contentType?: InputMaybe<Scalars['String']>;
};

export type LogEntryUpdateInput = {
  content?: InputMaybe<Scalars['String']>;
  startedAt?: InputMaybe<Scalars['Time']>;
  contentType?: InputMaybe<Scalars['String']>;
};

export enum Market {
  B2B = 'B2B',
  B2C = 'B2C',
  Marketplace = 'MARKETPLACE',
}

export type Meeting = Node & {
  id: Scalars['ID'];
  note: Array<Note>;
  source: DataSource;
  status: MeetingStatus;
  __typename?: 'Meeting';
  sourceOfTruth: DataSource;
  createdAt: Scalars['Time'];
  updatedAt: Scalars['Time'];
  includes: Array<Attachment>;
  appSource: Scalars['String'];
  describedBy: Array<Analysis>;
  recording?: Maybe<Attachment>;
  events: Array<InteractionEvent>;
  name?: Maybe<Scalars['String']>;
  endedAt?: Maybe<Scalars['Time']>;
  agenda?: Maybe<Scalars['String']>;
  startedAt?: Maybe<Scalars['Time']>;
  createdBy: Array<MeetingParticipant>;
  attendedBy: Array<MeetingParticipant>;
  externalSystem: Array<ExternalSystem>;
  conferenceUrl?: Maybe<Scalars['String']>;
  agendaContentType?: Maybe<Scalars['String']>;
  meetingExternalUrl?: Maybe<Scalars['String']>;
};

export type MeetingInput = {
  note?: InputMaybe<NoteInput>;
  status?: InputMaybe<MeetingStatus>;
  name?: InputMaybe<Scalars['String']>;
  endedAt?: InputMaybe<Scalars['Time']>;
  agenda?: InputMaybe<Scalars['String']>;
  createdAt?: InputMaybe<Scalars['Time']>;
  startedAt?: InputMaybe<Scalars['Time']>;
  appSource?: InputMaybe<Scalars['String']>;
  conferenceUrl?: InputMaybe<Scalars['String']>;
  agendaContentType?: InputMaybe<Scalars['String']>;
  meetingExternalUrl?: InputMaybe<Scalars['String']>;
  createdBy?: InputMaybe<Array<MeetingParticipantInput>>;
  attendedBy?: InputMaybe<Array<MeetingParticipantInput>>;
  externalSystem?: InputMaybe<ExternalSystemReferenceInput>;
};

export type MeetingParticipant =
  | ContactParticipant
  | EmailParticipant
  | OrganizationParticipant
  | UserParticipant;

export type MeetingParticipantInput = {
  userId?: InputMaybe<Scalars['ID']>;
  contactId?: InputMaybe<Scalars['ID']>;
  organizationId?: InputMaybe<Scalars['ID']>;
};

export enum MeetingStatus {
  Accepted = 'ACCEPTED',
  Canceled = 'CANCELED',
  Undefined = 'UNDEFINED',
}

export type MeetingUpdateInput = {
  note?: InputMaybe<NoteUpdateInput>;
  status?: InputMaybe<MeetingStatus>;
  name?: InputMaybe<Scalars['String']>;
  endedAt?: InputMaybe<Scalars['Time']>;
  agenda?: InputMaybe<Scalars['String']>;
  startedAt?: InputMaybe<Scalars['Time']>;
  appSource?: InputMaybe<Scalars['String']>;
  conferenceUrl?: InputMaybe<Scalars['String']>;
  agendaContentType?: InputMaybe<Scalars['String']>;
  meetingExternalUrl?: InputMaybe<Scalars['String']>;
  externalSystem?: InputMaybe<ExternalSystemReferenceInput>;
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
  /**
   * Total number of pages in the query response.
   * **Required.**
   */
  totalPages: Scalars['Int'];
  __typename?: 'MeetingsPage';
  /**
   * Total number of elements in the query response.
   * **Required.**
   */
  totalElements: Scalars['Int64'];
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
  player_Merge: Result;
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
  jobRole_Create: JobRole;
  jobRole_Update: JobRole;
  meeting_Create: Meeting;
  meeting_Update: Meeting;
  tag_Update?: Maybe<Tag>;
  workspace_Merge: Result;
  emailUpdateInUser: Email;
  meeting_AddNote: Meeting;
  analysis_Create: Analysis;
  contact_AddSocial: Social;
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
  user_RemoveRoleInTenant: User;
  contact_RemoveTagById: Contact;
  emailRemoveFromContact: Result;
  logEntry_AddTag: Scalars['ID'];
  logEntry_Update: Scalars['ID'];
  meeting_LinkRecording: Meeting;
  organization_AddSocial: Social;
  contact_RemoveLocation: Contact;
  emailMergeToOrganization: Email;
  emailRemoveFromUserById: Result;
  meeting_LinkAttachment: Meeting;
  meeting_LinkAttendedBy: Meeting;
  tenant_Merge: Scalars['String'];
  workspace_MergeToTenant: Result;
  contact_AddNewLocation: Location;
  emailUpdateInOrganization: Email;
  meeting_AddNewLocation: Location;
  meeting_UnlinkRecording: Meeting;
  note_CreateForOrganization: Note;
  organization_Hide: Scalars['ID'];
  organization_Merge: Organization;
  organization_Show: Scalars['ID'];
  fieldSetDeleteFromContact: Result;
  logEntry_RemoveTag: Scalars['ID'];
  logEntry_ResetTags: Scalars['ID'];
  meeting_UnlinkAttachment: Meeting;
  meeting_UnlinkAttendedBy: Meeting;
  organization_Create: Organization;
  organization_Update: Organization;
  contact_RestoreFromArchive: Result;
  emailRemoveFromContactById: Result;
  emailRemoveFromOrganization: Result;
  location_RemoveFromContact: Contact;
  organization_SetOwner: Organization;
  phoneNumberMergeToUser: PhoneNumber;
  contact_AddOrganizationById: Contact;
  entityTemplateCreate: EntityTemplate;
  organization_Archive?: Maybe<Result>;
  organization_HideAll?: Maybe<Result>;
  organization_ShowAll?: Maybe<Result>;
  phoneNumberUpdateInUser: PhoneNumber;
  organization_AddNewLocation: Location;
  organization_UnsetOwner: Organization;
  phoneNumberRemoveFromUserById: Result;
  customFieldMergeToContact: CustomField;
  customer_user_AddJobRole: CustomerUser;
  phoneNumberMergeToContact: PhoneNumber;
  contact_RemoveOrganizationById: Contact;
  customFieldMergeToFieldSet: CustomField;
  customFieldUpdateInContact: CustomField;
  emailRemoveFromOrganizationById: Result;
  organization_ArchiveAll?: Maybe<Result>;
  phoneNumberRemoveFromUserByE164: Result;
  phoneNumberUpdateInContact: PhoneNumber;
  customFieldDeleteFromContactById: Result;
  customFieldUpdateInFieldSet: CustomField;
  customer_contact_Create: CustomerContact;
  fieldSetMergeToContact?: Maybe<FieldSet>;
  organization_AddSubsidiary: Organization;
  phoneNumberRemoveFromContactById: Result;
  customFieldDeleteFromFieldSetById: Result;
  fieldSetUpdateInContact?: Maybe<FieldSet>;
  interactionEvent_Create: InteractionEvent;
  customFieldDeleteFromContactByName: Result;
  organization_AddRelationship: Organization;
  phoneNumberRemoveFromContactByE164: Result;
  organization_RemoveSubsidiary: Organization;
  phoneNumberMergeToOrganization: PhoneNumber;
  customFieldsMergeAndUpdateInContact: Contact;
  phoneNumberUpdateInOrganization: PhoneNumber;
  interactionSession_Create: InteractionSession;
  location_RemoveFromOrganization: Organization;
  logEntry_CreateForOrganization: Scalars['ID'];
  organization_RemoveRelationship: Organization;
  phoneNumberRemoveFromOrganizationById: Result;
  customFieldTemplate_Create: CustomFieldTemplate;
  organization_SetRelationshipStage: Organization;
  phoneNumberRemoveFromOrganizationByE164: Result;
  organization_UpdateBillingDetails: Scalars['ID'];
  interactionEvent_LinkAttachment: InteractionEvent;
  organization_UpdateRenewalForecast: Scalars['ID'];
  organization_RemoveRelationshipStage: Organization;
  organization_UpdateRenewalLikelihood: Scalars['ID'];
  interactionSession_LinkAttachment: InteractionSession;
  /** @deprecated Use organization_UpdateBillingDetails instead */
  organization_UpdateBillingDetailsAsync: Scalars['ID'];
  /** @deprecated Use organization_UpdateRenewalForecast instead */
  organization_UpdateRenewalForecastAsync: Scalars['ID'];
  /** @deprecated Use organization_UpdateRenewalLikelihood instead */
  organization_UpdateRenewalLikelihoodAsync: Scalars['ID'];
};

export type MutationAnalysis_CreateArgs = {
  analysis: AnalysisInput;
};

export type MutationAttachment_CreateArgs = {
  input: AttachmentInput;
};

export type MutationContact_AddNewLocationArgs = {
  contactId: Scalars['ID'];
};

export type MutationContact_AddOrganizationByIdArgs = {
  input: ContactOrganizationInput;
};

export type MutationContact_AddSocialArgs = {
  input: SocialInput;
  contactId: Scalars['ID'];
};

export type MutationContact_AddTagByIdArgs = {
  input: ContactTagInput;
};

export type MutationContact_ArchiveArgs = {
  contactId: Scalars['ID'];
};

export type MutationContact_CreateArgs = {
  input: ContactInput;
};

export type MutationContact_HardDeleteArgs = {
  contactId: Scalars['ID'];
};

export type MutationContact_MergeArgs = {
  primaryContactId: Scalars['ID'];
  mergedContactIds: Array<Scalars['ID']>;
};

export type MutationContact_RemoveLocationArgs = {
  contactId: Scalars['ID'];
  locationId: Scalars['ID'];
};

export type MutationContact_RemoveOrganizationByIdArgs = {
  input: ContactOrganizationInput;
};

export type MutationContact_RemoveTagByIdArgs = {
  input: ContactTagInput;
};

export type MutationContact_RestoreFromArchiveArgs = {
  contactId: Scalars['ID'];
};

export type MutationContact_UpdateArgs = {
  input: ContactUpdateInput;
};

export type MutationCustomFieldDeleteFromContactByIdArgs = {
  id: Scalars['ID'];
  contactId: Scalars['ID'];
};

export type MutationCustomFieldDeleteFromContactByNameArgs = {
  contactId: Scalars['ID'];
  fieldName: Scalars['String'];
};

export type MutationCustomFieldDeleteFromFieldSetByIdArgs = {
  id: Scalars['ID'];
  contactId: Scalars['ID'];
  fieldSetId: Scalars['ID'];
};

export type MutationCustomFieldMergeToContactArgs = {
  input: CustomFieldInput;
  contactId: Scalars['ID'];
};

export type MutationCustomFieldMergeToFieldSetArgs = {
  input: CustomFieldInput;
  contactId: Scalars['ID'];
  fieldSetId: Scalars['ID'];
};

export type MutationCustomFieldTemplate_CreateArgs = {
  input: CustomFieldTemplateInput;
};

export type MutationCustomFieldUpdateInContactArgs = {
  contactId: Scalars['ID'];
  input: CustomFieldUpdateInput;
};

export type MutationCustomFieldUpdateInFieldSetArgs = {
  contactId: Scalars['ID'];
  fieldSetId: Scalars['ID'];
  input: CustomFieldUpdateInput;
};

export type MutationCustomFieldsMergeAndUpdateInContactArgs = {
  contactId: Scalars['ID'];
  fieldSets?: InputMaybe<Array<FieldSetInput>>;
  customFields?: InputMaybe<Array<CustomFieldInput>>;
};

export type MutationCustomer_Contact_CreateArgs = {
  input: CustomerContactInput;
};

export type MutationCustomer_User_AddJobRoleArgs = {
  id: Scalars['ID'];
  jobRoleInput: JobRoleInput;
};

export type MutationEmailDeleteArgs = {
  id: Scalars['ID'];
};

export type MutationEmailMergeToContactArgs = {
  input: EmailInput;
  contactId: Scalars['ID'];
};

export type MutationEmailMergeToOrganizationArgs = {
  input: EmailInput;
  organizationId: Scalars['ID'];
};

export type MutationEmailMergeToUserArgs = {
  input: EmailInput;
  userId: Scalars['ID'];
};

export type MutationEmailRemoveFromContactArgs = {
  contactId: Scalars['ID'];
  email: Scalars['String'];
};

export type MutationEmailRemoveFromContactByIdArgs = {
  id: Scalars['ID'];
  contactId: Scalars['ID'];
};

export type MutationEmailRemoveFromOrganizationArgs = {
  email: Scalars['String'];
  organizationId: Scalars['ID'];
};

export type MutationEmailRemoveFromOrganizationByIdArgs = {
  id: Scalars['ID'];
  organizationId: Scalars['ID'];
};

export type MutationEmailRemoveFromUserArgs = {
  userId: Scalars['ID'];
  email: Scalars['String'];
};

export type MutationEmailRemoveFromUserByIdArgs = {
  id: Scalars['ID'];
  userId: Scalars['ID'];
};

export type MutationEmailUpdateInContactArgs = {
  input: EmailUpdateInput;
  contactId: Scalars['ID'];
};

export type MutationEmailUpdateInOrganizationArgs = {
  input: EmailUpdateInput;
  organizationId: Scalars['ID'];
};

export type MutationEmailUpdateInUserArgs = {
  userId: Scalars['ID'];
  input: EmailUpdateInput;
};

export type MutationEntityTemplateCreateArgs = {
  input: EntityTemplateInput;
};

export type MutationFieldSetDeleteFromContactArgs = {
  id: Scalars['ID'];
  contactId: Scalars['ID'];
};

export type MutationFieldSetMergeToContactArgs = {
  input: FieldSetInput;
  contactId: Scalars['ID'];
};

export type MutationFieldSetUpdateInContactArgs = {
  contactId: Scalars['ID'];
  input: FieldSetUpdateInput;
};

export type MutationInteractionEvent_CreateArgs = {
  event: InteractionEventInput;
};

export type MutationInteractionEvent_LinkAttachmentArgs = {
  eventId: Scalars['ID'];
  attachmentId: Scalars['ID'];
};

export type MutationInteractionSession_CreateArgs = {
  session: InteractionSessionInput;
};

export type MutationInteractionSession_LinkAttachmentArgs = {
  sessionId: Scalars['ID'];
  attachmentId: Scalars['ID'];
};

export type MutationJobRole_CreateArgs = {
  input: JobRoleInput;
  contactId: Scalars['ID'];
};

export type MutationJobRole_DeleteArgs = {
  roleId: Scalars['ID'];
  contactId: Scalars['ID'];
};

export type MutationJobRole_UpdateArgs = {
  contactId: Scalars['ID'];
  input: JobRoleUpdateInput;
};

export type MutationLocation_RemoveFromContactArgs = {
  contactId: Scalars['ID'];
  locationId: Scalars['ID'];
};

export type MutationLocation_RemoveFromOrganizationArgs = {
  locationId: Scalars['ID'];
  organizationId: Scalars['ID'];
};

export type MutationLocation_UpdateArgs = {
  input: LocationUpdateInput;
};

export type MutationLogEntry_AddTagArgs = {
  id: Scalars['ID'];
  input: TagIdOrNameInput;
};

export type MutationLogEntry_CreateForOrganizationArgs = {
  input: LogEntryInput;
  organizationId: Scalars['ID'];
};

export type MutationLogEntry_RemoveTagArgs = {
  id: Scalars['ID'];
  input: TagIdOrNameInput;
};

export type MutationLogEntry_ResetTagsArgs = {
  id: Scalars['ID'];
  input?: InputMaybe<Array<TagIdOrNameInput>>;
};

export type MutationLogEntry_UpdateArgs = {
  id: Scalars['ID'];
  input: LogEntryUpdateInput;
};

export type MutationMeeting_AddNewLocationArgs = {
  meetingId: Scalars['ID'];
};

export type MutationMeeting_AddNoteArgs = {
  meetingId: Scalars['ID'];
  note?: InputMaybe<NoteInput>;
};

export type MutationMeeting_CreateArgs = {
  meeting: MeetingInput;
};

export type MutationMeeting_LinkAttachmentArgs = {
  meetingId: Scalars['ID'];
  attachmentId: Scalars['ID'];
};

export type MutationMeeting_LinkAttendedByArgs = {
  meetingId: Scalars['ID'];
  participant: MeetingParticipantInput;
};

export type MutationMeeting_LinkRecordingArgs = {
  meetingId: Scalars['ID'];
  attachmentId: Scalars['ID'];
};

export type MutationMeeting_UnlinkAttachmentArgs = {
  meetingId: Scalars['ID'];
  attachmentId: Scalars['ID'];
};

export type MutationMeeting_UnlinkAttendedByArgs = {
  meetingId: Scalars['ID'];
  participant: MeetingParticipantInput;
};

export type MutationMeeting_UnlinkRecordingArgs = {
  meetingId: Scalars['ID'];
  attachmentId: Scalars['ID'];
};

export type MutationMeeting_UpdateArgs = {
  meetingId: Scalars['ID'];
  meeting: MeetingUpdateInput;
};

export type MutationNote_CreateForContactArgs = {
  input: NoteInput;
  contactId: Scalars['ID'];
};

export type MutationNote_CreateForOrganizationArgs = {
  input: NoteInput;
  organizationId: Scalars['ID'];
};

export type MutationNote_DeleteArgs = {
  id: Scalars['ID'];
};

export type MutationNote_LinkAttachmentArgs = {
  noteId: Scalars['ID'];
  attachmentId: Scalars['ID'];
};

export type MutationNote_UnlinkAttachmentArgs = {
  noteId: Scalars['ID'];
  attachmentId: Scalars['ID'];
};

export type MutationNote_UpdateArgs = {
  input: NoteUpdateInput;
};

export type MutationOrganization_AddNewLocationArgs = {
  organizationId: Scalars['ID'];
};

export type MutationOrganization_AddRelationshipArgs = {
  organizationId: Scalars['ID'];
  relationship: OrganizationRelationship;
};

export type MutationOrganization_AddSocialArgs = {
  input: SocialInput;
  organizationId: Scalars['ID'];
};

export type MutationOrganization_AddSubsidiaryArgs = {
  input: LinkOrganizationsInput;
};

export type MutationOrganization_ArchiveArgs = {
  id: Scalars['ID'];
};

export type MutationOrganization_ArchiveAllArgs = {
  ids: Array<Scalars['ID']>;
};

export type MutationOrganization_CreateArgs = {
  input: OrganizationInput;
};

export type MutationOrganization_HideArgs = {
  id: Scalars['ID'];
};

export type MutationOrganization_HideAllArgs = {
  ids: Array<Scalars['ID']>;
};

export type MutationOrganization_MergeArgs = {
  primaryOrganizationId: Scalars['ID'];
  mergedOrganizationIds: Array<Scalars['ID']>;
};

export type MutationOrganization_RemoveRelationshipArgs = {
  organizationId: Scalars['ID'];
  relationship: OrganizationRelationship;
};

export type MutationOrganization_RemoveRelationshipStageArgs = {
  organizationId: Scalars['ID'];
  relationship: OrganizationRelationship;
};

export type MutationOrganization_RemoveSubsidiaryArgs = {
  subsidiaryId: Scalars['ID'];
  organizationId: Scalars['ID'];
};

export type MutationOrganization_SetOwnerArgs = {
  userId: Scalars['ID'];
  organizationId: Scalars['ID'];
};

export type MutationOrganization_SetRelationshipStageArgs = {
  stage: Scalars['String'];
  organizationId: Scalars['ID'];
  relationship: OrganizationRelationship;
};

export type MutationOrganization_ShowArgs = {
  id: Scalars['ID'];
};

export type MutationOrganization_ShowAllArgs = {
  ids: Array<Scalars['ID']>;
};

export type MutationOrganization_UnsetOwnerArgs = {
  organizationId: Scalars['ID'];
};

export type MutationOrganization_UpdateArgs = {
  input: OrganizationUpdateInput;
};

export type MutationOrganization_UpdateBillingDetailsArgs = {
  input: BillingDetailsInput;
};

export type MutationOrganization_UpdateBillingDetailsAsyncArgs = {
  input: BillingDetailsInput;
};

export type MutationOrganization_UpdateRenewalForecastArgs = {
  input: RenewalForecastInput;
};

export type MutationOrganization_UpdateRenewalForecastAsyncArgs = {
  input: RenewalForecastInput;
};

export type MutationOrganization_UpdateRenewalLikelihoodArgs = {
  input: RenewalLikelihoodInput;
};

export type MutationOrganization_UpdateRenewalLikelihoodAsyncArgs = {
  input: RenewalLikelihoodInput;
};

export type MutationPhoneNumberMergeToContactArgs = {
  input: PhoneNumberInput;
  contactId: Scalars['ID'];
};

export type MutationPhoneNumberMergeToOrganizationArgs = {
  input: PhoneNumberInput;
  organizationId: Scalars['ID'];
};

export type MutationPhoneNumberMergeToUserArgs = {
  userId: Scalars['ID'];
  input: PhoneNumberInput;
};

export type MutationPhoneNumberRemoveFromContactByE164Args = {
  e164: Scalars['String'];
  contactId: Scalars['ID'];
};

export type MutationPhoneNumberRemoveFromContactByIdArgs = {
  id: Scalars['ID'];
  contactId: Scalars['ID'];
};

export type MutationPhoneNumberRemoveFromOrganizationByE164Args = {
  e164: Scalars['String'];
  organizationId: Scalars['ID'];
};

export type MutationPhoneNumberRemoveFromOrganizationByIdArgs = {
  id: Scalars['ID'];
  organizationId: Scalars['ID'];
};

export type MutationPhoneNumberRemoveFromUserByE164Args = {
  userId: Scalars['ID'];
  e164: Scalars['String'];
};

export type MutationPhoneNumberRemoveFromUserByIdArgs = {
  id: Scalars['ID'];
  userId: Scalars['ID'];
};

export type MutationPhoneNumberUpdateInContactArgs = {
  contactId: Scalars['ID'];
  input: PhoneNumberUpdateInput;
};

export type MutationPhoneNumberUpdateInOrganizationArgs = {
  input: PhoneNumberUpdateInput;
  organizationId: Scalars['ID'];
};

export type MutationPhoneNumberUpdateInUserArgs = {
  userId: Scalars['ID'];
  input: PhoneNumberUpdateInput;
};

export type MutationPlayer_MergeArgs = {
  input: PlayerInput;
  userId: Scalars['ID'];
};

export type MutationSocial_RemoveArgs = {
  socialId: Scalars['ID'];
};

export type MutationSocial_UpdateArgs = {
  input: SocialUpdateInput;
};

export type MutationTag_CreateArgs = {
  input: TagInput;
};

export type MutationTag_DeleteArgs = {
  id: Scalars['ID'];
};

export type MutationTag_UpdateArgs = {
  input: TagUpdateInput;
};

export type MutationTenant_MergeArgs = {
  tenant: TenantInput;
};

export type MutationUser_AddRoleArgs = {
  role: Role;
  id: Scalars['ID'];
};

export type MutationUser_AddRoleInTenantArgs = {
  role: Role;
  id: Scalars['ID'];
  tenant: Scalars['String'];
};

export type MutationUser_CreateArgs = {
  input: UserInput;
};

export type MutationUser_DeleteArgs = {
  id: Scalars['ID'];
};

export type MutationUser_DeleteInTenantArgs = {
  id: Scalars['ID'];
  tenant: Scalars['String'];
};

export type MutationUser_RemoveRoleArgs = {
  role: Role;
  id: Scalars['ID'];
};

export type MutationUser_RemoveRoleInTenantArgs = {
  role: Role;
  id: Scalars['ID'];
  tenant: Scalars['String'];
};

export type MutationUser_UpdateArgs = {
  input: UserUpdateInput;
};

export type MutationWorkspace_MergeArgs = {
  workspace: WorkspaceInput;
};

export type MutationWorkspace_MergeToTenantArgs = {
  tenant: Scalars['String'];
  workspace: WorkspaceInput;
};

export type Node = {
  id: Scalars['ID'];
};

export type Note = {
  id: Scalars['ID'];
  source: DataSource;
  __typename?: 'Note';
  createdBy?: Maybe<User>;
  noted: Array<NotedEntity>;
  sourceOfTruth: DataSource;
  createdAt: Scalars['Time'];
  updatedAt: Scalars['Time'];
  includes: Array<Attachment>;
  appSource: Scalars['String'];
  content?: Maybe<Scalars['String']>;
  contentType?: Maybe<Scalars['String']>;
};

export type NoteInput = {
  content?: InputMaybe<Scalars['String']>;
  appSource?: InputMaybe<Scalars['String']>;
  contentType?: InputMaybe<Scalars['String']>;
};

export type NotePage = Pages & {
  content: Array<Note>;
  __typename?: 'NotePage';
  totalPages: Scalars['Int'];
  totalElements: Scalars['Int64'];
};

export type NoteUpdateInput = {
  id: Scalars['ID'];
  content?: InputMaybe<Scalars['String']>;
  contentType?: InputMaybe<Scalars['String']>;
};

export type NotedEntity = Contact | Organization;

export type OrgAccountDetails = {
  __typename?: 'OrgAccountDetails';
  billingDetails?: Maybe<BillingDetails>;
  renewalForecast?: Maybe<RenewalForecast>;
  renewalLikelihood?: Maybe<RenewalLikelihood>;
};

export type Organization = Node & {
  notes: NotePage;
  id: Scalars['ID'];
  source: DataSource;
  owner?: Maybe<User>;
  emails: Array<Email>;
  contacts: ContactsPage;
  market?: Maybe<Market>;
  socials: Array<Social>;
  name: Scalars['String'];
  jobRoles: Array<JobRole>;
  tags?: Maybe<Array<Tag>>;
  sourceOfTruth: DataSource;
  createdAt: Scalars['Time'];
  fieldSets: Array<FieldSet>;
  locations: Array<Location>;
  updatedAt: Scalars['Time'];
  __typename?: 'Organization';
  appSource: Scalars['String'];
  customerOsId: Scalars['String'];
  note?: Maybe<Scalars['String']>;
  customFields: Array<CustomField>;
  phoneNumbers: Array<PhoneNumber>;
  domains: Array<Scalars['String']>;
  website?: Maybe<Scalars['String']>;
  employees?: Maybe<Scalars['Int64']>;
  industry?: Maybe<Scalars['String']>;
  externalLinks: Array<ExternalSystem>;
  isPublic?: Maybe<Scalars['Boolean']>;
  timelineEvents: Array<TimelineEvent>;
  description?: Maybe<Scalars['String']>;
  entityTemplate?: Maybe<EntityTemplate>;
  isCustomer?: Maybe<Scalars['Boolean']>;
  lastFundingRound?: Maybe<FundingRound>;
  referenceId?: Maybe<Scalars['String']>;
  subIndustry?: Maybe<Scalars['String']>;
  subsidiaries: Array<LinkedOrganization>;
  subsidiaryOf: Array<LinkedOrganization>;
  industryGroup?: Maybe<Scalars['String']>;
  accountDetails?: Maybe<OrgAccountDetails>;
  lastTouchPointAt?: Maybe<Scalars['Time']>;
  targetAudience?: Maybe<Scalars['String']>;
  timelineEventsTotalCount: Scalars['Int64'];
  valueProposition?: Maybe<Scalars['String']>;
  lastFundingAmount?: Maybe<Scalars['String']>;
  relationships: Array<OrganizationRelationship>;
  issueSummaryByStatus: Array<IssueSummaryByStatus>;
  lastTouchPointTimelineEvent?: Maybe<TimelineEvent>;
  suggestedMergeTo: Array<SuggestedMergeOrganization>;
  lastTouchPointTimelineEventId?: Maybe<Scalars['ID']>;
  relationshipStages: Array<OrganizationRelationshipStage>;
};

export type OrganizationContactsArgs = {
  where?: InputMaybe<Filter>;
  sort?: InputMaybe<Array<SortBy>>;
  pagination?: InputMaybe<Pagination>;
};

export type OrganizationNotesArgs = {
  pagination?: InputMaybe<Pagination>;
};

export type OrganizationTimelineEventsArgs = {
  size: Scalars['Int'];
  from?: InputMaybe<Scalars['Time']>;
  timelineEventTypes?: InputMaybe<Array<TimelineEventType>>;
};

export type OrganizationTimelineEventsTotalCountArgs = {
  timelineEventTypes?: InputMaybe<Array<TimelineEventType>>;
};

export type OrganizationInput = {
  name: Scalars['String'];
  market?: InputMaybe<Market>;
  note?: InputMaybe<Scalars['String']>;
  templateId?: InputMaybe<Scalars['ID']>;
  website?: InputMaybe<Scalars['String']>;
  employees?: InputMaybe<Scalars['Int64']>;
  industry?: InputMaybe<Scalars['String']>;
  appSource?: InputMaybe<Scalars['String']>;
  isPublic?: InputMaybe<Scalars['Boolean']>;
  description?: InputMaybe<Scalars['String']>;
  isCustomer?: InputMaybe<Scalars['Boolean']>;
  /**
   * The name of the organization.
   * **Required.**
   */
  referenceId?: InputMaybe<Scalars['String']>;
  subIndustry?: InputMaybe<Scalars['String']>;
  fieldSets?: InputMaybe<Array<FieldSetInput>>;
  industryGroup?: InputMaybe<Scalars['String']>;
  domains?: InputMaybe<Array<Scalars['String']>>;
  customFields?: InputMaybe<Array<CustomFieldInput>>;
};

export type OrganizationPage = Pages & {
  totalPages: Scalars['Int'];
  content: Array<Organization>;
  __typename?: 'OrganizationPage';
  totalElements: Scalars['Int64'];
};

export type OrganizationParticipant = {
  type?: Maybe<Scalars['String']>;
  organizationParticipant: Organization;
  __typename?: 'OrganizationParticipant';
};

export enum OrganizationRelationship {
  Affiliate = 'AFFILIATE',
  CertificationBody = 'CERTIFICATION_BODY',
  Competitor = 'COMPETITOR',
  Consultant = 'CONSULTANT',
  ContractManufacturer = 'CONTRACT_MANUFACTURER',
  Customer = 'CUSTOMER',
  DataProvider = 'DATA_PROVIDER',
  Distributor = 'DISTRIBUTOR',
  Franchisee = 'FRANCHISEE',
  Franchisor = 'FRANCHISOR',
  IndustryAnalyst = 'INDUSTRY_ANALYST',
  InfluencerOrContentCreator = 'INFLUENCER_OR_CONTENT_CREATOR',
  InsourcingPartner = 'INSOURCING_PARTNER',
  Investor = 'INVESTOR',
  JointVenture = 'JOINT_VENTURE',
  LicensingPartner = 'LICENSING_PARTNER',
  LogisticsPartner = 'LOGISTICS_PARTNER',
  MediaPartner = 'MEDIA_PARTNER',
  MergerOrAcquisitionTarget = 'MERGER_OR_ACQUISITION_TARGET',
  OriginalDesignManufacturer = 'ORIGINAL_DESIGN_MANUFACTURER',
  OriginalEquipmentManufacturer = 'ORIGINAL_EQUIPMENT_MANUFACTURER',
  OutsourcingProvider = 'OUTSOURCING_PROVIDER',
  ParentCompany = 'PARENT_COMPANY',
  Partner = 'PARTNER',
  PrivateLabelManufacturer = 'PRIVATE_LABEL_MANUFACTURER',
  ProfessionalEmployerOrganization = 'PROFESSIONAL_EMPLOYER_ORGANIZATION',
  RealEstatePartner = 'REAL_ESTATE_PARTNER',
  RegulatoryBody = 'REGULATORY_BODY',
  ResearchCollaborator = 'RESEARCH_COLLABORATOR',
  Reseller = 'RESELLER',
  ServiceProvider = 'SERVICE_PROVIDER',
  Sponsor = 'SPONSOR',
  StandardsOrganization = 'STANDARDS_ORGANIZATION',
  Subsidiary = 'SUBSIDIARY',
  Supplier = 'SUPPLIER',
  TalentAcquisitionPartner = 'TALENT_ACQUISITION_PARTNER',
  TechnologyProvider = 'TECHNOLOGY_PROVIDER',
  TradeAssociationMember = 'TRADE_ASSOCIATION_MEMBER',
  Vendor = 'VENDOR',
}

export type OrganizationRelationshipStage = {
  stage?: Maybe<Scalars['String']>;
  relationship: OrganizationRelationship;
  __typename?: 'OrganizationRelationshipStage';
};

export type OrganizationUpdateInput = {
  id: Scalars['ID'];
  name: Scalars['String'];
  market?: InputMaybe<Market>;
  note?: InputMaybe<Scalars['String']>;
  /** Set to true when partial update is needed. Empty or missing fields will not be ignored. */
  patch?: InputMaybe<Scalars['Boolean']>;
  website?: InputMaybe<Scalars['String']>;
  employees?: InputMaybe<Scalars['Int64']>;
  industry?: InputMaybe<Scalars['String']>;
  isPublic?: InputMaybe<Scalars['Boolean']>;
  description?: InputMaybe<Scalars['String']>;
  isCustomer?: InputMaybe<Scalars['Boolean']>;
  lastFundingRound?: InputMaybe<FundingRound>;
  referenceId?: InputMaybe<Scalars['String']>;
  subIndustry?: InputMaybe<Scalars['String']>;
  industryGroup?: InputMaybe<Scalars['String']>;
  domains?: InputMaybe<Array<Scalars['String']>>;
  targetAudience?: InputMaybe<Scalars['String']>;
  valueProposition?: InputMaybe<Scalars['String']>;
  lastFundingAmount?: InputMaybe<Scalars['String']>;
};

export type PageView = Node &
  SourceFields & {
    id: Scalars['ID'];
    source: DataSource;
    __typename?: 'PageView';
    endedAt: Scalars['Time'];
    sessionId: Scalars['ID'];
    sourceOfTruth: DataSource;
    pageUrl: Scalars['String'];
    startedAt: Scalars['Time'];
    appSource: Scalars['String'];
    pageTitle: Scalars['String'];
    engagedTime: Scalars['Int64'];
    application: Scalars['String'];
    orderInSession: Scalars['Int64'];
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
  totalPages: Scalars['Int'];
  /**
   * The total number of elements included in the query response.
   * **Required.**
   */
  totalElements: Scalars['Int64'];
};

/** If provided as part of the request, results will be filtered down to the `page` and `limit` specified. */
export type Pagination = {
  /**
   * The results page to return in the response.
   * **Required.**
   */
  page: Scalars['Int'];
  /**
   * The maximum number of results in the response.
   * **Required.**
   */
  limit: Scalars['Int'];
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
  /**
   * The unique ID associated with the phone number.
   * **Required**
   */
  id: Scalars['ID'];
  source: DataSource;
  users: Array<User>;
  contacts: Array<Contact>;
  country?: Maybe<Country>;
  __typename?: 'PhoneNumber';
  createdAt: Scalars['Time'];
  updatedAt: Scalars['Time'];
  /**
   * Determines if the phone number is primary or not.
   * **Required**
   */
  primary: Scalars['Boolean'];
  /** The phone number in e164 format.  */
  e164?: Maybe<Scalars['String']>;
  /** Defines the type of phone number. */
  label?: Maybe<PhoneNumberLabel>;
  organizations: Array<Organization>;
  appSource?: Maybe<Scalars['String']>;
  validated?: Maybe<Scalars['Boolean']>;
  rawPhoneNumber?: Maybe<Scalars['String']>;
};

/**
 * Describes a phone number associated with a `Contact` in customerOS.
 * **A `create` object.**
 */
export type PhoneNumberInput = {
  /**
   * The phone number in e164 format.
   * **Required**
   */
  phoneNumber: Scalars['String'];
  /** Defines the type of phone number. */
  label?: InputMaybe<PhoneNumberLabel>;
  /**
   * Determines if the phone number is primary or not.
   * **Required**
   */
  primary?: InputMaybe<Scalars['Boolean']>;
  countryCodeA2?: InputMaybe<Scalars['String']>;
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
  type?: Maybe<Scalars['String']>;
  phoneNumberParticipant: PhoneNumber;
  __typename?: 'PhoneNumberParticipant';
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
  id: Scalars['ID'];
  /** Defines the type of phone number. */
  label?: InputMaybe<PhoneNumberLabel>;
  /**
   * Determines if the phone number is primary or not.
   * **Required**
   */
  primary?: InputMaybe<Scalars['Boolean']>;
  phoneNumber?: InputMaybe<Scalars['String']>;
  countryCodeA2?: InputMaybe<Scalars['String']>;
};

export type Player = {
  id: Scalars['ID'];
  source: DataSource;
  __typename?: 'Player';
  users: Array<PlayerUser>;
  authId: Scalars['String'];
  sourceOfTruth: DataSource;
  createdAt: Scalars['Time'];
  updatedAt: Scalars['Time'];
  provider: Scalars['String'];
  appSource: Scalars['String'];
  identityId?: Maybe<Scalars['String']>;
};

export type PlayerInput = {
  authId: Scalars['String'];
  provider: Scalars['String'];
  appSource?: InputMaybe<Scalars['String']>;
  identityId?: InputMaybe<Scalars['String']>;
};

export type PlayerUpdate = {
  appSource?: InputMaybe<Scalars['String']>;
  identityId?: InputMaybe<Scalars['String']>;
};

export type PlayerUser = {
  user: User;
  __typename?: 'PlayerUser';
  tenant: Scalars['String'];
  default: Scalars['Boolean'];
};

export type Query = {
  user: User;
  email: Email;
  issue: Issue;
  users: UserPage;
  meeting: Meeting;
  tags: Array<Tag>;
  analysis: Analysis;
  logEntry: LogEntry;
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
  /** Fetch a single contact from customerOS by contact ID. */
  contact?: Maybe<Contact>;
  contact_ByEmail: Contact;
  contact_ByPhone: Contact;
  phoneNumber: PhoneNumber;
  global_Cache: GlobalCache;
  tenant: Scalars['String'];
  gcli_Search: Array<GCliItem>;
  externalMeetings: MeetingsPage;
  organizations: OrganizationPage;
  player_ByAuthIdProvider: Player;
  billableInfo: TenantBillableInfo;
  interactionEvent: InteractionEvent;
  organization?: Maybe<Organization>;
  timelineEvents: Array<TimelineEvent>;
  entityTemplates: Array<EntityTemplate>;
  interactionSession: InteractionSession;
  organization_DistinctOwners: Array<User>;
  tenant_ByEmail?: Maybe<Scalars['String']>;
  /** sort.By available options: CONTACT, EMAIL, ORGANIZATION, LOCATION */
  dashboardView_Contacts?: Maybe<ContactsPage>;
  tenant_ByWorkspace?: Maybe<Scalars['String']>;
  interactionEvent_ByEventIdentifier: InteractionEvent;
  /** sort.By available options: ORGANIZATION, IS_CUSTOMER, DOMAIN, LOCATION, OWNER, LAST_TOUCHPOINT, FORECAST_AMOUNT, RENEWAL_LIKELIHOOD, RENEWAL_CYCLE_NEXT */
  dashboardView_Organizations?: Maybe<OrganizationPage>;
  interactionSession_BySessionIdentifier: InteractionSession;
};

export type QueryAnalysisArgs = {
  id: Scalars['ID'];
};

export type QueryAttachmentArgs = {
  id: Scalars['ID'];
};

export type QueryContactArgs = {
  id: Scalars['ID'];
};

export type QueryContact_ByEmailArgs = {
  email: Scalars['String'];
};

export type QueryContact_ByPhoneArgs = {
  e164: Scalars['String'];
};

export type QueryContactsArgs = {
  where?: InputMaybe<Filter>;
  sort?: InputMaybe<Array<SortBy>>;
  pagination?: InputMaybe<Pagination>;
};

export type QueryDashboardView_ContactsArgs = {
  pagination: Pagination;
  sort?: InputMaybe<SortBy>;
  where?: InputMaybe<Filter>;
};

export type QueryDashboardView_OrganizationsArgs = {
  pagination: Pagination;
  sort?: InputMaybe<SortBy>;
  where?: InputMaybe<Filter>;
};

export type QueryEmailArgs = {
  id: Scalars['ID'];
};

export type QueryEntityTemplatesArgs = {
  extends?: InputMaybe<EntityTemplateExtension>;
};

export type QueryExternalMeetingsArgs = {
  where?: InputMaybe<Filter>;
  sort?: InputMaybe<Array<SortBy>>;
  externalSystemId: Scalars['String'];
  pagination?: InputMaybe<Pagination>;
  externalId?: InputMaybe<Scalars['ID']>;
};

export type QueryGcli_SearchArgs = {
  keyword: Scalars['String'];
  limit?: InputMaybe<Scalars['Int']>;
};

export type QueryInteractionEventArgs = {
  id: Scalars['ID'];
};

export type QueryInteractionEvent_ByEventIdentifierArgs = {
  eventIdentifier: Scalars['String'];
};

export type QueryInteractionSessionArgs = {
  id: Scalars['ID'];
};

export type QueryInteractionSession_BySessionIdentifierArgs = {
  sessionIdentifier: Scalars['String'];
};

export type QueryIssueArgs = {
  id: Scalars['ID'];
};

export type QueryLogEntryArgs = {
  id: Scalars['ID'];
};

export type QueryMeetingArgs = {
  id: Scalars['ID'];
};

export type QueryOrganizationArgs = {
  id: Scalars['ID'];
};

export type QueryOrganizationsArgs = {
  where?: InputMaybe<Filter>;
  sort?: InputMaybe<Array<SortBy>>;
  pagination?: InputMaybe<Pagination>;
};

export type QueryPhoneNumberArgs = {
  id: Scalars['ID'];
};

export type QueryPlayer_ByAuthIdProviderArgs = {
  authId: Scalars['String'];
  provider: Scalars['String'];
};

export type QueryTenant_ByEmailArgs = {
  email: Scalars['String'];
};

export type QueryTenant_ByWorkspaceArgs = {
  workspace: WorkspaceInput;
};

export type QueryTimelineEventsArgs = {
  ids: Array<Scalars['ID']>;
};

export type QueryUserArgs = {
  id: Scalars['ID'];
};

export type QueryUser_ByEmailArgs = {
  email: Scalars['String'];
};

export type QueryUsersArgs = {
  where?: InputMaybe<Filter>;
  sort?: InputMaybe<Array<SortBy>>;
  pagination?: InputMaybe<Pagination>;
};

export enum RenewalCycle {
  Annually = 'ANNUALLY',
  Biannually = 'BIANNUALLY',
  Biweekly = 'BIWEEKLY',
  Monthly = 'MONTHLY',
  Quarterly = 'QUARTERLY',
  Weekly = 'WEEKLY',
}

export type RenewalForecast = {
  updatedBy?: Maybe<User>;
  __typename?: 'RenewalForecast';
  amount?: Maybe<Scalars['Float']>;
  comment?: Maybe<Scalars['String']>;
  updatedAt?: Maybe<Scalars['Time']>;
  updatedById?: Maybe<Scalars['String']>;
  potentialAmount?: Maybe<Scalars['Float']>;
};

export type RenewalForecastInput = {
  id: Scalars['ID'];
  amount?: InputMaybe<Scalars['Float']>;
  comment?: InputMaybe<Scalars['String']>;
};

export type RenewalLikelihood = {
  updatedBy?: Maybe<User>;
  __typename?: 'RenewalLikelihood';
  comment?: Maybe<Scalars['String']>;
  updatedAt?: Maybe<Scalars['Time']>;
  updatedById?: Maybe<Scalars['String']>;
  probability?: Maybe<RenewalLikelihoodProbability>;
  previousProbability?: Maybe<RenewalLikelihoodProbability>;
};

export type RenewalLikelihoodInput = {
  id: Scalars['ID'];
  comment?: InputMaybe<Scalars['String']>;
  probability?: InputMaybe<RenewalLikelihoodProbability>;
};

export enum RenewalLikelihoodProbability {
  High = 'HIGH',
  Low = 'LOW',
  Medium = 'MEDIUM',
  Zero = 'ZERO',
}

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
  result: Scalars['Boolean'];
};

export enum Role {
  Admin = 'ADMIN',
  CustomerOsPlatformOwner = 'CUSTOMER_OS_PLATFORM_OWNER',
  Owner = 'OWNER',
  User = 'USER',
}

export type Social = Node &
  SourceFields & {
    id: Scalars['ID'];
    source: DataSource;
    __typename?: 'Social';
    url: Scalars['String'];
    sourceOfTruth: DataSource;
    createdAt: Scalars['Time'];
    updatedAt: Scalars['Time'];
    appSource: Scalars['String'];
    platformName?: Maybe<Scalars['String']>;
  };

export type SocialInput = {
  url: Scalars['String'];
  appSource?: InputMaybe<Scalars['String']>;
  platformName?: InputMaybe<Scalars['String']>;
};

export type SocialUpdateInput = {
  id: Scalars['ID'];
  url: Scalars['String'];
  platformName?: InputMaybe<Scalars['String']>;
};

export type SortBy = {
  by: Scalars['String'];
  direction?: SortingDirection;
  caseSensitive?: InputMaybe<Scalars['Boolean']>;
};

export enum SortingDirection {
  Asc = 'ASC',
  Desc = 'DESC',
}

export type SourceFields = {
  id: Scalars['ID'];
  source: DataSource;
  sourceOfTruth: DataSource;
  appSource: Scalars['String'];
};

export type State = {
  country: Country;
  id: Scalars['ID'];
  __typename?: 'State';
  code: Scalars['String'];
  name: Scalars['String'];
};

export type SuggestedMergeOrganization = {
  organization: Organization;
  confidence?: Maybe<Scalars['Float']>;
  suggestedAt?: Maybe<Scalars['Time']>;
  suggestedBy?: Maybe<Scalars['String']>;
  __typename?: 'SuggestedMergeOrganization';
};

export type Tag = {
  id: Scalars['ID'];
  __typename?: 'Tag';
  source: DataSource;
  name: Scalars['String'];
  createdAt: Scalars['Time'];
  updatedAt: Scalars['Time'];
  appSource: Scalars['String'];
};

export type TagIdOrNameInput = {
  id?: InputMaybe<Scalars['ID']>;
  name?: InputMaybe<Scalars['String']>;
};

export type TagInput = {
  name: Scalars['String'];
  appSource?: InputMaybe<Scalars['String']>;
};

export type TagUpdateInput = {
  id: Scalars['ID'];
  name: Scalars['String'];
};

export type TenantBillableInfo = {
  __typename?: 'TenantBillableInfo';
  greylistedContacts: Scalars['Int64'];
  whitelistedContacts: Scalars['Int64'];
  greylistedOrganizations: Scalars['Int64'];
  whitelistedOrganizations: Scalars['Int64'];
};

export type TenantInput = {
  name: Scalars['String'];
  appSource?: InputMaybe<Scalars['String']>;
};

export type TimeRange = {
  /**
   * The end time of the time range.
   * **Required.**
   */
  to: Scalars['Time'];
  /**
   * The start time of the time range.
   * **Required.**
   */
  from: Scalars['Time'];
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
  PageView = 'PAGE_VIEW',
}

/**
 * Describes the User of customerOS.  A user is the person who logs into the Openline platform.
 * **A `return` object**
 */
export type User = {
  player: Player;
  /**
   * The unique ID associated with the customerOS user.
   * **Required**
   */
  id: Scalars['ID'];
  roles: Array<Role>;
  source: DataSource;
  __typename?: 'User';
  jobRoles: Array<JobRole>;
  sourceOfTruth: DataSource;
  calendars: Array<Calendar>;
  /**
   * Timestamp of user creation.
   * **Required**
   */
  createdAt: Scalars['Time'];
  updatedAt: Scalars['Time'];
  /**
   * The last name of the customerOS user.
   * **Required**
   */
  lastName: Scalars['String'];
  appSource: Scalars['String'];
  /**
   * All email addresses associated with a user in customerOS.
   * **Required.  If no values it returns an empty array.**
   */
  emails?: Maybe<Array<Email>>;
  /**
   * The first name of the customerOS user.
   * **Required**
   */
  firstName: Scalars['String'];
  internal: Scalars['Boolean'];
  name?: Maybe<Scalars['String']>;
  phoneNumbers: Array<PhoneNumber>;
  timezone?: Maybe<Scalars['String']>;
  profilePhotoUrl?: Maybe<Scalars['String']>;
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
  lastName: Scalars['String'];
  /**
   * The first name of the customerOS user.
   * **Required**
   */
  firstName: Scalars['String'];
  name?: InputMaybe<Scalars['String']>;
  timezone?: InputMaybe<Scalars['String']>;
  /**
   * The name of the app performing the create.
   * **Optional**
   */
  appSource?: InputMaybe<Scalars['String']>;
  /**
   * The Job Roles of the user.
   * **Optional**
   */
  jobRoles?: InputMaybe<Array<JobRoleInput>>;
  profilePhotoUrl?: InputMaybe<Scalars['String']>;
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
  totalPages: Scalars['Int'];
  /**
   * Total number of elements in the query response.
   * **Required.**
   */
  totalElements: Scalars['Int64'];
};

export type UserParticipant = {
  userParticipant: User;
  __typename?: 'UserParticipant';
  type?: Maybe<Scalars['String']>;
};

export type UserUpdateInput = {
  id: Scalars['ID'];
  /**
   * The last name of the customerOS user.
   * **Required**
   */
  lastName: Scalars['String'];
  /**
   * The first name of the customerOS user.
   * **Required**
   */
  firstName: Scalars['String'];
  name?: InputMaybe<Scalars['String']>;
  timezone?: InputMaybe<Scalars['String']>;
  profilePhotoUrl?: InputMaybe<Scalars['String']>;
};

export type Workspace = {
  id: Scalars['ID'];
  source: DataSource;
  name: Scalars['String'];
  __typename?: 'Workspace';
  sourceOfTruth: DataSource;
  createdAt: Scalars['Time'];
  updatedAt: Scalars['Time'];
  provider: Scalars['String'];
  appSource: Scalars['String'];
};

export type WorkspaceInput = {
  name: Scalars['String'];
  provider: Scalars['String'];
  appSource?: InputMaybe<Scalars['String']>;
};
