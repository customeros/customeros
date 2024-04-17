import type { RootStore } from '@store/root';
import type { TransportLayer } from '@store/transport';

import { values, makeAutoObservable } from 'mobx';

import type { TableViewDef } from '@graphql/types';

import { TableViewDefStore } from './TableViewDef.store';

export class TableViewDefsStore {
  data: Map<string, TableViewDefStore> = new Map();

  constructor(
    private rootStore: RootStore,
    private transportLayer: TransportLayer,
  ) {
    makeAutoObservable(this);
  }

  load(tableViewDefs: TableViewDef[]) {
    tableViewDefs.forEach((tableViewDef) => {
      if (this.data.has(tableViewDef.id)) return;

      const tableViewDefStore = new TableViewDefStore(
        this.rootStore,
        this.transportLayer,
      );
      tableViewDefStore.load(tableViewDef);

      this.data.set(tableViewDef.id, tableViewDefStore);
    });
  }

  getById(id: string) {
    return this.data.get(id);
  }

  toArray() {
    return values(this.data);
  }
}
