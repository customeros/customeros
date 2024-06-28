import { Channel } from 'phoenix';
import { P, match } from 'ts-pattern';
import { RootStore } from '@store/root';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';

import { Filter, FilterItem } from '@shared/types/__generated__/graphql.types';

enum WorkflowType {
  'IDEAL_CUSTOMER_PROFILE',
  'IDEAL_CONTACT_PERSONA',
}

type WorkFlow = {
  id: string;
  name: string;
  condition: string;
  type: WorkflowType;
  status: 'running' | 'stopped';
};

export class WorkFlowStore implements Store<WorkFlow> {
  value: WorkFlow = defaultValue;
  version = 0;
  isLoading = false;
  error: string | null = null;
  history: Operation[] = [];
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  load = makeAutoSyncable.load<WorkFlow>();
  update = makeAutoSyncable.update<WorkFlow>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, { channelName: 'Workflow', mutator: this.save });
    makeAutoObservable(this);
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

  private async save() {}
}

const defaultValue: WorkFlow = {
  id: crypto.randomUUID(),
  name: '',
  condition: '',
  type: WorkflowType.IDEAL_CUSTOMER_PROFILE,
  status: 'stopped',
};
