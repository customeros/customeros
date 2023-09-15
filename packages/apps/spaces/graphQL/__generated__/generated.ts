import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
const defaultOptions = {} as const;
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
  RenewalLikelihoodUpdated = 'RENEWAL_LIKELIHOOD_UPDATED'
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
  Google = 'GOOGLE'
}

export enum ComparisonOperator {
  Contains = 'CONTAINS',
  Eq = 'EQ',
  StartsWith = 'STARTS_WITH'
}

/**
 * A contact represents an individual in customerOS.
 * **A `response` object.**
 */
export type Contact = ExtensibleEntity & Node & {
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
  Text = 'TEXT'
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
  Text = 'TEXT'
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
  ZendeskSupport = 'ZENDESK_SUPPORT'
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
  Work = 'WORK'
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
  Organization = 'ORGANIZATION'
}

export type EntityTemplateInput = {
  customFieldTemplateInputs?: InputMaybe<Array<CustomFieldTemplateInput>>;
  extends?: InputMaybe<EntityTemplateExtension>;
  fieldSetTemplateInputs?: InputMaybe<Array<FieldSetTemplateInput>>;
  name: Scalars['String'];
};

export enum EntityType {
  Contact = 'Contact',
  Organization = 'Organization'
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
  ZendeskSupport = 'ZENDESK_SUPPORT'
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
  SeriesF = 'SERIES_F'
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
  State = 'STATE'
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
  State = 'STATE'
}

export type GlobalCache = {
  __typename?: 'GlobalCache';
  gCliCache: Array<GCliItem>;
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

export type InteractionEventParticipant = ContactParticipant | EmailParticipant | JobRoleParticipant | OrganizationParticipant | PhoneNumberParticipant | UserParticipant;

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

export type InteractionSessionParticipant = ContactParticipant | EmailParticipant | PhoneNumberParticipant | UserParticipant;

export type InteractionSessionParticipantInput = {
  contactID?: InputMaybe<Scalars['ID']>;
  email?: InputMaybe<Scalars['String']>;
  phoneNumber?: InputMaybe<Scalars['String']>;
  type?: InputMaybe<Scalars['String']>;
  userID?: InputMaybe<Scalars['ID']>;
};

export type Issue = Node & SourceFields & {
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

export type Location = Node & SourceFields & {
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

export enum Market {
  B2B = 'B2B',
  B2C = 'B2C',
  Marketplace = 'MARKETPLACE'
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

export type MeetingParticipant = ContactParticipant | OrganizationParticipant | UserParticipant;

export type MeetingParticipantInput = {
  contactId?: InputMaybe<Scalars['ID']>;
  organizationId?: InputMaybe<Scalars['ID']>;
  userId?: InputMaybe<Scalars['ID']>;
};

export enum MeetingStatus {
  Accepted = 'ACCEPTED',
  Canceled = 'CANCELED',
  Undefined = 'UNDEFINED'
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
  organization_UpdateBillingDetailsAsync: Scalars['ID'];
  organization_UpdateRenewalForecastAsync: Scalars['ID'];
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


export type MutationOrganization_UpdateBillingDetailsAsyncArgs = {
  input: BillingDetailsInput;
};


export type MutationOrganization_UpdateRenewalForecastAsyncArgs = {
  input: RenewalForecastInput;
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
  Vendor = 'VENDOR'
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

export type PageView = Node & SourceFields & {
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
  Ms = 'MS'
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
  Work = 'WORK'
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
  Weekly = 'WEEKLY'
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
  Zero = 'ZERO'
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
  User = 'USER'
}

export type Social = Node & SourceFields & {
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
  Desc = 'DESC'
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

export type TimelineEvent = Action | Analysis | InteractionEvent | InteractionSession | Issue | Meeting | Note | PageView;

export enum TimelineEventType {
  Action = 'ACTION',
  Analysis = 'ANALYSIS',
  InteractionEvent = 'INTERACTION_EVENT',
  InteractionSession = 'INTERACTION_SESSION',
  Issue = 'ISSUE',
  Meeting = 'MEETING',
  Note = 'NOTE',
  PageView = 'PAGE_VIEW'
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

export type CreateTagMutationVariables = Exact<{
  input: TagInput;
}>;


export type CreateTagMutation = { __typename?: 'Mutation', tag_Create: { __typename?: 'Tag', id: string, name: string, createdAt: any, updatedAt: any, source: DataSource } };

export type DeleteTagMutationVariables = Exact<{
  id: Scalars['ID'];
}>;


export type DeleteTagMutation = { __typename?: 'Mutation', tag_Delete?: { __typename?: 'Result', result: boolean } | null };

export type GetTagsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetTagsQuery = { __typename?: 'Query', tags: Array<{ __typename?: 'Tag', id: string, name: string }> };

export type UpdateTagMutationVariables = Exact<{
  input: TagUpdateInput;
}>;


export type UpdateTagMutation = { __typename?: 'Mutation', tag_Update?: { __typename?: 'Tag', id: string, name: string } | null };

export type AddEmailToContactMutationVariables = Exact<{
  contactId: Scalars['ID'];
  input: EmailInput;
}>;


export type AddEmailToContactMutation = { __typename?: 'Mutation', emailMergeToContact: { __typename?: 'Email', label?: EmailLabel | null, id: string, primary: boolean, email?: string | null } };

export type AddLocationToContactMutationVariables = Exact<{
  contactId: Scalars['ID'];
}>;


export type AddLocationToContactMutation = { __typename?: 'Mutation', contact_AddNewLocation: { __typename?: 'Location', id: string } };

export type AddPhoneToContactMutationVariables = Exact<{
  contactId: Scalars['ID'];
  input: PhoneNumberInput;
}>;


export type AddPhoneToContactMutation = { __typename?: 'Mutation', phoneNumberMergeToContact: { __typename?: 'PhoneNumber', label?: PhoneNumberLabel | null, id: string, primary: boolean, e164?: string | null, rawPhoneNumber?: string | null } };

export type AddTagToContactMutationVariables = Exact<{
  input: ContactTagInput;
}>;


export type AddTagToContactMutation = { __typename?: 'Mutation', contact_AddTagById: { __typename?: 'Contact', id: string, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null } };

export type ArchiveContactMutationVariables = Exact<{
  id: Scalars['ID'];
}>;


export type ArchiveContactMutation = { __typename?: 'Mutation', contact_Archive: { __typename?: 'Result', result: boolean } };

export type AttachOrganizationToContactMutationVariables = Exact<{
  input: ContactOrganizationInput;
}>;


export type AttachOrganizationToContactMutation = { __typename?: 'Mutation', contact_AddOrganizationById: { __typename?: 'Contact', id: string, source: DataSource, firstName?: string | null, lastName?: string | null, name?: string | null, jobRoles: Array<{ __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string, organization?: { __typename?: 'Organization', id: string, name: string } | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null } };

export type CreateContactMutationVariables = Exact<{
  input: ContactInput;
}>;


export type CreateContactMutation = { __typename?: 'Mutation', contact_Create: { __typename?: 'Contact', id: string, source: DataSource, firstName?: string | null, lastName?: string | null, name?: string | null, jobRoles: Array<{ __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string, organization?: { __typename?: 'Organization', id: string, name: string } | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null, emails: Array<{ __typename?: 'Email', label?: EmailLabel | null, id: string, primary: boolean, email?: string | null, emailValidationDetails: { __typename?: 'EmailValidationDetails', isReachable?: string | null, isValidSyntax?: boolean | null, canConnectSmtp?: boolean | null, acceptsMail?: boolean | null, hasFullInbox?: boolean | null, isCatchAll?: boolean | null, isDeliverable?: boolean | null, validated?: boolean | null, isDisabled?: boolean | null } }>, phoneNumbers: Array<{ __typename?: 'PhoneNumber', label?: PhoneNumberLabel | null, id: string, primary: boolean, e164?: string | null, rawPhoneNumber?: string | null }> } };

export type CreateContactJobRoleMutationVariables = Exact<{
  contactId: Scalars['ID'];
  input: JobRoleInput;
}>;


export type CreateContactJobRoleMutation = { __typename?: 'Mutation', jobRole_Create: { __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string, organization?: { __typename?: 'Organization', id: string, name: string } | null } };

export type CreateContactNoteMutationVariables = Exact<{
  contactId: Scalars['ID'];
  input: NoteInput;
}>;


export type CreateContactNoteMutation = { __typename?: 'Mutation', note_CreateForContact: { __typename?: 'Note', id: string, createdAt: any, updatedAt: any, source: DataSource, sourceOfTruth: DataSource, appSource: string, createdBy?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null, includes: Array<{ __typename?: 'Attachment', id: string, name: string, mimeType: string, extension: string, size: any }> } };

export type CreatePhoneCallInteractionEventMutationVariables = Exact<{
  contactId?: InputMaybe<Scalars['ID']>;
  sentBy?: InputMaybe<Scalars['String']>;
  content?: InputMaybe<Scalars['String']>;
  contentType?: InputMaybe<Scalars['String']>;
}>;


export type CreatePhoneCallInteractionEventMutation = { __typename?: 'Mutation', interactionEvent_Create: { __typename?: 'InteractionEvent', id: string, createdAt: any, channel?: string | null, content?: string | null, contentType?: string | null, interactionSession?: { __typename?: 'InteractionSession', name: string } | null, issue?: { __typename?: 'Issue', externalLinks: Array<{ __typename?: 'ExternalSystem', type: ExternalSystemType, externalId?: string | null, externalUrl?: string | null }> } | null, sentBy: Array<{ __typename: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null } } | { __typename: 'EmailParticipant', emailParticipant: { __typename?: 'Email', email?: string | null, id: string, contacts: Array<{ __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null, emails: Array<{ __typename?: 'Email', email?: string | null }> }>, users: Array<{ __typename?: 'User', id: string, firstName: string, lastName: string }>, organizations: Array<{ __typename?: 'Organization', id: string, name: string }> } } | { __typename?: 'JobRoleParticipant' } | { __typename?: 'OrganizationParticipant' } | { __typename: 'PhoneNumberParticipant', phoneNumberParticipant: { __typename?: 'PhoneNumber', e164?: string | null, id: string, contacts: Array<{ __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null }>, users: Array<{ __typename?: 'User', id: string, firstName: string, lastName: string }>, organizations: Array<{ __typename?: 'Organization', id: string, name: string }> } } | { __typename: 'UserParticipant', userParticipant: { __typename?: 'User', id: string, firstName: string, lastName: string } }>, sentTo: Array<{ __typename: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', name?: string | null, id: string, firstName?: string | null, lastName?: string | null } } | { __typename: 'EmailParticipant', emailParticipant: { __typename?: 'Email', email?: string | null, id: string, contacts: Array<{ __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null }>, users: Array<{ __typename?: 'User', id: string, firstName: string, lastName: string }>, organizations: Array<{ __typename?: 'Organization', id: string, name: string }> } } | { __typename: 'JobRoleParticipant' } | { __typename: 'OrganizationParticipant' } | { __typename: 'PhoneNumberParticipant', phoneNumberParticipant: { __typename?: 'PhoneNumber', e164?: string | null, id: string, contacts: Array<{ __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null }>, users: Array<{ __typename?: 'User', id: string, firstName: string, lastName: string }>, organizations: Array<{ __typename?: 'Organization', id: string, name: string }> } } | { __typename: 'UserParticipant', type?: string | null, userParticipant: { __typename?: 'User', id: string, firstName: string, lastName: string } }> } };

export type RemoveContactJobRoleMutationVariables = Exact<{
  contactId: Scalars['ID'];
  roleId: Scalars['ID'];
}>;


export type RemoveContactJobRoleMutation = { __typename?: 'Mutation', jobRole_Delete: { __typename?: 'Result', result: boolean } };

export type GetContactQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetContactQuery = { __typename?: 'Query', contact?: { __typename?: 'Contact', source: DataSource, id: string, firstName?: string | null, lastName?: string | null, name?: string | null, owner?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null, jobRoles: Array<{ __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string, organization?: { __typename?: 'Organization', id: string, name: string } | null }>, locations: Array<{ __typename?: 'Location', rawAddress?: string | null, id: string, name?: string | null, country?: string | null, region?: string | null, locality?: string | null, zip?: string | null, street?: string | null, postalCode?: string | null, houseNumber?: string | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null, emails: Array<{ __typename?: 'Email', label?: EmailLabel | null, id: string, primary: boolean, email?: string | null, emailValidationDetails: { __typename?: 'EmailValidationDetails', isReachable?: string | null, isValidSyntax?: boolean | null, canConnectSmtp?: boolean | null, acceptsMail?: boolean | null, hasFullInbox?: boolean | null, isCatchAll?: boolean | null, isDeliverable?: boolean | null, validated?: boolean | null, isDisabled?: boolean | null } }>, phoneNumbers: Array<{ __typename?: 'PhoneNumber', label?: PhoneNumberLabel | null, id: string, primary: boolean, e164?: string | null, rawPhoneNumber?: string | null }> } | null };

export type GetContactCommunicationChannelsQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetContactCommunicationChannelsQuery = { __typename?: 'Query', contact?: { __typename?: 'Contact', firstName?: string | null, lastName?: string | null, name?: string | null, id: string, emails: Array<{ __typename?: 'Email', label?: EmailLabel | null, id: string, primary: boolean, email?: string | null, emailValidationDetails: { __typename?: 'EmailValidationDetails', isReachable?: string | null, isValidSyntax?: boolean | null, canConnectSmtp?: boolean | null, acceptsMail?: boolean | null, hasFullInbox?: boolean | null, isCatchAll?: boolean | null, isDeliverable?: boolean | null, validated?: boolean | null, isDisabled?: boolean | null } }>, phoneNumbers: Array<{ __typename?: 'PhoneNumber', label?: PhoneNumberLabel | null, id: string, primary: boolean, e164?: string | null, rawPhoneNumber?: string | null }> } | null };

export type GetContactListQueryVariables = Exact<{
  pagination: Pagination;
  where?: InputMaybe<Filter>;
  sort?: InputMaybe<Array<SortBy> | SortBy>;
}>;


export type GetContactListQuery = { __typename?: 'Query', contacts: { __typename?: 'ContactsPage', totalElements: any, content: Array<{ __typename?: 'Contact', id: string, firstName?: string | null, lastName?: string | null, name?: string | null, emails: Array<{ __typename?: 'Email', id: string, email?: string | null }> }> } };

export type GetContactLocationsQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetContactLocationsQuery = { __typename?: 'Query', contact?: { __typename?: 'Contact', locations: Array<{ __typename?: 'Location', rawAddress?: string | null, id: string, name?: string | null, country?: string | null, region?: string | null, locality?: string | null, zip?: string | null, street?: string | null, postalCode?: string | null, houseNumber?: string | null }> } | null };

export type GetContactMentionSuggestionsQueryVariables = Exact<{
  pagination: Pagination;
  where?: InputMaybe<Filter>;
  sort?: InputMaybe<Array<SortBy> | SortBy>;
}>;


export type GetContactMentionSuggestionsQuery = { __typename?: 'Query', contacts: { __typename?: 'ContactsPage', content: Array<{ __typename?: 'Contact', id: string, firstName?: string | null, lastName?: string | null }> } };

export type GetContactNameByEmailQueryVariables = Exact<{
  email: Scalars['String'];
}>;


export type GetContactNameByEmailQuery = { __typename?: 'Query', contact_ByEmail: { __typename?: 'Contact', id: string, firstName?: string | null, lastName?: string | null, name?: string | null } };

export type GetContactNameByIdQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetContactNameByIdQuery = { __typename?: 'Query', contact?: { __typename?: 'Contact', id: string, firstName?: string | null, lastName?: string | null, name?: string | null } | null };

export type GetContactNameByPhoneNumberQueryVariables = Exact<{
  e164: Scalars['String'];
}>;


export type GetContactNameByPhoneNumberQuery = { __typename?: 'Query', contact_ByPhone: { __typename?: 'Contact', id: string, firstName?: string | null, lastName?: string | null, name?: string | null } };

export type GetContactNotesQueryVariables = Exact<{
  id: Scalars['ID'];
  pagination?: InputMaybe<Pagination>;
}>;


export type GetContactNotesQuery = { __typename?: 'Query', contact?: { __typename?: 'Contact', notes: { __typename?: 'NotePage', content: Array<{ __typename?: 'Note', id: string, createdAt: any, updatedAt: any, source: DataSource, sourceOfTruth: DataSource, appSource: string, createdBy?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null, includes: Array<{ __typename?: 'Attachment', id: string, name: string, mimeType: string, extension: string, size: any }> }> } } | null };

export type GetContactPersonalDetailsQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetContactPersonalDetailsQuery = { __typename?: 'Query', contact?: { __typename?: 'Contact', id: string, source: DataSource, firstName?: string | null, lastName?: string | null, name?: string | null, owner?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null, jobRoles: Array<{ __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string, organization?: { __typename?: 'Organization', id: string, name: string } | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null } | null };

export type GetContactPersonalDetailsWithOrganizationsQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetContactPersonalDetailsWithOrganizationsQuery = { __typename?: 'Query', contact?: { __typename?: 'Contact', id: string, source: DataSource, firstName?: string | null, lastName?: string | null, name?: string | null, organizations: { __typename?: 'OrganizationPage', content: Array<{ __typename?: 'Organization', id: string, name: string }> }, jobRoles: Array<{ __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string, organization?: { __typename?: 'Organization', id: string, name: string } | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null } | null };

export type GetContactTagsQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetContactTagsQuery = { __typename?: 'Query', contact?: { __typename?: 'Contact', id: string, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null } | null };

export type MergeContactsMutationVariables = Exact<{
  primaryContactId: Scalars['ID'];
  mergedContactIds: Array<Scalars['ID']> | Scalars['ID'];
}>;


export type MergeContactsMutation = { __typename?: 'Mutation', contact_Merge: { __typename?: 'Contact', id: string, source: DataSource, firstName?: string | null, lastName?: string | null, name?: string | null, jobRoles: Array<{ __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string, organization?: { __typename?: 'Organization', id: string, name: string } | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null } };

export type RemoveEmailFromContactMutationVariables = Exact<{
  contactId: Scalars['ID'];
  id: Scalars['ID'];
}>;


export type RemoveEmailFromContactMutation = { __typename?: 'Mutation', emailRemoveFromContactById: { __typename?: 'Result', result: boolean } };

export type RemoveLocationFromContactMutationVariables = Exact<{
  locationId: Scalars['ID'];
  contactId: Scalars['ID'];
}>;


export type RemoveLocationFromContactMutation = { __typename?: 'Mutation', location_RemoveFromContact: { __typename?: 'Contact', id: string, locations: Array<{ __typename?: 'Location', id: string, name?: string | null, country?: string | null, region?: string | null, locality?: string | null, zip?: string | null, street?: string | null, postalCode?: string | null, houseNumber?: string | null }> } };

export type RemoveOrganizationFromContactMutationVariables = Exact<{
  input: ContactOrganizationInput;
}>;


export type RemoveOrganizationFromContactMutation = { __typename?: 'Mutation', contact_RemoveOrganizationById: { __typename?: 'Contact', id: string, source: DataSource, firstName?: string | null, lastName?: string | null, name?: string | null, jobRoles: Array<{ __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string, organization?: { __typename?: 'Organization', id: string, name: string } | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null } };

export type RemovePhoneNumberFromContactMutationVariables = Exact<{
  contactId: Scalars['ID'];
  id: Scalars['ID'];
}>;


export type RemovePhoneNumberFromContactMutation = { __typename?: 'Mutation', phoneNumberRemoveFromContactById: { __typename?: 'Result', result: boolean } };

export type RemoveTagFromContactMutationVariables = Exact<{
  input: ContactTagInput;
}>;


export type RemoveTagFromContactMutation = { __typename?: 'Mutation', contact_RemoveTagById: { __typename?: 'Contact', id: string } };

export type UpdateContactEmailMutationVariables = Exact<{
  contactId: Scalars['ID'];
  input: EmailUpdateInput;
}>;


export type UpdateContactEmailMutation = { __typename?: 'Mutation', emailUpdateInContact: { __typename?: 'Email', primary: boolean, label?: EmailLabel | null, email?: string | null, id: string } };

export type UpdateJobRoleMutationVariables = Exact<{
  contactId: Scalars['ID'];
  input: JobRoleUpdateInput;
}>;


export type UpdateJobRoleMutation = { __typename?: 'Mutation', jobRole_Update: { __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string, organization?: { __typename?: 'Organization', id: string, name: string } | null } };

export type UpdateContactPersonalDetailsMutationVariables = Exact<{
  input: ContactUpdateInput;
}>;


export type UpdateContactPersonalDetailsMutation = { __typename?: 'Mutation', contact_Update: { __typename?: 'Contact', id: string, title?: string | null, firstName?: string | null, lastName?: string | null } };

export type UpdateContactPhoneNumberMutationVariables = Exact<{
  contactId: Scalars['ID'];
  input: PhoneNumberUpdateInput;
}>;


export type UpdateContactPhoneNumberMutation = { __typename?: 'Mutation', phoneNumberUpdateInContact: { __typename?: 'PhoneNumber', label?: PhoneNumberLabel | null, id: string, primary: boolean, e164?: string | null, rawPhoneNumber?: string | null } };

export type DashboardView_ContactsQueryVariables = Exact<{
  pagination: Pagination;
  where?: InputMaybe<Filter>;
  sort?: InputMaybe<SortBy>;
}>;


export type DashboardView_ContactsQuery = { __typename?: 'Query', dashboardView_Contacts?: { __typename?: 'ContactsPage', totalElements: any, content: Array<{ __typename?: 'Contact', id: string, source: DataSource, firstName?: string | null, lastName?: string | null, name?: string | null, locations: Array<{ __typename?: 'Location', rawAddress?: string | null, id: string, name?: string | null, country?: string | null, region?: string | null, locality?: string | null, zip?: string | null, street?: string | null, postalCode?: string | null, houseNumber?: string | null }>, jobRoles: Array<{ __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string, organization?: { __typename?: 'Organization', id: string, name: string } | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null, emails: Array<{ __typename?: 'Email', label?: EmailLabel | null, id: string, primary: boolean, email?: string | null, emailValidationDetails: { __typename?: 'EmailValidationDetails', isReachable?: string | null, isValidSyntax?: boolean | null, canConnectSmtp?: boolean | null, acceptsMail?: boolean | null, hasFullInbox?: boolean | null, isCatchAll?: boolean | null, isDeliverable?: boolean | null, validated?: boolean | null, isDisabled?: boolean | null } }>, phoneNumbers: Array<{ __typename?: 'PhoneNumber', label?: PhoneNumberLabel | null, id: string, primary: boolean, e164?: string | null, rawPhoneNumber?: string | null }> }> } | null };

export type DashboardView_OrganizationsQueryVariables = Exact<{
  pagination: Pagination;
  where?: InputMaybe<Filter>;
  sort?: InputMaybe<SortBy>;
}>;


export type DashboardView_OrganizationsQuery = { __typename?: 'Query', dashboardView_Organizations?: { __typename?: 'OrganizationPage', totalElements: any, content: Array<{ __typename?: 'Organization', id: string, name: string, description?: string | null, industry?: string | null, website?: string | null, domains: Array<string>, lastTouchPointTimelineEventId?: string | null, lastTouchPointAt?: any | null, subsidiaryOf: Array<{ __typename?: 'LinkedOrganization', organization: { __typename?: 'Organization', id: string, name: string } }>, owner?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null, accountDetails?: { __typename?: 'OrgAccountDetails', renewalForecast?: { __typename?: 'RenewalForecast', amount?: number | null, potentialAmount?: number | null, comment?: string | null, updatedAt?: any | null, updatedById?: string | null, updatedBy?: { __typename?: 'User', id: string, firstName: string, lastName: string, emails?: Array<{ __typename?: 'Email', email?: string | null }> | null } | null } | null, renewalLikelihood?: { __typename?: 'RenewalLikelihood', probability?: RenewalLikelihoodProbability | null, previousProbability?: RenewalLikelihoodProbability | null, comment?: string | null, updatedById?: string | null, updatedAt?: any | null, updatedBy?: { __typename?: 'User', id: string, firstName: string, lastName: string, emails?: Array<{ __typename?: 'Email', email?: string | null }> | null } | null } | null, billingDetails?: { __typename?: 'BillingDetails', renewalCycle?: RenewalCycle | null, frequency?: RenewalCycle | null, amount?: number | null, renewalCycleNext?: any | null } | null } | null, locations: Array<{ __typename?: 'Location', rawAddress?: string | null, id: string, name?: string | null, country?: string | null, region?: string | null, locality?: string | null, zip?: string | null, street?: string | null, postalCode?: string | null, houseNumber?: string | null }>, relationshipStages: Array<{ __typename?: 'OrganizationRelationshipStage', relationship: OrganizationRelationship, stage?: string | null }>, lastTouchPointTimelineEvent?: { __typename?: 'Action', id: string, actionType: ActionType, createdAt: any, source: DataSource, createdBy?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null } | { __typename?: 'Analysis', id: string } | { __typename?: 'InteractionEvent', id: string, channel?: string | null, eventType?: string | null, externalLinks: Array<{ __typename?: 'ExternalSystem', type: ExternalSystemType }>, sentBy: Array<{ __typename: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null } } | { __typename: 'EmailParticipant', type?: string | null, emailParticipant: { __typename?: 'Email', id: string, email?: string | null, rawEmail?: string | null } } | { __typename: 'JobRoleParticipant', jobRoleParticipant: { __typename?: 'JobRole', contact?: { __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null } | null } } | { __typename: 'OrganizationParticipant' } | { __typename: 'PhoneNumberParticipant' } | { __typename: 'UserParticipant', userParticipant: { __typename?: 'User', id: string, firstName: string, lastName: string } }> } | { __typename?: 'InteractionSession' } | { __typename?: 'Issue', id: string } | { __typename?: 'Meeting', id: string, name?: string | null, attendedBy: Array<{ __typename: 'ContactParticipant' } | { __typename: 'OrganizationParticipant' } | { __typename: 'UserParticipant' }> } | { __typename?: 'Note', id: string, createdBy?: { __typename?: 'User', firstName: string, lastName: string } | null } | { __typename?: 'PageView', id: string } | null }> } | null };

export type LocationBaseDetailsFragment = { __typename?: 'Location', id: string, name?: string | null, country?: string | null, region?: string | null, locality?: string | null, zip?: string | null, street?: string | null, postalCode?: string | null, houseNumber?: string | null };

export type LocationTotalFragment = { __typename?: 'Location', id: string, name?: string | null, createdAt: any, updatedAt: any, source: DataSource, appSource: string, country?: string | null, region?: string | null, locality?: string | null, address?: string | null, address2?: string | null, zip?: string | null, addressType?: string | null, houseNumber?: string | null, postalCode?: string | null, plusFour?: string | null, commercial?: boolean | null, predirection?: string | null, district?: string | null, street?: string | null, rawAddress?: string | null, latitude?: number | null, longitude?: number | null };

export type JobRoleFragment = { __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string };

export type NoteContentFragment = { __typename?: 'Note', id: string, createdAt: any, updatedAt: any, source: DataSource, sourceOfTruth: DataSource, appSource: string, createdBy?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null, includes: Array<{ __typename?: 'Attachment', id: string, name: string, mimeType: string, extension: string, size: any }> };

export type TagFragment = { __typename?: 'Tag', id: string, name: string };

export type EmailFragment = { __typename?: 'Email', id: string, primary: boolean, email?: string | null };

export type EmailWithValidationFragment = { __typename?: 'Email', id: string, primary: boolean, email?: string | null, emailValidationDetails: { __typename?: 'EmailValidationDetails', isReachable?: string | null, isValidSyntax?: boolean | null, canConnectSmtp?: boolean | null, acceptsMail?: boolean | null, hasFullInbox?: boolean | null, isCatchAll?: boolean | null, isDeliverable?: boolean | null, validated?: boolean | null, isDisabled?: boolean | null } };

export type PhoneNumberFragment = { __typename?: 'PhoneNumber', id: string, primary: boolean, e164?: string | null, rawPhoneNumber?: string | null };

export type InteractionSessionFragmentFragment = { __typename?: 'InteractionSession', id: string, startedAt: any, name: string, status: string, type?: string | null, events: Array<{ __typename?: 'InteractionEvent', content?: string | null, contentType?: string | null }> };

export type InteractionEventFragmentFragment = { __typename?: 'InteractionEvent', id: string, createdAt: any, channel?: string | null, content?: string | null, contentType?: string | null, interactionSession?: { __typename?: 'InteractionSession', name: string } | null, issue?: { __typename?: 'Issue', externalLinks: Array<{ __typename?: 'ExternalSystem', type: ExternalSystemType, externalId?: string | null, externalUrl?: string | null }> } | null, sentBy: Array<{ __typename: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null } } | { __typename: 'EmailParticipant', emailParticipant: { __typename?: 'Email', email?: string | null, id: string, contacts: Array<{ __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null, emails: Array<{ __typename?: 'Email', email?: string | null }> }>, users: Array<{ __typename?: 'User', id: string, firstName: string, lastName: string }>, organizations: Array<{ __typename?: 'Organization', id: string, name: string }> } } | { __typename?: 'JobRoleParticipant' } | { __typename?: 'OrganizationParticipant' } | { __typename: 'PhoneNumberParticipant', phoneNumberParticipant: { __typename?: 'PhoneNumber', e164?: string | null, id: string, contacts: Array<{ __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null }>, users: Array<{ __typename?: 'User', id: string, firstName: string, lastName: string }>, organizations: Array<{ __typename?: 'Organization', id: string, name: string }> } } | { __typename: 'UserParticipant', userParticipant: { __typename?: 'User', id: string, firstName: string, lastName: string } }>, sentTo: Array<{ __typename: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', name?: string | null, id: string, firstName?: string | null, lastName?: string | null } } | { __typename: 'EmailParticipant', emailParticipant: { __typename?: 'Email', email?: string | null, id: string, contacts: Array<{ __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null }>, users: Array<{ __typename?: 'User', id: string, firstName: string, lastName: string }>, organizations: Array<{ __typename?: 'Organization', id: string, name: string }> } } | { __typename: 'JobRoleParticipant' } | { __typename: 'OrganizationParticipant' } | { __typename: 'PhoneNumberParticipant', phoneNumberParticipant: { __typename?: 'PhoneNumber', e164?: string | null, id: string, contacts: Array<{ __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null }>, users: Array<{ __typename?: 'User', id: string, firstName: string, lastName: string }>, organizations: Array<{ __typename?: 'Organization', id: string, name: string }> } } | { __typename: 'UserParticipant', type?: string | null, userParticipant: { __typename?: 'User', id: string, firstName: string, lastName: string } }> };

export type MeetingTimelineEventFragmentFragment = { __typename?: 'Meeting', id: string, createdAt: any, agenda?: string | null, agendaContentType?: string | null, conferenceUrl?: string | null, meetingStartedAt?: any | null, meetingEndedAt?: any | null, attendedBy: Array<{ __typename?: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', id: string, firstName?: string | null, lastName?: string | null, name?: string | null } } | { __typename?: 'OrganizationParticipant' } | { __typename?: 'UserParticipant', userParticipant: { __typename?: 'User', id: string, firstName: string, lastName: string } }>, meetingCreatedBy: Array<{ __typename?: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', id: string } } | { __typename?: 'OrganizationParticipant' } | { __typename?: 'UserParticipant', userParticipant: { __typename?: 'User', id: string } }>, describedBy: Array<{ __typename?: 'Analysis', id: string, analysisType?: string | null, content?: string | null, contentType?: string | null }>, events: Array<{ __typename?: 'InteractionEvent', id: string, createdAt: any, channel?: string | null, content?: string | null, contentType?: string | null, sentBy: Array<{ __typename?: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', id: string, firstName?: string | null, lastName?: string | null, name?: string | null } } | { __typename?: 'EmailParticipant' } | { __typename?: 'JobRoleParticipant' } | { __typename?: 'OrganizationParticipant' } | { __typename?: 'PhoneNumberParticipant' } | { __typename?: 'UserParticipant', userParticipant: { __typename?: 'User', id: string, firstName: string, lastName: string } }>, sentTo: Array<{ __typename?: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', id: string, firstName?: string | null, lastName?: string | null, name?: string | null } } | { __typename?: 'EmailParticipant' } | { __typename?: 'JobRoleParticipant' } | { __typename?: 'OrganizationParticipant' } | { __typename?: 'PhoneNumberParticipant' } | { __typename?: 'UserParticipant', userParticipant: { __typename?: 'User', id: string, firstName: string, lastName: string } }>, includes: Array<{ __typename?: 'Attachment', id: string, name: string, mimeType: string, extension: string, size: any }> }>, recording?: { __typename?: 'Attachment', id: string } | null, includes: Array<{ __typename?: 'Attachment', id: string, name: string, mimeType: string, extension: string, size: any }>, note: Array<{ __typename?: 'Note', id: string, appSource: string }> };

export type ContactNameFragmentFragment = { __typename?: 'Contact', firstName?: string | null, lastName?: string | null, name?: string | null };

export type OrganizationBaseDetailsFragment = { __typename?: 'Organization', id: string, name: string, industry?: string | null };

export type ContactPersonalDetailsFragment = { __typename?: 'Contact', id: string, source: DataSource, firstName?: string | null, lastName?: string | null, name?: string | null, jobRoles: Array<{ __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string, organization?: { __typename?: 'Organization', id: string, name: string } | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null };

export type ContactCommunicationChannelsDetailsFragment = { __typename?: 'Contact', id: string, emails: Array<{ __typename?: 'Email', label?: EmailLabel | null, id: string, primary: boolean, email?: string | null, emailValidationDetails: { __typename?: 'EmailValidationDetails', isReachable?: string | null, isValidSyntax?: boolean | null, canConnectSmtp?: boolean | null, acceptsMail?: boolean | null, hasFullInbox?: boolean | null, isCatchAll?: boolean | null, isDeliverable?: boolean | null, validated?: boolean | null, isDisabled?: boolean | null } }>, phoneNumbers: Array<{ __typename?: 'PhoneNumber', label?: PhoneNumberLabel | null, id: string, primary: boolean, e164?: string | null, rawPhoneNumber?: string | null }> };

export type OrganizationDetailsFragment = { __typename?: 'Organization', id: string, name: string, description?: string | null, source: DataSource, industry?: string | null, website?: string | null, domains: Array<string>, updatedAt: any, emails: Array<{ __typename?: 'Email', id: string, primary: boolean, email?: string | null }>, locations: Array<{ __typename?: 'Location', rawAddress?: string | null, id: string, name?: string | null, country?: string | null, region?: string | null, locality?: string | null, zip?: string | null, street?: string | null, postalCode?: string | null, houseNumber?: string | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null };

export type OrganizationContactsFragment = { __typename?: 'Organization', contacts: { __typename?: 'ContactsPage', content: Array<{ __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null, jobRoles: Array<{ __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string }>, emails: Array<{ __typename?: 'Email', label?: EmailLabel | null, id: string, primary: boolean, email?: string | null, emailValidationDetails: { __typename?: 'EmailValidationDetails', isReachable?: string | null, isValidSyntax?: boolean | null, canConnectSmtp?: boolean | null, acceptsMail?: boolean | null, hasFullInbox?: boolean | null, isCatchAll?: boolean | null, isDeliverable?: boolean | null, validated?: boolean | null, isDisabled?: boolean | null } }>, phoneNumbers: Array<{ __typename?: 'PhoneNumber', label?: PhoneNumberLabel | null, id: string, primary: boolean, e164?: string | null, rawPhoneNumber?: string | null }> }> } };

export type GCliSearchQueryVariables = Exact<{
  limit?: InputMaybe<Scalars['Int']>;
  keyword: Scalars['String'];
}>;


export type GCliSearchQuery = { __typename?: 'Query', gcli_Search: Array<{ __typename?: 'GCliItem', id: string, type: GCliSearchResultType, display: string, data?: Array<{ __typename?: 'GCliAttributeKeyValuePair', key: string, value: string, display?: string | null }> | null }> };

export type Global_CacheQueryVariables = Exact<{ [key: string]: never; }>;


export type Global_CacheQuery = { __typename?: 'Query', global_Cache: { __typename?: 'GlobalCache', isOwner: boolean, user: { __typename?: 'User', id: string, firstName: string, lastName: string, emails?: Array<{ __typename?: 'Email', email?: string | null, rawEmail?: string | null, primary: boolean }> | null }, gCliCache: Array<{ __typename?: 'GCliItem', id: string, type: GCliSearchResultType, display: string, data?: Array<{ __typename?: 'GCliAttributeKeyValuePair', key: string, value: string, display?: string | null }> | null }> } };

export type AddEmailToOrganizationMutationVariables = Exact<{
  organizationId: Scalars['ID'];
  input: EmailInput;
}>;


export type AddEmailToOrganizationMutation = { __typename?: 'Mutation', emailMergeToOrganization: { __typename?: 'Email', label?: EmailLabel | null, id: string, primary: boolean, email?: string | null } };

export type AddLocationToOrganizationMutationVariables = Exact<{
  organzationId: Scalars['ID'];
}>;


export type AddLocationToOrganizationMutation = { __typename?: 'Mutation', organization_AddNewLocation: { __typename?: 'Location', id: string } };

export type AddPhoneToOrganizationMutationVariables = Exact<{
  organizationId: Scalars['ID'];
  input: PhoneNumberInput;
}>;


export type AddPhoneToOrganizationMutation = { __typename?: 'Mutation', phoneNumberMergeToOrganization: { __typename?: 'PhoneNumber', label?: PhoneNumberLabel | null, id: string, primary: boolean, e164?: string | null, rawPhoneNumber?: string | null } };

export type AddRelationshipToOrganizationMutationVariables = Exact<{
  organizationId: Scalars['ID'];
  relationship: OrganizationRelationship;
}>;


export type AddRelationshipToOrganizationMutation = { __typename?: 'Mutation', organization_AddRelationship: { __typename?: 'Organization', id: string } };

export type AddOrganizationSubsidiaryMutationVariables = Exact<{
  input: LinkOrganizationsInput;
}>;


export type AddOrganizationSubsidiaryMutation = { __typename?: 'Mutation', organization_AddSubsidiary: { __typename?: 'Organization', id: string, subsidiaries: Array<{ __typename?: 'LinkedOrganization', organization: { __typename?: 'Organization', id: string, name: string } }> } };

export type CreateOrganizationMutationVariables = Exact<{
  input: OrganizationInput;
}>;


export type CreateOrganizationMutation = { __typename?: 'Mutation', organization_Create: { __typename?: 'Organization', id: string, name: string } };

export type CreateOrganizationNoteMutationVariables = Exact<{
  organizationId: Scalars['ID'];
  input: NoteInput;
}>;


export type CreateOrganizationNoteMutation = { __typename?: 'Mutation', note_CreateForOrganization: { __typename?: 'Note', id: string, createdAt: any, updatedAt: any, source: DataSource, sourceOfTruth: DataSource, appSource: string, createdBy?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null, includes: Array<{ __typename?: 'Attachment', id: string, name: string, mimeType: string, extension: string, size: any }> } };

export type GetOrganizationQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetOrganizationQuery = { __typename?: 'Query', organization?: { __typename?: 'Organization', industry?: string | null, industryGroup?: string | null, subIndustry?: string | null, id: string, name: string, description?: string | null, source: DataSource, website?: string | null, domains: Array<string>, updatedAt: any, owner?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null, subsidiaryOf: Array<{ __typename?: 'LinkedOrganization', organization: { __typename?: 'Organization', id: string, name: string } }>, subsidiaries: Array<{ __typename?: 'LinkedOrganization', organization: { __typename?: 'Organization', name: string, id: string } }>, emails: Array<{ __typename?: 'Email', id: string, email?: string | null, primary: boolean, label?: EmailLabel | null }>, phoneNumbers: Array<{ __typename?: 'PhoneNumber', id: string, e164?: string | null, rawPhoneNumber?: string | null, label?: PhoneNumberLabel | null }>, customFields: Array<{ __typename?: 'CustomField', id: string, name: string, datatype: CustomFieldDataType, value: any, template?: { __typename?: 'CustomFieldTemplate', type: CustomFieldTemplateType } | null }>, locations: Array<{ __typename?: 'Location', rawAddress?: string | null, id: string, name?: string | null, country?: string | null, region?: string | null, locality?: string | null, zip?: string | null, street?: string | null, postalCode?: string | null, houseNumber?: string | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null, contacts: { __typename?: 'ContactsPage', content: Array<{ __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null, jobRoles: Array<{ __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string }>, emails: Array<{ __typename?: 'Email', label?: EmailLabel | null, id: string, primary: boolean, email?: string | null, emailValidationDetails: { __typename?: 'EmailValidationDetails', isReachable?: string | null, isValidSyntax?: boolean | null, canConnectSmtp?: boolean | null, acceptsMail?: boolean | null, hasFullInbox?: boolean | null, isCatchAll?: boolean | null, isDeliverable?: boolean | null, validated?: boolean | null, isDisabled?: boolean | null } }>, phoneNumbers: Array<{ __typename?: 'PhoneNumber', label?: PhoneNumberLabel | null, id: string, primary: boolean, e164?: string | null, rawPhoneNumber?: string | null }> }> } } | null };

export type GetOrganizationCommunicationChannelsQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetOrganizationCommunicationChannelsQuery = { __typename?: 'Query', organization?: { __typename?: 'Organization', id: string, name: string, emails: Array<{ __typename?: 'Email', id: string, email?: string | null, primary: boolean, label?: EmailLabel | null }>, phoneNumbers: Array<{ __typename?: 'PhoneNumber', id: string, e164?: string | null, rawPhoneNumber?: string | null, label?: PhoneNumberLabel | null }> } | null };

export type GetOrganizationContactsQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetOrganizationContactsQuery = { __typename?: 'Query', organization?: { __typename?: 'Organization', contacts: { __typename?: 'ContactsPage', content: Array<{ __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null, jobRoles: Array<{ __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string }>, emails: Array<{ __typename?: 'Email', label?: EmailLabel | null, id: string, primary: boolean, email?: string | null, emailValidationDetails: { __typename?: 'EmailValidationDetails', isReachable?: string | null, isValidSyntax?: boolean | null, canConnectSmtp?: boolean | null, acceptsMail?: boolean | null, hasFullInbox?: boolean | null, isCatchAll?: boolean | null, isDeliverable?: boolean | null, validated?: boolean | null, isDisabled?: boolean | null } }>, phoneNumbers: Array<{ __typename?: 'PhoneNumber', label?: PhoneNumberLabel | null, id: string, primary: boolean, e164?: string | null, rawPhoneNumber?: string | null }> }> } } | null };

export type GetOrganizationCustomFieldsQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetOrganizationCustomFieldsQuery = { __typename?: 'Query', organization?: { __typename?: 'Organization', customFields: Array<{ __typename?: 'CustomField', id: string, name: string, datatype: CustomFieldDataType, value: any, template?: { __typename?: 'CustomFieldTemplate', type: CustomFieldTemplateType } | null }> } | null };

export type GetOrganizationDetailsQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetOrganizationDetailsQuery = { __typename?: 'Query', organization?: { __typename?: 'Organization', id: string, name: string, description?: string | null, source: DataSource, industry?: string | null, website?: string | null, domains: Array<string>, updatedAt: any, subsidiaryOf: Array<{ __typename?: 'LinkedOrganization', organization: { __typename?: 'Organization', id: string, name: string } }>, emails: Array<{ __typename?: 'Email', id: string, primary: boolean, email?: string | null }>, locations: Array<{ __typename?: 'Location', rawAddress?: string | null, id: string, name?: string | null, country?: string | null, region?: string | null, locality?: string | null, zip?: string | null, street?: string | null, postalCode?: string | null, houseNumber?: string | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null } | null };

export type GetOrganizationLocationsQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetOrganizationLocationsQuery = { __typename?: 'Query', organization?: { __typename?: 'Organization', locations: Array<{ __typename?: 'Location', rawAddress?: string | null, id: string, name?: string | null, country?: string | null, region?: string | null, locality?: string | null, zip?: string | null, street?: string | null, postalCode?: string | null, houseNumber?: string | null }> } | null };

export type GetOrganizationMentionSuggestionsQueryVariables = Exact<{
  pagination: Pagination;
  where?: InputMaybe<Filter>;
  sort?: InputMaybe<Array<SortBy> | SortBy>;
}>;


export type GetOrganizationMentionSuggestionsQuery = { __typename?: 'Query', organizations: { __typename?: 'OrganizationPage', content: Array<{ __typename?: 'Organization', id: string, name: string }> } };

export type GetOrganizationNameQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetOrganizationNameQuery = { __typename?: 'Query', organization?: { __typename?: 'Organization', id: string, name: string } | null };

export type GetOrganizationNotesQueryVariables = Exact<{
  id: Scalars['ID'];
  pagination?: InputMaybe<Pagination>;
}>;


export type GetOrganizationNotesQuery = { __typename?: 'Query', organization?: { __typename?: 'Organization', notes: { __typename?: 'NotePage', content: Array<{ __typename?: 'Note', id: string, createdAt: any, updatedAt: any, source: DataSource, sourceOfTruth: DataSource, appSource: string, createdBy?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null, includes: Array<{ __typename?: 'Attachment', id: string, name: string, mimeType: string, extension: string, size: any }> }> } } | null };

export type GetOrganizationOwnerQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetOrganizationOwnerQuery = { __typename?: 'Query', organization?: { __typename?: 'Organization', id: string, owner?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null } | null };

export type GetOrganizationSubsidiariesQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetOrganizationSubsidiariesQuery = { __typename?: 'Query', organization?: { __typename?: 'Organization', subsidiaries: Array<{ __typename?: 'LinkedOrganization', organization: { __typename?: 'Organization', name: string, id: string } }> } | null };

export type GetOrganizationTableDataQueryVariables = Exact<{
  pagination?: InputMaybe<Pagination>;
  where?: InputMaybe<Filter>;
  sort?: InputMaybe<Array<SortBy> | SortBy>;
}>;


export type GetOrganizationTableDataQuery = { __typename?: 'Query', organizations: { __typename?: 'OrganizationPage', totalElements: any, totalPages: number, content: Array<{ __typename?: 'Organization', id: string, name: string, industry?: string | null, locations: Array<{ __typename?: 'Location', id: string, name?: string | null, country?: string | null, region?: string | null, locality?: string | null, zip?: string | null, street?: string | null, postalCode?: string | null, houseNumber?: string | null }>, subsidiaryOf: Array<{ __typename?: 'LinkedOrganization', type?: string | null, organization: { __typename?: 'Organization', name: string } }> }> } };

export type GetOrganizationsOptionsQueryVariables = Exact<{
  pagination?: InputMaybe<Pagination>;
}>;


export type GetOrganizationsOptionsQuery = { __typename?: 'Query', organizations: { __typename?: 'OrganizationPage', content: Array<{ __typename?: 'Organization', id: string, name: string }> } };

export type HideOrganizationsMutationVariables = Exact<{
  ids: Array<Scalars['ID']> | Scalars['ID'];
}>;


export type HideOrganizationsMutation = { __typename?: 'Mutation', organization_HideAll?: { __typename?: 'Result', result: boolean } | null };

export type MergeOrganizationsMutationVariables = Exact<{
  primaryOrganizationId: Scalars['ID'];
  mergedOrganizationIds: Array<Scalars['ID']> | Scalars['ID'];
}>;


export type MergeOrganizationsMutation = { __typename?: 'Mutation', organization_Merge: { __typename?: 'Organization', id: string, name: string, description?: string | null, source: DataSource, industry?: string | null, website?: string | null, domains: Array<string>, updatedAt: any, emails: Array<{ __typename?: 'Email', id: string, primary: boolean, email?: string | null }>, locations: Array<{ __typename?: 'Location', rawAddress?: string | null, id: string, name?: string | null, country?: string | null, region?: string | null, locality?: string | null, zip?: string | null, street?: string | null, postalCode?: string | null, houseNumber?: string | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string }> | null } };

export type RemoveEmailFromOrganizationMutationVariables = Exact<{
  organizationId: Scalars['ID'];
  id: Scalars['ID'];
}>;


export type RemoveEmailFromOrganizationMutation = { __typename?: 'Mutation', emailRemoveFromOrganizationById: { __typename?: 'Result', result: boolean } };

export type RemoveLocationFromOrganizationMutationVariables = Exact<{
  locationId: Scalars['ID'];
  organizationId: Scalars['ID'];
}>;


export type RemoveLocationFromOrganizationMutation = { __typename?: 'Mutation', location_RemoveFromOrganization: { __typename?: 'Organization', id: string, locations: Array<{ __typename?: 'Location', id: string, name?: string | null, country?: string | null, region?: string | null, locality?: string | null, zip?: string | null, street?: string | null, postalCode?: string | null, houseNumber?: string | null }> } };

export type RemoveOrganizationOwnerMutationVariables = Exact<{
  organizationId: Scalars['ID'];
}>;


export type RemoveOrganizationOwnerMutation = { __typename?: 'Mutation', organization_UnsetOwner: { __typename?: 'Organization', id: string, owner?: { __typename?: 'User', id: string } | null } };

export type RemoveOrganizationRelationshipMutationVariables = Exact<{
  organizationId: Scalars['ID'];
  relationship: OrganizationRelationship;
}>;


export type RemoveOrganizationRelationshipMutation = { __typename?: 'Mutation', organization_RemoveRelationship: { __typename?: 'Organization', id: string } };

export type RemovePhoneNumberFromOrganizationMutationVariables = Exact<{
  organizationId: Scalars['ID'];
  id: Scalars['ID'];
}>;


export type RemovePhoneNumberFromOrganizationMutation = { __typename?: 'Mutation', phoneNumberRemoveFromOrganizationById: { __typename?: 'Result', result: boolean } };

export type RemoveStageFromOrganizationRelationshipMutationVariables = Exact<{
  organizationId: Scalars['ID'];
  relationship: OrganizationRelationship;
}>;


export type RemoveStageFromOrganizationRelationshipMutation = { __typename?: 'Mutation', organization_RemoveRelationshipStage: { __typename?: 'Organization', id: string } };

export type RemoveOrganizationSubsidiaryMutationVariables = Exact<{
  organizationId: Scalars['ID'];
  subsidiaryId: Scalars['ID'];
}>;


export type RemoveOrganizationSubsidiaryMutation = { __typename?: 'Mutation', organization_RemoveSubsidiary: { __typename?: 'Organization', id: string, subsidiaries: Array<{ __typename?: 'LinkedOrganization', organization: { __typename?: 'Organization', id: string, name: string } }> } };

export type SetStageToOrganizationRelationshipMutationVariables = Exact<{
  organizationId: Scalars['ID'];
  relationship: OrganizationRelationship;
  stage: Scalars['String'];
}>;


export type SetStageToOrganizationRelationshipMutation = { __typename?: 'Mutation', organization_SetRelationshipStage: { __typename?: 'Organization', id: string } };

export type UpdateOrganizationDescriptionMutationVariables = Exact<{
  input: OrganizationUpdateInput;
}>;


export type UpdateOrganizationDescriptionMutation = { __typename?: 'Mutation', organization_Update: { __typename?: 'Organization', id: string, description?: string | null } };

export type UpdateOrganizationEmailMutationVariables = Exact<{
  organizationId: Scalars['ID'];
  input: EmailUpdateInput;
}>;


export type UpdateOrganizationEmailMutation = { __typename?: 'Mutation', emailUpdateInOrganization: { __typename?: 'Email', primary: boolean, label?: EmailLabel | null, id: string, email?: string | null } };

export type UpdateOrganizationIndustryMutationVariables = Exact<{
  input: OrganizationUpdateInput;
}>;


export type UpdateOrganizationIndustryMutation = { __typename?: 'Mutation', organization_Update: { __typename?: 'Organization', id: string, industry?: string | null } };

export type UpdateOrganizationNameMutationVariables = Exact<{
  input: OrganizationUpdateInput;
}>;


export type UpdateOrganizationNameMutation = { __typename?: 'Mutation', organization_Update: { __typename?: 'Organization', id: string, name: string } };

export type UpdateOrganizationOwnerMutationVariables = Exact<{
  organizationId: Scalars['ID'];
  userId: Scalars['ID'];
}>;


export type UpdateOrganizationOwnerMutation = { __typename?: 'Mutation', organization_SetOwner: { __typename?: 'Organization', id: string, owner?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null } };

export type UpdateOrganizationPhoneNumberMutationVariables = Exact<{
  organizationId: Scalars['ID'];
  input: PhoneNumberUpdateInput;
}>;


export type UpdateOrganizationPhoneNumberMutation = { __typename?: 'Mutation', phoneNumberUpdateInOrganization: { __typename?: 'PhoneNumber', label?: PhoneNumberLabel | null, primary: boolean, id: string, e164?: string | null, rawPhoneNumber?: string | null } };

export type UpdateOrganizationWebsiteMutationVariables = Exact<{
  input: OrganizationUpdateInput;
}>;


export type UpdateOrganizationWebsiteMutation = { __typename?: 'Mutation', organization_Update: { __typename?: 'Organization', id: string, website?: string | null } };

export type UpdateRenewalForecastMutationVariables = Exact<{
  input: RenewalForecastInput;
}>;


export type UpdateRenewalForecastMutation = { __typename?: 'Mutation', organization_UpdateRenewalForecastAsync: string };

export type UpdateRenewalLikelihoodMutationVariables = Exact<{
  input: RenewalLikelihoodInput;
}>;


export type UpdateRenewalLikelihoodMutation = { __typename?: 'Mutation', organization_UpdateRenewalLikelihoodAsync: string };

export type CreateMeetingMutationVariables = Exact<{
  meeting: MeetingInput;
}>;


export type CreateMeetingMutation = { __typename?: 'Mutation', meeting_Create: { __typename?: 'Meeting', id: string, conferenceUrl?: string | null, name?: string | null, agenda?: string | null, agendaContentType?: string | null, meetingStartedAt?: any | null, meetingEndedAt?: any | null, attendedBy: Array<{ __typename?: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null } } | { __typename?: 'OrganizationParticipant' } | { __typename?: 'UserParticipant', userParticipant: { __typename?: 'User', id: string, lastName: string, firstName: string } }>, note: Array<{ __typename?: 'Note', id: string, appSource: string }>, createdBy: Array<{ __typename?: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', id: string } } | { __typename?: 'OrganizationParticipant' } | { __typename?: 'UserParticipant', userParticipant: { __typename?: 'User', id: string } }> } };

export type GetEmailValidationQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetEmailValidationQuery = { __typename?: 'Query', email: { __typename?: 'Email', id: string, emailValidationDetails: { __typename?: 'EmailValidationDetails', isReachable?: string | null, isValidSyntax?: boolean | null, canConnectSmtp?: boolean | null, acceptsMail?: boolean | null, hasFullInbox?: boolean | null, isCatchAll?: boolean | null, isDeliverable?: boolean | null, validated?: boolean | null, isDisabled?: boolean | null } } };

export type GetTenantNameQueryVariables = Exact<{ [key: string]: never; }>;


export type GetTenantNameQuery = { __typename?: 'Query', tenant: string };

export type LinkMeetingAttachmentMutationVariables = Exact<{
  meetingId: Scalars['ID'];
  attachmentId: Scalars['ID'];
}>;


export type LinkMeetingAttachmentMutation = { __typename?: 'Mutation', meeting_LinkAttachment: { __typename?: 'Meeting', id: string } };

export type MeetingLinkAttachmentMutationVariables = Exact<{
  meetingId: Scalars['ID'];
  attachmentId: Scalars['ID'];
}>;


export type MeetingLinkAttachmentMutation = { __typename?: 'Mutation', meeting_LinkAttachment: { __typename?: 'Meeting', id: string, includes: Array<{ __typename?: 'Attachment', id: string, name: string, mimeType: string }> } };

export type LinkMeetingAttendeeMutationVariables = Exact<{
  meetingId: Scalars['ID'];
  participant: MeetingParticipantInput;
}>;


export type LinkMeetingAttendeeMutation = { __typename?: 'Mutation', meeting_LinkAttendedBy: { __typename?: 'Meeting', id: string, attendedBy: Array<{ __typename?: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null } } | { __typename?: 'OrganizationParticipant' } | { __typename?: 'UserParticipant', userParticipant: { __typename?: 'User', id: string, lastName: string, firstName: string } }> } };

export type MeetingLinkRecordingMutationVariables = Exact<{
  meetingId: Scalars['ID'];
  attachmentId: Scalars['ID'];
}>;


export type MeetingLinkRecordingMutation = { __typename?: 'Mutation', meeting_LinkRecording: { __typename?: 'Meeting', id: string, agenda?: string | null, meetingStartedAt?: any | null, attendedBy: Array<{ __typename?: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', id: string, firstName?: string | null, lastName?: string | null, name?: string | null } } | { __typename?: 'OrganizationParticipant' } | { __typename?: 'UserParticipant', userParticipant: { __typename?: 'User', id: string, firstName: string, lastName: string } }>, recording?: { __typename?: 'Attachment', id: string } | null } };

export type MeetingUnlinkAttachmentMutationVariables = Exact<{
  meetingId: Scalars['ID'];
  attachmentId: Scalars['ID'];
}>;


export type MeetingUnlinkAttachmentMutation = { __typename?: 'Mutation', meeting_UnlinkAttachment: { __typename?: 'Meeting', id: string, includes: Array<{ __typename?: 'Attachment', id: string, name: string, mimeType: string }> } };

export type UnlinkMeetingAttendeeMutationVariables = Exact<{
  meetingId: Scalars['ID'];
  participant: MeetingParticipantInput;
}>;


export type UnlinkMeetingAttendeeMutation = { __typename?: 'Mutation', meeting_UnlinkAttendedBy: { __typename?: 'Meeting', id: string, attendedBy: Array<{ __typename?: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null } } | { __typename?: 'OrganizationParticipant' } | { __typename?: 'UserParticipant', userParticipant: { __typename?: 'User', id: string, lastName: string, firstName: string } }> } };

export type MeetingUnlinkRecordingMutationVariables = Exact<{
  meetingId: Scalars['ID'];
  attachmentId: Scalars['ID'];
}>;


export type MeetingUnlinkRecordingMutation = { __typename?: 'Mutation', meeting_UnlinkRecording: { __typename?: 'Meeting', id: string, includes: Array<{ __typename?: 'Attachment', id: string }> } };

export type NoteLinkAttachmentMutationVariables = Exact<{
  noteId: Scalars['ID'];
  attachmentId: Scalars['ID'];
}>;


export type NoteLinkAttachmentMutation = { __typename?: 'Mutation', note_LinkAttachment: { __typename?: 'Note', id: string, includes: Array<{ __typename?: 'Attachment', id: string, name: string, mimeType: string }> } };

export type NoteUnlinkAttachmentMutationVariables = Exact<{
  noteId: Scalars['ID'];
  attachmentId: Scalars['ID'];
}>;


export type NoteUnlinkAttachmentMutation = { __typename?: 'Mutation', note_UnlinkAttachment: { __typename?: 'Note', id: string, includes: Array<{ __typename?: 'Attachment', id: string, name: string, mimeType: string }> } };

export type RemoveNoteMutationVariables = Exact<{
  id: Scalars['ID'];
}>;


export type RemoveNoteMutation = { __typename?: 'Mutation', note_Delete: { __typename?: 'Result', result: boolean } };

export type UpdateLocationMutationVariables = Exact<{
  input: LocationUpdateInput;
}>;


export type UpdateLocationMutation = { __typename?: 'Mutation', location_Update: { __typename?: 'Location', locality?: string | null, rawAddress?: string | null, postalCode?: string | null, street?: string | null } };

export type UpdateMeetingMutationVariables = Exact<{
  meetingId: Scalars['ID'];
  meetingInput: MeetingUpdateInput;
}>;


export type UpdateMeetingMutation = { __typename?: 'Mutation', meeting_Update: { __typename?: 'Meeting', id: string, createdAt: any, agenda?: string | null, agendaContentType?: string | null, conferenceUrl?: string | null, meetingStartedAt?: any | null, meetingEndedAt?: any | null, attendedBy: Array<{ __typename?: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', id: string, firstName?: string | null, lastName?: string | null, name?: string | null } } | { __typename?: 'OrganizationParticipant' } | { __typename?: 'UserParticipant', userParticipant: { __typename?: 'User', id: string, firstName: string, lastName: string } }>, meetingCreatedBy: Array<{ __typename?: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', id: string } } | { __typename?: 'OrganizationParticipant' } | { __typename?: 'UserParticipant', userParticipant: { __typename?: 'User', id: string } }>, describedBy: Array<{ __typename?: 'Analysis', id: string, analysisType?: string | null, content?: string | null, contentType?: string | null }>, events: Array<{ __typename?: 'InteractionEvent', id: string, createdAt: any, channel?: string | null, content?: string | null, contentType?: string | null, sentBy: Array<{ __typename?: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', id: string, firstName?: string | null, lastName?: string | null, name?: string | null } } | { __typename?: 'EmailParticipant' } | { __typename?: 'JobRoleParticipant' } | { __typename?: 'OrganizationParticipant' } | { __typename?: 'PhoneNumberParticipant' } | { __typename?: 'UserParticipant', userParticipant: { __typename?: 'User', id: string, firstName: string, lastName: string } }>, sentTo: Array<{ __typename?: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', id: string, firstName?: string | null, lastName?: string | null, name?: string | null } } | { __typename?: 'EmailParticipant' } | { __typename?: 'JobRoleParticipant' } | { __typename?: 'OrganizationParticipant' } | { __typename?: 'PhoneNumberParticipant' } | { __typename?: 'UserParticipant', userParticipant: { __typename?: 'User', id: string, firstName: string, lastName: string } }>, includes: Array<{ __typename?: 'Attachment', id: string, name: string, mimeType: string, extension: string, size: any }> }>, recording?: { __typename?: 'Attachment', id: string } | null, includes: Array<{ __typename?: 'Attachment', id: string, name: string, mimeType: string, extension: string, size: any }>, note: Array<{ __typename?: 'Note', id: string, appSource: string }> } };

export type UpdateNoteMutationVariables = Exact<{
  input: NoteUpdateInput;
}>;


export type UpdateNoteMutation = { __typename?: 'Mutation', note_Update: { __typename?: 'Note', id: string, createdAt: any, updatedAt: any, source: DataSource, sourceOfTruth: DataSource, appSource: string, createdBy?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null, includes: Array<{ __typename?: 'Attachment', id: string, name: string, mimeType: string, extension: string, size: any }> } };

export type GetUserByEmailQueryVariables = Exact<{
  email: Scalars['String'];
}>;


export type GetUserByEmailQuery = { __typename?: 'Query', user_ByEmail: { __typename?: 'User', id: string, firstName: string, lastName: string } };

export type GetUsersQueryVariables = Exact<{
  pagination: Pagination;
  where?: InputMaybe<Filter>;
}>;


export type GetUsersQuery = { __typename?: 'Query', users: { __typename?: 'UserPage', totalElements: any, content: Array<{ __typename?: 'User', id: string, firstName: string, lastName: string }> } };

export const LocationTotalFragmentDoc = gql`
    fragment LocationTotal on Location {
  id
  name
  createdAt
  updatedAt
  source
  appSource
  country
  region
  locality
  address
  address2
  zip
  addressType
  houseNumber
  postalCode
  plusFour
  commercial
  predirection
  district
  street
  rawAddress
  latitude
  longitude
}
    `;
export const NoteContentFragmentDoc = gql`
    fragment NoteContent on Note {
  id
  createdAt
  updatedAt
  createdBy {
    id
    firstName
    lastName
  }
  source
  sourceOfTruth
  appSource
  includes {
    id
    name
    mimeType
    extension
    size
  }
}
    `;
export const InteractionSessionFragmentFragmentDoc = gql`
    fragment InteractionSessionFragment on InteractionSession {
  id
  startedAt
  name
  status
  type
  events {
    content
    contentType
  }
}
    `;
export const InteractionEventFragmentFragmentDoc = gql`
    fragment InteractionEventFragment on InteractionEvent {
  id
  createdAt
  channel
  interactionSession {
    name
  }
  content
  contentType
  issue {
    externalLinks {
      type
      externalId
      externalUrl
    }
  }
  sentBy {
    ... on EmailParticipant {
      __typename
      emailParticipant {
        email
        id
        contacts {
          id
          name
          firstName
          lastName
          emails {
            email
          }
        }
        users {
          id
          firstName
          lastName
        }
        organizations {
          id
          name
        }
      }
    }
    ... on PhoneNumberParticipant {
      __typename
      phoneNumberParticipant {
        contacts {
          id
          name
          firstName
          lastName
        }
        users {
          id
          firstName
          lastName
        }
        organizations {
          id
          name
        }
        e164
        id
      }
    }
    ... on ContactParticipant {
      __typename
      contactParticipant {
        id
        name
        firstName
        lastName
      }
    }
    ... on UserParticipant {
      __typename
      userParticipant {
        id
        firstName
        lastName
      }
    }
  }
  sentTo {
    __typename
    ... on EmailParticipant {
      __typename
      emailParticipant {
        email
        contacts {
          id
          name
          firstName
          lastName
        }
        users {
          id
          firstName
          lastName
        }
        organizations {
          id
          name
        }
        id
      }
    }
    ... on PhoneNumberParticipant {
      __typename
      phoneNumberParticipant {
        e164
        id
        contacts {
          id
          name
          firstName
          lastName
        }
        users {
          id
          firstName
          lastName
        }
        organizations {
          id
          name
        }
      }
    }
    ... on ContactParticipant {
      __typename
      contactParticipant {
        name
        id
        firstName
        lastName
      }
    }
    ... on UserParticipant {
      __typename
      type
      userParticipant {
        id
        firstName
        lastName
      }
    }
  }
}
    `;
export const MeetingTimelineEventFragmentFragmentDoc = gql`
    fragment MeetingTimelineEventFragment on Meeting {
  id
  attendedBy {
    ... on UserParticipant {
      userParticipant {
        id
        firstName
        lastName
      }
    }
    ... on ContactParticipant {
      contactParticipant {
        id
        firstName
        lastName
        name
      }
    }
  }
  meetingCreatedBy: createdBy {
    ... on UserParticipant {
      userParticipant {
        id
      }
    }
    ... on ContactParticipant {
      contactParticipant {
        id
      }
    }
  }
  describedBy {
    id
    analysisType
    content
    contentType
  }
  events {
    id
    createdAt
    channel
    content
    contentType
    sentBy {
      ... on UserParticipant {
        userParticipant {
          id
          firstName
          lastName
        }
      }
      ... on ContactParticipant {
        contactParticipant {
          id
          firstName
          lastName
          name
        }
      }
    }
    sentTo {
      ... on UserParticipant {
        userParticipant {
          id
          firstName
          lastName
        }
      }
      ... on ContactParticipant {
        contactParticipant {
          id
          firstName
          lastName
          name
        }
      }
    }
    includes {
      id
      name
      mimeType
      extension
      size
    }
  }
  meetingStartedAt: startedAt
  meetingEndedAt: endedAt
  createdAt
  agenda
  agendaContentType
  recording {
    id
  }
  includes {
    id
    name
    mimeType
    extension
    size
  }
  conferenceUrl
  note {
    id
    appSource
  }
}
    `;
export const OrganizationBaseDetailsFragmentDoc = gql`
    fragment organizationBaseDetails on Organization {
  id
  name
  industry
}
    `;
export const ContactNameFragmentFragmentDoc = gql`
    fragment ContactNameFragment on Contact {
  firstName
  lastName
  name
}
    `;
export const JobRoleFragmentDoc = gql`
    fragment JobRole on JobRole {
  jobTitle
  primary
  id
}
    `;
export const TagFragmentDoc = gql`
    fragment Tag on Tag {
  id
  name
}
    `;
export const ContactPersonalDetailsFragmentDoc = gql`
    fragment ContactPersonalDetails on Contact {
  id
  ...ContactNameFragment
  source
  jobRoles {
    ...JobRole
    organization {
      id
      name
    }
  }
  tags {
    ...Tag
  }
}
    ${ContactNameFragmentFragmentDoc}
${JobRoleFragmentDoc}
${TagFragmentDoc}`;
export const EmailFragmentDoc = gql`
    fragment Email on Email {
  id
  primary
  email
}
    `;
export const LocationBaseDetailsFragmentDoc = gql`
    fragment LocationBaseDetails on Location {
  id
  name
  country
  region
  locality
  zip
  street
  postalCode
  houseNumber
}
    `;
export const OrganizationDetailsFragmentDoc = gql`
    fragment OrganizationDetails on Organization {
  id
  name
  description
  source
  industry
  emails {
    ...Email
  }
  locations {
    ...LocationBaseDetails
    rawAddress
  }
  website
  domains
  updatedAt
  tags {
    ...Tag
  }
}
    ${EmailFragmentDoc}
${LocationBaseDetailsFragmentDoc}
${TagFragmentDoc}`;
export const EmailWithValidationFragmentDoc = gql`
    fragment EmailWithValidation on Email {
  id
  primary
  email
  emailValidationDetails {
    isReachable
    isValidSyntax
    canConnectSmtp
    acceptsMail
    hasFullInbox
    isCatchAll
    isDeliverable
    validated
    isDisabled
  }
}
    `;
export const PhoneNumberFragmentDoc = gql`
    fragment PhoneNumber on PhoneNumber {
  id
  primary
  e164
  rawPhoneNumber
}
    `;
export const ContactCommunicationChannelsDetailsFragmentDoc = gql`
    fragment ContactCommunicationChannelsDetails on Contact {
  id
  emails {
    label
    ...EmailWithValidation
  }
  phoneNumbers {
    label
    ...PhoneNumber
  }
}
    ${EmailWithValidationFragmentDoc}
${PhoneNumberFragmentDoc}`;
export const OrganizationContactsFragmentDoc = gql`
    fragment OrganizationContacts on Organization {
  contacts {
    content {
      id
      name
      firstName
      lastName
      jobRoles {
        ...JobRole
      }
      ...ContactCommunicationChannelsDetails
    }
  }
}
    ${JobRoleFragmentDoc}
${ContactCommunicationChannelsDetailsFragmentDoc}`;
export const CreateTagDocument = gql`
    mutation CreateTag($input: TagInput!) {
  tag_Create(input: $input) {
    id
    name
    createdAt
    updatedAt
    source
  }
}
    `;
export type CreateTagMutationFn = Apollo.MutationFunction<CreateTagMutation, CreateTagMutationVariables>;

/**
 * __useCreateTagMutation__
 *
 * To run a mutation, you first call `useCreateTagMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateTagMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createTagMutation, { data, loading, error }] = useCreateTagMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useCreateTagMutation(baseOptions?: Apollo.MutationHookOptions<CreateTagMutation, CreateTagMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreateTagMutation, CreateTagMutationVariables>(CreateTagDocument, options);
      }
export type CreateTagMutationHookResult = ReturnType<typeof useCreateTagMutation>;
export type CreateTagMutationResult = Apollo.MutationResult<CreateTagMutation>;
export type CreateTagMutationOptions = Apollo.BaseMutationOptions<CreateTagMutation, CreateTagMutationVariables>;
export const DeleteTagDocument = gql`
    mutation DeleteTag($id: ID!) {
  tag_Delete(id: $id) {
    result
  }
}
    `;
export type DeleteTagMutationFn = Apollo.MutationFunction<DeleteTagMutation, DeleteTagMutationVariables>;

/**
 * __useDeleteTagMutation__
 *
 * To run a mutation, you first call `useDeleteTagMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteTagMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteTagMutation, { data, loading, error }] = useDeleteTagMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useDeleteTagMutation(baseOptions?: Apollo.MutationHookOptions<DeleteTagMutation, DeleteTagMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteTagMutation, DeleteTagMutationVariables>(DeleteTagDocument, options);
      }
export type DeleteTagMutationHookResult = ReturnType<typeof useDeleteTagMutation>;
export type DeleteTagMutationResult = Apollo.MutationResult<DeleteTagMutation>;
export type DeleteTagMutationOptions = Apollo.BaseMutationOptions<DeleteTagMutation, DeleteTagMutationVariables>;
export const GetTagsDocument = gql`
    query GetTags {
  tags {
    id
    name
  }
}
    `;

/**
 * __useGetTagsQuery__
 *
 * To run a query within a React component, call `useGetTagsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetTagsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetTagsQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetTagsQuery(baseOptions?: Apollo.QueryHookOptions<GetTagsQuery, GetTagsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetTagsQuery, GetTagsQueryVariables>(GetTagsDocument, options);
      }
export function useGetTagsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetTagsQuery, GetTagsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetTagsQuery, GetTagsQueryVariables>(GetTagsDocument, options);
        }
export type GetTagsQueryHookResult = ReturnType<typeof useGetTagsQuery>;
export type GetTagsLazyQueryHookResult = ReturnType<typeof useGetTagsLazyQuery>;
export type GetTagsQueryResult = Apollo.QueryResult<GetTagsQuery, GetTagsQueryVariables>;
export const UpdateTagDocument = gql`
    mutation UpdateTag($input: TagUpdateInput!) {
  tag_Update(input: $input) {
    id
    name
  }
}
    `;
export type UpdateTagMutationFn = Apollo.MutationFunction<UpdateTagMutation, UpdateTagMutationVariables>;

/**
 * __useUpdateTagMutation__
 *
 * To run a mutation, you first call `useUpdateTagMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateTagMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateTagMutation, { data, loading, error }] = useUpdateTagMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateTagMutation(baseOptions?: Apollo.MutationHookOptions<UpdateTagMutation, UpdateTagMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateTagMutation, UpdateTagMutationVariables>(UpdateTagDocument, options);
      }
export type UpdateTagMutationHookResult = ReturnType<typeof useUpdateTagMutation>;
export type UpdateTagMutationResult = Apollo.MutationResult<UpdateTagMutation>;
export type UpdateTagMutationOptions = Apollo.BaseMutationOptions<UpdateTagMutation, UpdateTagMutationVariables>;
export const AddEmailToContactDocument = gql`
    mutation addEmailToContact($contactId: ID!, $input: EmailInput!) {
  emailMergeToContact(contactId: $contactId, input: $input) {
    ...Email
    label
  }
}
    ${EmailFragmentDoc}`;
export type AddEmailToContactMutationFn = Apollo.MutationFunction<AddEmailToContactMutation, AddEmailToContactMutationVariables>;

/**
 * __useAddEmailToContactMutation__
 *
 * To run a mutation, you first call `useAddEmailToContactMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddEmailToContactMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addEmailToContactMutation, { data, loading, error }] = useAddEmailToContactMutation({
 *   variables: {
 *      contactId: // value for 'contactId'
 *      input: // value for 'input'
 *   },
 * });
 */
export function useAddEmailToContactMutation(baseOptions?: Apollo.MutationHookOptions<AddEmailToContactMutation, AddEmailToContactMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddEmailToContactMutation, AddEmailToContactMutationVariables>(AddEmailToContactDocument, options);
      }
export type AddEmailToContactMutationHookResult = ReturnType<typeof useAddEmailToContactMutation>;
export type AddEmailToContactMutationResult = Apollo.MutationResult<AddEmailToContactMutation>;
export type AddEmailToContactMutationOptions = Apollo.BaseMutationOptions<AddEmailToContactMutation, AddEmailToContactMutationVariables>;
export const AddLocationToContactDocument = gql`
    mutation addLocationToContact($contactId: ID!) {
  contact_AddNewLocation(contactId: $contactId) {
    id
  }
}
    `;
export type AddLocationToContactMutationFn = Apollo.MutationFunction<AddLocationToContactMutation, AddLocationToContactMutationVariables>;

/**
 * __useAddLocationToContactMutation__
 *
 * To run a mutation, you first call `useAddLocationToContactMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddLocationToContactMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addLocationToContactMutation, { data, loading, error }] = useAddLocationToContactMutation({
 *   variables: {
 *      contactId: // value for 'contactId'
 *   },
 * });
 */
export function useAddLocationToContactMutation(baseOptions?: Apollo.MutationHookOptions<AddLocationToContactMutation, AddLocationToContactMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddLocationToContactMutation, AddLocationToContactMutationVariables>(AddLocationToContactDocument, options);
      }
export type AddLocationToContactMutationHookResult = ReturnType<typeof useAddLocationToContactMutation>;
export type AddLocationToContactMutationResult = Apollo.MutationResult<AddLocationToContactMutation>;
export type AddLocationToContactMutationOptions = Apollo.BaseMutationOptions<AddLocationToContactMutation, AddLocationToContactMutationVariables>;
export const AddPhoneToContactDocument = gql`
    mutation addPhoneToContact($contactId: ID!, $input: PhoneNumberInput!) {
  phoneNumberMergeToContact(contactId: $contactId, input: $input) {
    ...PhoneNumber
    label
  }
}
    ${PhoneNumberFragmentDoc}`;
export type AddPhoneToContactMutationFn = Apollo.MutationFunction<AddPhoneToContactMutation, AddPhoneToContactMutationVariables>;

/**
 * __useAddPhoneToContactMutation__
 *
 * To run a mutation, you first call `useAddPhoneToContactMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddPhoneToContactMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addPhoneToContactMutation, { data, loading, error }] = useAddPhoneToContactMutation({
 *   variables: {
 *      contactId: // value for 'contactId'
 *      input: // value for 'input'
 *   },
 * });
 */
export function useAddPhoneToContactMutation(baseOptions?: Apollo.MutationHookOptions<AddPhoneToContactMutation, AddPhoneToContactMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddPhoneToContactMutation, AddPhoneToContactMutationVariables>(AddPhoneToContactDocument, options);
      }
export type AddPhoneToContactMutationHookResult = ReturnType<typeof useAddPhoneToContactMutation>;
export type AddPhoneToContactMutationResult = Apollo.MutationResult<AddPhoneToContactMutation>;
export type AddPhoneToContactMutationOptions = Apollo.BaseMutationOptions<AddPhoneToContactMutation, AddPhoneToContactMutationVariables>;
export const AddTagToContactDocument = gql`
    mutation addTagToContact($input: ContactTagInput!) {
  contact_AddTagById(input: $input) {
    id
    tags {
      ...Tag
    }
  }
}
    ${TagFragmentDoc}`;
export type AddTagToContactMutationFn = Apollo.MutationFunction<AddTagToContactMutation, AddTagToContactMutationVariables>;

/**
 * __useAddTagToContactMutation__
 *
 * To run a mutation, you first call `useAddTagToContactMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddTagToContactMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addTagToContactMutation, { data, loading, error }] = useAddTagToContactMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useAddTagToContactMutation(baseOptions?: Apollo.MutationHookOptions<AddTagToContactMutation, AddTagToContactMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddTagToContactMutation, AddTagToContactMutationVariables>(AddTagToContactDocument, options);
      }
export type AddTagToContactMutationHookResult = ReturnType<typeof useAddTagToContactMutation>;
export type AddTagToContactMutationResult = Apollo.MutationResult<AddTagToContactMutation>;
export type AddTagToContactMutationOptions = Apollo.BaseMutationOptions<AddTagToContactMutation, AddTagToContactMutationVariables>;
export const ArchiveContactDocument = gql`
    mutation archiveContact($id: ID!) {
  contact_Archive(contactId: $id) {
    result
  }
}
    `;
export type ArchiveContactMutationFn = Apollo.MutationFunction<ArchiveContactMutation, ArchiveContactMutationVariables>;

/**
 * __useArchiveContactMutation__
 *
 * To run a mutation, you first call `useArchiveContactMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useArchiveContactMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [archiveContactMutation, { data, loading, error }] = useArchiveContactMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useArchiveContactMutation(baseOptions?: Apollo.MutationHookOptions<ArchiveContactMutation, ArchiveContactMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ArchiveContactMutation, ArchiveContactMutationVariables>(ArchiveContactDocument, options);
      }
export type ArchiveContactMutationHookResult = ReturnType<typeof useArchiveContactMutation>;
export type ArchiveContactMutationResult = Apollo.MutationResult<ArchiveContactMutation>;
export type ArchiveContactMutationOptions = Apollo.BaseMutationOptions<ArchiveContactMutation, ArchiveContactMutationVariables>;
export const AttachOrganizationToContactDocument = gql`
    mutation attachOrganizationToContact($input: ContactOrganizationInput!) {
  contact_AddOrganizationById(input: $input) {
    ...ContactPersonalDetails
  }
}
    ${ContactPersonalDetailsFragmentDoc}`;
export type AttachOrganizationToContactMutationFn = Apollo.MutationFunction<AttachOrganizationToContactMutation, AttachOrganizationToContactMutationVariables>;

/**
 * __useAttachOrganizationToContactMutation__
 *
 * To run a mutation, you first call `useAttachOrganizationToContactMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAttachOrganizationToContactMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [attachOrganizationToContactMutation, { data, loading, error }] = useAttachOrganizationToContactMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useAttachOrganizationToContactMutation(baseOptions?: Apollo.MutationHookOptions<AttachOrganizationToContactMutation, AttachOrganizationToContactMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AttachOrganizationToContactMutation, AttachOrganizationToContactMutationVariables>(AttachOrganizationToContactDocument, options);
      }
export type AttachOrganizationToContactMutationHookResult = ReturnType<typeof useAttachOrganizationToContactMutation>;
export type AttachOrganizationToContactMutationResult = Apollo.MutationResult<AttachOrganizationToContactMutation>;
export type AttachOrganizationToContactMutationOptions = Apollo.BaseMutationOptions<AttachOrganizationToContactMutation, AttachOrganizationToContactMutationVariables>;
export const CreateContactDocument = gql`
    mutation createContact($input: ContactInput!) {
  contact_Create(input: $input) {
    ...ContactPersonalDetails
    ...ContactCommunicationChannelsDetails
  }
}
    ${ContactPersonalDetailsFragmentDoc}
${ContactCommunicationChannelsDetailsFragmentDoc}`;
export type CreateContactMutationFn = Apollo.MutationFunction<CreateContactMutation, CreateContactMutationVariables>;

/**
 * __useCreateContactMutation__
 *
 * To run a mutation, you first call `useCreateContactMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateContactMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createContactMutation, { data, loading, error }] = useCreateContactMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useCreateContactMutation(baseOptions?: Apollo.MutationHookOptions<CreateContactMutation, CreateContactMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreateContactMutation, CreateContactMutationVariables>(CreateContactDocument, options);
      }
export type CreateContactMutationHookResult = ReturnType<typeof useCreateContactMutation>;
export type CreateContactMutationResult = Apollo.MutationResult<CreateContactMutation>;
export type CreateContactMutationOptions = Apollo.BaseMutationOptions<CreateContactMutation, CreateContactMutationVariables>;
export const CreateContactJobRoleDocument = gql`
    mutation createContactJobRole($contactId: ID!, $input: JobRoleInput!) {
  jobRole_Create(contactId: $contactId, input: $input) {
    ...JobRole
    organization {
      id
      name
    }
  }
}
    ${JobRoleFragmentDoc}`;
export type CreateContactJobRoleMutationFn = Apollo.MutationFunction<CreateContactJobRoleMutation, CreateContactJobRoleMutationVariables>;

/**
 * __useCreateContactJobRoleMutation__
 *
 * To run a mutation, you first call `useCreateContactJobRoleMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateContactJobRoleMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createContactJobRoleMutation, { data, loading, error }] = useCreateContactJobRoleMutation({
 *   variables: {
 *      contactId: // value for 'contactId'
 *      input: // value for 'input'
 *   },
 * });
 */
export function useCreateContactJobRoleMutation(baseOptions?: Apollo.MutationHookOptions<CreateContactJobRoleMutation, CreateContactJobRoleMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreateContactJobRoleMutation, CreateContactJobRoleMutationVariables>(CreateContactJobRoleDocument, options);
      }
export type CreateContactJobRoleMutationHookResult = ReturnType<typeof useCreateContactJobRoleMutation>;
export type CreateContactJobRoleMutationResult = Apollo.MutationResult<CreateContactJobRoleMutation>;
export type CreateContactJobRoleMutationOptions = Apollo.BaseMutationOptions<CreateContactJobRoleMutation, CreateContactJobRoleMutationVariables>;
export const CreateContactNoteDocument = gql`
    mutation createContactNote($contactId: ID!, $input: NoteInput!) {
  note_CreateForContact(contactId: $contactId, input: $input) {
    ...NoteContent
  }
}
    ${NoteContentFragmentDoc}`;
export type CreateContactNoteMutationFn = Apollo.MutationFunction<CreateContactNoteMutation, CreateContactNoteMutationVariables>;

/**
 * __useCreateContactNoteMutation__
 *
 * To run a mutation, you first call `useCreateContactNoteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateContactNoteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createContactNoteMutation, { data, loading, error }] = useCreateContactNoteMutation({
 *   variables: {
 *      contactId: // value for 'contactId'
 *      input: // value for 'input'
 *   },
 * });
 */
export function useCreateContactNoteMutation(baseOptions?: Apollo.MutationHookOptions<CreateContactNoteMutation, CreateContactNoteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreateContactNoteMutation, CreateContactNoteMutationVariables>(CreateContactNoteDocument, options);
      }
export type CreateContactNoteMutationHookResult = ReturnType<typeof useCreateContactNoteMutation>;
export type CreateContactNoteMutationResult = Apollo.MutationResult<CreateContactNoteMutation>;
export type CreateContactNoteMutationOptions = Apollo.BaseMutationOptions<CreateContactNoteMutation, CreateContactNoteMutationVariables>;
export const CreatePhoneCallInteractionEventDocument = gql`
    mutation CreatePhoneCallInteractionEvent($contactId: ID, $sentBy: String, $content: String, $contentType: String) {
  interactionEvent_Create(
    event: {channel: "VOICE", sentTo: [{contactID: $contactId}], sentBy: [{email: $sentBy}], appSource: "Openline", content: $content, contentType: $contentType}
  ) {
    ...InteractionEventFragment
  }
}
    ${InteractionEventFragmentFragmentDoc}`;
export type CreatePhoneCallInteractionEventMutationFn = Apollo.MutationFunction<CreatePhoneCallInteractionEventMutation, CreatePhoneCallInteractionEventMutationVariables>;

/**
 * __useCreatePhoneCallInteractionEventMutation__
 *
 * To run a mutation, you first call `useCreatePhoneCallInteractionEventMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreatePhoneCallInteractionEventMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createPhoneCallInteractionEventMutation, { data, loading, error }] = useCreatePhoneCallInteractionEventMutation({
 *   variables: {
 *      contactId: // value for 'contactId'
 *      sentBy: // value for 'sentBy'
 *      content: // value for 'content'
 *      contentType: // value for 'contentType'
 *   },
 * });
 */
export function useCreatePhoneCallInteractionEventMutation(baseOptions?: Apollo.MutationHookOptions<CreatePhoneCallInteractionEventMutation, CreatePhoneCallInteractionEventMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreatePhoneCallInteractionEventMutation, CreatePhoneCallInteractionEventMutationVariables>(CreatePhoneCallInteractionEventDocument, options);
      }
export type CreatePhoneCallInteractionEventMutationHookResult = ReturnType<typeof useCreatePhoneCallInteractionEventMutation>;
export type CreatePhoneCallInteractionEventMutationResult = Apollo.MutationResult<CreatePhoneCallInteractionEventMutation>;
export type CreatePhoneCallInteractionEventMutationOptions = Apollo.BaseMutationOptions<CreatePhoneCallInteractionEventMutation, CreatePhoneCallInteractionEventMutationVariables>;
export const RemoveContactJobRoleDocument = gql`
    mutation removeContactJobRole($contactId: ID!, $roleId: ID!) {
  jobRole_Delete(contactId: $contactId, roleId: $roleId) {
    result
  }
}
    `;
export type RemoveContactJobRoleMutationFn = Apollo.MutationFunction<RemoveContactJobRoleMutation, RemoveContactJobRoleMutationVariables>;

/**
 * __useRemoveContactJobRoleMutation__
 *
 * To run a mutation, you first call `useRemoveContactJobRoleMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemoveContactJobRoleMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removeContactJobRoleMutation, { data, loading, error }] = useRemoveContactJobRoleMutation({
 *   variables: {
 *      contactId: // value for 'contactId'
 *      roleId: // value for 'roleId'
 *   },
 * });
 */
export function useRemoveContactJobRoleMutation(baseOptions?: Apollo.MutationHookOptions<RemoveContactJobRoleMutation, RemoveContactJobRoleMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RemoveContactJobRoleMutation, RemoveContactJobRoleMutationVariables>(RemoveContactJobRoleDocument, options);
      }
export type RemoveContactJobRoleMutationHookResult = ReturnType<typeof useRemoveContactJobRoleMutation>;
export type RemoveContactJobRoleMutationResult = Apollo.MutationResult<RemoveContactJobRoleMutation>;
export type RemoveContactJobRoleMutationOptions = Apollo.BaseMutationOptions<RemoveContactJobRoleMutation, RemoveContactJobRoleMutationVariables>;
export const GetContactDocument = gql`
    query GetContact($id: ID!) {
  contact(id: $id) {
    ...ContactPersonalDetails
    owner {
      id
      firstName
      lastName
    }
    ...ContactCommunicationChannelsDetails
    jobRoles {
      ...JobRole
    }
    source
    locations {
      ...LocationBaseDetails
      rawAddress
    }
  }
}
    ${ContactPersonalDetailsFragmentDoc}
${ContactCommunicationChannelsDetailsFragmentDoc}
${JobRoleFragmentDoc}
${LocationBaseDetailsFragmentDoc}`;

/**
 * __useGetContactQuery__
 *
 * To run a query within a React component, call `useGetContactQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetContactQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetContactQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetContactQuery(baseOptions: Apollo.QueryHookOptions<GetContactQuery, GetContactQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetContactQuery, GetContactQueryVariables>(GetContactDocument, options);
      }
export function useGetContactLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetContactQuery, GetContactQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetContactQuery, GetContactQueryVariables>(GetContactDocument, options);
        }
export type GetContactQueryHookResult = ReturnType<typeof useGetContactQuery>;
export type GetContactLazyQueryHookResult = ReturnType<typeof useGetContactLazyQuery>;
export type GetContactQueryResult = Apollo.QueryResult<GetContactQuery, GetContactQueryVariables>;
export const GetContactCommunicationChannelsDocument = gql`
    query GetContactCommunicationChannels($id: ID!) {
  contact(id: $id) {
    ...ContactNameFragment
    ...ContactCommunicationChannelsDetails
  }
}
    ${ContactNameFragmentFragmentDoc}
${ContactCommunicationChannelsDetailsFragmentDoc}`;

/**
 * __useGetContactCommunicationChannelsQuery__
 *
 * To run a query within a React component, call `useGetContactCommunicationChannelsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetContactCommunicationChannelsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetContactCommunicationChannelsQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetContactCommunicationChannelsQuery(baseOptions: Apollo.QueryHookOptions<GetContactCommunicationChannelsQuery, GetContactCommunicationChannelsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetContactCommunicationChannelsQuery, GetContactCommunicationChannelsQueryVariables>(GetContactCommunicationChannelsDocument, options);
      }
export function useGetContactCommunicationChannelsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetContactCommunicationChannelsQuery, GetContactCommunicationChannelsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetContactCommunicationChannelsQuery, GetContactCommunicationChannelsQueryVariables>(GetContactCommunicationChannelsDocument, options);
        }
export type GetContactCommunicationChannelsQueryHookResult = ReturnType<typeof useGetContactCommunicationChannelsQuery>;
export type GetContactCommunicationChannelsLazyQueryHookResult = ReturnType<typeof useGetContactCommunicationChannelsLazyQuery>;
export type GetContactCommunicationChannelsQueryResult = Apollo.QueryResult<GetContactCommunicationChannelsQuery, GetContactCommunicationChannelsQueryVariables>;
export const GetContactListDocument = gql`
    query GetContactList($pagination: Pagination!, $where: Filter, $sort: [SortBy!]) {
  contacts(pagination: $pagination, where: $where, sort: $sort) {
    content {
      id
      ...ContactNameFragment
      emails {
        id
        email
      }
    }
    totalElements
  }
}
    ${ContactNameFragmentFragmentDoc}`;

/**
 * __useGetContactListQuery__
 *
 * To run a query within a React component, call `useGetContactListQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetContactListQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetContactListQuery({
 *   variables: {
 *      pagination: // value for 'pagination'
 *      where: // value for 'where'
 *      sort: // value for 'sort'
 *   },
 * });
 */
export function useGetContactListQuery(baseOptions: Apollo.QueryHookOptions<GetContactListQuery, GetContactListQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetContactListQuery, GetContactListQueryVariables>(GetContactListDocument, options);
      }
export function useGetContactListLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetContactListQuery, GetContactListQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetContactListQuery, GetContactListQueryVariables>(GetContactListDocument, options);
        }
export type GetContactListQueryHookResult = ReturnType<typeof useGetContactListQuery>;
export type GetContactListLazyQueryHookResult = ReturnType<typeof useGetContactListLazyQuery>;
export type GetContactListQueryResult = Apollo.QueryResult<GetContactListQuery, GetContactListQueryVariables>;
export const GetContactLocationsDocument = gql`
    query GetContactLocations($id: ID!) {
  contact(id: $id) {
    locations {
      ...LocationBaseDetails
      rawAddress
    }
  }
}
    ${LocationBaseDetailsFragmentDoc}`;

/**
 * __useGetContactLocationsQuery__
 *
 * To run a query within a React component, call `useGetContactLocationsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetContactLocationsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetContactLocationsQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetContactLocationsQuery(baseOptions: Apollo.QueryHookOptions<GetContactLocationsQuery, GetContactLocationsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetContactLocationsQuery, GetContactLocationsQueryVariables>(GetContactLocationsDocument, options);
      }
export function useGetContactLocationsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetContactLocationsQuery, GetContactLocationsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetContactLocationsQuery, GetContactLocationsQueryVariables>(GetContactLocationsDocument, options);
        }
export type GetContactLocationsQueryHookResult = ReturnType<typeof useGetContactLocationsQuery>;
export type GetContactLocationsLazyQueryHookResult = ReturnType<typeof useGetContactLocationsLazyQuery>;
export type GetContactLocationsQueryResult = Apollo.QueryResult<GetContactLocationsQuery, GetContactLocationsQueryVariables>;
export const GetContactMentionSuggestionsDocument = gql`
    query GetContactMentionSuggestions($pagination: Pagination!, $where: Filter, $sort: [SortBy!]) {
  contacts(pagination: $pagination, where: $where, sort: $sort) {
    content {
      id
      firstName
      lastName
    }
  }
}
    `;

/**
 * __useGetContactMentionSuggestionsQuery__
 *
 * To run a query within a React component, call `useGetContactMentionSuggestionsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetContactMentionSuggestionsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetContactMentionSuggestionsQuery({
 *   variables: {
 *      pagination: // value for 'pagination'
 *      where: // value for 'where'
 *      sort: // value for 'sort'
 *   },
 * });
 */
export function useGetContactMentionSuggestionsQuery(baseOptions: Apollo.QueryHookOptions<GetContactMentionSuggestionsQuery, GetContactMentionSuggestionsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetContactMentionSuggestionsQuery, GetContactMentionSuggestionsQueryVariables>(GetContactMentionSuggestionsDocument, options);
      }
export function useGetContactMentionSuggestionsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetContactMentionSuggestionsQuery, GetContactMentionSuggestionsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetContactMentionSuggestionsQuery, GetContactMentionSuggestionsQueryVariables>(GetContactMentionSuggestionsDocument, options);
        }
export type GetContactMentionSuggestionsQueryHookResult = ReturnType<typeof useGetContactMentionSuggestionsQuery>;
export type GetContactMentionSuggestionsLazyQueryHookResult = ReturnType<typeof useGetContactMentionSuggestionsLazyQuery>;
export type GetContactMentionSuggestionsQueryResult = Apollo.QueryResult<GetContactMentionSuggestionsQuery, GetContactMentionSuggestionsQueryVariables>;
export const GetContactNameByEmailDocument = gql`
    query GetContactNameByEmail($email: String!) {
  contact_ByEmail(email: $email) {
    id
    ...ContactNameFragment
  }
}
    ${ContactNameFragmentFragmentDoc}`;

/**
 * __useGetContactNameByEmailQuery__
 *
 * To run a query within a React component, call `useGetContactNameByEmailQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetContactNameByEmailQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetContactNameByEmailQuery({
 *   variables: {
 *      email: // value for 'email'
 *   },
 * });
 */
export function useGetContactNameByEmailQuery(baseOptions: Apollo.QueryHookOptions<GetContactNameByEmailQuery, GetContactNameByEmailQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetContactNameByEmailQuery, GetContactNameByEmailQueryVariables>(GetContactNameByEmailDocument, options);
      }
export function useGetContactNameByEmailLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetContactNameByEmailQuery, GetContactNameByEmailQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetContactNameByEmailQuery, GetContactNameByEmailQueryVariables>(GetContactNameByEmailDocument, options);
        }
export type GetContactNameByEmailQueryHookResult = ReturnType<typeof useGetContactNameByEmailQuery>;
export type GetContactNameByEmailLazyQueryHookResult = ReturnType<typeof useGetContactNameByEmailLazyQuery>;
export type GetContactNameByEmailQueryResult = Apollo.QueryResult<GetContactNameByEmailQuery, GetContactNameByEmailQueryVariables>;
export const GetContactNameByIdDocument = gql`
    query GetContactNameById($id: ID!) {
  contact(id: $id) {
    id
    ...ContactNameFragment
  }
}
    ${ContactNameFragmentFragmentDoc}`;

/**
 * __useGetContactNameByIdQuery__
 *
 * To run a query within a React component, call `useGetContactNameByIdQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetContactNameByIdQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetContactNameByIdQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetContactNameByIdQuery(baseOptions: Apollo.QueryHookOptions<GetContactNameByIdQuery, GetContactNameByIdQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetContactNameByIdQuery, GetContactNameByIdQueryVariables>(GetContactNameByIdDocument, options);
      }
export function useGetContactNameByIdLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetContactNameByIdQuery, GetContactNameByIdQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetContactNameByIdQuery, GetContactNameByIdQueryVariables>(GetContactNameByIdDocument, options);
        }
export type GetContactNameByIdQueryHookResult = ReturnType<typeof useGetContactNameByIdQuery>;
export type GetContactNameByIdLazyQueryHookResult = ReturnType<typeof useGetContactNameByIdLazyQuery>;
export type GetContactNameByIdQueryResult = Apollo.QueryResult<GetContactNameByIdQuery, GetContactNameByIdQueryVariables>;
export const GetContactNameByPhoneNumberDocument = gql`
    query GetContactNameByPhoneNumber($e164: String!) {
  contact_ByPhone(e164: $e164) {
    id
    ...ContactNameFragment
  }
}
    ${ContactNameFragmentFragmentDoc}`;

/**
 * __useGetContactNameByPhoneNumberQuery__
 *
 * To run a query within a React component, call `useGetContactNameByPhoneNumberQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetContactNameByPhoneNumberQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetContactNameByPhoneNumberQuery({
 *   variables: {
 *      e164: // value for 'e164'
 *   },
 * });
 */
export function useGetContactNameByPhoneNumberQuery(baseOptions: Apollo.QueryHookOptions<GetContactNameByPhoneNumberQuery, GetContactNameByPhoneNumberQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetContactNameByPhoneNumberQuery, GetContactNameByPhoneNumberQueryVariables>(GetContactNameByPhoneNumberDocument, options);
      }
export function useGetContactNameByPhoneNumberLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetContactNameByPhoneNumberQuery, GetContactNameByPhoneNumberQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetContactNameByPhoneNumberQuery, GetContactNameByPhoneNumberQueryVariables>(GetContactNameByPhoneNumberDocument, options);
        }
export type GetContactNameByPhoneNumberQueryHookResult = ReturnType<typeof useGetContactNameByPhoneNumberQuery>;
export type GetContactNameByPhoneNumberLazyQueryHookResult = ReturnType<typeof useGetContactNameByPhoneNumberLazyQuery>;
export type GetContactNameByPhoneNumberQueryResult = Apollo.QueryResult<GetContactNameByPhoneNumberQuery, GetContactNameByPhoneNumberQueryVariables>;
export const GetContactNotesDocument = gql`
    query GetContactNotes($id: ID!, $pagination: Pagination) {
  contact(id: $id) {
    notes(pagination: $pagination) {
      content {
        ...NoteContent
      }
    }
  }
}
    ${NoteContentFragmentDoc}`;

/**
 * __useGetContactNotesQuery__
 *
 * To run a query within a React component, call `useGetContactNotesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetContactNotesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetContactNotesQuery({
 *   variables: {
 *      id: // value for 'id'
 *      pagination: // value for 'pagination'
 *   },
 * });
 */
export function useGetContactNotesQuery(baseOptions: Apollo.QueryHookOptions<GetContactNotesQuery, GetContactNotesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetContactNotesQuery, GetContactNotesQueryVariables>(GetContactNotesDocument, options);
      }
export function useGetContactNotesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetContactNotesQuery, GetContactNotesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetContactNotesQuery, GetContactNotesQueryVariables>(GetContactNotesDocument, options);
        }
export type GetContactNotesQueryHookResult = ReturnType<typeof useGetContactNotesQuery>;
export type GetContactNotesLazyQueryHookResult = ReturnType<typeof useGetContactNotesLazyQuery>;
export type GetContactNotesQueryResult = Apollo.QueryResult<GetContactNotesQuery, GetContactNotesQueryVariables>;
export const GetContactPersonalDetailsDocument = gql`
    query GetContactPersonalDetails($id: ID!) {
  contact(id: $id) {
    ...ContactPersonalDetails
    owner {
      id
      firstName
      lastName
    }
  }
}
    ${ContactPersonalDetailsFragmentDoc}`;

/**
 * __useGetContactPersonalDetailsQuery__
 *
 * To run a query within a React component, call `useGetContactPersonalDetailsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetContactPersonalDetailsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetContactPersonalDetailsQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetContactPersonalDetailsQuery(baseOptions: Apollo.QueryHookOptions<GetContactPersonalDetailsQuery, GetContactPersonalDetailsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetContactPersonalDetailsQuery, GetContactPersonalDetailsQueryVariables>(GetContactPersonalDetailsDocument, options);
      }
export function useGetContactPersonalDetailsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetContactPersonalDetailsQuery, GetContactPersonalDetailsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetContactPersonalDetailsQuery, GetContactPersonalDetailsQueryVariables>(GetContactPersonalDetailsDocument, options);
        }
export type GetContactPersonalDetailsQueryHookResult = ReturnType<typeof useGetContactPersonalDetailsQuery>;
export type GetContactPersonalDetailsLazyQueryHookResult = ReturnType<typeof useGetContactPersonalDetailsLazyQuery>;
export type GetContactPersonalDetailsQueryResult = Apollo.QueryResult<GetContactPersonalDetailsQuery, GetContactPersonalDetailsQueryVariables>;
export const GetContactPersonalDetailsWithOrganizationsDocument = gql`
    query getContactPersonalDetailsWithOrganizations($id: ID!) {
  contact(id: $id) {
    ...ContactPersonalDetails
    organizations(pagination: {limit: 99999, page: 1}) {
      content {
        id
        name
      }
    }
  }
}
    ${ContactPersonalDetailsFragmentDoc}`;

/**
 * __useGetContactPersonalDetailsWithOrganizationsQuery__
 *
 * To run a query within a React component, call `useGetContactPersonalDetailsWithOrganizationsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetContactPersonalDetailsWithOrganizationsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetContactPersonalDetailsWithOrganizationsQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetContactPersonalDetailsWithOrganizationsQuery(baseOptions: Apollo.QueryHookOptions<GetContactPersonalDetailsWithOrganizationsQuery, GetContactPersonalDetailsWithOrganizationsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetContactPersonalDetailsWithOrganizationsQuery, GetContactPersonalDetailsWithOrganizationsQueryVariables>(GetContactPersonalDetailsWithOrganizationsDocument, options);
      }
export function useGetContactPersonalDetailsWithOrganizationsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetContactPersonalDetailsWithOrganizationsQuery, GetContactPersonalDetailsWithOrganizationsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetContactPersonalDetailsWithOrganizationsQuery, GetContactPersonalDetailsWithOrganizationsQueryVariables>(GetContactPersonalDetailsWithOrganizationsDocument, options);
        }
export type GetContactPersonalDetailsWithOrganizationsQueryHookResult = ReturnType<typeof useGetContactPersonalDetailsWithOrganizationsQuery>;
export type GetContactPersonalDetailsWithOrganizationsLazyQueryHookResult = ReturnType<typeof useGetContactPersonalDetailsWithOrganizationsLazyQuery>;
export type GetContactPersonalDetailsWithOrganizationsQueryResult = Apollo.QueryResult<GetContactPersonalDetailsWithOrganizationsQuery, GetContactPersonalDetailsWithOrganizationsQueryVariables>;
export const GetContactTagsDocument = gql`
    query GetContactTags($id: ID!) {
  contact(id: $id) {
    id
    tags {
      ...Tag
    }
  }
}
    ${TagFragmentDoc}`;

/**
 * __useGetContactTagsQuery__
 *
 * To run a query within a React component, call `useGetContactTagsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetContactTagsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetContactTagsQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetContactTagsQuery(baseOptions: Apollo.QueryHookOptions<GetContactTagsQuery, GetContactTagsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetContactTagsQuery, GetContactTagsQueryVariables>(GetContactTagsDocument, options);
      }
export function useGetContactTagsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetContactTagsQuery, GetContactTagsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetContactTagsQuery, GetContactTagsQueryVariables>(GetContactTagsDocument, options);
        }
export type GetContactTagsQueryHookResult = ReturnType<typeof useGetContactTagsQuery>;
export type GetContactTagsLazyQueryHookResult = ReturnType<typeof useGetContactTagsLazyQuery>;
export type GetContactTagsQueryResult = Apollo.QueryResult<GetContactTagsQuery, GetContactTagsQueryVariables>;
export const MergeContactsDocument = gql`
    mutation mergeContacts($primaryContactId: ID!, $mergedContactIds: [ID!]!) {
  contact_Merge(
    primaryContactId: $primaryContactId
    mergedContactIds: $mergedContactIds
  ) {
    id
    ...ContactPersonalDetails
  }
}
    ${ContactPersonalDetailsFragmentDoc}`;
export type MergeContactsMutationFn = Apollo.MutationFunction<MergeContactsMutation, MergeContactsMutationVariables>;

/**
 * __useMergeContactsMutation__
 *
 * To run a mutation, you first call `useMergeContactsMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useMergeContactsMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [mergeContactsMutation, { data, loading, error }] = useMergeContactsMutation({
 *   variables: {
 *      primaryContactId: // value for 'primaryContactId'
 *      mergedContactIds: // value for 'mergedContactIds'
 *   },
 * });
 */
export function useMergeContactsMutation(baseOptions?: Apollo.MutationHookOptions<MergeContactsMutation, MergeContactsMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<MergeContactsMutation, MergeContactsMutationVariables>(MergeContactsDocument, options);
      }
export type MergeContactsMutationHookResult = ReturnType<typeof useMergeContactsMutation>;
export type MergeContactsMutationResult = Apollo.MutationResult<MergeContactsMutation>;
export type MergeContactsMutationOptions = Apollo.BaseMutationOptions<MergeContactsMutation, MergeContactsMutationVariables>;
export const RemoveEmailFromContactDocument = gql`
    mutation removeEmailFromContact($contactId: ID!, $id: ID!) {
  emailRemoveFromContactById(contactId: $contactId, id: $id) {
    result
  }
}
    `;
export type RemoveEmailFromContactMutationFn = Apollo.MutationFunction<RemoveEmailFromContactMutation, RemoveEmailFromContactMutationVariables>;

/**
 * __useRemoveEmailFromContactMutation__
 *
 * To run a mutation, you first call `useRemoveEmailFromContactMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemoveEmailFromContactMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removeEmailFromContactMutation, { data, loading, error }] = useRemoveEmailFromContactMutation({
 *   variables: {
 *      contactId: // value for 'contactId'
 *      id: // value for 'id'
 *   },
 * });
 */
export function useRemoveEmailFromContactMutation(baseOptions?: Apollo.MutationHookOptions<RemoveEmailFromContactMutation, RemoveEmailFromContactMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RemoveEmailFromContactMutation, RemoveEmailFromContactMutationVariables>(RemoveEmailFromContactDocument, options);
      }
export type RemoveEmailFromContactMutationHookResult = ReturnType<typeof useRemoveEmailFromContactMutation>;
export type RemoveEmailFromContactMutationResult = Apollo.MutationResult<RemoveEmailFromContactMutation>;
export type RemoveEmailFromContactMutationOptions = Apollo.BaseMutationOptions<RemoveEmailFromContactMutation, RemoveEmailFromContactMutationVariables>;
export const RemoveLocationFromContactDocument = gql`
    mutation removeLocationFromContact($locationId: ID!, $contactId: ID!) {
  location_RemoveFromContact(locationId: $locationId, contactId: $contactId) {
    id
    locations {
      ...LocationBaseDetails
    }
  }
}
    ${LocationBaseDetailsFragmentDoc}`;
export type RemoveLocationFromContactMutationFn = Apollo.MutationFunction<RemoveLocationFromContactMutation, RemoveLocationFromContactMutationVariables>;

/**
 * __useRemoveLocationFromContactMutation__
 *
 * To run a mutation, you first call `useRemoveLocationFromContactMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemoveLocationFromContactMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removeLocationFromContactMutation, { data, loading, error }] = useRemoveLocationFromContactMutation({
 *   variables: {
 *      locationId: // value for 'locationId'
 *      contactId: // value for 'contactId'
 *   },
 * });
 */
export function useRemoveLocationFromContactMutation(baseOptions?: Apollo.MutationHookOptions<RemoveLocationFromContactMutation, RemoveLocationFromContactMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RemoveLocationFromContactMutation, RemoveLocationFromContactMutationVariables>(RemoveLocationFromContactDocument, options);
      }
export type RemoveLocationFromContactMutationHookResult = ReturnType<typeof useRemoveLocationFromContactMutation>;
export type RemoveLocationFromContactMutationResult = Apollo.MutationResult<RemoveLocationFromContactMutation>;
export type RemoveLocationFromContactMutationOptions = Apollo.BaseMutationOptions<RemoveLocationFromContactMutation, RemoveLocationFromContactMutationVariables>;
export const RemoveOrganizationFromContactDocument = gql`
    mutation removeOrganizationFromContact($input: ContactOrganizationInput!) {
  contact_RemoveOrganizationById(input: $input) {
    ...ContactPersonalDetails
  }
}
    ${ContactPersonalDetailsFragmentDoc}`;
export type RemoveOrganizationFromContactMutationFn = Apollo.MutationFunction<RemoveOrganizationFromContactMutation, RemoveOrganizationFromContactMutationVariables>;

/**
 * __useRemoveOrganizationFromContactMutation__
 *
 * To run a mutation, you first call `useRemoveOrganizationFromContactMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemoveOrganizationFromContactMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removeOrganizationFromContactMutation, { data, loading, error }] = useRemoveOrganizationFromContactMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useRemoveOrganizationFromContactMutation(baseOptions?: Apollo.MutationHookOptions<RemoveOrganizationFromContactMutation, RemoveOrganizationFromContactMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RemoveOrganizationFromContactMutation, RemoveOrganizationFromContactMutationVariables>(RemoveOrganizationFromContactDocument, options);
      }
export type RemoveOrganizationFromContactMutationHookResult = ReturnType<typeof useRemoveOrganizationFromContactMutation>;
export type RemoveOrganizationFromContactMutationResult = Apollo.MutationResult<RemoveOrganizationFromContactMutation>;
export type RemoveOrganizationFromContactMutationOptions = Apollo.BaseMutationOptions<RemoveOrganizationFromContactMutation, RemoveOrganizationFromContactMutationVariables>;
export const RemovePhoneNumberFromContactDocument = gql`
    mutation removePhoneNumberFromContact($contactId: ID!, $id: ID!) {
  phoneNumberRemoveFromContactById(contactId: $contactId, id: $id) {
    result
  }
}
    `;
export type RemovePhoneNumberFromContactMutationFn = Apollo.MutationFunction<RemovePhoneNumberFromContactMutation, RemovePhoneNumberFromContactMutationVariables>;

/**
 * __useRemovePhoneNumberFromContactMutation__
 *
 * To run a mutation, you first call `useRemovePhoneNumberFromContactMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemovePhoneNumberFromContactMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removePhoneNumberFromContactMutation, { data, loading, error }] = useRemovePhoneNumberFromContactMutation({
 *   variables: {
 *      contactId: // value for 'contactId'
 *      id: // value for 'id'
 *   },
 * });
 */
export function useRemovePhoneNumberFromContactMutation(baseOptions?: Apollo.MutationHookOptions<RemovePhoneNumberFromContactMutation, RemovePhoneNumberFromContactMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RemovePhoneNumberFromContactMutation, RemovePhoneNumberFromContactMutationVariables>(RemovePhoneNumberFromContactDocument, options);
      }
export type RemovePhoneNumberFromContactMutationHookResult = ReturnType<typeof useRemovePhoneNumberFromContactMutation>;
export type RemovePhoneNumberFromContactMutationResult = Apollo.MutationResult<RemovePhoneNumberFromContactMutation>;
export type RemovePhoneNumberFromContactMutationOptions = Apollo.BaseMutationOptions<RemovePhoneNumberFromContactMutation, RemovePhoneNumberFromContactMutationVariables>;
export const RemoveTagFromContactDocument = gql`
    mutation RemoveTagFromContact($input: ContactTagInput!) {
  contact_RemoveTagById(input: $input) {
    id
  }
}
    `;
export type RemoveTagFromContactMutationFn = Apollo.MutationFunction<RemoveTagFromContactMutation, RemoveTagFromContactMutationVariables>;

/**
 * __useRemoveTagFromContactMutation__
 *
 * To run a mutation, you first call `useRemoveTagFromContactMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemoveTagFromContactMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removeTagFromContactMutation, { data, loading, error }] = useRemoveTagFromContactMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useRemoveTagFromContactMutation(baseOptions?: Apollo.MutationHookOptions<RemoveTagFromContactMutation, RemoveTagFromContactMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RemoveTagFromContactMutation, RemoveTagFromContactMutationVariables>(RemoveTagFromContactDocument, options);
      }
export type RemoveTagFromContactMutationHookResult = ReturnType<typeof useRemoveTagFromContactMutation>;
export type RemoveTagFromContactMutationResult = Apollo.MutationResult<RemoveTagFromContactMutation>;
export type RemoveTagFromContactMutationOptions = Apollo.BaseMutationOptions<RemoveTagFromContactMutation, RemoveTagFromContactMutationVariables>;
export const UpdateContactEmailDocument = gql`
    mutation updateContactEmail($contactId: ID!, $input: EmailUpdateInput!) {
  emailUpdateInContact(contactId: $contactId, input: $input) {
    primary
    label
    email
    id
  }
}
    `;
export type UpdateContactEmailMutationFn = Apollo.MutationFunction<UpdateContactEmailMutation, UpdateContactEmailMutationVariables>;

/**
 * __useUpdateContactEmailMutation__
 *
 * To run a mutation, you first call `useUpdateContactEmailMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateContactEmailMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateContactEmailMutation, { data, loading, error }] = useUpdateContactEmailMutation({
 *   variables: {
 *      contactId: // value for 'contactId'
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateContactEmailMutation(baseOptions?: Apollo.MutationHookOptions<UpdateContactEmailMutation, UpdateContactEmailMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateContactEmailMutation, UpdateContactEmailMutationVariables>(UpdateContactEmailDocument, options);
      }
export type UpdateContactEmailMutationHookResult = ReturnType<typeof useUpdateContactEmailMutation>;
export type UpdateContactEmailMutationResult = Apollo.MutationResult<UpdateContactEmailMutation>;
export type UpdateContactEmailMutationOptions = Apollo.BaseMutationOptions<UpdateContactEmailMutation, UpdateContactEmailMutationVariables>;
export const UpdateJobRoleDocument = gql`
    mutation updateJobRole($contactId: ID!, $input: JobRoleUpdateInput!) {
  jobRole_Update(contactId: $contactId, input: $input) {
    ...JobRole
    organization {
      id
      name
    }
  }
}
    ${JobRoleFragmentDoc}`;
export type UpdateJobRoleMutationFn = Apollo.MutationFunction<UpdateJobRoleMutation, UpdateJobRoleMutationVariables>;

/**
 * __useUpdateJobRoleMutation__
 *
 * To run a mutation, you first call `useUpdateJobRoleMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateJobRoleMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateJobRoleMutation, { data, loading, error }] = useUpdateJobRoleMutation({
 *   variables: {
 *      contactId: // value for 'contactId'
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateJobRoleMutation(baseOptions?: Apollo.MutationHookOptions<UpdateJobRoleMutation, UpdateJobRoleMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateJobRoleMutation, UpdateJobRoleMutationVariables>(UpdateJobRoleDocument, options);
      }
export type UpdateJobRoleMutationHookResult = ReturnType<typeof useUpdateJobRoleMutation>;
export type UpdateJobRoleMutationResult = Apollo.MutationResult<UpdateJobRoleMutation>;
export type UpdateJobRoleMutationOptions = Apollo.BaseMutationOptions<UpdateJobRoleMutation, UpdateJobRoleMutationVariables>;
export const UpdateContactPersonalDetailsDocument = gql`
    mutation updateContactPersonalDetails($input: ContactUpdateInput!) {
  contact_Update(input: $input) {
    id
    title
    firstName
    lastName
  }
}
    `;
export type UpdateContactPersonalDetailsMutationFn = Apollo.MutationFunction<UpdateContactPersonalDetailsMutation, UpdateContactPersonalDetailsMutationVariables>;

/**
 * __useUpdateContactPersonalDetailsMutation__
 *
 * To run a mutation, you first call `useUpdateContactPersonalDetailsMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateContactPersonalDetailsMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateContactPersonalDetailsMutation, { data, loading, error }] = useUpdateContactPersonalDetailsMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateContactPersonalDetailsMutation(baseOptions?: Apollo.MutationHookOptions<UpdateContactPersonalDetailsMutation, UpdateContactPersonalDetailsMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateContactPersonalDetailsMutation, UpdateContactPersonalDetailsMutationVariables>(UpdateContactPersonalDetailsDocument, options);
      }
export type UpdateContactPersonalDetailsMutationHookResult = ReturnType<typeof useUpdateContactPersonalDetailsMutation>;
export type UpdateContactPersonalDetailsMutationResult = Apollo.MutationResult<UpdateContactPersonalDetailsMutation>;
export type UpdateContactPersonalDetailsMutationOptions = Apollo.BaseMutationOptions<UpdateContactPersonalDetailsMutation, UpdateContactPersonalDetailsMutationVariables>;
export const UpdateContactPhoneNumberDocument = gql`
    mutation updateContactPhoneNumber($contactId: ID!, $input: PhoneNumberUpdateInput!) {
  phoneNumberUpdateInContact(contactId: $contactId, input: $input) {
    ...PhoneNumber
    label
  }
}
    ${PhoneNumberFragmentDoc}`;
export type UpdateContactPhoneNumberMutationFn = Apollo.MutationFunction<UpdateContactPhoneNumberMutation, UpdateContactPhoneNumberMutationVariables>;

/**
 * __useUpdateContactPhoneNumberMutation__
 *
 * To run a mutation, you first call `useUpdateContactPhoneNumberMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateContactPhoneNumberMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateContactPhoneNumberMutation, { data, loading, error }] = useUpdateContactPhoneNumberMutation({
 *   variables: {
 *      contactId: // value for 'contactId'
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateContactPhoneNumberMutation(baseOptions?: Apollo.MutationHookOptions<UpdateContactPhoneNumberMutation, UpdateContactPhoneNumberMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateContactPhoneNumberMutation, UpdateContactPhoneNumberMutationVariables>(UpdateContactPhoneNumberDocument, options);
      }
export type UpdateContactPhoneNumberMutationHookResult = ReturnType<typeof useUpdateContactPhoneNumberMutation>;
export type UpdateContactPhoneNumberMutationResult = Apollo.MutationResult<UpdateContactPhoneNumberMutation>;
export type UpdateContactPhoneNumberMutationOptions = Apollo.BaseMutationOptions<UpdateContactPhoneNumberMutation, UpdateContactPhoneNumberMutationVariables>;
export const DashboardView_ContactsDocument = gql`
    query dashboardView_Contacts($pagination: Pagination!, $where: Filter, $sort: SortBy) {
  dashboardView_Contacts(pagination: $pagination, where: $where, sort: $sort) {
    content {
      ...ContactPersonalDetails
      ...ContactCommunicationChannelsDetails
      locations {
        ...LocationBaseDetails
        rawAddress
      }
    }
    totalElements
  }
}
    ${ContactPersonalDetailsFragmentDoc}
${ContactCommunicationChannelsDetailsFragmentDoc}
${LocationBaseDetailsFragmentDoc}`;

/**
 * __useDashboardView_ContactsQuery__
 *
 * To run a query within a React component, call `useDashboardView_ContactsQuery` and pass it any options that fit your needs.
 * When your component renders, `useDashboardView_ContactsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useDashboardView_ContactsQuery({
 *   variables: {
 *      pagination: // value for 'pagination'
 *      where: // value for 'where'
 *      sort: // value for 'sort'
 *   },
 * });
 */
export function useDashboardView_ContactsQuery(baseOptions: Apollo.QueryHookOptions<DashboardView_ContactsQuery, DashboardView_ContactsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<DashboardView_ContactsQuery, DashboardView_ContactsQueryVariables>(DashboardView_ContactsDocument, options);
      }
export function useDashboardView_ContactsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<DashboardView_ContactsQuery, DashboardView_ContactsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<DashboardView_ContactsQuery, DashboardView_ContactsQueryVariables>(DashboardView_ContactsDocument, options);
        }
export type DashboardView_ContactsQueryHookResult = ReturnType<typeof useDashboardView_ContactsQuery>;
export type DashboardView_ContactsLazyQueryHookResult = ReturnType<typeof useDashboardView_ContactsLazyQuery>;
export type DashboardView_ContactsQueryResult = Apollo.QueryResult<DashboardView_ContactsQuery, DashboardView_ContactsQueryVariables>;
export const DashboardView_OrganizationsDocument = gql`
    query dashboardView_Organizations($pagination: Pagination!, $where: Filter, $sort: SortBy) {
  dashboardView_Organizations(pagination: $pagination, where: $where, sort: $sort) {
    content {
      id
      name
      subsidiaryOf {
        organization {
          id
          name
        }
      }
      owner {
        id
        firstName
        lastName
      }
      description
      industry
      website
      domains
      accountDetails {
        renewalForecast {
          amount
          potentialAmount
          comment
          updatedAt
          updatedById
          updatedBy {
            id
            firstName
            lastName
            emails {
              email
            }
          }
        }
        renewalLikelihood {
          probability
          previousProbability
          comment
          updatedById
          updatedBy {
            id
            firstName
            lastName
            emails {
              email
            }
          }
          updatedAt
        }
        billingDetails {
          renewalCycle
          frequency
          amount
          renewalCycleNext
        }
      }
      locations {
        ...LocationBaseDetails
        rawAddress
      }
      relationshipStages {
        relationship
        stage
      }
      lastTouchPointTimelineEventId
      lastTouchPointAt
      lastTouchPointTimelineEvent {
        ... on PageView {
          id
        }
        ... on Issue {
          id
        }
        ... on Note {
          id
          createdBy {
            firstName
            lastName
          }
        }
        ... on InteractionEvent {
          id
          channel
          eventType
          externalLinks {
            type
          }
          sentBy {
            __typename
            ... on EmailParticipant {
              type
              emailParticipant {
                id
                email
                rawEmail
              }
            }
            ... on ContactParticipant {
              contactParticipant {
                id
                name
                firstName
                lastName
              }
            }
            ... on JobRoleParticipant {
              jobRoleParticipant {
                contact {
                  id
                  name
                  firstName
                  lastName
                }
              }
            }
            ... on UserParticipant {
              userParticipant {
                id
                firstName
                lastName
              }
            }
          }
        }
        ... on Analysis {
          id
        }
        ... on Meeting {
          id
          name
          attendedBy {
            __typename
          }
        }
        ... on Action {
          id
          actionType
          createdAt
          source
          createdBy {
            id
            firstName
            lastName
          }
        }
      }
    }
    totalElements
  }
}
    ${LocationBaseDetailsFragmentDoc}`;

/**
 * __useDashboardView_OrganizationsQuery__
 *
 * To run a query within a React component, call `useDashboardView_OrganizationsQuery` and pass it any options that fit your needs.
 * When your component renders, `useDashboardView_OrganizationsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useDashboardView_OrganizationsQuery({
 *   variables: {
 *      pagination: // value for 'pagination'
 *      where: // value for 'where'
 *      sort: // value for 'sort'
 *   },
 * });
 */
export function useDashboardView_OrganizationsQuery(baseOptions: Apollo.QueryHookOptions<DashboardView_OrganizationsQuery, DashboardView_OrganizationsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<DashboardView_OrganizationsQuery, DashboardView_OrganizationsQueryVariables>(DashboardView_OrganizationsDocument, options);
      }
export function useDashboardView_OrganizationsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<DashboardView_OrganizationsQuery, DashboardView_OrganizationsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<DashboardView_OrganizationsQuery, DashboardView_OrganizationsQueryVariables>(DashboardView_OrganizationsDocument, options);
        }
export type DashboardView_OrganizationsQueryHookResult = ReturnType<typeof useDashboardView_OrganizationsQuery>;
export type DashboardView_OrganizationsLazyQueryHookResult = ReturnType<typeof useDashboardView_OrganizationsLazyQuery>;
export type DashboardView_OrganizationsQueryResult = Apollo.QueryResult<DashboardView_OrganizationsQuery, DashboardView_OrganizationsQueryVariables>;
export const GCliSearchDocument = gql`
    query gCliSearch($limit: Int, $keyword: String!) {
  gcli_Search(limit: $limit, keyword: $keyword) {
    id
    type
    display
    data {
      key
      value
      display
    }
  }
}
    `;

/**
 * __useGCliSearchQuery__
 *
 * To run a query within a React component, call `useGCliSearchQuery` and pass it any options that fit your needs.
 * When your component renders, `useGCliSearchQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGCliSearchQuery({
 *   variables: {
 *      limit: // value for 'limit'
 *      keyword: // value for 'keyword'
 *   },
 * });
 */
export function useGCliSearchQuery(baseOptions: Apollo.QueryHookOptions<GCliSearchQuery, GCliSearchQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GCliSearchQuery, GCliSearchQueryVariables>(GCliSearchDocument, options);
      }
export function useGCliSearchLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GCliSearchQuery, GCliSearchQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GCliSearchQuery, GCliSearchQueryVariables>(GCliSearchDocument, options);
        }
export type GCliSearchQueryHookResult = ReturnType<typeof useGCliSearchQuery>;
export type GCliSearchLazyQueryHookResult = ReturnType<typeof useGCliSearchLazyQuery>;
export type GCliSearchQueryResult = Apollo.QueryResult<GCliSearchQuery, GCliSearchQueryVariables>;
export const Global_CacheDocument = gql`
    query global_Cache {
  global_Cache {
    user {
      id
      emails {
        email
        rawEmail
        primary
      }
      firstName
      lastName
    }
    isOwner
    gCliCache {
      id
      type
      display
      data {
        key
        value
        display
      }
    }
  }
}
    `;

/**
 * __useGlobal_CacheQuery__
 *
 * To run a query within a React component, call `useGlobal_CacheQuery` and pass it any options that fit your needs.
 * When your component renders, `useGlobal_CacheQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGlobal_CacheQuery({
 *   variables: {
 *   },
 * });
 */
export function useGlobal_CacheQuery(baseOptions?: Apollo.QueryHookOptions<Global_CacheQuery, Global_CacheQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<Global_CacheQuery, Global_CacheQueryVariables>(Global_CacheDocument, options);
      }
export function useGlobal_CacheLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<Global_CacheQuery, Global_CacheQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<Global_CacheQuery, Global_CacheQueryVariables>(Global_CacheDocument, options);
        }
export type Global_CacheQueryHookResult = ReturnType<typeof useGlobal_CacheQuery>;
export type Global_CacheLazyQueryHookResult = ReturnType<typeof useGlobal_CacheLazyQuery>;
export type Global_CacheQueryResult = Apollo.QueryResult<Global_CacheQuery, Global_CacheQueryVariables>;
export const AddEmailToOrganizationDocument = gql`
    mutation addEmailToOrganization($organizationId: ID!, $input: EmailInput!) {
  emailMergeToOrganization(organizationId: $organizationId, input: $input) {
    ...Email
    label
  }
}
    ${EmailFragmentDoc}`;
export type AddEmailToOrganizationMutationFn = Apollo.MutationFunction<AddEmailToOrganizationMutation, AddEmailToOrganizationMutationVariables>;

/**
 * __useAddEmailToOrganizationMutation__
 *
 * To run a mutation, you first call `useAddEmailToOrganizationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddEmailToOrganizationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addEmailToOrganizationMutation, { data, loading, error }] = useAddEmailToOrganizationMutation({
 *   variables: {
 *      organizationId: // value for 'organizationId'
 *      input: // value for 'input'
 *   },
 * });
 */
export function useAddEmailToOrganizationMutation(baseOptions?: Apollo.MutationHookOptions<AddEmailToOrganizationMutation, AddEmailToOrganizationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddEmailToOrganizationMutation, AddEmailToOrganizationMutationVariables>(AddEmailToOrganizationDocument, options);
      }
export type AddEmailToOrganizationMutationHookResult = ReturnType<typeof useAddEmailToOrganizationMutation>;
export type AddEmailToOrganizationMutationResult = Apollo.MutationResult<AddEmailToOrganizationMutation>;
export type AddEmailToOrganizationMutationOptions = Apollo.BaseMutationOptions<AddEmailToOrganizationMutation, AddEmailToOrganizationMutationVariables>;
export const AddLocationToOrganizationDocument = gql`
    mutation addLocationToOrganization($organzationId: ID!) {
  organization_AddNewLocation(organizationId: $organzationId) {
    id
  }
}
    `;
export type AddLocationToOrganizationMutationFn = Apollo.MutationFunction<AddLocationToOrganizationMutation, AddLocationToOrganizationMutationVariables>;

/**
 * __useAddLocationToOrganizationMutation__
 *
 * To run a mutation, you first call `useAddLocationToOrganizationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddLocationToOrganizationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addLocationToOrganizationMutation, { data, loading, error }] = useAddLocationToOrganizationMutation({
 *   variables: {
 *      organzationId: // value for 'organzationId'
 *   },
 * });
 */
export function useAddLocationToOrganizationMutation(baseOptions?: Apollo.MutationHookOptions<AddLocationToOrganizationMutation, AddLocationToOrganizationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddLocationToOrganizationMutation, AddLocationToOrganizationMutationVariables>(AddLocationToOrganizationDocument, options);
      }
export type AddLocationToOrganizationMutationHookResult = ReturnType<typeof useAddLocationToOrganizationMutation>;
export type AddLocationToOrganizationMutationResult = Apollo.MutationResult<AddLocationToOrganizationMutation>;
export type AddLocationToOrganizationMutationOptions = Apollo.BaseMutationOptions<AddLocationToOrganizationMutation, AddLocationToOrganizationMutationVariables>;
export const AddPhoneToOrganizationDocument = gql`
    mutation addPhoneToOrganization($organizationId: ID!, $input: PhoneNumberInput!) {
  phoneNumberMergeToOrganization(organizationId: $organizationId, input: $input) {
    ...PhoneNumber
    label
  }
}
    ${PhoneNumberFragmentDoc}`;
export type AddPhoneToOrganizationMutationFn = Apollo.MutationFunction<AddPhoneToOrganizationMutation, AddPhoneToOrganizationMutationVariables>;

/**
 * __useAddPhoneToOrganizationMutation__
 *
 * To run a mutation, you first call `useAddPhoneToOrganizationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddPhoneToOrganizationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addPhoneToOrganizationMutation, { data, loading, error }] = useAddPhoneToOrganizationMutation({
 *   variables: {
 *      organizationId: // value for 'organizationId'
 *      input: // value for 'input'
 *   },
 * });
 */
export function useAddPhoneToOrganizationMutation(baseOptions?: Apollo.MutationHookOptions<AddPhoneToOrganizationMutation, AddPhoneToOrganizationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddPhoneToOrganizationMutation, AddPhoneToOrganizationMutationVariables>(AddPhoneToOrganizationDocument, options);
      }
export type AddPhoneToOrganizationMutationHookResult = ReturnType<typeof useAddPhoneToOrganizationMutation>;
export type AddPhoneToOrganizationMutationResult = Apollo.MutationResult<AddPhoneToOrganizationMutation>;
export type AddPhoneToOrganizationMutationOptions = Apollo.BaseMutationOptions<AddPhoneToOrganizationMutation, AddPhoneToOrganizationMutationVariables>;
export const AddRelationshipToOrganizationDocument = gql`
    mutation addRelationshipToOrganization($organizationId: ID!, $relationship: OrganizationRelationship!) {
  organization_AddRelationship(
    organizationId: $organizationId
    relationship: $relationship
  ) {
    id
  }
}
    `;
export type AddRelationshipToOrganizationMutationFn = Apollo.MutationFunction<AddRelationshipToOrganizationMutation, AddRelationshipToOrganizationMutationVariables>;

/**
 * __useAddRelationshipToOrganizationMutation__
 *
 * To run a mutation, you first call `useAddRelationshipToOrganizationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddRelationshipToOrganizationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addRelationshipToOrganizationMutation, { data, loading, error }] = useAddRelationshipToOrganizationMutation({
 *   variables: {
 *      organizationId: // value for 'organizationId'
 *      relationship: // value for 'relationship'
 *   },
 * });
 */
export function useAddRelationshipToOrganizationMutation(baseOptions?: Apollo.MutationHookOptions<AddRelationshipToOrganizationMutation, AddRelationshipToOrganizationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddRelationshipToOrganizationMutation, AddRelationshipToOrganizationMutationVariables>(AddRelationshipToOrganizationDocument, options);
      }
export type AddRelationshipToOrganizationMutationHookResult = ReturnType<typeof useAddRelationshipToOrganizationMutation>;
export type AddRelationshipToOrganizationMutationResult = Apollo.MutationResult<AddRelationshipToOrganizationMutation>;
export type AddRelationshipToOrganizationMutationOptions = Apollo.BaseMutationOptions<AddRelationshipToOrganizationMutation, AddRelationshipToOrganizationMutationVariables>;
export const AddOrganizationSubsidiaryDocument = gql`
    mutation addOrganizationSubsidiary($input: LinkOrganizationsInput!) {
  organization_AddSubsidiary(input: $input) {
    id
    subsidiaries {
      organization {
        id
        name
      }
    }
  }
}
    `;
export type AddOrganizationSubsidiaryMutationFn = Apollo.MutationFunction<AddOrganizationSubsidiaryMutation, AddOrganizationSubsidiaryMutationVariables>;

/**
 * __useAddOrganizationSubsidiaryMutation__
 *
 * To run a mutation, you first call `useAddOrganizationSubsidiaryMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddOrganizationSubsidiaryMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addOrganizationSubsidiaryMutation, { data, loading, error }] = useAddOrganizationSubsidiaryMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useAddOrganizationSubsidiaryMutation(baseOptions?: Apollo.MutationHookOptions<AddOrganizationSubsidiaryMutation, AddOrganizationSubsidiaryMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddOrganizationSubsidiaryMutation, AddOrganizationSubsidiaryMutationVariables>(AddOrganizationSubsidiaryDocument, options);
      }
export type AddOrganizationSubsidiaryMutationHookResult = ReturnType<typeof useAddOrganizationSubsidiaryMutation>;
export type AddOrganizationSubsidiaryMutationResult = Apollo.MutationResult<AddOrganizationSubsidiaryMutation>;
export type AddOrganizationSubsidiaryMutationOptions = Apollo.BaseMutationOptions<AddOrganizationSubsidiaryMutation, AddOrganizationSubsidiaryMutationVariables>;
export const CreateOrganizationDocument = gql`
    mutation createOrganization($input: OrganizationInput!) {
  organization_Create(input: $input) {
    id
    name
  }
}
    `;
export type CreateOrganizationMutationFn = Apollo.MutationFunction<CreateOrganizationMutation, CreateOrganizationMutationVariables>;

/**
 * __useCreateOrganizationMutation__
 *
 * To run a mutation, you first call `useCreateOrganizationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateOrganizationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createOrganizationMutation, { data, loading, error }] = useCreateOrganizationMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useCreateOrganizationMutation(baseOptions?: Apollo.MutationHookOptions<CreateOrganizationMutation, CreateOrganizationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreateOrganizationMutation, CreateOrganizationMutationVariables>(CreateOrganizationDocument, options);
      }
export type CreateOrganizationMutationHookResult = ReturnType<typeof useCreateOrganizationMutation>;
export type CreateOrganizationMutationResult = Apollo.MutationResult<CreateOrganizationMutation>;
export type CreateOrganizationMutationOptions = Apollo.BaseMutationOptions<CreateOrganizationMutation, CreateOrganizationMutationVariables>;
export const CreateOrganizationNoteDocument = gql`
    mutation createOrganizationNote($organizationId: ID!, $input: NoteInput!) {
  note_CreateForOrganization(organizationId: $organizationId, input: $input) {
    ...NoteContent
  }
}
    ${NoteContentFragmentDoc}`;
export type CreateOrganizationNoteMutationFn = Apollo.MutationFunction<CreateOrganizationNoteMutation, CreateOrganizationNoteMutationVariables>;

/**
 * __useCreateOrganizationNoteMutation__
 *
 * To run a mutation, you first call `useCreateOrganizationNoteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateOrganizationNoteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createOrganizationNoteMutation, { data, loading, error }] = useCreateOrganizationNoteMutation({
 *   variables: {
 *      organizationId: // value for 'organizationId'
 *      input: // value for 'input'
 *   },
 * });
 */
export function useCreateOrganizationNoteMutation(baseOptions?: Apollo.MutationHookOptions<CreateOrganizationNoteMutation, CreateOrganizationNoteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreateOrganizationNoteMutation, CreateOrganizationNoteMutationVariables>(CreateOrganizationNoteDocument, options);
      }
export type CreateOrganizationNoteMutationHookResult = ReturnType<typeof useCreateOrganizationNoteMutation>;
export type CreateOrganizationNoteMutationResult = Apollo.MutationResult<CreateOrganizationNoteMutation>;
export type CreateOrganizationNoteMutationOptions = Apollo.BaseMutationOptions<CreateOrganizationNoteMutation, CreateOrganizationNoteMutationVariables>;
export const GetOrganizationDocument = gql`
    query GetOrganization($id: ID!) {
  organization(id: $id) {
    ...OrganizationDetails
    ...OrganizationContacts
    owner {
      id
      firstName
      lastName
    }
    subsidiaryOf {
      organization {
        id
        name
      }
    }
    subsidiaries {
      organization {
        name
        id
      }
    }
    emails {
      id
      email
      primary
      label
    }
    phoneNumbers {
      id
      e164
      rawPhoneNumber
      label
    }
    industry
    industryGroup
    subIndustry
    customFields {
      id
      name
      datatype
      value
      template {
        type
      }
    }
  }
}
    ${OrganizationDetailsFragmentDoc}
${OrganizationContactsFragmentDoc}`;

/**
 * __useGetOrganizationQuery__
 *
 * To run a query within a React component, call `useGetOrganizationQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetOrganizationQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetOrganizationQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetOrganizationQuery(baseOptions: Apollo.QueryHookOptions<GetOrganizationQuery, GetOrganizationQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetOrganizationQuery, GetOrganizationQueryVariables>(GetOrganizationDocument, options);
      }
export function useGetOrganizationLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetOrganizationQuery, GetOrganizationQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetOrganizationQuery, GetOrganizationQueryVariables>(GetOrganizationDocument, options);
        }
export type GetOrganizationQueryHookResult = ReturnType<typeof useGetOrganizationQuery>;
export type GetOrganizationLazyQueryHookResult = ReturnType<typeof useGetOrganizationLazyQuery>;
export type GetOrganizationQueryResult = Apollo.QueryResult<GetOrganizationQuery, GetOrganizationQueryVariables>;
export const GetOrganizationCommunicationChannelsDocument = gql`
    query GetOrganizationCommunicationChannels($id: ID!) {
  organization(id: $id) {
    id
    name
    emails {
      id
      email
      primary
      label
    }
    phoneNumbers {
      id
      e164
      rawPhoneNumber
      label
    }
  }
}
    `;

/**
 * __useGetOrganizationCommunicationChannelsQuery__
 *
 * To run a query within a React component, call `useGetOrganizationCommunicationChannelsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetOrganizationCommunicationChannelsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetOrganizationCommunicationChannelsQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetOrganizationCommunicationChannelsQuery(baseOptions: Apollo.QueryHookOptions<GetOrganizationCommunicationChannelsQuery, GetOrganizationCommunicationChannelsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetOrganizationCommunicationChannelsQuery, GetOrganizationCommunicationChannelsQueryVariables>(GetOrganizationCommunicationChannelsDocument, options);
      }
export function useGetOrganizationCommunicationChannelsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetOrganizationCommunicationChannelsQuery, GetOrganizationCommunicationChannelsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetOrganizationCommunicationChannelsQuery, GetOrganizationCommunicationChannelsQueryVariables>(GetOrganizationCommunicationChannelsDocument, options);
        }
export type GetOrganizationCommunicationChannelsQueryHookResult = ReturnType<typeof useGetOrganizationCommunicationChannelsQuery>;
export type GetOrganizationCommunicationChannelsLazyQueryHookResult = ReturnType<typeof useGetOrganizationCommunicationChannelsLazyQuery>;
export type GetOrganizationCommunicationChannelsQueryResult = Apollo.QueryResult<GetOrganizationCommunicationChannelsQuery, GetOrganizationCommunicationChannelsQueryVariables>;
export const GetOrganizationContactsDocument = gql`
    query GetOrganizationContacts($id: ID!) {
  organization(id: $id) {
    ...OrganizationContacts
  }
}
    ${OrganizationContactsFragmentDoc}`;

/**
 * __useGetOrganizationContactsQuery__
 *
 * To run a query within a React component, call `useGetOrganizationContactsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetOrganizationContactsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetOrganizationContactsQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetOrganizationContactsQuery(baseOptions: Apollo.QueryHookOptions<GetOrganizationContactsQuery, GetOrganizationContactsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetOrganizationContactsQuery, GetOrganizationContactsQueryVariables>(GetOrganizationContactsDocument, options);
      }
export function useGetOrganizationContactsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetOrganizationContactsQuery, GetOrganizationContactsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetOrganizationContactsQuery, GetOrganizationContactsQueryVariables>(GetOrganizationContactsDocument, options);
        }
export type GetOrganizationContactsQueryHookResult = ReturnType<typeof useGetOrganizationContactsQuery>;
export type GetOrganizationContactsLazyQueryHookResult = ReturnType<typeof useGetOrganizationContactsLazyQuery>;
export type GetOrganizationContactsQueryResult = Apollo.QueryResult<GetOrganizationContactsQuery, GetOrganizationContactsQueryVariables>;
export const GetOrganizationCustomFieldsDocument = gql`
    query GetOrganizationCustomFields($id: ID!) {
  organization(id: $id) {
    customFields {
      id
      name
      datatype
      value
      template {
        type
      }
    }
  }
}
    `;

/**
 * __useGetOrganizationCustomFieldsQuery__
 *
 * To run a query within a React component, call `useGetOrganizationCustomFieldsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetOrganizationCustomFieldsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetOrganizationCustomFieldsQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetOrganizationCustomFieldsQuery(baseOptions: Apollo.QueryHookOptions<GetOrganizationCustomFieldsQuery, GetOrganizationCustomFieldsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetOrganizationCustomFieldsQuery, GetOrganizationCustomFieldsQueryVariables>(GetOrganizationCustomFieldsDocument, options);
      }
export function useGetOrganizationCustomFieldsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetOrganizationCustomFieldsQuery, GetOrganizationCustomFieldsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetOrganizationCustomFieldsQuery, GetOrganizationCustomFieldsQueryVariables>(GetOrganizationCustomFieldsDocument, options);
        }
export type GetOrganizationCustomFieldsQueryHookResult = ReturnType<typeof useGetOrganizationCustomFieldsQuery>;
export type GetOrganizationCustomFieldsLazyQueryHookResult = ReturnType<typeof useGetOrganizationCustomFieldsLazyQuery>;
export type GetOrganizationCustomFieldsQueryResult = Apollo.QueryResult<GetOrganizationCustomFieldsQuery, GetOrganizationCustomFieldsQueryVariables>;
export const GetOrganizationDetailsDocument = gql`
    query GetOrganizationDetails($id: ID!) {
  organization(id: $id) {
    ...OrganizationDetails
    subsidiaryOf {
      organization {
        id
        name
      }
    }
  }
}
    ${OrganizationDetailsFragmentDoc}`;

/**
 * __useGetOrganizationDetailsQuery__
 *
 * To run a query within a React component, call `useGetOrganizationDetailsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetOrganizationDetailsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetOrganizationDetailsQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetOrganizationDetailsQuery(baseOptions: Apollo.QueryHookOptions<GetOrganizationDetailsQuery, GetOrganizationDetailsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetOrganizationDetailsQuery, GetOrganizationDetailsQueryVariables>(GetOrganizationDetailsDocument, options);
      }
export function useGetOrganizationDetailsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetOrganizationDetailsQuery, GetOrganizationDetailsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetOrganizationDetailsQuery, GetOrganizationDetailsQueryVariables>(GetOrganizationDetailsDocument, options);
        }
export type GetOrganizationDetailsQueryHookResult = ReturnType<typeof useGetOrganizationDetailsQuery>;
export type GetOrganizationDetailsLazyQueryHookResult = ReturnType<typeof useGetOrganizationDetailsLazyQuery>;
export type GetOrganizationDetailsQueryResult = Apollo.QueryResult<GetOrganizationDetailsQuery, GetOrganizationDetailsQueryVariables>;
export const GetOrganizationLocationsDocument = gql`
    query GetOrganizationLocations($id: ID!) {
  organization(id: $id) {
    locations {
      ...LocationBaseDetails
      rawAddress
    }
  }
}
    ${LocationBaseDetailsFragmentDoc}`;

/**
 * __useGetOrganizationLocationsQuery__
 *
 * To run a query within a React component, call `useGetOrganizationLocationsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetOrganizationLocationsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetOrganizationLocationsQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetOrganizationLocationsQuery(baseOptions: Apollo.QueryHookOptions<GetOrganizationLocationsQuery, GetOrganizationLocationsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetOrganizationLocationsQuery, GetOrganizationLocationsQueryVariables>(GetOrganizationLocationsDocument, options);
      }
export function useGetOrganizationLocationsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetOrganizationLocationsQuery, GetOrganizationLocationsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetOrganizationLocationsQuery, GetOrganizationLocationsQueryVariables>(GetOrganizationLocationsDocument, options);
        }
export type GetOrganizationLocationsQueryHookResult = ReturnType<typeof useGetOrganizationLocationsQuery>;
export type GetOrganizationLocationsLazyQueryHookResult = ReturnType<typeof useGetOrganizationLocationsLazyQuery>;
export type GetOrganizationLocationsQueryResult = Apollo.QueryResult<GetOrganizationLocationsQuery, GetOrganizationLocationsQueryVariables>;
export const GetOrganizationMentionSuggestionsDocument = gql`
    query GetOrganizationMentionSuggestions($pagination: Pagination!, $where: Filter, $sort: [SortBy!]) {
  organizations(pagination: $pagination, where: $where, sort: $sort) {
    content {
      id
      name
    }
  }
}
    `;

/**
 * __useGetOrganizationMentionSuggestionsQuery__
 *
 * To run a query within a React component, call `useGetOrganizationMentionSuggestionsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetOrganizationMentionSuggestionsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetOrganizationMentionSuggestionsQuery({
 *   variables: {
 *      pagination: // value for 'pagination'
 *      where: // value for 'where'
 *      sort: // value for 'sort'
 *   },
 * });
 */
export function useGetOrganizationMentionSuggestionsQuery(baseOptions: Apollo.QueryHookOptions<GetOrganizationMentionSuggestionsQuery, GetOrganizationMentionSuggestionsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetOrganizationMentionSuggestionsQuery, GetOrganizationMentionSuggestionsQueryVariables>(GetOrganizationMentionSuggestionsDocument, options);
      }
export function useGetOrganizationMentionSuggestionsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetOrganizationMentionSuggestionsQuery, GetOrganizationMentionSuggestionsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetOrganizationMentionSuggestionsQuery, GetOrganizationMentionSuggestionsQueryVariables>(GetOrganizationMentionSuggestionsDocument, options);
        }
export type GetOrganizationMentionSuggestionsQueryHookResult = ReturnType<typeof useGetOrganizationMentionSuggestionsQuery>;
export type GetOrganizationMentionSuggestionsLazyQueryHookResult = ReturnType<typeof useGetOrganizationMentionSuggestionsLazyQuery>;
export type GetOrganizationMentionSuggestionsQueryResult = Apollo.QueryResult<GetOrganizationMentionSuggestionsQuery, GetOrganizationMentionSuggestionsQueryVariables>;
export const GetOrganizationNameDocument = gql`
    query GetOrganizationName($id: ID!) {
  organization(id: $id) {
    id
    name
  }
}
    `;

/**
 * __useGetOrganizationNameQuery__
 *
 * To run a query within a React component, call `useGetOrganizationNameQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetOrganizationNameQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetOrganizationNameQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetOrganizationNameQuery(baseOptions: Apollo.QueryHookOptions<GetOrganizationNameQuery, GetOrganizationNameQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetOrganizationNameQuery, GetOrganizationNameQueryVariables>(GetOrganizationNameDocument, options);
      }
export function useGetOrganizationNameLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetOrganizationNameQuery, GetOrganizationNameQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetOrganizationNameQuery, GetOrganizationNameQueryVariables>(GetOrganizationNameDocument, options);
        }
export type GetOrganizationNameQueryHookResult = ReturnType<typeof useGetOrganizationNameQuery>;
export type GetOrganizationNameLazyQueryHookResult = ReturnType<typeof useGetOrganizationNameLazyQuery>;
export type GetOrganizationNameQueryResult = Apollo.QueryResult<GetOrganizationNameQuery, GetOrganizationNameQueryVariables>;
export const GetOrganizationNotesDocument = gql`
    query GetOrganizationNotes($id: ID!, $pagination: Pagination) {
  organization(id: $id) {
    notes(pagination: $pagination) {
      content {
        ...NoteContent
      }
    }
  }
}
    ${NoteContentFragmentDoc}`;

/**
 * __useGetOrganizationNotesQuery__
 *
 * To run a query within a React component, call `useGetOrganizationNotesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetOrganizationNotesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetOrganizationNotesQuery({
 *   variables: {
 *      id: // value for 'id'
 *      pagination: // value for 'pagination'
 *   },
 * });
 */
export function useGetOrganizationNotesQuery(baseOptions: Apollo.QueryHookOptions<GetOrganizationNotesQuery, GetOrganizationNotesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetOrganizationNotesQuery, GetOrganizationNotesQueryVariables>(GetOrganizationNotesDocument, options);
      }
export function useGetOrganizationNotesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetOrganizationNotesQuery, GetOrganizationNotesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetOrganizationNotesQuery, GetOrganizationNotesQueryVariables>(GetOrganizationNotesDocument, options);
        }
export type GetOrganizationNotesQueryHookResult = ReturnType<typeof useGetOrganizationNotesQuery>;
export type GetOrganizationNotesLazyQueryHookResult = ReturnType<typeof useGetOrganizationNotesLazyQuery>;
export type GetOrganizationNotesQueryResult = Apollo.QueryResult<GetOrganizationNotesQuery, GetOrganizationNotesQueryVariables>;
export const GetOrganizationOwnerDocument = gql`
    query GetOrganizationOwner($id: ID!) {
  organization(id: $id) {
    id
    owner {
      id
      firstName
      lastName
    }
  }
}
    `;

/**
 * __useGetOrganizationOwnerQuery__
 *
 * To run a query within a React component, call `useGetOrganizationOwnerQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetOrganizationOwnerQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetOrganizationOwnerQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetOrganizationOwnerQuery(baseOptions: Apollo.QueryHookOptions<GetOrganizationOwnerQuery, GetOrganizationOwnerQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetOrganizationOwnerQuery, GetOrganizationOwnerQueryVariables>(GetOrganizationOwnerDocument, options);
      }
export function useGetOrganizationOwnerLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetOrganizationOwnerQuery, GetOrganizationOwnerQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetOrganizationOwnerQuery, GetOrganizationOwnerQueryVariables>(GetOrganizationOwnerDocument, options);
        }
export type GetOrganizationOwnerQueryHookResult = ReturnType<typeof useGetOrganizationOwnerQuery>;
export type GetOrganizationOwnerLazyQueryHookResult = ReturnType<typeof useGetOrganizationOwnerLazyQuery>;
export type GetOrganizationOwnerQueryResult = Apollo.QueryResult<GetOrganizationOwnerQuery, GetOrganizationOwnerQueryVariables>;
export const GetOrganizationSubsidiariesDocument = gql`
    query GetOrganizationSubsidiaries($id: ID!) {
  organization(id: $id) {
    subsidiaries {
      organization {
        name
        id
      }
    }
  }
}
    `;

/**
 * __useGetOrganizationSubsidiariesQuery__
 *
 * To run a query within a React component, call `useGetOrganizationSubsidiariesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetOrganizationSubsidiariesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetOrganizationSubsidiariesQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetOrganizationSubsidiariesQuery(baseOptions: Apollo.QueryHookOptions<GetOrganizationSubsidiariesQuery, GetOrganizationSubsidiariesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetOrganizationSubsidiariesQuery, GetOrganizationSubsidiariesQueryVariables>(GetOrganizationSubsidiariesDocument, options);
      }
export function useGetOrganizationSubsidiariesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetOrganizationSubsidiariesQuery, GetOrganizationSubsidiariesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetOrganizationSubsidiariesQuery, GetOrganizationSubsidiariesQueryVariables>(GetOrganizationSubsidiariesDocument, options);
        }
export type GetOrganizationSubsidiariesQueryHookResult = ReturnType<typeof useGetOrganizationSubsidiariesQuery>;
export type GetOrganizationSubsidiariesLazyQueryHookResult = ReturnType<typeof useGetOrganizationSubsidiariesLazyQuery>;
export type GetOrganizationSubsidiariesQueryResult = Apollo.QueryResult<GetOrganizationSubsidiariesQuery, GetOrganizationSubsidiariesQueryVariables>;
export const GetOrganizationTableDataDocument = gql`
    query getOrganizationTableData($pagination: Pagination, $where: Filter, $sort: [SortBy!]) {
  organizations(pagination: $pagination, where: $where, sort: $sort) {
    content {
      id
      name
      industry
      locations {
        ...LocationBaseDetails
      }
      subsidiaryOf {
        type
        organization {
          name
        }
      }
    }
    totalElements
    totalPages
  }
}
    ${LocationBaseDetailsFragmentDoc}`;

/**
 * __useGetOrganizationTableDataQuery__
 *
 * To run a query within a React component, call `useGetOrganizationTableDataQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetOrganizationTableDataQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetOrganizationTableDataQuery({
 *   variables: {
 *      pagination: // value for 'pagination'
 *      where: // value for 'where'
 *      sort: // value for 'sort'
 *   },
 * });
 */
export function useGetOrganizationTableDataQuery(baseOptions?: Apollo.QueryHookOptions<GetOrganizationTableDataQuery, GetOrganizationTableDataQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetOrganizationTableDataQuery, GetOrganizationTableDataQueryVariables>(GetOrganizationTableDataDocument, options);
      }
export function useGetOrganizationTableDataLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetOrganizationTableDataQuery, GetOrganizationTableDataQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetOrganizationTableDataQuery, GetOrganizationTableDataQueryVariables>(GetOrganizationTableDataDocument, options);
        }
export type GetOrganizationTableDataQueryHookResult = ReturnType<typeof useGetOrganizationTableDataQuery>;
export type GetOrganizationTableDataLazyQueryHookResult = ReturnType<typeof useGetOrganizationTableDataLazyQuery>;
export type GetOrganizationTableDataQueryResult = Apollo.QueryResult<GetOrganizationTableDataQuery, GetOrganizationTableDataQueryVariables>;
export const GetOrganizationsOptionsDocument = gql`
    query getOrganizationsOptions($pagination: Pagination) {
  organizations(pagination: $pagination) {
    content {
      id
      name
    }
  }
}
    `;

/**
 * __useGetOrganizationsOptionsQuery__
 *
 * To run a query within a React component, call `useGetOrganizationsOptionsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetOrganizationsOptionsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetOrganizationsOptionsQuery({
 *   variables: {
 *      pagination: // value for 'pagination'
 *   },
 * });
 */
export function useGetOrganizationsOptionsQuery(baseOptions?: Apollo.QueryHookOptions<GetOrganizationsOptionsQuery, GetOrganizationsOptionsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetOrganizationsOptionsQuery, GetOrganizationsOptionsQueryVariables>(GetOrganizationsOptionsDocument, options);
      }
export function useGetOrganizationsOptionsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetOrganizationsOptionsQuery, GetOrganizationsOptionsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetOrganizationsOptionsQuery, GetOrganizationsOptionsQueryVariables>(GetOrganizationsOptionsDocument, options);
        }
export type GetOrganizationsOptionsQueryHookResult = ReturnType<typeof useGetOrganizationsOptionsQuery>;
export type GetOrganizationsOptionsLazyQueryHookResult = ReturnType<typeof useGetOrganizationsOptionsLazyQuery>;
export type GetOrganizationsOptionsQueryResult = Apollo.QueryResult<GetOrganizationsOptionsQuery, GetOrganizationsOptionsQueryVariables>;
export const HideOrganizationsDocument = gql`
    mutation hideOrganizations($ids: [ID!]!) {
  organization_HideAll(ids: $ids) {
    result
  }
}
    `;
export type HideOrganizationsMutationFn = Apollo.MutationFunction<HideOrganizationsMutation, HideOrganizationsMutationVariables>;

/**
 * __useHideOrganizationsMutation__
 *
 * To run a mutation, you first call `useHideOrganizationsMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useHideOrganizationsMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [hideOrganizationsMutation, { data, loading, error }] = useHideOrganizationsMutation({
 *   variables: {
 *      ids: // value for 'ids'
 *   },
 * });
 */
export function useHideOrganizationsMutation(baseOptions?: Apollo.MutationHookOptions<HideOrganizationsMutation, HideOrganizationsMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<HideOrganizationsMutation, HideOrganizationsMutationVariables>(HideOrganizationsDocument, options);
      }
export type HideOrganizationsMutationHookResult = ReturnType<typeof useHideOrganizationsMutation>;
export type HideOrganizationsMutationResult = Apollo.MutationResult<HideOrganizationsMutation>;
export type HideOrganizationsMutationOptions = Apollo.BaseMutationOptions<HideOrganizationsMutation, HideOrganizationsMutationVariables>;
export const MergeOrganizationsDocument = gql`
    mutation mergeOrganizations($primaryOrganizationId: ID!, $mergedOrganizationIds: [ID!]!) {
  organization_Merge(
    primaryOrganizationId: $primaryOrganizationId
    mergedOrganizationIds: $mergedOrganizationIds
  ) {
    id
    ...OrganizationDetails
  }
}
    ${OrganizationDetailsFragmentDoc}`;
export type MergeOrganizationsMutationFn = Apollo.MutationFunction<MergeOrganizationsMutation, MergeOrganizationsMutationVariables>;

/**
 * __useMergeOrganizationsMutation__
 *
 * To run a mutation, you first call `useMergeOrganizationsMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useMergeOrganizationsMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [mergeOrganizationsMutation, { data, loading, error }] = useMergeOrganizationsMutation({
 *   variables: {
 *      primaryOrganizationId: // value for 'primaryOrganizationId'
 *      mergedOrganizationIds: // value for 'mergedOrganizationIds'
 *   },
 * });
 */
export function useMergeOrganizationsMutation(baseOptions?: Apollo.MutationHookOptions<MergeOrganizationsMutation, MergeOrganizationsMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<MergeOrganizationsMutation, MergeOrganizationsMutationVariables>(MergeOrganizationsDocument, options);
      }
export type MergeOrganizationsMutationHookResult = ReturnType<typeof useMergeOrganizationsMutation>;
export type MergeOrganizationsMutationResult = Apollo.MutationResult<MergeOrganizationsMutation>;
export type MergeOrganizationsMutationOptions = Apollo.BaseMutationOptions<MergeOrganizationsMutation, MergeOrganizationsMutationVariables>;
export const RemoveEmailFromOrganizationDocument = gql`
    mutation removeEmailFromOrganization($organizationId: ID!, $id: ID!) {
  emailRemoveFromOrganizationById(organizationId: $organizationId, id: $id) {
    result
  }
}
    `;
export type RemoveEmailFromOrganizationMutationFn = Apollo.MutationFunction<RemoveEmailFromOrganizationMutation, RemoveEmailFromOrganizationMutationVariables>;

/**
 * __useRemoveEmailFromOrganizationMutation__
 *
 * To run a mutation, you first call `useRemoveEmailFromOrganizationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemoveEmailFromOrganizationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removeEmailFromOrganizationMutation, { data, loading, error }] = useRemoveEmailFromOrganizationMutation({
 *   variables: {
 *      organizationId: // value for 'organizationId'
 *      id: // value for 'id'
 *   },
 * });
 */
export function useRemoveEmailFromOrganizationMutation(baseOptions?: Apollo.MutationHookOptions<RemoveEmailFromOrganizationMutation, RemoveEmailFromOrganizationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RemoveEmailFromOrganizationMutation, RemoveEmailFromOrganizationMutationVariables>(RemoveEmailFromOrganizationDocument, options);
      }
export type RemoveEmailFromOrganizationMutationHookResult = ReturnType<typeof useRemoveEmailFromOrganizationMutation>;
export type RemoveEmailFromOrganizationMutationResult = Apollo.MutationResult<RemoveEmailFromOrganizationMutation>;
export type RemoveEmailFromOrganizationMutationOptions = Apollo.BaseMutationOptions<RemoveEmailFromOrganizationMutation, RemoveEmailFromOrganizationMutationVariables>;
export const RemoveLocationFromOrganizationDocument = gql`
    mutation removeLocationFromOrganization($locationId: ID!, $organizationId: ID!) {
  location_RemoveFromOrganization(
    locationId: $locationId
    organizationId: $organizationId
  ) {
    id
    locations {
      ...LocationBaseDetails
    }
  }
}
    ${LocationBaseDetailsFragmentDoc}`;
export type RemoveLocationFromOrganizationMutationFn = Apollo.MutationFunction<RemoveLocationFromOrganizationMutation, RemoveLocationFromOrganizationMutationVariables>;

/**
 * __useRemoveLocationFromOrganizationMutation__
 *
 * To run a mutation, you first call `useRemoveLocationFromOrganizationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemoveLocationFromOrganizationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removeLocationFromOrganizationMutation, { data, loading, error }] = useRemoveLocationFromOrganizationMutation({
 *   variables: {
 *      locationId: // value for 'locationId'
 *      organizationId: // value for 'organizationId'
 *   },
 * });
 */
export function useRemoveLocationFromOrganizationMutation(baseOptions?: Apollo.MutationHookOptions<RemoveLocationFromOrganizationMutation, RemoveLocationFromOrganizationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RemoveLocationFromOrganizationMutation, RemoveLocationFromOrganizationMutationVariables>(RemoveLocationFromOrganizationDocument, options);
      }
export type RemoveLocationFromOrganizationMutationHookResult = ReturnType<typeof useRemoveLocationFromOrganizationMutation>;
export type RemoveLocationFromOrganizationMutationResult = Apollo.MutationResult<RemoveLocationFromOrganizationMutation>;
export type RemoveLocationFromOrganizationMutationOptions = Apollo.BaseMutationOptions<RemoveLocationFromOrganizationMutation, RemoveLocationFromOrganizationMutationVariables>;
export const RemoveOrganizationOwnerDocument = gql`
    mutation removeOrganizationOwner($organizationId: ID!) {
  organization_UnsetOwner(organizationId: $organizationId) {
    id
    owner {
      id
    }
  }
}
    `;
export type RemoveOrganizationOwnerMutationFn = Apollo.MutationFunction<RemoveOrganizationOwnerMutation, RemoveOrganizationOwnerMutationVariables>;

/**
 * __useRemoveOrganizationOwnerMutation__
 *
 * To run a mutation, you first call `useRemoveOrganizationOwnerMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemoveOrganizationOwnerMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removeOrganizationOwnerMutation, { data, loading, error }] = useRemoveOrganizationOwnerMutation({
 *   variables: {
 *      organizationId: // value for 'organizationId'
 *   },
 * });
 */
export function useRemoveOrganizationOwnerMutation(baseOptions?: Apollo.MutationHookOptions<RemoveOrganizationOwnerMutation, RemoveOrganizationOwnerMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RemoveOrganizationOwnerMutation, RemoveOrganizationOwnerMutationVariables>(RemoveOrganizationOwnerDocument, options);
      }
export type RemoveOrganizationOwnerMutationHookResult = ReturnType<typeof useRemoveOrganizationOwnerMutation>;
export type RemoveOrganizationOwnerMutationResult = Apollo.MutationResult<RemoveOrganizationOwnerMutation>;
export type RemoveOrganizationOwnerMutationOptions = Apollo.BaseMutationOptions<RemoveOrganizationOwnerMutation, RemoveOrganizationOwnerMutationVariables>;
export const RemoveOrganizationRelationshipDocument = gql`
    mutation removeOrganizationRelationship($organizationId: ID!, $relationship: OrganizationRelationship!) {
  organization_RemoveRelationship(
    organizationId: $organizationId
    relationship: $relationship
  ) {
    id
  }
}
    `;
export type RemoveOrganizationRelationshipMutationFn = Apollo.MutationFunction<RemoveOrganizationRelationshipMutation, RemoveOrganizationRelationshipMutationVariables>;

/**
 * __useRemoveOrganizationRelationshipMutation__
 *
 * To run a mutation, you first call `useRemoveOrganizationRelationshipMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemoveOrganizationRelationshipMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removeOrganizationRelationshipMutation, { data, loading, error }] = useRemoveOrganizationRelationshipMutation({
 *   variables: {
 *      organizationId: // value for 'organizationId'
 *      relationship: // value for 'relationship'
 *   },
 * });
 */
export function useRemoveOrganizationRelationshipMutation(baseOptions?: Apollo.MutationHookOptions<RemoveOrganizationRelationshipMutation, RemoveOrganizationRelationshipMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RemoveOrganizationRelationshipMutation, RemoveOrganizationRelationshipMutationVariables>(RemoveOrganizationRelationshipDocument, options);
      }
export type RemoveOrganizationRelationshipMutationHookResult = ReturnType<typeof useRemoveOrganizationRelationshipMutation>;
export type RemoveOrganizationRelationshipMutationResult = Apollo.MutationResult<RemoveOrganizationRelationshipMutation>;
export type RemoveOrganizationRelationshipMutationOptions = Apollo.BaseMutationOptions<RemoveOrganizationRelationshipMutation, RemoveOrganizationRelationshipMutationVariables>;
export const RemovePhoneNumberFromOrganizationDocument = gql`
    mutation removePhoneNumberFromOrganization($organizationId: ID!, $id: ID!) {
  phoneNumberRemoveFromOrganizationById(organizationId: $organizationId, id: $id) {
    result
  }
}
    `;
export type RemovePhoneNumberFromOrganizationMutationFn = Apollo.MutationFunction<RemovePhoneNumberFromOrganizationMutation, RemovePhoneNumberFromOrganizationMutationVariables>;

/**
 * __useRemovePhoneNumberFromOrganizationMutation__
 *
 * To run a mutation, you first call `useRemovePhoneNumberFromOrganizationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemovePhoneNumberFromOrganizationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removePhoneNumberFromOrganizationMutation, { data, loading, error }] = useRemovePhoneNumberFromOrganizationMutation({
 *   variables: {
 *      organizationId: // value for 'organizationId'
 *      id: // value for 'id'
 *   },
 * });
 */
export function useRemovePhoneNumberFromOrganizationMutation(baseOptions?: Apollo.MutationHookOptions<RemovePhoneNumberFromOrganizationMutation, RemovePhoneNumberFromOrganizationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RemovePhoneNumberFromOrganizationMutation, RemovePhoneNumberFromOrganizationMutationVariables>(RemovePhoneNumberFromOrganizationDocument, options);
      }
export type RemovePhoneNumberFromOrganizationMutationHookResult = ReturnType<typeof useRemovePhoneNumberFromOrganizationMutation>;
export type RemovePhoneNumberFromOrganizationMutationResult = Apollo.MutationResult<RemovePhoneNumberFromOrganizationMutation>;
export type RemovePhoneNumberFromOrganizationMutationOptions = Apollo.BaseMutationOptions<RemovePhoneNumberFromOrganizationMutation, RemovePhoneNumberFromOrganizationMutationVariables>;
export const RemoveStageFromOrganizationRelationshipDocument = gql`
    mutation removeStageFromOrganizationRelationship($organizationId: ID!, $relationship: OrganizationRelationship!) {
  organization_RemoveRelationshipStage(
    organizationId: $organizationId
    relationship: $relationship
  ) {
    id
  }
}
    `;
export type RemoveStageFromOrganizationRelationshipMutationFn = Apollo.MutationFunction<RemoveStageFromOrganizationRelationshipMutation, RemoveStageFromOrganizationRelationshipMutationVariables>;

/**
 * __useRemoveStageFromOrganizationRelationshipMutation__
 *
 * To run a mutation, you first call `useRemoveStageFromOrganizationRelationshipMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemoveStageFromOrganizationRelationshipMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removeStageFromOrganizationRelationshipMutation, { data, loading, error }] = useRemoveStageFromOrganizationRelationshipMutation({
 *   variables: {
 *      organizationId: // value for 'organizationId'
 *      relationship: // value for 'relationship'
 *   },
 * });
 */
export function useRemoveStageFromOrganizationRelationshipMutation(baseOptions?: Apollo.MutationHookOptions<RemoveStageFromOrganizationRelationshipMutation, RemoveStageFromOrganizationRelationshipMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RemoveStageFromOrganizationRelationshipMutation, RemoveStageFromOrganizationRelationshipMutationVariables>(RemoveStageFromOrganizationRelationshipDocument, options);
      }
export type RemoveStageFromOrganizationRelationshipMutationHookResult = ReturnType<typeof useRemoveStageFromOrganizationRelationshipMutation>;
export type RemoveStageFromOrganizationRelationshipMutationResult = Apollo.MutationResult<RemoveStageFromOrganizationRelationshipMutation>;
export type RemoveStageFromOrganizationRelationshipMutationOptions = Apollo.BaseMutationOptions<RemoveStageFromOrganizationRelationshipMutation, RemoveStageFromOrganizationRelationshipMutationVariables>;
export const RemoveOrganizationSubsidiaryDocument = gql`
    mutation removeOrganizationSubsidiary($organizationId: ID!, $subsidiaryId: ID!) {
  organization_RemoveSubsidiary(
    organizationId: $organizationId
    subsidiaryId: $subsidiaryId
  ) {
    id
    subsidiaries {
      organization {
        id
        name
      }
    }
  }
}
    `;
export type RemoveOrganizationSubsidiaryMutationFn = Apollo.MutationFunction<RemoveOrganizationSubsidiaryMutation, RemoveOrganizationSubsidiaryMutationVariables>;

/**
 * __useRemoveOrganizationSubsidiaryMutation__
 *
 * To run a mutation, you first call `useRemoveOrganizationSubsidiaryMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemoveOrganizationSubsidiaryMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removeOrganizationSubsidiaryMutation, { data, loading, error }] = useRemoveOrganizationSubsidiaryMutation({
 *   variables: {
 *      organizationId: // value for 'organizationId'
 *      subsidiaryId: // value for 'subsidiaryId'
 *   },
 * });
 */
export function useRemoveOrganizationSubsidiaryMutation(baseOptions?: Apollo.MutationHookOptions<RemoveOrganizationSubsidiaryMutation, RemoveOrganizationSubsidiaryMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RemoveOrganizationSubsidiaryMutation, RemoveOrganizationSubsidiaryMutationVariables>(RemoveOrganizationSubsidiaryDocument, options);
      }
export type RemoveOrganizationSubsidiaryMutationHookResult = ReturnType<typeof useRemoveOrganizationSubsidiaryMutation>;
export type RemoveOrganizationSubsidiaryMutationResult = Apollo.MutationResult<RemoveOrganizationSubsidiaryMutation>;
export type RemoveOrganizationSubsidiaryMutationOptions = Apollo.BaseMutationOptions<RemoveOrganizationSubsidiaryMutation, RemoveOrganizationSubsidiaryMutationVariables>;
export const SetStageToOrganizationRelationshipDocument = gql`
    mutation setStageToOrganizationRelationship($organizationId: ID!, $relationship: OrganizationRelationship!, $stage: String!) {
  organization_SetRelationshipStage(
    organizationId: $organizationId
    relationship: $relationship
    stage: $stage
  ) {
    id
  }
}
    `;
export type SetStageToOrganizationRelationshipMutationFn = Apollo.MutationFunction<SetStageToOrganizationRelationshipMutation, SetStageToOrganizationRelationshipMutationVariables>;

/**
 * __useSetStageToOrganizationRelationshipMutation__
 *
 * To run a mutation, you first call `useSetStageToOrganizationRelationshipMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useSetStageToOrganizationRelationshipMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [setStageToOrganizationRelationshipMutation, { data, loading, error }] = useSetStageToOrganizationRelationshipMutation({
 *   variables: {
 *      organizationId: // value for 'organizationId'
 *      relationship: // value for 'relationship'
 *      stage: // value for 'stage'
 *   },
 * });
 */
export function useSetStageToOrganizationRelationshipMutation(baseOptions?: Apollo.MutationHookOptions<SetStageToOrganizationRelationshipMutation, SetStageToOrganizationRelationshipMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<SetStageToOrganizationRelationshipMutation, SetStageToOrganizationRelationshipMutationVariables>(SetStageToOrganizationRelationshipDocument, options);
      }
export type SetStageToOrganizationRelationshipMutationHookResult = ReturnType<typeof useSetStageToOrganizationRelationshipMutation>;
export type SetStageToOrganizationRelationshipMutationResult = Apollo.MutationResult<SetStageToOrganizationRelationshipMutation>;
export type SetStageToOrganizationRelationshipMutationOptions = Apollo.BaseMutationOptions<SetStageToOrganizationRelationshipMutation, SetStageToOrganizationRelationshipMutationVariables>;
export const UpdateOrganizationDescriptionDocument = gql`
    mutation updateOrganizationDescription($input: OrganizationUpdateInput!) {
  organization_Update(input: $input) {
    id
    description
  }
}
    `;
export type UpdateOrganizationDescriptionMutationFn = Apollo.MutationFunction<UpdateOrganizationDescriptionMutation, UpdateOrganizationDescriptionMutationVariables>;

/**
 * __useUpdateOrganizationDescriptionMutation__
 *
 * To run a mutation, you first call `useUpdateOrganizationDescriptionMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateOrganizationDescriptionMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateOrganizationDescriptionMutation, { data, loading, error }] = useUpdateOrganizationDescriptionMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateOrganizationDescriptionMutation(baseOptions?: Apollo.MutationHookOptions<UpdateOrganizationDescriptionMutation, UpdateOrganizationDescriptionMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateOrganizationDescriptionMutation, UpdateOrganizationDescriptionMutationVariables>(UpdateOrganizationDescriptionDocument, options);
      }
export type UpdateOrganizationDescriptionMutationHookResult = ReturnType<typeof useUpdateOrganizationDescriptionMutation>;
export type UpdateOrganizationDescriptionMutationResult = Apollo.MutationResult<UpdateOrganizationDescriptionMutation>;
export type UpdateOrganizationDescriptionMutationOptions = Apollo.BaseMutationOptions<UpdateOrganizationDescriptionMutation, UpdateOrganizationDescriptionMutationVariables>;
export const UpdateOrganizationEmailDocument = gql`
    mutation updateOrganizationEmail($organizationId: ID!, $input: EmailUpdateInput!) {
  emailUpdateInOrganization(organizationId: $organizationId, input: $input) {
    primary
    label
    id
    email
  }
}
    `;
export type UpdateOrganizationEmailMutationFn = Apollo.MutationFunction<UpdateOrganizationEmailMutation, UpdateOrganizationEmailMutationVariables>;

/**
 * __useUpdateOrganizationEmailMutation__
 *
 * To run a mutation, you first call `useUpdateOrganizationEmailMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateOrganizationEmailMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateOrganizationEmailMutation, { data, loading, error }] = useUpdateOrganizationEmailMutation({
 *   variables: {
 *      organizationId: // value for 'organizationId'
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateOrganizationEmailMutation(baseOptions?: Apollo.MutationHookOptions<UpdateOrganizationEmailMutation, UpdateOrganizationEmailMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateOrganizationEmailMutation, UpdateOrganizationEmailMutationVariables>(UpdateOrganizationEmailDocument, options);
      }
export type UpdateOrganizationEmailMutationHookResult = ReturnType<typeof useUpdateOrganizationEmailMutation>;
export type UpdateOrganizationEmailMutationResult = Apollo.MutationResult<UpdateOrganizationEmailMutation>;
export type UpdateOrganizationEmailMutationOptions = Apollo.BaseMutationOptions<UpdateOrganizationEmailMutation, UpdateOrganizationEmailMutationVariables>;
export const UpdateOrganizationIndustryDocument = gql`
    mutation updateOrganizationIndustry($input: OrganizationUpdateInput!) {
  organization_Update(input: $input) {
    id
    industry
  }
}
    `;
export type UpdateOrganizationIndustryMutationFn = Apollo.MutationFunction<UpdateOrganizationIndustryMutation, UpdateOrganizationIndustryMutationVariables>;

/**
 * __useUpdateOrganizationIndustryMutation__
 *
 * To run a mutation, you first call `useUpdateOrganizationIndustryMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateOrganizationIndustryMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateOrganizationIndustryMutation, { data, loading, error }] = useUpdateOrganizationIndustryMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateOrganizationIndustryMutation(baseOptions?: Apollo.MutationHookOptions<UpdateOrganizationIndustryMutation, UpdateOrganizationIndustryMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateOrganizationIndustryMutation, UpdateOrganizationIndustryMutationVariables>(UpdateOrganizationIndustryDocument, options);
      }
export type UpdateOrganizationIndustryMutationHookResult = ReturnType<typeof useUpdateOrganizationIndustryMutation>;
export type UpdateOrganizationIndustryMutationResult = Apollo.MutationResult<UpdateOrganizationIndustryMutation>;
export type UpdateOrganizationIndustryMutationOptions = Apollo.BaseMutationOptions<UpdateOrganizationIndustryMutation, UpdateOrganizationIndustryMutationVariables>;
export const UpdateOrganizationNameDocument = gql`
    mutation updateOrganizationName($input: OrganizationUpdateInput!) {
  organization_Update(input: $input) {
    id
    name
  }
}
    `;
export type UpdateOrganizationNameMutationFn = Apollo.MutationFunction<UpdateOrganizationNameMutation, UpdateOrganizationNameMutationVariables>;

/**
 * __useUpdateOrganizationNameMutation__
 *
 * To run a mutation, you first call `useUpdateOrganizationNameMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateOrganizationNameMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateOrganizationNameMutation, { data, loading, error }] = useUpdateOrganizationNameMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateOrganizationNameMutation(baseOptions?: Apollo.MutationHookOptions<UpdateOrganizationNameMutation, UpdateOrganizationNameMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateOrganizationNameMutation, UpdateOrganizationNameMutationVariables>(UpdateOrganizationNameDocument, options);
      }
export type UpdateOrganizationNameMutationHookResult = ReturnType<typeof useUpdateOrganizationNameMutation>;
export type UpdateOrganizationNameMutationResult = Apollo.MutationResult<UpdateOrganizationNameMutation>;
export type UpdateOrganizationNameMutationOptions = Apollo.BaseMutationOptions<UpdateOrganizationNameMutation, UpdateOrganizationNameMutationVariables>;
export const UpdateOrganizationOwnerDocument = gql`
    mutation updateOrganizationOwner($organizationId: ID!, $userId: ID!) {
  organization_SetOwner(organizationId: $organizationId, userId: $userId) {
    id
    owner {
      id
      firstName
      lastName
    }
  }
}
    `;
export type UpdateOrganizationOwnerMutationFn = Apollo.MutationFunction<UpdateOrganizationOwnerMutation, UpdateOrganizationOwnerMutationVariables>;

/**
 * __useUpdateOrganizationOwnerMutation__
 *
 * To run a mutation, you first call `useUpdateOrganizationOwnerMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateOrganizationOwnerMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateOrganizationOwnerMutation, { data, loading, error }] = useUpdateOrganizationOwnerMutation({
 *   variables: {
 *      organizationId: // value for 'organizationId'
 *      userId: // value for 'userId'
 *   },
 * });
 */
export function useUpdateOrganizationOwnerMutation(baseOptions?: Apollo.MutationHookOptions<UpdateOrganizationOwnerMutation, UpdateOrganizationOwnerMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateOrganizationOwnerMutation, UpdateOrganizationOwnerMutationVariables>(UpdateOrganizationOwnerDocument, options);
      }
export type UpdateOrganizationOwnerMutationHookResult = ReturnType<typeof useUpdateOrganizationOwnerMutation>;
export type UpdateOrganizationOwnerMutationResult = Apollo.MutationResult<UpdateOrganizationOwnerMutation>;
export type UpdateOrganizationOwnerMutationOptions = Apollo.BaseMutationOptions<UpdateOrganizationOwnerMutation, UpdateOrganizationOwnerMutationVariables>;
export const UpdateOrganizationPhoneNumberDocument = gql`
    mutation updateOrganizationPhoneNumber($organizationId: ID!, $input: PhoneNumberUpdateInput!) {
  phoneNumberUpdateInOrganization(organizationId: $organizationId, input: $input) {
    ...PhoneNumber
    label
    primary
  }
}
    ${PhoneNumberFragmentDoc}`;
export type UpdateOrganizationPhoneNumberMutationFn = Apollo.MutationFunction<UpdateOrganizationPhoneNumberMutation, UpdateOrganizationPhoneNumberMutationVariables>;

/**
 * __useUpdateOrganizationPhoneNumberMutation__
 *
 * To run a mutation, you first call `useUpdateOrganizationPhoneNumberMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateOrganizationPhoneNumberMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateOrganizationPhoneNumberMutation, { data, loading, error }] = useUpdateOrganizationPhoneNumberMutation({
 *   variables: {
 *      organizationId: // value for 'organizationId'
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateOrganizationPhoneNumberMutation(baseOptions?: Apollo.MutationHookOptions<UpdateOrganizationPhoneNumberMutation, UpdateOrganizationPhoneNumberMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateOrganizationPhoneNumberMutation, UpdateOrganizationPhoneNumberMutationVariables>(UpdateOrganizationPhoneNumberDocument, options);
      }
export type UpdateOrganizationPhoneNumberMutationHookResult = ReturnType<typeof useUpdateOrganizationPhoneNumberMutation>;
export type UpdateOrganizationPhoneNumberMutationResult = Apollo.MutationResult<UpdateOrganizationPhoneNumberMutation>;
export type UpdateOrganizationPhoneNumberMutationOptions = Apollo.BaseMutationOptions<UpdateOrganizationPhoneNumberMutation, UpdateOrganizationPhoneNumberMutationVariables>;
export const UpdateOrganizationWebsiteDocument = gql`
    mutation updateOrganizationWebsite($input: OrganizationUpdateInput!) {
  organization_Update(input: $input) {
    id
    website
  }
}
    `;
export type UpdateOrganizationWebsiteMutationFn = Apollo.MutationFunction<UpdateOrganizationWebsiteMutation, UpdateOrganizationWebsiteMutationVariables>;

/**
 * __useUpdateOrganizationWebsiteMutation__
 *
 * To run a mutation, you first call `useUpdateOrganizationWebsiteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateOrganizationWebsiteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateOrganizationWebsiteMutation, { data, loading, error }] = useUpdateOrganizationWebsiteMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateOrganizationWebsiteMutation(baseOptions?: Apollo.MutationHookOptions<UpdateOrganizationWebsiteMutation, UpdateOrganizationWebsiteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateOrganizationWebsiteMutation, UpdateOrganizationWebsiteMutationVariables>(UpdateOrganizationWebsiteDocument, options);
      }
export type UpdateOrganizationWebsiteMutationHookResult = ReturnType<typeof useUpdateOrganizationWebsiteMutation>;
export type UpdateOrganizationWebsiteMutationResult = Apollo.MutationResult<UpdateOrganizationWebsiteMutation>;
export type UpdateOrganizationWebsiteMutationOptions = Apollo.BaseMutationOptions<UpdateOrganizationWebsiteMutation, UpdateOrganizationWebsiteMutationVariables>;
export const UpdateRenewalForecastDocument = gql`
    mutation updateRenewalForecast($input: RenewalForecastInput!) {
  organization_UpdateRenewalForecastAsync(input: $input)
}
    `;
export type UpdateRenewalForecastMutationFn = Apollo.MutationFunction<UpdateRenewalForecastMutation, UpdateRenewalForecastMutationVariables>;

/**
 * __useUpdateRenewalForecastMutation__
 *
 * To run a mutation, you first call `useUpdateRenewalForecastMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateRenewalForecastMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateRenewalForecastMutation, { data, loading, error }] = useUpdateRenewalForecastMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateRenewalForecastMutation(baseOptions?: Apollo.MutationHookOptions<UpdateRenewalForecastMutation, UpdateRenewalForecastMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateRenewalForecastMutation, UpdateRenewalForecastMutationVariables>(UpdateRenewalForecastDocument, options);
      }
export type UpdateRenewalForecastMutationHookResult = ReturnType<typeof useUpdateRenewalForecastMutation>;
export type UpdateRenewalForecastMutationResult = Apollo.MutationResult<UpdateRenewalForecastMutation>;
export type UpdateRenewalForecastMutationOptions = Apollo.BaseMutationOptions<UpdateRenewalForecastMutation, UpdateRenewalForecastMutationVariables>;
export const UpdateRenewalLikelihoodDocument = gql`
    mutation updateRenewalLikelihood($input: RenewalLikelihoodInput!) {
  organization_UpdateRenewalLikelihoodAsync(input: $input)
}
    `;
export type UpdateRenewalLikelihoodMutationFn = Apollo.MutationFunction<UpdateRenewalLikelihoodMutation, UpdateRenewalLikelihoodMutationVariables>;

/**
 * __useUpdateRenewalLikelihoodMutation__
 *
 * To run a mutation, you first call `useUpdateRenewalLikelihoodMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateRenewalLikelihoodMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateRenewalLikelihoodMutation, { data, loading, error }] = useUpdateRenewalLikelihoodMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateRenewalLikelihoodMutation(baseOptions?: Apollo.MutationHookOptions<UpdateRenewalLikelihoodMutation, UpdateRenewalLikelihoodMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateRenewalLikelihoodMutation, UpdateRenewalLikelihoodMutationVariables>(UpdateRenewalLikelihoodDocument, options);
      }
export type UpdateRenewalLikelihoodMutationHookResult = ReturnType<typeof useUpdateRenewalLikelihoodMutation>;
export type UpdateRenewalLikelihoodMutationResult = Apollo.MutationResult<UpdateRenewalLikelihoodMutation>;
export type UpdateRenewalLikelihoodMutationOptions = Apollo.BaseMutationOptions<UpdateRenewalLikelihoodMutation, UpdateRenewalLikelihoodMutationVariables>;
export const CreateMeetingDocument = gql`
    mutation createMeeting($meeting: MeetingInput!) {
  meeting_Create(meeting: $meeting) {
    id
    attendedBy {
      ... on ContactParticipant {
        contactParticipant {
          id
          name
          firstName
          lastName
        }
      }
      ... on UserParticipant {
        userParticipant {
          id
          lastName
          firstName
        }
      }
    }
    conferenceUrl
    meetingStartedAt: startedAt
    meetingEndedAt: endedAt
    name
    agenda
    agendaContentType
    note {
      id
      appSource
    }
    createdBy {
      ... on ContactParticipant {
        contactParticipant {
          id
        }
      }
      ... on UserParticipant {
        userParticipant {
          id
        }
      }
    }
  }
}
    `;
export type CreateMeetingMutationFn = Apollo.MutationFunction<CreateMeetingMutation, CreateMeetingMutationVariables>;

/**
 * __useCreateMeetingMutation__
 *
 * To run a mutation, you first call `useCreateMeetingMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateMeetingMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createMeetingMutation, { data, loading, error }] = useCreateMeetingMutation({
 *   variables: {
 *      meeting: // value for 'meeting'
 *   },
 * });
 */
export function useCreateMeetingMutation(baseOptions?: Apollo.MutationHookOptions<CreateMeetingMutation, CreateMeetingMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreateMeetingMutation, CreateMeetingMutationVariables>(CreateMeetingDocument, options);
      }
export type CreateMeetingMutationHookResult = ReturnType<typeof useCreateMeetingMutation>;
export type CreateMeetingMutationResult = Apollo.MutationResult<CreateMeetingMutation>;
export type CreateMeetingMutationOptions = Apollo.BaseMutationOptions<CreateMeetingMutation, CreateMeetingMutationVariables>;
export const GetEmailValidationDocument = gql`
    query GetEmailValidation($id: ID!) {
  email(id: $id) {
    id
    emailValidationDetails {
      isReachable
      isValidSyntax
      canConnectSmtp
      acceptsMail
      hasFullInbox
      isCatchAll
      isDeliverable
      validated
      isDisabled
    }
  }
}
    `;

/**
 * __useGetEmailValidationQuery__
 *
 * To run a query within a React component, call `useGetEmailValidationQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetEmailValidationQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetEmailValidationQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetEmailValidationQuery(baseOptions: Apollo.QueryHookOptions<GetEmailValidationQuery, GetEmailValidationQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetEmailValidationQuery, GetEmailValidationQueryVariables>(GetEmailValidationDocument, options);
      }
export function useGetEmailValidationLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetEmailValidationQuery, GetEmailValidationQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetEmailValidationQuery, GetEmailValidationQueryVariables>(GetEmailValidationDocument, options);
        }
export type GetEmailValidationQueryHookResult = ReturnType<typeof useGetEmailValidationQuery>;
export type GetEmailValidationLazyQueryHookResult = ReturnType<typeof useGetEmailValidationLazyQuery>;
export type GetEmailValidationQueryResult = Apollo.QueryResult<GetEmailValidationQuery, GetEmailValidationQueryVariables>;
export const GetTenantNameDocument = gql`
    query GetTenantName {
  tenant
}
    `;

/**
 * __useGetTenantNameQuery__
 *
 * To run a query within a React component, call `useGetTenantNameQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetTenantNameQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetTenantNameQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetTenantNameQuery(baseOptions?: Apollo.QueryHookOptions<GetTenantNameQuery, GetTenantNameQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetTenantNameQuery, GetTenantNameQueryVariables>(GetTenantNameDocument, options);
      }
export function useGetTenantNameLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetTenantNameQuery, GetTenantNameQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetTenantNameQuery, GetTenantNameQueryVariables>(GetTenantNameDocument, options);
        }
export type GetTenantNameQueryHookResult = ReturnType<typeof useGetTenantNameQuery>;
export type GetTenantNameLazyQueryHookResult = ReturnType<typeof useGetTenantNameLazyQuery>;
export type GetTenantNameQueryResult = Apollo.QueryResult<GetTenantNameQuery, GetTenantNameQueryVariables>;
export const LinkMeetingAttachmentDocument = gql`
    mutation linkMeetingAttachment($meetingId: ID!, $attachmentId: ID!) {
  meeting_LinkAttachment(meetingId: $meetingId, attachmentId: $attachmentId) {
    id
  }
}
    `;
export type LinkMeetingAttachmentMutationFn = Apollo.MutationFunction<LinkMeetingAttachmentMutation, LinkMeetingAttachmentMutationVariables>;

/**
 * __useLinkMeetingAttachmentMutation__
 *
 * To run a mutation, you first call `useLinkMeetingAttachmentMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useLinkMeetingAttachmentMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [linkMeetingAttachmentMutation, { data, loading, error }] = useLinkMeetingAttachmentMutation({
 *   variables: {
 *      meetingId: // value for 'meetingId'
 *      attachmentId: // value for 'attachmentId'
 *   },
 * });
 */
export function useLinkMeetingAttachmentMutation(baseOptions?: Apollo.MutationHookOptions<LinkMeetingAttachmentMutation, LinkMeetingAttachmentMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<LinkMeetingAttachmentMutation, LinkMeetingAttachmentMutationVariables>(LinkMeetingAttachmentDocument, options);
      }
export type LinkMeetingAttachmentMutationHookResult = ReturnType<typeof useLinkMeetingAttachmentMutation>;
export type LinkMeetingAttachmentMutationResult = Apollo.MutationResult<LinkMeetingAttachmentMutation>;
export type LinkMeetingAttachmentMutationOptions = Apollo.BaseMutationOptions<LinkMeetingAttachmentMutation, LinkMeetingAttachmentMutationVariables>;
export const MeetingLinkAttachmentDocument = gql`
    mutation meetingLinkAttachment($meetingId: ID!, $attachmentId: ID!) {
  meeting_LinkAttachment(meetingId: $meetingId, attachmentId: $attachmentId) {
    id
    includes {
      id
      name
      mimeType
    }
  }
}
    `;
export type MeetingLinkAttachmentMutationFn = Apollo.MutationFunction<MeetingLinkAttachmentMutation, MeetingLinkAttachmentMutationVariables>;

/**
 * __useMeetingLinkAttachmentMutation__
 *
 * To run a mutation, you first call `useMeetingLinkAttachmentMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useMeetingLinkAttachmentMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [meetingLinkAttachmentMutation, { data, loading, error }] = useMeetingLinkAttachmentMutation({
 *   variables: {
 *      meetingId: // value for 'meetingId'
 *      attachmentId: // value for 'attachmentId'
 *   },
 * });
 */
export function useMeetingLinkAttachmentMutation(baseOptions?: Apollo.MutationHookOptions<MeetingLinkAttachmentMutation, MeetingLinkAttachmentMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<MeetingLinkAttachmentMutation, MeetingLinkAttachmentMutationVariables>(MeetingLinkAttachmentDocument, options);
      }
export type MeetingLinkAttachmentMutationHookResult = ReturnType<typeof useMeetingLinkAttachmentMutation>;
export type MeetingLinkAttachmentMutationResult = Apollo.MutationResult<MeetingLinkAttachmentMutation>;
export type MeetingLinkAttachmentMutationOptions = Apollo.BaseMutationOptions<MeetingLinkAttachmentMutation, MeetingLinkAttachmentMutationVariables>;
export const LinkMeetingAttendeeDocument = gql`
    mutation linkMeetingAttendee($meetingId: ID!, $participant: MeetingParticipantInput!) {
  meeting_LinkAttendedBy(meetingId: $meetingId, participant: $participant) {
    id
    attendedBy {
      ... on ContactParticipant {
        contactParticipant {
          id
          name
          firstName
          lastName
        }
      }
      ... on UserParticipant {
        userParticipant {
          id
          lastName
          firstName
        }
      }
    }
  }
}
    `;
export type LinkMeetingAttendeeMutationFn = Apollo.MutationFunction<LinkMeetingAttendeeMutation, LinkMeetingAttendeeMutationVariables>;

/**
 * __useLinkMeetingAttendeeMutation__
 *
 * To run a mutation, you first call `useLinkMeetingAttendeeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useLinkMeetingAttendeeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [linkMeetingAttendeeMutation, { data, loading, error }] = useLinkMeetingAttendeeMutation({
 *   variables: {
 *      meetingId: // value for 'meetingId'
 *      participant: // value for 'participant'
 *   },
 * });
 */
export function useLinkMeetingAttendeeMutation(baseOptions?: Apollo.MutationHookOptions<LinkMeetingAttendeeMutation, LinkMeetingAttendeeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<LinkMeetingAttendeeMutation, LinkMeetingAttendeeMutationVariables>(LinkMeetingAttendeeDocument, options);
      }
export type LinkMeetingAttendeeMutationHookResult = ReturnType<typeof useLinkMeetingAttendeeMutation>;
export type LinkMeetingAttendeeMutationResult = Apollo.MutationResult<LinkMeetingAttendeeMutation>;
export type LinkMeetingAttendeeMutationOptions = Apollo.BaseMutationOptions<LinkMeetingAttendeeMutation, LinkMeetingAttendeeMutationVariables>;
export const MeetingLinkRecordingDocument = gql`
    mutation meetingLinkRecording($meetingId: ID!, $attachmentId: ID!) {
  meeting_LinkRecording(meetingId: $meetingId, attachmentId: $attachmentId) {
    id
    attendedBy {
      ... on UserParticipant {
        userParticipant {
          id
          firstName
          lastName
        }
      }
      ... on ContactParticipant {
        contactParticipant {
          id
          firstName
          lastName
          name
        }
      }
    }
    recording {
      id
    }
    meetingStartedAt: startedAt
    agenda
  }
}
    `;
export type MeetingLinkRecordingMutationFn = Apollo.MutationFunction<MeetingLinkRecordingMutation, MeetingLinkRecordingMutationVariables>;

/**
 * __useMeetingLinkRecordingMutation__
 *
 * To run a mutation, you first call `useMeetingLinkRecordingMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useMeetingLinkRecordingMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [meetingLinkRecordingMutation, { data, loading, error }] = useMeetingLinkRecordingMutation({
 *   variables: {
 *      meetingId: // value for 'meetingId'
 *      attachmentId: // value for 'attachmentId'
 *   },
 * });
 */
export function useMeetingLinkRecordingMutation(baseOptions?: Apollo.MutationHookOptions<MeetingLinkRecordingMutation, MeetingLinkRecordingMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<MeetingLinkRecordingMutation, MeetingLinkRecordingMutationVariables>(MeetingLinkRecordingDocument, options);
      }
export type MeetingLinkRecordingMutationHookResult = ReturnType<typeof useMeetingLinkRecordingMutation>;
export type MeetingLinkRecordingMutationResult = Apollo.MutationResult<MeetingLinkRecordingMutation>;
export type MeetingLinkRecordingMutationOptions = Apollo.BaseMutationOptions<MeetingLinkRecordingMutation, MeetingLinkRecordingMutationVariables>;
export const MeetingUnlinkAttachmentDocument = gql`
    mutation meetingUnlinkAttachment($meetingId: ID!, $attachmentId: ID!) {
  meeting_UnlinkAttachment(meetingId: $meetingId, attachmentId: $attachmentId) {
    id
    includes {
      id
      name
      mimeType
    }
  }
}
    `;
export type MeetingUnlinkAttachmentMutationFn = Apollo.MutationFunction<MeetingUnlinkAttachmentMutation, MeetingUnlinkAttachmentMutationVariables>;

/**
 * __useMeetingUnlinkAttachmentMutation__
 *
 * To run a mutation, you first call `useMeetingUnlinkAttachmentMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useMeetingUnlinkAttachmentMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [meetingUnlinkAttachmentMutation, { data, loading, error }] = useMeetingUnlinkAttachmentMutation({
 *   variables: {
 *      meetingId: // value for 'meetingId'
 *      attachmentId: // value for 'attachmentId'
 *   },
 * });
 */
export function useMeetingUnlinkAttachmentMutation(baseOptions?: Apollo.MutationHookOptions<MeetingUnlinkAttachmentMutation, MeetingUnlinkAttachmentMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<MeetingUnlinkAttachmentMutation, MeetingUnlinkAttachmentMutationVariables>(MeetingUnlinkAttachmentDocument, options);
      }
export type MeetingUnlinkAttachmentMutationHookResult = ReturnType<typeof useMeetingUnlinkAttachmentMutation>;
export type MeetingUnlinkAttachmentMutationResult = Apollo.MutationResult<MeetingUnlinkAttachmentMutation>;
export type MeetingUnlinkAttachmentMutationOptions = Apollo.BaseMutationOptions<MeetingUnlinkAttachmentMutation, MeetingUnlinkAttachmentMutationVariables>;
export const UnlinkMeetingAttendeeDocument = gql`
    mutation unlinkMeetingAttendee($meetingId: ID!, $participant: MeetingParticipantInput!) {
  meeting_UnlinkAttendedBy(meetingId: $meetingId, participant: $participant) {
    id
    attendedBy {
      ... on ContactParticipant {
        contactParticipant {
          id
          name
          firstName
          lastName
        }
      }
      ... on UserParticipant {
        userParticipant {
          id
          lastName
          firstName
        }
      }
    }
  }
}
    `;
export type UnlinkMeetingAttendeeMutationFn = Apollo.MutationFunction<UnlinkMeetingAttendeeMutation, UnlinkMeetingAttendeeMutationVariables>;

/**
 * __useUnlinkMeetingAttendeeMutation__
 *
 * To run a mutation, you first call `useUnlinkMeetingAttendeeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUnlinkMeetingAttendeeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [unlinkMeetingAttendeeMutation, { data, loading, error }] = useUnlinkMeetingAttendeeMutation({
 *   variables: {
 *      meetingId: // value for 'meetingId'
 *      participant: // value for 'participant'
 *   },
 * });
 */
export function useUnlinkMeetingAttendeeMutation(baseOptions?: Apollo.MutationHookOptions<UnlinkMeetingAttendeeMutation, UnlinkMeetingAttendeeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UnlinkMeetingAttendeeMutation, UnlinkMeetingAttendeeMutationVariables>(UnlinkMeetingAttendeeDocument, options);
      }
export type UnlinkMeetingAttendeeMutationHookResult = ReturnType<typeof useUnlinkMeetingAttendeeMutation>;
export type UnlinkMeetingAttendeeMutationResult = Apollo.MutationResult<UnlinkMeetingAttendeeMutation>;
export type UnlinkMeetingAttendeeMutationOptions = Apollo.BaseMutationOptions<UnlinkMeetingAttendeeMutation, UnlinkMeetingAttendeeMutationVariables>;
export const MeetingUnlinkRecordingDocument = gql`
    mutation meetingUnlinkRecording($meetingId: ID!, $attachmentId: ID!) {
  meeting_UnlinkRecording(meetingId: $meetingId, attachmentId: $attachmentId) {
    id
    includes {
      id
    }
  }
}
    `;
export type MeetingUnlinkRecordingMutationFn = Apollo.MutationFunction<MeetingUnlinkRecordingMutation, MeetingUnlinkRecordingMutationVariables>;

/**
 * __useMeetingUnlinkRecordingMutation__
 *
 * To run a mutation, you first call `useMeetingUnlinkRecordingMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useMeetingUnlinkRecordingMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [meetingUnlinkRecordingMutation, { data, loading, error }] = useMeetingUnlinkRecordingMutation({
 *   variables: {
 *      meetingId: // value for 'meetingId'
 *      attachmentId: // value for 'attachmentId'
 *   },
 * });
 */
export function useMeetingUnlinkRecordingMutation(baseOptions?: Apollo.MutationHookOptions<MeetingUnlinkRecordingMutation, MeetingUnlinkRecordingMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<MeetingUnlinkRecordingMutation, MeetingUnlinkRecordingMutationVariables>(MeetingUnlinkRecordingDocument, options);
      }
export type MeetingUnlinkRecordingMutationHookResult = ReturnType<typeof useMeetingUnlinkRecordingMutation>;
export type MeetingUnlinkRecordingMutationResult = Apollo.MutationResult<MeetingUnlinkRecordingMutation>;
export type MeetingUnlinkRecordingMutationOptions = Apollo.BaseMutationOptions<MeetingUnlinkRecordingMutation, MeetingUnlinkRecordingMutationVariables>;
export const NoteLinkAttachmentDocument = gql`
    mutation noteLinkAttachment($noteId: ID!, $attachmentId: ID!) {
  note_LinkAttachment(noteId: $noteId, attachmentId: $attachmentId) {
    id
    includes {
      id
      name
      mimeType
    }
  }
}
    `;
export type NoteLinkAttachmentMutationFn = Apollo.MutationFunction<NoteLinkAttachmentMutation, NoteLinkAttachmentMutationVariables>;

/**
 * __useNoteLinkAttachmentMutation__
 *
 * To run a mutation, you first call `useNoteLinkAttachmentMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useNoteLinkAttachmentMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [noteLinkAttachmentMutation, { data, loading, error }] = useNoteLinkAttachmentMutation({
 *   variables: {
 *      noteId: // value for 'noteId'
 *      attachmentId: // value for 'attachmentId'
 *   },
 * });
 */
export function useNoteLinkAttachmentMutation(baseOptions?: Apollo.MutationHookOptions<NoteLinkAttachmentMutation, NoteLinkAttachmentMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<NoteLinkAttachmentMutation, NoteLinkAttachmentMutationVariables>(NoteLinkAttachmentDocument, options);
      }
export type NoteLinkAttachmentMutationHookResult = ReturnType<typeof useNoteLinkAttachmentMutation>;
export type NoteLinkAttachmentMutationResult = Apollo.MutationResult<NoteLinkAttachmentMutation>;
export type NoteLinkAttachmentMutationOptions = Apollo.BaseMutationOptions<NoteLinkAttachmentMutation, NoteLinkAttachmentMutationVariables>;
export const NoteUnlinkAttachmentDocument = gql`
    mutation noteUnlinkAttachment($noteId: ID!, $attachmentId: ID!) {
  note_UnlinkAttachment(noteId: $noteId, attachmentId: $attachmentId) {
    id
    includes {
      id
      name
      mimeType
    }
  }
}
    `;
export type NoteUnlinkAttachmentMutationFn = Apollo.MutationFunction<NoteUnlinkAttachmentMutation, NoteUnlinkAttachmentMutationVariables>;

/**
 * __useNoteUnlinkAttachmentMutation__
 *
 * To run a mutation, you first call `useNoteUnlinkAttachmentMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useNoteUnlinkAttachmentMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [noteUnlinkAttachmentMutation, { data, loading, error }] = useNoteUnlinkAttachmentMutation({
 *   variables: {
 *      noteId: // value for 'noteId'
 *      attachmentId: // value for 'attachmentId'
 *   },
 * });
 */
export function useNoteUnlinkAttachmentMutation(baseOptions?: Apollo.MutationHookOptions<NoteUnlinkAttachmentMutation, NoteUnlinkAttachmentMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<NoteUnlinkAttachmentMutation, NoteUnlinkAttachmentMutationVariables>(NoteUnlinkAttachmentDocument, options);
      }
export type NoteUnlinkAttachmentMutationHookResult = ReturnType<typeof useNoteUnlinkAttachmentMutation>;
export type NoteUnlinkAttachmentMutationResult = Apollo.MutationResult<NoteUnlinkAttachmentMutation>;
export type NoteUnlinkAttachmentMutationOptions = Apollo.BaseMutationOptions<NoteUnlinkAttachmentMutation, NoteUnlinkAttachmentMutationVariables>;
export const RemoveNoteDocument = gql`
    mutation removeNote($id: ID!) {
  note_Delete(id: $id) {
    result
  }
}
    `;
export type RemoveNoteMutationFn = Apollo.MutationFunction<RemoveNoteMutation, RemoveNoteMutationVariables>;

/**
 * __useRemoveNoteMutation__
 *
 * To run a mutation, you first call `useRemoveNoteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemoveNoteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removeNoteMutation, { data, loading, error }] = useRemoveNoteMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useRemoveNoteMutation(baseOptions?: Apollo.MutationHookOptions<RemoveNoteMutation, RemoveNoteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RemoveNoteMutation, RemoveNoteMutationVariables>(RemoveNoteDocument, options);
      }
export type RemoveNoteMutationHookResult = ReturnType<typeof useRemoveNoteMutation>;
export type RemoveNoteMutationResult = Apollo.MutationResult<RemoveNoteMutation>;
export type RemoveNoteMutationOptions = Apollo.BaseMutationOptions<RemoveNoteMutation, RemoveNoteMutationVariables>;
export const UpdateLocationDocument = gql`
    mutation updateLocation($input: LocationUpdateInput!) {
  location_Update(input: $input) {
    locality
    rawAddress
    postalCode
    street
  }
}
    `;
export type UpdateLocationMutationFn = Apollo.MutationFunction<UpdateLocationMutation, UpdateLocationMutationVariables>;

/**
 * __useUpdateLocationMutation__
 *
 * To run a mutation, you first call `useUpdateLocationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateLocationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateLocationMutation, { data, loading, error }] = useUpdateLocationMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateLocationMutation(baseOptions?: Apollo.MutationHookOptions<UpdateLocationMutation, UpdateLocationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateLocationMutation, UpdateLocationMutationVariables>(UpdateLocationDocument, options);
      }
export type UpdateLocationMutationHookResult = ReturnType<typeof useUpdateLocationMutation>;
export type UpdateLocationMutationResult = Apollo.MutationResult<UpdateLocationMutation>;
export type UpdateLocationMutationOptions = Apollo.BaseMutationOptions<UpdateLocationMutation, UpdateLocationMutationVariables>;
export const UpdateMeetingDocument = gql`
    mutation updateMeeting($meetingId: ID!, $meetingInput: MeetingUpdateInput!) {
  meeting_Update(meetingId: $meetingId, meeting: $meetingInput) {
    ...MeetingTimelineEventFragment
  }
}
    ${MeetingTimelineEventFragmentFragmentDoc}`;
export type UpdateMeetingMutationFn = Apollo.MutationFunction<UpdateMeetingMutation, UpdateMeetingMutationVariables>;

/**
 * __useUpdateMeetingMutation__
 *
 * To run a mutation, you first call `useUpdateMeetingMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateMeetingMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateMeetingMutation, { data, loading, error }] = useUpdateMeetingMutation({
 *   variables: {
 *      meetingId: // value for 'meetingId'
 *      meetingInput: // value for 'meetingInput'
 *   },
 * });
 */
export function useUpdateMeetingMutation(baseOptions?: Apollo.MutationHookOptions<UpdateMeetingMutation, UpdateMeetingMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateMeetingMutation, UpdateMeetingMutationVariables>(UpdateMeetingDocument, options);
      }
export type UpdateMeetingMutationHookResult = ReturnType<typeof useUpdateMeetingMutation>;
export type UpdateMeetingMutationResult = Apollo.MutationResult<UpdateMeetingMutation>;
export type UpdateMeetingMutationOptions = Apollo.BaseMutationOptions<UpdateMeetingMutation, UpdateMeetingMutationVariables>;
export const UpdateNoteDocument = gql`
    mutation updateNote($input: NoteUpdateInput!) {
  note_Update(input: $input) {
    ...NoteContent
  }
}
    ${NoteContentFragmentDoc}`;
export type UpdateNoteMutationFn = Apollo.MutationFunction<UpdateNoteMutation, UpdateNoteMutationVariables>;

/**
 * __useUpdateNoteMutation__
 *
 * To run a mutation, you first call `useUpdateNoteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateNoteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateNoteMutation, { data, loading, error }] = useUpdateNoteMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateNoteMutation(baseOptions?: Apollo.MutationHookOptions<UpdateNoteMutation, UpdateNoteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateNoteMutation, UpdateNoteMutationVariables>(UpdateNoteDocument, options);
      }
export type UpdateNoteMutationHookResult = ReturnType<typeof useUpdateNoteMutation>;
export type UpdateNoteMutationResult = Apollo.MutationResult<UpdateNoteMutation>;
export type UpdateNoteMutationOptions = Apollo.BaseMutationOptions<UpdateNoteMutation, UpdateNoteMutationVariables>;
export const GetUserByEmailDocument = gql`
    query getUserByEmail($email: String!) {
  user_ByEmail(email: $email) {
    id
    firstName
    lastName
  }
}
    `;

/**
 * __useGetUserByEmailQuery__
 *
 * To run a query within a React component, call `useGetUserByEmailQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetUserByEmailQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetUserByEmailQuery({
 *   variables: {
 *      email: // value for 'email'
 *   },
 * });
 */
export function useGetUserByEmailQuery(baseOptions: Apollo.QueryHookOptions<GetUserByEmailQuery, GetUserByEmailQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetUserByEmailQuery, GetUserByEmailQueryVariables>(GetUserByEmailDocument, options);
      }
export function useGetUserByEmailLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetUserByEmailQuery, GetUserByEmailQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetUserByEmailQuery, GetUserByEmailQueryVariables>(GetUserByEmailDocument, options);
        }
export type GetUserByEmailQueryHookResult = ReturnType<typeof useGetUserByEmailQuery>;
export type GetUserByEmailLazyQueryHookResult = ReturnType<typeof useGetUserByEmailLazyQuery>;
export type GetUserByEmailQueryResult = Apollo.QueryResult<GetUserByEmailQuery, GetUserByEmailQueryVariables>;
export const GetUsersDocument = gql`
    query getUsers($pagination: Pagination!, $where: Filter) {
  users(pagination: $pagination, where: $where) {
    content {
      id
      firstName
      lastName
    }
    totalElements
  }
}
    `;

/**
 * __useGetUsersQuery__
 *
 * To run a query within a React component, call `useGetUsersQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetUsersQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetUsersQuery({
 *   variables: {
 *      pagination: // value for 'pagination'
 *      where: // value for 'where'
 *   },
 * });
 */
export function useGetUsersQuery(baseOptions: Apollo.QueryHookOptions<GetUsersQuery, GetUsersQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetUsersQuery, GetUsersQueryVariables>(GetUsersDocument, options);
      }
export function useGetUsersLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetUsersQuery, GetUsersQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetUsersQuery, GetUsersQueryVariables>(GetUsersDocument, options);
        }
export type GetUsersQueryHookResult = ReturnType<typeof useGetUsersQuery>;
export type GetUsersLazyQueryHookResult = ReturnType<typeof useGetUsersLazyQuery>;
export type GetUsersQueryResult = Apollo.QueryResult<GetUsersQuery, GetUsersQueryVariables>;