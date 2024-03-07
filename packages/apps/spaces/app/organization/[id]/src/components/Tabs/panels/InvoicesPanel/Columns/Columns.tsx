import React from 'react';

import { Invoice } from '@graphql/types';
import { Text } from '@ui/typography/Text';
import { DateTimeUtils } from '@spaces/utils/date';
import { Skeleton } from '@ui/presentation/Skeleton';
import { createColumnHelper } from '@ui/presentation/Table';
import { StatusCell } from '@shared/components/Invoice/Cells';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';

const columnHelper = createColumnHelper<Invoice>();

export const columns = [
  columnHelper.accessor('invoiceNumber', {
    id: 'NUMBER',
    minSize: 90,
    maxSize: 90,
    enableSorting: false,
    enableColumnFilter: false,
    cell: (props) => <Text overflow='hidden'>{props?.getValue()}</Text>,
    header: (props) => (
      <THead id='number' title='N°' py='1' {...getTHeadProps<Invoice>(props)} />
    ),
    skeleton: () => (
      <Skeleton
        width='20%'
        height='18px'
        startColor='gray.300'
        endColor='gray.300'
      />
    ),
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
    skeleton: () => (
      <Skeleton
        width='100%'
        height='18px'
        startColor='gray.300'
        endColor='gray.300'
      />
    ),
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
        <Text>
          {DateTimeUtils.format(
            props.getValue().created,
            DateTimeUtils.defaultFormatShortString,
          )}
        </Text>
      );
    },
    skeleton: () => (
      <Skeleton
        width='100%'
        height='18px'
        startColor='gray.300'
        endColor='gray.300'
      />
    ),
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
        <Text textAlign='center'>
          {formatCurrency(props.getValue(), 2, props.row.original.currency)}
        </Text>
      );
    },
    skeleton: () => (
      <Skeleton
        width='100%'
        height='18px'
        startColor='gray.300'
        endColor='gray.300'
      />
    ),
  }),
];
