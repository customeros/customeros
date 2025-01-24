import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type RemoveDomainsMutationVariables = Types.Exact<{
  organizationId: Types.Scalars['ID']['input'];
  domains?: Types.InputMaybe<
    Array<Types.Scalars['String']['input']> | Types.Scalars['String']['input']
  >;
}>;

export type RemoveDomainsMutation = {
  __typename?: 'Mutation';
  organization_RemoveDomains: {
    __typename?: 'ActionResponse';
    accepted: boolean;
  };
};
