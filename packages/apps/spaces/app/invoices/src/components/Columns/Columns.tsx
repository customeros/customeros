import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';

import { Invoice, TableViewDef } from '@graphql/types';
import { Skeleton } from '@ui/feedback/Skeleton/Skeleton2';
import { createColumnHelper } from '@ui/presentation/Table';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';

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

type ColumnDatum = Invoice;

// REASON: we do not care about exhaustively typing this TValue type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Column = ColumnDefinition<ColumnDatum, any>;

const columnHelper = createColumnHelper<ColumnDatum>();

const columns: Record<string, Column> = {
  INVOICES_ISSUE_DATE: columnHelper.accessor('issued', {
    id: 'INVOICE_ISSUED_DATE',
    minSize: 50,
    maxSize: 50,
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
    cell: (props) => <IssueDateCell value={props.getValue()} />,
    skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  }),
  // this needs to be removed - INVOICES_ISSUE_DATE is the good one.
  INVOICES_ISSUE_DATE_PAST: columnHelper.accessor('issued', {
    id: 'INVOICE_CREATED_AT',
    minSize: 50,
    maxSize: 50,
    enableColumnFilter: true,
    enableSorting: true,
    filterFn: filterIssueDatePastFn,
    header: (props) => (
      <THead
        id='issueDate'
        title='Issue date'
        renderFilter={() => (
          <IssueDateFilter
            isPast
            onFilterValueChange={props.column.setFilterValue}
          />
        )}
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => <IssueDateCell value={props.getValue()} />,
    skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  }),
  INVOICES_DUE_DATE: columnHelper.accessor('due', {
    id: 'INVOICE_DUE_DATE',
    minSize: 50,
    maxSize: 50,
    enableColumnFilter: false,
    enableSorting: true,
    header: (props) => (
      <THead id='dueDate' title='Due date' {...getTHeadProps(props)} />
    ),
    cell: (props) => <DueDateCell value={props.getValue()} />,
    skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  }),
  INVOICES_CONTRACT: columnHelper.accessor((row) => row, {
    id: 'CONTRACT',
    minSize: 200,
    enableColumnFilter: false,
    enableSorting: false,
    header: (props) => (
      <THead id='contract' title='Contract' {...getTHeadProps(props)} />
    ),
    cell: (props) => (
      <ContractCell
        value={props.getValue()?.contract?.name}
        organizationName={props.getValue()?.organization?.name}
        organizationId={props.getValue()?.organization?.metadata?.id}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),
  INVOICES_BILLING_CYCLE: columnHelper.accessor('contract.billingDetails', {
    id: 'CONTRACT_BILLING_CYCLE',
    minSize: 100,
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
    cell: (props) => <BillingCycleCell value={props.getValue().billingCycle} />,
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),
  INVOICES_PAYMENT_STATUS: columnHelper.accessor((row) => row, {
    id: 'INVOICE_STATUS',
    minSize: 70,
    maxSize: 70,
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
        value={props.getValue()?.status}
        invoiceId={props.getValue()?.metadata?.id}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),
  INVOICES_AMOUNT: columnHelper.accessor('amountDue', {
    id: 'AMOUNT',
    minSize: 100,
    maxSize: 100,
    enableColumnFilter: false,
    enableSorting: false,
    header: (props) => (
      <THead id='amount' title='Amount' {...getTHeadProps(props)} />
    ),
    cell: (props) => (
      <AmountCell
        value={props.getValue()}
        currency={props.row.original.currency}
      />
    ),
    skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  }),
  INVOICES_INVOICE_NUMBER: columnHelper.accessor((row) => row, {
    id: 'INVOICE_NUMBER',
    minSize: 100,
    maxSize: 100,
    enableColumnFilter: false,
    enableSorting: false,
    header: (props) => (
      <THead id='invoiceNumber' title='Invoice' {...getTHeadProps(props)} />
    ),
    cell: (props) => (
      <InvoiceNumberCell
        value={props.getValue()?.invoiceNumber}
        invoiceId={props.getValue()?.metadata?.id}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),
  INVOICES_INVOICE_STATUS: columnHelper.accessor('contract.contractEnded', {
    id: 'CONTRACT_ENDED_AT',
    minSize: 100,
    maxSize: 100,
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
    cell: (props) => <InvoiceStatusCell isOutOfContract={props.getValue()} />,
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),
  INVOICES_INVOICE_PREVIEW: columnHelper.accessor((row) => row, {
    id: 'INVOICE_PREVIEW',
    minSize: 100,
    maxSize: 100,
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
        value={props.getValue()?.invoiceNumber}
        invoiceId={props.getValue()?.metadata?.id}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),
  INVOICES_PLACEHOLDER: columnHelper.accessor((row) => row, {
    id: 'PLACEHOLDER',
    minSize: 32,
    maxSize: 32,
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
