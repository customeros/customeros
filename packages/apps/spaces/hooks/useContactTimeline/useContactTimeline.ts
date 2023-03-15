import { ApolloError, NetworkStatus } from 'apollo-client';
import { GetContactTimelineQuery, useGetContactTimelineQuery } from './types';

interface Props {
  contactId: string;
}

interface Result {
  //@ts-expect-error fixme
  data: GetContactTimelineQuery['contact']['timelineEvents'] | null | undefined;
  loading: boolean;
  error: ApolloError | null;
  fetchMore: (data: { variables: any }) => void;
  variables: any;
  networkStatus?: NetworkStatus;
}

const DATE_NOW = new Date().toISOString();
export const useContactTimeline = ({ contactId }: Props): Result => {
  const { data, loading, error, refetch, variables, fetchMore, networkStatus } =
    useGetContactTimelineQuery({
      variables: {
        contactId,
        from: DATE_NOW,
        size: 10,
      },
      fetchPolicy: 'cache-first',
      notifyOnNetworkStatusChange: true,
    });

  const test = [...(data?.contact?.timelineEvents || [])].sort((a, b) => {
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
