import type { RootStore } from '@store/root';

import { makeAutoObservable } from 'mobx';
import { AbstractStore } from '@store/abstract';
import { TransportLayer } from '@store/transport';
import { gql, GraphQLClient } from 'graphql-request';
import { GroupMeta, AbstractGroupStore } from '@store/abstract-group';

import type { TableViewDef } from '@graphql/types';

import { TableViewDefStore } from './TableViewDef.store';

export class TableViewDefsStore implements AbstractGroupStore<TableViewDef> {
  value: Map<string, TableViewDefStore> = new Map();
  isLoading = false;
  meta: GroupMeta<TableViewDef>;
  isBootstrapped = false;
  error: string | null = null;

  constructor(
    private rootStore: RootStore,
    private transportLayer: TransportLayer,
  ) {
    this.meta = new GroupMeta(this, this.rootStore, this.transportLayer);
    makeAutoObservable(this);

    this.bootstrap();
  }

  async bootstrap() {
    if (this.isBootstrapped || this.isLoading) return;

    try {
      this.isLoading = true;
      const res =
        await this.transportLayer.client.request<TABLE_VIEW_DEFS_QUERY_RESULT>(
          TABLE_VIEW_DEFS_QUERY,
        );

      this.load(res?.tableViewDefs);
    } catch (e) {
      this.error = (e as Error)?.message;
    } finally {
      this.isLoading = false;
    }
  }

  async load(tableViewDefs: TableViewDef[]) {
    this.meta.load(tableViewDefs, TableViewDefStore);
  }

  update(
    updater: (
      prev: Map<string, AbstractStore<TableViewDef>>,
    ) => Map<string, AbstractStore<TableViewDef>>,
  ): void {
    this.meta.update(updater);
  }

  getById(id: string) {
    return this.value.get(id);
  }
  toArray(): TableViewDefStore[] {
    return Array.from(this.value).flatMap(
      ([, tableViewDefStore]) => tableViewDefStore,
    );
  }

  static async serverSideBootstrap(client: GraphQLClient) {
    return client.request<TABLE_VIEW_DEFS_QUERY_RESULT>(TABLE_VIEW_DEFS_QUERY);
  }
}

type TABLE_VIEW_DEFS_QUERY_RESULT = { tableViewDefs: TableViewDef[] };
const TABLE_VIEW_DEFS_QUERY = gql`
  query tableViewDefs {
    tableViewDefs {
      id
      name
      tableType
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
