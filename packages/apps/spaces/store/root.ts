import { makeAutoObservable } from 'mobx';

import { TableViewDef } from '@graphql/types';

import { TransportLayer } from './transport';

export class RootStore {
  tableViewDefStore: TableViewDefStore;

  constructor(private transportLayer: TransportLayer) {
    this.tableViewDefStore = new TableViewDefStore(this, transportLayer);
  }
}

class TableViewDefStore {
  data: Map<string, TableViewDef> = new Map();

  constructor(
    private rootStore: RootStore,
    private transportLayer: TransportLayer,
  ) {
    makeAutoObservable(this);
    transportLayer?.join('test');
  }

  load(tableViewDefs: TableViewDef[]) {
    tableViewDefs.forEach((tableViewDef) => {
      this.data.set(tableViewDef.id, tableViewDef);
    });
  }

  getTableViewDef(id: string) {
    return this.data.get(id);
  }

  get all() {
    return Array.from(this.data.values());
  }
}
