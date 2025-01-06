import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type DeleteTagMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;


export type DeleteTagMutation = { __typename?: 'Mutation', tag_Delete?: { __typename?: 'Result', result: boolean } | null };
