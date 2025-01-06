import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type GetLogEntryQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;


export type GetLogEntryQuery = { __typename?: 'Query', logEntry: { __typename?: 'LogEntry', id: string, content?: string | null, contentType?: string | null, createdAt: any, updatedAt: any, tags: Array<{ __typename?: 'Tag', name: string, metadata: { __typename?: 'Metadata', id: string } }>, createdBy?: { __typename?: 'User', id: string, firstName: string, lastName: string, name?: string | null } | null } };
