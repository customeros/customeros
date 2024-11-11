import { Building07 } from '@ui/media/icons/Building07';
import {
  InternalStage,
  ColumnViewType,
  ComparisonOperator,
} from '@shared/types/__generated__/graphql.types';

export type FilterType = {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  options?: any[];
  icon: JSX.Element;
  filterName: string;
  filterAccesor: ColumnViewType;
  filterOperators: ComparisonOperator[];
  filterType: 'text' | 'date' | 'number' | 'list';
  groupOptions?: { label: string; options: { id: string; label: string }[] }[];
};

import { type RootStore } from '@store/root';

import { Key01 } from '@ui/media/icons/Key01';
import { Calendar } from '@ui/media/icons/Calendar';
import { Columns03 } from '@ui/media/icons/Columns03';
import { ClockCheck } from '@ui/media/icons/ClockCheck';
import { Calculator } from '@ui/media/icons/Calculator';
import { ArrowsRight } from '@ui/media/icons/ArrowsRight';
import { CoinsStacked01 } from '@ui/media/icons/CoinsStacked01';

export const getFilterTypes = (store?: RootStore) => {
  const preset = store?.tableViewDefs.opportunitiesPreset;
  const stageLabels = store?.tableViewDefs
    .getById(preset ?? '1')
    ?.value?.columns.reduce((acc, curr) => {
      if (!curr.filter?.includes('STAGE')) return acc;

      return {
        ...acc,
        ['STAGE' + curr.filter.split('STAGE')[1][0]]: curr.name,
        [InternalStage.ClosedLost]: 'Closed Lost',
        [InternalStage.ClosedWon]: 'Closed Won',
      };
    }, {} as Record<string, string>);

  const filterTypes: Partial<Record<ColumnViewType, FilterType>> = {
    [ColumnViewType.OpportunitiesName]: {
      filterType: 'text',
      filterName: 'Opportunity name',
      filterAccesor: ColumnViewType.OpportunitiesName,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <CoinsStacked01 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.OpportunitiesOrganization]: {
      filterType: 'text',
      filterName: 'Organization name',
      filterAccesor: ColumnViewType.OpportunitiesOrganization,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Building07 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.OpportunitiesStage]: {
      filterType: 'list',
      filterName: 'Stage',
      filterAccesor: ColumnViewType.OpportunitiesStage,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Columns03 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
      options: Object.entries(stageLabels ?? {}).map(([key, value]) => ({
        id: key,
        label: value,
      })),
    },

    [ColumnViewType.OpportunitiesNextStep]: {
      filterType: 'text',
      filterName: 'Next step',
      filterAccesor: ColumnViewType.OpportunitiesNextStep,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <ArrowsRight className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },

    [ColumnViewType.OpportunitiesEstimatedArr]: {
      filterType: 'number',
      filterName: 'Estimated ARR',
      filterAccesor: ColumnViewType.OpportunitiesEstimatedArr,
      filterOperators: [
        ComparisonOperator.Gt,
        ComparisonOperator.Lt,
        ComparisonOperator.Eq,
        ComparisonOperator.NotEqual,
      ],
      icon: (
        <Calculator className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.OpportunitiesOwner]: {
      filterType: 'list',
      filterName: 'Owner',
      filterAccesor: ColumnViewType.OpportunitiesOwner,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: <Key01 className='grou-hover:text-gray-700 text-gray-500' />,
      options: store?.users.toArray().map((user) => ({
        id: user?.id,
        label: user?.name,
        avatar: user?.value?.profilePhotoUrl,
      })),
    },
    [ColumnViewType.OpportunitiesTimeInStage]: {
      filterType: 'number',
      filterName: 'Time in stage',
      filterAccesor: ColumnViewType.OpportunitiesTimeInStage,
      filterOperators: [
        ComparisonOperator.Gt,
        ComparisonOperator.Lt,
        ComparisonOperator.Eq,
        ComparisonOperator.NotEqual,
      ],
      icon: (
        <ClockCheck className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.OpportunitiesCreatedDate]: {
      filterType: 'date',
      filterName: 'Created',
      filterAccesor: ColumnViewType.OpportunitiesCreatedDate,
      filterOperators: [ComparisonOperator.Lt, ComparisonOperator.Gt],
      icon: (
        <Calendar className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
  };

  return filterTypes;
};
