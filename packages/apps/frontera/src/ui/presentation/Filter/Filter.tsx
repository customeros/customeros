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
        operatorValue={operatorValue}
        onChangeFilterValue={onChangeFilterValue}
      />
      <ClearFilter onClearFilter={onClearFilter} />
    </ButtonGroup>
  );
};
