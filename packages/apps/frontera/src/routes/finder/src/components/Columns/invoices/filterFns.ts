import { match } from 'ts-pattern';
import { FilterItem } from '@store/types.ts';
import { isAfter, isBefore } from 'date-fns';
import { InvoiceStore } from '@store/Invoices/Invoice.store.ts';

import { Filter, ColumnViewType, ComparisonOperator } from '@graphql/types';

export const getFilterFn = (serverFilter: FilterItem | null | undefined) => {
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

export const getFilterV2Fn = (serverFilter: FilterItem | null | undefined) => {
  const noop = (_row: InvoiceStore) => true;

  if (!serverFilter) return noop;

  return match(serverFilter)
    .with(
      { property: ColumnViewType.InvoicesBillingCycle },
      (filter) => (row: InvoiceStore) => {
        if (!filter.active) return true;
        const value = row.value?.contract?.billingDetails?.billingCycleInMonths;

        if (!value)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(filter, [value]);
      },
    )
    .with(
      { property: ColumnViewType.InvoicesInvoicePreview },
      (filter) => (row: InvoiceStore) => {
        if (!filter.active) return true;
        const value = row.value?.invoiceNumber.trim();

        if (!value) return false;

        return filterTypeText(filter, value);
      },
    )

    .with(
      { property: ColumnViewType.InvoicesInvoiceNumber },
      (filter) => (row: InvoiceStore) => {
        if (!filter.active) return true;
        const value = row.value?.invoiceNumber.trim();

        if (!value) return false;

        return filterTypeText(filter, value);
      },
    )

    .with(
      { property: ColumnViewType.InvoicesContract },
      (filter) => (row: InvoiceStore) => {
        if (!filter.active) return true;

        const value = row.value.contract.contractName;

        if (!value)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeText(filter, value);
      },
    )

    .with(
      { property: ColumnViewType.InvoicesIssueDate },
      (filter) => (row: InvoiceStore) => {
        if (!filter.active) return true;
        const value = row?.value?.issued.split('T')[0];

        if (!value) return false;

        return filterTypeDate(filter, value);
      },
    )
    .with(
      { property: ColumnViewType.InvoicesIssueDatePast },
      (filter) => (row: InvoiceStore) => {
        if (!filter.active) return true;
        const value = row?.value?.issued.split('T')[0];

        if (!value) return false;

        return filterTypeDate(filter, value);
      },
    )
    .with(
      { property: ColumnViewType.InvoicesDueDate },
      (filter) => (row: InvoiceStore) => {
        if (!filter.active) return true;
        const value = row?.value?.due.split('T')[0];

        if (!value) return false;

        return filterTypeDate(filter, value);
      },
    )
    .with(
      { property: ColumnViewType.InvoicesInvoiceStatus },
      (filter) => (row: InvoiceStore) => {
        if (!filter.active) return true;

        const value = row.value.status;

        if (!value)
          return (
            filter.operation === ComparisonOperator.IsEmpty ||
            filter.operation === ComparisonOperator.NotContains
          );

        return filterTypeList(filter, [value]);
      },
    )

    .with(
      { property: ColumnViewType.InvoicesAmount },
      (filter) => (row: InvoiceStore) => {
        if (!filter.active) return true;
        const value = row.value.amountDue;

        if (!value) return false;

        return filterTypeNumber(filter, value);
      },
    )
    .with({ property: 'INVOICE_DRY_RUN' }, (filter) => (row: InvoiceStore) => {
      const filterValues = filter?.value;

      return row.value?.dryRun === filterValues;
    })

    .otherwise(() => noop);
};

const filterTypeText = (filter: FilterItem, value: string | undefined) => {
  const filterValue = filter?.value?.toLowerCase();
  const filterOperator = filter?.operation;
  const valueLower = value?.toLowerCase();

  return match(filterOperator)
    .with(ComparisonOperator.IsEmpty, () => !value)
    .with(ComparisonOperator.IsNotEmpty, () => value)
    .with(
      ComparisonOperator.NotContains,
      () => !valueLower?.includes(filterValue),
    )
    .with(ComparisonOperator.Contains, () => valueLower?.includes(filterValue))
    .otherwise(() => false);
};

const filterTypeNumber = (filter: FilterItem, value: number | undefined) => {
  const filterValue = filter?.value;
  const filterOperator = filter?.operation;

  if (value === undefined || value === null) return false;

  return match(filterOperator)
    .with(ComparisonOperator.Lt, () => value < Number(filterValue))
    .with(ComparisonOperator.Gt, () => value > Number(filterValue))
    .with(ComparisonOperator.Eq, () => value === Number(filterValue))
    .with(ComparisonOperator.NotEqual, () => value !== Number(filterValue))
    .otherwise(() => true);
};

const filterTypeList = (filter: FilterItem, value: string[] | undefined) => {
  const filterValue = filter?.value;
  const filterOperator = filter?.operation;

  return match(filterOperator)
    .with(ComparisonOperator.IsEmpty, () => !value?.length)
    .with(ComparisonOperator.IsNotEmpty, () => value?.length)
    .with(
      ComparisonOperator.NotContains,
      () =>
        !value?.length ||
        (value?.length && !value.some((v) => filterValue?.includes(v))),
    )
    .with(
      ComparisonOperator.Contains,
      () => value?.length && value.some((v) => filterValue?.includes(v)),
    )
    .otherwise(() => false);
};

const filterTypeDate = (filter: FilterItem, value: string | undefined) => {
  const filterValue = filter?.value;
  const filterOperator = filter?.operation;

  if (!value) return false;

  return match(filterOperator)
    .with(ComparisonOperator.Lt, () =>
      isBefore(new Date(value), new Date(filterValue)),
    )
    .with(ComparisonOperator.Gt, () =>
      isAfter(new Date(value), new Date(filterValue)),
    )

    .otherwise(() => true);
};

export const getInvoiceDefaultFilterFns = (
  filters: Filter | null,
  isFeatureEnabled: boolean,
) => {
  if (!filters || !filters.AND) return [];

  const data = filters?.AND;

  if (isFeatureEnabled) {
    return data.map(({ filter }) => getFilterV2Fn(filter));
  }

  return data.map(({ filter }) => getFilterFn(filter));
};

export const getInvoiceFilterFns = (
  filters: Filter | null,
  isFeatureEnabled: boolean,
) => {
  if (!filters || !filters.AND) return [];

  const data = filters?.AND;

  if (isFeatureEnabled) {
    return data.map(({ filter }) => getFilterV2Fn(filter));
  }

  return data.map(({ filter }) => getFilterFn(filter));
};
