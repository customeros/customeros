'use client';
import type { HeaderContext } from '@tanstack/react-table';

import { memo, useRef, RefObject } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { ArrowUp } from '@ui/media/icons/ArrowUp';
import { ArrowDown } from '@ui/media/icons/ArrowDown';
import { FilterLines } from '@ui/media/icons/FilterLines';
import {
  Popover,
  PopoverBody,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover';

import { useTHeadState } from './THead.atom';

interface THeadProps<
  InitialRefType extends { focus(): void } = HTMLButtonElement,
> {
  id: string;
  title: string;
  subTitle?: string;
  canSort?: boolean;
  canFilter?: boolean;
  isFiltered?: boolean;
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
}: THeadProps<InitialRefType>) => {
  const [isOpen, setIsOpen] = useTHeadState(id);
  const initialFocusRef = useRef<InitialRefType>(null);

  const sortIconProps = {
    mx: '1',
    boxSize: '3',
    id: 'sort-icon',
    w: isSorted ? '3' : '0',
    opacity: isSorted ? 1 : 0,
    color: !isSorted ? 'gray.400' : 'gray.700',
  };

  return (
    <Popover
      isOpen={isOpen}
      placement='bottom-start'
      onOpen={() => setIsOpen(true)}
      onClose={() => setIsOpen(false)}
      initialFocusRef={initialFocusRef}
    >
      {({ isOpen }) => (
        <Flex
          w='full'
          ml='-22px'
          flexDir='column'
          justify='flex-start'
          alignItems='flex-start'
        >
          <Flex
            align='center'
            border='1px solid'
            borderRadius='4px'
            transition='opacity 0.2s ease-in-out'
            ml={isOpen || isSorted || !canSort ? '0' : '3'}
            borderColor={isFiltered || isOpen ? 'gray.300' : 'transparent'}
            boxShadow={
              isFiltered || isOpen
                ? '0px 1px 2px 0px rgba(16, 24, 40, 0.05)'
                : 'unset'
            }
            _hover={{
              ml: 0,
              '& #sort-icon': {
                transition: 'opacity 0.2s ease-in-out',
                opacity: 1,
                w: '3',
              },
              '& .filter-icon-button': {
                transition: 'opacity 0.2s ease-in-out',
                opacity: 1,
              },
            }}
          >
            {canSort ? (
              isSorted === 'asc' ? (
                <ArrowUp
                  {...sortIconProps}
                  w={isSorted || isOpen ? '3' : '0'}
                  opacity={isSorted || isOpen ? 1 : 0}
                />
              ) : (
                <ArrowDown
                  {...sortIconProps}
                  w={isSorted || isOpen ? '3' : '0'}
                  opacity={isSorted || isOpen ? 1 : 0}
                />
              )
            ) : (
              <Flex w='3' mx='1' />
            )}
            <Text
              fontSize='sm'
              color='gray.700'
              cursor='pointer'
              onClick={onToggleSort}
              mt={isSorted ? '-2px' : '0px'}
              fontWeight={!isSorted ? 'normal' : 'medium'}
              letterSpacing={isSorted ? '-0.3px' : 'unset'}
            >
              {title}
            </Text>
            {canFilter && (
              <>
                <PopoverTrigger>
                  <IconButton
                    mx='1'
                    size='14px'
                    variant='ghost'
                    borderRadius='2px'
                    aria-label='filter'
                    opacity={isFiltered || isOpen ? 1 : 0}
                    className='filter-icon-button'
                    icon={
                      <FilterLines
                        boxSize='3'
                        color={isFiltered || isOpen ? 'gray.700' : 'gray.400'}
                      />
                    }
                  />
                </PopoverTrigger>
                <PopoverContent maxW={filterWidth ?? '12rem'}>
                  <PopoverBody>{renderFilter?.(initialFocusRef)}</PopoverBody>
                </PopoverContent>
              </>
            )}
          </Flex>
          {subTitle && (
            <Text fontSize='xs' color='gray.500'>
              {subTitle}
            </Text>
          )}
        </Flex>
      )}
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
