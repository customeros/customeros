import { Node } from '@xyflow/react';
import { runInAction, makeAutoObservable } from 'mobx';

export type FlowStepCommandMenuType = 'EmailAction';

export type Context = {
  id: string;
  node: Node | null;
  hasUnsavedChanges?: boolean;
};

const makeDefaultContext = () => ({
  id: '',
  node: null,
});

export class FlowActionSidePanelStore {
  isOpen = false;
  type: FlowStepCommandMenuType = 'EmailAction';
  context: Context = makeDefaultContext();

  constructor() {
    makeAutoObservable(this);
  }

  setOpen(
    open: boolean,
    options: {
      context: Context | null;
      type: FlowStepCommandMenuType | null;
    } = {
      type: null,
      context: null,
    },
  ) {
    runInAction(() => {
      this.isOpen = open;
      this.type = options?.type ?? this.type;
      this.context = options?.context ?? makeDefaultContext();
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
      this.type = type ?? 'EmailAction';

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

  reset() {
    runInAction(() => {
      this.isOpen = false;
      this.type = 'EmailAction';
      this.clearContext();
    });
  }

  clearContext() {
    runInAction(() => {
      Object.assign(this.context, makeDefaultContext());
    });
  }
}
