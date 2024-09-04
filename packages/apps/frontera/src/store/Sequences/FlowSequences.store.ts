import { Channel } from 'phoenix';
import { Store } from '@store/store';
import { gql } from 'graphql-request';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';
import { FlowSequenceStore } from '@store/Sequences/FlowSequence.store';
import { FlowSequenceService } from '@store/Sequences/FlowSequence.service.ts';

import { FlowSequence } from '@graphql/types';

export class FlowSequencesStore implements GroupStore<FlowSequence> {
  version = 0;
  isLoading = false;
  history: GroupOperation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  isBootstrapped: boolean = false;
  value: Map<string, Store<FlowSequence>> = new Map();
  sync = makeAutoSyncableGroup.sync;
  subscribe = makeAutoSyncableGroup.subscribe;
  load = makeAutoSyncableGroup.load<FlowSequence>();
  totalElements = 0;
  private service: FlowSequenceService;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'FlowSequences',
      getItemId: (item) => item?.metadata?.id,
      ItemStore: FlowSequenceStore,
    });
    this.service = FlowSequenceService.getInstance(transport);
  }

  toArray() {
    return Array.from(this.value.values());
  }

  toComputedArray(compute: (arr: FlowSequenceStore[]) => FlowSequenceStore[]) {
    const arr = this.toArray();

    return compute(arr as FlowSequenceStore[]);
  }

  async bootstrap() {
    if (this.root.demoMode) {
      this.isBootstrapped = true;

      return;
    }

    if (this.isBootstrapped || this.isLoading) return;

    try {
      const { sequences } =
        await this.transport.graphql.request<SEQUENCES_RESPONSE>(
          SEQUENCES_QUERY,
        );

      runInAction(() => {
        this.load(sequences);
      });
      runInAction(() => {
        this.isBootstrapped = true;
        this.totalElements = sequences.length;
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

  async invalidate() {
    this.isLoading = true;
  }

  async create() {
    // todo
  }

  async remove() {
    // todo
  }
}

type SEQUENCES_RESPONSE = {
  sequences: FlowSequence[];
};
const SEQUENCES_QUERY = gql`
  query getFlowSequences {
    sequences {
      metadata {
        id
      }
      contacts {
        metadata {
          id
        }
      }
      name
      status
      description
      steps {
        Type
        Text
        status
        email {
          email
        }

        metadata {
          id
        }
      }
      flow {
        description
        name
        status
        metadata {
          id
        }
      }
      mailboxes {
        email
        metadata {
          id
        }
      }
    }
  }
`;
