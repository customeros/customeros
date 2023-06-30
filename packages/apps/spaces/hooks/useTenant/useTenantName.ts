import { ApolloError } from '@apollo/client';
import { useGetTenantNameLazyQuery } from '@spaces/graphql';
import { useEffect } from 'react';
import { tenantName } from '../../state/userData';
import { useRecoilState } from 'recoil';

interface Result {
  loading: boolean;
  error?: ApolloError | null;
  loadTenantName: () => void;
}
export const useTenantName = (): Result => {
  const [loadTenantName, { loading, error }] = useGetTenantNameLazyQuery();

  const [tenant, setTenantName] = useRecoilState(tenantName);

  useEffect(() => {
    if (!tenant.length) {
      loadTenantName().then((res) => {
        setTenantName(res.data?.tenant ?? '');
      });
    }
  }, []);

  return {
    loading,
    error,
    loadTenantName,
  };
};
