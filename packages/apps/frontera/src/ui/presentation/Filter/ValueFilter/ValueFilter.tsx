import { useState, useEffect } from 'react';

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

  useEffect(() => {
    if (!filterValue) {
      if (filterName) {
        setTimeout(() => {
          setIsOpen(true);
        }, 100);
      }
    }
  }, [filterName]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;

    setInputValue(newValue);
    onChangeFilterValue(newValue);
  };

  return (
    <Popover
      modal={true}
      open={isOpen}
      onOpenChange={(value) => setIsOpen(value)}
    >
      <PopoverTrigger>
        <Button
          size='xs'
          colorScheme='grayModern'
          onClick={() => setIsOpen(!isOpen)}
          className='border-l-0 rounded-none text-gray-700 bg-white font-normal'
        >
          {filterValue ? filterValue : '...'}
        </Button>
      </PopoverTrigger>
      <PopoverContent className='py-1 w-[240px]'>
        <Input
          size='sm'
          variant='unstyled'
          value={inputValue}
          onChange={handleInputChange}
          placeholder={`${filterName} ${operatorValue}`}
        />
      </PopoverContent>
    </Popover>
  );
};
