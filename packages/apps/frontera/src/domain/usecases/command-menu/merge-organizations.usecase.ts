import { RootStore } from '@store/root.ts';
import { action, observable, runInAction } from 'mobx';
import { OrganizationService } from '@domain/services';

export class MergeOrganizationsCase {
  @observable accessor primaryId: string = '';
  @observable accessor secondaryId: string = '';
  @observable accessor error: string = '';
  private service = new OrganizationService();

  private root = RootStore.getInstance();

  constructor() {
    this.setIds = this.setIds.bind(this);
  }

  @action
  setIds(primary: string, secondary: string) {
    this.primaryId = primary;
    this.secondaryId = secondary;
  }

  @action
  async merge() {
    try {
      await this.service.merge(this.primaryId, [this.secondaryId]);
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
        this.root.ui.toastError(
          `We were unable to merge this organization`,
          this.primaryId,
        );
      });
    } finally {
      this.primaryId = '';
      this.secondaryId = '';
    }
  }
}
