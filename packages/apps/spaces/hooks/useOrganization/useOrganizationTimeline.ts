import {
  GetOrganizationDetailsQuery,
  useGetOrganizationDetailsQuery,
} from './types';
import { ApolloError } from 'apollo-client';
import {
  LoadTimelineForOrganizationQuery,
  useLoadTimelineForOrganizationQuery,
} from '../../graphQL/__generated__/generated';

interface Props {
  id: string;
}

interface Result {
  data: LoadTimelineForOrganizationQuery['organization'] | undefined | null;
  loading: boolean;
  error: ApolloError | null;
}
export const useOrganizationTimelineData = ({ id }: Props): Result => {
  const { data, loading, error } = useLoadTimelineForOrganizationQuery({
    variables: { id },
  });

  if (loading) {
    return {
      loading: true,
      error: null,
      data: null,
    };
  }

  if (error) {
    return {
      error,
      loading: false,
      data: null,
    };
  }

  return {
    data: data?.organization,
    loading,
    error: null,
  };
};
