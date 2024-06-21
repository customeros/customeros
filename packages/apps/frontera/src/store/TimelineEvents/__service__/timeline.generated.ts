import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type TimelineQueryVariables = Types.Exact<{
  organizationId: Types.Scalars['ID']['input'];
  from: Types.Scalars['Time']['input'];
  size: Types.Scalars['Int']['input'];
}>;

export type TimelineQuery = {
  __typename?: 'Query';
  organization?: {
    __typename?: 'Organization';
    timelineEventsTotalCount: any;
    timelineEvents: Array<
      | {
          __typename: 'Action';
          id: string;
          actionType: Types.ActionType;
          appSource: string;
          createdAt: any;
          metadata?: string | null;
          content?: string | null;
          actionCreatedBy?: {
            __typename: 'User';
            id: string;
            name?: string | null;
            firstName: string;
            lastName: string;
            profilePhotoUrl?: string | null;
          } | null;
        }
      | { __typename: 'Analysis' }
      | {
          __typename: 'InteractionEvent';
          id: string;
          channel?: string | null;
          content?: string | null;
          contentType?: string | null;
          source: Types.DataSource;
          date: any;
          includes: Array<{
            __typename?: 'Attachment';
            id: string;
            mimeType: string;
            fileName: string;
            size: any;
          }>;
          issue?: {
            __typename?: 'Issue';
            id: string;
            externalLinks: Array<{
              __typename?: 'ExternalSystem';
              type: Types.ExternalSystemType;
              externalId?: string | null;
              externalUrl?: string | null;
            }>;
          } | null;
          externalLinks: Array<{
            __typename?: 'ExternalSystem';
            externalUrl?: string | null;
            type: Types.ExternalSystemType;
          }>;
          repliesTo?: { __typename?: 'InteractionEvent'; id: string } | null;
          summary?: {
            __typename?: 'Analysis';
            id: string;
            content?: string | null;
            contentType?: string | null;
          } | null;
          actionItems?: Array<{
            __typename?: 'ActionItem';
            id: string;
            content: string;
          }> | null;
          sentBy: Array<
            | {
                __typename: 'ContactParticipant';
                contactParticipant: {
                  __typename?: 'Contact';
                  id: string;
                  name?: string | null;
                  firstName?: string | null;
                  lastName?: string | null;
                  profilePhotoUrl?: string | null;
                };
              }
            | {
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
                    profilePhotoUrl?: string | null;
                  }>;
                  users: Array<{
                    __typename?: 'User';
                    id: string;
                    firstName: string;
                    lastName: string;
                    profilePhotoUrl?: string | null;
                  }>;
                  organizations: Array<{
                    __typename?: 'Organization';
                    id: string;
                    name: string;
                  }>;
                };
              }
            | {
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
              }
            | {
                __typename: 'OrganizationParticipant';
                organizationParticipant: {
                  __typename?: 'Organization';
                  id: string;
                  name: string;
                };
              }
            | { __typename: 'PhoneNumberParticipant' }
            | {
                __typename: 'UserParticipant';
                userParticipant: {
                  __typename?: 'User';
                  id: string;
                  name?: string | null;
                  firstName: string;
                  lastName: string;
                  profilePhotoUrl?: string | null;
                };
              }
          >;
          sentTo: Array<
            | {
                __typename: 'ContactParticipant';
                contactParticipant: {
                  __typename?: 'Contact';
                  id: string;
                  name?: string | null;
                  firstName?: string | null;
                  lastName?: string | null;
                  profilePhotoUrl?: string | null;
                };
              }
            | {
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
                    profilePhotoUrl?: string | null;
                  }>;
                  users: Array<{
                    __typename?: 'User';
                    id: string;
                    firstName: string;
                    lastName: string;
                    profilePhotoUrl?: string | null;
                  }>;
                  organizations: Array<{
                    __typename?: 'Organization';
                    id: string;
                    name: string;
                  }>;
                };
              }
            | {
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
              }
            | {
                __typename: 'OrganizationParticipant';
                organizationParticipant: {
                  __typename?: 'Organization';
                  id: string;
                  name: string;
                };
              }
            | { __typename: 'PhoneNumberParticipant' }
            | {
                __typename: 'UserParticipant';
                userParticipant: {
                  __typename?: 'User';
                  id: string;
                  name?: string | null;
                  firstName: string;
                  lastName: string;
                  profilePhotoUrl?: string | null;
                };
              }
          >;
          interactionSession?: {
            __typename?: 'InteractionSession';
            id: string;
            name: string;
            events: Array<{
              __typename?: 'InteractionEvent';
              id: string;
              channel?: string | null;
              date: any;
              sentBy: Array<
                | {
                    __typename: 'ContactParticipant';
                    contactParticipant: {
                      __typename?: 'Contact';
                      id: string;
                      name?: string | null;
                      firstName?: string | null;
                      lastName?: string | null;
                      profilePhotoUrl?: string | null;
                    };
                  }
                | {
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
                        profilePhotoUrl?: string | null;
                      }>;
                      users: Array<{
                        __typename?: 'User';
                        id: string;
                        firstName: string;
                        lastName: string;
                        profilePhotoUrl?: string | null;
                      }>;
                      organizations: Array<{
                        __typename?: 'Organization';
                        id: string;
                        name: string;
                      }>;
                    };
                  }
                | {
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
                  }
                | {
                    __typename: 'OrganizationParticipant';
                    organizationParticipant: {
                      __typename?: 'Organization';
                      id: string;
                      name: string;
                    };
                  }
                | { __typename?: 'PhoneNumberParticipant' }
                | {
                    __typename: 'UserParticipant';
                    userParticipant: {
                      __typename?: 'User';
                      id: string;
                      name?: string | null;
                      firstName: string;
                      lastName: string;
                      profilePhotoUrl?: string | null;
                    };
                  }
              >;
            }>;
            attendedBy: Array<
              | {
                  __typename: 'ContactParticipant';
                  contactParticipant: {
                    __typename?: 'Contact';
                    id: string;
                    name?: string | null;
                    firstName?: string | null;
                    lastName?: string | null;
                    profilePhotoUrl?: string | null;
                  };
                }
              | {
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
                      profilePhotoUrl?: string | null;
                    }>;
                    users: Array<{
                      __typename?: 'User';
                      id: string;
                      firstName: string;
                      lastName: string;
                      profilePhotoUrl?: string | null;
                    }>;
                    organizations: Array<{
                      __typename?: 'Organization';
                      id: string;
                      name: string;
                    }>;
                  };
                }
              | { __typename?: 'PhoneNumberParticipant' }
              | {
                  __typename: 'UserParticipant';
                  userParticipant: {
                    __typename?: 'User';
                    id: string;
                    name?: string | null;
                    firstName: string;
                    lastName: string;
                    profilePhotoUrl?: string | null;
                  };
                }
            >;
          } | null;
        }
      | { __typename: 'InteractionSession' }
      | {
          __typename: 'Issue';
          id: string;
          subject?: string | null;
          priority?: string | null;
          appSource: string;
          updatedAt: any;
          createdAt: any;
          description?: string | null;
          issueStatus: string;
          externalLinks: Array<{
            __typename?: 'ExternalSystem';
            type: Types.ExternalSystemType;
            externalId?: string | null;
            externalUrl?: string | null;
          }>;
          submittedBy?:
            | {
                __typename: 'ContactParticipant';
                contactParticipant: {
                  __typename?: 'Contact';
                  id: string;
                  name?: string | null;
                  firstName?: string | null;
                  lastName?: string | null;
                  profilePhotoUrl?: string | null;
                };
              }
            | {
                __typename: 'OrganizationParticipant';
                organizationParticipant: {
                  __typename?: 'Organization';
                  id: string;
                  name: string;
                };
              }
            | {
                __typename: 'UserParticipant';
                userParticipant: {
                  __typename?: 'User';
                  id: string;
                  name?: string | null;
                  firstName: string;
                  lastName: string;
                  profilePhotoUrl?: string | null;
                };
              }
            | null;
          followedBy: Array<
            | {
                __typename: 'ContactParticipant';
                contactParticipant: {
                  __typename?: 'Contact';
                  id: string;
                  name?: string | null;
                  firstName?: string | null;
                  lastName?: string | null;
                  profilePhotoUrl?: string | null;
                };
              }
            | {
                __typename: 'OrganizationParticipant';
                organizationParticipant: {
                  __typename?: 'Organization';
                  id: string;
                  name: string;
                };
              }
            | {
                __typename: 'UserParticipant';
                userParticipant: {
                  __typename?: 'User';
                  id: string;
                  name?: string | null;
                  firstName: string;
                  lastName: string;
                  profilePhotoUrl?: string | null;
                };
              }
          >;
          assignedTo: Array<
            | {
                __typename: 'ContactParticipant';
                contactParticipant: {
                  __typename?: 'Contact';
                  id: string;
                  name?: string | null;
                  firstName?: string | null;
                  lastName?: string | null;
                  profilePhotoUrl?: string | null;
                };
              }
            | {
                __typename: 'OrganizationParticipant';
                organizationParticipant: {
                  __typename?: 'Organization';
                  id: string;
                  name: string;
                };
              }
            | {
                __typename: 'UserParticipant';
                userParticipant: {
                  __typename?: 'User';
                  id: string;
                  name?: string | null;
                  firstName: string;
                  lastName: string;
                  profilePhotoUrl?: string | null;
                };
              }
          >;
          reportedBy?:
            | {
                __typename: 'ContactParticipant';
                contactParticipant: {
                  __typename?: 'Contact';
                  id: string;
                  name?: string | null;
                  firstName?: string | null;
                  lastName?: string | null;
                  profilePhotoUrl?: string | null;
                };
              }
            | {
                __typename: 'OrganizationParticipant';
                organizationParticipant: {
                  __typename?: 'Organization';
                  id: string;
                  name: string;
                };
              }
            | {
                __typename: 'UserParticipant';
                userParticipant: {
                  __typename?: 'User';
                  id: string;
                  name?: string | null;
                  firstName: string;
                  lastName: string;
                  profilePhotoUrl?: string | null;
                };
              }
            | null;
          interactionEvents: Array<{
            __typename?: 'InteractionEvent';
            content?: string | null;
            contentType?: string | null;
            createdAt: any;
            sentBy: Array<
              | {
                  __typename: 'ContactParticipant';
                  contactParticipant: {
                    __typename?: 'Contact';
                    id: string;
                    name?: string | null;
                    firstName?: string | null;
                    lastName?: string | null;
                    profilePhotoUrl?: string | null;
                  };
                }
              | {
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
                      profilePhotoUrl?: string | null;
                    }>;
                    users: Array<{
                      __typename?: 'User';
                      id: string;
                      firstName: string;
                      lastName: string;
                      profilePhotoUrl?: string | null;
                    }>;
                    organizations: Array<{
                      __typename?: 'Organization';
                      id: string;
                      name: string;
                    }>;
                  };
                }
              | {
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
                }
              | {
                  __typename: 'OrganizationParticipant';
                  organizationParticipant: {
                    __typename?: 'Organization';
                    id: string;
                    name: string;
                  };
                }
              | { __typename?: 'PhoneNumberParticipant' }
              | {
                  __typename: 'UserParticipant';
                  userParticipant: {
                    __typename?: 'User';
                    id: string;
                    name?: string | null;
                    firstName: string;
                    lastName: string;
                    profilePhotoUrl?: string | null;
                  };
                }
            >;
          }>;
          comments: Array<{
            __typename?: 'Comment';
            content?: string | null;
            contentType?: string | null;
            createdAt: any;
            createdBy?: {
              __typename?: 'User';
              id: string;
              name?: string | null;
              firstName: string;
              lastName: string;
            } | null;
          }>;
          issueTags?: Array<{
            __typename?: 'Tag';
            id: string;
            name: string;
          } | null> | null;
        }
      | {
          __typename: 'LogEntry';
          id: string;
          createdAt: any;
          updatedAt: any;
          source: Types.DataSource;
          content?: string | null;
          contentType?: string | null;
          logEntryStartedAt: any;
          logEntryCreatedBy?: {
            __typename: 'User';
            id: string;
            name?: string | null;
            firstName: string;
            lastName: string;
            profilePhotoUrl?: string | null;
            emails?: Array<{
              __typename?: 'Email';
              email?: string | null;
            }> | null;
          } | null;
          tags: Array<{ __typename?: 'Tag'; id: string; name: string }>;
          externalLinks: Array<{
            __typename?: 'ExternalSystem';
            type: Types.ExternalSystemType;
            externalUrl?: string | null;
            externalSource?: string | null;
          }>;
        }
      | {
          __typename: 'Meeting';
          id: string;
          name?: string | null;
          createdAt: any;
          updatedAt: any;
          startedAt?: any | null;
          endedAt?: any | null;
          agenda?: string | null;
          status: Types.MeetingStatus;
          attendedBy: Array<
            | {
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
              }
            | {
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
                  users: Array<{
                    __typename?: 'User';
                    firstName: string;
                    lastName: string;
                  }>;
                };
              }
            | {
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
              }
            | {
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
              }
          >;
          createdBy: Array<
            | {
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
              }
            | {
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
                  users: Array<{
                    __typename?: 'User';
                    firstName: string;
                    lastName: string;
                  }>;
                };
              }
            | {
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
              }
            | {
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
              }
          >;
          note: Array<{
            __typename?: 'Note';
            id: string;
            content?: string | null;
          }>;
        }
      | { __typename: 'Note' }
      | {
          __typename: 'Order';
          id: string;
          confirmedAt?: any | null;
          fulfilledAt?: any | null;
          createdAt: any;
          cancelledAt?: any | null;
        }
      | { __typename: 'PageView' }
    >;
  } | null;
};

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
      profilePhotoUrl?: string | null;
    }>;
    users: Array<{
      __typename?: 'User';
      id: string;
      firstName: string;
      lastName: string;
      profilePhotoUrl?: string | null;
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
    name?: string | null;
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
