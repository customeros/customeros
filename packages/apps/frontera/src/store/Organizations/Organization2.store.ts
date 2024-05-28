import type { RootStore } from '@store/root';

import { Channel } from 'phoenix';
import { P, match } from 'ts-pattern';
import { gql } from 'graphql-request';
import { Operation } from '@store/types';
import { makePayload } from '@store/util';
import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';
import { Store, makeAutoSyncable } from '@store/store';

import {
  Market,
  DataSource,
  ActionType,
  FundingRound,
  Organization,
  OnboardingStatus,
  OrganizationStage,
  LastTouchpointType,
  OrganizationUpdateInput,
  OrganizationRelationship,
  OpportunityRenewalLikelihood,
  OpportunityRenewalUpdateAllForOrganizationInput,
} from '@graphql/types';

export class OrganizationStore implements Store<Organization> {
  value: Organization = defaultValue;
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  load = makeAutoSyncable.load<Organization>();
  update = makeAutoSyncable.update<Organization>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, {
      channelName: 'Organization',
      mutator: this.save,
      getId: (d) => d?.metadata?.id,
      storeMapper: {
        contracts: {
          storeName: 'contracts',
          getItemId: (item) => item?.metadata?.id,
        },
      },
    });
    makeAutoObservable(this);
  }

  get id() {
    return this.value.metadata.id;
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
  private async save(operation: Operation) {
    const diff = operation.diff?.[0];
    const type = diff?.op;
    const path = diff?.path;
    const value = diff?.val;

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
      .otherwise(() => {
        const payload = makePayload<OrganizationUpdateInput>(operation);
        this.updateOrganization(payload);
      });
  }
}

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
    lastTouchPointTimelineEvent: {
      __typename: 'Action',
      id: crypto.randomUUID(),
      actionType: ActionType.Created,
      appSource: DataSource.Openline,
      createdAt: new Date().toISOString(),
      source: DataSource.Openline,
      createdBy: null,
    },
  }, // nested defaults ignored for now -> should be converted into a Store
  lastTouchPointTimelineEventId: '',
  leadSource: '',
  market: Market.B2B,
  notes: '',
  public: false,
  relationship: OrganizationRelationship.Stranger,
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
