import { ApolloError } from '@apollo/client';
import {
  useGetHealthIndicatorsQuery,
} from '@spaces/graphql';

interface Result {
  data: Array<{ label: string; value: string }> | undefined | null;
  loading: boolean;
  error: ApolloError | undefined;
}
export const useHealthIndicators = (): Result => {
  const { data, loading, error } = useGetHealthIndicatorsQuery();

  return {
    data: (data?.healthIndicators || []).map(({ id, name }) => ({
      label: name,
      value: id,
    })),
    loading,
    error,
  };
};
