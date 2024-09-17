import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type FlowChangeStatusMutationVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
  status: Types.FlowStatus;
}>;

export type FlowChangeStatusMutation = {
  __typename?: 'Mutation';
  flow_ChangeStatus: {
    __typename?: 'Flow';
    metadata: { __typename?: 'Metadata'; id: string };
  };
};
