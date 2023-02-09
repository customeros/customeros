import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
import { Exact, Scalars } from './types';
const defaultOptions = {} as const;
export type GetContactDetailsQueryVariables = Exact<{
  email: Scalars['String'];
}>;

export type GetContactDetailsQuery = {
  __typename?: 'Query';
  contact_ByEmail: {
    __typename?: 'Contact';
    id: string;
    firstName?: string | null;
    lastName?: string | null;
    emails: Array<{ __typename?: 'Email'; email: string }>;
    phoneNumbers: Array<{ __typename?: 'PhoneNumber'; e164: string }>;
  };
};

export const GetContactDetailsDocument = gql`
  query GetContactDetails($email: String!) {
    contact_ByEmail(email: $email) {
      id
      firstName
      lastName
      emails {
        email
      }
      phoneNumbers {
        e164
      }
    }
  }
`;

/**
 * __useGetContactDetailsQuery__
 *
 * To run a query within a React component, call `useGetContactDetailsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetContactDetailsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetContactDetailsQuery({
 *   variables: {
 *      email: // value for 'email'
 *   },
 * });
 */
export function useGetContactDetailsQuery(
  baseOptions: Apollo.QueryHookOptions<
    GetContactDetailsQuery,
    GetContactDetailsQueryVariables
  >,
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useQuery<
    GetContactDetailsQuery,
    GetContactDetailsQueryVariables
  >(GetContactDetailsDocument, options);
}
export function useGetContactDetailsLazyQuery(
  baseOptions?: Apollo.LazyQueryHookOptions<
    GetContactDetailsQuery,
    GetContactDetailsQueryVariables
  >,
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useLazyQuery<
    GetContactDetailsQuery,
    GetContactDetailsQueryVariables
  >(GetContactDetailsDocument, options);
}
export type GetContactDetailsQueryHookResult = ReturnType<
  typeof useGetContactDetailsQuery
>;
export type GetContactDetailsLazyQueryHookResult = ReturnType<
  typeof useGetContactDetailsLazyQuery
>;
export type GetContactDetailsQueryResult = Apollo.QueryResult<
  GetContactDetailsQuery,
  GetContactDetailsQueryVariables
>;
