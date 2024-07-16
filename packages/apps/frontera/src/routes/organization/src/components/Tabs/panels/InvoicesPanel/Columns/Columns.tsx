import { InvoiceStore } from '@store/Invoices/Invoice.store.ts';

import { DateTimeUtils } from '@utils/date';
import { ColumnViewType } from '@graphql/types';
import { Skeleton } from '@ui/feedback/Skeleton';
import { createColumnHelper } from '@ui/presentation/Table';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';

import { PaymentStatusSelect } from '../../../../../../../invoices/src/components/shared';
import { InvoicePreviewCell } from '../../../../../../../invoices/src/components/Columns/Cells';

type ColumnDatum = InvoiceStore;

const columnHelper = createColumnHelper<ColumnDatum>();

export const columns = [
  columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesInvoiceNumber,
    size: 100,
    enableColumnFilter: false,
    enableSorting: false,
    header: (props) => (
      <THead
        id='invoiceNumber'
        title='Invoice number'
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
  columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesIssueDate,
    size: 110,
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
    id: ColumnViewType.InvoicesPaymentStatus,
    size: 120,
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
      <PaymentStatusSelect
        value={props.getValue()?.value?.status || null}
        invoiceId={props.getValue()?.value?.metadata?.id}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[18px]' />,
  }),
  columnHelper.accessor((row) => row, {
    id: ColumnViewType.InvoicesAmount,
    size: 100,
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
