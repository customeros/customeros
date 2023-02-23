import {
  GetContactCommunicationChannelsQuery,
  useGetContactCommunicationChannelsQuery,
} from './types';
import { ApolloError } from 'apollo-client';

interface Props {
  id: string;
}

interface Result {
  data: GetContactCommunicationChannelsQuery['contact'] | null | undefined;
  loading: boolean;
  error: ApolloError | null;
}
export const useContactCommunicationChannelsDetails = ({
  id,
}: Props): Result => {
  const { data, loading, error } = useGetContactCommunicationChannelsQuery({
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
    data: data?.contact,
    loading,
    error: null,
  };
};
