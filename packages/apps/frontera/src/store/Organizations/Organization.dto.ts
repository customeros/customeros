import type { UserStore } from '@store/Users/User.store';

import merge from 'lodash/merge';
import { Entity } from '@store/record';
import { ContactDatum } from '@store/Contacts/Contact.dto';
import { countryMap } from '@assets/countries/countriesMap';
import { action, computed, observable, runInAction } from 'mobx';

import {
  FundingRound,
  type Contract,
  OnboardingStatus,
  OrganizationStage,
  LastTouchpointType,
  OrganizationRelationship,
} from '@graphql/types';

import type { GetOrganizationsByIdsQuery } from './__service__/getOrganizationsByIds.generated';
import type { SaveOrganizationMutationVariables } from './__service__/saveOrganization.generated';

import { OrganizationsStore } from './Organizations.store';

export type OrganizationDatum = NonNullable<
  GetOrganizationsByIdsQuery['ui_organizations'][number]
>;

export class Organization extends Entity<OrganizationDatum> {
  @observable accessor value: OrganizationDatum = Organization.default();

  constructor(store: OrganizationsStore, data: OrganizationDatum) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    super(store as any, data);
  }

  @computed
  get id() {
    return this.value.id;
  }

  @computed
  get name() {
    return this.value.name;
  }

  set id(value: string) {
    runInAction(() => {
      this.value.id = value;
    });
  }

  @computed
  get owner(): UserStore | null {
    if (!this.value.owner) return null;
    const user = this.store.root.users.value.get(
      this.value?.owner.id as string,
    );

    return user ?? null;
  }

  @computed
  get primaryDomains() {
    if (!this.value.domainsDetails) return [];

    return this.value.domainsDetails
      .filter((e) => e.primary)
      .map((e) => e.domain);
  }

  @computed
  get isEnriching(): boolean {
    return (
      this.value?.enrichedRequestedAt &&
      !this.value?.enrichedAt &&
      !this.value?.enrichedFailedAt
    );
  }

  @computed
  get contacts() {
    this.store.root.contacts.retrieve(this.value.contacts);

    return this.value.contacts.reduce((acc, id) => {
      const record = this.store.root.contacts.getById(id);

      if (record) acc.push(record.value);

      return acc;
    }, [] as ContactDatum[]);
  }

  @computed
  get contracts() {
    return this.value.contracts?.reduce((acc, id) => {
      const store = this.store.root.contracts.value.get(id);

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
    return this.value.parentId
      ? [this.store.getById(this.value.parentId)?.value]
      : [null];
  }

  @computed
  get subsidiaries() {
    return this.value.subsidiaries.reduce((acc, id) => {
      const record = this.store.getById(id);

      if (record) acc.push(record.value);

      return acc;
    }, [] as OrganizationDatum[]);
  }

  @computed
  get tagCount() {
    return this.value.tags?.length ?? 0;
  }

  @action
  public addSubsidiary(id: string) {
    this.value.subsidiaries.push(id);
  }

  @action
  public removeSubsidiary(id: string) {
    const removeIndex = this.value.subsidiaries.indexOf(id);

    this.value.subsidiaries.splice(removeIndex, 1);
  }

  @action
  public addParent(id: string) {
    const record = this.store.getById(id);

    if (!record) return;

    this.value.parentId = id;
    this.value.parentName = record.value.name;
  }

  @action
  public clearParent() {
    this.value.parentId = null;
    this.value.parentName = null;
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

  @action
  public addSocial(url: string) {
    this.value.socialMedia.push({
      id: crypto.randomUUID(),
      url,
      followersCount: 0,
      __typename: 'Social',
      alias: '',
    });
  }

  @action
  public deleteTag(id: string) {
    const idx = this.value.tags?.findIndex((t) => t.metadata.id === id);

    if (idx === -1 || idx === undefined || idx === null) return;
    this.value.tags?.splice(idx, 1);
  }

  static default(
    payload?: OrganizationDatum | SaveOrganizationMutationVariables['input'],
  ): OrganizationDatum {
    return merge(
      {
        id: crypto.randomUUID(),
        name: 'Unnamed',
        notes: '',
        description: '',
        industry: '',
        market: '',
        website: '',
        logoUrl: '',
        iconUrl: '',
        public: false,
        stage: OrganizationStage.Target,
        relationship: OrganizationRelationship.Prospect,
        lastFundingRound: FundingRound.PreSeed,
        leadSource: '',
        valueProposition: '',
        slackChannelId: '',
        employees: 0,
        yearFounded: '',
        enrichedAt: null,
        enrichedFailedAt: null,
        enrichedRequestedAt: null,
        ltv: 0,
        hide: false,
        domains: [],
        domainsDetails: [],
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        churnedAt: null,
        customerOsId: '',
        referenceId: '',
        renewalSummaryArrForecast: null,
        renewalSummaryMaxArrForecast: null,
        renewalSummaryRenewalLikelihood: null,
        renewalSummaryNextRenewalAt: '',
        onboardingStatus: OnboardingStatus.NotApplicable,
        onboardingStatusUpdatedAt: '',
        onboardingComments: '',
        lastTouchPointAt: new Date().toISOString(),
        lastTouchPointType: LastTouchpointType.ActionCreated,
        contactCount: 0,
        parentId: null,
        parentName: null,
        contracts: [],
        contacts: [],
        subsidiaries: [],
        owner: null,
        tags: [],
        socialMedia: [],
        locations: [],
      },
      payload ?? {},
    );
  }
}
