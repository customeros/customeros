import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type RemindersQueryVariables = Types.Exact<{
  organizationId: Types.Scalars['ID']['input'];
}>;


export type RemindersQuery = { __typename?: 'Query', remindersForOrganization: Array<{ __typename?: 'Reminder', content?: string | null, dueDate?: any | null, dismissed?: boolean | null, metadata: { __typename?: 'Metadata', id: string, created: any, lastUpdated: any }, owner?: { __typename?: 'User', id: string, firstName: string, lastName: string, name?: string | null } | null }> };
