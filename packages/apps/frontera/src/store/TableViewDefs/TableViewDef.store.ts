import type { RootStore } from '@store/root';

import omit from 'lodash/omit';
import { gql } from 'graphql-request';
import { Meta } from '@store/abstract';
import { makeAutoObservable } from 'mobx';
import { AbstractStore } from '@store/abstract';
import { TransportLayer } from '@store/transport';

import {
  TableViewDef,
  TableViewType,
  TableViewDefUpdateInput,
} from '@graphql/types';

export class TableViewDefStore implements AbstractStore<TableViewDef> {
  value: TableViewDef = defaultValue;
  meta: Meta<TableViewDef>;

  constructor(
    private rootStore: RootStore,
    private transportLayer: TransportLayer,
  ) {
    this.meta = new Meta(this, this.transportLayer, {
      channelName: 'TableViewDef',
    });

    makeAutoObservable(this);
  }

  update(updater: (prev: TableViewDef) => TableViewDef) {
    this.meta.update(updater, this.save.bind(this));
  }
  async load(data: TableViewDef): Promise<void> {
    this.meta.load(data);
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
    this.update((value) => {
      value.columns.sort((a, b) => {
        if (a.visible === b.visible) return 0;
        if (a.visible) return -1;

        return 1;
      });

      return value;
    });
  }

  private async save() {
    const payload: PAYLOAD = {
      input: omit(this.value, 'updatedAt', 'createdAt', 'tableType'),
    };

    try {
      this.meta.isLoading = true;
      await this.transportLayer.client.request(UPDATE_TABLE_VIEW_DEF, payload);
    } catch (e) {
      this.meta.error = (e as Error)?.message;
    } finally {
      this.meta.isLoading = false;
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
