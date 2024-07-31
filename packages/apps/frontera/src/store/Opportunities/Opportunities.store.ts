import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { when, runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import { Pagination, Opportunity } from '@graphql/types';

import mock from './mock.json';
import { OpportunityStore } from './Opportunity.store';

export class OpportunitiesStore implements GroupStore<Opportunity> {
  version = 0;
  isLoading = false;
  totalElements = 0;
  history: GroupOperation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  isBootstrapped: boolean = false;
  value: Map<string, OpportunityStore> = new Map();
  sync = makeAutoSyncableGroup.sync;
  subscribe = makeAutoSyncableGroup.subscribe;
  load = makeAutoSyncableGroup.load<Opportunity>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'Opportunities',
      getItemId: (item) => item?.metadata?.id,
      ItemStore: OpportunityStore,
    });

    when(
      () => this.isBootstrapped && this.totalElements > 0,
      async () => {
        await this.bootstrapRest();
      },
    );
  }

  async bootstrap() {
    if (this.root.demoMode) {
      this.load(
        mock.data.opportunities_LinkedToOrganizations
          .content as unknown as Opportunity[],
      );
      this.totalElements =
        mock.data.opportunities_LinkedToOrganizations.totalElements;
      this.isBootstrapped = true;

      return;
    }

    if (this.isBootstrapped || this.isLoading) return;

    try {
      this.isLoading = true;

      const { opportunities_LinkedToOrganizations } =
        await this.transport.graphql.request<
          OPPORTUNITIES_QUERY_RESPONSE,
          OPPORTUNITIES_QUERY_PAYLOAD
        >(OPPORTUNITIES_QUERY, {
          pagination: { limit: 1000, page: 1 },
        });

      this.load(opportunities_LinkedToOrganizations.content);
      runInAction(() => {
        this.isBootstrapped = true;
        this.totalElements = opportunities_LinkedToOrganizations.totalElements;
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

  async bootstrapRest() {
    let page = 1;

    while (this.totalElements > this.value.size) {
      try {
        const { opportunities_LinkedToOrganizations } =
          await this.transport.graphql.request<
            OPPORTUNITIES_QUERY_RESPONSE,
            OPPORTUNITIES_QUERY_PAYLOAD
          >(OPPORTUNITIES_QUERY, {
            pagination: { limit: 1000, page },
          });

        runInAction(() => {
          page++;
          this.load(opportunities_LinkedToOrganizations.content);
        });
      } catch (e) {
        runInAction(() => {
          this.error = (e as Error)?.message;
        });
        break;
      }
    }
  }

  toArray() {
    return Array.from(this.value.values());
  }

  toComputedArray(compute: (arr: OpportunityStore[]) => OpportunityStore[]) {
    const arr = this.toArray();

    return compute(arr);
  }
}

type OPPORTUNITIES_QUERY_PAYLOAD = {
  pagination: Pagination;
};

type OPPORTUNITIES_QUERY_RESPONSE = {
  opportunities_LinkedToOrganizations: {
    content: [];
    totalElements: number;
    totalAvailable: number;
  };
};

const OPPORTUNITIES_QUERY = gql`
  query getOpportunities($pagination: Pagination!) {
    opportunities_LinkedToOrganizations(pagination: $pagination) {
      content {
        metadata {
          id
          created
          lastUpdated
          source
          sourceOfTruth
          appSource
        }
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
        currency
        stageLastUpdated
        renewalApproved
        renewalLikelihood
        renewalUpdatedByUserId
        renewalUpdatedByUserAt
        renewalAdjustedRate
        comments
        organization {
          metadata {
            id
            created
            lastUpdated
            sourceOfTruth
          }
        }
        createdBy {
          id
          firstName
          lastName
          name
        }
        owner {
          id
          firstName
          lastName
          name
        }
        externalLinks {
          externalUrl
          externalId
        }
        id
        createdAt
        updatedAt
        source
        appSource
      }
      totalElements
      totalAvailable
    }
  }
`;
