import { when, makeAutoObservable } from 'mobx';
import { FlowSendersStore } from '@store/FlowSenders/FlowSenders.store.ts';

import type { Transport } from './transport';

import { Persister } from './persister';
import { UIStore } from './UI/UI.store';
import { MailStore } from './Mail/Mail.store';
import { TagsStore } from './Tags/Tags.store';
import { WindowManager } from './window-manager';
import { UsersStore } from './Users/Users.store';
import { FilesStore } from './Files/Files.store';
import { FlowsStore } from './Flows/Flows.store';
import { TransactionService } from './transaction';
import { SessionStore } from './Session/Session.store';
import { SettingsStore } from './Settings/Settings.store';
import { InvoicesStore } from './Invoices/Invoices.store';
import { ContactsStore } from './Contacts/Contacts.store';
import { ContractsStore } from './Contracts/Contracts.store';
import { RemindersStore } from './Reminders/Reminders.store';
import { WorkFlowsStore } from './WorkFlows/WorkFlows.store';
import { CustomFieldsStore } from './Settings/CustomFields.store';
import { GlobalCacheStore } from './GlobalCache/GlobalCache.store';
import { FlowContactsStore } from './FlowContacts/FlowContacts.store';
import { TableViewDefsStore } from './TableViewDefs/TableViewDefs.store';
import { OrganizationsStore } from './Organizations/Organizations.store';
import { OpportunitiesStore } from './Opportunities/Opportunities.store';
import { TimelineEventsStore } from './TimelineEvents/TimelineEvents.store';
import { ContractLineItemsStore } from './ContractLineItems/ContractLineItems.store';
import { FlowEmailVariablesStore } from './FlowEmailVariables/FlowEmailVariables.store';
import { ExternalSystemInstancesStore } from './ExternalSystemInstances/ExternalSystemInstances.store';

export class RootStore {
  demoMode = false;
  transactions: TransactionService;

  ui: UIStore;
  mail: MailStore;
  tags: TagsStore;
  files: FilesStore;
  users: UsersStore;
  flows: FlowsStore;
  session: SessionStore;
  settings: SettingsStore;
  invoices: InvoicesStore;
  contacts: ContactsStore;
  flowSenders: FlowSendersStore;
  contracts: ContractsStore;
  reminders: RemindersStore;
  workFlows: WorkFlowsStore;
  windowManager: WindowManager;
  globalCache: GlobalCacheStore;
  flowContacts: FlowContactsStore;
  customFields: CustomFieldsStore;
  tableViewDefs: TableViewDefsStore;
  organizations: OrganizationsStore;
  opportunities: OpportunitiesStore;
  timelineEvents: TimelineEventsStore;
  contractLineItems: ContractLineItemsStore;
  flowEmailVariables: FlowEmailVariablesStore;
  externalSystemInstances: ExternalSystemInstancesStore;

  constructor(private transport: Transport, demoMode: boolean = false) {
    makeAutoObservable(this);

    this.demoMode = demoMode;
    this.transactions = new TransactionService(this, this.transport);

    this.ui = new UIStore(this, this.transport);
    this.windowManager = new WindowManager(this);
    this.mail = new MailStore(this, this.transport);
    this.tags = new TagsStore(this, this.transport);
    this.files = new FilesStore(this, this.transport);
    this.users = new UsersStore(this, this.transport);
    this.flows = new FlowsStore(this, this.transport);
    this.flows = new FlowsStore(this, this.transport);
    this.session = new SessionStore(this, this.transport);
    this.settings = new SettingsStore(this, this.transport);
    this.invoices = new InvoicesStore(this, this.transport);
    this.contacts = new ContactsStore(this, this.transport);
    this.contracts = new ContractsStore(this, this.transport);
    this.reminders = new RemindersStore(this, this.transport);
    this.workFlows = new WorkFlowsStore(this, this.transport);
    this.customFields = new CustomFieldsStore(this, this.transport);
    this.globalCache = new GlobalCacheStore(this, this.transport);
    this.flowSenders = new FlowSendersStore(this, this.transport);
    this.flowContacts = new FlowContactsStore(this, this.transport);
    this.tableViewDefs = new TableViewDefsStore(this, this.transport);
    this.organizations = new OrganizationsStore(this, this.transport);
    this.opportunities = new OpportunitiesStore(this, this.transport);
    this.timelineEvents = new TimelineEventsStore(this, this.transport);
    this.contractLineItems = new ContractLineItemsStore(this, this.transport);
    this.flowEmailVariables = new FlowEmailVariablesStore(this, this.transport);

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
      () => this.isAuthenticated && !this.isHydrated,
      async () => {
        await Persister.attemptPurge();
        await this.bootstrap();
      },
    );

    when(
      () => this.isBootstrapped,
      () => this.transactions.startRunners(),
    );
  }

  async bootstrap() {
    await Promise.all([
      this.tableViewDefs.bootstrap(),
      this.globalCache.bootstrap(),
      this.settings.bootstrap(),
      this.customFields.bootstrap(),
      // this.organizations.bootstrapStream(),
      this.organizations.bootstrap(),
      this.tags.bootstrap(),
      this.opportunities.bootstrap(),
      this.invoices.bootstrap(),
      this.contracts.bootstrap(),
      this.externalSystemInstances.bootstrap(),
      this.users.bootstrap(),
      this.contacts.bootstrap(),
      this.workFlows.bootstrap(),
      this.flows.bootstrap(),
      this.flowEmailVariables.bootstrap(),
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

  get isHydrated() {
    if (this.demoMode) return true;

    return this.organizations.isHydrated;
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
      this.organizations.isBootstrapping ||
      this.tableViewDefs.isLoading ||
      this.settings.isBootstrapping ||
      this.globalCache.isLoading
    );
  }

  get isSyncing() {
    if (this.demoMode) return false;

    return this.organizations.isLoading;
  }
}
