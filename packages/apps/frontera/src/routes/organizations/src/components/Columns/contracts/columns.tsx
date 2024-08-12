import React from 'react';

import { ContractStore } from '@store/Contracts/Contract.store';
import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';

import { DateTimeUtils } from '@utils/date.ts';
import { createColumnHelper } from '@ui/presentation/Table';
import { TableViewDef, ColumnViewType } from '@graphql/types';
import { Skeleton } from '@ui/feedback/Skeleton/Skeleton.tsx';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';
import { PeriodCell } from '@organizations/components/Columns/contracts/period';
import { StatusCell } from '@organizations/components/Columns/contracts/status';
import { TextCell } from '@organizations/components/Columns/shared/Cells/TextCell';

import { getColumnConfig } from '../shared/util/getColumnConfig';
import { SearchTextFilter } from '../shared/Filters/SearchTextFilter';

type ColumnDatum = ContractStore;

// REASON: we do not care about exhaustively typing this TValue type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Column = ColumnDefinition<ColumnDatum, any>;

const columnHelper = createColumnHelper<ColumnDatum>();

const columns: Record<string, Column> = {
  [ColumnViewType.ContractsName]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsName,
    minSize: 150,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    cell: (props) => {
      return (
        <TextCell {...props} text={props.getValue()?.value?.contractName} />
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        filterWidth='14rem'
        title='  Contract Name'
        id={ColumnViewType.ContactsName}
        renderFilter={(initialFocusRef) => (
          <SearchTextFilter
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContractsName}
            placeholder={'e.g. CustomerOS contract'}
          />
        )}
        {...getTHeadProps<ContractStore>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,
  }),

  [ColumnViewType.ContractsEnded]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsEnded,
    minSize: 150,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    cell: (props) => {
      const contractEnded = props.getValue()?.value?.contractEnded;

      if (!contractEnded) {
        return <p className='text-gray-400'>Unknown</p>;
      }
      const formatted = DateTimeUtils.format(
        contractEnded,
        DateTimeUtils.defaultFormatShortString,
      );

      return <TextCell {...props} text={formatted} />;
    },
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,

    header: (props) => (
      <THead<HTMLInputElement>
        title='Ended'
        filterWidth='14rem'
        id={ColumnViewType.ContractsEnded}
        renderFilter={(initialFocusRef) => (
          <SearchTextFilter
            placeholder={'e.g. Yes/No'}
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContractsEnded}
          />
        )}
        {...getTHeadProps<ContractStore>(props)}
      />
    ),
  }),

  [ColumnViewType.ContractsPeriod]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsPeriod,
    minSize: 150,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    cell: (props) => {
      const committedPeriodInMonths =
        props.getValue()?.value?.committedPeriodInMonths;

      return <PeriodCell committedPeriodInMonths={committedPeriodInMonths} />;
    },
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,

    header: (props) => (
      <THead<HTMLInputElement>
        title='Period'
        filterWidth='14rem'
        id={ColumnViewType.ContractsPeriod}
        renderFilter={(initialFocusRef) => (
          <SearchTextFilter
            placeholder={'e.g. 1 year'}
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContractsPeriod}
          />
        )}
        {...getTHeadProps<ContractStore>(props)}
      />
    ),
  }),

  [ColumnViewType.ContractsSignDate]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsSignDate,
    minSize: 150,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    cell: (props) => {
      const contractEnded = props.getValue()?.value?.contractSigned;

      if (!contractEnded) {
        return <p className='text-gray-400'>Unknown</p>;
      }
      const formatted = DateTimeUtils.format(
        contractEnded,
        DateTimeUtils.defaultFormatShortString,
      );

      return <TextCell {...props} text={formatted} />;
    },
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,
    header: (props) => (
      <THead<HTMLInputElement>
        title='Sign Date'
        filterWidth='14rem'
        id={ColumnViewType.ContractsSignDate}
        renderFilter={(initialFocusRef) => (
          <SearchTextFilter
            placeholder={'e.g. 2023-01-01'}
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContractsSignDate}
          />
        )}
        {...getTHeadProps<ContractStore>(props)}
      />
    ),
  }),

  [ColumnViewType.ContractsCurrency]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsCurrency,
    minSize: 150,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    cell: (props) => {
      return <TextCell {...props} text={props.getValue().value.currency} />;
    },
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,

    header: (props) => (
      <THead<HTMLInputElement>
        title='Currency'
        filterWidth='14rem'
        id={ColumnViewType.ContractsCurrency}
        renderFilter={(initialFocusRef) => (
          <SearchTextFilter
            placeholder={'e.g. USD'}
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContractsCurrency}
          />
        )}
        {...getTHeadProps<ContractStore>(props)}
      />
    ),
  }),

  [ColumnViewType.ContractsStatus]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsStatus,
    minSize: 100,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    cell: (props) => {
      return <StatusCell status={props.getValue().value.contractStatus} />;
    },
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,

    header: (props) => (
      <THead<HTMLInputElement>
        title='Status'
        filterWidth='14rem'
        id={ColumnViewType.ContractsStatus}
        renderFilter={(initialFocusRef) => (
          <SearchTextFilter
            placeholder={'e.g. Active'}
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContractsStatus}
          />
        )}
        {...getTHeadProps<ContractStore>(props)}
      />
    ),
  }),

  [ColumnViewType.ContractsRenewal]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsRenewal,
    minSize: 150,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    cell: (props) => {
      return (
        <TextCell
          text={
            props.getValue().value.autoRenew
              ? 'Auto-renews'
              : 'Non auto-renewing'
          }
        />
      );
    },
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,
    header: (props) => (
      <THead<HTMLInputElement>
        title='Renewal'
        filterWidth='14rem'
        id={ColumnViewType.ContractsRenewal}
        renderFilter={(initialFocusRef) => (
          <SearchTextFilter
            placeholder={'e.g. Auto-renews'}
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContractsRenewal}
          />
        )}
        {...getTHeadProps<ContractStore>(props)}
      />
    ),
  }),

  [ColumnViewType.ContractsLtv]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsLtv,
    minSize: 150,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    cell: (props) => {
      return <TextCell {...props} text={props.getValue().value.ltv} />;
    },
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,

    header: (props) => (
      <THead<HTMLInputElement>
        title='LTV'
        filterWidth='14rem'
        id={ColumnViewType.ContractsLtv}
        renderFilter={(initialFocusRef) => (
          <SearchTextFilter
            placeholder={'e.g. 10000'}
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContractsLtv}
          />
        )}
        {...getTHeadProps<ContractStore>(props)}
      />
    ),
  }),
};
export const getContractColumnsConfig = (
  tableViewDef?: Array<TableViewDef>[0],
) => getColumnConfig<ColumnDatum>(columns, tableViewDef);
