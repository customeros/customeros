import type { RootStore } from '@store/root';

import set from 'lodash/set';
import get from 'lodash/get';
import omit from 'lodash/omit';
import { countryMap } from '@assets/countries/countriesMap';
import { ActionStore } from '@store/TimelineEvents/Actions/Action.store';

import {
  User,
  Market,
  type Contact,
  FundingRound,
  type Contract,
  OnboardingStatus,
  OrganizationStage,
  LastTouchpointType,
  OrganizationRelationship,
  OpportunityRenewalLikelihood,
} from '@graphql/types';

import type { OrganizationQuery } from './__service__/getOrganization.generated';

type OrganizationDatum = NonNullable<OrganizationQuery['organization']>;

export type Organization = OrganizationDTO & OrganizationDatum;

function makePersistablePayload<T extends object>(
  instance: T,
  fields: string[],
) {
  const out = omit(instance, ['root', ...fields.map((f) => '_' + f)]);

  fields.forEach((field) => set(out, field, get(instance, '_' + field)));

  return out;
}

const privateFields = [
  'owner',
  'contacts',
  'contracts',
  'subsidiaries',
  'parentCompanies',
];

export class OrganizationDTO {
  private _owner: OrganizationDatum['owner'];
  private _contacts: OrganizationDatum['contacts'];
  private _contracts: OrganizationDatum['contracts'];
  private _subsidiaries: OrganizationDatum['subsidiaries'];
  private _parentCompanies: OrganizationDatum['parentCompanies'];

  constructor(private root: RootStore, raw: OrganizationDatum) {
    Object.assign(this, omit(raw, privateFields) as OrganizationDatum);

    this._owner = raw?.owner;
    this._contacts = raw?.contacts ?? [];
    this._contracts = raw?.contracts ?? [];
    this._subsidiaries = raw?.subsidiaries ?? [];
    this._parentCompanies = raw?.parentCompanies ?? [];
  }

  get id() {
    return get(this as unknown as Organization, 'metadata.id');
  }

  set id(value: string) {
    set(this, 'metadata.id', value);
  }

  get owner() {
    return this.root.users.value.get(this._owner?.id as string);
  }

  set owner(user: User) {
    console.log('Setez pla', user);
    this._owner = user;
  }

  get contacts() {
    return this._contacts.content.reduce((acc, { metadata }) => {
      const store = this.root.contacts.value.get(metadata?.id);

      if (store) acc.push(store.value);

      return acc;
    }, [] as Contact[]);
  }

  get contracts() {
    return this._contracts?.reduce((acc, { metadata }) => {
      const store = this.root.contracts.value.get(metadata.id);

      if (store) acc.push(store.value);

      return acc;
    }, [] as Contract[]);
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
    const self = this as unknown as Organization & OrganizationDatum;
    const code = self.locations?.[0]?.countryCodeA2;

    if (!code) return;

    return countryMap.get(code.toLowerCase());
  }

  get parentCompanies() {
    return this._parentCompanies.reduce((acc, curr) => {
      // TODO: see if we need to getById if this becomes observable
      const id = curr?.organization?.metadata?.id;
      const store = this.root.organizations.value.get(id);

      if (store) acc.push(store);

      return acc;
    }, [] as Organization[]);
  }

  set parentCompanies(data: Organization[]) {
    this._parentCompanies = data.map((obj) => {
      return {
        organization: obj,
      };
    });
  }

  get subsidiaries() {
    return this._subsidiaries.reduce((acc, curr) => {
      // TODO: see if we need to getById if this becomes observable
      const id = curr?.organization?.metadata?.id;
      const store = this.root.organizations.value.get(id);

      if (store) acc.push(store);

      return acc;
    }, [] as Organization[]);
  }

  set subsidiaries(data: Organization[]) {
    this._subsidiaries = data.map((obj) => {
      return {
        organization: obj,
      };
    });
  }

  public addSubsidiary(id: string) {
    const organization = this.root.organizations.value.get(id);

    if (!organization) return;

    this._subsidiaries.push({
      organization,
    });
  }

  public removeSubsidiary(id: string) {
    const removeIndex = this._subsidiaries.findIndex(
      (org) => org.organization.metadata.id === id,
    );

    this._subsidiaries.splice(removeIndex, 1);
  }

  public clearParentCompanies() {
    this._parentCompanies = [];
  }

  public setOwner(userId: string) {
    const user = this.root.users.value.get(userId);

    if (!user) return;

    this._owner = user.value;
  }

  public clearOwner() {
    this._owner = null;
  }

  static of(root: RootStore, raw: OrganizationDatum) {
    return new OrganizationDTO(root, raw) as Organization;
  }

  static toPersistable(instance: Organization) {
    return makePersistablePayload(instance, privateFields);
  }

  static default(): OrganizationDatum {
    return {
      name: 'Unnamed',
      metadata: {
        id: crypto.randomUUID(),
        created: new Date().toISOString(),
      },
      owner: null,
      contactCount: 0,
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
        // @ts-expect-error ignore for now
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
    };
  }
}
