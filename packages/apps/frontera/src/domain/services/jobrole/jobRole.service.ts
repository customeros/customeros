import { action } from 'mobx';
import { RootStore } from '@store/root';
import { JobRolesService as JobRoleRepo } from '@store/JobRoles/__service__/JobRoles.service';
import { SaveJobRolesMutationVariables } from '@store/JobRoles/__service__/saveJobRole.generated';

import { unwrap } from '@utils/unwrap';

type SaveJobRolePayload = SaveJobRolesMutationVariables['input'];

export class JobRoleService {
  private root = RootStore.getInstance();
  private jobRoleRepo = JobRoleRepo.getInstance();

  constructor() {}

  @action
  async create(jobRole: SaveJobRolePayload) {
    if (!jobRole) return;

    const [res, err] = await unwrap(
      this.jobRoleRepo.saveJobRoles({ input: jobRole }),
    );

    if (err) {
      console.error(err);

      return;
    }

    if (!res) {
      console.error('No response from saveJobRoles');
    }

    const serverId = res?.jobRole_Save;

    this.root.jobRoles.createNew({ id: serverId, ...jobRole });

    if (serverId) {
      this.root.contacts.getById(jobRole.contactId || '')?.addJobRole(serverId);
      jobRole.contactId && this.root.contacts.retrieve([jobRole.contactId]);
    }
  }

  @action
  async update(jobRole: SaveJobRolePayload) {
    const [res, err] = await unwrap(
      this.jobRoleRepo.saveJobRoles({ input: { ...jobRole, id: jobRole.id } }),
    );

    if (err) {
      console.error(err);

      return;
    }

    if (!res) {
      console.error('No response from saveJobRoles');
    } else {
      jobRole.contactId && this.root.contacts.retrieve([jobRole.contactId]);
    }
  }
}
