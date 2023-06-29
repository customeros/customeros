import { ApolloError } from '@apollo/client';
import { GetOrganizationQuery, useGetOrganizationQuery } from '@spaces/graphql';

interface Props {
  id: string;
}

interface Result {
  data: GetOrganizationQuery['organization'] | undefined | null;
  loading: boolean;
  error: ApolloError | null | undefined;
}
export const useOrganization = ({ id }: Props): Result => {
  const { data, loading, error } = useGetOrganizationQuery({
    variables: { id },
  });

  return {
    data: data?.organization,
    loading,
    error,
  };
};
