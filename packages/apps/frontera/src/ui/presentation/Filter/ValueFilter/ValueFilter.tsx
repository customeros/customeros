import { useState } from 'react';

import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover/Popover';

interface ValueFilterProps {
  filterName: string;
  filterType: string;
  filterValue: string;
  operatorValue: string;
  onChangeFilterValue: (value: string) => void;
}

export const ValueFilter = ({
  filterType,
  filterName,
  operatorValue,
  onChangeFilterValue,
  filterValue,
}: ValueFilterProps) => {
  const [isOpen, setIsOpen] = useState(false);
  const [inputValue, setInputValue] = useState(filterValue);

  const handleSetFilter = (value: boolean) => {
    setIsOpen(value);

    if (!value) {
      onChangeFilterValue(inputValue);
    }
  };

  return (
    <Popover open={isOpen} onOpenChange={(value) => handleSetFilter(value)}>
      <PopoverTrigger>
        <Button
          size='xs'
          colorScheme='grayModern'
          className='border-transparent rounded-none text-gray-400'
        >
          {filterValue ? filterValue : '...'}
        </Button>
      </PopoverTrigger>
      <PopoverContent className='py-1 w-[240px]'>
        <Input
          size='sm'
          variant='unstyled'
          value={inputValue}
          placeholder={`${filterName} ${operatorValue}`}
          onChange={(e) => setInputValue(e.target.value)}
        />
      </PopoverContent>
    </Popover>
  );
};
