import { ButtonGroup } from '@ui/form/ButtonGroup';
import { ComparisonOperator } from '@shared/types/__generated__/graphql.types';

import { DateFilter } from './components/DateFilter';
import { PropertyFilter } from './components/PropertyFilter';
import { OperatorFilter } from './components/OperatorFilter';
import { TextFilter } from './components/TextFilter/TextFilter';
import { ClearFilter } from './components/ClearFilter/ClearFilter';

interface FilterProps {
  filterName: string;
  filterType: string;
  operators: string[];
  filterValue: string;
  onClearFilter: () => void;
  operatorValue: ComparisonOperator;
  onChangeOperator: (operator: string) => void;
  onChangeFilterValue: (value: string | Date) => void;
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
      {filterType === 'text' && (
        <TextFilter
          filterName={filterName}
          filterValue={filterValue}
          operatorValue={operatorValue}
          onChangeFilterValue={onChangeFilterValue}
        />
      )}

      {filterType === 'date' && (
        <DateFilter
          filterName={filterName}
          operatorValue={operatorValue}
          filterValue={filterValue as string | [string | null, string | null]}
          onChangeFilterValue={
            onChangeFilterValue as (
              value: string | [string | null, string | null],
            ) => void
          }
        />
      )}

      <ClearFilter onClearFilter={onClearFilter} />
    </ButtonGroup>
  );
};
