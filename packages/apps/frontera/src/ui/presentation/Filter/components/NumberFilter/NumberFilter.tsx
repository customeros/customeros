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

interface NumberFilterProps {
  filterName: string;
  filterValue: string;
  operatorValue: string;
  onChangeFilterValue: (value: string | [number, number]) => void;
}

const formatNumberWithCommas = (value: string | number | undefined): string => {
  if (value === undefined || value === '' || value === 0) return '';
  const numString = value?.toString();

  return numString?.replace(/\B(?=(\d{3})+(?!\d))/g, ',');
};

export const NumberFilter = ({
  filterName,
  operatorValue,
  onChangeFilterValue,
  filterValue,
}: NumberFilterProps) => {
  const [isOpen, setIsOpen] = useState(false);
  const [inputValue, setInputValue] = useState(filterValue);
  const [minValue, setMinValue] = useState<number | string>('');
  const [maxValue, setMaxValue] = useState<number | string>('');

  useEffect(() => {
    if (!filterValue) {
      if (filterName) {
        setTimeout(() => {
          setIsOpen(true);
        }, 100);
      }
    }
  }, [filterName]);

  useEffect(() => {
    if (operatorValue === ComparisonOperator.Between) {
      onChangeFilterValue([Number(minValue), Number(maxValue)]);
    }
  }, [minValue, maxValue, operatorValue]);

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

  const handleMinValueChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;

    setMinValue(Number(newValue));
  };

  const handleMaxValueChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;

    setMaxValue(Number(newValue));
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
          className='rounded-none text-gray-700 bg-white font-normal border-l-0'
        >
          <span className=' max-w-[160px] text-ellipsis whitespace-nowrap overflow-hidden'>
            {formatNumberWithCommas(filterValue) || '...'}
          </span>
        </Button>
      </PopoverTrigger>
      <PopoverContent
        side='bottom'
        align='start'
        className='py-1 min-w-[254px]'
      >
        {ComparisonOperator.Between === operatorValue ? (
          <>
            <div className='flex items-center'>
              <Input
                size='sm'
                type='number'
                value={minValue}
                placeholder='Min'
                variant='unstyled'
                onChange={handleMinValueChange}
              />
              <Input
                size='sm'
                type='number'
                value={maxValue}
                placeholder='Max'
                variant='unstyled'
                onChange={handleMaxValueChange}
              />
            </div>
          </>
        ) : (
          <>
            <Input
              size='sm'
              type='number'
              variant='unstyled'
              value={inputValue}
              onChange={handleInputChange}
              placeholder={`${filterName} ${handleOperatorName(
                operatorValue as ComparisonOperator,
              )}`}
            />
          </>
        )}
      </PopoverContent>
    </Popover>
  );
};
