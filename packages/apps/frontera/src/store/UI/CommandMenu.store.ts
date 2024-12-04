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
  | 'RenameOrganizationProperty'
  | 'DeleteConfirmationModal'
  | 'OrganizationBulkCommands'
  | 'ChooseOpportunityOrganization'
  | 'ChooseOpportunityStage'
  | 'AddNewOrganization'
  | 'SetOpportunityNextSteps'
  | 'MergeConfirmationModal'
  | 'EditPersonaTag'
  | 'ContactHub'
  | 'ContactCommands'
  | 'EditEmail'
  | 'EditName'
  | 'EditPhoneNumber'
  | 'EditJobTitle'
  | 'ChangeOrAddJobRoles'
  | 'EditTimeZone'
  | 'RenameTableViewDef'
  | 'ContactEmailVerificationInfoModal'
  | 'DuplicateView'
  | 'OpportunityBulkCommands'
  | 'ChangeBulkArrEstimate'
  | 'CreateNewFlow'
  | 'DuplicateFlow'
  | 'RenameFlow'
  | 'FlowCommands'
  | 'ChangeFlowStatus'
  | 'FlowsBulkCommands'
  | 'FlowHub'
  | 'StartFlow'
  | 'StopFlow'
  | 'EditContactFlow'
  | 'AddContactsToFlow'
  | 'ActiveFlowUpdateInfo'
  | 'UnlinkContactFromFlow'
  | 'ConfirmBulkFlowEdit'
  | 'ConfirmSingleFlowEdit'
  | 'GetBrowserExtensionLink'
  | 'InstallLinkedInExtension'
  | 'FlowValidationMessage'
  | 'ConfirmEmailContentChanges'
  | 'AddContactsBulk'
  | 'ContactBulkCommands';

export type Context = {
  ids: Array<string>;
  callback?: () => void;
  selectId?: string | null;
  property?: string | null;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  meta?: Record<string, any>;
  entity:
    | 'Opportunity'
    | 'Organization'
    | 'Organizations'
    | 'Opportunities'
    | 'Contact'
    | 'ContactFlows'
    | 'TableViewDef'
    | 'Flow'
    | 'Flows'
    | null;
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
    options: { context: Context | null; type: CommandMenuType | null } = {
      type: null,
      context: null,
    },
  ) {
    runInAction(() => {
      this.isOpen = open;
      this.type = options?.type ?? this.type;
      this.context = options?.context ?? this.context;
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

  setCallback(callback: () => void) {
    runInAction(() => {
      this.context.callback = callback;
    });
  }

  clearCallback() {
    runInAction(() => {
      this.context.callback = undefined;
    });
  }
}
