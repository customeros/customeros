import { FlowStore } from '@store/Flows/Flow.store.ts';
import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';
import { StatusFilter } from '@finder/components/Columns/flows/Filters';
import { TextCell } from '@finder/components/Columns/shared/Cells/TextCell';
import { getColumnConfig } from '@finder/components/Columns/shared/util/getColumnConfig';
import {
  FlowNameCell,
  FlowStatusCell,
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
  // Todo uncomment when necessary, for now should stay hidden
  // [ColumnViewType.FlowName]: columnHelper.accessor((row) => row, {
  //   id: ColumnViewType.FlowName,
  //   size: 150,
  //   minSize: 150,
  //   maxSize: 300,
  //   enableResizing: true,
  //   enableColumnFilter: false,
  //   enableSorting: true,
  //   header: (props) => (
  //     <THead
  //       title='Sequence'
  //       filterWidth={250}
  //       id={ColumnViewType.FlowName}
  //       {...getTHeadProps(props)}
  //     />
  //   ),
  //   cell: (props) => <TextCell text={props.row?.original?.value?.name ?? ''} />,
  //   skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  // }),

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

  // [ColumnViewType.FlowSequenceContactCount]: columnHelper.accessor(
  //   (row) => row,
  //   {
  //     id: ColumnViewType.FlowSequenceContactCount,
  //     size: 150,
  //     minSize: 150,
  //     maxSize: 300,
  //     enableResizing: true,
  //     enableColumnFilter: false,
  //     enableSorting: true,
  //     header: (props) => (
  //       <THead
  //         title='Contacts'
  //         filterWidth={250}
  //         id={ColumnViewType.FlowSequenceContactCount}
  //         {...getTHeadProps(props)}
  //       />
  //     ),
  //     cell: (props) => (
  //       <TextCell
  //         text={props.row?.original?.value?.contacts?.length?.toString() ?? ''}
  //       />
  //     ),
  //     skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
  //   },
  // ),
  [ColumnViewType.FlowSequenceStatusInProgressCount]: columnHelper.accessor(
    (row) => row,
    {
      id: ColumnViewType.FlowSequenceStatusInProgressCount,
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
          id={ColumnViewType.FlowSequenceStatusInProgressCount}
          {...getTHeadProps(props)}
        />
      ),
      cell: () => (
        <TextCell
          text={''}
          unknownText='No data yet'
          dataTest='flow-in-progress-in-flows-table'
        />
      ),
      skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
    },
  ),
  [ColumnViewType.FlowSequenceStatusPendingCount]: columnHelper.accessor(
    (row) => row,
    {
      id: ColumnViewType.FlowSequenceStatusPendingCount,
      size: 150,
      minSize: 150,
      maxSize: 300,
      enableResizing: true,
      enableColumnFilter: false,
      enableSorting: true,
      header: (props) => (
        <THead
          filterWidth={250}
          title='Not Started'
          id={ColumnViewType.FlowSequenceStatusPendingCount}
          {...getTHeadProps(props)}
        />
      ),
      cell: () => (
        <TextCell
          text={''}
          unknownText='No data yet'
          dataTest={'flow-not-started-in-flows-table'}
        />
      ),
      skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
    },
  ),
  [ColumnViewType.FlowSequenceStatusSuccessfulCount]: columnHelper.accessor(
    (row) => row,
    {
      id: ColumnViewType.FlowSequenceStatusSuccessfulCount,
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
          id={ColumnViewType.FlowSequenceStatusSuccessfulCount}
          {...getTHeadProps(props)}
        />
      ),
      cell: () => (
        <TextCell
          text={''}
          unknownText='No data yet'
          dataTest='flow-completed-in-flows-table'
        />
      ),
      skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
    },
  ),
  [ColumnViewType.FlowSequenceStatusUnsuccessfulCount]: columnHelper.accessor(
    (row) => row,
    {
      id: ColumnViewType.FlowSequenceStatusUnsuccessfulCount,
      size: 150,
      minSize: 150,
      maxSize: 300,
      enableResizing: true,
      enableColumnFilter: false,
      enableSorting: true,
      header: (props) => (
        <THead
          filterWidth={250}
          title='Ended Early'
          id={ColumnViewType.FlowSequenceStatusUnsuccessfulCount}
          {...getTHeadProps(props)}
        />
      ),
      cell: () => (
        <TextCell
          text={''}
          unknownText='No data yet'
          dataTest={'flow-ended-early-in-flows-table'}
        />
      ),
      skeleton: () => <Skeleton className='w-[200px] h-[18px]' />,
    },
  ),
};

export const getFlowColumnsConfig = (tableViewDef?: Array<TableViewDef>[0]) =>
  getColumnConfig<ColumnDatum>(columns, tableViewDef);
