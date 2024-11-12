import { OpportunityStore } from '@store/Opportunities/Opportunity.store';
import { DateCell } from '@finder/components/Columns/shared/Cells/DateCell';
import { TextCell } from '@finder/components/Columns/shared/Cells/TextCell';
import {
  ColumnDef,
  ColumnDef as ColumnDefinition,
} from '@tanstack/react-table';
import { OrganizationCell } from '@finder/components/Columns/shared/Cells/organization';
import { getColumnConfig } from '@finder/components/Columns/shared/util/getColumnConfig.ts';

import { DateTimeUtils } from '@utils/date.ts';
import { createColumnHelper } from '@ui/presentation/Table';
import { Skeleton } from '@ui/feedback/Skeleton/Skeleton.tsx';
import { TableViewDef, ColumnViewType } from '@graphql/types';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead.tsx';

import { OpportunityName } from './Cells/opportunityName';
import { OwnerCell, StageCell, ArrEstimateCell } from './Cells';

type ColumnDatum = OpportunityStore;

// REASON: we do not care about exhaustively typing this TValue type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Column = ColumnDefinition<ColumnDatum, any>;

const columnHelper = createColumnHelper<ColumnDatum>();

export const columns: Record<string, Column> = {
  [ColumnViewType.OpportunitiesName]: columnHelper.accessor((row) => row, {
    id: ColumnViewType.OpportunitiesName,
    minSize: 160,
    size: 160,
    maxSize: 400,
    enableColumnFilter: false,
    enableSorting: true,
    enableResizing: true,
    cell: (props) => {
      const name = props.row.original.value.name;

      if (!name) return <div className='text-gray-400'>Unnamed</div>;

      return <OpportunityName opportunityId={props.row.original.value.id} />;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Name'
        id={ColumnViewType.OpportunitiesName}
        {...getTHeadProps<OpportunityStore>(props)}
      />
    ),
    skeleton: () => <Skeleton className='w-[100px] h-[14px]' />,
  }),
  [ColumnViewType.OpportunitiesOrganization]: columnHelper.accessor(
    (row) => row,
    {
      id: ColumnViewType.OpportunitiesOrganization,
      minSize: 125,
      maxSize: 400,
      enableColumnFilter: false,
      enableResizing: true,
      enableSorting: true,
      cell: (props) => {
        if (!props.row.original.value.organization?.metadata?.id) {
          return <p className='text-gray-400'>Unknown</p>;
        }

        return (
          <OrganizationCell
            id={props.row.original.value.organization.metadata.id}
          />
        );
      },
      header: (props) => (
        <THead<HTMLInputElement>
          title='Organization'
          id={ColumnViewType.OpportunitiesOrganization}
          {...getTHeadProps<OpportunityStore>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[50%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OpportunitiesStage]: columnHelper.accessor(
    'value.externalStage',
    {
      id: ColumnViewType.OpportunitiesStage,
      minSize: 160,
      maxSize: 400,
      enableColumnFilter: false,
      enableResizing: true,
      enableSorting: true,
      header: (props) => (
        <THead
          title='Stage'
          id={ColumnViewType.OpportunitiesStage}
          {...getTHeadProps<OpportunityStore>(props)}
        />
      ),
      cell: (props) => {
        return (
          <StageCell stage={props.getValue()} id={props.row.original.id} />
        );
      },
      skeleton: () => <Skeleton className='w-[100%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OpportunitiesEstimatedArr]: columnHelper.accessor(
    'value.metadata.id',
    {
      id: ColumnViewType.OpportunitiesEstimatedArr,
      minSize: 125,
      maxSize: 400,
      enableColumnFilter: false,
      enableResizing: true,
      enableSorting: true,
      cell: (props) => {
        return <ArrEstimateCell opportunityId={props.getValue()} />;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          title='ARR Estimate'
          id={ColumnViewType.OpportunitiesEstimatedArr}
          {...getTHeadProps<OpportunityStore>(props)}
        />
      ),
      skeleton: () => (
        <div className='flex flex-col gap-1'>
          <Skeleton className='w-[33%] h-[14px]' />
        </div>
      ),
    },
  ),
  [ColumnViewType.OpportunitiesOwner]: columnHelper.accessor('value.owner', {
    id: ColumnViewType.OpportunitiesOwner,
    minSize: 110,
    size: 110,
    maxSize: 400,
    enableColumnFilter: false,
    enableResizing: true,
    enableSorting: true,
    cell: (props) => {
      return (
        <OwnerCell
          owner={props.getValue()}
          id={props.row.original.value.metadata?.id}
        />
      );
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Owner'
        data-test='owner'
        id={ColumnViewType.OpportunitiesOwner}
        {...getTHeadProps<OpportunityStore>(props)}
      />
    ),
    skeleton: () => (
      <div className='flex flex-col gap-1'>
        <Skeleton className='w-[25%] h-[14px]' />
      </div>
    ),
  }),
  [ColumnViewType.OpportunitiesTimeInStage]: columnHelper.accessor(
    'value.stageLastUpdated',
    {
      id: ColumnViewType.OpportunitiesTimeInStage,
      minSize: 140,
      size: 140,
      maxSize: 400,
      enableColumnFilter: false,
      enableResizing: true,
      enableSorting: true,
      cell: (props) => {
        const timeInStage = props.getValue()
          ? DateTimeUtils.getDistanceToNowStrict(props.getValue(), 'day')
          : '';

        return <TextCell text={timeInStage} />;
      },

      header: (props) => (
        <THead<HTMLInputElement>
          title='Time in Stage'
          id={ColumnViewType.OpportunitiesTimeInStage}
          {...getTHeadProps<OpportunityStore>(props)}
        />
      ),
      skeleton: () => <Skeleton className='w-[50%] h-[14px]' />,
    },
  ),
  [ColumnViewType.OpportunitiesCreatedDate]: columnHelper.accessor(
    'value.metadata.created',
    {
      id: ColumnViewType.OpportunitiesCreatedDate,
      minSize: 154,
      size: 154,
      maxSize: 400,
      enableColumnFilter: false,
      enableResizing: true,
      enableSorting: true,
      cell: (props) => {
        return <DateCell value={props.getValue()} />;
      },
      header: (props) => (
        <THead
          title='Created'
          id={ColumnViewType.OpportunitiesCreatedDate}
          {...getTHeadProps<OpportunityStore>(props)}
        />
      ),
      skeleton: () => (
        <div className='flex flex-col gap-1'>
          <Skeleton className='w-[50%] h-[14px]' />
          <Skeleton className='w-[25%] h-[14px]' />
        </div>
      ),
    },
  ),
  [ColumnViewType.OpportunitiesNextStep]: columnHelper.accessor(
    'value.nextSteps',
    {
      id: ColumnViewType.OpportunitiesNextStep,
      minSize: 154,
      size: 154,
      maxSize: 400,
      enableColumnFilter: false,
      enableResizing: true,
      enableSorting: true,
      cell: (props) => {
        return (
          <TextCell text={props.getValue()} unknownText={'No next step'} />
        );
      },
      header: (props) => (
        <THead<HTMLInputElement>
          title='Next Steps'
          id={ColumnViewType.OpportunitiesNextStep}
          {...getTHeadProps<OpportunityStore>(props)}
        />
      ),
      skeleton: () => (
        <div className='flex flex-col gap-1'>
          <Skeleton className='w-[50%] h-[14px]' />
          <Skeleton className='w-[25%] h-[14px]' />
        </div>
      ),
    },
  ),
};

export const getOpportunityColumnsConfig = (
  tableViewDef?: Array<TableViewDef>[0],
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
): ColumnDef<ColumnDatum, any>[] =>
  getColumnConfig<ColumnDatum>(columns, tableViewDef);
