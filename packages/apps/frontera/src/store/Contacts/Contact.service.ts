import type { Transport } from '@store/transport';

import { gql } from 'graphql-request';

import type {
  EmailInput,
  SocialInput,
  JobRoleInput,
  ContactTagInput,
  PhoneNumberInput,
  SocialUpdateInput,
  ContactUpdateInput,
  JobRoleUpdateInput,
  PhoneNumberUpdateInput,
  EmailUpdateAddressInput,
  ContactOrganizationInput,
} from '@graphql/types';

class ContactService {
  private static instance: ContactService | null = null;
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport): ContactService {
    if (!ContactService.instance) {
      ContactService.instance = new ContactService(transport);
    }

    return ContactService.instance;
  }

  async linkOrganization(
    payload: LINK_ORGANIZATION_PAYLOAD,
  ): Promise<LINK_ORGANIZATION_RESPONSE> {
    return this.transport.graphql.request<
      LINK_ORGANIZATION_RESPONSE,
      LINK_ORGANIZATION_PAYLOAD
    >(LINK_ORGANIZATION_MUTATION, payload);
  }

  async updateContact(
    payload: UPDATE_CONTACT_PAYLOAD,
  ): Promise<UPDATE_CONTACT_RESPONSE> {
    return this.transport.graphql.request<
      UPDATE_CONTACT_RESPONSE,
      UPDATE_CONTACT_PAYLOAD
    >(UPDATE_CONTACT_MUTATION, payload);
  }

  async addJobRole(
    payload: ADD_JOB_ROLE_PAYLOAD,
  ): Promise<ADD_JOB_ROLE_RESPONSE> {
    return this.transport.graphql.request<
      ADD_JOB_ROLE_RESPONSE,
      ADD_JOB_ROLE_PAYLOAD
    >(ADD_JOB_ROLE_MUTATION, payload);
  }

  async updateJobRole(
    payload: UPDATE_JOB_ROLE_PAYLOAD,
  ): Promise<UPDATE_JOB_ROLE_RESPONSE> {
    return this.transport.graphql.request<
      UPDATE_JOB_ROLE_RESPONSE,
      UPDATE_JOB_ROLE_PAYLOAD
    >(UPDATE_JOB_ROLE_MUTATION, payload);
  }

  async addContactEmail(
    payload: ADD_CONTACT_EMAIL_PAYLOAD,
  ): Promise<ADD_CONTACT_EMAIL_RESPONSE> {
    return this.transport.graphql.request<
      ADD_CONTACT_EMAIL_RESPONSE,
      ADD_CONTACT_EMAIL_PAYLOAD
    >(ADD_CONTACT_EMAIL_MUTATION, payload);
  }

  async updateContactEmail(
    payload: UPDATE_CONTACT_EMAIL_PAYLOAD,
  ): Promise<UPDATE_CONTACT_EMAIL_RESPONSE> {
    return this.transport.graphql.request<
      UPDATE_CONTACT_EMAIL_RESPONSE,
      UPDATE_CONTACT_EMAIL_PAYLOAD
    >(UPDATE_CONTACT_EMAIL_MUTATION, payload);
  }

  async removeContactEmail(
    payload: REMOVE_CONTACT_EMAIL_PAYLOAD,
  ): Promise<REMOVE_CONTACT_EMAIL_RESPONSE> {
    return this.transport.graphql.request<
      REMOVE_CONTACT_EMAIL_RESPONSE,
      REMOVE_CONTACT_EMAIL_PAYLOAD
    >(REMOVE_CONTACT_EMAIL_MUTATION, payload);
  }

  async addPhoneNumber(
    payload: ADD_PHONE_NUMBER_PAYLOAD,
  ): Promise<ADD_PHONE_NUMBER_RESPONSE> {
    return this.transport.graphql.request<
      ADD_PHONE_NUMBER_RESPONSE,
      ADD_PHONE_NUMBER_PAYLOAD
    >(ADD_PHONE_NUMBER_MUTATION, payload);
  }

  async updatePhoneNumber(
    payload: UPDATE_PHONE_NUMBER_PAYLOAD,
  ): Promise<UPDATE_PHONE_NUMBER_RESPONSE> {
    return this.transport.graphql.request<
      UPDATE_PHONE_NUMBER_RESPONSE,
      UPDATE_PHONE_NUMBER_PAYLOAD
    >(UPDATE_PHONE_NUMBER_MUTATION, payload);
  }

  async removePhoneNumber(
    payload: REMOVE_PHONE_NUMBER_PAYLOAD,
  ): Promise<REMOVE_PHONE_NUMBER_RESPONSE> {
    return this.transport.graphql.request<
      REMOVE_PHONE_NUMBER_RESPONSE,
      REMOVE_PHONE_NUMBER_PAYLOAD
    >(REMOVE_PHONE_NUMBER_MUTATION, payload);
  }

  async addSocial(payload: ADD_SOCIAL_PAYLOAD): Promise<ADD_SOCIAL_RESPONSE> {
    return this.transport.graphql.request<
      ADD_SOCIAL_RESPONSE,
      ADD_SOCIAL_PAYLOAD
    >(ADD_SOCIAL_MUTATION, payload);
  }

  async updateSocial(
    payload: UPDATE_SOCIAL_PAYLOAD,
  ): Promise<UPDATE_SOCIAL_RESPONSE> {
    return this.transport.graphql.request<
      UPDATE_SOCIAL_RESPONSE,
      UPDATE_SOCIAL_PAYLOAD
    >(UPDATE_SOCIAL_MUTATION, payload);
  }

  async findEmail(payload: FIND_EMAIL_PAYLOAD): Promise<FIND_EMAIL_RESPONSE> {
    return this.transport.graphql.request<
      FIND_EMAIL_RESPONSE,
      FIND_EMAIL_PAYLOAD
    >(FIND_EMAIL_MUTATION, payload);
  }

  async deleteContact(
    payload: DELETE_CONTACT_PAYLOAD,
  ): Promise<DELETE_CONTACT_RESPONSE> {
    return this.transport.graphql.request<
      DELETE_CONTACT_RESPONSE,
      DELETE_CONTACT_PAYLOAD
    >(DELETE_CONTACT_MUTATION, payload);
  }

  async addTagsToContact(payload: ADD_TAGS_TO_CONTACT_PAYLOAD): Promise<void> {
    return this.transport.graphql.request<void, ADD_TAGS_TO_CONTACT_PAYLOAD>(
      ADD_TAGS_TO_CONTACT_MUTATION,
      payload,
    );
  }

  async removeTagsFromContact(
    payload: REMOVE_TAGS_FROM_CONTACT_PAYLOAD,
  ): Promise<void> {
    return this.transport.graphql.request<
      void,
      REMOVE_TAGS_FROM_CONTACT_PAYLOAD
    >(REMOVE_TAGS_FROM_CONTACT_MUTATION, payload);
  }
}

type LINK_ORGANIZATION_RESPONSE = {
  contact_AddOrganizationById: {
    id: string;
  };
};
type LINK_ORGANIZATION_PAYLOAD = {
  input: ContactOrganizationInput;
};
const LINK_ORGANIZATION_MUTATION = gql`
  mutation linkOrganization($input: ContactOrganizationInput!) {
    contact_AddOrganizationById(input: $input) {
      id
    }
  }
`;

type UPDATE_CONTACT_RESPONSE = {
  contact_Update: {
    id: string;
  };
};
type UPDATE_CONTACT_PAYLOAD = {
  input: ContactUpdateInput;
};
const UPDATE_CONTACT_MUTATION = gql`
  mutation updateContact($input: ContactUpdateInput!) {
    contact_Update(input: $input) {
      id
    }
  }
`;

type ADD_JOB_ROLE_RESPONSE = {
  jobRole_Create: {
    id: string;
  };
};
type ADD_JOB_ROLE_PAYLOAD = {
  contactId: string;
  input: JobRoleInput;
};
const ADD_JOB_ROLE_MUTATION = gql`
  mutation addContactSocial($contactId: ID!, $input: JobRoleInput!) {
    jobRole_Create(contactId: $contactId, input: $input) {
      id
    }
  }
`;

type UPDATE_JOB_ROLE_RESPONSE = {
  jobRole_Update: {
    id: string;
  };
};
type UPDATE_JOB_ROLE_PAYLOAD = {
  contactId: string;
  input: JobRoleUpdateInput;
};
const UPDATE_JOB_ROLE_MUTATION = gql`
  mutation updateContactRole($contactId: ID!, $input: JobRoleUpdateInput!) {
    jobRole_Update(contactId: $contactId, input: $input) {
      id
    }
  }
`;

type ADD_CONTACT_EMAIL_RESPONSE = {
  emailMergeToContact: {
    id: string;
  };
};
type ADD_CONTACT_EMAIL_PAYLOAD = {
  contactId: string;
  input: EmailInput;
};
const ADD_CONTACT_EMAIL_MUTATION = gql`
  mutation addContactEmail($contactId: ID!, $input: EmailInput!) {
    emailMergeToContact(contactId: $contactId, input: $input) {
      id
    }
  }
`;

type UPDATE_CONTACT_EMAIL_RESPONSE = {
  emailUpdate: {
    id: string;
  };
};
type UPDATE_CONTACT_EMAIL_PAYLOAD = {
  input: EmailUpdateAddressInput;
};
const UPDATE_CONTACT_EMAIL_MUTATION = gql`
  mutation updateContactEmail($input: EmailUpdateAddressInput!) {
    emailUpdate(input: $input) {
      id
    }
  }
`;

type REMOVE_CONTACT_EMAIL_RESPONSE = {
  emailRemoveFromContact: {
    result: boolean;
  };
};
type REMOVE_CONTACT_EMAIL_PAYLOAD = {
  email: string;
  contactId: string;
};
const REMOVE_CONTACT_EMAIL_MUTATION = gql`
  mutation removeContactEmail($contactId: ID!, $email: String!) {
    emailRemoveFromContact(contactId: $contactId, email: $email) {
      result
    }
  }
`;

type ADD_PHONE_NUMBER_RESPONSE = {
  phoneNumberMergeToContact: {
    id: string;
    rawPhoneNumber: string;
  };
};
type ADD_PHONE_NUMBER_PAYLOAD = {
  contactId: string;
  input: PhoneNumberInput;
};
const ADD_PHONE_NUMBER_MUTATION = gql`
  mutation addContactPhoneNumber($contactId: ID!, $input: PhoneNumberInput!) {
    phoneNumberMergeToContact(contactId: $contactId, input: $input) {
      id
      rawPhoneNumber
    }
  }
`;

type UPDATE_PHONE_NUMBER_RESPONSE = {
  phoneNumberUpdateInContact: {
    id: string;
  };
};
type UPDATE_PHONE_NUMBER_PAYLOAD = {
  input: PhoneNumberUpdateInput;
};
const UPDATE_PHONE_NUMBER_MUTATION = gql`
  mutation updateContactPhoneNumber($input: PhoneNumberUpdateInput!) {
    phoneNumber_Update(input: $input) {
      id
    }
  }
`;

type REMOVE_PHONE_NUMBER_RESPONSE = {
  phoneNumberRemoveFromContactById: {
    result: boolean;
  };
};
type REMOVE_PHONE_NUMBER_PAYLOAD = {
  id: string;
  contactId: string;
};
const REMOVE_PHONE_NUMBER_MUTATION = gql`
  mutation removeContactPhoneNumber($contactId: ID!, $id: ID!) {
    phoneNumberRemoveFromContactById(contactId: $contactId, id: $id) {
      result
    }
  }
`;

type ADD_SOCIAL_RESPONSE = {
  contact_AddSocial: {
    id: string;
  };
};
type ADD_SOCIAL_PAYLOAD = {
  contactId: string;
  input: SocialInput;
};
const ADD_SOCIAL_MUTATION = gql`
  mutation addContactSocial($contactId: ID!, $input: SocialInput!) {
    contact_AddSocial(contactId: $contactId, input: $input) {
      id
    }
  }
`;

type UPDATE_SOCIAL_RESPONSE = {
  social_Update: {
    id: string;
  };
};
type UPDATE_SOCIAL_PAYLOAD = {
  input: SocialUpdateInput;
};
const UPDATE_SOCIAL_MUTATION = gql`
  mutation updateContactSocial($input: SocialUpdateInput!) {
    social_Update(input: $input) {
      id
    }
  }
`;

type FIND_EMAIL_RESPONSE = void;
type FIND_EMAIL_PAYLOAD = {
  contactId: string;
  organizationId: string;
};
const FIND_EMAIL_MUTATION = gql`
  mutation findContactEmail($contactId: ID!, $organizationId: ID!) {
    contact_FindEmail(contactId: $contactId, organizationId: $organizationId)
  }
`;

type DELETE_CONTACT_PAYLOAD = {
  contactId: string;
};
type DELETE_CONTACT_RESPONSE = {
  contact_HardDelete: {
    result: boolean;
  };
};
const DELETE_CONTACT_MUTATION = gql`
  mutation deleteContact($contactId: ID!) {
    contact_HardDelete(contactId: $contactId) {
      result
    }
  }
`;

type ADD_TAGS_TO_CONTACT_PAYLOAD = {
  input: ContactTagInput;
};

const ADD_TAGS_TO_CONTACT_MUTATION = gql`
  mutation addTagsToContact($input: ContactTagInput!) {
    contact_AddTag(input: $input) {
      accepted
    }
  }
`;

type REMOVE_TAGS_FROM_CONTACT_PAYLOAD = {
  input: ContactTagInput;
};

const REMOVE_TAGS_FROM_CONTACT_MUTATION = gql`
  mutation removeTagFromContact($input: ContactTagInput!) {
    contact_RemoveTag(input: $input) {
      accepted
    }
  }
`;

export { ContactService };
