import localforage from 'localforage';
import { when, makeAutoObservable } from 'mobx';
import { configurePersistable } from 'mobx-persist-store';
import { InvoicesStore } from '@store/Invoices/Invoices.store.ts';
import { ContractLineItemsStore } from '@store/ContractLineItems/ContractLineItems.store.ts';
import { ExternalSystemInstancesStore } from '@store/ExternalSystemInstances/ExternalSystemInstances.store.ts';

import { UIStore } from './UI/UI.store';
import { Transport } from './transport';
import { MailStore } from './Mail/Mail.store.ts';
import { TagsStore } from './Tags/Tags.store.ts';
import { UsersStore } from './Users/Users.store.ts';
import { FilesStore } from './Files/Files.store.ts';
import { SessionStore } from './Session/Session.store';
import { SettingsStore } from './Settings/Settings.store';
import { ContactsStore } from './Contacts/Contacts.store.ts';
import { ContractsStore } from './Contracts/Contracts.store.ts';
import { RemindersStore } from './Reminders/Reminders.store.ts';
import { GlobalCacheStore } from './GlobalCache/GlobalCache.store';
import { TableViewDefsStore } from './TableViewDefs/TableViewDefs.store';
import { OrganizationsStore } from './Organizations/Organizations.store.ts';
import { OpportunitiesStore } from './Opportunities/Opportunities.store.ts';
import { TimelineEventsStore } from './TimelineEvents/TimelineEvents.store.ts';

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
  demoMode = false;
  ui: UIStore;
  mail: MailStore;
  files: FilesStore;
  users: UsersStore;
  session: SessionStore;
  settings: SettingsStore;
  invoices: InvoicesStore;
  contacts: ContactsStore;
  contracts: ContractsStore;
  reminders: RemindersStore;
  globalCache: GlobalCacheStore;
  tableViewDefs: TableViewDefsStore;
  organizations: OrganizationsStore;
  opportunities: OpportunitiesStore;
  timelineEvents: TimelineEventsStore;
  contractLineItems: ContractLineItemsStore;
  tags: TagsStore;
  externalSystemInstances: ExternalSystemInstancesStore;

  constructor(private transport: Transport, demoMode: boolean = false) {
    makeAutoObservable(this);
    this.demoMode = demoMode;

    this.ui = new UIStore();
    this.mail = new MailStore(this, this.transport);
    this.files = new FilesStore(this, this.transport);
    this.users = new UsersStore(this, this.transport);
    this.tags = new TagsStore(this, this.transport);
    this.session = new SessionStore(this, this.transport);
    this.settings = new SettingsStore(this, this.transport);
    this.invoices = new InvoicesStore(this, this.transport);
    this.contacts = new ContactsStore(this, this.transport);
    this.contracts = new ContractsStore(this, this.transport);
    this.reminders = new RemindersStore(this, this.transport);
    this.globalCache = new GlobalCacheStore(this, this.transport);
    this.tableViewDefs = new TableViewDefsStore(this, this.transport);
    this.organizations = new OrganizationsStore(this, this.transport);
    this.opportunities = new OpportunitiesStore(this, this.transport);
    this.timelineEvents = new TimelineEventsStore(this, this.transport);
    this.contractLineItems = new ContractLineItemsStore(this, this.transport);

    this.externalSystemInstances = new ExternalSystemInstancesStore(
      this,
      this.transport,
    );

    when(
      () => this.demoMode,
      () => {
        console.info('Demo mode enabled');
      },
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
      this.tags.bootstrap(),
      this.opportunities.bootstrap(),
      this.invoices.bootstrap(),
      this.contracts.bootstrap(),
      this.externalSystemInstances.bootstrap(),
      this.users.bootstrap(),
      this.contacts.bootstrap(),
    ]);
  }

  get isAuthenticating() {
    if (this.demoMode) return false;

    return this.session.isLoading !== null || this.session.isBootstrapping;
  }
  get isAuthenticated() {
    if (this.demoMode) return true;

    return Boolean(this.session.sessionToken);
  }
  get isBootstrapped() {
    if (this.demoMode) return true;

    return (
      this.tableViewDefs.isBootstrapped &&
      this.settings.isBootstrapped &&
      this.globalCache.isBootstrapped
    );
  }

  get isBootstrapping() {
    if (this.demoMode) return false;

    return (
      this.tableViewDefs.isLoading ||
      this.settings.isBootstrapping ||
      this.globalCache.isLoading
    );
  }
}
