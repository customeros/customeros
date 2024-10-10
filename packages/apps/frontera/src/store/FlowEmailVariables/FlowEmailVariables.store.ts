import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';
import { FlowEmailVariablesService } from '@store/FlowEmailVariables/__service__';
import { FlowEmailVariableStore } from '@store/FlowEmailVariables/FlowEmailVariable.store.ts';

import { EmailVariableEntity } from '@graphql/types';

export class FlowEmailVariablesStore
  implements GroupStore<EmailVariableEntity>
{
  version = 0;
  isLoading = false;
  history: GroupOperation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  isBootstrapped: boolean = false;
  value: Map<string, FlowEmailVariableStore> = new Map();
  sync = makeAutoSyncableGroup.sync;
  subscribe = makeAutoSyncableGroup.subscribe;
  load = makeAutoSyncableGroup.load<EmailVariableEntity>();
  totalElements = 0;
  private service: FlowEmailVariablesService;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'FlowEmailVariables',
      getItemId: (item) => item?.type,
      ItemStore: FlowEmailVariableStore,
    });
    this.service = FlowEmailVariablesService.getInstance(transport);
  }

  toArray() {
    return Array.from(this.value.values());
  }

  async bootstrap() {
    if (this.root.demoMode) {
      this.isBootstrapped = true;

      return;
    }

    if (this.isBootstrapped || this.isLoading) return;

    try {
      const { flow_emailVariables } =
        await this.service.getFlowEmailVariables();

      runInAction(() => {
        this.load(flow_emailVariables);
      });
      runInAction(() => {
        this.isBootstrapped = true;
        this.totalElements = flow_emailVariables.length;
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

    try {
      const { flow_emailVariables } =
        await this.service.getFlowEmailVariables();

      runInAction(() => {
        this.load(flow_emailVariables);
      });
      runInAction(() => {
        this.isBootstrapped = true;
        this.totalElements = flow_emailVariables.length;
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
}
