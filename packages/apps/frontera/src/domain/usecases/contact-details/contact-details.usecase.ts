import { RootStore } from '@store/root';
import { ContactService } from '@domain/services';
import { action, computed, observable } from 'mobx';
import { JobRoleService } from '@domain/services/jobrole/jobRole.service';
import { SaveJobRolesMutationVariables } from '@store/JobRoles/__service__/saveJobRole.generated';
type SaveJobRolePayload = SaveJobRolesMutationVariables['input'];

export class ContactDetailsUseCase {
  private root = RootStore.getInstance();
  private jobRoleService = new JobRoleService();
  private contactService = new ContactService();
  private contactId: string;
  @observable accessor jobRole: string = '';

  constructor(contactId: string) {
    this.contactId = contactId;
    this.jobRole = this.root.jobRoles.getjobTitleByContactId(contactId) ?? '';
    this.setJobRole = this.setJobRole.bind(this);
  }

  @action
  setJobRole(jobRole: string) {
    this.jobRole = jobRole;
  }

  @computed
  get getJobRole() {
    return (
      this.jobRole ?? this.root.jobRoles.getjobTitleByContactId(this.contactId)
    );
  }

  async createJobRole(contactId: string, orgId: string) {
    await this.jobRoleService.create({
      jobTitle: this.getJobRole,
      contactId,
      primary: true,
      company: orgId,
    });
  }

  @action
  async updateJobRole(jobRole: SaveJobRolePayload) {
    await this.jobRoleService.update({
      primary: jobRole.primary ?? true,
      ...jobRole,
    });

    this.root.jobRoles.getById(jobRole.id!)!.value.jobTitle = jobRole.jobTitle;
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
      await this.jobRoleService.create({
        jobTitle: this.getJobRole,
        contactId,
        primary: true,
        company: orgId,
      });
    }
  }

  @action
  async changeContactName() {
    const contact = this.root.contacts.getById(this.contactId);

    if (!contact) return;
    await this.contactService.changeContactName(contact, contact.name);
  }
}
