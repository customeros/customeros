import { ApolloError, NetworkStatus } from 'apollo-client';
import {
  GetOrganizationTimelineQuery,
  useGetOrganizationTimelineQuery,
} from './types';

interface Props {
  organizationId: string;
}

interface Result {
  data: //@ts-expect-error fixme
  | GetOrganizationTimelineQuery['organization']['timelineEvents']
    | null
    | undefined;
  loading: boolean;
  error: ApolloError | null;
  fetchMore: (data: { variables: any }) => void;
  variables: any;
  networkStatus?: NetworkStatus;
}

const x = new Date().toISOString();
export const useOrganizationTimeline = ({ organizationId }: Props): Result => {
  const { data, loading, error, refetch, variables, fetchMore, networkStatus } =
    useGetOrganizationTimelineQuery({
      variables: {
        organizationId,
        from: x,
        size: 10,
      },
      fetchPolicy: 'cache-first',
      notifyOnNetworkStatusChange: true,
    });

  const test = [...(data?.organization?.timelineEvents || [])].sort((a, b) => {
    return (
      //@ts-expect-error fixme
      Date.parse(a?.createdAt || a?.startedAt) -
      //@ts-expect-error fixme
      Date.parse(b?.createdAt || b?.startedAt)
    );
  });
  if (loading) {
    return {
      loading: true,
      error: null,
      data: test,
      fetchMore,
      variables: variables,
      networkStatus,
    };
  }

  if (error) {
    return {
      error,
      loading: false,
      variables: variables,
      networkStatus,
      data: null,
      fetchMore,
    };
  }

  return {
    data: test,
    fetchMore,
    loading,
    error: null,
    variables: variables,
    networkStatus,
  };
};
