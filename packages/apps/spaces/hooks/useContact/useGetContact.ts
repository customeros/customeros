import { ApolloError } from '@apollo/client';
import { GetContactQuery, useGetContactQuery } from '@spaces/graphql';

interface Props {
  id: string;
}

export type ContactResponse = GetContactQuery['contact'] | undefined | null;
interface Result {
  data: ContactResponse;
  loading: boolean;
  error?: ApolloError | null;
}
export const useContact = ({ id }: Props): Result => {
  const { data, loading, error } = useGetContactQuery({
    variables: { id },
  });

  return {
    data: data?.contact,
    loading,
    error,
  };
};
