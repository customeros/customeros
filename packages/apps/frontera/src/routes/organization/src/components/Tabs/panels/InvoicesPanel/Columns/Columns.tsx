import { Store } from '@store/store.ts';
import { ColumnDef as ColumnDefinition } from '@tanstack/table-core/build/lib/types';

import { Invoice } from '@graphql/types';
import { DateTimeUtils } from '@utils/date';
import { Skeleton } from '@ui/feedback/Skeleton';
import { createColumnHelper } from '@ui/presentation/Table';
import { StatusCell } from '@shared/components/Invoice/Cells';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';
type ColumnDatum = Store<Invoice>;
// REASON: we do not care about exhaustively typing this TValue type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Column = ColumnDefinition<ColumnDatum, any>;
const columnHelper = createColumnHelper<ColumnDatum>();

export const columns: Record<string, Column> = [
  columnHelper.accessor('invoiceNumber', {
    id: 'NUMBER',
    minSize: 90,
    maxSize: 90,
    enableSorting: false,
    enableColumnFilter: false,
    cell: (props) => (
      <p className='overflow-hidden'>
        {props?.row.original?.value?.invoiceNumber}
      </p>
    ),
    header: (props) => (
      <THead id='number' title='N°' py='1' {...getTHeadProps<Invoice>(props)} />
    ),
    skeleton: () => <Skeleton className='w-[20px] h-[18px]' />,
  }),

  columnHelper.accessor('status', {
    id: 'STATUS',
    minSize: 105,
    maxSize: 105,
    enableSorting: false,
    enableColumnFilter: false,
    header: (props) => (
      <THead
        id='status'
        title='Status'
        py='1'
        {...getTHeadProps<Invoice>(props)}
      />
    ),
    cell: (props) => {
      return <StatusCell status={props?.row.original?.value?.status} />;
    },
    skeleton: () => <Skeleton className='w-[100%] h-[18px]' />,
  }),
  columnHelper.accessor('issued', {
    id: 'DATE_ISSUED',
    minSize: 10,
    maxSize: 10,
    enableSorting: false,
    enableColumnFilter: false,
    header: (props) => (
      <THead
        id='issued'
        title='Issued'
        py='1'
        {...getTHeadProps<Invoice>(props)}
      />
    ),
    cell: (props) => {
      return (
        <p>
          {DateTimeUtils.format(
            props?.row.original?.value?.issued,
            DateTimeUtils.defaultFormatShortString,
          )}
        </p>
      );
    },
    skeleton: () => <Skeleton className='w-[100%] h-[18px]' />,
  }),
  columnHelper.accessor('amountDue', {
    id: 'AMOUNT_DUE',
    minSize: 20,
    maxSize: 20,
    enableSorting: false,
    enableColumnFilter: false,
    header: (props) => (
      <THead
        id='amount'
        title='Amount'
        py='1'
        {...getTHeadProps<Invoice>(props)}
      />
    ),
    cell: (props) => {
      return (
        <p className='text-center'>
          {formatCurrency(
            props?.row.original?.value?.amountDue,
            2,
            props.row.original.currency,
          )}
        </p>
      );
    },
    skeleton: () => <Skeleton className='w-[100%] h-[18px]' />,
  }),
];
