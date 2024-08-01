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
  | 'ChangeArrEstimate'
  | 'OpportunityCommands'
  | 'AddContactViaLinkedInUrl'
  | 'RenameOrganizationProperty';

type Context = {
  id: string | null;
  property?: string | null;
  entity: 'Opportunity' | 'Organization' | null;
};

const makeDefaultContext = () => ({ id: null, entity: null, property: null });

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

  clearContextId() {
    runInAction(() => {
      if (this.context) this.context.id = null;
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
