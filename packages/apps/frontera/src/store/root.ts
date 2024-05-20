import localforage from 'localforage';
import { when, makeAutoObservable } from 'mobx';
import { configurePersistable } from 'mobx-persist-store';
import { NewBusinessTableStore } from '@store/Organizations/NewBusinessTable.store.ts';

import { UIStore } from './UI/UI.store';
import { Transport } from './transport';
import { MailStore } from './Mail/Mail.store.ts';
import { FilesStore } from './Files/Files.store.ts';
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
  ui: UIStore;
  mail: MailStore;
  files: FilesStore;
  session: SessionStore;
  settings: SettingsStore;
  globalCache: GlobalCacheStore;
  tableViewDefs: TableViewDefsStore;
  newBusiness: NewBusinessTableStore;

  constructor(private transport: Transport) {
    makeAutoObservable(this);

    this.ui = new UIStore();
    this.mail = new MailStore(this, this.transport);
    this.files = new FilesStore(this, this.transport);
    this.session = new SessionStore(this, this.transport);
    this.settings = new SettingsStore(this, this.transport);
    this.globalCache = new GlobalCacheStore(this, this.transport);
    this.tableViewDefs = new TableViewDefsStore(this, this.transport);
    this.newBusiness = new NewBusinessTableStore(this, this.transport);

    when(
      () => this.isAuthenticated,
      async () => {
        await this.bootstrap();
      },
    );
  }

  async bootstrap() {
    await Promise.all([
      this.globalCache.bootstrap(),
      this.settings.bootstrap(),
      this.tableViewDefs.bootstrap(),
    ]);
  }

  get isAuthenticating() {
    return this.session.isLoading !== null || this.session.isBootstrapping;
  }
  get isAuthenticated() {
    return Boolean(this.session.sessionToken);
  }
  get isBootstrapped() {
    return (
      this.tableViewDefs.isBootstrapped &&
      this.settings.isBootstrapped &&
      this.globalCache.isBootstrapped
    );
  }

  get isBootstrapping() {
    return (
      this.tableViewDefs.isLoading ||
      this.settings.isBootstrapping ||
      this.globalCache.isLoading
    );
  }
}
