import { Transport } from '@infra/transport';
import { when, makeAutoObservable } from 'mobx';
import { SkusStore } from '@store/Sku/Skus.store.ts';

import { Persister } from './persister';
import { UIStore } from './UI/UI.store';
import { MailStore } from './Mail/Mail.store';
import { TagsStore } from './Tags/Tags.store';
import { WindowManager } from './window-manager';
import { UsersStore } from './Users/Users.store';
import { FilesStore } from './Files/Files.store';
import { FlowsStore } from './Flows/Flows.store';
import { AgentStore } from './Agents/Agent.store';
import { TransactionService } from './transaction';
import { CommonStore } from './Common/Common.store';
import { SessionStore } from './Session/Session.store';
import { SettingsStore } from './Settings/Settings.store';
import { InvoicesStore } from './Invoices/Invoices.store';
import { ContactsStore } from './Contacts/Contacts.store';
import { MailboxesStore } from './Settings/Mailboxes.store';
import { ContractsStore } from './Contracts/Contracts.store';
import { RemindersStore } from './Reminders/Reminders.store';
import { JobRolesStore } from './JobRoles/JobRoles.store.ts';
import { IndustriesStore } from './Industries/Industries.store';
import { CustomFieldsStore } from './Settings/CustomFields.store';
import { GlobalCacheStore } from './GlobalCache/GlobalCache.store';
import { FlowSendersStore } from './FlowSenders/FlowSenders.store';
import { TableViewDefStore } from './TableViewDefs/TableViewDef.store';
import { OpportunitiesStore } from './Opportunities/Opportunities.store';
import { OrganizationsStore } from './Organizations/Organizations.store';
import { TimelineEventsStore } from './TimelineEvents/TimelineEvents.store';
import { FlowParticipantsStore } from './FlowParticipants/FlowParticipants.store';
import { ContractLineItemsStore } from './ContractLineItems/ContractLineItems.store';
import { FlowEmailVariablesStore } from './FlowEmailVariables/FlowEmailVariables.store';
import { ExternalSystemInstancesStore } from './ExternalSystemInstances/ExternalSystemInstances.store';

export class RootStore {
  demoMode = false;
  transactions: TransactionService;
  private transport = Transport.getInstance();

  ui: UIStore;
  mail: MailStore;
  tags: TagsStore;
  files: FilesStore;
  users: UsersStore;
  flows: FlowsStore;
  agents: AgentStore;
  common: CommonStore;
  settings: SettingsStore;
  session: SessionStore;
  invoices: InvoicesStore;
  contacts: ContactsStore;
  flowSenders: FlowSendersStore;
  contracts: ContractsStore;
  reminders: RemindersStore;
  windowManager: WindowManager;
  globalCache: GlobalCacheStore;
  flowParticipants: FlowParticipantsStore;
  customFields: CustomFieldsStore;
  tableViewDefs: TableViewDefStore;
  organizations: OrganizationsStore;
  industries: IndustriesStore;
  opportunities: OpportunitiesStore;
  timelineEvents: TimelineEventsStore;
  contractLineItems: ContractLineItemsStore;
  flowEmailVariables: FlowEmailVariablesStore;
  mailboxes: MailboxesStore;
  externalSystemInstances: ExternalSystemInstancesStore;
  jobRoles: JobRolesStore;
  skus: SkusStore;

  static instance: RootStore;

  constructor() {
    makeAutoObservable(this);

    this.transactions = new TransactionService(this, this.transport);

    this.common = new CommonStore(this);
    this.tableViewDefs = new TableViewDefStore(this, this.transport);
    this.settings = new SettingsStore(this, this.transport);
    this.ui = new UIStore(this, this.transport);
    this.windowManager = new WindowManager(this);
    this.mail = new MailStore(this, this.transport);
    this.tags = new TagsStore(this, this.transport);
    this.files = new FilesStore(this, this.transport);
    this.users = new UsersStore(this, this.transport);
    this.flows = new FlowsStore(this, this.transport);
    this.agents = new AgentStore(this, this.transport);
    this.session = new SessionStore(this, this.transport);
    this.industries = new IndustriesStore(this, this.transport);
    this.mailboxes = new MailboxesStore(this, this.transport);
    this.invoices = new InvoicesStore(this, this.transport);
    this.jobRoles = new JobRolesStore(this, this.transport);
    this.contacts = new ContactsStore(this, this.transport);
    this.contracts = new ContractsStore(this, this.transport);
    this.reminders = new RemindersStore(this, this.transport);
    this.customFields = new CustomFieldsStore(this, this.transport);
    this.globalCache = new GlobalCacheStore(this, this.transport);
    this.flowSenders = new FlowSendersStore(this, this.transport);
    this.flowParticipants = new FlowParticipantsStore(this, this.transport);
    this.organizations = new OrganizationsStore(this, this.transport);
    this.opportunities = new OpportunitiesStore(this, this.transport);
    this.timelineEvents = new TimelineEventsStore(this, this.transport);
    this.contractLineItems = new ContractLineItemsStore(this, this.transport);
    this.flowEmailVariables = new FlowEmailVariablesStore(this, this.transport);
    this.skus = new SkusStore(this, this.transport);
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

    when(
      () => this.isHydrated,
      async () => {
        await Persister.attemptPurge();
      },
    );

    this.transactions.startRunners();
  }

  async bootstrap() {
    await Promise.all([
      this.tableViewDefs.bootstrap(),
      this.globalCache.bootstrap(),
      this.industries.bootstrap(),
      this.settings.bootstrap(),
      this.customFields.bootstrap(),
      this.mailboxes.bootstrap(),
      this.tags.bootstrap(),
      this.opportunities.bootstrap(),
      this.invoices.bootstrap(),
      this.contracts.bootstrap(),
      this.externalSystemInstances.bootstrap(),
      this.users.bootstrap(),
      this.flows.bootstrap(),
      this.flowEmailVariables.bootstrap(),
      this.skus.bootstrap(),
      this.agents.bootstrap(),
      this.common.bootstrap(),
    ]);
  }

  public static getInstance() {
    if (!RootStore.instance) {
      RootStore.instance = new RootStore();
    }

    return RootStore.instance;
  }

  get isAuthenticating() {
    if (this.demoMode) return false;

    return this.session.isLoading !== null || this.session.isBootstrapping;
  }

  get isAuthenticated() {
    return Boolean(this.session.sessionToken);
  }

  get isHydrated() {
    return (
      this.organizations.isHydrated &&
      this.tableViewDefs.isHydrated &&
      this.contacts.isHydrated
    );
  }

  get isBootstrapped() {
    return (
      (this.tableViewDefs.isHydrated || this.tableViewDefs.isBootstrapped) &&
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

  get isSyncing() {
    if (this.demoMode) return false;

    return this.organizations.isLoading;
  }
}
