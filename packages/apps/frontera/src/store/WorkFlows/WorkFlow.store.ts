import { Channel } from 'phoenix';
import { P, match } from 'ts-pattern';
import { RootStore } from '@store/root';
import { Operation } from '@store/types';
import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';
import { Store, makeAutoSyncable } from '@store/store';
import { makeAutoSyncableGroup } from '@store/group-store';

import {
  Filter,
  Workflow,
  FilterItem,
  WorkflowType,
} from '@shared/types/__generated__/graphql.types';

import { WorkFlowService } from './WorkFlow.service';

export class WorkFlowStore implements Store<Workflow> {
  value: Workflow = defaultValue;
  version = 0;
  isLoading = false;
  error: string | null = null;
  history: Operation[] = [];
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncable.load<Workflow>();
  update = makeAutoSyncable.update<Workflow>();
  private service: WorkFlowService;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, { channelName: 'Workflow', mutator: this.save });
    makeAutoObservable(this);
    this.service = WorkFlowService.getInstance(transport);
  }
  async invalidate() {}

  set id(id: string) {
    this.value.id = id;
  }

  getFilters() {
    try {
      return match(this.value.condition)
        .with(P.string.includes('AND'), (data) => JSON.parse(data))
        .otherwise(() => null);
    } catch (err) {
      console.error('Error parsing filters', err);

      return null;
    }
  }

  getFilter(id: string) {
    const filters = this.getFilters();

    return (filters?.AND as Filter[])?.find((f) => f.filter?.property === id)
      ?.filter;
  }

  appendFilter(filter: FilterItem) {
    this.update((value) => {
      let draft = this.getFilters() as Filter;

      if (
        draft &&
        draft?.AND?.findIndex((f) => f.filter?.property === filter.property) !==
          -1
      ) {
        return value;
      }

      if (draft) {
        (draft as Filter).AND?.push({ filter });
      } else {
        draft = { AND: [{ filter }] };
      }

      value.condition = JSON.stringify(draft);

      return value;
    });
  }

  removeFilter(id: string) {
    this.update((value) => {
      const draft = this.getFilters();
      if (draft) {
        draft.AND = (draft.AND as Filter[])?.filter(
          (f) => f.filter?.property !== id,
        );
        value.condition = JSON.stringify(draft);
      }

      return value;
    });
  }

  setFilter(filter: FilterItem) {
    this.update((value) => {
      const draft = this.getFilters();
      value.live = false;

      if (!draft) {
        this.appendFilter({ ...filter });

        return value;
      }

      const foundIndex = (draft.AND as Filter[])?.findIndex(
        (f) => f.filter?.property === filter.property,
      );

      if (foundIndex !== -1) {
        draft.AND[foundIndex].filter = filter;
        value.condition = JSON.stringify(draft);
      } else {
        this.appendFilter({ ...filter });
      }

      return value;
    });
  }

  async getWorfklowByType() {
    try {
      await this.service.getWorkFlowsByType();
    } catch (e) {
      this.error = (e as Error).message;
    }
  }

  async updateWorkflow(value?: boolean) {
    const payload = {
      id: this.value.id,
      name: this.value.name,
      condition: this.value.condition,
      live: value,
    };
    try {
      await this.service.updateWorkFlow(payload);
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    }
  }

  private async save(operation: Operation) {
    const diff = operation.diff?.[0];
    const type = diff?.op;
    const path = diff?.path;
    const value = diff?.val;

    match(path)
      .with(['live', ...P.array()], () => {
        if (type === 'update') {
          this.updateWorkflow(value);
        }
      })
      .with(['condition', ...P.array()], () => {
        if (type === 'update') {
          this.updateWorkflow(false);
        }
      });
  }
}

const defaultValue: Workflow = {
  condition: '',
  name: '',
  id: crypto.randomUUID(),
  live: false,
  type: WorkflowType.IdealCustomerProfile,
};
