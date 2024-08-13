import React from 'react';

import { Store } from '@store/store.ts';
import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';
import { currencyIcon } from '@settings/components/Tabs/panels/BillingPanel/components/utils.tsx';

import { DateTimeUtils } from '@utils/date.ts';
import { createColumnHelper } from '@ui/presentation/Table';
import { Skeleton } from '@ui/feedback/Skeleton/Skeleton.tsx';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';
import { Contract, TableViewDef, ColumnViewType } from '@graphql/types';
import { TextCell } from '@organizations/components/Columns/shared/Cells/TextCell';
import {
  LtvFilter,
  DateFilter,
  StatusFilter,
  RenewalFilter,
  CurrencyFilter,
} from '@organizations/components/Columns/contracts/Filters';

import { getColumnConfig } from '../shared/util/getColumnConfig';
import { SearchTextFilter } from '../shared/Filters/SearchTextFilter';
import { LtvCell, PeriodCell, StatusCell, ContractCell } from './Cells';

type ColumnDatum = Store<Contract>;

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
    enableColumnFilter: true,
    enableSorting: true,
    cell: (props) => {
      return (
        <ContractCell contractId={props.getValue()?.value?.metadata?.id} />
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        filterWidth='14rem'
        title='  Contract Name'
        id={ColumnViewType.ContractsName}
        renderFilter={(initialFocusRef) => (
          <SearchTextFilter
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContractsName}
            placeholder={'e.g. CustomerOS contract'}
          />
        )}
        {...getTHeadProps<Store<Contract>>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,
  }),

  [ColumnViewType.ContractsEnded]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsEnded,
    minSize: 150,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: true,
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

      return <TextCell text={formatted} />;
    },
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,

    header: (props) => (
      <THead<HTMLInputElement>
        title='Ended'
        filterWidth='14rem'
        id={ColumnViewType.ContractsEnded}
        renderFilter={() => (
          <DateFilter property={ColumnViewType.ContractsEnded} />
        )}
        {...getTHeadProps<Store<Contract>>(props)}
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
        {...getTHeadProps<Store<Contract>>(props)}
      />
    ),
  }),

  [ColumnViewType.ContractsSignDate]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsSignDate,
    minSize: 150,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: true,
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

      return <TextCell text={formatted} />;
    },
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,
    header: (props) => (
      <THead<HTMLInputElement>
        title='Sign Date'
        filterWidth='14rem'
        id={ColumnViewType.ContractsSignDate}
        renderFilter={() => (
          <DateFilter property={ColumnViewType.ContractsSignDate} />
        )}
        {...getTHeadProps<Store<Contract>>(props)}
      />
    ),
  }),

  [ColumnViewType.ContractsCurrency]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsCurrency,
    minSize: 150,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: true,
    enableSorting: true,
    cell: (props) => {
      const currency = props.getValue().value.currency;

      return (
        <TextCell text={currency} leftIcon={currencyIcon?.[currency || '']} />
      );
    },
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,

    header: (props) => (
      <THead<HTMLInputElement>
        title='Currency'
        filterWidth='14rem'
        id={ColumnViewType.ContractsCurrency}
        renderFilter={() => <CurrencyFilter />}
        {...getTHeadProps<Store<Contract>>(props)}
      />
    ),
  }),

  [ColumnViewType.ContractsStatus]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsStatus,
    minSize: 100,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: true,
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
        renderFilter={() => <StatusFilter />}
        {...getTHeadProps<Store<Contract>>(props)}
      />
    ),
  }),

  [ColumnViewType.ContractsRenewal]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsRenewal,
    minSize: 230,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: true,
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
        renderFilter={() => <RenewalFilter />}
        {...getTHeadProps<Store<Contract>>(props)}
      />
    ),
  }),

  [ColumnViewType.ContractsLtv]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsLtv,
    minSize: 150,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: true,
    enableSorting: true,
    cell: (props) => {
      return (
        <LtvCell
          ltv={props.getValue()?.value?.ltv}
          currency={props.getValue()?.value?.currency}
        />
      );
    },
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,

    header: (props) => (
      <THead<HTMLInputElement>
        title='LTV'
        filterWidth='14rem'
        id={ColumnViewType.ContractsLtv}
        renderFilter={(initialFocusRef) => (
          <LtvFilter initialFocusRef={initialFocusRef} />
        )}
        {...getTHeadProps<Store<Contract>>(props)}
      />
    ),
  }),
};
export const getContractColumnsConfig = (
  tableViewDef?: Array<TableViewDef>[0],
) => getColumnConfig<ColumnDatum>(columns, tableViewDef);
