import merge from 'lodash/merge';
import { Entity } from '@store/record';
import { computed, observable, runInAction } from 'mobx';

import { DataSource } from '@shared/types/__generated__/graphql.types';

import { JobRolesStore } from './JobRoles.store';
import { GetJobRolesQuery } from './__service__/getJobRoles.generated';
import { SaveJobRolesMutationVariables } from './__service__/saveJobRole.generated';

export type JobRoleDatum = NonNullable<GetJobRolesQuery['jobRoles'][number]>;

export class JobRole extends Entity<JobRoleDatum> {
  @observable accessor value: JobRoleDatum = JobRole.default();

  constructor(store: JobRolesStore, data: JobRoleDatum) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    super(store as any, data);
  }

  @computed
  get id() {
    return this.value.id;
  }

  @computed
  get jobTitle() {
    return this.value.jobTitle;
  }

  set id(value: string) {
    runInAction(() => {
      this.value.id = value;
    });
  }

  @computed
  get primary() {
    return this.value.primary;
  }

  static default(
    payload?: JobRoleDatum | SaveJobRolesMutationVariables['input'],
  ): JobRoleDatum {
    return merge(
      {
        id: crypto.randomUUID(),
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        jobTitle: '',
        primary: false,
        description: '',
        company: '',
        startedAt: null,
        endedAt: null,
        source: DataSource.Openline,
        appSource: '',
        contact: null,
      },
      payload ?? {},
    );
  }
}
