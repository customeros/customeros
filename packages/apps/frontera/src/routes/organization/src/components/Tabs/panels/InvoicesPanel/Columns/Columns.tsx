import { Invoice } from '@graphql/types';
import { DateTimeUtils } from '@utils/date';
import { Skeleton } from '@ui/feedback/Skeleton';
import { createColumnHelper } from '@ui/presentation/Table';
import { StatusCell } from '@shared/components/Invoice/Cells';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';

const columnHelper = createColumnHelper<Invoice>();

export const columns = [
  columnHelper.accessor('invoiceNumber', {
    id: 'NUMBER',
    minSize: 90,
    maxSize: 90,
    enableSorting: false,
    enableColumnFilter: false,
    cell: (props) => <p className='overflow-hidden'>{props?.getValue()}</p>,
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
      return <StatusCell status={props.getValue()} />;
    },
    skeleton: () => <Skeleton className='w-[100%] h-[18px]' />,
  }),
  columnHelper.accessor('metadata', {
    id: 'DATE',
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
            props.getValue().created,
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
          {formatCurrency(props.getValue(), 2, props.row.original.currency)}
        </p>
      );
    },
    skeleton: () => <Skeleton className='w-[100%] h-[18px]' />,
  }),
];
