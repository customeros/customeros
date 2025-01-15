import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type GetJobRolesQueryVariables = Types.Exact<{
  ids: Array<Types.Scalars['ID']['input']> | Types.Scalars['ID']['input'];
}>;

export type GetJobRolesQuery = {
  __typename?: 'Query';
  jobRoles: Array<{
    __typename?: 'JobRole';
    id: string;
    createdAt: any;
    updatedAt: any;
    jobTitle?: string | null;
    primary: boolean;
    description?: string | null;
    company?: string | null;
    startedAt?: any | null;
    endedAt?: any | null;
    source: Types.DataSource;
    appSource: string;
    contact?: {
      __typename?: 'Contact';
      metadata: { __typename?: 'Metadata'; id: string };
    } | null;
  }>;
};
