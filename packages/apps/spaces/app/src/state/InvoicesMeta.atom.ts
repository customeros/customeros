import { atom, useRecoilState } from 'recoil';

import { ComparisonOperator } from '@graphql/types';
import { filterOutDryRunInvoices } from '@shared/components/Invoice/utils';

interface InvoicesMeta {
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

export const InvoicesMetaAtom = atom<InvoicesMeta>({
  key: 'InvoicesMeta',
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

export const useInvoicesMeta = () => {
  return useRecoilState(InvoicesMetaAtom);
};
