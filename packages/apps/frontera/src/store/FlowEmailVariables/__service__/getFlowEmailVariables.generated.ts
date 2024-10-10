import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type GetFlowEmailVariablesQueryVariables = Types.Exact<{
  [key: string]: never;
}>;

export type GetFlowEmailVariablesQuery = {
  __typename?: 'Query';
  flow_emailVariables: Array<{
    __typename?: 'EmailVariableEntity';
    type: Types.EmailVariableEntityType;
    variables: Array<Types.EmailVariableName>;
  }>;
};
