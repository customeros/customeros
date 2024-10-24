import type { RootStore } from '@store/root';
import type { Transport } from '@store/transport';

import { Syncable } from '@store/syncable';
import { countryMap } from '@assets/countries/countriesMap';
import { ActionStore } from '@store/TimelineEvents/Actions/Action.store';
import { action, override, computed, runInAction, makeObservable } from 'mobx';

import {
  Market,
  Contact,
  Contract,
  DataSource,
  FundingRound,
  Organization,
  OnboardingStatus,
  OrganizationStage,
  LastTouchpointType,
  OrganizationRelationship,
  OpportunityRenewalLikelihood,
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

    makeObservable<OrganizationStore>(this, {
      id: override,
      save: override,
      getId: override,
      setId: override,
      invalidate: action,
      contacts: computed,
      contracts: computed,
      subsidiaries: computed,
      country: computed,
      getChannelName: override,
      parentCompanies: computed,
      invoices: computed,
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
    return 'Organizations';
  }

  deleteTag(tagId: string) {
    this.value.tags = (this.value?.tags || []).filter(
      (tag) => tag.id !== tagId,
    );
  }

  async invalidate() {
    try {
      this.isLoading = true;

      const { organization } = await this.service.getOrganization(this.id);

      this.load(organization as Organization);
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

  static getDefaultValue(): Organization {
    return {
      name: 'Unnamed',
      metadata: {
        id: crypto.randomUUID(),
        created: new Date().toISOString(),
        lastUpdated: new Date().toISOString(),
        appSource: DataSource.Openline,
        source: DataSource.Openline,
        sourceOfTruth: DataSource.Openline,
      },
      icpFit: false,
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
      hide: false,
      inboundCommsCount: 0,
      issueSummaryByStatus: [],
      jobRoles: [],
      locations: [],
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
      enrichDetails: {
        enrichedAt: '',
        failedAt: '',
        requestedAt: '',
      },
    };
  }
}

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
  icpFit: false,
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
  hide: false,
  inboundCommsCount: 0,
  issueSummaryByStatus: [],
  jobRoles: [],
  locations: [],
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
  enrichDetails: {
    enrichedAt: '',
    failedAt: '',
    requestedAt: '',
  },
});
