import { action } from 'mobx';
import { RootStore } from '@store/root.ts';
import { TenantService } from '@domain/services/tenant/tenant.service';

export class SwitchWorkspaceUsecase {
  private root = RootStore.getInstance();
  private service = new TenantService();

  @action
  execute(tenant: string) {
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
