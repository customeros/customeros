import { action } from 'mobx';
import { RootStore } from '@store/root.ts';
import { SettingsService } from '@domain/services/settings/settings.service';

export class ToggleCustomerBillingUsecase {
  private root = RootStore.getInstance();
  private service = new SettingsService();

  @action
  async execute(newStatus: boolean) {
    await this.service.toggleBilling(newStatus);
  }
}
