import { ButtonGroup } from '@ui/form/ButtonGroup';

import { PropertyFilter } from './PropertyFilter';
import { OperatorFilter } from './OperatorFilter';
import { ValueFilter } from './ValueFilter/ValueFilter';
import { ClearFilter } from './ClearFilter/ClearFilter';

interface FilterProps {
  filterName: string;
  filterType: string;
  operators: string[];
  operatorValue: string;
  onChangeOperator: (operator: string) => void;
}

export const Filter = ({
  onChangeOperator,
  operators,
  operatorValue,
  filterType,
  filterName,
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
      <ValueFilter />
      <ClearFilter />
    </ButtonGroup>
  );
};
