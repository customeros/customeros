import { RootStore } from '@store/root';
import { action, observable } from 'mobx';
import { JobRoleService } from '@domain/services/jobrole/jobRole.service';
import { SaveJobRolesMutationVariables } from '@store/JobRoles/__service__/saveJobRole.generated';
type SaveJobRolePayload = SaveJobRolesMutationVariables['input'];

export class EditJobRole {
  private root = RootStore.getInstance();

  private jobRoleService = new JobRoleService();
  @observable accessor jobRole: string = '';

  constructor() {
    this.setJobRole = this.setJobRole.bind(this);
  }

  setJobRole(jobRole: string) {
    this.jobRole = jobRole;
  }

  async createJobRole(contactId: string, orgId: string) {
    try {
      await this.jobRoleService.create({
        jobTitle: this.jobRole,
        contactId: contactId,
        primary: true,
        organizationId: orgId,
      });
    } catch (e) {
      throw new Error(e instanceof Error ? e.message : String(e));
    }
  }

  @action
  async updateJobRole(jobRole: SaveJobRolePayload) {
    try {
      await this.jobRoleService.update({
        ...jobRole,
        primary: jobRole.primary,
      });
    } catch (e) {
      throw new Error(e instanceof Error ? e.message : String(e));
    }
  }

  @action
  async submitJobRole(contactId: string, orgId: string) {
    const jobTitle = this.root.jobRoles.getjobTitleByContactId(contactId);

    if (jobTitle) {
      await this.updateJobRole({ jobTitle: this.jobRole });
    } else {
      await this.createJobRole(contactId, orgId);
    }
    this.jobRole = '';
  }
}
