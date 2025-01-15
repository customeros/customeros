import { observable } from 'mobx';
import { RootStore } from '@store/root';
import { JobRoleService } from '@domain/services/jobrole/jobRole.service';

export class AddJobRole {
  private root = RootStore.getInstance();

  private jobRoleService = new JobRoleService();
  @observable accessor jobRole: string = '';

  constructor() {
    this.setJobRole = this.setJobRole.bind(this);
  }

  setJobRole(jobRole: string) {
    this.jobRole = jobRole;
  }

  async createJobRole(contactId: string) {
    const contactStore = this.root.contacts.getById(contactId);

    if (!contactStore) {
      throw new Error('Contact not found');
    }

    try {
      await this.jobRoleService.create({
        jobTitle: this.jobRole,
        contactId: contactId,
      });
    } catch (e) {
      throw new Error(e instanceof Error ? e.message : String(e));
    } finally {
      contactStore?.draft();
      contactStore.value.primaryOrganizationJobRoleTitle = this.jobRole;
      contactStore?.commit({ syncOnly: true });
    }
  }
}
