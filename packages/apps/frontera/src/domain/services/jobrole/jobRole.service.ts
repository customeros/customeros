import { RootStore } from '@store/root';
import { action, runInAction } from 'mobx';
import { JobRolesService as JobRoleRepo } from '@store/JobRoles/__service__/JobRoles.service';
import { SaveJobRolesMutationVariables } from '@store/JobRoles/__service__/saveJobRole.generated';

import { unwrap } from '@shared/util/unwrap';

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
  }

  @action
  async update(jobRole: SaveJobRolePayload) {
    try {
      await this.jobRoleRepo.saveJobRoles({
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
      this.root.contacts.retrieve([jobRole.id || '']);
    }
  }
}
