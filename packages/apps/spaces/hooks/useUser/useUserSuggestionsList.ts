import { ApolloError } from '@apollo/client';
import {
  ComparisonOperator,
  GetUsersQueryVariables,
  useGetUsersLazyQuery,
} from '../../graphQL/__generated__/generated';
import { toast } from 'react-toastify';

interface Result {
  loading: boolean;
  error: ApolloError | undefined;
  onLoadUsersSuggestionsList: ({
    variables,
  }: {
    variables: GetUsersQueryVariables;
  }) => Promise<any>;
  getUsersSuggestions: (
    filter: string,
  ) => Promise<Array<{ label: string; value: string }>>;
}

export const useUserSuggestionsList = (): Result => {
  const [onLoadUsersSuggestionsList, { loading, error }] =
    useGetUsersLazyQuery();

  const getCondition = (filter = '') => {
    const splittedFilter = filter.split(' ');

    if (splittedFilter.length === 1) {
      return {
        OR: [
          {
            filter: {
              property: 'FIRST_NAME',
              value: filter,
              operation: ComparisonOperator.Contains,
            },
          },
          {
            filter: {
              property: 'LAST_NAME',
              value: filter,
              operation: ComparisonOperator.Contains,
            },
          },
        ],
      };
    }

    return {
      OR: [
        {
          AND: [
            {
              filter: {
                property: 'FIRST_NAME',
                value: splittedFilter[1],
                operation: ComparisonOperator.Contains,
              },
            },
            {
              filter: {
                property: 'LAST_NAME',
                value: splittedFilter[0],
                operation: ComparisonOperator.Contains,
              },
            },
          ],
        },
        {
          AND: [
            {
              filter: {
                property: 'FIRST_NAME',
                value: splittedFilter[0],
                operation: ComparisonOperator.Contains,
              },
            },
            {
              filter: {
                property: 'LAST_NAME',
                value: splittedFilter[1],
                operation: ComparisonOperator.Contains,
              },
            },
          ],
        },
      ],
    };
  };

  const getUsersSuggestions: Result['getUsersSuggestions'] = async (filter) => {
    try {
      const response = await onLoadUsersSuggestionsList({
        variables: {
          pagination: { page: 0, limit: 10 },
          where: {
            ...getCondition(filter),
          },
        },
      });
      if (response?.data) {
        return (response.data?.users?.content || []).map((e) => ({
          label: e.firstName + ' ' + e.lastName,
          value: e.id,
        }));
      }
      return [];
    } catch (e) {
      toast.error('Something went wrong while loading users suggestion list');
      return [];
    }
  };

  return {
    getUsersSuggestions,
    onLoadUsersSuggestionsList,
    loading,
    error,
  };
};
