import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { ComparisonOperator } from '@shared/types/__generated__/graphql.types';

import { handleOperatorIcon, handleOperatorName } from '../../utils/utils';

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
  type,
}: OperatorFilterProps) => {
  const defaultOperator = handleOperatorName(
    operators[0] as ComparisonOperator,
    type,
  );

  return (
    <Menu>
      <MenuButton>
        {type === 'number' ? (
          <IconButton
            size='xs'
            variant='outline'
            colorScheme='grayModern'
            className='rounded-none bg-white'
            aria-label={`filter type ${type}`}
            icon={
              value
                ? (handleOperatorIcon(
                    value as ComparisonOperator,
                    type,
                  ) as React.ReactElement)
                : (handleOperatorIcon(
                    operators[0] as ComparisonOperator,
                    type,
                  ) as React.ReactElement)
            }
          />
        ) : (
          <Button
            size='xs'
            variant='outline'
            colorScheme='grayModern'
            className='rounded-none font-normal bg-white text-gray-500'
          >
            {value
              ? handleOperatorName(value as ComparisonOperator, type)
              : defaultOperator}
          </Button>
        )}
      </MenuButton>
      <MenuList>
        {operators?.map((operator) => (
          <MenuItem
            key={operator}
            className='group'
            onClick={() => onSelect(operator)}
          >
            <div className='flex items-center gap-2'>
              <span>
                {handleOperatorIcon(operator as ComparisonOperator, type)}
              </span>
              {handleOperatorName(operator as ComparisonOperator, type)}
            </div>
          </MenuItem>
        ))}
      </MenuList>
    </Menu>
  );
};
