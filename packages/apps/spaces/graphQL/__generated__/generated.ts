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

export type Action = Conversation | InteractionSession | Note | PageViewAction | Ticket;

export enum ActionType {
  Conversation = 'CONVERSATION',
  InteractionSession = 'INTERACTION_SESSION',
  Note = 'NOTE',
  PageView = 'PAGE_VIEW',
  Ticket = 'TICKET'
}

export enum ComparisonOperator {
  Contains = 'CONTAINS',
  Eq = 'EQ'
}

/**
 * A contact represents an individual in customerOS.
 * **A `response` object.**
 */
export type Contact = ExtensibleEntity & Node & {
  __typename?: 'Contact';
  /** @deprecated Use `timelineEvents` instead */
  actions: Array<Action>;
  appSource?: Maybe<Scalars['String']>;
  conversations: ConversationPage;
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
  /**
   * All email addresses associated with a contact in customerOS.
   * **Required.  If no values it returns an empty array.**
   */
  emails: Array<Email>;
  fieldSets: Array<FieldSet>;
  /** The first name of the contact in customerOS. */
  firstName?: Maybe<Scalars['String']>;
  /**
   * Identifies any contact groups the contact is associated with.
   * **Required.  If no values it returns an empty array.**
   */
  groups: Array<ContactGroup>;
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
  source: DataSource;
  sourceOfTruth: DataSource;
  tags?: Maybe<Array<Tag>>;
  /** Template of the contact in customerOS. */
  template?: Maybe<EntityTemplate>;
  ticketSummaryByStatus: Array<TicketSummaryByStatus>;
  tickets: Array<Maybe<Ticket>>;
  timelineEvents: Array<TimelineEvent>;
  timelineEventsTotalCount: Scalars['Int64'];
  /** The title associate with the contact in customerOS. */
  title?: Maybe<PersonTitle>;
  updatedAt: Scalars['Time'];
};


/**
 * A contact represents an individual in customerOS.
 * **A `response` object.**
 */
export type ContactActionsArgs = {
  actionTypes?: InputMaybe<Array<ActionType>>;
  from: Scalars['Time'];
  to: Scalars['Time'];
};


/**
 * A contact represents an individual in customerOS.
 * **A `response` object.**
 */
export type ContactConversationsArgs = {
  pagination?: InputMaybe<Pagination>;
  sort?: InputMaybe<Array<SortBy>>;
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
 * A collection of groups that a Contact belongs to.  Groups are user-defined entities.
 * **A `return` object.**
 */
export type ContactGroup = {
  __typename?: 'ContactGroup';
  contacts: ContactsPage;
  createdAt: Scalars['Time'];
  /**
   * The unique ID associated with the `ContactGroup` in customerOS.
   * **Required**
   */
  id: Scalars['ID'];
  /**
   * The name of the `ContactGroup`.
   * **Required**
   */
  name: Scalars['String'];
  source: DataSource;
};


/**
 * A collection of groups that a Contact belongs to.  Groups are user-defined entities.
 * **A `return` object.**
 */
export type ContactGroupContactsArgs = {
  pagination?: InputMaybe<Pagination>;
  sort?: InputMaybe<Array<SortBy>>;
  where?: InputMaybe<Filter>;
};

/**
 * Create a groups that can be associated with a `Contact` in customerOS.
 * **A `create` object.**
 */
export type ContactGroupInput = {
  /**
   * The name of the `ContactGroup`.
   * **Required**
   */
  name: Scalars['String'];
};

/**
 * Specifies how many pages of `ContactGroup` information has been returned in the query response.
 * **A `response` object.**
 */
export type ContactGroupPage = Pages & {
  __typename?: 'ContactGroupPage';
  /**
   * A collection of groups that a Contact belongs to.  Groups are user-defined entities.
   * **Required.  If no values it returns an empty array.**
   */
  content: Array<ContactGroup>;
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

/**
 * Update a group that can be associated with a `Contact` in customerOS.
 * **A `update` object.**
 */
export type ContactGroupUpdateInput = {
  /**
   * The unique ID associated with the `ContactGroup` in customerOS.
   * **Required**
   */
  id: Scalars['ID'];
  /**
   * The name of the `ContactGroup`.
   * **Required**
   */
  name: Scalars['String'];
};

/**
 * Create an individual in customerOS.
 * **A `create` object.**
 */
export type ContactInput = {
  /** An ISO8601 timestamp recording when the contact was created in customerOS. */
  createdAt?: InputMaybe<Scalars['Time']>;
  /**
   * User defined metadata appended to contact.
   * **Required.**
   */
  customFields?: InputMaybe<Array<CustomFieldInput>>;
  /** An email addresses associted with the contact. */
  email?: InputMaybe<EmailInput>;
  externalReference?: InputMaybe<ExternalSystemReferenceInput>;
  fieldSets?: InputMaybe<Array<FieldSetInput>>;
  /** The first name of the contact. */
  firstName?: InputMaybe<Scalars['String']>;
  label?: InputMaybe<Scalars['String']>;
  /** The last name of the contact. */
  lastName?: InputMaybe<Scalars['String']>;
  /** Id of the contact owner (user) */
  ownerId?: InputMaybe<Scalars['ID']>;
  /** A phone number associated with the contact. */
  phoneNumber?: InputMaybe<PhoneNumberInput>;
  /** The unique ID associated with the template of the contact in customerOS. */
  templateId?: InputMaybe<Scalars['ID']>;
  /** The title of the contact. */
  title?: InputMaybe<PersonTitle>;
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
  /** Id of the contact owner (user) */
  ownerId?: InputMaybe<Scalars['ID']>;
  /** The title associate with the contact in customerOS. */
  title?: InputMaybe<PersonTitle>;
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

export type Conversation = Node & {
  __typename?: 'Conversation';
  appSource?: Maybe<Scalars['String']>;
  channel?: Maybe<Scalars['String']>;
  contacts?: Maybe<Array<Contact>>;
  endedAt?: Maybe<Scalars['Time']>;
  id: Scalars['ID'];
  initiatorFirstName?: Maybe<Scalars['String']>;
  initiatorLastName?: Maybe<Scalars['String']>;
  initiatorType?: Maybe<Scalars['String']>;
  initiatorUsername?: Maybe<Scalars['String']>;
  messageCount: Scalars['Int64'];
  source: DataSource;
  sourceOfTruth: DataSource;
  startedAt: Scalars['Time'];
  status: ConversationStatus;
  subject?: Maybe<Scalars['String']>;
  threadId?: Maybe<Scalars['String']>;
  updatedAt: Scalars['Time'];
  users?: Maybe<Array<User>>;
};

export type ConversationInput = {
  appSource?: InputMaybe<Scalars['String']>;
  channel?: InputMaybe<Scalars['String']>;
  contactIds?: InputMaybe<Array<Scalars['ID']>>;
  id?: InputMaybe<Scalars['ID']>;
  startedAt?: InputMaybe<Scalars['Time']>;
  status?: ConversationStatus;
  userIds?: InputMaybe<Array<Scalars['ID']>>;
};

export type ConversationPage = Pages & {
  __typename?: 'ConversationPage';
  content: Array<Conversation>;
  totalElements: Scalars['Int64'];
  totalPages: Scalars['Int'];
};

export enum ConversationStatus {
  Active = 'ACTIVE',
  Closed = 'CLOSED'
}

export type ConversationUpdateInput = {
  channel?: InputMaybe<Scalars['String']>;
  contactIds?: InputMaybe<Array<Scalars['ID']>>;
  id: Scalars['ID'];
  skipMessageCountIncrement?: Scalars['Boolean'];
  status?: InputMaybe<ConversationStatus>;
  userIds?: InputMaybe<Array<Scalars['ID']>>;
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

export type DashboardViewItem = {
  __typename?: 'DashboardViewItem';
  contact?: Maybe<Contact>;
  organization?: Maybe<Organization>;
};

export type DashboardViewItemPage = Pages & {
  __typename?: 'DashboardViewItemPage';
  content: Array<DashboardViewItem>;
  totalElements: Scalars['Int64'];
  totalPages: Scalars['Int'];
};

export enum DataSource {
  Hubspot = 'HUBSPOT',
  Na = 'NA',
  Openline = 'OPENLINE',
  ZendeskSupport = 'ZENDESK_SUPPORT'
}

/**
 * Describes an email address associated with a `Contact` in customerOS.
 * **A `return` object.**
 */
export type Email = {
  __typename?: 'Email';
  appSource: Scalars['String'];
  createdAt: Scalars['Time'];
  /** An email address assocaited with the contact in customerOS. */
  email?: Maybe<Scalars['String']>;
  /**
   * The unique ID associated with the contact in customerOS.
   * **Required**
   */
  id: Scalars['ID'];
  /** Describes the type of email address (WORK, PERSONAL, etc). */
  label?: Maybe<EmailLabel>;
  /**
   * Identifies whether the email address is primary or not.
   * **Required.**
   */
  primary: Scalars['Boolean'];
  rawEmail?: Maybe<Scalars['String']>;
  source: DataSource;
  sourceOfTruth: DataSource;
  updatedAt: Scalars['Time'];
  validated?: Maybe<Scalars['Boolean']>;
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
  Home = 'HOME',
  Main = 'MAIN',
  Other = 'OTHER',
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

export type EntityTemplate = Node & {
  __typename?: 'EntityTemplate';
  createdAt: Scalars['Time'];
  customFields: Array<CustomFieldTemplate>;
  extends?: Maybe<EntityTemplateExtension>;
  fieldSets: Array<FieldSetTemplate>;
  id: Scalars['ID'];
  name: Scalars['String'];
  updatedAt: Scalars['Time'];
  version: Scalars['Int'];
};

export enum EntityTemplateExtension {
  Contact = 'CONTACT'
}

export type EntityTemplateInput = {
  customFields?: InputMaybe<Array<CustomFieldTemplateInput>>;
  extends?: InputMaybe<EntityTemplateExtension>;
  fieldSets?: InputMaybe<Array<FieldSetTemplateInput>>;
  name: Scalars['String'];
};

export type ExtensibleEntity = {
  id: Scalars['ID'];
  template?: Maybe<EntityTemplate>;
};

export type ExternalSystemReferenceInput = {
  id: Scalars['ID'];
  syncDate?: InputMaybe<Scalars['Time']>;
  type: ExternalSystemType;
};

export enum ExternalSystemType {
  Hubspot = 'HUBSPOT',
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
  customFields: Array<CustomFieldTemplate>;
  id: Scalars['ID'];
  name: Scalars['String'];
  order: Scalars['Int'];
  updatedAt: Scalars['Time'];
};

export type FieldSetTemplateInput = {
  customFields?: InputMaybe<Array<CustomFieldTemplateInput>>;
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

export type InteractionEvent = Node & {
  __typename?: 'InteractionEvent';
  appSource: Scalars['String'];
  channel?: Maybe<Scalars['String']>;
  content?: Maybe<Scalars['String']>;
  contentType?: Maybe<Scalars['String']>;
  createdAt: Scalars['Time'];
  eventIdentifier?: Maybe<Scalars['String']>;
  id: Scalars['ID'];
  interactionSession?: Maybe<InteractionSession>;
  sentBy: Array<InteractionEventParticipant>;
  sentTo: Array<InteractionEventParticipant>;
  source: DataSource;
  sourceOfTruth: DataSource;
};

export type InteractionEventParticipant = ContactParticipant | EmailParticipant | PhoneNumberParticipant | UserParticipant;

export type InteractionSession = Node & {
  __typename?: 'InteractionSession';
  appSource: Scalars['String'];
  channel?: Maybe<Scalars['String']>;
  endedAt?: Maybe<Scalars['Time']>;
  events: Array<InteractionEvent>;
  id: Scalars['ID'];
  name?: Maybe<Scalars['String']>;
  source: DataSource;
  sourceOfTruth: DataSource;
  startedAt: Scalars['Time'];
  status?: Maybe<Scalars['String']>;
  type?: Maybe<Scalars['String']>;
};

/**
 * Describes the relationship a Contact has with a Organization.
 * **A `return` object**
 */
export type JobRole = {
  __typename?: 'JobRole';
  appSource: Scalars['String'];
  contact?: Maybe<Contact>;
  createdAt: Scalars['Time'];
  id: Scalars['ID'];
  /** The Contact's job title. */
  jobTitle?: Maybe<Scalars['String']>;
  /**
   * Organization associated with a Contact.
   * **Required.**
   */
  organization?: Maybe<Organization>;
  primary: Scalars['Boolean'];
  responsibilityLevel: Scalars['Int64'];
  source: DataSource;
  sourceOfTruth: DataSource;
  updatedAt: Scalars['Time'];
};

/**
 * Describes the relationship a Contact has with an Organization.
 * **A `create` object**
 */
export type JobRoleInput = {
  appSource?: InputMaybe<Scalars['String']>;
  /** The Contact's job title. */
  jobTitle?: InputMaybe<Scalars['String']>;
  organizationId?: InputMaybe<Scalars['ID']>;
  primary?: InputMaybe<Scalars['Boolean']>;
  responsibilityLevel?: InputMaybe<Scalars['Int64']>;
};

/**
 * Describes the relationship a Contact has with an Organization.
 * **A `create` object**
 */
export type JobRoleUpdateInput = {
  id: Scalars['ID'];
  /** The Contact's job title. */
  jobTitle?: InputMaybe<Scalars['String']>;
  organizationId?: InputMaybe<Scalars['ID']>;
  primary?: InputMaybe<Scalars['Boolean']>;
  responsibilityLevel?: InputMaybe<Scalars['Int64']>;
};

export type Location = {
  __typename?: 'Location';
  address?: Maybe<Scalars['String']>;
  address2?: Maybe<Scalars['String']>;
  addressType?: Maybe<Scalars['String']>;
  appSource?: Maybe<Scalars['String']>;
  commercial?: Maybe<Scalars['Boolean']>;
  country?: Maybe<Scalars['String']>;
  createdAt: Scalars['Time'];
  district?: Maybe<Scalars['String']>;
  houseNumber?: Maybe<Scalars['String']>;
  id: Scalars['ID'];
  latitude?: Maybe<Scalars['Float']>;
  locality?: Maybe<Scalars['String']>;
  longitude?: Maybe<Scalars['Float']>;
  name: Scalars['String'];
  /** @deprecated Use location instead */
  place?: Maybe<Place>;
  plusFour?: Maybe<Scalars['String']>;
  postalCode?: Maybe<Scalars['String']>;
  predirection?: Maybe<Scalars['String']>;
  rawAddress?: Maybe<Scalars['String']>;
  region?: Maybe<Scalars['String']>;
  source?: Maybe<DataSource>;
  street?: Maybe<Scalars['String']>;
  updatedAt: Scalars['Time'];
  zip?: Maybe<Scalars['String']>;
};

export type Mutation = {
  __typename?: 'Mutation';
  contactGroupAddContact: Result;
  contactGroupCreate: ContactGroup;
  contactGroupDeleteAndUnlinkAllContacts: Result;
  contactGroupRemoveContact: Result;
  contactGroupUpdate: ContactGroup;
  contact_AddOrganizationById: Contact;
  contact_AddTagById: Contact;
  contact_Create: Contact;
  contact_HardDelete: Result;
  contact_Merge: Contact;
  contact_RemoveOrganizationById: Contact;
  contact_RemoveTagById: Contact;
  contact_SoftDelete: Result;
  contact_Update: Contact;
  conversation_Close: Conversation;
  conversation_Create: Conversation;
  conversation_Update: Conversation;
  customFieldDeleteFromContactById: Result;
  customFieldDeleteFromContactByName: Result;
  customFieldDeleteFromFieldSetById: Result;
  customFieldMergeToContact: CustomField;
  customFieldMergeToFieldSet: CustomField;
  customFieldUpdateInContact: CustomField;
  customFieldUpdateInFieldSet: CustomField;
  customFieldsMergeAndUpdateInContact: Contact;
  emailDelete: Result;
  emailMergeToContact: Email;
  emailMergeToUser: Email;
  emailRemoveFromContact: Result;
  emailRemoveFromContactById: Result;
  emailRemoveFromUser: Result;
  emailRemoveFromUserById: Result;
  emailUpdateInContact: Email;
  emailUpdateInUser: Email;
  entityTemplateCreate: EntityTemplate;
  fieldSetDeleteFromContact: Result;
  fieldSetMergeToContact?: Maybe<FieldSet>;
  fieldSetUpdateInContact?: Maybe<FieldSet>;
  jobRole_Create: JobRole;
  jobRole_Delete: Result;
  jobRole_Update: JobRole;
  note_CreateForContact: Note;
  note_CreateForOrganization: Note;
  note_Delete: Result;
  note_Update: Note;
  organizationType_Create: OrganizationType;
  organizationType_Delete?: Maybe<Result>;
  organizationType_Update?: Maybe<OrganizationType>;
  organization_Create: Organization;
  organization_Delete?: Maybe<Result>;
  organization_Merge: Organization;
  organization_Update: Organization;
  phoneNumberMergeToContact: PhoneNumber;
  phoneNumberRemoveFromContactByE164: Result;
  phoneNumberRemoveFromContactById: Result;
  phoneNumberUpdateInContact: PhoneNumber;
  tag_Create: Tag;
  tag_Delete?: Maybe<Result>;
  tag_Update?: Maybe<Tag>;
  user_Create: User;
  user_Update: User;
};


export type MutationContactGroupAddContactArgs = {
  contactId: Scalars['ID'];
  groupId: Scalars['ID'];
};


export type MutationContactGroupCreateArgs = {
  input: ContactGroupInput;
};


export type MutationContactGroupDeleteAndUnlinkAllContactsArgs = {
  id: Scalars['ID'];
};


export type MutationContactGroupRemoveContactArgs = {
  contactId: Scalars['ID'];
  groupId: Scalars['ID'];
};


export type MutationContactGroupUpdateArgs = {
  input: ContactGroupUpdateInput;
};


export type MutationContact_AddOrganizationByIdArgs = {
  input: ContactOrganizationInput;
};


export type MutationContact_AddTagByIdArgs = {
  input: ContactTagInput;
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


export type MutationContact_RemoveOrganizationByIdArgs = {
  input: ContactOrganizationInput;
};


export type MutationContact_RemoveTagByIdArgs = {
  input: ContactTagInput;
};


export type MutationContact_SoftDeleteArgs = {
  contactId: Scalars['ID'];
};


export type MutationContact_UpdateArgs = {
  input: ContactUpdateInput;
};


export type MutationConversation_CloseArgs = {
  conversationId: Scalars['ID'];
};


export type MutationConversation_CreateArgs = {
  input: ConversationInput;
};


export type MutationConversation_UpdateArgs = {
  input: ConversationUpdateInput;
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


export type MutationEmailDeleteArgs = {
  id: Scalars['ID'];
};


export type MutationEmailMergeToContactArgs = {
  contactId: Scalars['ID'];
  input: EmailInput;
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


export type MutationNote_UpdateArgs = {
  input: NoteUpdateInput;
};


export type MutationOrganizationType_CreateArgs = {
  input: OrganizationTypeInput;
};


export type MutationOrganizationType_DeleteArgs = {
  id: Scalars['ID'];
};


export type MutationOrganizationType_UpdateArgs = {
  input: OrganizationTypeUpdateInput;
};


export type MutationOrganization_CreateArgs = {
  input: OrganizationInput;
};


export type MutationOrganization_DeleteArgs = {
  id: Scalars['ID'];
};


export type MutationOrganization_MergeArgs = {
  mergedOrganizationIds: Array<Scalars['ID']>;
  primaryOrganizationId: Scalars['ID'];
};


export type MutationOrganization_UpdateArgs = {
  input: OrganizationUpdateInput;
};


export type MutationPhoneNumberMergeToContactArgs = {
  contactId: Scalars['ID'];
  input: PhoneNumberInput;
};


export type MutationPhoneNumberRemoveFromContactByE164Args = {
  contactId: Scalars['ID'];
  e164: Scalars['String'];
};


export type MutationPhoneNumberRemoveFromContactByIdArgs = {
  contactId: Scalars['ID'];
  id: Scalars['ID'];
};


export type MutationPhoneNumberUpdateInContactArgs = {
  contactId: Scalars['ID'];
  input: PhoneNumberUpdateInput;
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


export type MutationUser_CreateArgs = {
  input: UserInput;
};


export type MutationUser_UpdateArgs = {
  input: UserUpdateInput;
};

export type Node = {
  id: Scalars['ID'];
};

export type Note = {
  __typename?: 'Note';
  appSource: Scalars['String'];
  createdAt: Scalars['Time'];
  createdBy?: Maybe<User>;
  html: Scalars['String'];
  id: Scalars['ID'];
  source: DataSource;
  sourceOfTruth: DataSource;
  updatedAt: Scalars['Time'];
};

export type NoteInput = {
  appSource?: InputMaybe<Scalars['String']>;
  html: Scalars['String'];
};

export type NotePage = Pages & {
  __typename?: 'NotePage';
  content: Array<Note>;
  totalElements: Scalars['Int64'];
  totalPages: Scalars['Int'];
};

export type NoteUpdateInput = {
  html: Scalars['String'];
  id: Scalars['ID'];
};

export type Organization = Node & {
  __typename?: 'Organization';
  appSource: Scalars['String'];
  contacts: ContactsPage;
  createdAt: Scalars['Time'];
  description?: Maybe<Scalars['String']>;
  /** @deprecated Deprecated in favor of domains */
  domain?: Maybe<Scalars['String']>;
  domains: Array<Scalars['String']>;
  emails: Array<Email>;
  id: Scalars['ID'];
  industry?: Maybe<Scalars['String']>;
  isPublic?: Maybe<Scalars['Boolean']>;
  jobRoles: Array<JobRole>;
  /**
   * All addresses associated with an organization in customerOS.
   * **Required.  If no values it returns an empty array.**
   */
  locations: Array<Location>;
  name: Scalars['String'];
  notes: NotePage;
  organizationType?: Maybe<OrganizationType>;
  phoneNumbers: Array<PhoneNumber>;
  source: DataSource;
  sourceOfTruth: DataSource;
  tags?: Maybe<Array<Tag>>;
  ticketSummaryByStatus: Array<TicketSummaryByStatus>;
  timelineEvents: Array<TimelineEvent>;
  timelineEventsTotalCount: Scalars['Int64'];
  updatedAt: Scalars['Time'];
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
  description?: InputMaybe<Scalars['String']>;
  domain?: InputMaybe<Scalars['String']>;
  domains?: InputMaybe<Array<Scalars['String']>>;
  industry?: InputMaybe<Scalars['String']>;
  isPublic?: InputMaybe<Scalars['Boolean']>;
  /**
   * The name of the organization.
   * **Required.**
   */
  name: Scalars['String'];
  organizationTypeId?: InputMaybe<Scalars['ID']>;
  website?: InputMaybe<Scalars['String']>;
};

export type OrganizationPage = Pages & {
  __typename?: 'OrganizationPage';
  content: Array<Organization>;
  totalElements: Scalars['Int64'];
  totalPages: Scalars['Int'];
};

export type OrganizationType = {
  __typename?: 'OrganizationType';
  createdAt: Scalars['Time'];
  id: Scalars['ID'];
  name: Scalars['String'];
  updatedAt: Scalars['Time'];
};

export type OrganizationTypeInput = {
  name: Scalars['String'];
};

export type OrganizationTypeUpdateInput = {
  id: Scalars['ID'];
  name: Scalars['String'];
};

export type OrganizationUpdateInput = {
  description?: InputMaybe<Scalars['String']>;
  domain?: InputMaybe<Scalars['String']>;
  domains?: InputMaybe<Array<Scalars['String']>>;
  id: Scalars['ID'];
  industry?: InputMaybe<Scalars['String']>;
  isPublic?: InputMaybe<Scalars['Boolean']>;
  name: Scalars['String'];
  organizationTypeId?: InputMaybe<Scalars['ID']>;
  website?: InputMaybe<Scalars['String']>;
};

export type PageView = Node & {
  __typename?: 'PageView';
  application: Scalars['String'];
  endedAt: Scalars['Time'];
  engagedTime: Scalars['Int64'];
  id: Scalars['ID'];
  orderInSession: Scalars['Int64'];
  pageTitle: Scalars['String'];
  pageUrl: Scalars['String'];
  sessionId: Scalars['ID'];
  startedAt: Scalars['Time'];
};

export type PageViewAction = Node & {
  __typename?: 'PageViewAction';
  application: Scalars['String'];
  endedAt: Scalars['Time'];
  engagedTime: Scalars['Int64'];
  id: Scalars['ID'];
  orderInSession: Scalars['Int64'];
  pageTitle: Scalars['String'];
  pageUrl: Scalars['String'];
  sessionId: Scalars['ID'];
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
  /**
   * Determines if the phone number is primary or not.
   * **Required**
   */
  primary: Scalars['Boolean'];
  rawPhoneNumber?: Maybe<Scalars['String']>;
  source: DataSource;
  updatedAt: Scalars['Time'];
  validated?: Maybe<Scalars['Boolean']>;
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

export type Place = {
  __typename?: 'Place';
  address?: Maybe<Scalars['String']>;
  address2?: Maybe<Scalars['String']>;
  appSource?: Maybe<Scalars['String']>;
  city?: Maybe<Scalars['String']>;
  country?: Maybe<Scalars['String']>;
  createdAt: Scalars['Time'];
  fax?: Maybe<Scalars['String']>;
  id: Scalars['ID'];
  phone?: Maybe<Scalars['String']>;
  source?: Maybe<DataSource>;
  state?: Maybe<Scalars['String']>;
  updatedAt: Scalars['Time'];
  zip?: Maybe<Scalars['String']>;
};

export type Query = {
  __typename?: 'Query';
  /** Fetch a single contact from customerOS by contact ID. */
  contact?: Maybe<Contact>;
  /** Fetch a specific contact group associated with a `Contact` in customerOS */
  contactGroup?: Maybe<ContactGroup>;
  /**
   * Fetch paginated list of contact groups
   * Possible values for sort:
   * - NAME
   */
  contactGroups: ContactGroupPage;
  contact_ByEmail: Contact;
  contact_ByPhone: Contact;
  /**
   * Fetch paginated list of contacts
   * Possible values for sort:
   * - TITLE
   * - FIRST_NAME
   * - LAST_NAME
   * - CREATED_AT
   */
  contacts: ContactsPage;
  dashboardView?: Maybe<DashboardViewItemPage>;
  entityTemplates: Array<EntityTemplate>;
  organization?: Maybe<Organization>;
  organizationTypes: Array<OrganizationType>;
  organizations: OrganizationPage;
  search_Basic: Array<SearchBasicResultItem>;
  tags: Array<Tag>;
  user: User;
  user_ByEmail: User;
  users: UserPage;
};


export type QueryContactArgs = {
  id: Scalars['ID'];
};


export type QueryContactGroupArgs = {
  id: Scalars['ID'];
};


export type QueryContactGroupsArgs = {
  pagination?: InputMaybe<Pagination>;
  sort?: InputMaybe<Array<SortBy>>;
  where?: InputMaybe<Filter>;
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


export type QueryDashboardViewArgs = {
  pagination: Pagination;
  searchTerm?: InputMaybe<Scalars['String']>;
};


export type QueryEntityTemplatesArgs = {
  extends?: InputMaybe<EntityTemplateExtension>;
};


export type QueryOrganizationArgs = {
  id: Scalars['ID'];
};


export type QueryOrganizationsArgs = {
  pagination?: InputMaybe<Pagination>;
  sort?: InputMaybe<Array<SortBy>>;
  where?: InputMaybe<Filter>;
};


export type QuerySearch_BasicArgs = {
  keyword: Scalars['String'];
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

export type SearchBasicResult = Contact | Email | Organization;

export type SearchBasicResultItem = {
  __typename?: 'SearchBasicResultItem';
  result: SearchBasicResult;
  score: Scalars['Float'];
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

export type Ticket = Node & {
  __typename?: 'Ticket';
  createdAt: Scalars['Time'];
  description?: Maybe<Scalars['String']>;
  id: Scalars['ID'];
  notes?: Maybe<Array<Maybe<Note>>>;
  priority?: Maybe<Scalars['String']>;
  status?: Maybe<Scalars['String']>;
  subject?: Maybe<Scalars['String']>;
  tags?: Maybe<Array<Maybe<Tag>>>;
  updatedAt: Scalars['Time'];
};

export type TicketSummaryByStatus = {
  __typename?: 'TicketSummaryByStatus';
  count: Scalars['Int64'];
  status: Scalars['String'];
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

export type TimelineEvent = Conversation | InteractionEvent | InteractionSession | Note | PageView | Ticket;

export enum TimelineEventType {
  Conversation = 'CONVERSATION',
  InteractionEvent = 'INTERACTION_EVENT',
  InteractionSession = 'INTERACTION_SESSION',
  Note = 'NOTE',
  PageView = 'PAGE_VIEW',
  Ticket = 'TICKET'
}

/**
 * Describes the User of customerOS.  A user is the person who logs into the Openline platform.
 * **A `return` object**
 */
export type User = {
  __typename?: 'User';
  conversations: ConversationPage;
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
  /**
   * The last name of the customerOS user.
   * **Required**
   */
  lastName: Scalars['String'];
  source: DataSource;
  updatedAt: Scalars['Time'];
};


/**
 * Describes the User of customerOS.  A user is the person who logs into the Openline platform.
 * **A `return` object**
 */
export type UserConversationsArgs = {
  pagination?: InputMaybe<Pagination>;
  sort?: InputMaybe<Array<SortBy>>;
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
   * The first name of the customerOS user.
   * **Required**
   */
  firstName: Scalars['String'];
  /**
   * The last name of the customerOS user.
   * **Required**
   */
  lastName: Scalars['String'];
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

export type AddOrganizationToContactMutationVariables = Exact<{
  input: ContactOrganizationInput;
}>;


export type AddOrganizationToContactMutation = { __typename?: 'Mutation', contact_AddOrganizationById: { __typename?: 'Contact', id: string, source: DataSource, firstName?: string | null, lastName?: string | null, name?: string | null, jobRoles: Array<{ __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string, organization?: { __typename?: 'Organization', id: string, name: string } | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string, createdAt: any, source: DataSource }> | null } };

export type AddPhoneToContactMutationVariables = Exact<{
  contactId: Scalars['ID'];
  input: PhoneNumberInput;
}>;


export type AddPhoneToContactMutation = { __typename?: 'Mutation', phoneNumberMergeToContact: { __typename?: 'PhoneNumber', label?: PhoneNumberLabel | null, id: string, primary: boolean, e164?: string | null, rawPhoneNumber?: string | null } };

export type AddTagToContactMutationVariables = Exact<{
  input: ContactTagInput;
}>;


export type AddTagToContactMutation = { __typename?: 'Mutation', contact_AddTagById: { __typename?: 'Contact', id: string, tags?: Array<{ __typename?: 'Tag', id: string, name: string, createdAt: any, source: DataSource }> | null } };

export type CreateContactMutationVariables = Exact<{
  input: ContactInput;
}>;


export type CreateContactMutation = { __typename?: 'Mutation', contact_Create: { __typename?: 'Contact', id: string, source: DataSource, firstName?: string | null, lastName?: string | null, name?: string | null, jobRoles: Array<{ __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string, organization?: { __typename?: 'Organization', id: string, name: string } | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string, createdAt: any, source: DataSource }> | null, emails: Array<{ __typename?: 'Email', label?: EmailLabel | null, id: string, primary: boolean, email?: string | null }>, phoneNumbers: Array<{ __typename?: 'PhoneNumber', label?: PhoneNumberLabel | null, id: string, primary: boolean, e164?: string | null, rawPhoneNumber?: string | null }> } };

export type CreateContactJobRoleMutationVariables = Exact<{
  contactId: Scalars['ID'];
  input: JobRoleInput;
}>;


export type CreateContactJobRoleMutation = { __typename?: 'Mutation', jobRole_Create: { __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string } };

export type CreateContactNoteMutationVariables = Exact<{
  contactId: Scalars['ID'];
  input: NoteInput;
}>;


export type CreateContactNoteMutation = { __typename?: 'Mutation', note_CreateForContact: { __typename?: 'Note', id: string, html: string, createdAt: any, updatedAt: any, source: DataSource, sourceOfTruth: DataSource, appSource: string, createdBy?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null } };

export type RemoveContactJobRoleMutationVariables = Exact<{
  contactId: Scalars['ID'];
  roleId: Scalars['ID'];
}>;


export type RemoveContactJobRoleMutation = { __typename?: 'Mutation', jobRole_Delete: { __typename?: 'Result', result: boolean } };

export type GetContactCommunicationChannelsQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetContactCommunicationChannelsQuery = { __typename?: 'Query', contact?: { __typename?: 'Contact', firstName?: string | null, lastName?: string | null, name?: string | null, id: string, emails: Array<{ __typename?: 'Email', label?: EmailLabel | null, id: string, primary: boolean, email?: string | null }>, phoneNumbers: Array<{ __typename?: 'PhoneNumber', label?: PhoneNumberLabel | null, id: string, primary: boolean, e164?: string | null, rawPhoneNumber?: string | null }> } | null };

export type GetContactConversationsQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetContactConversationsQuery = { __typename?: 'Query', contact?: { __typename?: 'Contact', conversations: { __typename?: 'ConversationPage', content: Array<{ __typename?: 'Conversation', id: string, startedAt: any }> } } | null };

export type GetContactListQueryVariables = Exact<{
  pagination: Pagination;
  where?: InputMaybe<Filter>;
  sort?: InputMaybe<Array<SortBy> | SortBy>;
}>;


export type GetContactListQuery = { __typename?: 'Query', contacts: { __typename?: 'ContactsPage', content: Array<{ __typename?: 'Contact', firstName?: string | null, lastName?: string | null, name?: string | null, emails: Array<{ __typename?: 'Email', id: string, email?: string | null }> }> } };

export type GetContactNameQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetContactNameQuery = { __typename?: 'Query', contact?: { __typename?: 'Contact', id: string, firstName?: string | null, lastName?: string | null, name?: string | null } | null };

export type GetContactNotesQueryVariables = Exact<{
  id: Scalars['ID'];
  pagination?: InputMaybe<Pagination>;
}>;


export type GetContactNotesQuery = { __typename?: 'Query', contact?: { __typename?: 'Contact', notes: { __typename?: 'NotePage', content: Array<{ __typename?: 'Note', id: string, html: string, createdAt: any, updatedAt: any, source: DataSource, sourceOfTruth: DataSource, appSource: string, createdBy?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null }> } } | null };

export type GetContactPersonalDetailsQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetContactPersonalDetailsQuery = { __typename?: 'Query', contact?: { __typename?: 'Contact', id: string, source: DataSource, firstName?: string | null, lastName?: string | null, name?: string | null, owner?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null, jobRoles: Array<{ __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string, organization?: { __typename?: 'Organization', id: string, name: string } | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string, createdAt: any, source: DataSource }> | null } | null };

export type GetContactTagsQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetContactTagsQuery = { __typename?: 'Query', contact?: { __typename?: 'Contact', id: string, tags?: Array<{ __typename?: 'Tag', id: string, name: string, createdAt: any, source: DataSource }> | null } | null };

export type GetContactTicketsQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetContactTicketsQuery = { __typename?: 'Query', contact?: { __typename?: 'Contact', tickets: Array<{ __typename?: 'Ticket', id: string, createdAt: any, updatedAt: any, subject?: string | null, status?: string | null, priority?: string | null, description?: string | null, tags?: Array<{ __typename?: 'Tag', id: string, name: string } | null> | null, notes?: Array<{ __typename?: 'Note', id: string, createdAt: any, html: string } | null> | null } | null> } | null };

export type GetContactTimelineQueryVariables = Exact<{
  contactId: Scalars['ID'];
  from: Scalars['Time'];
  size: Scalars['Int'];
}>;


export type GetContactTimelineQuery = { __typename?: 'Query', contact?: { __typename?: 'Contact', id: string, firstName?: string | null, lastName?: string | null, name?: string | null, timelineEvents: Array<{ __typename?: 'Conversation', id: string, startedAt: any, subject?: string | null, channel?: string | null, updatedAt: any, messageCount: any, source: DataSource, appSource?: string | null, initiatorFirstName?: string | null, initiatorLastName?: string | null, initiatorUsername?: string | null, initiatorType?: string | null, threadId?: string | null, contacts?: Array<{ __typename?: 'Contact', id: string, lastName?: string | null, firstName?: string | null }> | null, users?: Array<{ __typename?: 'User', lastName: string, firstName: string, emails?: Array<{ __typename?: 'Email', email?: string | null }> | null }> | null } | { __typename?: 'InteractionEvent', id: string, createdAt: any, channel?: string | null, content?: string | null, contentType?: string | null, interactionSession?: { __typename?: 'InteractionSession', name?: string | null } | null, sentBy: Array<{ __typename: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', name?: string | null } } | { __typename: 'EmailParticipant', emailParticipant: { __typename?: 'Email', email?: string | null } } | { __typename: 'PhoneNumberParticipant', phoneNumberParticipant: { __typename?: 'PhoneNumber', e164?: string | null } } | { __typename: 'UserParticipant', userParticipant: { __typename?: 'User', firstName: string, lastName: string } }>, sentTo: Array<{ __typename: 'ContactParticipant', type?: string | null, contactParticipant: { __typename?: 'Contact', name?: string | null } } | { __typename: 'EmailParticipant', type?: string | null, emailParticipant: { __typename?: 'Email', email?: string | null } } | { __typename: 'PhoneNumberParticipant', type?: string | null, phoneNumberParticipant: { __typename?: 'PhoneNumber', e164?: string | null } } | { __typename: 'UserParticipant', type?: string | null, userParticipant: { __typename?: 'User', firstName: string, lastName: string } }> } | { __typename?: 'InteractionSession', id: string, startedAt: any, name?: string | null, status?: string | null, type?: string | null, events: Array<{ __typename?: 'InteractionEvent', content?: string | null, contentType?: string | null }> } | { __typename?: 'Note', id: string, html: string, createdAt: any, createdBy?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null } | { __typename?: 'PageView', id: string, application: string, startedAt: any, endedAt: any, engagedTime: any, pageUrl: string, pageTitle: string, orderInSession: any, sessionId: string } | { __typename?: 'Ticket', id: string, createdAt: any, updatedAt: any, subject?: string | null, status?: string | null, priority?: string | null, description?: string | null, notes?: Array<{ __typename?: 'Note', id: string, html: string } | null> | null, tags?: Array<{ __typename?: 'Tag', id: string, name: string } | null> | null }> } | null };

export type MergeContactsMutationVariables = Exact<{
  primaryContactId: Scalars['ID'];
  mergedContactIds: Array<Scalars['ID']> | Scalars['ID'];
}>;


export type MergeContactsMutation = { __typename?: 'Mutation', contact_Merge: { __typename?: 'Contact', id: string, source: DataSource, firstName?: string | null, lastName?: string | null, name?: string | null, jobRoles: Array<{ __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string, organization?: { __typename?: 'Organization', id: string, name: string } | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string, createdAt: any, source: DataSource }> | null } };

export type RemoveEmailFromContactMutationVariables = Exact<{
  contactId: Scalars['ID'];
  id: Scalars['ID'];
}>;


export type RemoveEmailFromContactMutation = { __typename?: 'Mutation', emailRemoveFromContactById: { __typename?: 'Result', result: boolean } };

export type RemoveOrganizationToContactMutationVariables = Exact<{
  input: ContactOrganizationInput;
}>;


export type RemoveOrganizationToContactMutation = { __typename?: 'Mutation', contact_RemoveOrganizationById: { __typename?: 'Contact', id: string, source: DataSource, firstName?: string | null, lastName?: string | null, name?: string | null, jobRoles: Array<{ __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string, organization?: { __typename?: 'Organization', id: string, name: string } | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string, createdAt: any, source: DataSource }> | null } };

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


export type UpdateContactEmailMutation = { __typename?: 'Mutation', emailUpdateInContact: { __typename?: 'Email', primary: boolean, label?: EmailLabel | null } };

export type UpdateJobRoleMutationVariables = Exact<{
  contactId: Scalars['ID'];
  input: JobRoleUpdateInput;
}>;


export type UpdateJobRoleMutation = { __typename?: 'Mutation', jobRole_Update: { __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string } };

export type UpdateContactPersonalDetailsMutationVariables = Exact<{
  input: ContactUpdateInput;
}>;


export type UpdateContactPersonalDetailsMutation = { __typename?: 'Mutation', contact_Update: { __typename?: 'Contact', id: string, title?: PersonTitle | null, firstName?: string | null, lastName?: string | null } };

export type UpdateContactPhoneNumberMutationVariables = Exact<{
  contactId: Scalars['ID'];
  input: PhoneNumberUpdateInput;
}>;


export type UpdateContactPhoneNumberMutation = { __typename?: 'Mutation', phoneNumberUpdateInContact: { __typename?: 'PhoneNumber', label?: PhoneNumberLabel | null, id: string, primary: boolean, e164?: string | null, rawPhoneNumber?: string | null } };

export type GetDashboardDataQueryVariables = Exact<{
  pagination: Pagination;
  searchTerm?: InputMaybe<Scalars['String']>;
}>;


export type GetDashboardDataQuery = { __typename?: 'Query', dashboardView?: { __typename?: 'DashboardViewItemPage', totalElements: any, content: Array<{ __typename?: 'DashboardViewItem', contact?: { __typename?: 'Contact', id: string, firstName?: string | null, lastName?: string | null, name?: string | null, jobRoles: Array<{ __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string }>, emails: Array<{ __typename?: 'Email', id: string, primary: boolean, email?: string | null }>, locations: Array<{ __typename?: 'Location', id: string, name: string, country?: string | null, region?: string | null, locality?: string | null }> } | null, organization?: { __typename?: 'Organization', id: string, name: string, industry?: string | null } | null }> } | null };

export type LocationBaseDetailsFragment = { __typename?: 'Location', id: string, name: string, country?: string | null, region?: string | null, locality?: string | null };

export type LocationTotalFragment = { __typename?: 'Location', id: string, name: string, createdAt: any, updatedAt: any, source?: DataSource | null, appSource?: string | null, country?: string | null, region?: string | null, locality?: string | null, address?: string | null, address2?: string | null, zip?: string | null, addressType?: string | null, houseNumber?: string | null, postalCode?: string | null, plusFour?: string | null, commercial?: boolean | null, predirection?: string | null, district?: string | null, street?: string | null, rawAddress?: string | null, latitude?: number | null, longitude?: number | null };

export type JobRoleFragment = { __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string };

export type NoteContentFragment = { __typename?: 'Note', id: string, html: string, createdAt: any, updatedAt: any, source: DataSource, sourceOfTruth: DataSource, appSource: string, createdBy?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null };

export type TagFragment = { __typename?: 'Tag', id: string, name: string, createdAt: any, source: DataSource };

export type EmailFragment = { __typename?: 'Email', id: string, primary: boolean, email?: string | null };

export type PhoneNumberFragment = { __typename?: 'PhoneNumber', id: string, primary: boolean, e164?: string | null, rawPhoneNumber?: string | null };

export type ConversationFragment = { __typename?: 'Conversation', id: string, startedAt: any, updatedAt: any };

export type TicketFragment = { __typename?: 'Ticket', id: string, createdAt: any, updatedAt: any, subject?: string | null, status?: string | null, priority?: string | null, description?: string | null, tags?: Array<{ __typename?: 'Tag', id: string, name: string, createdAt: any, source: DataSource } | null> | null, notes?: Array<{ __typename?: 'Note', id: string, html: string, createdAt: any, updatedAt: any, source: DataSource, sourceOfTruth: DataSource, appSource: string, createdBy?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null } | null> | null };

export type ContactNameFragmentFragment = { __typename?: 'Contact', firstName?: string | null, lastName?: string | null, name?: string | null };

export type OrganizationBaseDetailsFragment = { __typename?: 'Organization', id: string, name: string, industry?: string | null };

export type ContactPersonalDetailsFragment = { __typename?: 'Contact', id: string, source: DataSource, firstName?: string | null, lastName?: string | null, name?: string | null, jobRoles: Array<{ __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string, organization?: { __typename?: 'Organization', id: string, name: string } | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string, createdAt: any, source: DataSource }> | null };

export type ContactCommunicationChannelsDetailsFragment = { __typename?: 'Contact', id: string, emails: Array<{ __typename?: 'Email', label?: EmailLabel | null, id: string, primary: boolean, email?: string | null }>, phoneNumbers: Array<{ __typename?: 'PhoneNumber', label?: PhoneNumberLabel | null, id: string, primary: boolean, e164?: string | null, rawPhoneNumber?: string | null }> };

export type OrganizationDetailsFragment = { __typename?: 'Organization', id: string, name: string, description?: string | null, source: DataSource, industry?: string | null, website?: string | null, updatedAt: any, locations: Array<{ __typename?: 'Location', id: string, name: string, country?: string | null, region?: string | null, locality?: string | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string, createdAt: any, source: DataSource }> | null };

export type OrganizationContactsFragment = { __typename?: 'Organization', contacts: { __typename?: 'ContactsPage', content: Array<{ __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null, jobRoles: Array<{ __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string }>, emails: Array<{ __typename?: 'Email', label?: EmailLabel | null, id: string, primary: boolean, email?: string | null }>, phoneNumbers: Array<{ __typename?: 'PhoneNumber', label?: PhoneNumberLabel | null, id: string, primary: boolean, e164?: string | null, rawPhoneNumber?: string | null }> }> } };

export type CreateOrganizationMutationVariables = Exact<{
  input: OrganizationInput;
}>;


export type CreateOrganizationMutation = { __typename?: 'Mutation', organization_Create: { __typename?: 'Organization', id: string } };

export type CreateOrganizationNoteMutationVariables = Exact<{
  organizationId: Scalars['ID'];
  input: NoteInput;
}>;


export type CreateOrganizationNoteMutation = { __typename?: 'Mutation', note_CreateForOrganization: { __typename?: 'Note', id: string, html: string, createdAt: any, updatedAt: any, source: DataSource, sourceOfTruth: DataSource, appSource: string, createdBy?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null } };

export type GetOrganizationContactsQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetOrganizationContactsQuery = { __typename?: 'Query', organization?: { __typename?: 'Organization', contacts: { __typename?: 'ContactsPage', content: Array<{ __typename?: 'Contact', id: string, name?: string | null, firstName?: string | null, lastName?: string | null, jobRoles: Array<{ __typename?: 'JobRole', jobTitle?: string | null, primary: boolean, id: string }>, emails: Array<{ __typename?: 'Email', label?: EmailLabel | null, id: string, primary: boolean, email?: string | null }>, phoneNumbers: Array<{ __typename?: 'PhoneNumber', label?: PhoneNumberLabel | null, id: string, primary: boolean, e164?: string | null, rawPhoneNumber?: string | null }> }> } } | null };

export type GetOrganizationDetailsQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetOrganizationDetailsQuery = { __typename?: 'Query', organization?: { __typename?: 'Organization', id: string, name: string, description?: string | null, source: DataSource, industry?: string | null, website?: string | null, updatedAt: any, locations: Array<{ __typename?: 'Location', id: string, name: string, country?: string | null, region?: string | null, locality?: string | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string, createdAt: any, source: DataSource }> | null } | null };

export type GetOrganizationNameQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetOrganizationNameQuery = { __typename?: 'Query', organization?: { __typename?: 'Organization', id: string, name: string } | null };

export type GetOrganizationNotesQueryVariables = Exact<{
  id: Scalars['ID'];
  pagination?: InputMaybe<Pagination>;
}>;


export type GetOrganizationNotesQuery = { __typename?: 'Query', organization?: { __typename?: 'Organization', notes: { __typename?: 'NotePage', content: Array<{ __typename?: 'Note', id: string, html: string, createdAt: any, updatedAt: any, source: DataSource, sourceOfTruth: DataSource, appSource: string, createdBy?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null }> } } | null };

export type GetOrganizationTimelineQueryVariables = Exact<{
  organizationId: Scalars['ID'];
  from: Scalars['Time'];
  size: Scalars['Int'];
}>;


export type GetOrganizationTimelineQuery = { __typename?: 'Query', organization?: { __typename?: 'Organization', id: string, timelineEvents: Array<{ __typename?: 'Conversation', id: string, startedAt: any, subject?: string | null, channel?: string | null, updatedAt: any, messageCount: any, source: DataSource, appSource?: string | null, initiatorFirstName?: string | null, initiatorLastName?: string | null, initiatorUsername?: string | null, initiatorType?: string | null, threadId?: string | null, contacts?: Array<{ __typename?: 'Contact', id: string, lastName?: string | null, firstName?: string | null }> | null, users?: Array<{ __typename?: 'User', lastName: string, firstName: string, emails?: Array<{ __typename?: 'Email', email?: string | null }> | null }> | null } | { __typename?: 'InteractionEvent', id: string, createdAt: any, channel?: string | null, content?: string | null, contentType?: string | null, interactionSession?: { __typename?: 'InteractionSession', name?: string | null } | null, sentBy: Array<{ __typename: 'ContactParticipant', contactParticipant: { __typename?: 'Contact', name?: string | null } } | { __typename: 'EmailParticipant', emailParticipant: { __typename?: 'Email', email?: string | null } } | { __typename: 'PhoneNumberParticipant', phoneNumberParticipant: { __typename?: 'PhoneNumber', e164?: string | null } } | { __typename: 'UserParticipant', userParticipant: { __typename?: 'User', firstName: string, lastName: string } }>, sentTo: Array<{ __typename: 'ContactParticipant', type?: string | null, contactParticipant: { __typename?: 'Contact', name?: string | null } } | { __typename: 'EmailParticipant', type?: string | null, emailParticipant: { __typename?: 'Email', email?: string | null } } | { __typename: 'PhoneNumberParticipant', type?: string | null, phoneNumberParticipant: { __typename?: 'PhoneNumber', e164?: string | null } } | { __typename: 'UserParticipant', type?: string | null, userParticipant: { __typename?: 'User', firstName: string, lastName: string } }> } | { __typename?: 'InteractionSession', id: string, startedAt: any, name?: string | null, status?: string | null, type?: string | null, events: Array<{ __typename?: 'InteractionEvent', content?: string | null, contentType?: string | null }> } | { __typename?: 'Note', id: string, html: string, createdAt: any, createdBy?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null } | { __typename?: 'PageView', id: string, application: string, startedAt: any, endedAt: any, engagedTime: any, pageUrl: string, pageTitle: string, orderInSession: any, sessionId: string } | { __typename?: 'Ticket', id: string, createdAt: any, updatedAt: any, subject?: string | null, status?: string | null, priority?: string | null, description?: string | null, notes?: Array<{ __typename?: 'Note', id: string, html: string } | null> | null, tags?: Array<{ __typename?: 'Tag', id: string, name: string } | null> | null }> } | null };

export type GetOrganizationsOptionsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetOrganizationsOptionsQuery = { __typename?: 'Query', organizations: { __typename?: 'OrganizationPage', content: Array<{ __typename?: 'Organization', id: string, name: string }> } };

export type MergeOrganizationsMutationVariables = Exact<{
  primaryOrganizationId: Scalars['ID'];
  mergedOrganizationIds: Array<Scalars['ID']> | Scalars['ID'];
}>;


export type MergeOrganizationsMutation = { __typename?: 'Mutation', organization_Merge: { __typename?: 'Organization', id: string, name: string, description?: string | null, source: DataSource, industry?: string | null, website?: string | null, updatedAt: any, locations: Array<{ __typename?: 'Location', id: string, name: string, country?: string | null, region?: string | null, locality?: string | null }>, tags?: Array<{ __typename?: 'Tag', id: string, name: string, createdAt: any, source: DataSource }> | null } };

export type UpdateOrganizationMutationVariables = Exact<{
  input: OrganizationUpdateInput;
}>;


export type UpdateOrganizationMutation = { __typename?: 'Mutation', organization_Update: { __typename?: 'Organization', id: string } };

export type RemoveNoteMutationVariables = Exact<{
  id: Scalars['ID'];
}>;


export type RemoveNoteMutation = { __typename?: 'Mutation', note_Delete: { __typename?: 'Result', result: boolean } };

export type UpdateNoteMutationVariables = Exact<{
  input: NoteUpdateInput;
}>;


export type UpdateNoteMutation = { __typename?: 'Mutation', note_Update: { __typename?: 'Note', id: string, html: string, createdAt: any, updatedAt: any, source: DataSource, sourceOfTruth: DataSource, appSource: string, createdBy?: { __typename?: 'User', id: string, firstName: string, lastName: string } | null } };

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
export const ConversationFragmentDoc = gql`
    fragment Conversation on Conversation {
  id
  startedAt
  updatedAt
}
    `;
export const TagFragmentDoc = gql`
    fragment Tag on Tag {
  id
  name
  createdAt
  source
}
    `;
export const NoteContentFragmentDoc = gql`
    fragment NoteContent on Note {
  id
  html
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
}
    `;
export const TicketFragmentDoc = gql`
    fragment Ticket on Ticket {
  id
  createdAt
  updatedAt
  subject
  status
  priority
  description
  tags {
    ...Tag
  }
  notes {
    ...NoteContent
  }
}
    ${TagFragmentDoc}
${NoteContentFragmentDoc}`;
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
export const LocationBaseDetailsFragmentDoc = gql`
    fragment LocationBaseDetails on Location {
  id
  name
  country
  region
  locality
}
    `;
export const OrganizationDetailsFragmentDoc = gql`
    fragment OrganizationDetails on Organization {
  id
  name
  description
  source
  industry
  locations {
    ...LocationBaseDetails
  }
  website
  updatedAt
  tags {
    ...Tag
  }
}
    ${LocationBaseDetailsFragmentDoc}
${TagFragmentDoc}`;
export const EmailFragmentDoc = gql`
    fragment Email on Email {
  id
  primary
  email
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
    ...Email
  }
  phoneNumbers {
    label
    ...PhoneNumber
  }
}
    ${EmailFragmentDoc}
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
export const AddOrganizationToContactDocument = gql`
    mutation addOrganizationToContact($input: ContactOrganizationInput!) {
  contact_AddOrganizationById(input: $input) {
    ...ContactPersonalDetails
  }
}
    ${ContactPersonalDetailsFragmentDoc}`;
export type AddOrganizationToContactMutationFn = Apollo.MutationFunction<AddOrganizationToContactMutation, AddOrganizationToContactMutationVariables>;

/**
 * __useAddOrganizationToContactMutation__
 *
 * To run a mutation, you first call `useAddOrganizationToContactMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddOrganizationToContactMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addOrganizationToContactMutation, { data, loading, error }] = useAddOrganizationToContactMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useAddOrganizationToContactMutation(baseOptions?: Apollo.MutationHookOptions<AddOrganizationToContactMutation, AddOrganizationToContactMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddOrganizationToContactMutation, AddOrganizationToContactMutationVariables>(AddOrganizationToContactDocument, options);
      }
export type AddOrganizationToContactMutationHookResult = ReturnType<typeof useAddOrganizationToContactMutation>;
export type AddOrganizationToContactMutationResult = Apollo.MutationResult<AddOrganizationToContactMutation>;
export type AddOrganizationToContactMutationOptions = Apollo.BaseMutationOptions<AddOrganizationToContactMutation, AddOrganizationToContactMutationVariables>;
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
export const GetContactConversationsDocument = gql`
    query GetContactConversations($id: ID!) {
  contact(id: $id) {
    conversations(
      pagination: {page: 0, limit: 25}
      sort: {by: "STARTED_AT", direction: DESC}
    ) {
      content {
        id
        startedAt
      }
    }
  }
}
    `;

/**
 * __useGetContactConversationsQuery__
 *
 * To run a query within a React component, call `useGetContactConversationsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetContactConversationsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetContactConversationsQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetContactConversationsQuery(baseOptions: Apollo.QueryHookOptions<GetContactConversationsQuery, GetContactConversationsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetContactConversationsQuery, GetContactConversationsQueryVariables>(GetContactConversationsDocument, options);
      }
export function useGetContactConversationsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetContactConversationsQuery, GetContactConversationsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetContactConversationsQuery, GetContactConversationsQueryVariables>(GetContactConversationsDocument, options);
        }
export type GetContactConversationsQueryHookResult = ReturnType<typeof useGetContactConversationsQuery>;
export type GetContactConversationsLazyQueryHookResult = ReturnType<typeof useGetContactConversationsLazyQuery>;
export type GetContactConversationsQueryResult = Apollo.QueryResult<GetContactConversationsQuery, GetContactConversationsQueryVariables>;
export const GetContactListDocument = gql`
    query GetContactList($pagination: Pagination!, $where: Filter, $sort: [SortBy!]) {
  contacts(pagination: $pagination, where: $where, sort: $sort) {
    content {
      ...ContactNameFragment
      emails {
        id
        email
      }
    }
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
export const GetContactNameDocument = gql`
    query GetContactName($id: ID!) {
  contact(id: $id) {
    id
    ...ContactNameFragment
  }
}
    ${ContactNameFragmentFragmentDoc}`;

/**
 * __useGetContactNameQuery__
 *
 * To run a query within a React component, call `useGetContactNameQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetContactNameQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetContactNameQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetContactNameQuery(baseOptions: Apollo.QueryHookOptions<GetContactNameQuery, GetContactNameQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetContactNameQuery, GetContactNameQueryVariables>(GetContactNameDocument, options);
      }
export function useGetContactNameLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetContactNameQuery, GetContactNameQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetContactNameQuery, GetContactNameQueryVariables>(GetContactNameDocument, options);
        }
export type GetContactNameQueryHookResult = ReturnType<typeof useGetContactNameQuery>;
export type GetContactNameLazyQueryHookResult = ReturnType<typeof useGetContactNameLazyQuery>;
export type GetContactNameQueryResult = Apollo.QueryResult<GetContactNameQuery, GetContactNameQueryVariables>;
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
export const GetContactTicketsDocument = gql`
    query GetContactTickets($id: ID!) {
  contact(id: $id) {
    tickets {
      id
      createdAt
      updatedAt
      subject
      status
      priority
      description
      tags {
        id
        name
      }
      notes {
        id
        createdAt
        html
      }
    }
  }
}
    `;

/**
 * __useGetContactTicketsQuery__
 *
 * To run a query within a React component, call `useGetContactTicketsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetContactTicketsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetContactTicketsQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetContactTicketsQuery(baseOptions: Apollo.QueryHookOptions<GetContactTicketsQuery, GetContactTicketsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetContactTicketsQuery, GetContactTicketsQueryVariables>(GetContactTicketsDocument, options);
      }
export function useGetContactTicketsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetContactTicketsQuery, GetContactTicketsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetContactTicketsQuery, GetContactTicketsQueryVariables>(GetContactTicketsDocument, options);
        }
export type GetContactTicketsQueryHookResult = ReturnType<typeof useGetContactTicketsQuery>;
export type GetContactTicketsLazyQueryHookResult = ReturnType<typeof useGetContactTicketsLazyQuery>;
export type GetContactTicketsQueryResult = Apollo.QueryResult<GetContactTicketsQuery, GetContactTicketsQueryVariables>;
export const GetContactTimelineDocument = gql`
    query GetContactTimeline($contactId: ID!, $from: Time!, $size: Int!) {
  contact(id: $contactId) {
    id
    ...ContactNameFragment
    timelineEvents(from: $from, size: $size) {
      ... on PageView {
        id
        application
        startedAt
        endedAt
        engagedTime
        pageUrl
        pageTitle
        orderInSession
        sessionId
      }
      ... on Ticket {
        id
        createdAt
        updatedAt
        subject
        status
        priority
        description
        notes {
          id
          html
        }
        tags {
          id
          name
        }
      }
      ... on Conversation {
        id
        startedAt
        subject
        channel
        updatedAt
        messageCount
        contacts {
          id
          lastName
          firstName
        }
        users {
          lastName
          firstName
          emails {
            email
          }
        }
        source
        appSource
        initiatorFirstName
        initiatorLastName
        initiatorUsername
        initiatorType
        threadId
      }
      ... on InteractionSession {
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
      ... on InteractionEvent {
        id
        createdAt
        channel
        interactionSession {
          name
        }
        content
        contentType
        sentBy {
          __typename
          ... on EmailParticipant {
            __typename
            emailParticipant {
              email
            }
          }
          ... on PhoneNumberParticipant {
            __typename
            phoneNumberParticipant {
              e164
            }
          }
          ... on ContactParticipant {
            __typename
            contactParticipant {
              name
            }
          }
          ... on UserParticipant {
            __typename
            userParticipant {
              firstName
              lastName
            }
          }
        }
        sentTo {
          __typename
          ... on EmailParticipant {
            __typename
            type
            emailParticipant {
              email
            }
          }
          ... on PhoneNumberParticipant {
            __typename
            type
            phoneNumberParticipant {
              e164
            }
          }
          ... on ContactParticipant {
            __typename
            type
            contactParticipant {
              name
            }
          }
          ... on UserParticipant {
            __typename
            type
            userParticipant {
              firstName
              lastName
            }
          }
        }
      }
      ... on Note {
        id
        html
        createdAt
        createdBy {
          id
          firstName
          lastName
        }
      }
    }
  }
}
    ${ContactNameFragmentFragmentDoc}`;

/**
 * __useGetContactTimelineQuery__
 *
 * To run a query within a React component, call `useGetContactTimelineQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetContactTimelineQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetContactTimelineQuery({
 *   variables: {
 *      contactId: // value for 'contactId'
 *      from: // value for 'from'
 *      size: // value for 'size'
 *   },
 * });
 */
export function useGetContactTimelineQuery(baseOptions: Apollo.QueryHookOptions<GetContactTimelineQuery, GetContactTimelineQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetContactTimelineQuery, GetContactTimelineQueryVariables>(GetContactTimelineDocument, options);
      }
export function useGetContactTimelineLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetContactTimelineQuery, GetContactTimelineQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetContactTimelineQuery, GetContactTimelineQueryVariables>(GetContactTimelineDocument, options);
        }
export type GetContactTimelineQueryHookResult = ReturnType<typeof useGetContactTimelineQuery>;
export type GetContactTimelineLazyQueryHookResult = ReturnType<typeof useGetContactTimelineLazyQuery>;
export type GetContactTimelineQueryResult = Apollo.QueryResult<GetContactTimelineQuery, GetContactTimelineQueryVariables>;
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
export const RemoveOrganizationToContactDocument = gql`
    mutation removeOrganizationToContact($input: ContactOrganizationInput!) {
  contact_RemoveOrganizationById(input: $input) {
    ...ContactPersonalDetails
  }
}
    ${ContactPersonalDetailsFragmentDoc}`;
export type RemoveOrganizationToContactMutationFn = Apollo.MutationFunction<RemoveOrganizationToContactMutation, RemoveOrganizationToContactMutationVariables>;

/**
 * __useRemoveOrganizationToContactMutation__
 *
 * To run a mutation, you first call `useRemoveOrganizationToContactMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemoveOrganizationToContactMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removeOrganizationToContactMutation, { data, loading, error }] = useRemoveOrganizationToContactMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useRemoveOrganizationToContactMutation(baseOptions?: Apollo.MutationHookOptions<RemoveOrganizationToContactMutation, RemoveOrganizationToContactMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RemoveOrganizationToContactMutation, RemoveOrganizationToContactMutationVariables>(RemoveOrganizationToContactDocument, options);
      }
export type RemoveOrganizationToContactMutationHookResult = ReturnType<typeof useRemoveOrganizationToContactMutation>;
export type RemoveOrganizationToContactMutationResult = Apollo.MutationResult<RemoveOrganizationToContactMutation>;
export type RemoveOrganizationToContactMutationOptions = Apollo.BaseMutationOptions<RemoveOrganizationToContactMutation, RemoveOrganizationToContactMutationVariables>;
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
export const GetDashboardDataDocument = gql`
    query GetDashboardData($pagination: Pagination!, $searchTerm: String) {
  dashboardView(pagination: $pagination, searchTerm: $searchTerm) {
    content {
      contact {
        id
        ...ContactNameFragment
        jobRoles {
          ...JobRole
        }
        emails {
          ...Email
        }
        locations {
          ...LocationBaseDetails
        }
      }
      organization {
        ...organizationBaseDetails
      }
    }
    totalElements
  }
}
    ${ContactNameFragmentFragmentDoc}
${JobRoleFragmentDoc}
${EmailFragmentDoc}
${LocationBaseDetailsFragmentDoc}
${OrganizationBaseDetailsFragmentDoc}`;

/**
 * __useGetDashboardDataQuery__
 *
 * To run a query within a React component, call `useGetDashboardDataQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetDashboardDataQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetDashboardDataQuery({
 *   variables: {
 *      pagination: // value for 'pagination'
 *      searchTerm: // value for 'searchTerm'
 *   },
 * });
 */
export function useGetDashboardDataQuery(baseOptions: Apollo.QueryHookOptions<GetDashboardDataQuery, GetDashboardDataQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetDashboardDataQuery, GetDashboardDataQueryVariables>(GetDashboardDataDocument, options);
      }
export function useGetDashboardDataLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetDashboardDataQuery, GetDashboardDataQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetDashboardDataQuery, GetDashboardDataQueryVariables>(GetDashboardDataDocument, options);
        }
export type GetDashboardDataQueryHookResult = ReturnType<typeof useGetDashboardDataQuery>;
export type GetDashboardDataLazyQueryHookResult = ReturnType<typeof useGetDashboardDataLazyQuery>;
export type GetDashboardDataQueryResult = Apollo.QueryResult<GetDashboardDataQuery, GetDashboardDataQueryVariables>;
export const CreateOrganizationDocument = gql`
    mutation createOrganization($input: OrganizationInput!) {
  organization_Create(input: $input) {
    id
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
export const GetOrganizationDetailsDocument = gql`
    query GetOrganizationDetails($id: ID!) {
  organization(id: $id) {
    ...OrganizationDetails
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
export const GetOrganizationTimelineDocument = gql`
    query GetOrganizationTimeline($organizationId: ID!, $from: Time!, $size: Int!) {
  organization(id: $organizationId) {
    id
    timelineEvents(from: $from, size: $size) {
      ... on PageView {
        id
        application
        startedAt
        endedAt
        engagedTime
        pageUrl
        pageTitle
        orderInSession
        sessionId
      }
      ... on Ticket {
        id
        createdAt
        updatedAt
        subject
        status
        priority
        description
        notes {
          id
          html
        }
        tags {
          id
          name
        }
      }
      ... on Conversation {
        id
        startedAt
        subject
        channel
        updatedAt
        messageCount
        contacts {
          id
          lastName
          firstName
        }
        users {
          lastName
          firstName
          emails {
            email
          }
        }
        source
        appSource
        initiatorFirstName
        initiatorLastName
        initiatorUsername
        initiatorType
        threadId
      }
      ... on InteractionSession {
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
      ... on InteractionEvent {
        id
        createdAt
        channel
        interactionSession {
          name
        }
        content
        contentType
        sentBy {
          __typename
          ... on EmailParticipant {
            __typename
            emailParticipant {
              email
            }
          }
          ... on PhoneNumberParticipant {
            __typename
            phoneNumberParticipant {
              e164
            }
          }
          ... on ContactParticipant {
            __typename
            contactParticipant {
              name
            }
          }
          ... on UserParticipant {
            __typename
            userParticipant {
              firstName
              lastName
            }
          }
        }
        sentTo {
          __typename
          ... on EmailParticipant {
            __typename
            type
            emailParticipant {
              email
            }
          }
          ... on PhoneNumberParticipant {
            __typename
            type
            phoneNumberParticipant {
              e164
            }
          }
          ... on ContactParticipant {
            __typename
            type
            contactParticipant {
              name
            }
          }
          ... on UserParticipant {
            __typename
            type
            userParticipant {
              firstName
              lastName
            }
          }
        }
      }
      ... on Note {
        id
        html
        createdAt
        createdBy {
          id
          firstName
          lastName
        }
      }
    }
  }
}
    `;

/**
 * __useGetOrganizationTimelineQuery__
 *
 * To run a query within a React component, call `useGetOrganizationTimelineQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetOrganizationTimelineQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetOrganizationTimelineQuery({
 *   variables: {
 *      organizationId: // value for 'organizationId'
 *      from: // value for 'from'
 *      size: // value for 'size'
 *   },
 * });
 */
export function useGetOrganizationTimelineQuery(baseOptions: Apollo.QueryHookOptions<GetOrganizationTimelineQuery, GetOrganizationTimelineQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetOrganizationTimelineQuery, GetOrganizationTimelineQueryVariables>(GetOrganizationTimelineDocument, options);
      }
export function useGetOrganizationTimelineLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetOrganizationTimelineQuery, GetOrganizationTimelineQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetOrganizationTimelineQuery, GetOrganizationTimelineQueryVariables>(GetOrganizationTimelineDocument, options);
        }
export type GetOrganizationTimelineQueryHookResult = ReturnType<typeof useGetOrganizationTimelineQuery>;
export type GetOrganizationTimelineLazyQueryHookResult = ReturnType<typeof useGetOrganizationTimelineLazyQuery>;
export type GetOrganizationTimelineQueryResult = Apollo.QueryResult<GetOrganizationTimelineQuery, GetOrganizationTimelineQueryVariables>;
export const GetOrganizationsOptionsDocument = gql`
    query getOrganizationsOptions {
  organizations {
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
export const UpdateOrganizationDocument = gql`
    mutation updateOrganization($input: OrganizationUpdateInput!) {
  organization_Update(input: $input) {
    id
  }
}
    `;
export type UpdateOrganizationMutationFn = Apollo.MutationFunction<UpdateOrganizationMutation, UpdateOrganizationMutationVariables>;

/**
 * __useUpdateOrganizationMutation__
 *
 * To run a mutation, you first call `useUpdateOrganizationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateOrganizationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateOrganizationMutation, { data, loading, error }] = useUpdateOrganizationMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateOrganizationMutation(baseOptions?: Apollo.MutationHookOptions<UpdateOrganizationMutation, UpdateOrganizationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateOrganizationMutation, UpdateOrganizationMutationVariables>(UpdateOrganizationDocument, options);
      }
export type UpdateOrganizationMutationHookResult = ReturnType<typeof useUpdateOrganizationMutation>;
export type UpdateOrganizationMutationResult = Apollo.MutationResult<UpdateOrganizationMutation>;
export type UpdateOrganizationMutationOptions = Apollo.BaseMutationOptions<UpdateOrganizationMutation, UpdateOrganizationMutationVariables>;
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