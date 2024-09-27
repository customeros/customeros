import React from 'react';

import { Store } from '@store/store.ts';
import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';
import { TextCell } from '@finder/components/Columns/shared/Cells/TextCell';
import { OwnerFilter } from '@finder/components/Columns/shared/Filters/Owner';
import { ForecastFilter } from '@finder/components/Columns/shared/Filters/Forecast';
import { currencyIcon } from '@settings/components/Tabs/panels/BillingPanel/components/utils.tsx';
import { RenewalLikelihoodFilter } from '@finder/components/Columns/shared/Filters/RenewalLikelihood';
import {
  LtvFilter,
  DateFilter,
  StatusFilter,
  PeriodFilter,
  RenewalFilter,
  CurrencyFilter,
} from '@finder/components/Columns/contracts/Filters';

import { DateTimeUtils } from '@utils/date.ts';
import { createColumnHelper } from '@ui/presentation/Table';
import { Skeleton } from '@ui/feedback/Skeleton/Skeleton.tsx';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';
import {
  Contract,
  Opportunity,
  TableViewDef,
  ColumnViewType,
} from '@graphql/types';

import { getColumnConfig } from '../shared/util/getColumnConfig';
import { SearchTextFilter } from '../shared/Filters/SearchTextFilter';
import {
  LtvCell,
  OwnerCell,
  HealthCell,
  PeriodCell,
  StatusCell,
  ContractCell,
  ArrForecastCell,
} from './Cells';

type ColumnDatum = Store<Contract>;

// REASON: we do not care about exhaustively typing this TValue type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Column = ColumnDefinition<ColumnDatum, any>;

const columnHelper = createColumnHelper<ColumnDatum>();

const columns: Record<string, Column> = {
  [ColumnViewType.ContractsName]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsName,
    minSize: 160,
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
        DateTimeUtils.dateWithAbreviatedMonth,
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
    enableColumnFilter: true,
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
          <PeriodFilter initialFocusRef={initialFocusRef} />
        )}
        {...getTHeadProps<Store<Contract>>(props)}
      />
    ),
  }),

  [ColumnViewType.ContractsCurrency]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsCurrency,
    minSize: 120,
    maxSize: 350,
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
    minSize: 150,
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
              : 'Not auto-renewing'
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
  [ColumnViewType.ContractsOwner]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsOwner,
    minSize: 150,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: true,
    enableSorting: true,
    cell: (props) => {
      return <OwnerCell id={props.getValue()?.value?.metadata?.id} />;
    },
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,

    header: (props) => (
      <THead<HTMLInputElement>
        title='Owner'
        filterWidth='14rem'
        id={ColumnViewType.ContractsOwner}
        renderFilter={(initialFocusRef) => (
          <OwnerFilter
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContractsOwner}
          />
        )}
        {...getTHeadProps<Store<Contract>>(props)}
      />
    ),
  }),
  [ColumnViewType.ContractsRenewalDate]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsRenewalDate,
    minSize: 150,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: true,
    enableSorting: true,
    cell: (props) => {
      const renewsAt = props
        .getValue()
        ?.value?.opportunities?.find(
          (e: Opportunity) => e.internalStage === 'OPEN',
        )?.renewedAt;

      if (!renewsAt) {
        return <p className='text-gray-400'>Unknown</p>;
      }
      const formatted = DateTimeUtils.format(
        renewsAt,
        DateTimeUtils.dateWithAbreviatedMonth,
      );

      return <TextCell text={formatted} />;
    },
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,

    header: (props) => (
      <THead<HTMLInputElement>
        filterWidth='14rem'
        title='Renewal Date'
        id={ColumnViewType.ContractsRenewalDate}
        renderFilter={() => (
          <DateFilter property={ColumnViewType.ContractsRenewalDate} />
        )}
        {...getTHeadProps<Store<Contract>>(props)}
      />
    ),
  }),
  [ColumnViewType.ContractsHealth]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsHealth,
    minSize: 150,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: true,
    enableSorting: true,
    cell: (props) => {
      return <HealthCell id={props.getValue()?.value?.metadata?.id} />;
    },
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,

    header: (props) => (
      <THead<HTMLInputElement>
        title='Health'
        filterWidth='14rem'
        id={ColumnViewType.ContractsHealth}
        renderFilter={() => (
          <RenewalLikelihoodFilter property={ColumnViewType.ContractsHealth} />
        )}
        {...getTHeadProps<Store<Contract>>(props)}
      />
    ),
  }),
  [ColumnViewType.ContractsForecastArr]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsForecastArr,
    minSize: 150,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: true,
    enableSorting: true,
    cell: (props) => {
      return <ArrForecastCell id={props.getValue()?.value?.metadata?.id} />;
    },
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,

    header: (props) => (
      <THead<HTMLInputElement>
        filterWidth='14rem'
        title='ARR Forecast'
        id={ColumnViewType.ContractsForecastArr}
        renderFilter={(initialFocusRef) => (
          <ForecastFilter
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.ContractsForecastArr}
          />
        )}
        {...getTHeadProps<Store<Contract>>(props)}
      />
    ),
  }),
};
export const getContractColumnsConfig = (
  tableViewDef?: Array<TableViewDef>[0],
) => getColumnConfig<ColumnDatum>(columns, tableViewDef);
