import { match } from 'ts-pattern';

import { Equal } from '@ui/media/icons/Equal';
import { Button } from '@ui/form/Button/Button';
import { ChevronLeft } from '@ui/media/icons/ChevronLeft';
import { CheckCircle } from '@ui/media/icons/CheckCircle';
import { CubeOutline } from '@ui/media/icons/CubeOutline';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { CalendarAfter } from '@ui/media/icons/CalendarAfter';
import { SpacingWidth01 } from '@ui/media/icons/SpacingWidth01';
import { CalendarBefore } from '@ui/media/icons/CalendarBefore';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { ComparisonOperator } from '@shared/types/__generated__/graphql.types';

interface OperatorFilterProps {
  type: string;
  value: string;
  operators: string[];
  onSelect: (operator: string) => void;
}

export const OperatorFilter = ({
  operators,
  onSelect,
  value,
}: OperatorFilterProps) => {
  const defaultOperator = handleOperatorName(
    operators[0] as ComparisonOperator,
  );

  return (
    <Menu>
      <MenuButton>
        <Button
          size='xs'
          colorScheme='grayModern'
          className='border-transparent'
        >
          {handleOperatorName(value as ComparisonOperator) ?? defaultOperator}
        </Button>
      </MenuButton>
      <MenuList>
        {operators?.map((operator) => (
          <MenuItem
            key={operator}
            className='group'
            onClick={() => onSelect(operator)}
          >
            <div className='flex items-center gap-2'>
              <span className='text-gray-500 group-hover:text-gray-700'>
                {handleOperatorIcon(operator as ComparisonOperator, 'date')}
              </span>
              {handleOperatorName(operator as ComparisonOperator, 'date')}
            </div>
          </MenuItem>
        ))}
      </MenuList>
    </Menu>
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

const handleOperatorIcon = (operator: ComparisonOperator, type?: string) => {
  return match(operator)
    .with(ComparisonOperator.Between, () => <SpacingWidth01 />)
    .with(ComparisonOperator.Contains, () => <CheckCircle />)
    .with(ComparisonOperator.Eq, () => <Equal />)
    .with(ComparisonOperator.Gt, () =>
      type === 'date' ? <CalendarAfter /> : <ChevronRight />,
    )
    .with(ComparisonOperator.Gte, () => 'greater than or equal to')
    .with(ComparisonOperator.In, () => 'in')
    .with(ComparisonOperator.IsEmpty, () => <CubeOutline />)
    .with(ComparisonOperator.IsNull, () => 'is null')
    .with(ComparisonOperator.Lt, () =>
      type === 'date' ? <CalendarBefore /> : <ChevronLeft />,
    )
    .with(ComparisonOperator.Lte, () => 'less than or equal to')
    .with(ComparisonOperator.StartsWith, () => 'starts with')
    .otherwise(() => 'unknown');
};
