import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { RootStore } from '@store/root';
import { Operation } from '@store/types';
import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';
import { Store, makeAutoSyncable } from '@store/store';

import {
  DataSource,
  Opportunity,
  InternalType,
  InternalStage,
  ServiceLineItem,
  OpportunityRenewalLikelihood,
} from '@graphql/types';

export class OpportunityStore implements Store<ServiceLineItem> {
  value: Opportunity = defaultValue;
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  load = makeAutoSyncable.load<Opportunity>();
  update = makeAutoSyncable.update<Opportunity>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, {
      channelName: 'Opportunity',
      mutator: this.save,
      getId: (d) => d?.id,
    });
    makeAutoObservable(this);
  }

  async invalidate() {
    try {
      this.isLoading = true;
      const { opportunity } = await this.transport.graphql.request<
        OPPORTUNITY_QUERY_RESULT,
        { id: string }
      >(OPPORTUNITY_QUERY, { id: this.id });

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

  set id(id: string) {
    this.value.id = id;
  }

  private async save() {
    // const payload: PAYLOAD = {
    //   input: {
    //     ...omit(this.value, 'metadata', 'owner'),
    //     contractId: this.value.metadata.id,
    //   },
    // };
    try {
      this.isLoading = true;
      // await this.transport.graphql.request(UPDATE_CONTRACT_DEF, payload);
    } catch (e) {
      this.error = (e as Error)?.message;
    } finally {
      this.isLoading = false;
    }
  }
}

type OPPORTUNITY_QUERY_RESULT = {
  opportunity: Opportunity;
};
const OPPORTUNITY_QUERY = gql`
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
      estimatedClosed
      generalNotes
      nextSteps
      renewed
      renewalApproved
      renewalLikelihood
      renewalUpdatedByUserId
      renewalUpdatedByUserAt
      renewalAdjustedRate
      comments
      createdBy {
        id
        name
      }
      owner {
        id
        name
      }
      source
      sourceOfTruth
      appSource
      externalLinks {
        type
        syncDate
        externalId
        externalUrl
        externalSource
      }
    }
  }
`;

const defaultValue: Opportunity = {
  id: crypto.randomUUID(),
  appSource: DataSource.Openline,
  created: new Date().toISOString(),
  lastUpdated: new Date().toISOString(),
  source: DataSource.Openline,
  sourceOfTruth: DataSource.Openline,
  name: '',
  amount: 0,
  maxAmount: 0,
  internalType: InternalType.Nbo,
  externalType: '',
  internalStage: InternalStage.ClosedLost,
  externalStage: '',
  estimatedClosedAt: new Date().toISOString(),
  generalNotes: '',
  nextSteps: '',
  renewedAt: new Date().toISOString(),
  renewalApproved: false,
  renewalLikelihood: OpportunityRenewalLikelihood.LowRenewal,
  renewalUpdatedByUserId: '',
  renewalUpdatedByUserAt: new Date().toISOString(),
  renewalAdjustedRate: 0,
  comments: '',
  createdBy: {
    id: crypto.randomUUID(),
    name: '',
  },
  owner: {
    id: crypto.randomUUID(),
    name: '',
  },
  externalLinks: [],
};
