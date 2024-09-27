import { InvoiceStore } from '@store/Invoices/Invoice.store.ts';
import { PaymentStatusSelect } from '@invoices/components/shared';
import { InvoicePreviewCell } from '@finder/components/Columns/invoices/Cells';

import { DateTimeUtils } from '@utils/date';
import { ColumnViewType } from '@graphql/types';
import { Skeleton } from '@ui/feedback/Skeleton';
import { createColumnHelper } from '@ui/presentation/Table';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';

type ColumnDatum = InvoiceStore;

const columnHelper = createColumnHelper<ColumnDatum>();

export const columns = [
  columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesInvoiceNumber,
    size: 100,
    enableColumnFilter: false,
    enableSorting: false,
    header: (props) => (
      <THead title='N°' id='invoiceNumber' {...getTHeadProps(props)} />
    ),
    cell: (props) => (
      <InvoicePreviewCell value={props.getValue()?.value?.invoiceNumber} />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),
  columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesIssueDate,
    size: 120,
    enableColumnFilter: false,
    enableSorting: true,
    header: (props) => (
      <THead id='issueDate' title='Issue date' {...getTHeadProps(props)} />
    ),
    cell: (props) => {
      return (
        <p>
          {DateTimeUtils.format(
            props?.getValue()?.value?.issued,
            DateTimeUtils.defaultFormatShortString,
          )}
        </p>
      );
    },
    skeleton: () => <Skeleton className='w-[50px] h-[18px]' />,
  }),

  columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesInvoiceStatus,
    size: 110,
    enableColumnFilter: false,
    enableSorting: true,
    header: (props) => (
      <THead
        id='paymentStatus'
        title='Payment status'
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => (
      <PaymentStatusSelect invoiceNumber={props.getValue()?.number} />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),
  columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesAmount,
    size: 90,
    enableColumnFilter: false,
    enableSorting: false,
    header: (props) => (
      <THead id='amount' title='Amount' {...getTHeadProps(props)} />
    ),
    cell: (props) => {
      return (
        <p>
          {formatCurrency(
            props?.getValue()?.value?.amountDue,
            2,
            props.getValue()?.value?.currency,
          )}
        </p>
      );
    },
    skeleton: () => <Skeleton className='w-[50px] h-[18px]' />,
  }),
];
