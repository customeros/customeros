import { ApolloError } from 'apollo-client';
import { GetUserByEmailQuery, useGetUserByEmailQuery } from './types';

interface Props {
  email: string;
}

interface Result {
  data: GetUserByEmailQuery['user_ByEmail'] | undefined | null;
  loading: boolean;
  error: ApolloError | null;
}
export const useUser = ({ email }: Props): Result => {
  console.log('🏷️ ----- email: ', email);
  const { data, loading, error } = useGetUserByEmailQuery({
    variables: {
      email,
    },
  });

  return {
    data: data?.user_ByEmail,
    loading,
    error: null,
  };
};
