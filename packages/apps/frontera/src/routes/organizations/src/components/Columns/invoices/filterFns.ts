import { match } from 'ts-pattern';
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
    .with({ property: 'INVOICE_PREVIEW' }, (filter) => (row: InvoiceStore) => {
      const filterValues = filter?.value;

      return row.value?.preview === filterValues;
    })

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
        const issued = row?.value?.issued.split('T')[0];

        if (!filterValue) return true;
        if (filterValue?.[1] === null) return filterValue?.[0] <= issued;
        if (filterValue?.[0] === null) return filterValue?.[1] >= issued;

        return filterValue[0] <= issued && filterValue[1] >= issued;
      },
    )
    .with(
      { property: ColumnViewType.InvoicesIssueDatePast },
      (filter) => (row: InvoiceStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const issued = row?.value?.issued.split('T')[0];

        if (!filterValue) return true;
        if (filterValue?.[1] === null) return filterValue?.[0] <= issued;
        if (filterValue?.[0] === null) return filterValue?.[1] >= issued;

        return filterValue[0] <= issued && filterValue[1] >= issued;
      },
    )
    .with(
      { property: ColumnViewType.InvoicesDueDate },
      (filter) => (row: InvoiceStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const due = row?.value?.due.split('T')[0];

        if (!filterValue) return true;
        if (filterValue?.[1] === null) return filterValue?.[0] <= due;
        if (filterValue?.[0] === null) return filterValue?.[1] >= due;

        return filterValue[0] <= due && filterValue[1] >= due;
      },
    )
    .with(
      { property: ColumnViewType.InvoicesInvoiceStatus },
      (filter) => (row: InvoiceStore) => {
        const filterValues = filter?.value;

        if (!filter.active) return true;

        const value = row.value.status;

        return filterValues.includes(value);
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
