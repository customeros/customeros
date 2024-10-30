import type { RootStore } from '@store/root';

import set from 'lodash/set';
import omit from 'lodash/omit';
import { Channel } from 'phoenix';
import { P, match } from 'ts-pattern';
import { gql } from 'graphql-request';
import debounce from 'lodash/debounce';
import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';
import { Store, makeAutoSyncable } from '@store/store';
import { Filter, Operation, FilterItem } from '@store/types';

import {
  TableIdType,
  TableViewDef,
  TableViewType,
  TableViewDefUpdateInput,
} from '@graphql/types';

export class TableViewDefStore implements Store<TableViewDef> {
  value: TableViewDef = getDefaultValue();
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  load = makeAutoSyncable.load<TableViewDef>();
  update = makeAutoSyncable.update<TableViewDef>();
  private readonly debouncedSave: () => void;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, { channelName: 'TableViewDef', mutator: this.save });
    makeAutoObservable(this);
    this.debouncedSave = debounce(this.save, 500);
  }

  set id(id: string) {
    this.value.id = id;
  }

  reorderColumn(sourceColumnId: number, targetColumnId: number) {
    this.update((value) => {
      const fromIndex = value.columns.findIndex(
        (c) => c.columnId === sourceColumnId,
      );
      const toIndex = value.columns.findIndex(
        (c) => c.columnId === targetColumnId,
      );
      const column = value.columns[fromIndex];

      value.columns.splice(fromIndex, 1);
      value.columns.splice(toIndex, 0, column);

      return value;
    });
  }

  orderColumnsByVisibility() {
    const prevLastVisibleIndex = [
      ...this.value.columns.map((c) => c.visible),
    ].lastIndexOf(true);

    const orderedColumns = this.value.columns.sort((a, b) => {
      if (a.visible === b.visible) return 0;
      if (a.visible) return -1;

      return 1;
    });

    const currentLastVisibleIndex = orderedColumns
      .map((c) => c.visible)
      .lastIndexOf(true);

    if (prevLastVisibleIndex === currentLastVisibleIndex) return;

    this.update((value) => {
      value.columns.sort((a, b) => {
        if (a.visible === b.visible) return 0;
        if (a.visible) return -1;

        return 1;
      });

      return value;
    });
  }

  setColumnName(columnId: number, name: string) {
    this.update(
      (value) => {
        const columnIdx = value.columns.findIndex(
          (c) => c.columnId === columnId,
        );

        value.columns[columnIdx].name = name;

        return value;
      },
      { mutate: false },
    );
  }

  setColumnSize(columnType: string, size: number) {
    runInAction(() => {
      const columnIdx = this.value.columns.findIndex(
        (c) => c.columnType === columnType,
      );

      if (columnIdx !== -1) {
        this.value.columns[columnIdx].width = size;
      }
    });

    this.debouncedSave();
  }

  async invalidate() {}

  async save() {
    const mutation = this.value.isShared
      ? UPDATE_TABLE_VIEW_DEF_SHARED
      : UPDATE_TABLE_VIEW_DEF;

    const payload: PAYLOAD = {
      input: omit(
        this.value,
        'updatedAt',
        'createdAt',
        'tableType',
        'tableId',
        'isPreset',
        'isShared',
        'defaultFilters',
      ),
    };

    try {
      this.isLoading = true;
      await this.transport.graphql.request(mutation, payload);
    } catch (e) {
      this.error = (e as Error)?.message;
    } finally {
      this.isLoading = false;
    }
  }

  getFilters() {
    try {
      return match(this.value.filters)
        .with(P.string.includes('AND'), (data) => JSON.parse(data))
        .otherwise(() => null);
    } catch (err) {
      console.error('Error parsing filters', err);

      return null;
    }
  }

  getDefaultFilters() {
    try {
      return match(this.value.defaultFilters)
        .with(P.string.includes('AND'), (data) => JSON.parse(data))
        .otherwise(() => null);
    } catch (err) {
      console.error('Error parsing default filters', err);

      return null;
    }
  }

  setDefaultFilters(filter: FilterItem) {
    this.update((value) => {
      const draft = this.getDefaultFilters();

      if (!draft) {
        this.appendDefaultFilter({ ...filter, active: true });

        return value;
      }
      const foundIndex = (draft.AND as Filter[])?.findIndex(
        (f) => f.filter?.property === filter.property,
      );

      if (foundIndex !== -1) {
        draft.AND[foundIndex].filter = filter;
        value.filters = JSON.stringify(draft);
      } else {
        this.appendDefaultFilter({ ...filter, active: true });
      }

      return value;
    });
  }

  hasFilters() {
    return this.getFilters()?.AND?.length > 0;
  }

  getFilter(id: string) {
    const filters = this.getFilters();

    return (filters?.AND as Filter[])?.find((f) => f.filter?.property === id)
      ?.filter;
  }

  appendFilter(filter: FilterItem) {
    this.update((value) => {
      let draft = this.getFilters() as Filter;

      if (draft) {
        (draft as Filter).AND?.push({ filter });
      } else {
        draft = { AND: [{ filter }] };
      }

      value.filters = JSON.stringify(draft);

      return value;
    });
  }

  appendDefaultFilter(filter: FilterItem) {
    this.update((value) => {
      let draft = this.getDefaultFilters() as Filter;

      if (draft) {
        (draft as Filter).AND?.push({ filter });
      } else {
        draft = { AND: [{ filter }] };
      }

      value.filters = JSON.stringify(draft);

      return value;
    });
  }

  removeFilter(id: string, index?: number) {
    this.update((value) => {
      const draft = this.getFilters();

      if (draft) {
        if (index !== undefined) {
          draft.AND.splice(index, 1);
        } else {
          draft.AND = (draft.AND as Filter[])?.filter(
            (f) => f.filter?.property !== id,
          );
        }
        value.filters = JSON.stringify(draft);
      }

      return value;
    });
  }

  removeFilters() {
    this.update((value) => {
      value.filters = JSON.stringify({ AND: [] });

      return value;
    });
  }

  toggleFilter(filter: FilterItem) {
    this.update((value) => {
      const draft = this.getFilters();

      if (draft) {
        const foundFilter = (draft.AND as Filter[])?.find(
          (f) => f.filter?.property === filter.property,
        )?.filter;

        if (foundFilter) {
          set(foundFilter, 'active', !filter?.active);
          value.filters = JSON.stringify(draft);
        } else {
          this.appendFilter({ ...filter, active: true });
        }
      }

      return value;
    });
  }

  setFilterv2(filter: FilterItem, index: number) {
    this.update((value) => {
      const draft = this.getFilters();

      if (!draft) {
        this.appendFilter({ ...filter, active: true });

        return value;
      }

      if (draft.AND && draft.AND[index]) {
        draft.AND[index].filter = filter;
      } else {
        draft.AND?.push({ filter });
      }

      value.filters = JSON.stringify(draft);

      return value;
    });
  }

  setFilter(filter: FilterItem) {
    this.update((value) => {
      const draft = this.getFilters();

      if (!draft) {
        this.appendFilter({ ...filter, active: true });

        return value;
      }
      const foundIndex = (draft.AND as Filter[])?.findIndex(
        (f) => f.filter?.property === filter.property,
      );

      if (foundIndex !== -1) {
        draft.AND[foundIndex].filter = filter;
        value.filters = JSON.stringify(draft);
      } else {
        this.appendFilter({ ...filter, active: true });
      }

      return value;
    });
  }

  setPropertyFilter(property: string) {
    this.update((value) => {
      const draft = this.getFilters();

      if (!draft) {
        this.appendFilter({
          property,
          active: false,
          value: undefined,
        });

        return value;
      }

      const foundIndex = (draft.AND as Filter[])?.findIndex(
        (f) => f.filter?.property === property,
      );

      if (foundIndex !== -1) {
        draft.AND[foundIndex].filter = { property, active: false };
        value.filters = JSON.stringify(draft);
      } else {
        this.appendFilter({
          property,
          active: false,
          value: undefined,
        });
      }

      return value;
    });
  }

  getPayloadToCopy = () => {
    return omit(this.value, 'id', 'createdAt', 'updatedAt');
  };
}

type PAYLOAD = { input: TableViewDefUpdateInput };
const UPDATE_TABLE_VIEW_DEF = gql`
  mutation updateTableViewDef($input: TableViewDefUpdateInput!) {
    tableViewDef_Update(input: $input) {
      id
    }
  }
`;

const UPDATE_TABLE_VIEW_DEF_SHARED = gql`
  mutation updateTableViewDefShared($input: TableViewDefUpdateInput!) {
    tableViewDef_UpdateShared(input: $input) {
      id
    }
  }
`;

export const getDefaultValue = (): TableViewDef => ({
  tableId: TableIdType.Organizations,
  columns: [],
  createdAt: '',
  filters: '',
  icon: '',
  id: '',
  name: '',
  order: 0,
  defaultFilters: '',
  sorting: '',
  updatedAt: '',
  isPreset: false,
  isShared: false,
  tableType: TableViewType.Organizations,
});
