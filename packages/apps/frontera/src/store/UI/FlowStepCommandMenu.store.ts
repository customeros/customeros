import { runInAction, makeAutoObservable } from 'mobx';

export type FlowStepCommandMenuType =
  // triggers
  | 'TriggersHub'
  | 'RecordAddedManually'
  | 'RecordCreated'
  | 'RecordUpdated'
  | 'RecordMatchesCondition'
  | 'Webhook'
  // steps
  | 'StepsHub';

export type Context = {
  id: string;
  callback?: () => void;
  property?: string | null;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  meta?: Record<string, any>;
  entity: 'Trigger' | 'Step' | null;
};

const makeDefaultContext = () => ({
  entity: null,
  property: null,
  id: '',
});

export class FlowStepCommandMenuStore {
  isOpen = false;
  type: FlowStepCommandMenuType = 'TriggersHub';
  context: Context = makeDefaultContext();

  constructor() {
    makeAutoObservable(this);
  }

  setOpen(
    open: boolean,
    options: {
      context: string | null;
      type: FlowStepCommandMenuType | null;
    } = {
      type: null,
      context: null,
    },
  ) {
    runInAction(() => {
      this.isOpen = open;
      this.type = options?.type ?? this.type;
    });
  }

  setType(type: FlowStepCommandMenuType) {
    runInAction(() => {
      this.type = type;
    });
  }

  toggle(type?: FlowStepCommandMenuType, context?: Context) {
    runInAction(() => {
      this.isOpen = !this.isOpen;
      this.type = type ?? 'StepsHub';

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
      if (this.context) this.context.id = '';
    });
  }

  reset() {
    runInAction(() => {
      this.isOpen = false;
      this.type = 'TriggersHub';
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
