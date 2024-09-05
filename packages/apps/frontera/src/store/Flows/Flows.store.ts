import { Channel } from 'phoenix';
import { Store } from '@store/store';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { FlowStore } from '@store/Flows/Flow.store';
import { FlowService } from '@store/Flows/__service__';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';
import { FlowSequenceStore } from '@store/Sequences/FlowSequence.store.ts';
import { CreateSequenceMutationVariables } from '@store/Sequences/__service__/createSequence.generated.ts';

import { Flow } from '@graphql/types';

export class FlowsStore implements GroupStore<Flow> {
  version = 0;
  isLoading = false;
  history: GroupOperation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  isBootstrapped: boolean = false;
  value: Map<string, Store<Flow>> = new Map();
  sync = makeAutoSyncableGroup.sync;
  subscribe = makeAutoSyncableGroup.subscribe;
  load = makeAutoSyncableGroup.load<Flow>();
  totalElements = 0;
  private service: FlowService;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'Flows',
      getItemId: (item) => item?.metadata?.id,
      ItemStore: FlowStore,
    });
    this.service = FlowService.getInstance(transport);
  }

  toArray() {
    return Array.from(this.value.values());
  }

  get educationFlow() {
    return this.toArray().find(
      (flow) => flow.value.name?.toLowerCase() === 'education',
    );
  }

  toComputedArray(compute: (arr: FlowStore[]) => FlowStore[]) {
    const arr = this.toArray();

    return compute(arr as FlowStore[]);
  }

  async bootstrap() {
    if (this.root.demoMode) {
      this.isBootstrapped = true;

      return;
    }

    if (this.isBootstrapped || this.isLoading) return;

    try {
      const { flows } = await this.service.getFlows();

      runInAction(() => {
        this.load(flows);
      });
      runInAction(() => {
        this.isBootstrapped = true;
        this.totalElements = flows.length;
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
      const { flow_sequence_store } = await this.service.cre({
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
