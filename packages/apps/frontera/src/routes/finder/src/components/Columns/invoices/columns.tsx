import { InvoiceStore } from '@store/Invoices/Invoice.store.ts';
import { PaymentStatusSelect } from '@invoices/components/shared';
import {
  ColumnDef,
  ColumnDef as ColumnDefinition,
} from '@tanstack/react-table';
import { DateCell } from '@finder/components/Columns/shared/Cells/DateCell/DateCell.tsx';
import { getColumnConfig } from '@finder/components/Columns/shared/util/getColumnConfig.ts';

import { Skeleton } from '@ui/feedback/Skeleton';
import { createColumnHelper } from '@ui/presentation/Table';
import { TableViewDef, ColumnViewType } from '@graphql/types';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';

import { OrganizationName } from './Cells/OrganizationName';
import {
  AmountCell,
  ContractCell,
  BillingCycleCell,
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
    minSize: 150,
    maxSize: 300,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    header: (props) => (
      <THead
        filterWidth={250}
        title='Issue Date'
        id={ColumnViewType.InvoicesIssueDate}
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => <DateCell value={props.getValue()?.value?.issued} />,
    skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  }),
  // this needs to be removed - INVOICES_ISSUE_DATE is the good one.
  [ColumnViewType.InvoicesIssueDatePast]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesIssueDatePast,
    size: 150,
    minSize: 150,
    maxSize: 300,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    header: (props) => (
      <THead
        filterWidth={250}
        title='Created At'
        id={ColumnViewType.InvoicesIssueDatePast}
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => <DateCell value={props.getValue()?.value?.issued} />,
    skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  }),
  [ColumnViewType.InvoicesDueDate]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesDueDate,
    size: 150,
    minSize: 150,
    maxSize: 300,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    header: (props) => (
      <THead
        title='Due Date'
        filterWidth={250}
        id={ColumnViewType.InvoicesDueDate}
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => <DateCell value={props.getValue()?.value?.due} />,
    skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  }),
  [ColumnViewType.InvoicesContract]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesContract,
    size: 160,
    minSize: 160,
    maxSize: 600,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    header: (props) => (
      <THead
        title='Contract'
        id={ColumnViewType.InvoicesContract}
        {...getTHeadProps(props)}
      />
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

  [ColumnViewType.InvoicesOrganization]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesOrganization,
    size: 160,
    minSize: 160,
    maxSize: 600,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    header: (props) => (
      <THead
        title='Organization'
        id={ColumnViewType.InvoicesOrganization}
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => {
      const orgId = props.getValue()?.value?.organization?.metadata?.id;

      return <OrganizationName orgId={orgId} />;
    },
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),
  [ColumnViewType.InvoicesBillingCycle]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesBillingCycle,
    size: 150,
    minSize: 150,
    maxSize: 300,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: false,
    header: (props) => (
      <THead
        title='Billing Cycle'
        id={ColumnViewType.InvoicesBillingCycle}
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => (
      <BillingCycleCell id={props.getValue()?.value?.metadata.id} />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),
  [ColumnViewType.InvoicesInvoiceStatus]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesInvoiceStatus,
    size: 175,
    minSize: 175,
    maxSize: 300,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    header: (props) => (
      <THead
        title='Invoice Status'
        id={ColumnViewType.InvoicesInvoiceStatus}
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => (
      <PaymentStatusSelect invoiceNumber={props.getValue()?.number} />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),
  [ColumnViewType.InvoicesAmount]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesAmount,
    size: 100,
    minSize: 100,
    maxSize: 300,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: false,
    header: (props) => (
      <THead
        title='Amount'
        id={ColumnViewType.InvoicesAmount}
        {...getTHeadProps(props)}
      />
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
    size: 120,
    minSize: 120,
    maxSize: 300,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: false,
    header: (props) => (
      <THead
        title='Invoice'
        id={ColumnViewType.InvoicesInvoiceNumber}
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => (
      <InvoicePreviewCell value={props.getValue()?.value?.invoiceNumber} />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),

  [ColumnViewType.InvoicesInvoicePreview]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesInvoicePreview,
    size: 150,
    minSize: 150,
    maxSize: 300,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: false,
    header: (props) => (
      <THead
        title='Invoice'
        id={ColumnViewType.InvoicesInvoicePreview}
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => (
      <InvoicePreviewCell value={props.getValue()?.value?.invoiceNumber} />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),
  INVOICES_PLACEHOLDER: columnHelper.accessor((row) => row, {
    id: 'PLACEHOLDER',
    size: 32,
    minSize: 32,
    maxSize: 32,
    fixWidth: true,
    header: () => <></>,
    cell: () => <></>,
    skeleton: () => <></>,
  }),
};

export const getInvoiceColumnsConfig = (
  tableViewDef?: Array<TableViewDef>[0],
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
): ColumnDef<ColumnDatum, any>[] =>
  getColumnConfig<ColumnDatum>(columns, tableViewDef);
