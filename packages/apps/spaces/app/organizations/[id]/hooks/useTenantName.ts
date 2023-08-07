import { useEffect } from 'react';
import { useRecoilState } from 'recoil';
import { tenantName } from '@spaces/globalState/userData';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

export const TenantNameDocument = `
    query TenantName {
  tenant
}
    `;
interface Result {
  tenant: string;
}
export const useTenantName = (): Result => {
  const client = getGraphQLClient();

  const [tenant, setTenantName] = useRecoilState(tenantName);

  useEffect(() => {
    if (!tenant.length) {
      client.request<any>(TenantNameDocument).then((res) => {
        console.log('🏷️ ----- res: '
            , res);
        setTenantName(res?.tenant ?? '');
      });
    }
  }, []);

  return {
    tenant,
  };
};
