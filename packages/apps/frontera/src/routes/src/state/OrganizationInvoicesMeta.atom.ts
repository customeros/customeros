import { atom, useRecoilState } from 'recoil';

import { ComparisonOperator } from '@graphql/types';
import { filterOutDryRunInvoices } from '@shared/components/Invoice/utils';

interface OrganizationInvoicesMeta {
  getInvoices: {
    organizationId?: string;
    pagination: {
      page: number;
      limit: number;
    };
    where: {
      AND: Array<{
        filter: {
          value: boolean;
          property: string;
          operation: ComparisonOperator;
        };
      }>;
    };
  };
}

export const OrganizationInvoicesMetaAtom = atom<OrganizationInvoicesMeta>({
  key: 'OrganizationInvoicesMeta',
  default: {
    getInvoices: {
      pagination: {
        page: 1,
        limit: 40,
      },
      where: { ...filterOutDryRunInvoices },
    },
  },
});

export const useOrganizationInvoicesMeta = () => {
  return useRecoilState(OrganizationInvoicesMetaAtom);
};
