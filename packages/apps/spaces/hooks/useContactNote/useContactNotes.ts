import {
  GetContactCommunicationChannelsQuery,
  useGetContactCommunicationChannelsQuery,
} from './types';
import { ApolloError } from 'apollo-client';
import {GetContactNotesQuery, useGetContactNotesQuery} from "../../graphQL/__generated__/generated";

interface Props {
  id: string;
}

interface Result {
  data: GetContactNotesQuery['contact'] | null | undefined;
  loading: boolean;
  error: ApolloError | null;
}
export const useContactNotes = ({
  id,
}: Props): Result => {
  const { data, loading, error } = useGetContactNotesQuery({
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

  console.log('data loaded for notes')
  return {
    data: data?.contact ?? null,
    loading,
    error: null,
  };
};
