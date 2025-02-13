import { RootStore } from '@store/root';
import { action, observable } from 'mobx';

export class DisconnectSlackChannelUsecase {
  private root = RootStore.getInstance();
  @observable accessor isOpen = false;

  constructor() {
    this.toggleConfirmationDialog = this.toggleConfirmationDialog.bind(this);

    this.execute = this.execute.bind(this);
  }

  @action
  toggleConfirmationDialog(open: boolean) {
    this.isOpen = open;
  }

  execute() {
    this.root.settings.slack.disableSync();
    this.toggleConfirmationDialog(false);
  }
}
