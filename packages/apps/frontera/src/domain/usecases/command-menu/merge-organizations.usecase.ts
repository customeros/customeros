import { action, observable } from 'mobx';
import { OrganizationsService } from '@store/Organizations/__service__/Organizations.service.ts';

export class MergeOrganizationsCase {
  @observable accessor primaryId: string = '';
  @observable accessor secondaryId: string = '';
  @observable accessor error: string = '';
  private service = OrganizationsService.getInstance();

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
      await this.service.mergeOrganizations({
        primaryOrganizationId: this.primaryId,
        mergedOrganizationIds: [this.secondaryId],
      });
    } catch (e) {
      this.error = `Merge failed`;
    } finally {
      this.primaryId = '';
      this.secondaryId = '';
    }
  }
}
