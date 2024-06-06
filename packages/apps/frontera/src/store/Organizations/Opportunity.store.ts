import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { RootStore } from '@store/root';
import { Operation } from '@store/types';
import { rdiffResult } from 'recursive-diff';
import { makePayload } from '@store/util.ts';
import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';
import { Store, makeAutoSyncable } from '@store/store';

import {
  DataSource,
  Opportunity,
  InternalType,
  InternalStage,
  ServiceLineItem,
  OpportunityUpdateInput,
  OpportunityRenewalLikelihood,
  OpportunityRenewalUpdateInput,
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
    makeAutoObservable(this);

    makeAutoSyncable(this, {
      channelName: 'Opportunity',
      mutator: this.save,
      getId: (d: Opportunity) => d?.id,
    });
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
  get id() {
    return this.value.id;
  }

  private async updateOpportunity(payload: OpportunityUpdateInput) {
    try {
      this.isLoading = true;
      await this.transport.graphql.request<unknown, UPDATE_OPPORTUNITY_PAYLOAD>(
        UPDATE_OPPORTUNITY_MUTATION,
        {
          input: {
            ...payload,
            opportunityId: this.id,
          },
        },
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
  private async updateOpportunityRenewal(
    payload: OpportunityRenewalUpdateInput,
  ) {
    try {
      this.isLoading = true;
      const input = {
        ...payload,
        opportunityId: this.id,
      };
      await this.transport.graphql.request<
        unknown,
        UPDATE_OPPORTUNITY_RENEWAL_PAYLOAD
      >(UPDATE_OPPORTUNITY_RENEWAL_MUTATION, {
        input,
      });

      runInAction(() => {
        this.invalidate();
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
    const payload = makePayload<OpportunityUpdateInput>(operation);
    this.value.amount =
      this.value.maxAmount * (payload.renewalAdjustedRate / 100);
    this.value.renewalLikelihood =
      payload?.renewalLikelihood ||
      (payload.renewalAdjustedRate <= 25 &&
        OpportunityRenewalLikelihood.LowRenewal) ||
      (payload.renewalAdjustedRate <= 75 &&
        payload.renewalAdjustedRate > 25 &&
        OpportunityRenewalLikelihood.MediumRenewal) ||
      OpportunityRenewalLikelihood.HighRenewal;
    this.updateOpportunityRenewal(payload);
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

type UPDATE_OPPORTUNITY_RENEWAL_PAYLOAD = {
  input: OpportunityRenewalUpdateInput;
};

const UPDATE_OPPORTUNITY_RENEWAL_MUTATION = gql`
  mutation updateOpportunityRenewal($input: OpportunityRenewalUpdateInput!) {
    opportunityRenewalUpdate(input: $input) {
      id
    }
  }
`;

type UPDATE_OPPORTUNITY_PAYLOAD = {
  input: OpportunityUpdateInput;
};

const UPDATE_OPPORTUNITY_MUTATION = gql`
  mutation updateOpportunity($input: OpportunityUpdateInput) {
    opportunityUpdate(input: $input) {
      id
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
