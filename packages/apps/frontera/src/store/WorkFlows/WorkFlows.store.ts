import type { RootStore } from '@store/root';

import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import { Workflow } from '@shared/types/__generated__/graphql.types';

import { WorkFlowStore } from './WorkFlow.store';
import { WorkFlowsService } from './WorkFLows.service';

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
  private service: WorkFlowsService;

  constructor(public root: RootStore, public transport: Transport) {
    this.service = WorkFlowsService.getInstance(transport);

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
      this.isLoading = true;

      const res =
        await this.transport.graphql.request<WORKFLOWS_QUERY_RESPONSE>(
          WORKFLOWS_QUERY,
        );

      this.load(res?.workflows);
      runInAction(() => {
        this.isBootstrapped = true;
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

  getByType(id: string) {
    return this.value.get(id);
  }

  toArray(): WorkFlowStore[] {
    return Array.from(this.value)?.flatMap(
      ([, WorkFlowStore]) => WorkFlowStore,
    );
  }
}

type WORKFLOWS_QUERY_RESPONSE = {
  workflows: Workflow[];
};

const WORKFLOWS_QUERY = gql`
  query workFlows {
    workflows {
      id
      name
      type
      live
      condition
    }
  }
`;
