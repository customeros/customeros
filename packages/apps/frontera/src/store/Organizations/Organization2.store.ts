import type { RootStore } from '@store/root';

<<<<<<< HEAD
import set from 'lodash/set';
=======
import { set } from 'lodash';
import merge from 'lodash/merge';
>>>>>>> db7d2c9c5 (add branches changes)
import { Channel } from 'phoenix';
import { P, match } from 'ts-pattern';
import { gql } from 'graphql-request';
import { Operation } from '@store/types';
import { makePayload } from '@store/util';
import { Transport } from '@store/transport';
import { rdiffResult } from 'recursive-diff';
import { Store, makeAutoSyncable } from '@store/store';
import { runInAction, makeAutoObservable } from 'mobx';
import { makeAutoSyncableGroup } from '@store/group-store';

import {
  Market,
  DataSource,
  SocialInput,
  FundingRound,
  Organization,
  OnboardingStatus,
  OrganizationStage,
  SocialUpdateInput,
  OrganizationInput,
  LastTouchpointType,
  LinkOrganizationsInput,
  OrganizationUpdateInput,
  OrganizationRelationship,
  OpportunityRenewalLikelihood,
  OpportunityRenewalUpdateAllForOrganizationInput,
} from '@graphql/types';
import { merge } from 'lodash';

export class OrganizationStore implements Store<Organization> {
  value: Organization = defaultValue;
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncable.load<Organization>();
  update = makeAutoSyncable.update<Organization>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncable(this, {
      channelName: 'Organization',
      mutator: this.save,
      getId: (d) => d?.metadata?.id,
    });
  }

  get id() {
    return this.value.metadata?.id;
  }
  set id(id: string) {
    this.value.metadata.id = id;
  }

  async invalidate() {
    try {
      this.isLoading = true;
      const { organization } = await this.transport.graphql.request<
        ORGANIZATION_QUERY_RESULT,
        { id: string }
      >(ORGANIZATIONS_QUERY, { id: this.id });

      this.load(organization);
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  private async updateOwner() {
    try {
      this.isLoading = true;
      await this.transport.graphql.request<unknown, UPDATE_OWNER_PAYLOAD>(
        UPDATE_OWNER_MUTATION,
        {
          organizationId: this.id,
          userId: this.value.owner?.id || '',
        },
      );
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }
  private async removeOwner() {
    try {
      this.isLoading = true;
      await this.transport.graphql.request<unknown, REMOVE_OWNER_PAYLOAD>(
        REMOVE_OWNER_MUTATION,
        {
          organizationId: this.id,
        },
      );
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }
  private async updateAllOpportunityRenewals() {
    try {
      this.isLoading = true;

      const amount =
        this.value.accountDetails?.renewalSummary?.arrForecast ?? 0;
      const potentialAmount =
        this.value.accountDetails?.renewalSummary?.maxArrForecast ?? 0;
      const rate = (amount / potentialAmount) * 100;

      await this.transport.graphql.request<
        unknown,
        UPDATE_ALL_OPPORTUNITY_RENEWALS_PAYLOAD
      >(UPDATE_ALL_OPPORTUNITY_RENEWAlS_MUTATION, {
        input: {
          organizationId: this.id,
          renewalAdjustedRate: rate,
          renewalLikelihood:
            this.value.accountDetails?.renewalSummary?.renewalLikelihood,
        },
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }
  private async updateOrganization(payload: OrganizationUpdateInput) {
    try {
      this.isLoading = true;
      await this.transport.graphql.request<
        unknown,
        UPDATE_ORGANIZATION_PAYLOAD
      >(UPDATE_ORGANIZATION_MUTATION, {
        input: {
          ...payload,
          id: this.id,
          patch: true,
        },
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  private async updateSocialMedia(index: number) {
    try {
      this.isLoading = true;
      await this.transport.graphql.request<
        unknown,
        UPDATE_SOCIAL_MEDIA_PAYLOAD
      >(UPDATE_SOCIAL_MEDIA_MUTATION, { input: this.value.socialMedia[index] });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  private async removeSocialMedia(socialId: string) {
    try {
      this.isLoading = true;
      await this.transport.graphql.request<
        unknown,
        REMOVE_SOCIAL_MEDIA_PAYLOAD
      >(REMOVE_SOCIAL_MEDIA_MUTATION, {
        socialId,
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  private async addSocialMedia(index: number) {
    try {
      this.isLoading = true;
      const { organization_AddSocial } = await this.transport.graphql.request<
        ADD_SOCIAL_MEDIA_RESPONSE,
        ADD_SOCIAL_MEDIA_PAYLOAD
      >(ADD_SOCIAL_MEDIA_MUTATION, {
        organizationId: this.id,
        input: {
          url: this.value.socialMedia[index].url,
        },
      });

      this.update(
        (org) => {
          org.socialMedia[index].id = organization_AddSocial.id;

          return org;
        },
        { mutate: false },
      );
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  private async addSubsidiary(organizationId: string) {
    try {
      this.isLoading = true;
      const { organization_AddSubsidiary } =
        await this.transport.graphql.request<
          { organization_AddSubsidiary: Organization },
          ADD_SUBSIDIARY_TO_ORGANIZATION
        >(ADD_SUBSIDIARY_TO_ORGANIZATION_MUTATION, {
          input: {
            organizationId: organizationId,
            subsidiaryId: this.id,
          },
        });

      this.update(
        (org) => {
          org.subsidiaries = organization_AddSubsidiary.subsidiaries;

          return org;
        },
        { mutate: false },
      );
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  private async createSubsidiary(subsidiaryId: string) {
    try {
      this.isLoading = true;
      const { organization_AddSubsidiary } =
        await this.transport.graphql.request<
          { organization_AddSubsidiary: Organization },
          ADD_SUBSIDIARY_TO_ORGANIZATION
        >(ADD_SUBSIDIARY_TO_ORGANIZATION_MUTATION, {
          input: {
            organizationId: this.id,
            subsidiaryId: subsidiaryId,
          },
        });
      this.update(
        (org) => {
          org.subsidiaries = organization_AddSubsidiary.subsidiaries;

          return org;
        },
        { mutate: false },
      );

      runInAction(() => {
        this.root.organizations.value.get(subsidiaryId)?.invalidate();
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        setTimeout(() => {}, 1500);

        this.isLoading = false;
      });
    }
  }

  private async removeSubsidiary(organizationId: string) {
    try {
      this.isLoading = true;
      const { organization_RemoveSubsidiary } =
        await this.transport.graphql.request<{
          organization_RemoveSubsidiary: Organization;
          REMOVE_SUBSIDIARY_FROM_ORGANIZATION: Organization;
        }>(REMOVE_SUBSIDIARY_FROM_ORGANIZATION_MUTATION, {
          organizationId: organizationId,
          subsidiaryId: this.id,
        });

      this.update(
        (org) => {
          org.subsidiaries = organization_RemoveSubsidiary.subsidiaries;

          return org;
        },
        { mutate: false },
      );
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  create = async (payload?: OrganizationInput) => {
    const newOrganization = new OrganizationStore(this.root, this.transport);
    const tempId = newOrganization.value.metadata.id;
    let serverId = '';
    if (payload) {
      merge(newOrganization.value, payload);
    }
    set(this.root.organizations.value, tempId, newOrganization);

    try {
      const { organization_Create } = await this.transport.graphql.request<
        CREATE_ORGANIZATION_RESPONSE,
        CREATE_ORGANIZATION_PAYLOAD
      >(CREATE_ORGANIZATION_MUTATION, {
        input: {
          name: 'Unnamed',
        },
      });
      runInAction(() => {
        serverId = organization_Create.id;

        newOrganization.value.metadata.id = serverId;
        set(this.root.organizations.value, serverId, newOrganization);
        tempId && this.root.organizations.value.delete(tempId);
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    } finally {
      if (serverId) {
        // Invalidate the cache after 1 second to allow the server to process the data
        // invalidating immediately would cause the server to return the organization data without
        // lastTouchpoint properties populated
        setTimeout(() => {
          newOrganization.invalidate();
          this.root.organizations.value.get(serverId)?.invalidate();
          this.root.organizations.value
            .get(serverId)
            ?.sync({ action: 'APPEND', ids: [serverId] });
        }, 1100);

        setTimeout(() => {
          this.createSubsidiary(serverId);
          this.root.organizations.value.get(this.id)?.invalidate();
          // window.location.href = `/organization/${serverId}?tab=about`;
        }, 1500);
      }
    }
  };

  private async save(operation: Operation) {
    const diff = operation.diff?.[0];
    const type = diff?.op;
    const path = diff?.path;
    const value = diff?.val;
    const oldValue = (diff as rdiffResult & { oldVal: unknown })?.oldVal;

    match(path)
      .with(['owner', ...P.array()], () => {
        if (type === 'update') {
          match(value)
            .with(null, () => {
              this.removeOwner();
            })
            .otherwise(() => {
              this.updateOwner();
            });
        }
      })
      .with(['accountDetails', 'renewalSummary', ...P.array()], () => {
        this.updateAllOpportunityRenewals();
      })

      .with(['socialMedia', ...P.array()], () => {
        const index = path[1];

        if (type === 'add') {
          this.addSocialMedia(index as number);
        }
        if (type === 'update') {
          this.updateSocialMedia(index as number);
        }
        if (type === 'delete') {
          this.removeSocialMedia(oldValue?.id);
        }
      })
      .with(['parentCompanies', 'subsidiaries', ...P.array()], () => {
        if (type === 'add') {
          this.addSubsidiary(value.organization.metadata?.id as string);
        }
        if (type === 'delete') {
          this.removeSubsidiary(oldValue?.organization?.metadata?.id);
        }
      })

      .otherwise(() => {
        const payload = makePayload<OrganizationUpdateInput>(operation);
        this.updateOrganization(payload);
      });
  }
  init(data: Organization) {
    const contracts = data.contracts?.map((item) => {
      this.root.contracts.load([item]);

      return this.root.contracts.value.get(item.metadata.id)?.value;
    });

    set(data, 'contracts', contracts);

    return data;
  }
}

type CREATE_ORGANIZATION_PAYLOAD = {
  input: OrganizationInput;
};
type CREATE_ORGANIZATION_RESPONSE = {
  organization_Create: {
    id: string;
    name: string;
  };
};
const CREATE_ORGANIZATION_MUTATION = gql`
  mutation createOrganization($input: OrganizationInput!) {
    organization_Create(input: $input) {
      id
      name
    }
  }
`;

type ORGANIZATION_QUERY_RESULT = {
  organization: Organization;
};
const ORGANIZATIONS_QUERY = gql`
  query Organization($id: ID!) {
    organization(id: $id) {
      name
      metadata {
        id
        created
      }
      parentCompanies {
        organization {
          metadata {
            id
          }
          name
        }
      }
      owner {
        id
        firstName
        lastName
        name
      }
      stage
      description
      industry
      website
      domains
      isCustomer
      logo
      icon
      relationship
      leadSource
      valueProposition
      socialMedia {
        id
        url
      }
      employees
      yearFounded
      accountDetails {
        renewalSummary {
          arrForecast
          maxArrForecast
          renewalLikelihood
          nextRenewalDate
        }
        onboarding {
          status
          comments
          updatedAt
        }
      }
      locations {
        id
        name
        country
        region
        locality
        zip
        street
        postalCode
        houseNumber
        rawAddress
      }
      subsidiaries {
        organization {
          metadata {
            id
          }
        }
      }
      parentCompanies {
        organization {
          metadata {
            id
          }
        }
      }
      lastTouchpoint {
        lastTouchPointTimelineEventId
        lastTouchPointAt
        lastTouchPointType
        lastTouchPointTimelineEvent {
          __typename
          ... on PageView {
            id
          }
          ... on Issue {
            id
            createdAt
            updatedAt
          }
          ... on LogEntry {
            id
            createdBy {
              lastName
              firstName
            }
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
            actionType
            createdBy {
              id
              firstName
              lastName
            }
          }
        }
      }
    }
  }
`;
type UPDATE_OWNER_PAYLOAD = {
  userId: string;
  organizationId: string;
};
const UPDATE_OWNER_MUTATION = gql`
  mutation setOrganizationOwner($organizationId: ID!, $userId: ID!) {
    organization_SetOwner(organizationId: $organizationId, userId: $userId) {
      id
    }
  }
`;
type REMOVE_OWNER_PAYLOAD = {
  organizationId: string;
};
const REMOVE_OWNER_MUTATION = gql`
  mutation setOrganizationOwner($organizationId: ID!) {
    organization_UnsetOwner(organizationId: $organizationId) {
      id
    }
  }
`;
type UPDATE_ALL_OPPORTUNITY_RENEWALS_PAYLOAD = {
  input: OpportunityRenewalUpdateAllForOrganizationInput;
};
const UPDATE_ALL_OPPORTUNITY_RENEWAlS_MUTATION = gql`
  mutation bulkUpdateOpportunityRenewal(
    $input: OpportunityRenewalUpdateAllForOrganizationInput!
  ) {
    opportunityRenewal_UpdateAllForOrganization(input: $input) {
      metadata {
        id
      }
    }
  }
`;
type UPDATE_ORGANIZATION_PAYLOAD = {
  input: OrganizationUpdateInput;
};
const UPDATE_ORGANIZATION_MUTATION = gql`
  mutation updateOrganization($input: OrganizationUpdateInput!) {
    organization_Update(input: $input) {
      metadata {
        id
      }
    }
  }
`;

type UPDATE_SOCIAL_MEDIA_PAYLOAD = {
  input: SocialUpdateInput;
};

const UPDATE_SOCIAL_MEDIA_MUTATION = gql`
  mutation updateSocial($input: SocialUpdateInput!) {
    social_Update(input: $input) {
      id
      url
    }
  }
`;

type REMOVE_SOCIAL_MEDIA_PAYLOAD = {
  socialId: string;
};

const REMOVE_SOCIAL_MEDIA_MUTATION = gql`
  mutation removeSocial($socialId: ID!) {
    social_Remove(socialId: $socialId) {
      result
    }
  }
`;

type ADD_SOCIAL_MEDIA_PAYLOAD = {
  input: SocialInput;
  organizationId: string;
};

type ADD_SOCIAL_MEDIA_RESPONSE = {
  organization_AddSocial: {
    id: string;
    url: string;
  };
};

const ADD_SOCIAL_MEDIA_MUTATION = gql`
  mutation addSocial($organizationId: ID!, $input: SocialInput!) {
    organization_AddSocial(organizationId: $organizationId, input: $input) {
      id
      url
    }
  }
`;

type ADD_SUBSIDIARY_TO_ORGANIZATION = {
  input: LinkOrganizationsInput;
};

const ADD_SUBSIDIARY_TO_ORGANIZATION_MUTATION = gql`
  mutation addSubsidiaryToOrganization($input: LinkOrganizationsInput!) {
    organization_AddSubsidiary(input: $input) {
      metadata {
        id
      }
      subsidiaries {
        organization {
          metadata {
            id
          }
          name
          locations {
            id
            address
          }
        }
      }
    }
  }
`;

const REMOVE_SUBSIDIARY_FROM_ORGANIZATION_MUTATION = gql`
  mutation removeSubsidiaryToOrganization(
    $organizationId: ID!
    $subsidiaryId: ID!
  ) {
    organization_RemoveSubsidiary(
      organizationId: $organizationId
      subsidiaryId: $subsidiaryId
    ) {
      id
      subsidiaries {
        organization {
          id
          name
          locations {
            id
            address
          }
        }
      }
    }
  }
`;

const defaultValue: Organization = {
  name: 'Unnamed',
  metadata: {
    id: crypto.randomUUID(),
    created: new Date().toISOString(),
    lastUpdated: new Date().toISOString(),
    appSource: DataSource.Openline,
    source: DataSource.Openline,
    sourceOfTruth: DataSource.Openline,
  },
  owner: null,
  contactCount: 0,
  contacts: {
    content: [],
    totalElements: 0,
    totalPages: 0,
  },
  customerOsId: '',
  customFields: [],
  domains: [],
  emails: [],
  externalLinks: [],
  industry: '',
  fieldSets: [],
  hide: false,
  inboundCommsCount: 0,
  issueSummaryByStatus: [],
  jobRoles: [],
  locations: [],
  orders: [],
  outboundCommsCount: 0,
  phoneNumbers: [],
  parentCompanies: [],
  socialMedia: [],
  stage: OrganizationStage.Target,
  tags: [],
  subsidiaries: [],
  suggestedMergeTo: [],
  timelineEvents: [],
  website: '',
  timelineEventsTotalCount: 0,
  accountDetails: {
    onboarding: {
      status: OnboardingStatus.NotStarted,
      comments: '',
      updatedAt: '',
    },
    renewalSummary: {
      arrForecast: 0,
      maxArrForecast: 0,
      renewalLikelihood: OpportunityRenewalLikelihood.HighRenewal,
      nextRenewalDate: '',
    },
  },
  contracts: [],
  customId: '',
  description: '',
  employees: 0,
  employeeGrowthRate: '',
  // entityTemplate: {} -> ignored | unused
  headquarters: '',
  isCustomer: false,
  logo: '',
  industryGroup: '',
  lastFundingAmount: '',
  lastFundingRound: FundingRound.PreSeed,
  lastTouchpoint: {
    lastTouchPointTimelineEventId: crypto.randomUUID(),
    lastTouchPointAt: new Date().toISOString(),
    lastTouchPointType: LastTouchpointType.ActionCreated,
  }, // nested defaults ignored for now -> should be converted into a Store
  lastTouchPointTimelineEventId: '',
  leadSource: '',
  market: Market.B2B,
  notes: '',
  public: false,
  relationship: OrganizationRelationship.Prospect,
  slackChannelId: '',
  stageLastUpdated: '',
  subIndustry: '',
  targetAudience: '',
  valueProposition: '',
  yearFounded: 0,
  // deprecated field -> needed because they're required in the TS type
  id: '',
  appSource: '',
  source: DataSource.Na,
  socials: [],
  createdAt: '',
  sourceOfTruth: DataSource.Na,
  subsidiaryOf: [],
  updatedAt: '',
};
