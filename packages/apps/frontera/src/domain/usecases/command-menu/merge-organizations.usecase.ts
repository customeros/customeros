import { RootStore } from '@store/root.ts';
import { action, observable, runInAction } from 'mobx';
import { OrganizationsService } from '@store/Organizations/__service__/Organizations.service.ts';

export class MergeOrganizationsCase {
  @observable accessor primaryId: string = '';
  @observable accessor secondaryId: string = '';
  @observable accessor error: string = '';
  private service = OrganizationsService.getInstance();
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
      const { organization_Merge } = await this.service.mergeOrganizations({
        primaryOrganizationId: this.primaryId,
        mergedOrganizationIds: [this.secondaryId],
      });

      runInAction(() => {
        this.root.organizations.sync({
          action: 'DELETE',
          ids: [this.secondaryId],
        });
        this.root.organizations.sync({
          action: 'INVALIDATE',
          ids: [this.primaryId],
        });

        if (organization_Merge.id) {
          this.root.ui.toastSuccess(`Merged organizations`, this.primaryId);
        }
      });
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
