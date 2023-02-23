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

export type Action = PageViewAction;

export enum ActionType {
  PageView = 'PAGE_VIEW',
}

export enum ComparisonOperator {
  Contains = 'CONTAINS',
  Eq = 'EQ',
}

/**
 * A contact represents an individual in customerOS.
 * **A `response` object.**
 */
export type Contact = ExtensibleEntity &
  Node & {
    __typename?: 'Contact';
    actions: Array<Action>;
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
    /** A user-defined label applied against a contact in customerOS. */
    label?: Maybe<Scalars['String']>;
    /** The last name of the contact in customerOS. */
    lastName?: Maybe<Scalars['String']>;
    /**
     * All locations associated with a contact in customerOS.
     * **Required.  If no values it returns an empty array.**
     */
    locations: Array<Location>;
    /** Contact notes */
    notes: NotePage;
    organizations: OrganizationPage;
    /** Contact owner (user) */
    owner?: Maybe<User>;
    /**
     * All phone numbers associated with a contact in customerOS.
     * **Required.  If no values it returns an empty array.**
     */
    phoneNumbers: Array<PhoneNumber>;
    source: DataSource;
    tags?: Maybe<Array<Tag>>;
    /** Template of the contact in customerOS. */
    template?: Maybe<EntityTemplate>;
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
export type ContactOrganizationsArgs = {
  pagination?: InputMaybe<Pagination>;
  sort?: InputMaybe<Array<SortBy>>;
  where?: InputMaybe<Filter>;
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
  /** A user-defined label attached to contact. */
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
  /** A user-defined label applied against a contact in customerOS. */
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
  Closed = 'CLOSED',
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
  Text = 'TEXT',
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
  Zendesk = 'ZENDESK',
}

/**
 * Describes an email address associated with a `Contact` in customerOS.
 * **A `return` object.**
 */
export type Email = {
  __typename?: 'Email';
  appSource: Scalars['String'];
  createdAt: Scalars['Time'];
  /**
   * An email address assocaited with the contact in customerOS.
   * **Required.**
   */
  email: Scalars['String'];
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
  source: DataSource;
  sourceOfTruth: DataSource;
  updatedAt: Scalars['Time'];
};

/**
 * Describes an email address associated with a `Contact` in customerOS.
 * **A `create` object.**
 */
export type EmailInput = {
  appSource?: InputMaybe<Scalars['String']>;
  /**
   * An email address assocaited with the contact in customerOS.
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
  Work = 'WORK',
}

/**
 * Describes an email address associated with a `Contact` in customerOS.
 * **An `update` object.**
 */
export type EmailUpdateInput = {
  /**
   * An email address assocaited with the contact in customerOS.
   * **Required.**
   */
  email: Scalars['String'];
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
  Contact = 'CONTACT',
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
  Zendesk = 'ZENDESK',
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
  appSource?: Maybe<Scalars['String']>;
  country?: Maybe<Scalars['String']>;
  createdAt: Scalars['Time'];
  id: Scalars['ID'];
  locality?: Maybe<Scalars['String']>;
  name: Scalars['String'];
  place?: Maybe<Place>;
  region?: Maybe<Scalars['String']>;
  source?: Maybe<DataSource>;
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
  contact_AddTagById: Contact;
  contact_Create: Contact;
  contact_HardDelete: Result;
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
  organization_Update: Organization;
  phoneNumberDeleteFromContact: Result;
  phoneNumberDeleteFromContactById: Result;
  phoneNumberMergeToContact: PhoneNumber;
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

export type MutationContact_AddTagByIdArgs = {
  input?: InputMaybe<ContactTagInput>;
};

export type MutationContact_CreateArgs = {
  input: ContactInput;
};

export type MutationContact_HardDeleteArgs = {
  contactId: Scalars['ID'];
};

export type MutationContact_RemoveTagByIdArgs = {
  input?: InputMaybe<ContactTagInput>;
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

export type MutationOrganization_UpdateArgs = {
  input: OrganizationUpdateInput;
};

export type MutationPhoneNumberDeleteFromContactArgs = {
  contactId: Scalars['ID'];
  e164: Scalars['String'];
};

export type MutationPhoneNumberDeleteFromContactByIdArgs = {
  contactId: Scalars['ID'];
  id: Scalars['ID'];
};

export type MutationPhoneNumberMergeToContactArgs = {
  contactId: Scalars['ID'];
  input: PhoneNumberInput;
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
  domain?: Maybe<Scalars['String']>;
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
  /** Organization notes */
  notes: NotePage;
  organizationType?: Maybe<OrganizationType>;
  source: DataSource;
  sourceOfTruth: DataSource;
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

export type OrganizationInput = {
  appSource?: InputMaybe<Scalars['String']>;
  description?: InputMaybe<Scalars['String']>;
  domain?: InputMaybe<Scalars['String']>;
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
  id: Scalars['ID'];
  industry?: InputMaybe<Scalars['String']>;
  isPublic?: InputMaybe<Scalars['Boolean']>;
  name: Scalars['String'];
  organizationTypeId?: InputMaybe<Scalars['ID']>;
  website?: InputMaybe<Scalars['String']>;
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
  Ms = 'MS',
}

/**
 * Describes a phone number associated with a `Contact` in customerOS.
 * **A `return` object.**
 */
export type PhoneNumber = {
  __typename?: 'PhoneNumber';
  createdAt: Scalars['Time'];
  /**
   * The phone number in e164 format.
   * **Required**
   */
  e164: Scalars['String'];
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
  source: DataSource;
  updatedAt: Scalars['Time'];
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
  e164: Scalars['String'];
  /** Defines the type of phone number. */
  label?: InputMaybe<PhoneNumberLabel>;
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

/**
 * Describes a phone number associated with a `Contact` in customerOS.
 * **An `update` object.**
 */
export type PhoneNumberUpdateInput = {
  /**
   * The phone number in e164 format.
   * **Required**
   */
  e164: Scalars['String'];
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
   * - LABEL
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
  Desc = 'DESC',
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
