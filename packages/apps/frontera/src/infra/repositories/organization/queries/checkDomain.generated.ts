import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type CheckDomainQueryVariables = Types.Exact<{
  domain: Types.Scalars['String']['input'];
}>;

export type CheckDomainQuery = {
  __typename?: 'Query';
  checkDomain: {
    __typename?: 'DomainCheckDetails';
    domain: string;
    validSyntax: boolean;
    accessible: boolean;
    primary: boolean;
    primaryDomain: string;
    domainOrganizationId?: string | null;
    primaryDomainOrganizationId?: string | null;
    primaryDomainOrganizationName?: string | null;
    domainOrganizationName?: string | null;
    allowedForOrganization: boolean;
  };
};
