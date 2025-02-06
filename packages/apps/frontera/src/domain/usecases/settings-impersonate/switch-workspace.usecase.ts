import { action } from 'mobx';
import { RootStore } from '@store/root.ts';
import { CommonService } from '@domain/services/common/common.service';

export class SwitchWorkspaceUsecase {
  private root = RootStore.getInstance();
  private service = new CommonService();

  @action
  execute(tenant: string | undefined) {
    if (!tenant) {
      this.root.ui.toastError(
        'Please select a workspace',
        'switch-workspace-failed',
      );

      return;
    }

    return this.service.switchWorkspace(tenant);
  }
}
