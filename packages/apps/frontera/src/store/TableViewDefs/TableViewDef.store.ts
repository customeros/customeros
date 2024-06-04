import type { RootStore } from '@store/root';

import omit from 'lodash/omit';
import { Channel } from 'phoenix';
import { P, match } from 'ts-pattern';
import { gql } from 'graphql-request';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';

import {
  TableIdType,
  TableViewDef,
  TableViewType,
  TableViewDefUpdateInput,
} from '@graphql/types';

export class TableViewDefStore implements Store<TableViewDef> {
  value: TableViewDef = defaultValue;
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  load = makeAutoSyncable.load<TableViewDef>();
  update = makeAutoSyncable.update<TableViewDef>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, { channelName: 'TableViewDef', mutator: this.save });
    makeAutoObservable(this);
  }
  set id(id: string) {
    this.value.id = id;
  }
  reorderColumn(fromIndex: number, toIndex: number) {
    this.update((value) => {
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

  async invalidate() {}

  private async save() {
    const payload: PAYLOAD = {
      input: omit(this.value, 'updatedAt', 'createdAt', 'tableType', 'tableId'),
    };

    try {
      this.isLoading = true;
      await this.transport.graphql.request(UPDATE_TABLE_VIEW_DEF, payload);
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
}

type PAYLOAD = { input: TableViewDefUpdateInput };
const UPDATE_TABLE_VIEW_DEF = gql`
  mutation updateTableViewDef($input: TableViewDefUpdateInput!) {
    tableViewDef_Update(input: $input) {
      id
    }
  }
`;

const defaultValue: TableViewDef = {
  tableId: TableIdType.Organizations,
  columns: [],
  createdAt: '',
  filters: '',
  icon: '',
  id: '',
  name: '',
  order: 0,
  sorting: '',
  updatedAt: '',
  tableType: TableViewType.Organizations,
};
