import localforage from 'localforage';
import { configurePersistable } from 'mobx-persist-store';

import { TransportLayer } from './transport';
import { SessionStore } from './Session/Session.store';
import { TableViewDefsStore } from './TableViewDefs/TableViewDefs.store';

localforage.config({
  driver: localforage.INDEXEDDB,
  name: 'customerDB',
  version: 1.0,
  storeName: 'customer_os',
});

configurePersistable(
  {
    storage: localforage,
    expireIn: 1000 * 60 * 60 * 24, // 1 day
    version: 1.0,
  },
  {
    delay: 200,
    fireImmediately: false,
  },
);

export class RootStore {
  sessionStore: SessionStore;
  tableViewDefsStore: TableViewDefsStore;

  constructor(private transportLayer: TransportLayer) {
    this.sessionStore = new SessionStore();
    this.tableViewDefsStore = new TableViewDefsStore(this, transportLayer);
  }
}
