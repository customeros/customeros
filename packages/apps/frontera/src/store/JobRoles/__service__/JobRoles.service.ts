import { Transport } from '@infra/transport';

import GetJobRolesDocument from './getJobRoles.graphql';
import SaveJobRolesDocument from './saveJobRole.graphql';
import {
  GetJobRolesQuery,
  GetJobRolesQueryVariables,
} from './getJobRoles.generated';
import {
  SaveJobRolesMutation,
  SaveJobRolesMutationVariables,
} from './saveJobRole.generated';

export class JobRolesService {
  private static instance: JobRolesService | null = null;
  private transport = Transport.getInstance();

  constructor() {}

  public static getInstance(): JobRolesService {
    if (!JobRolesService.instance) {
      JobRolesService.instance = new JobRolesService();
    }

    return JobRolesService.instance;
  }

  async getJobRoles(payload: GetJobRolesQueryVariables) {
    return this.transport.graphql.request<
      GetJobRolesQuery,
      GetJobRolesQueryVariables
    >(GetJobRolesDocument, payload);
  }

  async saveJobRoles(payload: SaveJobRolesMutationVariables) {
    return this.transport.graphql.request<
      SaveJobRolesMutation,
      SaveJobRolesMutationVariables
    >(SaveJobRolesDocument, payload);
  }
}
