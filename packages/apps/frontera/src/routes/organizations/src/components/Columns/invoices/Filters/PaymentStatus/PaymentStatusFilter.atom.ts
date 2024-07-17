import { atom, selector, useRecoilState } from 'recoil';

import { InvoiceStatus } from '@graphql/types';
import { makeServerToAtomMapper } from '@shared/components/Filters/makeServerToAtomMapper';

export interface PaymentStatusFilterState {
  isActive: boolean;
  value: InvoiceStatus[];
}

export const defaultState: PaymentStatusFilterState = {
  value: [],
  isActive: false,
};

export const PaymentStatusFilterAtom = atom<PaymentStatusFilterState>({
  key: 'payment-status-filter',
  default: defaultState,
});

export const PaymentStatusFilterSelector = selector({
  key: 'payment-status-filter-selector',
  get: ({ get }) => get(PaymentStatusFilterAtom),
});

export const usePaymentStatusFilter = () => {
  return useRecoilState(PaymentStatusFilterAtom);
};

/**
 * Used for mapping server-side Filter data to client-side atom PaymentStatusFilterState
 */
export const mapPaymentStatusToAtom =
  makeServerToAtomMapper<PaymentStatusFilterState>(
    {
      filter: {
        property: 'INVOICES_PAYMENT_STATUS',
      },
    },
    ({ filter }) => ({
      isActive: true,
      value: filter?.value,
    }),
    defaultState,
  );
