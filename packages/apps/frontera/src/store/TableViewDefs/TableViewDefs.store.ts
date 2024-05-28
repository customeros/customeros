import type { RootStore } from '@store/root';

import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import { TableIdType, type TableViewDef } from '@graphql/types';

import { TableViewDefStore } from './TableViewDef.store';

export class TableViewDefsStore implements GroupStore<TableViewDef> {
  value: Map<string, TableViewDefStore> = new Map();
  isLoading = false;
  channel?: Channel;
  version: number = 0;
  history: GroupOperation[] = [];
  isBootstrapped = false;
  error: string | null = null;
  sync = makeAutoSyncableGroup.sync;
  subscribe = makeAutoSyncableGroup.subscribe;
  load = makeAutoSyncableGroup.load<TableViewDef>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncableGroup(this, {
      channelName: `TableViewDefs:${this.root.session.value.tenant}`,
      ItemStore: TableViewDefStore,
      getItemId: (item) => item.id,
    });
    makeAutoObservable(this);
  }

  async bootstrap() {
    if (this.isBootstrapped) return;

    try {
      this.isLoading = true;
      const res =
        await this.transport.graphql.request<TABLE_VIEW_DEFS_QUERY_RESULT>(
          TABLE_VIEW_DEFS_QUERY,
        );

      this.load(res?.tableViewDefs);
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

  getById(id: string) {
    return this.value.get(id);
  }

  toArray(): TableViewDefStore[] {
    return Array.from(this.value)?.flatMap(
      ([, tableViewDefStore]) => tableViewDefStore,
    );
  }

  get defaultPreset() {
    const preset = this?.toArray().find(
      (t) => t.value.tableId === TableIdType.Organizations,
    )?.value.id;

    return preset;
  }
}

type TABLE_VIEW_DEFS_QUERY_RESULT = { tableViewDefs: TableViewDef[] };
const TABLE_VIEW_DEFS_QUERY = gql`
  query tableViewDefs {
    tableViewDefs {
      id
      name
      tableType
      tableId
      order
      icon
      filters
      sorting
      columns {
        columnType
        width
        visible
      }
      createdAt
      updatedAt
    }
  }
`;
