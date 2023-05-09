import { ApolloError } from '@apollo/client';
import { useGetContactMentionSuggestionsLazyQuery } from './types';
import { GetContactMentionSuggestionsQueryVariables } from '../../graphQL/__generated__/generated';

interface Result {
  loading: boolean;
  error: ApolloError | undefined;
  onLoadContactMentionSuggestionsList: ({
    variables,
  }: {
    variables: GetContactMentionSuggestionsQueryVariables;
  }) => Promise<any>;
}
export const useContactMentionSuggestionsList = (): Result => {
  const [onLoadContactMentionSuggestionsList, { data, loading, error }] =
    useGetContactMentionSuggestionsLazyQuery();

  return {
    onLoadContactMentionSuggestionsList,
    loading,
    error,
  };
};
