import { Channel } from 'phoenix';
import { Store } from '@store/store';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { FlowStore } from '@store/Flows/Flow.store';
import { FlowService } from '@store/Flows/__service__';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

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
