import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';

import { Invoice } from '@graphql/types';
import { Skeleton } from '@ui/feedback/Skeleton/Skeleton2';
import { createColumnHelper } from '@ui/presentation/Table';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';
import { TableViewDefsQuery } from '@shared/graphql/tableViewDefs.generated';

import {
  AmountCell,
  DueDateCell,
  ContractCell,
  IssueDateCell,
  BillingCycleCell,
  PaymentTermsCell,
  InvoiceStatusCell,
  InvoicePreviewCell,
} from './Cells';

type ColumnDatum = Invoice;

// REASON: we do not care about exhaustively typing this TValue type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Column = ColumnDefinition<ColumnDatum, any>;

const columnHelper = createColumnHelper<ColumnDatum>();

const columns: Record<string, Column> = {
  ISSUE_DATE: columnHelper.accessor('invoicePeriodStart', {
    id: 'ISSUE_DATE',
    minSize: 200,
    enableColumnFilter: false,
    enableSorting: false,
    header: (props) => (
      <THead id='issueDate' title='Issue date' {...getTHeadProps(props)} />
    ),
    cell: (props) => <IssueDateCell value={props.getValue()} />,
    skeleton: () => <Skeleton className='w-[200px]' />,
  }),
  DUE_DATE: columnHelper.accessor('due', {
    id: 'DUE_DATE',
    minSize: 200,
    header: (props) => (
      <THead id='dueDate' title='Due date' {...getTHeadProps(props)} />
    ),
    cell: (props) => <DueDateCell value={props.getValue()} />,
    skeleton: () => <Skeleton className='w-[200px]' />,
  }),
  CONTRACT: columnHelper.accessor('invoiceNumber', {
    id: 'CONTRACT',
    minSize: 100,
    header: (props) => (
      <THead id='contract' title='Contract' {...getTHeadProps(props)} />
    ),
    cell: (props) => <ContractCell value={props.getValue()} />,
    skeleton: () => <Skeleton className='w-[100px]' />,
  }),
  BILLING_CYCLE: columnHelper.accessor('invoiceNumber', {
    id: 'BILLING_CYCLE',
    minSize: 100,
    header: (props) => (
      <THead
        id='billingCycle'
        title='Billing cycle'
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => <BillingCycleCell value={props.getValue()} />,
    skeleton: () => <Skeleton className='w-[100px]' />,
  }),
  PAYMENT_TERMS: columnHelper.accessor('invoiceNumber', {
    id: 'PAYMENT_TERMS',
    minSize: 100,
    header: (props) => (
      <THead
        id='paymentTerms'
        title='Payment terms'
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => <PaymentTermsCell value={props.getValue()} />,
    skeleton: () => <Skeleton className='w-[100px]' />,
  }),
  AMOUNT: columnHelper.accessor('amountDue', {
    id: 'AMOUNT',
    minSize: 100,
    header: (props) => (
      <THead id='amount' title='Amount' {...getTHeadProps(props)} />
    ),
    cell: (props) => <AmountCell value={props.getValue()} />,
    skeleton: () => <Skeleton className='w-[200px]' />,
  }),
  INVOICE_STATUS: columnHelper.accessor('status', {
    id: 'INVOICE_STATUS',
    minSize: 100,
    header: (props) => (
      <THead
        id='invoiceStatus'
        title='Invoice status'
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => <InvoiceStatusCell value={props.getValue()} />,
    skeleton: () => <Skeleton className='w-[100px]' />,
  }),
  INVOICE_PREVIEW: columnHelper.accessor('invoiceNumber', {
    id: 'INVOICE_PREVIEW',
    minSize: 100,
    header: (props) => (
      <THead
        id='invoicePreview'
        title='Invoice preview'
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => <InvoicePreviewCell value={props.getValue()} />,
    skeleton: () => <Skeleton className='w-[100px]' />,
  }),
  PLACEHOLDER: columnHelper.accessor((row) => row, {
    id: 'PLACEHOLDER',
    minSize: 32,
    maxSize: 32,
    fixWidth: true,
    header: () => <></>,
    cell: () => <></>,
    skeleton: () => <></>,
  }),
};

export const getColumnsConfig = (
  tableViewDef?: TableViewDefsQuery['tableViewDefs']['content'][0],
) => {
  if (!tableViewDef) return [];

  return (tableViewDef.columns ?? []).reduce((acc, curr) => {
    const columnTypeName = curr?.columnType?.name;

    if (!columnTypeName) return acc;

    const column = columns[columnTypeName];

    if (!column) return acc;

    return [...acc, column];
  }, [] as Column[]);
};
