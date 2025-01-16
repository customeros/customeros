import { Store } from '@store/_store';
import { RootStore } from '@store/root';
import { action, runInAction } from 'mobx';
import { Transport } from '@store/transport';

import { JobRole, type JobRoleDatum } from './JobRole.dto';
import { JobRolesService } from './__service__/JobRoles.service';

export class JobRolesStore extends Store<JobRoleDatum, JobRole> {
  private service = JobRolesService.getInstance();

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, {
      name: 'JobRoles',
      getId: (data) => data?.id,
      factory: JobRole,
    });
  }

  @action
  getjobTitleStoreByContactId(contactId: string) {
    return this.value.get(contactId)?.value.id;
  }

  @action
  getjobTitleByContactId(contactId: string) {
    return this.value.get(contactId)?.value.jobTitle;
  }

  @action
  async retrieveJobRoles(ids: string[]) {
    try {
      const { jobRoles } = await this.service.getJobRoles({ ids });

      runInAction(() => {
        jobRoles.forEach((jobRole) => {
          const foundJobRole = this.value.get(jobRole.id);

          if (foundJobRole) {
            Object.assign(foundJobRole.value, jobRole);
          } else {
            const newJobRole = new JobRole(this, jobRole);

            this.value.set(jobRole.id, newJobRole);
          }
        });
      });
      this.version++;
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    }
  }
}
