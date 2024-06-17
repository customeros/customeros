import localforage from 'localforage';
import { when, makeAutoObservable } from 'mobx';
import { configurePersistable } from 'mobx-persist-store';
import { InvoicesStore } from '@store/Invoices/Invoices.store.ts';
import { ContractLineItemsStore } from '@store/Organizations/ContractLineItems.store.ts';
import { ExternalSystemInstancesStore } from '@store/ExternalSystemInstances/ExternalSystemInstances.store.ts';

import { UIStore } from './UI/UI.store';
import { Transport } from './transport';
import { MailStore } from './Mail/Mail.store.ts';
import { UsersStore } from './Users/Users.store.ts';
import { FilesStore } from './Files/Files.store.ts';
import { SessionStore } from './Session/Session.store';
import { SettingsStore } from './Settings/Settings.store';
import { GlobalCacheStore } from './GlobalCache/GlobalCache.store';
import { ContractsStore } from './Organizations/Contracts.store.ts';
import { TableViewDefsStore } from './TableViewDefs/TableViewDefs.store';
import { OrganizationsStore } from './Organizations/Organizations.store.ts';
import { OpportunitiesStore } from './Opportunities/Opportunities.store.ts';

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
  users: UsersStore;
  session: SessionStore;
  settings: SettingsStore;
  contracts: ContractsStore;
  contractLineItems: ContractLineItemsStore;
  globalCache: GlobalCacheStore;
  tableViewDefs: TableViewDefsStore;
  organizations: OrganizationsStore;
  opportunities: OpportunitiesStore;
  invoices: InvoicesStore;
  externalSystemInstances: ExternalSystemInstancesStore;

  constructor(private transport: Transport) {
    makeAutoObservable(this);

    this.ui = new UIStore();
    this.mail = new MailStore(this, this.transport);
    this.files = new FilesStore(this, this.transport);
    this.users = new UsersStore(this, this.transport);
    this.session = new SessionStore(this, this.transport);
    this.settings = new SettingsStore(this, this.transport);
    this.contracts = new ContractsStore(this, this.transport);
    this.globalCache = new GlobalCacheStore(this, this.transport);
    this.tableViewDefs = new TableViewDefsStore(this, this.transport);
    this.organizations = new OrganizationsStore(this, this.transport);
    this.opportunities = new OpportunitiesStore(this, this.transport);
    this.contractLineItems = new ContractLineItemsStore(this, this.transport);
    this.invoices = new InvoicesStore(this, this.transport);
    this.opportunities = new OpportunitiesStore(this, this.transport);

    this.externalSystemInstances = new ExternalSystemInstancesStore(
      this,
      this.transport,
    );

    when(
      () => this.isAuthenticated,
      async () => {
        await this.bootstrap();
      },
    );
  }

  async bootstrap() {
    await Promise.all([
      this.tableViewDefs.bootstrap(),
      this.globalCache.bootstrap(),
      this.settings.bootstrap(),
      this.organizations.bootstrap(),

      this.opportunities.bootstrap(),
      this.invoices.bootstrap(),
      this.contracts.bootstrap(),
      this.externalSystemInstances.bootstrap(),

      this.users.bootstrap(),
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
