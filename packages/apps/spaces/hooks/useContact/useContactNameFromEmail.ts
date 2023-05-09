import { ApolloError } from '@apollo/client';
import {
  GetContactNameByEmailQuery,
  useGetContactNameByEmailQuery,
} from './types';

interface Props {
  email: string;
}

interface Result {
  data: GetContactNameByEmailQuery['contact_ByEmail'] | undefined | null;
  loading: boolean;
  error: ApolloError | null;
}
export const useContactNameFromEmail = ({ email }: Props): Result => {
  const { data, loading, error } = useGetContactNameByEmailQuery({
    variables: { email },
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
    data: data?.contact_ByEmail,
    loading,
    error: null,
  };
};
