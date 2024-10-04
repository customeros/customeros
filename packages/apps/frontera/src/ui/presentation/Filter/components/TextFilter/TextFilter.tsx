import { useState, useEffect } from 'react';

import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { ComparisonOperator } from '@shared/types/__generated__/graphql.types';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover/Popover';

import { handleOperatorName } from '../../utils/utils';

interface TextFilterProps {
  filterName: string;
  filterValue: string;
  operatorValue: string;
  onChangeFilterValue: (value: string) => void;
}

export const TextFilter = ({
  filterName,
  operatorValue,
  onChangeFilterValue,
  filterValue,
}: TextFilterProps) => {
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

  if (
    operatorValue === ComparisonOperator.IsEmpty ||
    operatorValue === ComparisonOperator.IsNotEmpty
  )
    return;

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
      <PopoverTrigger asChild>
        <Button
          size='xs'
          colorScheme='grayModern'
          onClick={() => setIsOpen(!isOpen)}
          className='rounded-none text-gray-700 bg-white font-normal boreder-r-2'
        >
          <span className=' max-w-[160px] text-ellipsis whitespace-nowrap overflow-hidden'>
            {filterValue ? filterValue : '...'}
          </span>
        </Button>
      </PopoverTrigger>
      <PopoverContent
        side='bottom'
        align='start'
        className='py-1 min-w-[254px]'
      >
        <Input
          size='sm'
          variant='unstyled'
          value={inputValue}
          onChange={handleInputChange}
          placeholder={`${filterName} ${handleOperatorName(
            operatorValue as ComparisonOperator,
          )}`}
          onKeyDown={(e) => {
            if (e.key === 'Escape') {
              setIsOpen(false);
            }
          }}
        />
      </PopoverContent>
    </Popover>
  );
};
