import { ApolloError } from 'apollo-client';
import {
  GetContactNameByPhoneNumberQuery,
  useGetContactNameByPhoneNumberQuery,
} from './types';

interface Props {
  phoneNumber: string;
}

interface Result {
  data: GetContactNameByPhoneNumberQuery['contact_ByPhone'] | undefined | null;
  loading: boolean;
  error: ApolloError | null;
}
export const useContactNameFromPhoneNumber = ({
  phoneNumber,
}: Props): Result => {
  const { data, loading, error } = useGetContactNameByPhoneNumberQuery({
    variables: { e164: phoneNumber },
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
    data: data?.contact_ByPhone,
    loading,
    error: null,
  };
};
