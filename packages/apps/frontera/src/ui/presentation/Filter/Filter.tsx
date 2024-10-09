import { ButtonGroup } from '@ui/form/ButtonGroup';
import { ComparisonOperator } from '@shared/types/__generated__/graphql.types';

import { DateFilter } from './components/DateFilter';
import { ListFilter } from './components/ListFilter';
import { ClearFilter } from './components/ClearFilter';
import { NumberFilter } from './components/NumberFilter';
import { PropertyFilter } from './components/PropertyFilter';
import { OperatorFilter } from './components/OperatorFilter';
import { TextFilter } from './components/TextFilter/TextFilter';

interface FilterProps {
  filterName: string;
  filterType: string;
  operators: string[];
  icon: React.ReactElement;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  listFilterOptions: any[];
  onClearFilter: () => void;
  filterValue: string | string[];
  operatorValue: ComparisonOperator;
  onChangeOperator: (operator: string) => void;
  onChangeFilterValue: (value: string | Date | string[]) => void;
  groupOptions?: { label: string; options: { id: string; label: string }[] }[];
}

export const Filter = ({
  onChangeOperator,
  operators,
  operatorValue,
  filterType,
  filterName,
  onChangeFilterValue,
  filterValue,
  listFilterOptions,
  groupOptions,
  onClearFilter,
  icon,
}: FilterProps) => {
  return (
    <ButtonGroup className='flex items-center'>
      <PropertyFilter icon={icon} name={filterName} />
      <OperatorFilter
        type={filterType}
        value={operatorValue}
        operators={operators}
        onSelect={onChangeOperator}
        isOperatorPlural={filterValue?.length > 1}
      />
      {filterType === 'text' && (
        <TextFilter
          filterName={filterName}
          operatorValue={operatorValue}
          filterValue={filterValue as string}
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

      {filterType === 'number' && (
        <NumberFilter
          filterName={filterName}
          operatorValue={operatorValue}
          filterValue={filterValue as string}
          onChangeFilterValue={
            onChangeFilterValue as (
              value: string | [string | null | number, string | null | number],
            ) => void
          }
        />
      )}

      {filterType === 'list' && (
        <ListFilter
          filterName={filterName}
          options={listFilterOptions}
          operatorName={operatorValue}
          groupOptions={groupOptions || []}
          filterValue={filterValue as string[]}
          onMultiSelectChange={onChangeFilterValue as (ids: string[]) => void}
        />
      )}

      <ClearFilter onClearFilter={onClearFilter} />
    </ButtonGroup>
  );
};
