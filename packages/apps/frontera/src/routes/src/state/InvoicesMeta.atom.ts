import { atom, useRecoilState } from 'recoil';

import { SortBy, Filter } from '@graphql/types';

interface InvoicesMeta {
  getInvoices: {
    where?: Filter;
    sort?: SortBy[];
    pagination: {
      page: number;
      limit: number;
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
    },
  },
});

export const useInvoicesMeta = () => {
  return useRecoilState(InvoicesMetaAtom);
};
