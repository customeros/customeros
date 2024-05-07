import localforage from 'localforage';
import { makeAutoObservable } from 'mobx';
import { configurePersistable } from 'mobx-persist-store';

import { UIStore } from './UI/UI.store';
import { TransportLayer } from './transport';
import { SessionStore } from './Session/Session.store';
import { SettingsStore } from './Settings/Settings.store';
import { GlobalCacheStore } from './GlobalCache/GlobalCache.store';
import { TableViewDefsStore } from './TableViewDefs/TableViewDefs.store';

localforage.config({
  driver: localforage.INDEXEDDB,
  name: 'customerDB',
  version: 1.0,
  storeName: 'customer_os',
});

configurePersistable({
  storage: localforage,
  expireIn: 1000 * 60 * 60 * 24, // 1 day
  version: 1.0,
  stringify: false,
});

export class RootStore {
  uiStore: UIStore;
  sessionStore: SessionStore;
  settingsStore: SettingsStore;
  globalCacheStore: GlobalCacheStore;
  tableViewDefsStore: TableViewDefsStore;

  constructor(private transportLayer: TransportLayer) {
    makeAutoObservable(this);

    this.uiStore = new UIStore();
    this.sessionStore = new SessionStore(this, this.transportLayer);
    this.settingsStore = new SettingsStore(this, this.transportLayer);
    this.globalCacheStore = new GlobalCacheStore(this, this.transportLayer);
    this.tableViewDefsStore = new TableViewDefsStore(this, this.transportLayer);
  }

  get isBootstrapped() {
    return (
      this.tableViewDefsStore.isBootstrapped &&
      this.settingsStore.isBootstrapped &&
      this.sessionStore.isBootstrapped &&
      this.globalCacheStore.isBootstrapped
    );
  }
}
