// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../../src/types/__generated__/graphql.types';

import type { InfiniteData } from '@tanstack/react-query';
import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import {
  useQuery,
  useInfiniteQuery,
  UseQueryOptions,
  UseInfiniteQueryOptions,
  InfiniteData,
} from '@tanstack/react-query';

function fetcher<TData, TVariables extends { [key: string]: any }>(
  client: GraphQLClient,
  query: string,
  variables?: TVariables,
  requestHeaders?: RequestInit['headers'],
) {
  return async (): Promise<TData> =>
    client.request({
      document: query,
      variables,
      requestHeaders,
    });
}
export type OrganizationPeoplePanelQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;

export type OrganizationPeoplePanelQuery = {
  __typename?: 'Query';
  organization?: {
    __typename?: 'Organization';
    id: string;
    name: string;
    contacts: {
      __typename?: 'ContactsPage';
      totalElements: any;
      content: Array<{
        __typename?: 'Contact';
        id: string;
        name?: string | null;
        firstName?: string | null;
        lastName?: string | null;
        prefix?: string | null;
        description?: string | null;
        timezone?: string | null;
        profilePhotoUrl?: string | null;
        jobRoles: Array<{
          __typename?: 'JobRole';
          id: string;
          primary: boolean;
          jobTitle?: string | null;
          description?: string | null;
          company?: string | null;
          startedAt?: any | null;
        }>;
        phoneNumbers: Array<{
          __typename?: 'PhoneNumber';
          id: string;
          e164?: string | null;
          rawPhoneNumber?: string | null;
          label?: Types.PhoneNumberLabel | null;
          primary: boolean;
        }>;
        emails: Array<{
          __typename?: 'Email';
          id: string;
          email?: string | null;
          emailValidationDetails: {
            __typename?: 'EmailValidationDetails';
            isReachable?: string | null;
            isValidSyntax?: boolean | null;
            canConnectSmtp?: boolean | null;
            acceptsMail?: boolean | null;
            hasFullInbox?: boolean | null;
            isCatchAll?: boolean | null;
            isDeliverable?: boolean | null;
            validated?: boolean | null;
            isDisabled?: boolean | null;
          };
        }>;
        socials: Array<{
          __typename?: 'Social';
          id: string;
          platformName?: string | null;
          url: string;
        }>;
      }>;
    };
  } | null;
};

export const OrganizationPeoplePanelDocument = `
    query OrganizationPeoplePanel($id: ID!) {
  organization(id: $id) {
    id
    name
    contacts(pagination: {page: 0, limit: 100}) {
      content {
        id
        name
        firstName
        lastName
        prefix
        description
        timezone
        jobRoles {
          id
          primary
          jobTitle
          description
          company
          startedAt
        }
        phoneNumbers {
          id
          e164
          rawPhoneNumber
          label
          primary
        }
        emails {
          id
          email
          emailValidationDetails {
            isReachable
            isValidSyntax
            canConnectSmtp
            acceptsMail
            hasFullInbox
            isCatchAll
            isDeliverable
            validated
            isDisabled
          }
        }
        socials {
          id
          platformName
          url
        }
        profilePhotoUrl
      }
      totalElements
    }
  }
}
    `;

export const useOrganizationPeoplePanelQuery = <
  TData = OrganizationPeoplePanelQuery,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: OrganizationPeoplePanelQueryVariables,
  options?: Omit<
    UseQueryOptions<OrganizationPeoplePanelQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseQueryOptions<
      OrganizationPeoplePanelQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useQuery<OrganizationPeoplePanelQuery, TError, TData>({
    queryKey: ['OrganizationPeoplePanel', variables],
    queryFn: fetcher<
      OrganizationPeoplePanelQuery,
      OrganizationPeoplePanelQueryVariables
    >(client, OrganizationPeoplePanelDocument, variables, headers),
    ...options,
  });
};

useOrganizationPeoplePanelQuery.document = OrganizationPeoplePanelDocument;

useOrganizationPeoplePanelQuery.getKey = (
  variables: OrganizationPeoplePanelQueryVariables,
) => ['OrganizationPeoplePanel', variables];

export const useInfiniteOrganizationPeoplePanelQuery = <
  TData = InfiniteData<OrganizationPeoplePanelQuery>,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: OrganizationPeoplePanelQueryVariables,
  options: Omit<
    UseInfiniteQueryOptions<OrganizationPeoplePanelQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseInfiniteQueryOptions<
      OrganizationPeoplePanelQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useInfiniteQuery<OrganizationPeoplePanelQuery, TError, TData>(
    (() => {
      const { queryKey: optionsQueryKey, ...restOptions } = options;
      return {
        queryKey: optionsQueryKey ?? [
          'OrganizationPeoplePanel.infinite',
          variables,
        ],
        queryFn: (metaData) =>
          fetcher<
            OrganizationPeoplePanelQuery,
            OrganizationPeoplePanelQueryVariables
          >(
            client,
            OrganizationPeoplePanelDocument,
            { ...variables, ...(metaData.pageParam ?? {}) },
            headers,
          )(),
        ...restOptions,
      };
    })(),
  );
};

useInfiniteOrganizationPeoplePanelQuery.getKey = (
  variables: OrganizationPeoplePanelQueryVariables,
) => ['OrganizationPeoplePanel.infinite', variables];

useOrganizationPeoplePanelQuery.fetcher = (
  client: GraphQLClient,
  variables: OrganizationPeoplePanelQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<OrganizationPeoplePanelQuery, OrganizationPeoplePanelQueryVariables>(
    client,
    OrganizationPeoplePanelDocument,
    variables,
    headers,
  );

useOrganizationPeoplePanelQuery.mutateCacheEntry =
  (
    queryClient: QueryClient,
    variables: OrganizationPeoplePanelQueryVariables,
  ) =>
  (
    mutator: (
      cacheEntry: OrganizationPeoplePanelQuery,
    ) => OrganizationPeoplePanelQuery,
  ) => {
    const cacheKey = useOrganizationPeoplePanelQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<OrganizationPeoplePanelQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<OrganizationPeoplePanelQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteOrganizationPeoplePanelQuery.mutateCacheEntry =
  (
    queryClient: QueryClient,
    variables: OrganizationPeoplePanelQueryVariables,
  ) =>
  (
    mutator: (
      cacheEntry: InfiniteData<OrganizationPeoplePanelQuery>,
    ) => InfiniteData<OrganizationPeoplePanelQuery>,
  ) => {
    const cacheKey = useInfiniteOrganizationPeoplePanelQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<OrganizationPeoplePanelQuery>>(
        cacheKey,
      );
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<OrganizationPeoplePanelQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
