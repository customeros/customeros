import React from 'react';

import { Store } from '@store/store.ts';
import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';
import { DateFilter } from '@finder/components/Columns/contracts/Filters';
import { DateCell } from '@finder/components/Columns/shared/Cells/DateCell';
import { TextCell } from '@finder/components/Columns/shared/Cells/TextCell';
import { StageFilter } from '@finder/components/Columns/opportunities/Filter';
import { OwnerFilter } from '@finder/components/Columns/shared/Filters/Owner';
import { ForecastFilter } from '@finder/components/Columns/shared/Filters/Forecast';
import { OrganizationCell } from '@finder/components/Columns/shared/Cells/organization';
import { getColumnConfig } from '@finder/components/Columns/shared/util/getColumnConfig.ts';
import { SearchTextFilter } from '@finder/components/Columns/shared/Filters/SearchTextFilter';
import { NumericValueFilter } from '@finder/components/Columns/shared/Filters/NumericValueFilter';

import { DateTimeUtils } from '@utils/date.ts';
import { createColumnHelper } from '@ui/presentation/Table';
import { Skeleton } from '@ui/feedback/Skeleton/Skeleton.tsx';
import THead, { getTHeadProps } from '@ui/presentation/Table/THead.tsx';
import { Opportunity, TableViewDef, ColumnViewType } from '@graphql/types';

import { OwnerCell, StageCell, ArrEstimateCell } from './Cells';

type ColumnDatum = Store<Opportunity>;

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
    enableColumnFilter: true,
    enableSorting: true,
    enableResizing: true,
    cell: (props) => {
      const name = props.row.original.value.name;

      if (!name) return <div className='text-gray-400'>Unnamed</div>;

      return <p className='font-medium'>{name}</p>;
    },
    header: (props) => (
      <THead<HTMLInputElement>
        title='Name'
        filterWidth='14rem'
        id={ColumnViewType.OpportunitiesName}
        renderFilter={(initialFocusRef) => (
          <SearchTextFilter
            initialFocusRef={initialFocusRef}
            property={ColumnViewType.OpportunitiesName}
          />
        )}
        {...getTHeadProps<Store<Opportunity>>(props)}
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
      enableColumnFilter: true,
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
          filterWidth='14rem'
          title='Organization'
          id={ColumnViewType.OpportunitiesOrganization}
          renderFilter={(initialFocusRef) => (
            <SearchTextFilter
              initialFocusRef={initialFocusRef}
              property={ColumnViewType.OpportunitiesOrganization}
            />
          )}
          {...getTHeadProps<Store<Opportunity>>(props)}
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
      enableColumnFilter: true,
      enableResizing: true,
      enableSorting: true,
      header: (props) => (
        <THead
          title='Stage'
          renderFilter={() => <StageFilter />}
          id={ColumnViewType.OpportunitiesStage}
          {...getTHeadProps<Store<Opportunity>>(props)}
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
      enableColumnFilter: true,
      enableResizing: true,
      enableSorting: true,
      cell: (props) => {
        return <ArrEstimateCell opportunityId={props.getValue()} />;
      },
      header: (props) => (
        <THead<HTMLInputElement>
          title='ARR Estimate'
          id={ColumnViewType.OpportunitiesEstimatedArr}
          renderFilter={(initialFocusRef) => (
            <ForecastFilter
              initialFocusRef={initialFocusRef}
              property={ColumnViewType.OpportunitiesEstimatedArr}
            />
          )}
          {...getTHeadProps<Store<Opportunity>>(props)}
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
    enableColumnFilter: true,
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
        renderFilter={(initialFocusedRef) => (
          <OwnerFilter
            initialFocusRef={initialFocusedRef}
            property={ColumnViewType.OpportunitiesOwner}
          />
        )}
        {...getTHeadProps<Store<Opportunity>>(props)}
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
      enableColumnFilter: true,
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
          filterWidth='15rem'
          title='Time in Stage'
          id={ColumnViewType.OpportunitiesTimeInStage}
          renderFilter={(initialFocusRef) => (
            <NumericValueFilter
              label='days'
              suffix='day'
              initialFocusRef={initialFocusRef}
              property={ColumnViewType.OpportunitiesTimeInStage}
            />
          )}
          {...getTHeadProps<Store<Opportunity>>(props)}
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
      enableColumnFilter: true,
      enableResizing: true,
      enableSorting: true,
      cell: (props) => {
        return <DateCell value={props.getValue()} />;
      },
      header: (props) => (
        <THead
          title='Created'
          filterWidth='17rem'
          id={ColumnViewType.OpportunitiesCreatedDate}
          renderFilter={() => (
            <DateFilter property={ColumnViewType.OpportunitiesCreatedDate} />
          )}
          {...getTHeadProps<Store<Opportunity>>(props)}
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
      enableColumnFilter: true,
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
          filterWidth='17rem'
          id={ColumnViewType.OpportunitiesNextStep}
          renderFilter={(initialFocusRef) => (
            <SearchTextFilter
              initialFocusRef={initialFocusRef}
              property={ColumnViewType.OpportunitiesNextStep}
            />
          )}
          {...getTHeadProps<Store<Opportunity>>(props)}
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
) => getColumnConfig<ColumnDatum>(columns, tableViewDef);
