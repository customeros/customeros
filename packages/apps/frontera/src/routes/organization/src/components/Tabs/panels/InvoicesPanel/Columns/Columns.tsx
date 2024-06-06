import { Store } from '@store/store.ts';

import { Invoice } from '@graphql/types';
import { DateTimeUtils } from '@utils/date';
import { Skeleton } from '@ui/feedback/Skeleton';
import { createColumnHelper } from '@ui/presentation/Table';
import { StatusCell } from '@shared/components/Invoice/Cells';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';

type ColumnDatum = Store<Invoice>;

//@typescript-eslint/no-explicit-any
const columnHelper = createColumnHelper<ColumnDatum>();

export const columns = [
  columnHelper.accessor((row) => row, {
    id: 'NUMBER',
    minSize: 90,
    maxSize: 90,
    enableSorting: false,
    enableColumnFilter: false,
    cell: (props) => (
      <p className='overflow-hidden'>
        {props?.getValue()?.value?.invoiceNumber}
      </p>
    ),
    header: (props) => (
      <THead
        id='number'
        title='N°'
        py='1'
        {...getTHeadProps<Store<Invoice>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[20px] h-[18px]' />,
  }),

  columnHelper.accessor((row) => row, {
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
        {...getTHeadProps<Store<Invoice>>(props)}
      />
    ),
    cell: (props) => {
      return <StatusCell status={props?.getValue()?.value?.status} />;
    },
    skeleton: () => <Skeleton className='w-[100%] h-[18px]' />,
  }),
  columnHelper.accessor((row) => row, {
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
        {...getTHeadProps<Store<Invoice>>(props)}
      />
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
    skeleton: () => <Skeleton className='w-[100%] h-[18px]' />,
  }),
  columnHelper.accessor((row) => row, {
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
        {...getTHeadProps<Store<Invoice>>(props)}
      />
    ),
    cell: (props) => {
      return (
        <p className='text-center'>
          {formatCurrency(
            props?.getValue()?.value?.amountDue,
            2,
            props.getValue()?.value?.currency,
          )}
        </p>
      );
    },
    skeleton: () => <Skeleton className='w-[100%] h-[18px]' />,
  }),
];
