import { ApolloError } from 'apollo-client';
import {
  useGetOrganizationMentionSuggestionsLazyQuery,
  GetOrganizationMentionSuggestionsQueryVariables,
} from './types';
import {
  ComparisonOperator,
} from '../../graphQL/__generated__/generated';
import { toast } from 'react-toastify';

interface Result {
  loading: boolean;
  error: ApolloError | undefined;
  onLoadOrganizationSuggestionsList: ({
    variables,
  }: {
    variables: GetOrganizationMentionSuggestionsQueryVariables;
  }) => Promise<any>;
  getOrganizationSuggestions: (
    filter: string,
  ) => Promise<Array<{ label: string; value: string }>>;
}
export const useOrganizationSuggestionsList = (): Result => {
  const [onLoadOrganizationSuggestionsList, { loading, error }] =
    useGetOrganizationMentionSuggestionsLazyQuery();

  const getOrganizationSuggestions: Result['getOrganizationSuggestions'] =
    async (filter) => {
      try {
        const response = await onLoadOrganizationSuggestionsList({
          variables: {
            pagination: { page: 0, limit: 10 },
            where: {
              filter: {
                property: 'NAME',
                value: filter,
                operation: ComparisonOperator.Contains,
              },
            },
          },
        });
        if (response?.data) {
          return (response.data?.organizations?.content || []).map((e) => ({
            label: e.id,
            value: e.id,
          }));
        }
        return [];
      } catch (e) {
        toast.error(
          'Something went wrong while loading organization suggestion list',
        );
        return [];
      }
    };

  return {
    getOrganizationSuggestions,
    onLoadOrganizationSuggestionsList,
    loading,
    error,
  };
};
