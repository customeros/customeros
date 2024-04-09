import { atom, selector, useRecoilState } from 'recoil';

import { makeServerToAtomMapper } from '@shared/components/Filters/makeServerToAtomMapper';

export interface InvoiceStatusFilterState {
  isActive: boolean;
  value: ('ON_HOLD' | 'SCHEDULED')[];
}

export const defaultState: InvoiceStatusFilterState = {
  value: [],
  isActive: false,
};

export const InvoiceStatusFilterAtom = atom<InvoiceStatusFilterState>({
  key: 'invoice-status-filter',
  default: defaultState,
});

export const InvoiceStatusFilterSelector = selector({
  key: 'invoice-status-filter-selector',
  get: ({ get }) => get(InvoiceStatusFilterAtom),
});

export const useInvoiceStatusFilter = () => {
  return useRecoilState(InvoiceStatusFilterAtom);
};

/**
 * Used for mapping server-side Filter data to client-side atom InvoiceStatusFilterState
 */
export const mapInvoiceStatusToAtom =
  makeServerToAtomMapper<InvoiceStatusFilterState>(
    {
      filter: {
        property: 'INVOICE_STATUS',
      },
    },
    ({ filter }) => ({
      isActive: true,
      value: filter?.value,
    }),
    defaultState,
  );
