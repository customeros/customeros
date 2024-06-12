import type { RootStore } from '@store/root';

import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { Operation } from '@store/types';
import { Transport } from '@store/transport';
// import { rdiffResult } from 'recursive-diff';
import { Store, makeAutoSyncable } from '@store/store';
import { runInAction, makeAutoObservable } from 'mobx';
import { makeAutoSyncableGroup } from '@store/group-store';

import {
  DataSource,
  Opportunity,
  InternalType,
  InternalStage,
  OpportunityRenewalLikelihood,
} from '@graphql/types';

export class OpportunityStore implements Store<Opportunity> {
  value: Opportunity = defaultValue;
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncable.load<Opportunity>();
  update = makeAutoSyncable.update<Opportunity>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncable(this, {
      channelName: 'Opportunity',
      mutator: this.save,
      getId: (d) => d?.id,
    });
  }

  get id() {
    return this.value.id;
  }
  set id(id: string) {
    this.value.id = id;
  }

  async invalidate() {
    try {
      this.isLoading = true;
      const { opportunity } = await this.transport.graphql.request<
        OPPORTUNITY_QUERY_RESULT,
        { id: string }
      >(OPORTUNITY_QUERY, { id: this.id });

      this.load(opportunity);
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
    // const diff = operation.diff?.[0];
    // const type = diff?.op;
    // const path = diff?.path;
    // const value = diff?.val;
    // const oldValue = (diff as rdiffResult & { oldVal: unknown })?.oldVal;
  }
}

type OPPORTUNITY_QUERY_RESULT = {
  opportunity: Opportunity;
};

const OPORTUNITY_QUERY = gql`
  query Opportunity($id: ID!) {
    opportunity(id: $id) {
      id
      createdAt
      updatedAt
      name
      amount
      maxAmount
      internalType
      externalType
      internalStage
      externalStage
      estimatedClosedAt
      generalNotes
      nextSteps
      renewedAt
      renewalApproved
      renewalLikelihood
      renewalUpdatedByUserId
      renewalUpdatedByUserAt
      renewalAdjustedRate
      comments
      createdBy {
        id
        firstName
        lastName
      }
      owner {
        id
        firstName
        lastName
      }
      source
      sourceOfTruth
      appSource
      externalLinks {
        externalId
        externalUrl
      }
    }
  }
`;

const defaultValue: Opportunity = {
  amount: 0,
  appSource: '',
  comments: '',
  createdAt: '',
  externalLinks: [],
  externalStage: '',
  externalType: '',
  generalNotes: '',
  id: '',
  internalStage: InternalStage.Open,
  internalType: InternalType.Nbo,
  maxAmount: 0,
  name: '',
  nextSteps: '',
  renewalAdjustedRate: 0,
  renewalApproved: false,
  renewalLikelihood: OpportunityRenewalLikelihood.ZeroRenewal,
  renewalUpdatedByUserId: '',
  source: DataSource.Na,
  sourceOfTruth: DataSource.Na,
  updatedAt: new Date().toISOString(),
};
