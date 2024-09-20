import { match } from 'ts-pattern';

import { ButtonGroup } from '@ui/form/ButtonGroup';
import { ComparisonOperator } from '@shared/types/__generated__/graphql.types';

import { PropertyFilter } from './PropertyFilter';
import { OperatorFilter } from './OperatorFilter';
import { ValueFilter } from './ValueFilter/ValueFilter';
import { ClearFilter } from './ClearFilter/ClearFilter';

interface FilterProps {
  filterName: string;
  filterType: string;
  operators: string[];
  filterValue: string;
  onClearFilter: () => void;
  operatorValue: ComparisonOperator;
  onChangeOperator: (operator: string) => void;
  onChangeFilterValue: (value: string) => void;
}

export const Filter = ({
  onChangeOperator,
  operators,
  operatorValue,
  filterType,
  filterName,
  onChangeFilterValue,
  filterValue,
  onClearFilter,
}: FilterProps) => {
  return (
    <ButtonGroup className='flex items-center'>
      <PropertyFilter name={filterName} />
      <OperatorFilter
        type={filterType}
        value={operatorValue}
        operators={operators}
        onSelect={onChangeOperator}
      />
      <ValueFilter
        filterType={filterType}
        filterName={filterName}
        filterValue={filterValue}
        onChangeFilterValue={onChangeFilterValue}
        operatorValue={handleOperatorName(operatorValue)}
      />
      <ClearFilter onClearFilter={onClearFilter} />
    </ButtonGroup>
  );
};

const handleOperatorName = (operator: ComparisonOperator, type?: string) => {
  return match(operator)
    .with(ComparisonOperator.Between, () => 'between')
    .with(ComparisonOperator.Contains, () => 'contains')
    .with(ComparisonOperator.Eq, () => 'equals')
    .with(ComparisonOperator.Gt, () =>
      type === 'date' ? 'after' : 'more than',
    )
    .with(ComparisonOperator.Gte, () => 'greater than or equal to')
    .with(ComparisonOperator.In, () => 'in')
    .with(ComparisonOperator.IsEmpty, () => 'is empty')
    .with(ComparisonOperator.IsNull, () => 'is null')
    .with(ComparisonOperator.Lt, () =>
      type === 'date' ? 'before' : 'less than',
    )
    .with(ComparisonOperator.Lte, () => 'less than or equal to')
    .with(ComparisonOperator.StartsWith, () => 'starts with')
    .otherwise(() => 'unknown');
};
