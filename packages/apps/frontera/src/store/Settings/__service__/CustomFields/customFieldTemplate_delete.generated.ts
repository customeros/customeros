import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type DeleteCustomFieldTemplateMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;

export type DeleteCustomFieldTemplateMutation = {
  __typename?: 'Mutation';
  customFieldTemplate_Delete?: boolean | null;
};
