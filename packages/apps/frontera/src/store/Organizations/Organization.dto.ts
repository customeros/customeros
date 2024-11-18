import type { UserStore } from '@store/Users/User.store';

import { Store } from '@store/_store';
import { Record } from '@store/record';
import { action, computed, runInAction } from 'mobx';
import { countryMap } from '@assets/countries/countriesMap';

import {
  Market,
  FundingRound,
  type Contact,
  type Contract,
  OnboardingStatus,
  OrganizationStage,
  LastTouchpointType,
  OrganizationRelationship,
  OpportunityRenewalLikelihood,
} from '@graphql/types';

import type { OrganizationsStore } from './Organizations.store';
import type { OrganizationQuery } from './__service__/getOrganization.generated';

export type OrganizationDatum = NonNullable<OrganizationQuery['organization']>;

export class Organization extends Record<OrganizationDatum> {
  constructor(store: OrganizationsStore, data: OrganizationDatum) {
    super(store, data);
  }

  @computed
  get id() {
    return this.value.metadata.id;
  }

  set id(value: string) {
    runInAction(() => {
      this.value.metadata.id = value;
    });
  }

  @computed
  get owner(): UserStore | null {
    const user = this.store.root.users.value.get(this.owner?.id as string);

    return user ?? null;
  }

  @computed
  get contacts() {
    return this.value.contacts.content.reduce((acc, { metadata }) => {
      const store = this.store.root.contacts.value.get(metadata?.id);

      if (store) acc.push(store.value);

      return acc;
    }, [] as Contact[]);
  }

  @computed
  get contracts() {
    return this.value.contracts?.reduce((acc, { metadata }) => {
      const store = this.store.root.contracts.value.get(metadata.id);

      if (store) acc.push(store.value);

      return acc;
    }, [] as Contract[]);
  }

  @computed
  get invoices() {
    return this.store.root.invoices
      .toArray()
      .filter(
        (invoice) =>
          invoice?.value?.organization?.metadata?.id === this.id &&
          !invoice?.value?.dryRun,
      );
  }

  @computed
  get country() {
    const code = this.value.locations?.[0]?.countryCodeA2;

    if (!code) return;

    return countryMap.get(code.toLowerCase());
  }

  @computed
  get parentCompanies() {
    return this.value.parentCompanies.reduce((acc, curr) => {
      const id = curr?.organization?.metadata?.id;
      const store = this.store.getById(id);

      if (store) acc.push(store.value);

      return acc;
    }, [] as OrganizationDatum[]);
  }

  @computed
  get subsidiaries() {
    return this.value.subsidiaries.reduce((acc, curr) => {
      const id = curr?.organization?.metadata?.id;
      const record = this.store.getById(id);

      if (record) acc.push(record.value);

      return acc;
    }, [] as OrganizationDatum[]);
  }

  @action
  public addSubsidiary(id: string) {
    const record = this.store.getById(id);

    if (!record) return;

    this.value.subsidiaries.push({
      organization: record.value,
    });
  }

  @action
  public removeSubsidiary(id: string) {
    const removeIndex = this.value.subsidiaries.findIndex(
      (org) => org.organization.metadata.id === id,
    );

    this.value.subsidiaries.splice(removeIndex, 1);
  }

  @action
  public addParent(id: string) {
    const record = this.store.getById(id);

    if (!record) return;

    this.value.parentCompanies.push({ organization: record.value });
  }

  @action
  public clearParentCompanies() {
    this.value.parentCompanies = [];
  }

  @action
  public setOwner(userId: string) {
    const record = this.store.root.users.value.get(userId);

    if (!record) return;

    this.value.owner = record.value;
  }

  @action
  public clearOwner() {
    this.value.owner = null;
  }

  static default(): OrganizationDatum {
    return {
      name: 'Unnamed',
      metadata: {
        id: crypto.randomUUID(),
        created: new Date().toISOString(),
      },
      owner: null,
      contacts: {
        content: [],
      },
      icon: '',
      referenceId: '',
      yearFounded: '',
      enrichDetails: {
        failedAt: '',
        enrichedAt: '',
        requestedAt: '',
      },
      customerOsId: '',
      domains: [],
      industry: '',
      locations: [],
      parentCompanies: [],
      socialMedia: [],
      stage: OrganizationStage.Target,
      tags: [],
      subsidiaries: [],
      website: '',
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
      description: '',
      employees: 0,
      isCustomer: false,
      logo: '',
      lastFundingRound: FundingRound.PreSeed,
      lastTouchpoint: {
        lastTouchPointTimelineEventId: crypto.randomUUID(),
        lastTouchPointAt: new Date().toISOString(),
        lastTouchPointType: LastTouchpointType.ActionCreated,
        // @ts-expect-error ignore for now
        lastTouchPointTimelineEvent: ActionStore.getDefaultValue(),
      }, // nested defaults ignored for now -> should be converted into a Store
      leadSource: '',
      market: Market.B2B,
      public: false,
      relationship: OrganizationRelationship.Prospect,
      // slackChannelId: '',
      // stageLastUpdated: '',
      // subIndustry: '',
      // targetAudience: '',
      valueProposition: '',
    };
  }
}
