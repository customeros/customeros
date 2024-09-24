import { Button } from '@ui/form/Button/Button';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { ComparisonOperator } from '@shared/types/__generated__/graphql.types';

import { handleOperatorIcon, handleOperatorName } from '../utils/utils';

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
          variant='outline'
          colorScheme='grayModern'
          className='rounded-none font-normal bg-white text-gray-500'
        >
          {value
            ? handleOperatorName(value as ComparisonOperator)
            : defaultOperator}
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
