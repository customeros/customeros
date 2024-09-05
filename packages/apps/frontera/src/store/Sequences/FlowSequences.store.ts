import { Channel } from 'phoenix';
import { Store } from '@store/store';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';
import { FlowSequenceStore } from '@store/Sequences/FlowSequence.store';
import { FlowSequenceService } from '@store/Sequences/__service__/FlowSequence.service';
import { CreateSequenceMutationVariables } from '@store/Sequences/__service__/createSequence.generated';

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
      const { sequences } = await this.service.getSequences();

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

  async create(
    payload: CreateSequenceMutationVariables['input'],
    options?: { onSuccess?: (serverId: string) => void },
  ) {
    const newSequence = new FlowSequenceStore(this.root, this.transport);
    const tempId = newSequence.value.metadata?.id;

    newSequence.value = {
      ...newSequence.value,
    };

    let serverId: string | undefined;

    this.value.set(tempId, newSequence);

    try {
      const { flow_sequence_store } = await this.service.createSequence({
        input: payload,
      });

      runInAction(() => {
        serverId = flow_sequence_store?.metadata.id;
        newSequence.setId(serverId);
        newSequence.value = {
          ...newSequence.value,
          ...flow_sequence_store,
        };

        this.value.set(serverId, newSequence);
        this.value.delete(tempId);

        this.sync({ action: 'APPEND', ids: [serverId] });
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      serverId && options?.onSuccess?.(serverId);
      setTimeout(() => {
        if (serverId) {
          this.value.get(serverId)?.invalidate();
        }
      }, 1000);
    }
  }

  async remove() {
    // todo
  }
}
