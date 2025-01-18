import { RootStore } from '@store/root';
import { action, computed, observable } from 'mobx';
import { JobRoleService } from '@domain/services/jobrole/jobRole.service';
import { SaveJobRolesMutationVariables } from '@store/JobRoles/__service__/saveJobRole.generated';
type SaveJobRolePayload = SaveJobRolesMutationVariables['input'];

export class EditJobRole {
  private root = RootStore.getInstance();
  private jobRoleService = new JobRoleService();
  private contactId: string;
  @observable accessor jobRole: string | undefined = undefined;

  constructor(contactId: string) {
    this.contactId = contactId;
    this.setJobRole = this.setJobRole.bind(this);
  }

  @action
  setJobRole(jobRole: string | undefined) {
    this.jobRole = jobRole;
  }

  @computed
  get getJobRole() {
    return (
      this.jobRole ?? this.root.jobRoles.getjobTitleByContactId(this.contactId)
    );
  }

  async createJobRole(contactId: string, orgId: string) {
    try {
      await this.jobRoleService.create({
        jobTitle: this.getJobRole,
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
    const jobId = Array.from(this.root.jobRoles.value.values()).find(
      (job) => job.value.contact?.metadata.id === contactId,
    )?.id;

    if (!this.jobRole) return;

    if (jobTitle) {
      await this.updateJobRole({
        jobTitle: this.getJobRole,
        contactId,
        id: jobId,
      });
    } else {
      await this.createJobRole(contactId, orgId);
    }
  }
}
