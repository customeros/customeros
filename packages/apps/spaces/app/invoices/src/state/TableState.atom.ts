import { selector, useRecoilValue } from 'recoil';

import { IssueDateFilterSelector } from '../components/Columns/Filters/IssueDate';
import { BillingCycleFilterSelector } from '../components/Columns/Filters/BillingCycle';
import { InvoiceStatusFilterSelector } from '../components/Columns/Filters/InvoiceStatus';
import { PaymentStatusFilterSelector } from '../components/Columns/Filters/PaymentStatus';

const tableStateSelector = selector({
  key: 'invoices-tableState',
  get: ({ get }) => {
    const issueDate = get(IssueDateFilterSelector);
    const billingCycle = get(BillingCycleFilterSelector);
    const paymentStatus = get(PaymentStatusFilterSelector);

    const invoiceStatus = (() => {
      const state = get(InvoiceStatusFilterSelector);
      const value =
        state.value.length > 1 ? undefined : state.value.includes('ON_HOLD');

      return {
        ...state,
        value,
      };
    })();

    return {
      columnFilters: {
        issueDate,
        billingCycle,
        invoiceStatus,
        paymentStatus,
      },
    };
  },
});

export const useTableState = () => {
  return useRecoilValue(tableStateSelector);
};
