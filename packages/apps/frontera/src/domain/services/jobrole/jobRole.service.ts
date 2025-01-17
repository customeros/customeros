import { RootStore } from '@store/root';
import { action, runInAction } from 'mobx';
import { JobRole } from '@store/JobRoles/JobRole.dto';
import { JobRolesService } from '@store/JobRoles/__service__/JobRoles.service';
import { SaveJobRolesMutationVariables } from '@store/JobRoles/__service__/saveJobRole.generated';

type SaveJobRolePayload = SaveJobRolesMutationVariables['input'];

export class JobRoleService {
  private root = RootStore.getInstance();
  private service = JobRolesService.getInstance();

  constructor() {}

  @action
  async create(jobRole: SaveJobRolePayload) {
    let tempId = '';
    const contactStore = this.root.contacts.getById(jobRole.contactId || '');

    if (!contactStore) return;

    try {
      const draft = new JobRole(this.root.jobRoles, JobRole.default(jobRole));

      this.root.jobRoles.value.set(draft.id, draft);

      const { jobRole_Save } = await this.service.saveJobRoles({
        input: {
          jobTitle: jobRole.jobTitle,
          contactId: jobRole.contactId,
          primary: jobRole.primary,
          organizationId: jobRole.organizationId,
        },
      });

      runInAction(() => {
        draft.id = jobRole_Save;
        this.root.jobRoles.value.set(draft.id, draft);
        this.root.jobRoles.value.delete(tempId);

        tempId = draft.id;

        this.root.jobRoles.sync({
          action: 'APPEND',
          ids: [draft.id],
        });
      });
      this.root.jobRoles.version++;
      this.root.contacts.retrieve([jobRole.contactId || '']);
    } catch (e) {
      runInAction(() => {
        this.root.jobRoles.value.delete(tempId);
        this.root.jobRoles.error = (e as Error).message;
      });
    } finally {
      contactStore?.draft();
      contactStore.value.primaryOrganizationJobRoleTitle = jobRole.jobTitle;
      contactStore?.commit({ syncOnly: true });
    }
  }

  @action
  async update(jobRole: SaveJobRolePayload) {
    try {
      await this.service.saveJobRoles({
        input: {
          ...jobRole,
          id: jobRole.id,
        },
      });
    } catch (e) {
      runInAction(() => {
        this.root.jobRoles.error = (e as Error).message;
      });
    } finally {
      this.root.jobRoles.version++;
      this.root.contacts.retrieve([jobRole.contactId || '']);
    }
  }
}
