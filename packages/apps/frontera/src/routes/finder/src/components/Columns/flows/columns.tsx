import { FlowStore } from '@store/Flows/Flow.store.ts';
import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';
import { StatusFilter } from '@finder/components/Columns/flows/Filters';
import { getColumnConfig } from '@finder/components/Columns/shared/util/getColumnConfig';
import {
  FlowNameCell,
  FlowStatusCell,
  FlowStatisticsCell,
} from '@finder/components/Columns/flows/Cells';

import { Skeleton } from '@ui/feedback/Skeleton';
import { createColumnHelper } from '@ui/presentation/Table';
import { TableViewDef, ColumnViewType } from '@graphql/types';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead';

type ColumnDatum = FlowStore;

// REASON: we do not care about exhaustively typing this TValue type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Column = ColumnDefinition<ColumnDatum, any>;

const columnHelper = createColumnHelper<ColumnDatum>();

const columns: Record<string, Column> = {
  [ColumnViewType.FlowSequenceName]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.FlowSequenceName,
    size: 150,
    minSize: 150,
    maxSize: 300,
    enableResizing: true,
    enableColumnFilter: false,
    enableSorting: true,
    header: (props) => (
      <THead
        title='Flow'
        filterWidth={250}
        id={ColumnViewType.FlowSequenceName}
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => (
      <FlowNameCell id={props.row?.original?.value?.metadata?.id ?? ''} />
    ),
    skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  }),
  [ColumnViewType.FlowSequenceStatus]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.FlowSequenceStatus,
    size: 150,
    minSize: 150,
    maxSize: 300,
    enableResizing: true,
    enableColumnFilter: true,
    enableSorting: true,
    header: (props) => (
      <THead
        title='Status'
        filterWidth={250}
        renderFilter={() => <StatusFilter />}
        id={ColumnViewType.FlowSequenceStatus}
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => (
      <FlowStatusCell
        dataTest={'flow-status'}
        id={props.row?.original?.value?.metadata?.id ?? ''}
      />
    ),
    skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  }),

  [ColumnViewType.FlowSequencePendingCount]: columnHelper.accessor(
    (row) => row,
    {
      id: ColumnViewType.FlowSequencePendingCount,
      size: 150,
      minSize: 150,
      maxSize: 300,
      enableResizing: true,
      enableColumnFilter: false,
      enableSorting: true,
      header: (props) => (
        <THead
          filterWidth={250}
          title='In Progress'
          id={ColumnViewType.FlowSequencePendingCount}
          {...getTHeadProps(props)}
        />
      ),
      cell: (cell) => {
        const statistics = cell.getValue()?.value?.statistics;

        return (
          <FlowStatisticsCell
            total={statistics.total}
            value={statistics.pending}
          />
        );
      },
      skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
    },
  ),
  [ColumnViewType.FlowSequenceGoalAchievedCount]: columnHelper.accessor(
    (row) => row,
    {
      id: ColumnViewType.FlowSequenceGoalAchievedCount,
      size: 150,
      minSize: 150,
      maxSize: 300,
      enableResizing: true,
      enableColumnFilter: false,
      enableSorting: true,
      header: (props) => (
        <THead
          filterWidth={250}
          title='Goal Achieved'
          id={ColumnViewType.FlowSequenceGoalAchievedCount}
          {...getTHeadProps(props)}
        />
      ),
      cell: (cell) => {
        const statistics = cell.getValue()?.value?.statistics;

        return (
          <FlowStatisticsCell
            total={statistics.total}
            value={statistics.goalAchieved}
          />
        );
      },
      skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
    },
  ),
  // [ColumnViewType.FlowSequenceTotalCount]: columnHelper.accessor((row) => row, {
  //   id: ColumnViewType.FlowSequenceTotalCount,
  //   size: 150,
  //   minSize: 150,
  //   maxSize: 300,
  //   enableResizing: true,
  //   enableColumnFilter: false,
  //   enableSorting: true,
  //   header: (props) => (
  //     <THead
  //       title='Total '
  //       filterWidth={250}
  //       id={ColumnViewType.FlowSequenceTotalCount}
  //       {...getTHeadProps(props)}
  //     />
  //   ),
  //   cell: (e) => {
  //     const total = e.getValue()?.value?.statistics?.total;
  //
  //     return (
  //       <TextCell
  //         text={`${total}`}
  //         unknownText='No data yet'
  //         dataTest='flow-completed-in-flows-table'
  //       />
  //     );
  //   },
  //   skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  // }),
  [ColumnViewType.FlowSequenceCompletedCount]: columnHelper.accessor(
    (row) => row,
    {
      id: ColumnViewType.FlowSequenceCompletedCount,
      size: 150,
      minSize: 150,
      maxSize: 300,
      enableResizing: true,
      enableColumnFilter: false,
      enableSorting: true,
      header: (props) => (
        <THead
          filterWidth={250}
          title='Completed'
          id={ColumnViewType.FlowSequenceCompletedCount}
          {...getTHeadProps(props)}
        />
      ),
      cell: (cell) => {
        const statistics = cell.getValue()?.value?.statistics;

        return (
          <FlowStatisticsCell
            total={statistics.total}
            value={statistics.completed}
          />
        );
      },
      skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
    },
  ),
};

export const getFlowColumnsConfig = (tableViewDef?: Array<TableViewDef>[0]) =>
  getColumnConfig<ColumnDatum>(columns, tableViewDef);
