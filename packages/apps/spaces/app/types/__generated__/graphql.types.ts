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
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
  Any: any;
  Int64: any;
  Time: any;
};

export type Action = {
  __typename?: 'Action';
  actionType: ActionType;
  appSource: Scalars['String'];
  content?: Maybe<Scalars['String']>;
  createdAt: Scalars['Time'];
  createdBy?: Maybe<User>;
  id: Scalars['ID'];
  metadata?: Maybe<Scalars['String']>;
  source: DataSource;
};

export type ActionItem = {
  __typename?: 'ActionItem';
  appSource: Scalars['String'];
  content: Scalars['String'];
  createdAt: Scalars['Time'];
  id: Scalars['ID'];
  source: DataSource;
};

export enum ActionType {
  Created = 'CREATED',
  RenewalForecastUpdated = 'RENEWAL_FORECAST_UPDATED',
  RenewalLikelihoodUpdated = 'RENEWAL_LIKELIHOOD_UPDATED',
}

export type Analysis = Node & {
  __typename?: 'Analysis';
  analysisType?: Maybe<Scalars['String']>;
  appSource: Scalars['String'];
  content?: Maybe<Scalars['String']>;
  contentType?: Maybe<Scalars['String']>;
  createdAt: Scalars['Time'];
  describes: Array<DescriptionNode>;
  id: Scalars['ID'];
  source: DataSource;
  sourceOfTruth: DataSource;
};

export type AnalysisDescriptionInput = {
  interactionEventId?: InputMaybe<Scalars['ID']>;
  interactionSessionId?: InputMaybe<Scalars['ID']>;
  meetingId?: InputMaybe<Scalars['ID']>;
};

export type AnalysisInput = {
  analysisType?: InputMaybe<Scalars['String']>;
  appSource: Scalars['String'];
  content?: InputMaybe<Scalars['String']>;
  contentType?: InputMaybe<Scalars['String']>;
  describes: Array<AnalysisDescriptionInput>;
};

export type Attachment = Node & {
  __typename?: 'Attachment';
  appSource: Scalars['String'];
  createdAt: Scalars['Time'];
  extension: Scalars['String'];
  id: Scalars['ID'];
  mimeType: Scalars['String'];
  name: Scalars['String'];
  size: Scalars['Int64'];
  source: DataSource;
  sourceOfTruth: DataSource;
};

export type AttachmentInput = {
  appSource: Scalars['String'];
  extension: Scalars['String'];
  mimeType: Scalars['String'];
  name: Scalars['String'];
  size: Scalars['Int64'];
};

export type BillingDetails = {
  __typename?: 'BillingDetails';
  amount?: Maybe<Scalars['Float']>;
  frequency?: Maybe<RenewalCycle>;
  renewalCycle?: Maybe<RenewalCycle>;
  renewalCycleNext?: Maybe<Scalars['Time']>;
  renewalCycleStart?: Maybe<Scalars['Time']>;
};

export type BillingDetailsInput = {
  amount?: InputMaybe<Scalars['Float']>;
  frequency?: InputMaybe<RenewalCycle>;
  id: Scalars['ID'];
  renewalCycle?: InputMaybe<RenewalCycle>;
  renewalCycleStart?: InputMaybe<Scalars['Time']>;
};

/**
 * Describes the relationship a Contact has with a Organization.
 * **A `return` object**
 */
export type Calendar = {
  __typename?: 'Calendar';
  appSource: Scalars['String'];
  calType: CalendarType;
  createdAt: Scalars['Time'];
  id: Scalars['ID'];
  link?: Maybe<Scalars['String']>;
  primary: Scalars['Boolean'];
  source: DataSource;
  sourceOfTruth: DataSource;
  updatedAt: Scalars['Time'];
};

export enum CalendarType {
  Calcom = 'CALCOM',
  Google = 'GOOGLE',
}

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
    __typename?: 'Contact';
    appSource?: Maybe<Scalars['String']>;
    /**
     * An ISO8601 timestamp recording when the contact was created in customerOS.
     * **Required**
     */
    createdAt: Scalars['Time'];
    /**
     * User defined metadata appended to the contact record in customerOS.
     * **Required.  If no values it returns an empty array.**
     */
    customFields: Array<CustomField>;
    description?: Maybe<Scalars['String']>;
    /**
     * All email addresses associated with a contact in customerOS.
     * **Required.  If no values it returns an empty array.**
     */
    emails: Array<Email>;
    fieldSets: Array<FieldSet>;
    /** The first name of the contact in customerOS. */
    firstName?: Maybe<Scalars['String']>;
    /**
     * The unique ID associated with the contact in customerOS.
     * **Required**
     */
    id: Scalars['ID'];
    /**
     * `organizationName` and `jobTitle` of the contact if it has been associated with an organization.
     * **Required.  If no values it returns an empty array.**
     */
    jobRoles: Array<JobRole>;
    /** @deprecated Use `tags` instead */
    label?: Maybe<Scalars['String']>;
    /** The last name of the contact in customerOS. */
    lastName?: Maybe<Scalars['String']>;
    /**
     * All locations associated with a contact in customerOS.
     * **Required.  If no values it returns an empty array.**
     */
    locations: Array<Location>;
    /** The name of the contact in customerOS, alternative for firstName + lastName. */
    name?: Maybe<Scalars['String']>;
    /** Contact notes */
    notes: NotePage;
    notesByTime: Array<Note>;
    organizations: OrganizationPage;
    /** Contact owner (user) */
    owner?: Maybe<User>;
    /**
     * All phone numbers associated with a contact in customerOS.
     * **Required.  If no values it returns an empty array.**
     */
    phoneNumbers: Array<PhoneNumber>;
    prefix?: Maybe<Scalars['String']>;
    profilePhotoUrl?: Maybe<Scalars['String']>;
    socials: Array<Social>;
    source: DataSource;
    sourceOfTruth: DataSource;
    tags?: Maybe<Array<Tag>>;
    /** Template of the contact in customerOS. */
    template?: Maybe<EntityTemplate>;
    timelineEvents: Array<TimelineEvent>;
    timelineEventsTotalCount: Scalars['Int64'];
    timezone?: Maybe<Scalars['String']>;
    /**
     * The title associate with the contact in customerOS.
     * @deprecated Use `prefix` instead
     */
    title?: Maybe<Scalars['String']>;
    updatedAt: Scalars['Time'];
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
  pagination?: InputMaybe<Pagination>;
  sort?: InputMaybe<Array<SortBy>>;
  where?: InputMaybe<Filter>;
};

/**
 * A contact represents an individual in customerOS.
 * **A `response` object.**
 */
export type ContactTimelineEventsArgs = {
  from?: InputMaybe<Scalars['Time']>;
  size: Scalars['Int'];
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
  appSource?: InputMaybe<Scalars['String']>;
  /** An ISO8601 timestamp recording when the contact was created in customerOS. */
  createdAt?: InputMaybe<Scalars['Time']>;
  /**
   * User defined metadata appended to contact.
   * **Required.**
   */
  customFields?: InputMaybe<Array<CustomFieldInput>>;
  description?: InputMaybe<Scalars['String']>;
  /** An email addresses associted with the contact. */
  email?: InputMaybe<EmailInput>;
  externalReference?: InputMaybe<ExternalSystemReferenceInput>;
  fieldSets?: InputMaybe<Array<FieldSetInput>>;
  /** The first name of the contact. */
  firstName?: InputMaybe<Scalars['String']>;
  label?: InputMaybe<Scalars['String']>;
  /** The last name of the contact. */
  lastName?: InputMaybe<Scalars['String']>;
  name?: InputMaybe<Scalars['String']>;
  /** Id of the contact owner (user) */
  ownerId?: InputMaybe<Scalars['ID']>;
  /** A phone number associated with the contact. */
  phoneNumber?: InputMaybe<PhoneNumberInput>;
  /** The prefix of the contact. */
  prefix?: InputMaybe<Scalars['String']>;
  /** The unique ID associated with the template of the contact in customerOS. */
  templateId?: InputMaybe<Scalars['ID']>;
  timezone?: InputMaybe<Scalars['String']>;
};

export type ContactOrganizationInput = {
  contactId: Scalars['ID'];
  organizationId: Scalars['ID'];
};

export type ContactParticipant = {
  __typename?: 'ContactParticipant';
  contactParticipant: Contact;
  type?: Maybe<Scalars['String']>;
};

export type ContactTagInput = {
  contactId: Scalars['ID'];
  tagId: Scalars['ID'];
};

/**
 * Updates data fields associated with an existing customer record in customerOS.
 * **An `update` object.**
 */
export type ContactUpdateInput = {
  description?: InputMaybe<Scalars['String']>;
  /** The first name of the contact in customerOS. */
  firstName?: InputMaybe<Scalars['String']>;
  /**
   * The unique ID associated with the contact in customerOS.
   * **Required.**
   */
  id: Scalars['ID'];
  label?: InputMaybe<Scalars['String']>;
  /** The last name of the contact in customerOS. */
  lastName?: InputMaybe<Scalars['String']>;
  name?: InputMaybe<Scalars['String']>;
  /** Id of the contact owner (user) */
  ownerId?: InputMaybe<Scalars['ID']>;
  /** The prefix associate with the contact in customerOS. */
  prefix?: InputMaybe<Scalars['String']>;
  timezone?: InputMaybe<Scalars['String']>;
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
  totalElements: Scalars['Int64'];
  /**
   * Total number of pages in the query response.
   * **Required.**
   */
  totalPages: Scalars['Int'];
};

export type Country = {
  __typename?: 'Country';
  codeA2: Scalars['String'];
  codeA3: Scalars['String'];
  id: Scalars['ID'];
  name: Scalars['String'];
  phoneCode: Scalars['String'];
};

/**
 * Describes a custom, user-defined field associated with a `Contact`.
 * **A `return` object.**
 */
export type CustomField = Node & {
  __typename?: 'CustomField';
  createdAt: Scalars['Time'];
  /**
   * Datatype of the custom field.
   * **Required**
   */
  datatype: CustomFieldDataType;
  /**
   * The unique ID associated with the custom field.
   * **Required**
   */
  id: Scalars['ID'];
  /**
   * The name of the custom field.
   * **Required**
   */
  name: Scalars['String'];
  /** The source of the custom field value */
  source: DataSource;
  template?: Maybe<CustomFieldTemplate>;
  updatedAt: Scalars['Time'];
  /**
   * The value of the custom field.
   * **Required**
   */
  value: Scalars['Any'];
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
  id: Scalars['ID'];
};

/**
 * Describes a custom, user-defined field associated with a `Contact` of type String.
 * **A `create` object.**
 */
export type CustomFieldInput = {
  /**
   * Datatype of the custom field.
   * **Required**
   */
  datatype: CustomFieldDataType;
  /** The unique ID associated with the custom field. */
  id?: InputMaybe<Scalars['ID']>;
  /**
   * The name of the custom field.
   * **Required**
   */
  name: Scalars['String'];
  templateId?: InputMaybe<Scalars['ID']>;
  /**
   * The value of the custom field.
   * **Required**
   */
  value: Scalars['Any'];
};

export type CustomFieldTemplate = Node & {
  __typename?: 'CustomFieldTemplate';
  createdAt: Scalars['Time'];
  id: Scalars['ID'];
  length?: Maybe<Scalars['Int']>;
  mandatory: Scalars['Boolean'];
  max?: Maybe<Scalars['Int']>;
  min?: Maybe<Scalars['Int']>;
  name: Scalars['String'];
  order: Scalars['Int'];
  type: CustomFieldTemplateType;
  updatedAt: Scalars['Time'];
};

export type CustomFieldTemplateInput = {
  length?: InputMaybe<Scalars['Int']>;
  mandatory: Scalars['Boolean'];
  max?: InputMaybe<Scalars['Int']>;
  min?: InputMaybe<Scalars['Int']>;
  name: Scalars['String'];
  order: Scalars['Int'];
  type: CustomFieldTemplateType;
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
   * Datatype of the custom field.
   * **Required**
   */
  datatype: CustomFieldDataType;
  /**
   * The unique ID associated with the custom field.
   * **Required**
   */
  id: Scalars['ID'];
  /**
   * The name of the custom field.
   * **Required**
   */
  name: Scalars['String'];
  /**
   * The value of the custom field.
   * **Required**
   */
  value: Scalars['Any'];
};

export type CustomerContact = {
  __typename?: 'CustomerContact';
  email: CustomerEmail;
  id: Scalars['ID'];
};

export type CustomerContactInput = {
  appSource?: InputMaybe<Scalars['String']>;
  /** An ISO8601 timestamp recording when the contact was created in customerOS. */
  createdAt?: InputMaybe<Scalars['Time']>;
  description?: InputMaybe<Scalars['String']>;
  /** An email addresses associted with the contact. */
  email?: InputMaybe<EmailInput>;
  /** The first name of the contact. */
  firstName?: InputMaybe<Scalars['String']>;
  /** The last name of the contact. */
  lastName?: InputMaybe<Scalars['String']>;
  name?: InputMaybe<Scalars['String']>;
  /** The prefix of the contact. */
  prefix?: InputMaybe<Scalars['String']>;
  timezone?: InputMaybe<Scalars['String']>;
};

export type CustomerEmail = {
  __typename?: 'CustomerEmail';
  id: Scalars['ID'];
};

export type CustomerJobRole = {
  __typename?: 'CustomerJobRole';
  id: Scalars['ID'];
};

export type CustomerUser = {
  __typename?: 'CustomerUser';
  id: Scalars['ID'];
  jobRole: CustomerJobRole;
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
  __typename?: 'Email';
  appSource: Scalars['String'];
  contacts: Array<Contact>;
  createdAt: Scalars['Time'];
  /** An email address assocaited with the contact in customerOS. */
  email?: Maybe<Scalars['String']>;
  emailValidationDetails: EmailValidationDetails;
  /**
   * The unique ID associated with the contact in customerOS.
   * **Required**
   */
  id: Scalars['ID'];
  /** Describes the type of email address (WORK, PERSONAL, etc). */
  label?: Maybe<EmailLabel>;
  organizations: Array<Organization>;
  /**
   * Identifies whether the email address is primary or not.
   * **Required.**
   */
  primary: Scalars['Boolean'];
  rawEmail?: Maybe<Scalars['String']>;
  source: DataSource;
  sourceOfTruth: DataSource;
  updatedAt: Scalars['Time'];
  users: Array<User>;
};

/**
 * Describes an email address associated with a `Contact` in customerOS.
 * **A `create` object.**
 */
export type EmailInput = {
  appSource?: InputMaybe<Scalars['String']>;
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
  type?: Maybe<Scalars['String']>;
};

/**
 * Describes an email address associated with a `Contact` in customerOS.
 * **An `update` object.**
 */
export type EmailUpdateInput = {
  email?: InputMaybe<Scalars['String']>;
  /**
   * An email address assocaited with the contact in customerOS.
   * **Required.**
   */
  id: Scalars['ID'];
  /** Describes the type of email address (WORK, PERSONAL, etc). */
  label?: InputMaybe<EmailLabel>;
  /**
   * Identifies whether the email address is primary or not.
   * **Required.**
   */
  primary?: InputMaybe<Scalars['Boolean']>;
};

export type EmailValidationDetails = {
  __typename?: 'EmailValidationDetails';
  acceptsMail?: Maybe<Scalars['Boolean']>;
  canConnectSmtp?: Maybe<Scalars['Boolean']>;
  error?: Maybe<Scalars['String']>;
  hasFullInbox?: Maybe<Scalars['Boolean']>;
  isCatchAll?: Maybe<Scalars['Boolean']>;
  isDeliverable?: Maybe<Scalars['Boolean']>;
  isDisabled?: Maybe<Scalars['Boolean']>;
  isReachable?: Maybe<Scalars['String']>;
  isValidSyntax?: Maybe<Scalars['Boolean']>;
  validated?: Maybe<Scalars['Boolean']>;
};

export type EntityTemplate = Node & {
  __typename?: 'EntityTemplate';
  createdAt: Scalars['Time'];
  customFieldTemplates: Array<CustomFieldTemplate>;
  extends?: Maybe<EntityTemplateExtension>;
  fieldSetTemplates: Array<FieldSetTemplate>;
  id: Scalars['ID'];
  name: Scalars['String'];
  updatedAt: Scalars['Time'];
  version: Scalars['Int'];
};

export enum EntityTemplateExtension {
  Contact = 'CONTACT',
  Organization = 'ORGANIZATION',
}

export type EntityTemplateInput = {
  customFieldTemplateInputs?: InputMaybe<Array<CustomFieldTemplateInput>>;
  extends?: InputMaybe<EntityTemplateExtension>;
  fieldSetTemplateInputs?: InputMaybe<Array<FieldSetTemplateInput>>;
  name: Scalars['String'];
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
  __typename?: 'ExternalSystem';
  externalId?: Maybe<Scalars['String']>;
  externalSource?: Maybe<Scalars['String']>;
  externalUrl?: Maybe<Scalars['String']>;
  syncDate?: Maybe<Scalars['Time']>;
  type: ExternalSystemType;
};

export type ExternalSystemReferenceInput = {
  externalId: Scalars['ID'];
  externalSource?: InputMaybe<Scalars['String']>;
  externalUrl?: InputMaybe<Scalars['String']>;
  syncDate?: InputMaybe<Scalars['Time']>;
  type: ExternalSystemType;
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
  __typename?: 'FieldSet';
  createdAt: Scalars['Time'];
  customFields: Array<CustomField>;
  id: Scalars['ID'];
  name: Scalars['String'];
  source: DataSource;
  template?: Maybe<FieldSetTemplate>;
  updatedAt: Scalars['Time'];
};

export type FieldSetInput = {
  customFields?: InputMaybe<Array<CustomFieldInput>>;
  id?: InputMaybe<Scalars['ID']>;
  name: Scalars['String'];
  templateId?: InputMaybe<Scalars['ID']>;
};

export type FieldSetTemplate = Node & {
  __typename?: 'FieldSetTemplate';
  createdAt: Scalars['Time'];
  customFieldTemplates: Array<CustomFieldTemplate>;
  id: Scalars['ID'];
  name: Scalars['String'];
  order: Scalars['Int'];
  updatedAt: Scalars['Time'];
};

export type FieldSetTemplateInput = {
  customFieldTemplateInputs?: InputMaybe<Array<CustomFieldTemplateInput>>;
  name: Scalars['String'];
  order: Scalars['Int'];
};

export type FieldSetUpdateInput = {
  id: Scalars['ID'];
  name: Scalars['String'];
};

export type Filter = {
  AND?: InputMaybe<Array<Filter>>;
  NOT?: InputMaybe<Filter>;
  OR?: InputMaybe<Array<Filter>>;
  filter?: InputMaybe<FilterItem>;
};

export type FilterItem = {
  caseSensitive?: InputMaybe<Scalars['Boolean']>;
  operation?: ComparisonOperator;
  property: Scalars['String'];
  value: Scalars['Any'];
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
  __typename?: 'GCliAttributeKeyValuePair';
  display?: Maybe<Scalars['String']>;
  key: Scalars['String'];
  value: Scalars['String'];
};

export enum GCliCacheItemType {
  Contact = 'CONTACT',
  Organization = 'ORGANIZATION',
  State = 'STATE',
}

export type GCliItem = {
  __typename?: 'GCliItem';
  data?: Maybe<Array<GCliAttributeKeyValuePair>>;
  display: Scalars['String'];
  id: Scalars['ID'];
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
  gCliCache: Array<GCliItem>;
  gmailOauthTokenNeedsManualRefresh: Scalars['Boolean'];
  isOwner: Scalars['Boolean'];
  user: User;
};

export type InteractionEvent = Node & {
  __typename?: 'InteractionEvent';
  actionItems?: Maybe<Array<ActionItem>>;
  appSource: Scalars['String'];
  channel?: Maybe<Scalars['String']>;
  channelData?: Maybe<Scalars['String']>;
  content?: Maybe<Scalars['String']>;
  contentType?: Maybe<Scalars['String']>;
  createdAt: Scalars['Time'];
  eventIdentifier?: Maybe<Scalars['String']>;
  eventType?: Maybe<Scalars['String']>;
  externalLinks: Array<ExternalSystem>;
  id: Scalars['ID'];
  includes: Array<Attachment>;
  interactionSession?: Maybe<InteractionSession>;
  issue?: Maybe<Issue>;
  meeting?: Maybe<Meeting>;
  repliesTo?: Maybe<InteractionEvent>;
  sentBy: Array<InteractionEventParticipant>;
  sentTo: Array<InteractionEventParticipant>;
  source: DataSource;
  sourceOfTruth: DataSource;
  summary?: Maybe<Analysis>;
};

export type InteractionEventInput = {
  appSource: Scalars['String'];
  channel?: InputMaybe<Scalars['String']>;
  channelData?: InputMaybe<Scalars['String']>;
  content?: InputMaybe<Scalars['String']>;
  contentType?: InputMaybe<Scalars['String']>;
  createdAt?: InputMaybe<Scalars['Time']>;
  eventIdentifier?: InputMaybe<Scalars['String']>;
  eventType?: InputMaybe<Scalars['String']>;
  interactionSession?: InputMaybe<Scalars['ID']>;
  meetingId?: InputMaybe<Scalars['ID']>;
  repliesTo?: InputMaybe<Scalars['ID']>;
  sentBy: Array<InteractionEventParticipantInput>;
  sentTo: Array<InteractionEventParticipantInput>;
};

export type InteractionEventParticipant =
  | ContactParticipant
  | EmailParticipant
  | JobRoleParticipant
  | OrganizationParticipant
  | PhoneNumberParticipant
  | UserParticipant;

export type InteractionEventParticipantInput = {
  contactID?: InputMaybe<Scalars['ID']>;
  email?: InputMaybe<Scalars['String']>;
  phoneNumber?: InputMaybe<Scalars['String']>;
  type?: InputMaybe<Scalars['String']>;
  userID?: InputMaybe<Scalars['ID']>;
};

export type InteractionSession = Node & {
  __typename?: 'InteractionSession';
  appSource: Scalars['String'];
  attendedBy: Array<InteractionSessionParticipant>;
  channel?: Maybe<Scalars['String']>;
  channelData?: Maybe<Scalars['String']>;
  createdAt: Scalars['Time'];
  describedBy: Array<Analysis>;
  /** @deprecated Use updatedAt instead */
  endedAt?: Maybe<Scalars['Time']>;
  events: Array<InteractionEvent>;
  id: Scalars['ID'];
  includes: Array<Attachment>;
  name: Scalars['String'];
  sessionIdentifier?: Maybe<Scalars['String']>;
  source: DataSource;
  sourceOfTruth: DataSource;
  /** @deprecated Use createdAt instead */
  startedAt: Scalars['Time'];
  status: Scalars['String'];
  type?: Maybe<Scalars['String']>;
  updatedAt: Scalars['Time'];
};

export type InteractionSessionInput = {
  appSource: Scalars['String'];
  attendedBy?: InputMaybe<Array<InteractionSessionParticipantInput>>;
  channel?: InputMaybe<Scalars['String']>;
  channelData?: InputMaybe<Scalars['String']>;
  name: Scalars['String'];
  sessionIdentifier?: InputMaybe<Scalars['String']>;
  status: Scalars['String'];
  type?: InputMaybe<Scalars['String']>;
};

export type InteractionSessionParticipant =
  | ContactParticipant
  | EmailParticipant
  | PhoneNumberParticipant
  | UserParticipant;

export type InteractionSessionParticipantInput = {
  contactID?: InputMaybe<Scalars['ID']>;
  email?: InputMaybe<Scalars['String']>;
  phoneNumber?: InputMaybe<Scalars['String']>;
  type?: InputMaybe<Scalars['String']>;
  userID?: InputMaybe<Scalars['ID']>;
};

export type Issue = Node &
  SourceFields & {
    __typename?: 'Issue';
    appSource: Scalars['String'];
    createdAt: Scalars['Time'];
    description?: Maybe<Scalars['String']>;
    externalLinks: Array<ExternalSystem>;
    id: Scalars['ID'];
    interactionEvents: Array<InteractionEvent>;
    mentionedByNotes: Array<Note>;
    priority?: Maybe<Scalars['String']>;
    source: DataSource;
    sourceOfTruth: DataSource;
    status: Scalars['String'];
    subject?: Maybe<Scalars['String']>;
    tags?: Maybe<Array<Maybe<Tag>>>;
    updatedAt: Scalars['Time'];
  };

export type IssueSummaryByStatus = {
  __typename?: 'IssueSummaryByStatus';
  count: Scalars['Int64'];
  status: Scalars['String'];
};

/**
 * Describes the relationship a Contact has with a Organization.
 * **A `return` object**
 */
export type JobRole = {
  __typename?: 'JobRole';
  appSource: Scalars['String'];
  company?: Maybe<Scalars['String']>;
  contact?: Maybe<Contact>;
  createdAt: Scalars['Time'];
  description?: Maybe<Scalars['String']>;
  endedAt?: Maybe<Scalars['Time']>;
  id: Scalars['ID'];
  /** The Contact's job title. */
  jobTitle?: Maybe<Scalars['String']>;
  /**
   * Organization associated with a Contact.
   * **Required.**
   */
  organization?: Maybe<Organization>;
  primary: Scalars['Boolean'];
  source: DataSource;
  sourceOfTruth: DataSource;
  startedAt?: Maybe<Scalars['Time']>;
  updatedAt: Scalars['Time'];
};

/**
 * Describes the relationship a Contact has with an Organization.
 * **A `create` object**
 */
export type JobRoleInput = {
  appSource?: InputMaybe<Scalars['String']>;
  company?: InputMaybe<Scalars['String']>;
  description?: InputMaybe<Scalars['String']>;
  endedAt?: InputMaybe<Scalars['Time']>;
  jobTitle?: InputMaybe<Scalars['String']>;
  organizationId?: InputMaybe<Scalars['ID']>;
  primary?: InputMaybe<Scalars['Boolean']>;
  startedAt?: InputMaybe<Scalars['Time']>;
};

export type JobRoleParticipant = {
  __typename?: 'JobRoleParticipant';
  jobRoleParticipant: JobRole;
  type?: Maybe<Scalars['String']>;
};

/**
 * Describes the relationship a Contact has with an Organization.
 * **A `create` object**
 */
export type JobRoleUpdateInput = {
  company?: InputMaybe<Scalars['String']>;
  description?: InputMaybe<Scalars['String']>;
  endedAt?: InputMaybe<Scalars['Time']>;
  id: Scalars['ID'];
  jobTitle?: InputMaybe<Scalars['String']>;
  organizationId?: InputMaybe<Scalars['ID']>;
  primary?: InputMaybe<Scalars['Boolean']>;
  startedAt?: InputMaybe<Scalars['Time']>;
};

export type LinkOrganizationsInput = {
  organizationId: Scalars['ID'];
  subOrganizationId: Scalars['ID'];
  type?: InputMaybe<Scalars['String']>;
};

export type LinkedOrganization = {
  __typename?: 'LinkedOrganization';
  organization: Organization;
  type?: Maybe<Scalars['String']>;
};

export type Location = Node &
  SourceFields & {
    __typename?: 'Location';
    address?: Maybe<Scalars['String']>;
    address2?: Maybe<Scalars['String']>;
    addressType?: Maybe<Scalars['String']>;
    appSource: Scalars['String'];
    commercial?: Maybe<Scalars['Boolean']>;
    country?: Maybe<Scalars['String']>;
    createdAt: Scalars['Time'];
    district?: Maybe<Scalars['String']>;
    houseNumber?: Maybe<Scalars['String']>;
    id: Scalars['ID'];
    latitude?: Maybe<Scalars['Float']>;
    locality?: Maybe<Scalars['String']>;
    longitude?: Maybe<Scalars['Float']>;
    name?: Maybe<Scalars['String']>;
    plusFour?: Maybe<Scalars['String']>;
    postalCode?: Maybe<Scalars['String']>;
    predirection?: Maybe<Scalars['String']>;
    rawAddress?: Maybe<Scalars['String']>;
    region?: Maybe<Scalars['String']>;
    source: DataSource;
    sourceOfTruth: DataSource;
    street?: Maybe<Scalars['String']>;
    timeZone?: Maybe<Scalars['String']>;
    updatedAt: Scalars['Time'];
    utcOffset?: Maybe<Scalars['Int64']>;
    zip?: Maybe<Scalars['String']>;
  };

export type LocationUpdateInput = {
  address?: InputMaybe<Scalars['String']>;
  address2?: InputMaybe<Scalars['String']>;
  addressType?: InputMaybe<Scalars['String']>;
  commercial?: InputMaybe<Scalars['Boolean']>;
  country?: InputMaybe<Scalars['String']>;
  district?: InputMaybe<Scalars['String']>;
  houseNumber?: InputMaybe<Scalars['String']>;
  id: Scalars['ID'];
  latitude?: InputMaybe<Scalars['Float']>;
  locality?: InputMaybe<Scalars['String']>;
  longitude?: InputMaybe<Scalars['Float']>;
  name?: InputMaybe<Scalars['String']>;
  plusFour?: InputMaybe<Scalars['String']>;
  postalCode?: InputMaybe<Scalars['String']>;
  predirection?: InputMaybe<Scalars['String']>;
  rawAddress?: InputMaybe<Scalars['String']>;
  region?: InputMaybe<Scalars['String']>;
  street?: InputMaybe<Scalars['String']>;
  timeZone?: InputMaybe<Scalars['String']>;
  utcOffset?: InputMaybe<Scalars['Int64']>;
  zip?: InputMaybe<Scalars['String']>;
};

export type LogEntry = {
  __typename?: 'LogEntry';
  appSource: Scalars['String'];
  content?: Maybe<Scalars['String']>;
  contentType?: Maybe<Scalars['String']>;
  createdAt: Scalars['Time'];
  createdBy?: Maybe<User>;
  externalLinks: Array<ExternalSystem>;
  id: Scalars['ID'];
  source: DataSource;
  sourceOfTruth: DataSource;
  startedAt: Scalars['Time'];
  tags: Array<Tag>;
  updatedAt: Scalars['Time'];
};

export type LogEntryInput = {
  appSource?: InputMaybe<Scalars['String']>;
  content?: InputMaybe<Scalars['String']>;
  contentType?: InputMaybe<Scalars['String']>;
  startedAt?: InputMaybe<Scalars['Time']>;
  tags?: InputMaybe<Array<TagIdOrNameInput>>;
};

export type LogEntryUpdateInput = {
  content?: InputMaybe<Scalars['String']>;
  contentType?: InputMaybe<Scalars['String']>;
  startedAt?: InputMaybe<Scalars['Time']>;
};

export enum Market {
  B2B = 'B2B',
  B2C = 'B2C',
  Marketplace = 'MARKETPLACE',
}

export type Meeting = Node & {
  __typename?: 'Meeting';
  agenda?: Maybe<Scalars['String']>;
  agendaContentType?: Maybe<Scalars['String']>;
  appSource: Scalars['String'];
  attendedBy: Array<MeetingParticipant>;
  conferenceUrl?: Maybe<Scalars['String']>;
  createdAt: Scalars['Time'];
  createdBy: Array<MeetingParticipant>;
  describedBy: Array<Analysis>;
  endedAt?: Maybe<Scalars['Time']>;
  events: Array<InteractionEvent>;
  externalSystem: Array<ExternalSystem>;
  id: Scalars['ID'];
  includes: Array<Attachment>;
  meetingExternalUrl?: Maybe<Scalars['String']>;
  name?: Maybe<Scalars['String']>;
  note: Array<Note>;
  recording?: Maybe<Attachment>;
  source: DataSource;
  sourceOfTruth: DataSource;
  startedAt?: Maybe<Scalars['Time']>;
  status: MeetingStatus;
  updatedAt: Scalars['Time'];
};

export type MeetingInput = {
  agenda?: InputMaybe<Scalars['String']>;
  agendaContentType?: InputMaybe<Scalars['String']>;
  appSource?: InputMaybe<Scalars['String']>;
  attendedBy?: InputMaybe<Array<MeetingParticipantInput>>;
  conferenceUrl?: InputMaybe<Scalars['String']>;
  createdAt?: InputMaybe<Scalars['Time']>;
  createdBy?: InputMaybe<Array<MeetingParticipantInput>>;
  endedAt?: InputMaybe<Scalars['Time']>;
  externalSystem?: InputMaybe<ExternalSystemReferenceInput>;
  meetingExternalUrl?: InputMaybe<Scalars['String']>;
  name?: InputMaybe<Scalars['String']>;
  note?: InputMaybe<NoteInput>;
  startedAt?: InputMaybe<Scalars['Time']>;
  status?: InputMaybe<MeetingStatus>;
};

export type MeetingParticipant =
  | ContactParticipant
  | OrganizationParticipant
  | UserParticipant;

export type MeetingParticipantInput = {
  contactId?: InputMaybe<Scalars['ID']>;
  organizationId?: InputMaybe<Scalars['ID']>;
  userId?: InputMaybe<Scalars['ID']>;
};

export enum MeetingStatus {
  Accepted = 'ACCEPTED',
  Canceled = 'CANCELED',
  Undefined = 'UNDEFINED',
}

export type MeetingUpdateInput = {
  agenda?: InputMaybe<Scalars['String']>;
  agendaContentType?: InputMaybe<Scalars['String']>;
  appSource?: InputMaybe<Scalars['String']>;
  conferenceUrl?: InputMaybe<Scalars['String']>;
  endedAt?: InputMaybe<Scalars['Time']>;
  externalSystem?: InputMaybe<ExternalSystemReferenceInput>;
  meetingExternalUrl?: InputMaybe<Scalars['String']>;
  name?: InputMaybe<Scalars['String']>;
  note?: InputMaybe<NoteUpdateInput>;
  startedAt?: InputMaybe<Scalars['Time']>;
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
  totalElements: Scalars['Int64'];
  /**
   * Total number of pages in the query response.
   * **Required.**
   */
  totalPages: Scalars['Int'];
};

export type MentionedEntity = Issue;

export type Mutation = {
  __typename?: 'Mutation';
  analysis_Create: Analysis;
  attachment_Create: Attachment;
  contact_AddNewLocation: Location;
  contact_AddOrganizationById: Contact;
  contact_AddSocial: Social;
  contact_AddTagById: Contact;
  contact_Archive: Result;
  contact_Create: Contact;
  contact_HardDelete: Result;
  contact_Merge: Contact;
  contact_RemoveLocation: Contact;
  contact_RemoveOrganizationById: Contact;
  contact_RemoveTagById: Contact;
  contact_RestoreFromArchive: Result;
  contact_Update: Contact;
  customFieldDeleteFromContactById: Result;
  customFieldDeleteFromContactByName: Result;
  customFieldDeleteFromFieldSetById: Result;
  customFieldMergeToContact: CustomField;
  customFieldMergeToFieldSet: CustomField;
  customFieldUpdateInContact: CustomField;
  customFieldUpdateInFieldSet: CustomField;
  customFieldsMergeAndUpdateInContact: Contact;
  customer_contact_Create: CustomerContact;
  customer_user_AddJobRole: CustomerUser;
  emailDelete: Result;
  emailMergeToContact: Email;
  emailMergeToOrganization: Email;
  emailMergeToUser: Email;
  emailRemoveFromContact: Result;
  emailRemoveFromContactById: Result;
  emailRemoveFromOrganization: Result;
  emailRemoveFromOrganizationById: Result;
  emailRemoveFromUser: Result;
  emailRemoveFromUserById: Result;
  emailUpdateInContact: Email;
  emailUpdateInOrganization: Email;
  emailUpdateInUser: Email;
  entityTemplateCreate: EntityTemplate;
  fieldSetDeleteFromContact: Result;
  fieldSetMergeToContact?: Maybe<FieldSet>;
  fieldSetUpdateInContact?: Maybe<FieldSet>;
  interactionEvent_Create: InteractionEvent;
  interactionEvent_LinkAttachment: InteractionEvent;
  interactionSession_Create: InteractionSession;
  interactionSession_LinkAttachment: InteractionSession;
  jobRole_Create: JobRole;
  jobRole_Delete: Result;
  jobRole_Update: JobRole;
  location_RemoveFromContact: Contact;
  location_RemoveFromOrganization: Organization;
  location_Update: Location;
  logEntry_AddTag: Scalars['ID'];
  logEntry_CreateForOrganization: Scalars['ID'];
  logEntry_RemoveTag: Scalars['ID'];
  logEntry_Update: Scalars['ID'];
  meeting_AddNewLocation: Location;
  meeting_AddNote: Meeting;
  meeting_Create: Meeting;
  meeting_LinkAttachment: Meeting;
  meeting_LinkAttendedBy: Meeting;
  meeting_LinkRecording: Meeting;
  meeting_UnlinkAttachment: Meeting;
  meeting_UnlinkAttendedBy: Meeting;
  meeting_UnlinkRecording: Meeting;
  meeting_Update: Meeting;
  note_CreateForContact: Note;
  note_CreateForOrganization: Note;
  note_Delete: Result;
  note_LinkAttachment: Note;
  note_UnlinkAttachment: Note;
  note_Update: Note;
  organization_AddNewLocation: Location;
  organization_AddRelationship: Organization;
  organization_AddSocial: Social;
  organization_AddSubsidiary: Organization;
  organization_Archive?: Maybe<Result>;
  organization_ArchiveAll?: Maybe<Result>;
  organization_Create: Organization;
  organization_Hide: Scalars['ID'];
  organization_HideAll?: Maybe<Result>;
  organization_Merge: Organization;
  organization_RemoveRelationship: Organization;
  organization_RemoveRelationshipStage: Organization;
  organization_RemoveSubsidiary: Organization;
  organization_SetOwner: Organization;
  organization_SetRelationshipStage: Organization;
  organization_Show: Scalars['ID'];
  organization_ShowAll?: Maybe<Result>;
  organization_UnsetOwner: Organization;
  organization_Update: Organization;
  organization_UpdateBillingDetails: Scalars['ID'];
  /** @deprecated Use organization_UpdateBillingDetails instead */
  organization_UpdateBillingDetailsAsync: Scalars['ID'];
  organization_UpdateRenewalForecast: Scalars['ID'];
  /** @deprecated Use organization_UpdateRenewalForecast instead */
  organization_UpdateRenewalForecastAsync: Scalars['ID'];
  organization_UpdateRenewalLikelihood: Scalars['ID'];
  /** @deprecated Use organization_UpdateRenewalLikelihood instead */
  organization_UpdateRenewalLikelihoodAsync: Scalars['ID'];
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
  player_Merge: Player;
  player_SetDefaultUser: Player;
  player_Update: Player;
  social_Remove: Result;
  social_Update: Social;
  tag_Create: Tag;
  tag_Delete?: Maybe<Result>;
  tag_Update?: Maybe<Tag>;
  tenant_Merge: Scalars['String'];
  user_AddRole: User;
  user_AddRoleInTenant: User;
  user_Create: User;
  user_CreateInTenant: User;
  user_Delete: Result;
  user_DeleteInTenant: Result;
  user_RemoveRole: User;
  user_RemoveRoleInTenant: User;
  user_Update: User;
  workspace_Merge: Result;
  workspace_MergeToTenant: Result;
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
  contactId: Scalars['ID'];
  input: SocialInput;
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
  mergedContactIds: Array<Scalars['ID']>;
  primaryContactId: Scalars['ID'];
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
  contactId: Scalars['ID'];
  id: Scalars['ID'];
};

export type MutationCustomFieldDeleteFromContactByNameArgs = {
  contactId: Scalars['ID'];
  fieldName: Scalars['String'];
};

export type MutationCustomFieldDeleteFromFieldSetByIdArgs = {
  contactId: Scalars['ID'];
  fieldSetId: Scalars['ID'];
  id: Scalars['ID'];
};

export type MutationCustomFieldMergeToContactArgs = {
  contactId: Scalars['ID'];
  input: CustomFieldInput;
};

export type MutationCustomFieldMergeToFieldSetArgs = {
  contactId: Scalars['ID'];
  fieldSetId: Scalars['ID'];
  input: CustomFieldInput;
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
  customFields?: InputMaybe<Array<CustomFieldInput>>;
  fieldSets?: InputMaybe<Array<FieldSetInput>>;
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
  contactId: Scalars['ID'];
  input: EmailInput;
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
  contactId: Scalars['ID'];
  id: Scalars['ID'];
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
  email: Scalars['String'];
  userId: Scalars['ID'];
};

export type MutationEmailRemoveFromUserByIdArgs = {
  id: Scalars['ID'];
  userId: Scalars['ID'];
};

export type MutationEmailUpdateInContactArgs = {
  contactId: Scalars['ID'];
  input: EmailUpdateInput;
};

export type MutationEmailUpdateInOrganizationArgs = {
  input: EmailUpdateInput;
  organizationId: Scalars['ID'];
};

export type MutationEmailUpdateInUserArgs = {
  input: EmailUpdateInput;
  userId: Scalars['ID'];
};

export type MutationEntityTemplateCreateArgs = {
  input: EntityTemplateInput;
};

export type MutationFieldSetDeleteFromContactArgs = {
  contactId: Scalars['ID'];
  id: Scalars['ID'];
};

export type MutationFieldSetMergeToContactArgs = {
  contactId: Scalars['ID'];
  input: FieldSetInput;
};

export type MutationFieldSetUpdateInContactArgs = {
  contactId: Scalars['ID'];
  input: FieldSetUpdateInput;
};

export type MutationInteractionEvent_CreateArgs = {
  event: InteractionEventInput;
};

export type MutationInteractionEvent_LinkAttachmentArgs = {
  attachmentId: Scalars['ID'];
  eventId: Scalars['ID'];
};

export type MutationInteractionSession_CreateArgs = {
  session: InteractionSessionInput;
};

export type MutationInteractionSession_LinkAttachmentArgs = {
  attachmentId: Scalars['ID'];
  sessionId: Scalars['ID'];
};

export type MutationJobRole_CreateArgs = {
  contactId: Scalars['ID'];
  input: JobRoleInput;
};

export type MutationJobRole_DeleteArgs = {
  contactId: Scalars['ID'];
  roleId: Scalars['ID'];
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
  attachmentId: Scalars['ID'];
  meetingId: Scalars['ID'];
};

export type MutationMeeting_LinkAttendedByArgs = {
  meetingId: Scalars['ID'];
  participant: MeetingParticipantInput;
};

export type MutationMeeting_LinkRecordingArgs = {
  attachmentId: Scalars['ID'];
  meetingId: Scalars['ID'];
};

export type MutationMeeting_UnlinkAttachmentArgs = {
  attachmentId: Scalars['ID'];
  meetingId: Scalars['ID'];
};

export type MutationMeeting_UnlinkAttendedByArgs = {
  meetingId: Scalars['ID'];
  participant: MeetingParticipantInput;
};

export type MutationMeeting_UnlinkRecordingArgs = {
  attachmentId: Scalars['ID'];
  meetingId: Scalars['ID'];
};

export type MutationMeeting_UpdateArgs = {
  meeting: MeetingUpdateInput;
  meetingId: Scalars['ID'];
};

export type MutationNote_CreateForContactArgs = {
  contactId: Scalars['ID'];
  input: NoteInput;
};

export type MutationNote_CreateForOrganizationArgs = {
  input: NoteInput;
  organizationId: Scalars['ID'];
};

export type MutationNote_DeleteArgs = {
  id: Scalars['ID'];
};

export type MutationNote_LinkAttachmentArgs = {
  attachmentId: Scalars['ID'];
  noteId: Scalars['ID'];
};

export type MutationNote_UnlinkAttachmentArgs = {
  attachmentId: Scalars['ID'];
  noteId: Scalars['ID'];
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
  mergedOrganizationIds: Array<Scalars['ID']>;
  primaryOrganizationId: Scalars['ID'];
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
  organizationId: Scalars['ID'];
  subsidiaryId: Scalars['ID'];
};

export type MutationOrganization_SetOwnerArgs = {
  organizationId: Scalars['ID'];
  userId: Scalars['ID'];
};

export type MutationOrganization_SetRelationshipStageArgs = {
  organizationId: Scalars['ID'];
  relationship: OrganizationRelationship;
  stage: Scalars['String'];
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
  contactId: Scalars['ID'];
  input: PhoneNumberInput;
};

export type MutationPhoneNumberMergeToOrganizationArgs = {
  input: PhoneNumberInput;
  organizationId: Scalars['ID'];
};

export type MutationPhoneNumberMergeToUserArgs = {
  input: PhoneNumberInput;
  userId: Scalars['ID'];
};

export type MutationPhoneNumberRemoveFromContactByE164Args = {
  contactId: Scalars['ID'];
  e164: Scalars['String'];
};

export type MutationPhoneNumberRemoveFromContactByIdArgs = {
  contactId: Scalars['ID'];
  id: Scalars['ID'];
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
  e164: Scalars['String'];
  userId: Scalars['ID'];
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
  input: PhoneNumberUpdateInput;
  userId: Scalars['ID'];
};

export type MutationPlayer_MergeArgs = {
  input: PlayerInput;
};

export type MutationPlayer_SetDefaultUserArgs = {
  id: Scalars['ID'];
  userId: Scalars['ID'];
};

export type MutationPlayer_UpdateArgs = {
  id: Scalars['ID'];
  update: PlayerUpdate;
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
  id: Scalars['ID'];
  role: Role;
};

export type MutationUser_AddRoleInTenantArgs = {
  id: Scalars['ID'];
  role: Role;
  tenant: Scalars['String'];
};

export type MutationUser_CreateArgs = {
  input: UserInput;
};

export type MutationUser_CreateInTenantArgs = {
  input: UserInput;
  tenant: Scalars['String'];
};

export type MutationUser_DeleteArgs = {
  id: Scalars['ID'];
};

export type MutationUser_DeleteInTenantArgs = {
  id: Scalars['ID'];
  tenant: Scalars['String'];
};

export type MutationUser_RemoveRoleArgs = {
  id: Scalars['ID'];
  role: Role;
};

export type MutationUser_RemoveRoleInTenantArgs = {
  id: Scalars['ID'];
  role: Role;
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
  __typename?: 'Note';
  appSource: Scalars['String'];
  content?: Maybe<Scalars['String']>;
  contentType?: Maybe<Scalars['String']>;
  createdAt: Scalars['Time'];
  createdBy?: Maybe<User>;
  id: Scalars['ID'];
  includes: Array<Attachment>;
  mentioned: Array<MentionedEntity>;
  noted: Array<NotedEntity>;
  source: DataSource;
  sourceOfTruth: DataSource;
  updatedAt: Scalars['Time'];
};

export type NoteInput = {
  appSource?: InputMaybe<Scalars['String']>;
  content?: InputMaybe<Scalars['String']>;
  contentType?: InputMaybe<Scalars['String']>;
};

export type NotePage = Pages & {
  __typename?: 'NotePage';
  content: Array<Note>;
  totalElements: Scalars['Int64'];
  totalPages: Scalars['Int'];
};

export type NoteUpdateInput = {
  content?: InputMaybe<Scalars['String']>;
  contentType?: InputMaybe<Scalars['String']>;
  id: Scalars['ID'];
};

export type NotedEntity = Contact | Organization;

export type OrgAccountDetails = {
  __typename?: 'OrgAccountDetails';
  billingDetails?: Maybe<BillingDetails>;
  renewalForecast?: Maybe<RenewalForecast>;
  renewalLikelihood?: Maybe<RenewalLikelihood>;
};

export type Organization = Node & {
  __typename?: 'Organization';
  accountDetails?: Maybe<OrgAccountDetails>;
  appSource: Scalars['String'];
  contacts: ContactsPage;
  createdAt: Scalars['Time'];
  customFields: Array<CustomField>;
  description?: Maybe<Scalars['String']>;
  domains: Array<Scalars['String']>;
  emails: Array<Email>;
  employees?: Maybe<Scalars['Int64']>;
  entityTemplate?: Maybe<EntityTemplate>;
  externalLinks: Array<ExternalSystem>;
  fieldSets: Array<FieldSet>;
  id: Scalars['ID'];
  industry?: Maybe<Scalars['String']>;
  industryGroup?: Maybe<Scalars['String']>;
  isPublic?: Maybe<Scalars['Boolean']>;
  issueSummaryByStatus: Array<IssueSummaryByStatus>;
  jobRoles: Array<JobRole>;
  lastFundingAmount?: Maybe<Scalars['String']>;
  lastFundingRound?: Maybe<FundingRound>;
  lastTouchPointAt?: Maybe<Scalars['Time']>;
  lastTouchPointTimelineEvent?: Maybe<TimelineEvent>;
  lastTouchPointTimelineEventId?: Maybe<Scalars['ID']>;
  locations: Array<Location>;
  market?: Maybe<Market>;
  name: Scalars['String'];
  note?: Maybe<Scalars['String']>;
  notes: NotePage;
  owner?: Maybe<User>;
  phoneNumbers: Array<PhoneNumber>;
  relationshipStages: Array<OrganizationRelationshipStage>;
  relationships: Array<OrganizationRelationship>;
  socials: Array<Social>;
  source: DataSource;
  sourceOfTruth: DataSource;
  subIndustry?: Maybe<Scalars['String']>;
  subsidiaries: Array<LinkedOrganization>;
  subsidiaryOf: Array<LinkedOrganization>;
  suggestedMergeTo: Array<SuggestedMergeOrganization>;
  tags?: Maybe<Array<Tag>>;
  targetAudience?: Maybe<Scalars['String']>;
  timelineEvents: Array<TimelineEvent>;
  timelineEventsTotalCount: Scalars['Int64'];
  updatedAt: Scalars['Time'];
  valueProposition?: Maybe<Scalars['String']>;
  website?: Maybe<Scalars['String']>;
};

export type OrganizationContactsArgs = {
  pagination?: InputMaybe<Pagination>;
  sort?: InputMaybe<Array<SortBy>>;
  where?: InputMaybe<Filter>;
};

export type OrganizationNotesArgs = {
  pagination?: InputMaybe<Pagination>;
};

export type OrganizationTimelineEventsArgs = {
  from?: InputMaybe<Scalars['Time']>;
  size: Scalars['Int'];
  timelineEventTypes?: InputMaybe<Array<TimelineEventType>>;
};

export type OrganizationTimelineEventsTotalCountArgs = {
  timelineEventTypes?: InputMaybe<Array<TimelineEventType>>;
};

export type OrganizationInput = {
  appSource?: InputMaybe<Scalars['String']>;
  customFields?: InputMaybe<Array<CustomFieldInput>>;
  description?: InputMaybe<Scalars['String']>;
  domains?: InputMaybe<Array<Scalars['String']>>;
  employees?: InputMaybe<Scalars['Int64']>;
  fieldSets?: InputMaybe<Array<FieldSetInput>>;
  industry?: InputMaybe<Scalars['String']>;
  industryGroup?: InputMaybe<Scalars['String']>;
  isPublic?: InputMaybe<Scalars['Boolean']>;
  market?: InputMaybe<Market>;
  /**
   * The name of the organization.
   * **Required.**
   */
  name: Scalars['String'];
  note?: InputMaybe<Scalars['String']>;
  subIndustry?: InputMaybe<Scalars['String']>;
  templateId?: InputMaybe<Scalars['ID']>;
  website?: InputMaybe<Scalars['String']>;
};

export type OrganizationPage = Pages & {
  __typename?: 'OrganizationPage';
  content: Array<Organization>;
  totalElements: Scalars['Int64'];
  totalPages: Scalars['Int'];
};

export type OrganizationParticipant = {
  __typename?: 'OrganizationParticipant';
  organizationParticipant: Organization;
  type?: Maybe<Scalars['String']>;
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
  __typename?: 'OrganizationRelationshipStage';
  relationship: OrganizationRelationship;
  stage?: Maybe<Scalars['String']>;
};

export type OrganizationUpdateInput = {
  description?: InputMaybe<Scalars['String']>;
  domains?: InputMaybe<Array<Scalars['String']>>;
  employees?: InputMaybe<Scalars['Int64']>;
  id: Scalars['ID'];
  industry?: InputMaybe<Scalars['String']>;
  industryGroup?: InputMaybe<Scalars['String']>;
  isPublic?: InputMaybe<Scalars['Boolean']>;
  lastFundingAmount?: InputMaybe<Scalars['String']>;
  lastFundingRound?: InputMaybe<FundingRound>;
  market?: InputMaybe<Market>;
  name: Scalars['String'];
  note?: InputMaybe<Scalars['String']>;
  /** Set to true when partial update is needed. Empty or missing fields will not be ignored. */
  patch?: InputMaybe<Scalars['Boolean']>;
  subIndustry?: InputMaybe<Scalars['String']>;
  targetAudience?: InputMaybe<Scalars['String']>;
  valueProposition?: InputMaybe<Scalars['String']>;
  website?: InputMaybe<Scalars['String']>;
};

export type PageView = Node &
  SourceFields & {
    __typename?: 'PageView';
    appSource: Scalars['String'];
    application: Scalars['String'];
    endedAt: Scalars['Time'];
    engagedTime: Scalars['Int64'];
    id: Scalars['ID'];
    orderInSession: Scalars['Int64'];
    pageTitle: Scalars['String'];
    pageUrl: Scalars['String'];
    sessionId: Scalars['ID'];
    source: DataSource;
    sourceOfTruth: DataSource;
    startedAt: Scalars['Time'];
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
  totalElements: Scalars['Int64'];
  /**
   * The total number of pages included in the query response.
   * **Required.**
   */
  totalPages: Scalars['Int'];
};

/** If provided as part of the request, results will be filtered down to the `page` and `limit` specified. */
export type Pagination = {
  /**
   * The maximum number of results in the response.
   * **Required.**
   */
  limit: Scalars['Int'];
  /**
   * The results page to return in the response.
   * **Required.**
   */
  page: Scalars['Int'];
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
  appSource?: Maybe<Scalars['String']>;
  contacts: Array<Contact>;
  country?: Maybe<Country>;
  createdAt: Scalars['Time'];
  /** The phone number in e164 format.  */
  e164?: Maybe<Scalars['String']>;
  /**
   * The unique ID associated with the phone number.
   * **Required**
   */
  id: Scalars['ID'];
  /** Defines the type of phone number. */
  label?: Maybe<PhoneNumberLabel>;
  organizations: Array<Organization>;
  /**
   * Determines if the phone number is primary or not.
   * **Required**
   */
  primary: Scalars['Boolean'];
  rawPhoneNumber?: Maybe<Scalars['String']>;
  source: DataSource;
  updatedAt: Scalars['Time'];
  users: Array<User>;
  validated?: Maybe<Scalars['Boolean']>;
};

/**
 * Describes a phone number associated with a `Contact` in customerOS.
 * **A `create` object.**
 */
export type PhoneNumberInput = {
  countryCodeA2?: InputMaybe<Scalars['String']>;
  /** Defines the type of phone number. */
  label?: InputMaybe<PhoneNumberLabel>;
  /**
   * The phone number in e164 format.
   * **Required**
   */
  phoneNumber: Scalars['String'];
  /**
   * Determines if the phone number is primary or not.
   * **Required**
   */
  primary?: InputMaybe<Scalars['Boolean']>;
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
  type?: Maybe<Scalars['String']>;
};

/**
 * Describes a phone number associated with a `Contact` in customerOS.
 * **An `update` object.**
 */
export type PhoneNumberUpdateInput = {
  countryCodeA2?: InputMaybe<Scalars['String']>;
  /**
   * The unique ID associated with the phone number.
   * **Required**
   */
  id: Scalars['ID'];
  /** Defines the type of phone number. */
  label?: InputMaybe<PhoneNumberLabel>;
  phoneNumber?: InputMaybe<Scalars['String']>;
  /**
   * Determines if the phone number is primary or not.
   * **Required**
   */
  primary?: InputMaybe<Scalars['Boolean']>;
};

export type Player = {
  __typename?: 'Player';
  appSource: Scalars['String'];
  authId: Scalars['String'];
  createdAt: Scalars['Time'];
  id: Scalars['ID'];
  identityId?: Maybe<Scalars['String']>;
  provider: Scalars['String'];
  source: DataSource;
  sourceOfTruth: DataSource;
  updatedAt: Scalars['Time'];
  users: Array<PlayerUser>;
};

export type PlayerInput = {
  appSource?: InputMaybe<Scalars['String']>;
  authId: Scalars['String'];
  identityId?: InputMaybe<Scalars['String']>;
  provider: Scalars['String'];
};

export type PlayerUpdate = {
  appSource?: InputMaybe<Scalars['String']>;
  identityId?: InputMaybe<Scalars['String']>;
};

export type PlayerUser = {
  __typename?: 'PlayerUser';
  default: Scalars['Boolean'];
  tenant: Scalars['String'];
  user: User;
};

export type Query = {
  __typename?: 'Query';
  analysis: Analysis;
  attachment: Attachment;
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
  /** sort.By available options: CONTACT, EMAIL, ORGANIZATION, LOCATION, RELATIONSHIP, STAGE */
  dashboardView_Contacts?: Maybe<ContactsPage>;
  /** sort.By available options: ORGANIZATION, DOMAIN, LOCATION, OWNER, RELATIONSHIP, LAST_TOUCHPOINT, HEALTH_INDICATOR_ORDER, HEALTH_INDICATOR_NAME, FORECAST_AMOUNT, RENEWAL_LIKELIHOOD, RENEWAL_CYCLE_NEXT */
  dashboardView_Organizations?: Maybe<OrganizationPage>;
  email: Email;
  entityTemplates: Array<EntityTemplate>;
  externalMeetings: MeetingsPage;
  gcli_Search: Array<GCliItem>;
  global_Cache: GlobalCache;
  interactionEvent: InteractionEvent;
  interactionEvent_ByEventIdentifier: InteractionEvent;
  interactionSession: InteractionSession;
  interactionSession_BySessionIdentifier: InteractionSession;
  issue: Issue;
  logEntry: LogEntry;
  meeting: Meeting;
  organization?: Maybe<Organization>;
  organization_DistinctOwners: Array<User>;
  organizations: OrganizationPage;
  phoneNumber: PhoneNumber;
  player_ByAuthIdProvider: Player;
  player_GetUsers: Array<PlayerUser>;
  tags: Array<Tag>;
  tenant: Scalars['String'];
  tenant_ByEmail?: Maybe<Scalars['String']>;
  tenant_ByWorkspace?: Maybe<Scalars['String']>;
  timelineEvents: Array<TimelineEvent>;
  user: User;
  user_ByEmail: User;
  users: UserPage;
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
  pagination?: InputMaybe<Pagination>;
  sort?: InputMaybe<Array<SortBy>>;
  where?: InputMaybe<Filter>;
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
  externalId?: InputMaybe<Scalars['ID']>;
  externalSystemId: Scalars['String'];
  pagination?: InputMaybe<Pagination>;
  sort?: InputMaybe<Array<SortBy>>;
  where?: InputMaybe<Filter>;
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
  pagination?: InputMaybe<Pagination>;
  sort?: InputMaybe<Array<SortBy>>;
  where?: InputMaybe<Filter>;
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
  pagination?: InputMaybe<Pagination>;
  sort?: InputMaybe<Array<SortBy>>;
  where?: InputMaybe<Filter>;
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
  __typename?: 'RenewalForecast';
  amount?: Maybe<Scalars['Float']>;
  comment?: Maybe<Scalars['String']>;
  potentialAmount?: Maybe<Scalars['Float']>;
  updatedAt?: Maybe<Scalars['Time']>;
  updatedBy?: Maybe<User>;
  updatedById?: Maybe<Scalars['String']>;
};

export type RenewalForecastInput = {
  amount?: InputMaybe<Scalars['Float']>;
  comment?: InputMaybe<Scalars['String']>;
  id: Scalars['ID'];
};

export type RenewalLikelihood = {
  __typename?: 'RenewalLikelihood';
  comment?: Maybe<Scalars['String']>;
  previousProbability?: Maybe<RenewalLikelihoodProbability>;
  probability?: Maybe<RenewalLikelihoodProbability>;
  updatedAt?: Maybe<Scalars['Time']>;
  updatedBy?: Maybe<User>;
  updatedById?: Maybe<Scalars['String']>;
};

export type RenewalLikelihoodInput = {
  comment?: InputMaybe<Scalars['String']>;
  id: Scalars['ID'];
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
    __typename?: 'Social';
    appSource: Scalars['String'];
    createdAt: Scalars['Time'];
    id: Scalars['ID'];
    platformName?: Maybe<Scalars['String']>;
    source: DataSource;
    sourceOfTruth: DataSource;
    updatedAt: Scalars['Time'];
    url: Scalars['String'];
  };

export type SocialInput = {
  appSource?: InputMaybe<Scalars['String']>;
  platformName?: InputMaybe<Scalars['String']>;
  url: Scalars['String'];
};

export type SocialUpdateInput = {
  id: Scalars['ID'];
  platformName?: InputMaybe<Scalars['String']>;
  url: Scalars['String'];
};

export type SortBy = {
  by: Scalars['String'];
  caseSensitive?: InputMaybe<Scalars['Boolean']>;
  direction?: SortingDirection;
};

export enum SortingDirection {
  Asc = 'ASC',
  Desc = 'DESC',
}

export type SourceFields = {
  appSource: Scalars['String'];
  id: Scalars['ID'];
  source: DataSource;
  sourceOfTruth: DataSource;
};

export type State = {
  __typename?: 'State';
  code: Scalars['String'];
  country: Country;
  id: Scalars['ID'];
  name: Scalars['String'];
};

export type SuggestedMergeOrganization = {
  __typename?: 'SuggestedMergeOrganization';
  confidence?: Maybe<Scalars['Float']>;
  organization: Organization;
  suggestedAt?: Maybe<Scalars['Time']>;
  suggestedBy?: Maybe<Scalars['String']>;
};

export type Tag = {
  __typename?: 'Tag';
  appSource: Scalars['String'];
  createdAt: Scalars['Time'];
  id: Scalars['ID'];
  name: Scalars['String'];
  source: DataSource;
  updatedAt: Scalars['Time'];
};

export type TagIdOrNameInput = {
  id?: InputMaybe<Scalars['ID']>;
  name?: InputMaybe<Scalars['String']>;
};

export type TagInput = {
  appSource?: InputMaybe<Scalars['String']>;
  name: Scalars['String'];
};

export type TagUpdateInput = {
  id: Scalars['ID'];
  name: Scalars['String'];
};

export type TenantBillableInfo = {
  __typename?: 'TenantBillableInfo';
  greylistedContacts: Scalars['Int64'];
  greylistedOrganizations: Scalars['Int64'];
  whitelistedContacts: Scalars['Int64'];
  whitelistedOrganizations: Scalars['Int64'];
};

export type TenantInput = {
  appSource?: InputMaybe<Scalars['String']>;
  name: Scalars['String'];
};

export type TimeRange = {
  /**
   * The start time of the time range.
   * **Required.**
   */
  from: Scalars['Time'];
  /**
   * The end time of the time range.
   * **Required.**
   */
  to: Scalars['Time'];
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
  __typename?: 'User';
  appSource: Scalars['String'];
  calendars: Array<Calendar>;
  /**
   * Timestamp of user creation.
   * **Required**
   */
  createdAt: Scalars['Time'];
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
  /**
   * The unique ID associated with the customerOS user.
   * **Required**
   */
  id: Scalars['ID'];
  internal: Scalars['Boolean'];
  jobRoles: Array<JobRole>;
  /**
   * The last name of the customerOS user.
   * **Required**
   */
  lastName: Scalars['String'];
  phoneNumbers: Array<PhoneNumber>;
  player: Player;
  profilePhotoUrl?: Maybe<Scalars['String']>;
  roles: Array<Role>;
  source: DataSource;
  sourceOfTruth: DataSource;
  timezone?: Maybe<Scalars['String']>;
  updatedAt: Scalars['Time'];
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
  appSource?: InputMaybe<Scalars['String']>;
  /**
   * The email address of the customerOS user.
   * **Required**
   */
  email: EmailInput;
  /**
   * The first name of the customerOS user.
   * **Required**
   */
  firstName: Scalars['String'];
  /**
   * The Job Roles of the user.
   * **Optional**
   */
  jobRoles?: InputMaybe<Array<JobRoleInput>>;
  /**
   * The last name of the customerOS user.
   * **Required**
   */
  lastName: Scalars['String'];
  /**
   * Player to associate with the user with. If the person does not exist, it will be created.
   * **Required**
   */
  player: PlayerInput;
  timezone?: InputMaybe<Scalars['String']>;
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
  totalElements: Scalars['Int64'];
  /**
   * Total number of pages in the query response.
   * **Required.**
   */
  totalPages: Scalars['Int'];
};

export type UserParticipant = {
  __typename?: 'UserParticipant';
  type?: Maybe<Scalars['String']>;
  userParticipant: User;
};

export type UserUpdateInput = {
  /**
   * The first name of the customerOS user.
   * **Required**
   */
  firstName: Scalars['String'];
  id: Scalars['ID'];
  /**
   * The last name of the customerOS user.
   * **Required**
   */
  lastName: Scalars['String'];
  timezone?: InputMaybe<Scalars['String']>;
};

export type Workspace = {
  __typename?: 'Workspace';
  appSource: Scalars['String'];
  createdAt: Scalars['Time'];
  id: Scalars['ID'];
  name: Scalars['String'];
  provider: Scalars['String'];
  source: DataSource;
  sourceOfTruth: DataSource;
  updatedAt: Scalars['Time'];
};

export type WorkspaceInput = {
  appSource?: InputMaybe<Scalars['String']>;
  name: Scalars['String'];
  provider: Scalars['String'];
};
