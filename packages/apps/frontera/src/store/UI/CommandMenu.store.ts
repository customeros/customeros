import { runInAction, makeAutoObservable } from 'mobx';

export type CommandMenuType =
  | 'GlobalHub'
  | 'AssignOwner'
  | 'ChangeCurrency'
  | 'OpportunityHub'
  | 'OrganizationHub'
  | 'OrganizationCommands'
  | 'ChangeRelationship'
  | 'ChangeStage'
  | 'UpdateHealthStatus'
  | 'ChangeTags'
  | 'RenameOpportunityName'
  | 'ChangeArrEstimate'
  | 'OpportunityCommands'
  | 'AddContactViaLinkedInUrl'
  | 'RenameOrganizationProperty'
  | 'DeleteConfirmationModal'
  | 'OrganizationBulkCommands'
  | 'ChooseOpportunityOrganization'
  | 'ChooseOpportunityStage'
  | 'AddNewOrganization'
  | 'SetOpportunityNextSteps'
  | 'SetOpportunityNextSteps'
  | 'EditPersonaTag'
  | 'ContactHub'
  | 'ContactCommands'
  | 'EditEmail'
  | 'EditName'
  | 'EditPhoneNumber'
  | 'EditJobTitle'
  | 'ChangeOrAddJobRoles'
  | 'EditTimeZone';

type Context = {
  ids: Array<string>;
  property?: string | null;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  meta?: Record<string, any>;
  entity: 'Opportunity' | 'Organization' | 'Organizations' | 'Contact' | null;
};

const makeDefaultContext = () => ({
  entity: null,
  property: null,
  ids: [],
});

export class CommandMenuStore {
  isOpen = false;
  type: CommandMenuType = 'GlobalHub';
  context: Context = makeDefaultContext();

  constructor() {
    makeAutoObservable(this);
  }

  setOpen(
    open: boolean,
    options: { context: string | null; type: CommandMenuType | null } = {
      type: null,
      context: null,
    },
  ) {
    runInAction(() => {
      this.isOpen = open;
      this.type = options?.type ?? this.type;
    });
  }

  setType(type: CommandMenuType) {
    runInAction(() => {
      this.type = type;
    });
  }

  toggle(type?: CommandMenuType, context?: Context) {
    runInAction(() => {
      this.isOpen = !this.isOpen;
      this.type = type ?? 'GlobalHub';

      if (context) {
        Object.assign(this.context, context);
      }
    });
  }

  setContext(context: Context) {
    runInAction(() => {
      Object.assign(this.context, context);
    });
  }

  clearContextIds() {
    runInAction(() => {
      if (this.context) this.context.ids = [];
    });
  }

  reset() {
    runInAction(() => {
      this.isOpen = false;
      this.type = 'GlobalHub';
      this.clearContext();
    });
  }

  clearContext() {
    runInAction(() => {
      Object.assign(this.context, makeDefaultContext());
    });
  }
}
