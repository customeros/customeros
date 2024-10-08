import { FlowStore } from '@store/Flows/Flow.store.ts';
import { StatusFilter } from '@finder/components/Columns/flows/Filters';
import {
  ColumnDef,
  ColumnDef as ColumnDefinition,
} from '@tanstack/react-table';
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
  [ColumnViewType.FlowName]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.FlowName,
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
        id={ColumnViewType.FlowName}
        {...getTHeadProps(props)}
      />
    ),
    cell: (props) => (
      <FlowNameCell id={props.row?.original?.value?.metadata?.id ?? ''} />
    ),
    skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  }),
  // Temporary: Will be removed and replaced with FlowStatus
  [ColumnViewType.FlowActionName]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.FlowActionName,
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
        id={ColumnViewType.FlowActionName}
        renderFilter={() => <StatusFilter />}
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

  [ColumnViewType.FlowPendingCount]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.FlowPendingCount,
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
        id={ColumnViewType.FlowPendingCount}
        {...getTHeadProps(props)}
      />
    ),
    cell: (cell) => {
      const statistics = cell.getValue()?.value?.statistics;

      return (
        <FlowStatisticsCell
          total={statistics.total}
          value={statistics.pending}
          dataTest={'flow-in-progress'}
        />
      );
    },
    skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  }),
  [ColumnViewType.FlowGoalAchievedCount]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.FlowGoalAchievedCount,
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
        id={ColumnViewType.FlowGoalAchievedCount}
        {...getTHeadProps(props)}
      />
    ),
    cell: (cell) => {
      const statistics = cell.getValue()?.value?.statistics;

      return (
        <FlowStatisticsCell
          total={statistics.total}
          dataTest='flow-goal-achieved'
          value={statistics.goalAchieved}
        />
      );
    },
    skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  }),
  // [ColumnViewType.FlowActionTotalCount]: columnHelper.accessor((row) => row, {
  //   id: ColumnViewType.FlowActionTotalCount,
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
  //       id={ColumnViewType.FlowActionTotalCount}
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
  [ColumnViewType.FlowCompletedCount]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.FlowCompletedCount,
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
        id={ColumnViewType.FlowCompletedCount}
        {...getTHeadProps(props)}
      />
    ),
    cell: (cell) => {
      const statistics = cell.getValue()?.value?.statistics;

      return (
        <FlowStatisticsCell
          total={statistics.total}
          dataTest='flow-completed'
          value={statistics.completed}
        />
      );
    },
    skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  }),
};

export const getFlowColumnsConfig = (
  tableViewDef?: Array<TableViewDef>[0],
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
): ColumnDef<ColumnDatum, any>[] =>
  getColumnConfig<ColumnDatum>(columns, tableViewDef);
