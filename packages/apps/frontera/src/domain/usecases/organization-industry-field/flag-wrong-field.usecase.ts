import { action } from 'mobx';
import { RootStore } from '@store/root';
import { OrganizationService } from '@domain/services';

import { FlagWrongFields } from '@graphql/types';

export class FlagWrongFieldUsecase {
  private root = RootStore.getInstance();
  private service = new OrganizationService();

  constructor() {}

  @action
  public async flagWrongField(id: string, field: FlagWrongFields) {
    try {
      await this.service.flagWrongField(id, field);
    } catch (err) {
      this.root.ui.toastError(
        `Something went wrong. Please try again later.`,
        `flag-field-error-${field}`,
      );
    }
  }
}
