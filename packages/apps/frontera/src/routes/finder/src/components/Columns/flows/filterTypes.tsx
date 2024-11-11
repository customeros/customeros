import {
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
  groupOptions?: { label: string; options: { id: string; label: string }[] };
};

import { Shuffle01 } from '@ui/media/icons/Shuffle01';

export const getFilterTypes = () => {
  const filterTypes: Partial<Record<ColumnViewType, FilterType>> = {
    [ColumnViewType.FlowActionName]: {
      filterType: 'text',
      filterName: 'Flow name',
      filterAccesor: ColumnViewType.FlowActionName,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Shuffle01 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
  };

  return filterTypes;
};
