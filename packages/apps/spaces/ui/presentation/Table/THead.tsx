'use client';
import type { HeaderContext } from '@tanstack/react-table';

import { memo, useRef, RefObject } from 'react';

import { cn } from '@ui/utils/cn';
import { ArrowUp } from '@ui/media/icons/ArrowUp';
import { ArrowDown } from '@ui/media/icons/ArrowDown';
import { FilterLines } from '@ui/media/icons/FilterLines';
import { IconButton } from '@ui/form/IconButton/IconButton';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover/Popover';

import { useTHeadState } from './THead.atom';

interface THeadProps<
  InitialRefType extends { focus(): void } = HTMLButtonElement,
> {
  id: string;
  py?: string;
  title: string;
  padding?: string;
  subTitle?: string;
  canSort?: boolean;
  canFilter?: boolean;
  isFiltered?: boolean;
  borderTopColor?: string;
  isSorted?: string | boolean;
  filterWidth?: string | number;
  onToggleSort?: (e: unknown) => void;
  renderFilter?: (
    initialFocusRef: RefObject<InitialRefType>,
  ) => React.ReactNode;
}

const THead = <InitialRefType extends { focus(): void } = HTMLButtonElement>({
  id,
  title,
  canSort,
  isSorted,
  subTitle,
  canFilter,
  isFiltered,
  filterWidth,
  onToggleSort,
  renderFilter,
  py,
}: THeadProps<InitialRefType>) => {
  const [isOpen, setIsOpen] = useTHeadState(id);
  const initialFocusRef = useRef<InitialRefType>(null);

  return (
    <Popover open={isOpen} onOpenChange={(value) => setIsOpen(value)}>
      <div className='flex w-full ml-[-22px] flex-col justify-start items-start'>
        <div
          className={cn(
            isFiltered || isOpen
              ? 'border-gray-300 shadow-sm'
              : 'border-transparent',
            (canSort && isOpen) || isSorted ? 'ml-0' : 'ml-3',
            !canSort ? ' ' : 'hover:ml-0',
            'flex items-center border rounded-[4px] transition-opacity duration-200 ease-in-out group',
          )}
          style={{ paddingTop: py ?? '0', paddingBottom: py ?? '0' }}
        >
          {canSort ? (
            isSorted === 'asc' ? (
              <ArrowUp
                id='sort-icon'
                onClick={onToggleSort}
                className={cn(
                  isSorted || isOpen ? 'w-3 inline-block' : 'w-0 ',
                  !isSorted ? 'text-gray-400' : 'text-gray-700',
                  'mx-1 w-3 h-3 cursor-pointer group-hover:transition-opacity group-hover:opacity-100 group-hover:w-3 group-hover:duration-200 group-hover:ease-in-out',
                )}
              />
            ) : (
              <ArrowDown
                id='sort-icon'
                onClick={onToggleSort}
                className={cn(
                  isSorted || isOpen ? 'w-3 opacity-100' : 'w-0 opacity-0',
                  !isSorted ? 'text-gray-400' : 'text-gray-700',
                  'mx-1 h-3 cursor-pointer group-hover:transition-opacity group-hover:opacity-100 group-hover:w-3 group-hover:duration-200 group-hover:ease-in-out',
                )}
              />
            )
          ) : (
            <div className={cn(canSort ? 'w-3' : 'w-0', 'flex mx-1')} />
          )}
          <p
            className={cn(
              isSorted ? 'mt-[-2px] tracking-[-0.3px] ' : 'mt-0',
              canSort ? 'cursor-pointer' : 'cursor-default',
              !isSorted ? 'font-base' : 'font-medium',
              'text-sm text-gray-700',
            )}
            onClick={onToggleSort}
          >
            {title}
          </p>
          {canFilter && (
            <>
              <PopoverTrigger>
                <IconButton
                  size='sm'
                  variant='ghost'
                  aria-label='filter'
                  className={cn(
                    isFiltered || isOpen ? 'opacity-100' : 'opacity-0',
                    'filter-icon-button mr-1 rounded-sm group-hover:transition-opacity group-hover:opacity-100 group-hover:duration-200 group-hover:ease-in-out',
                  )}
                  icon={
                    <FilterLines
                      className={cn(
                        isFiltered || isOpen
                          ? 'text-gray-700'
                          : 'text-gray-400',
                        'size-3',
                      )}
                    />
                  }
                />
              </PopoverTrigger>
              <PopoverContent
                onFocus={() => setIsOpen(true)}
                side='bottom'
                align='start'
                style={{ width: filterWidth ?? '12rem' }}
              >
                {renderFilter?.(initialFocusRef)}
              </PopoverContent>
            </>
          )}
        </div>
        {subTitle && <p className='text-xs text-gray-500'>{subTitle}</p>}
      </div>
    </Popover>
  );
};

export function getTHeadProps<T extends object>(
  context: HeaderContext<T, unknown>,
) {
  const header = context.header;

  const canSort = header.column.getCanSort();
  const isSorted = header.column.getIsSorted();
  const canFilter = header.column.getCanFilter();
  const isFiltered = header.column.getIsFiltered();
  const onToggleSort = header.column.getToggleSortingHandler();

  return {
    canSort,
    isSorted,
    canFilter,
    isFiltered,
    onToggleSort,
  };
}

export default memo(THead) as typeof THead;
