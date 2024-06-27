import { match } from 'ts-pattern';
import { Store } from '@store/store.ts';
import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';

import { Skeleton } from '@ui/feedback/Skeleton';
import { createColumnHelper } from '@ui/presentation/Table';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';
import { Filter, Invoice, TableViewDef, InvoiceStatus } from '@graphql/types';

import {
  AmountCell,
  DueDateCell,
  ContractCell,
  IssueDateCell,
  BillingCycleCell,
  InvoiceStatusCell,
  PaymentStatusCell,
  InvoiceNumberCell,
  InvoicePreviewCell,
} from './Cells';
import {
  IssueDateFilter,
  filterIssueDateFn,
  BillingCycleFilter,
  PaymentStatusFilter,
  InvoiceStatusFilter,
  filterBillingCycleFn,
  filterPaymentStatusFn,
  filterInvoiceStatusFn,
  filterIssueDatePastFn,
} from './Filters';

type ColumnDatum = Store<Invoice>;

// REASON: we do not care about exhaustively typing this TValue type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Column = ColumnDefinition<ColumnDatum, any>;

const columnHelper = createColumnHelper<ColumnDatum>();

const columns: Record<string, Column> = {
  INVOICES_ISSUE_DATE: columnHelper.accessor((row) => row, {
    id: 'INVOICE_ISSUED_DATE',
    size: 150,
    enableColumnFilter: true,
    enableSorting: true,
    filterFn: filterIssueDateFn,
    header: (props) => (
      <THead
        id='issueDate'
        title='Issue date'
        renderFilter={() => (
          <IssueDateFilter onFilterValueChange={props.column.setFilterValue} />
        )}
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => <IssueDateCell value={props.getValue()?.value?.issued} />,
    skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  }),
  // this needs to be removed - INVOICES_ISSUE_DATE is the good one.
  INVOICES_ISSUE_DATE_PAST: columnHelper.accessor((row) => row, {
    id: 'INVOICE_CREATED_AT',
    size: 150,
    enableColumnFilter: true,
    enableSorting: true,
    filterFn: filterIssueDatePastFn,
    header: (props) => (
      <THead
        id='issueDate'
        title='Created at'
        renderFilter={() => (
          <IssueDateFilter
            isPast
            onFilterValueChange={props.column.setFilterValue}
          />
        )}
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => <IssueDateCell value={props.getValue()?.value?.issued} />,
    skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  }),
  INVOICES_DUE_DATE: columnHelper.accessor((row) => row, {
    id: 'INVOICE_DUE_DATE',
    size: 150,
    enableColumnFilter: false,
    enableSorting: true,
    header: (props) => (
      <THead id='dueDate' title='Due date' {...getTHeadProps(props)} />
    ),
    cell: (props) => <DueDateCell value={props.getValue()?.value?.due} />,
    skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  }),
  INVOICES_CONTRACT: columnHelper.accessor((row) => row, {
    id: 'CONTRACT',
    size: 200,
    enableColumnFilter: false,
    enableSorting: false,
    header: (props) => (
      <THead id='contract' title='Contract' {...getTHeadProps(props)} />
    ),
    cell: (props) => {
      return (
        <ContractCell
          contractId={props.getValue()?.value?.contract?.metadata?.id}
          organizationId={props.getValue()?.value?.organization?.metadata?.id}
        />
      );
    },
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),
  INVOICES_BILLING_CYCLE: columnHelper.accessor((row) => row, {
    id: 'CONTRACT_BILLING_CYCLE',
    size: 150,
    enableColumnFilter: true,
    enableSorting: false,
    filterFn: filterBillingCycleFn,
    header: (props) => (
      <THead
        id='billingCycle'
        title='Billing cycle'
        renderFilter={() => <BillingCycleFilter column={props?.column} />}
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => (
      <BillingCycleCell
        contractId={props.getValue()?.value?.contract?.metadata?.id}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),
  INVOICES_PAYMENT_STATUS: columnHelper.accessor((row) => row, {
    id: 'INVOICE_STATUS',
    size: 175,
    enableColumnFilter: true,
    enableSorting: true,
    filterFn: filterPaymentStatusFn,
    header: (props) => (
      <THead
        id='paymentStatus'
        title='Payment status'
        renderFilter={() => <PaymentStatusFilter column={props?.column} />}
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => (
      <PaymentStatusCell
        value={props.getValue()?.value?.status}
        invoiceId={props.getValue()?.value?.metadata?.id}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),
  INVOICES_AMOUNT: columnHelper.accessor((row) => row, {
    id: 'AMOUNT',
    size: 100,
    enableColumnFilter: false,
    enableSorting: false,
    header: (props) => (
      <THead id='amount' title='Amount' {...getTHeadProps(props)} />
    ),
    cell: (props) => (
      <AmountCell
        value={props.getValue()?.value?.amountDue}
        currency={props.getValue().value?.currency}
      />
    ),
    skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  }),
  INVOICES_INVOICE_NUMBER: columnHelper.accessor((row) => row, {
    id: 'INVOICES_INVOICE_NUMBER',
    size: 100,
    enableColumnFilter: false,
    enableSorting: false,
    header: (props) => (
      <THead id='invoiceNumber' title='Invoice' {...getTHeadProps(props)} />
    ),
    cell: (props) => (
      <InvoiceNumberCell
        value={props.getValue()?.value?.invoiceNumber}
        invoiceId={props.getValue()?.value?.metadata?.id}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),
  INVOICES_INVOICE_STATUS: columnHelper.accessor((row) => row, {
    id: 'INVOICES_INVOICE_STATUS',
    size: 150,
    enableColumnFilter: true,
    enableSorting: true,
    filterFn: filterInvoiceStatusFn,
    header: (props) => (
      <THead
        id='invoiceStatus'
        title='Invoice status'
        renderFilter={() => <InvoiceStatusFilter column={props?.column} />}
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => (
      <InvoiceStatusCell status={props.getValue()?.value?.status} />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),
  INVOICES_INVOICE_PREVIEW: columnHelper.accessor((row) => row, {
    id: 'INVOICE_PREVIEW',
    size: 100,
    enableColumnFilter: false,
    enableSorting: false,
    header: (props) => (
      <THead
        id='invoicePreview'
        title='Invoice preview'
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => (
      <InvoicePreviewCell
        value={props.getValue()?.value?.invoiceNumber}
        invoiceId={props.getValue()?.value?.metadata?.id}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),
  INVOICES_PLACEHOLDER: columnHelper.accessor((row) => row, {
    id: 'PLACEHOLDER',
    size: 32,
    fixWidth: true,
    header: () => <></>,
    cell: () => <></>,
    skeleton: () => <></>,
  }),
};

export const getColumnsConfig = (tableViewDef?: TableViewDef) => {
  if (!tableViewDef) return [];

  return (tableViewDef.columns ?? []).reduce((acc, curr) => {
    const columnTypeName = curr?.columnType;

    if (!columnTypeName) return acc;

    const column = { ...columns[columnTypeName], enableHiding: !curr.visible };

    if (!column) return acc;

    return [...acc, column];
  }, [] as Column[]);
};
export const getColumnSortFn = (columnId: string) =>
  match(columnId)
    .with(
      'INVOICE_STATUS',
      () => (row: Store<Invoice>) =>
        match(row.value?.status)
          .with(InvoiceStatus.Empty, () => null)
          .with(InvoiceStatus.Initialized, () => 1)
          .with(InvoiceStatus.OnHold, () => 2)
          .with(InvoiceStatus.Scheduled, () => 3)
          .with(InvoiceStatus.Void, () => 4)
          .with(InvoiceStatus.Paid, () => 5)
          .with(InvoiceStatus.Due, () => 6)
          .with(InvoiceStatus.Overdue, () => 7)
          .otherwise(() => null),
    )

    .with('INVOICE_DUE_DATE', () => (row: Store<Invoice>) => {
      const value = row.value?.due;

      return value ? new Date(value) : null;
    })
    .with('INVOICE_ISSUED_DATE', () => (row: Store<Invoice>) => {
      const value = row.value?.due;

      return value ? new Date(value) : null;
    })
    .with('INVOICE_CREATED_AT', () => (row: Store<Invoice>) => {
      const value = row.value?.due;

      return value ? new Date(value) : null;
    })
    .otherwise(() => (_row: Store<Invoice>) => null);

export const getPredefinedFilterFn = (serverFilter: Filter | null) => {
  if (!serverFilter) return null;

  const data = serverFilter?.AND?.[0];

  return match(data?.filter)
    .with(
      { property: 'INVOICE_PREVIEW' },
      (filter) => (row: Store<Invoice>) => {
        const filterValues = filter?.value;

        return row.value?.preview === filterValues;
      },
    )
    .with(
      { property: 'INVOICE_DRY_RUN' },
      (filter) => (row: Store<Invoice>) => {
        const filterValues = filter?.value;

        return row.value?.dryRun === filterValues;
      },
    )

    .otherwise(() => null);
};
