// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../../src/types/__generated__/graphql.types';

export type InteractionEventParticipantFragmentContactParticipantFragment = {
  __typename: 'ContactParticipant';
  contactParticipant: {
    __typename?: 'Contact';
    id: string;
    name?: string | null;
    firstName?: string | null;
    lastName?: string | null;
    profilePhotoUrl?: string | null;
  };
};

export type InteractionEventParticipantFragmentEmailParticipantFragment = {
  __typename: 'EmailParticipant';
  type?: string | null;
  emailParticipant: {
    __typename?: 'Email';
    email?: string | null;
    id: string;
    contacts: Array<{
      __typename?: 'Contact';
      id: string;
      name?: string | null;
      firstName?: string | null;
      lastName?: string | null;
    }>;
    users: Array<{
      __typename?: 'User';
      id: string;
      firstName: string;
      lastName: string;
    }>;
    organizations: Array<{
      __typename?: 'Organization';
      id: string;
      name: string;
    }>;
  };
};

export type InteractionEventParticipantFragmentJobRoleParticipantFragment = {
  __typename: 'JobRoleParticipant';
  jobRoleParticipant: {
    __typename?: 'JobRole';
    id: string;
    contact?: {
      __typename?: 'Contact';
      id: string;
      name?: string | null;
      firstName?: string | null;
      lastName?: string | null;
      profilePhotoUrl?: string | null;
    } | null;
  };
};

export type InteractionEventParticipantFragmentOrganizationParticipantFragment =
  {
    __typename: 'OrganizationParticipant';
    organizationParticipant: {
      __typename?: 'Organization';
      id: string;
      name: string;
    };
  };

export type InteractionEventParticipantFragmentPhoneNumberParticipantFragment =
  { __typename?: 'PhoneNumberParticipant' };

export type InteractionEventParticipantFragmentUserParticipantFragment = {
  __typename: 'UserParticipant';
  userParticipant: {
    __typename?: 'User';
    id: string;
    firstName: string;
    lastName: string;
    profilePhotoUrl?: string | null;
  };
};

export type InteractionEventParticipantFragmentFragment =
  | InteractionEventParticipantFragmentContactParticipantFragment
  | InteractionEventParticipantFragmentEmailParticipantFragment
  | InteractionEventParticipantFragmentJobRoleParticipantFragment
  | InteractionEventParticipantFragmentOrganizationParticipantFragment
  | InteractionEventParticipantFragmentPhoneNumberParticipantFragment
  | InteractionEventParticipantFragmentUserParticipantFragment;

export type MeetingParticipantFragmentContactParticipantFragment = {
  __typename: 'ContactParticipant';
  contactParticipant: {
    __typename?: 'Contact';
    id: string;
    name?: string | null;
    firstName?: string | null;
    lastName?: string | null;
    profilePhotoUrl?: string | null;
    timezone?: string | null;
    emails: Array<{
      __typename?: 'Email';
      id: string;
      email?: string | null;
      rawEmail?: string | null;
      primary: boolean;
    }>;
  };
};

export type MeetingParticipantFragmentEmailParticipantFragment = {
  __typename: 'EmailParticipant';
  emailParticipant: {
    __typename?: 'Email';
    rawEmail?: string | null;
    email?: string | null;
    contacts: Array<{
      __typename?: 'Contact';
      firstName?: string | null;
      lastName?: string | null;
      name?: string | null;
      timezone?: string | null;
    }>;
    users: Array<{ __typename?: 'User'; firstName: string; lastName: string }>;
  };
};

export type MeetingParticipantFragmentOrganizationParticipantFragment = {
  __typename: 'OrganizationParticipant';
  organizationParticipant: {
    __typename?: 'Organization';
    id: string;
    name: string;
    emails: Array<{
      __typename?: 'Email';
      id: string;
      email?: string | null;
      rawEmail?: string | null;
      primary: boolean;
    }>;
  };
};

export type MeetingParticipantFragmentUserParticipantFragment = {
  __typename: 'UserParticipant';
  userParticipant: {
    __typename?: 'User';
    id: string;
    firstName: string;
    lastName: string;
    profilePhotoUrl?: string | null;
    emails?: Array<{
      __typename?: 'Email';
      id: string;
      email?: string | null;
      rawEmail?: string | null;
      primary: boolean;
    }> | null;
  };
};

export type MeetingParticipantFragmentFragment =
  | MeetingParticipantFragmentContactParticipantFragment
  | MeetingParticipantFragmentEmailParticipantFragment
  | MeetingParticipantFragmentOrganizationParticipantFragment
  | MeetingParticipantFragmentUserParticipantFragment;

export const InteractionEventParticipantFragmentFragmentDoc = `
    fragment InteractionEventParticipantFragment on InteractionEventParticipant {
  ... on EmailParticipant {
    __typename
    type
    emailParticipant {
      email
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
      id
      name
      firstName
      lastName
      profilePhotoUrl
    }
  }
  ... on JobRoleParticipant {
    __typename
    jobRoleParticipant {
      id
      contact {
        id
        name
        firstName
        lastName
        profilePhotoUrl
      }
    }
  }
  ... on UserParticipant {
    __typename
    userParticipant {
      id
      firstName
      lastName
      profilePhotoUrl
    }
  }
  ... on OrganizationParticipant {
    __typename
    organizationParticipant {
      id
      name
    }
  }
}
    `;
export const MeetingParticipantFragmentFragmentDoc = `
    fragment MeetingParticipantFragment on MeetingParticipant {
  ... on ContactParticipant {
    __typename
    contactParticipant {
      id
      name
      firstName
      lastName
      profilePhotoUrl
      timezone
      emails {
        id
        email
        rawEmail
        primary
      }
    }
  }
  ... on UserParticipant {
    __typename
    userParticipant {
      id
      firstName
      lastName
      profilePhotoUrl
      emails {
        id
        email
        rawEmail
        primary
      }
    }
  }
  ... on OrganizationParticipant {
    __typename
    organizationParticipant {
      id
      name
      emails {
        id
        email
        rawEmail
        primary
      }
    }
  }
  ... on EmailParticipant {
    __typename
    emailParticipant {
      rawEmail
      email
      contacts {
        firstName
        lastName
        name
        timezone
      }
      users {
        firstName
        lastName
      }
    }
  }
}
    `;
