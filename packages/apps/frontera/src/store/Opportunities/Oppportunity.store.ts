import type { RootStore } from '@store/root';

import { Channel } from 'phoenix';
import { match } from 'ts-pattern';
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
  OpportunityUpdateInput,
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
      getId: (d) => d?.metadata.id,
    });
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

  private async updateOpportunityExternalStage(externalStage: string) {
    try {
      this.isLoading = true;
      await this.transport.graphql.request<
        unknown,
        OPPORTUNITY_UPDATE_STAGE_PAYLOAD
      >(OPPORTUNITY_UPDATE_STAGE, {
        input: { opportunityId: this.id, externalStage, internalStage: 'OPEN' },
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
  private async updateOpportunityCloseLost() {
    try {
      this.isLoading = true;
      await this.transport.graphql.request<
        unknown,
        OPPORTUNITY_UPDATE_CLOSE_LOST_PAYLOAD
      >(OPPORTUNITY_UPDATE_CLOSE_LOST, { opportunityId: this.id });
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
  private async updateOpportunityCloseWon() {
    try {
      this.isLoading = true;
      await this.transport.graphql.request<
        unknown,
        OPPORTUNITY_UPDATE_CLOSE_WON_PAYLOAD
      >(OPPORTUNITY_UPDATE_CLOSE_WON, { opportunityId: this.id });
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
    const path = diff?.path;
    const value = diff?.val;
    match(path)
      .with(['externalStage'], () => {
        this.updateOpportunityExternalStage(value as string);
      })
      .with(['internalStage'], () => {
        match(value)
          .with(InternalStage.ClosedLost, () => {
            this.updateOpportunityCloseLost();
          })
          .with(InternalStage.ClosedWon, () => {
            this.updateOpportunityCloseWon();
          });
      });
  }
}

type OPPORTUNITY_UPDATE_STAGE_PAYLOAD = {
  input: OpportunityUpdateInput;
};

const OPPORTUNITY_UPDATE_STAGE = gql`
  mutation OpportunityUpdateStage($input: OpportunityUpdateInput!) {
    opportunity_Update(input: $input) {
      id
    }
  }
`;

type OPPORTUNITY_UPDATE_CLOSE_WON_PAYLOAD = {
  opportunityId: string;
};

const OPPORTUNITY_UPDATE_CLOSE_WON = gql`
  mutation OpportunityUpdateCloseWon($opportunityId: ID!) {
    opportunity_CloseWon(opportunityId: $opportunityId) {
      accepted
    }
  }
`;

type OPPORTUNITY_UPDATE_CLOSE_LOST_PAYLOAD = {
  opportunityId: string;
};

const OPPORTUNITY_UPDATE_CLOSE_LOST = gql`
  mutation OpportunityUpdateCloseLost($opportunityId: ID!) {
    opportunity_CloseLost(opportunityId: $opportunityId) {
      accepted
    }
  }
`;

type OPPORTUNITY_QUERY_RESULT = {
  opportunity: Opportunity;
};

const OPORTUNITY_QUERY = gql`
  query Opportunity($id: ID!) {
    opportunity(id: $id) {
      metadata {
        id
        created
        lastUpdated
        source
        sourceOfTruth
        appSource
      }
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
      organization {
        metadata {
          id
          created
          lastUpdated
          sourceOfTruth
        }
      }
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
  metadata: {
    id: '',
    created: '',
    lastUpdated: '',
    source: DataSource.Na,
    sourceOfTruth: DataSource.Na,
    appSource: '',
  },
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
