import type { RootStore } from '@store/root';

import { P, match } from 'ts-pattern';
import { gql } from 'graphql-request';
import { Operation } from '@store/types';
import { makePayload } from '@store/util';
import { Syncable } from '@store/syncable';
import { Transport } from '@store/transport';
import { rdiffResult } from 'recursive-diff';
import { countryMap } from '@assets/countries/countriesMap';
import { ActionStore } from '@store/TimelineEvents/Actions/Action.store';
import { action, override, computed, runInAction, makeObservable } from 'mobx';

import {
  Tag,
  Market,
  Contact,
  Contract,
  DataSource,
  SocialInput,
  FundingRound,
  Organization,
  OnboardingStatus,
  OrganizationStage,
  SocialUpdateInput,
  LastTouchpointType,
  OrganizationTagInput,
  LinkOrganizationsInput,
  OrganizationUpdateInput,
  OrganizationRelationship,
  OpportunityRenewalLikelihood,
  OpportunityRenewalUpdateAllForOrganizationInput,
} from '@graphql/types';

import { OrganizationsService } from './__service__/Organizations.service';

export class OrganizationStore extends Syncable<Organization> {
  private service: OrganizationsService;

  constructor(
    public root: RootStore,
    public transport: Transport,
    data: Organization,
  ) {
    super(root, transport, data ?? getDefaultValue());
    this.service = OrganizationsService.getInstance(transport);

    makeObservable<
      OrganizationStore,
      | 'updateOwner'
      | 'removeOwner'
      | 'addSubsidiary'
      | 'addSocialMedia'
      | 'removeSubsidiary'
      | 'updateSocialMedia'
      | 'removeSocialMedia'
      | 'updateOrganization'
      | 'addTagsToOrganization'
      | 'updateOnboardingStatus'
      | 'removeTagsFromOrganization'
      | 'updateAllOpportunityRenewals'
    >(this, {
      id: override,
      save: override,
      getId: override,
      setId: override,
      invalidate: action,
      contacts: computed,
      contracts: computed,
      updateOwner: action,
      removeOwner: action,
      addSubsidiary: action,
      addSocialMedia: action,
      subsidiaries: computed,
      country: computed,
      getChannelName: override,
      removeSubsidiary: action,
      parentCompanies: computed,
      invoices: computed,
      removeSocialMedia: action,
      updateSocialMedia: action,
      updateOrganization: action,
      addTagsToOrganization: action,
      updateOnboardingStatus: action,
      removeTagsFromOrganization: action,
      updateAllOpportunityRenewals: action,
    });
  }

  get contacts() {
    const contactIds = this.value.contacts.content.map(
      ({ metadata }) => metadata?.id,
    );

    const result: Contact[] = [];

    contactIds.forEach((id) => {
      const contactStore = this.root.contacts.value.get(id);
      if (contactStore) {
        result.push(contactStore.value);
      }
    });

    return result;
  }

  get contracts() {
    const contractIds = this.value.contracts?.map(
      ({ metadata }) => metadata.id,
    );

    const result: Contract[] = [];

    contractIds?.forEach((id) => {
      const contractStore = this.root.contracts.value.get(id);
      if (contractStore) {
        result.push(contractStore.value);
      }
    });

    return result;
  }

  get invoices() {
    return this.root.invoices
      .toArray()
      .filter(
        (invoice) =>
          invoice?.value?.organization?.metadata?.id === this.id &&
          !invoice?.value?.dryRun,
      );
  }

  get country() {
    if (!this.value.locations?.[0]?.countryCodeA2) return undefined;

    return countryMap.get(this.value.locations[0].countryCodeA2.toLowerCase());
  }

  get parentCompanies() {
    const parentCompanyIds = this.value.parentCompanies.map(
      (o) => o.organization.metadata.id,
    );

    const result: Organization[] = [];

    parentCompanyIds.forEach((id) => {
      const organizationStore = this.root.organizations.value.get(id);
      if (organizationStore) {
        result.push(organizationStore.value);
      }
    });

    return result;
  }

  get subsidiaries() {
    const subsidiaryIds = this.value.subsidiaries.map(
      (o) => o.organization.metadata.id,
    );

    const result: Organization[] = [];

    subsidiaryIds.forEach((id) => {
      const organizationStore = this.root.organizations.value.get(id);
      if (organizationStore) {
        result.push(organizationStore.value);
      }
    });

    return result;
  }

  get owner() {
    return this.root.users.value.get(this.value.owner?.id as string);
  }

  get id() {
    return this.value.metadata?.id;
  }

  getId() {
    return this.value.metadata?.id;
  }

  setId(id: string) {
    this.value.metadata.id = id;
  }

  getChannelName() {
    return `Organization:${this.id}`;
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
    const { id, url } = this.value.socialMedia[index];
    try {
      this.isLoading = true;
      await this.transport.graphql.request<
        unknown,
        UPDATE_SOCIAL_MEDIA_PAYLOAD
      >(UPDATE_SOCIAL_MEDIA_MUTATION, {
        input: {
          id,
          url,
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

  private async addSubsidiary(subsidiaryId: string) {
    try {
      this.isLoading = true;
      await this.transport.graphql.request<
        { organization_AddSubsidiary: Organization },
        ADD_SUBSIDIARY_TO_ORGANIZATION
      >(ADD_SUBSIDIARY_TO_ORGANIZATION_MUTATION, {
        input: {
          organizationId: this.id,
          subsidiaryId: subsidiaryId,
        },
      });

      runInAction(() => {
        this.root.organizations.value.get(subsidiaryId)?.update(
          (org: Organization) => {
            org.parentCompanies.push({
              organization: this.value,
            });

            return org;
          },
          { mutate: false },
        );
        this.root.organizations.value.get(subsidiaryId)?.invalidate();
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

  private async removeSubsidiary(organizationId: string) {
    try {
      this.isLoading = true;
      await this.transport.graphql.request<{
        organization_RemoveSubsidiary: Organization;
        REMOVE_SUBSIDIARY_FROM_ORGANIZATION: Organization;
      }>(REMOVE_SUBSIDIARY_FROM_ORGANIZATION_MUTATION, {
        organizationId: organizationId,
        subsidiaryId: this.id,
      });

      runInAction(() => {
        this.root.organizations.value.get(organizationId)?.invalidate();
        this.root.organizations.value.get(organizationId)?.update(
          (org: Organization) => {
            org.subsidiaries = [];

            return org;
          },
          { mutate: false },
        );
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

  private async updateOnboardingStatus() {
    try {
      await this.service.updateOnboardingStatus({
        input: {
          organizationId: this.id,
          status:
            this.value?.accountDetails?.onboarding?.status ??
            OnboardingStatus.NotApplicable,
          comments: this.value?.accountDetails?.onboarding?.comments ?? '',
        },
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    }
  }

  private async addTagsToOrganization(tagId: string, tagName: string) {
    try {
      this.isLoading = true;
      await this.transport.graphql.request<
        unknown,
        ADD_TAGS_TO_ORGANIZATION_PAYLOAD
      >(ADD_TAGS_TO_ORGANIZATION_MUTATION, {
        input: {
          organizationId: this.id,
          tag: {
            id: tagId,
            name: tagName,
          },
        },
      });
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

  private async removeTagsFromOrganization(tagId: string) {
    try {
      this.isLoading = true;
      await this.transport.graphql.request<
        unknown,
        REMOVE_TAGS_FROM_ORGANIZATION_PAYLOAD
      >(REMOVE_TAGS_FROM_ORGANIZATION_MUTATION, {
        input: {
          organizationId: this.id,
          tag: {
            id: tagId,
          },
        },
      });
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

  async save(operation: Operation) {
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
      .with(['contracts', ...P.array()], () => {})
      .with(['accountDetails', 'renewalSummary', ...P.array()], () => {
        this.updateAllOpportunityRenewals();
      })
      .with(['accountDetails', 'onboarding', ...P.array()], () => {
        this.updateOnboardingStatus();
      })
      .with(['stage', ...P.array()], () => {
        const payload = makePayload<OrganizationUpdateInput>(operation);
        this.updateOrganization(payload);
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
      .with(['subsidiaries', ...P.array()], () => {
        if (type === 'add') {
          this.addSubsidiary(
            value[0]?.organization?.metadata?.id ||
              value?.organization?.metadata?.id,
          );
        }
      })
      .with(['parentCompanies', ...P.array()], () => {
        if (type === 'delete') {
          this.removeSubsidiary(oldValue?.organization?.metadata?.id);
        }
      })
      .with(['tags', ...P.array()], () => {
        if (type === 'add') {
          this.addTagsToOrganization(value.id, value.name);
        }
        if (type === 'delete') {
          if (typeof oldValue === 'object') {
            this.removeTagsFromOrganization(oldValue.id);
          }
        }
        // if tag with index different that last one is deleted it comes as an update, bulk creation updates also come as updates
        if (type === 'update') {
          if (!oldValue) {
            (value as Array<Tag>)?.forEach((tag: Tag) => {
              this.addTagsToOrganization(tag.id, tag.name);
            });
          }
          if (oldValue) {
            this.removeTagsFromOrganization(oldValue);
          }
        }
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
      contracts {
        metadata {
          id
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
      tags {
        id
        name
        createdAt
        updatedAt
        appSource
      }
      valueProposition
      socialMedia {
        id
        url
        followersCount
      }
      employees
      yearFounded
      public

      accountDetails {
        churned
        ltv
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
        locality
        countryCodeA2
        countryCodeA3
      }
      subsidiaries {
        organization {
          metadata {
            id
          }
          name
          parentCompanies {
            organization {
              name
              metadata {
                id
              }
            }
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

type ADD_TAGS_TO_ORGANIZATION_PAYLOAD = {
  input: OrganizationTagInput;
};

const ADD_TAGS_TO_ORGANIZATION_MUTATION = gql`
  mutation addTagsToOrganization($input: OrganizationTagInput!) {
    organization_AddTag(input: $input) {
      accepted
    }
  }
`;

type REMOVE_TAGS_FROM_ORGANIZATION_PAYLOAD = {
  input: OrganizationTagInput;
};

const REMOVE_TAGS_FROM_ORGANIZATION_MUTATION = gql`
  mutation removeTagFromOrganization($input: OrganizationTagInput!) {
    organization_RemoveTag(input: $input) {
      accepted
    }
  }
`;

export const getDefaultValue = (): Organization => ({
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
    ltv: 0,
    churned: new Date().toISOString(),
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
    lastTouchPointTimelineEvent: ActionStore.getDefaultValue(),
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
});
