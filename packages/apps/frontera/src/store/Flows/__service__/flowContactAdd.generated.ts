import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type FlowContactAddMutationVariables = Types.Exact<{
  flowId: Types.Scalars['ID']['input'];
  contactId: Types.Scalars['ID']['input'];
}>;

export type FlowContactAddMutation = {
  __typename?: 'Mutation';
  flowContact_Add: {
    __typename?: 'FlowContact';
    metadata: { __typename?: 'Metadata'; id: string };
  };
};
