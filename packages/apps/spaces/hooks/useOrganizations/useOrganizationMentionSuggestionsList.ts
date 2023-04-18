import { ApolloError } from 'apollo-client';
import {
  useGetOrganizationMentionSuggestionsLazyQuery,
  GetOrganizationMentionSuggestionsQueryVariables,
} from './types';

interface Result {
  loading: boolean;
  error: ApolloError | undefined;
  onLoadOrganizationMentionSuggestionsList: ({
    variables,
  }: {
    variables: GetOrganizationMentionSuggestionsQueryVariables;
  }) => Promise<any>;
}
export const useOrganizationMentionSuggestionsList = (): Result => {
  const [onLoadOrganizationMentionSuggestionsList, { data, loading, error }] =
    useGetOrganizationMentionSuggestionsLazyQuery();

  return {
    onLoadOrganizationMentionSuggestionsList,
    loading,
    error,
  };
};
