import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type ChangeFlowStatusMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
  status: Types.FlowStatus;
}>;

export type ChangeFlowStatusMutation = {
  __typename?: 'Mutation';
  flow_ChangeStatus: {
    __typename?: 'Flow';
    metadata: { __typename?: 'Metadata'; id: string };
  };
};
