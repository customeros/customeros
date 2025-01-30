import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type FlagWrongFieldMutationVariables = Types.Exact<{
  input: Types.FlagWrongFieldInput;
}>;

export type FlagWrongFieldMutation = {
  __typename?: 'Mutation';
  flagWrongField?: { __typename?: 'Result'; result: boolean } | null;
};
