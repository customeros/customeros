import { InvoiceStore } from '@store/Invoices/Invoice.store.ts';
import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';

import { Skeleton } from '@ui/feedback/Skeleton';
import { createColumnHelper } from '@ui/presentation/Table';
import { TableViewDef, ColumnViewType } from '@graphql/types';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';
import { getColumnConfig } from '@organizations/components/Columns/shared/util/getColumnConfig.ts';

import {
  IssueDateFilter,
  BillingCycleFilter,
  InvoiceStatusFilter,
  PaymentStatusFilter,
} from './Filters';
import {
  AmountCell,
  DueDateCell,
  ContractCell,
  IssueDateCell,
  BillingCycleCell,
  InvoiceNumberCell,
  InvoiceStatusCell,
  PaymentStatusCell,
  InvoicePreviewCell,
} from './Cells';

type ColumnDatum = InvoiceStore;

// REASON: we do not care about exhaustively typing this TValue type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Column = ColumnDefinition<ColumnDatum, any>;

const columnHelper = createColumnHelper<ColumnDatum>();

const columns: Record<string, Column> = {
  [ColumnViewType.InvoicesIssueDate]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesIssueDate,
    size: 150,
    enableColumnFilter: true,
    enableSorting: true,
    header: (props) => (
      <THead
        id='issueDate'
        filterWidth={250}
        title='Issue date'
        renderFilter={() => (
          <IssueDateFilter property={ColumnViewType.InvoicesIssueDate} />
        )}
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => <IssueDateCell value={props.getValue()?.value?.issued} />,
    skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  }),
  // this needs to be removed - INVOICES_ISSUE_DATE is the good one.
  [ColumnViewType.InvoicesIssueDatePast]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesIssueDatePast,
    size: 150,
    enableColumnFilter: true,
    enableSorting: true,
    header: (props) => (
      <THead
        id='issueDateP'
        filterWidth={250}
        title='Created at'
        renderFilter={() => (
          <IssueDateFilter property={ColumnViewType.InvoicesIssueDatePast} />
        )}
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => <IssueDateCell value={props.getValue()?.value?.issued} />,
    skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  }),
  [ColumnViewType.InvoicesDueDate]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesDueDate,
    size: 150,
    enableColumnFilter: true,
    enableSorting: true,
    header: (props) => (
      <THead
        id='dueDate'
        filterWidth={250}
        title='Due date'
        {...getTHeadProps(props)}
        renderFilter={() => (
          <IssueDateFilter property={ColumnViewType.InvoicesDueDate} />
        )}
      />
    ),
    cell: (props) => <DueDateCell value={props.getValue()?.value?.due} />,
    skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  }),
  [ColumnViewType.InvoicesContract]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesContract,
    size: 225,
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
  [ColumnViewType.InvoicesBillingCycle]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesBillingCycle,
    size: 150,
    enableColumnFilter: true,
    enableSorting: false,
    header: (props) => (
      <THead
        id='billingCycle'
        title='Billing cycle'
        renderFilter={() => <BillingCycleFilter />}
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => (
      <BillingCycleCell id={props.getValue()?.value?.metadata.id} />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),
  [ColumnViewType.InvoicesPaymentStatus]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesPaymentStatus,
    size: 175,
    enableColumnFilter: true,
    enableSorting: true,
    header: (props) => (
      <THead
        id='paymentStatus'
        title='Payment status'
        renderFilter={() => <PaymentStatusFilter />}
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
  [ColumnViewType.InvoicesAmount]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesAmount,
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
  [ColumnViewType.InvoicesInvoiceNumber]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesInvoiceNumber,
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
  [ColumnViewType.InvoicesInvoiceStatus]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesInvoiceStatus,
    size: 150,
    enableColumnFilter: true,
    enableSorting: true,
    header: (props) => (
      <THead
        id='invoiceStatus'
        title='Invoice status'
        renderFilter={() => <InvoiceStatusFilter />}
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => (
      <InvoiceStatusCell status={props.getValue()?.value?.status} />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),
  [ColumnViewType.InvoicesInvoicePreview]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesInvoicePreview,
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

export const getInvoiceColumnsConfig = (
  tableViewDef?: Array<TableViewDef>[0],
) => getColumnConfig<ColumnDatum>(columns, tableViewDef);
