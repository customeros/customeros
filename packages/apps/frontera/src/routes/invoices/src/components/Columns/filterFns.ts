import { match } from 'ts-pattern';
import { isAfter } from 'date-fns/isAfter';
import { FilterItem } from '@store/types.ts';
import { InvoiceStore } from '@store/Invoices/Invoice.store.ts';

import { Filter, ColumnViewType } from '@graphql/types';

export const getPredefinedFilterFn = (
  serverFilter: FilterItem | null | undefined,
) => {
  const noop = (_row: InvoiceStore) => true;
  if (!serverFilter) return noop;

  return match(serverFilter)
    .with(
      { property: ColumnViewType.InvoicesBillingCycle },
      (filter) => (row: InvoiceStore) => {
        const filterValues = filter?.value;
        if (!filter.active) return true;

        const billingCycle =
          row?.contract?.billingDetails?.billingCycleInMonths;

        if (!Array.isArray(filterValues)) {
          return false;
        }

        return filterValues.some((e) => e === billingCycle);
      },
    )

    .with(
      { property: ColumnViewType.InvoicesInvoicePreview },
      (filter) => (row: InvoiceStore) => {
        const filterValues = filter?.value;

        return row.value?.preview === filterValues;
      },
    )
    .with(
      { property: ColumnViewType.InvoicesPaymentStatus },
      (filter) => (row: InvoiceStore) => {
        const filterValues = filter?.value;
        if (!filter.active) return true;

        return filterValues.includes(row.value?.status);
      },
    )
    .with(
      { property: ColumnViewType.InvoicesIssueDate },
      (filter) => (row: InvoiceStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const issued = row?.value?.issued;

        return isAfter(new Date(issued), new Date(filterValue));
      },
    )
    .with(
      { property: ColumnViewType.InvoicesIssueDatePast },
      (filter) => (row: InvoiceStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const issued = row?.value?.issued;

        return isAfter(new Date(issued), new Date(filterValue));
      },
    )
    .with(
      { property: ColumnViewType.InvoicesInvoiceStatus },
      (filter) => (row: InvoiceStore) => {
        const filterValues = filter?.value;
        if (!filter.active) return true;

        const value = row.contract?.contractEnded;

        if (filterValues.length === 0 || filterValues.length === 2) return true;

        return (
          (filterValues[0] === 'ON_HOLD' && value) ||
          (filterValues[0] === 'SCHEDULED' && !value)
        );
      },
    )

    .with({ property: 'INVOICE_DRY_RUN' }, (filter) => (row: InvoiceStore) => {
      const filterValues = filter?.value;

      return row.value?.dryRun === filterValues;
    })

    .otherwise(() => noop);
};
export const getInvoiceFilterFns = (filters: Filter | null) => {
  if (!filters || !filters.AND) return [];

  const data = filters?.AND;

  return data.map(({ filter }) => getPredefinedFilterFn(filter));
};
