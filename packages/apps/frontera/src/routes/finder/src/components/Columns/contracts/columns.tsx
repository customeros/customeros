import { ContractStore } from '@store/Contracts/Contract.store';
import { TextCell } from '@finder/components/Columns/shared/Cells/TextCell';
import {
  ColumnDef,
  ColumnDef as ColumnDefinition,
} from '@tanstack/react-table';
import { currencyIcon } from '@settings/components/Tabs/panels/BillingPanel/components/utils.tsx';

import { DateTimeUtils } from '@utils/date.ts';
import { createColumnHelper } from '@ui/presentation/Table';
import { Skeleton } from '@ui/feedback/Skeleton/Skeleton.tsx';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';
import { Opportunity, TableViewDef, ColumnViewType } from '@graphql/types';

import { getColumnConfig } from '../shared/util/getColumnConfig';
import {
  LtvCell,
  OwnerCell,
  HealthCell,
  PeriodCell,
  StatusCell,
  ContractCell,
  ArrForecastCell,
} from './Cells';

type ColumnDatum = ContractStore;

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
    enableColumnFilter: false,
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
        {...getTHeadProps<ContractStore>(props)}
      />
    ),
  }),

  [ColumnViewType.ContractsCurrency]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsCurrency,
    minSize: 120,
    maxSize: 350,
    enableResizing: true,
    enableColumnFilter: false,
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
        {...getTHeadProps<ContractStore>(props)}
      />
    ),
  }),
  [ColumnViewType.ContractsOwner]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsOwner,
    minSize: 150,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: false,
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
        {...getTHeadProps<ContractStore>(props)}
      />
    ),
  }),
  [ColumnViewType.ContractsRenewalDate]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsRenewalDate,
    minSize: 150,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: false,
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
        {...getTHeadProps<ContractStore>(props)}
      />
    ),
  }),
  [ColumnViewType.ContractsHealth]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsHealth,
    minSize: 150,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: false,
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
        {...getTHeadProps<ContractStore>(props)}
      />
    ),
  }),
  [ColumnViewType.ContractsForecastArr]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.ContractsForecastArr,
    minSize: 150,
    maxSize: 650,
    enableResizing: true,
    enableColumnFilter: false,
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
        {...getTHeadProps<ContractStore>(props)}
      />
    ),
  }),
};
export const getContractColumnsConfig = (
  tableViewDef?: Array<TableViewDef>[0],
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
): ColumnDef<ColumnDatum, any>[] =>
  getColumnConfig<ColumnDatum>(columns, tableViewDef);
