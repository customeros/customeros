import type { RootStore } from '@store/root';

import { Channel } from 'phoenix';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import { Workflow } from '@shared/types/__generated__/graphql.types';

import { WorkFlowStore } from './WorkFlow.store';
import { WorkFlowService } from './WorkFlow.service';

export class WorkFlowsStore implements GroupStore<Workflow> {
  version = 0;
  isLoading = false;
  history: GroupOperation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  isBootstrapped: boolean = false;
  value: Map<string, WorkFlowStore> = new Map();
  sync = makeAutoSyncableGroup.sync;
  subscribe = makeAutoSyncableGroup.subscribe;
  load = makeAutoSyncableGroup.load<Workflow>();
  totalElements = 0;
  private service: WorkFlowService;

  constructor(public root: RootStore, public transport: Transport) {
    this.service = WorkFlowService.getInstance(transport);

    makeAutoSyncableGroup(this, {
      channelName: 'WorkFlows',
      ItemStore: WorkFlowStore,
      getItemId: (item) => item.id,
    });
    makeAutoObservable(this);
  }

  async bootstrap() {
    if (this.isBootstrapped) return;

    try {
      await this.service.getWorkFlowsByType();
    } catch (e) {
      this.error = (e as Error).message;
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  getById(id: string) {
    return this.value.get(id);
  }

  toArray(): WorkFlowStore[] {
    return Array.from(this.value)?.flatMap(
      ([, WorkFlowStore]) => WorkFlowStore,
    );
  }
}
