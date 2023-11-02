import type { HeaderContext } from '@tanstack/react-table';
import { Column } from '@tanstack/react-table';

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

interface THeadProps<T extends object> extends HeaderContext<T, unknown> {
  title: string;
  subTitle?: string;
  icon?: React.ReactNode;
  renderFilter?: (column: Column<T>) => React.ReactNode;
}

export const THead = <T extends object>({
  icon,
  title,
  header,
  subTitle,
  renderFilter,
}: THeadProps<T>) => {
  const canSort = header.column.getCanSort();
  const isSorted = header.column.getIsSorted();
  const canFilter = header.column.getCanFilter();
  const onToggleSort = header.column.getToggleSortingHandler();

  const sortIconProps = {
    mx: '1',
    boxSize: '3',
    id: 'sort-icon',
    opacity: isSorted ? 1 : 0,
    color: !isSorted ? 'gray.400' : 'gray.700',
  };

  return (
    <Popover placement='bottom-start'>
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
            transition='all 0.2s ease-in-out'
            borderColor={isOpen ? 'gray.300' : 'transparent'}
            boxShadow={
              isOpen ? '0px 1px 2px 0px rgba(16, 24, 40, 0.05)' : 'unset'
            }
            _hover={{
              '& #sort-icon': {
                transition: 'opacity 0.2s ease-in-out',
                opacity: 1,
              },
              '& .filter-icon-button': {
                transition: 'opacity 0.2s ease-in-out',
                opacity: 1,
              },
            }}
          >
            {canSort ? (
              isSorted === 'asc' ? (
                <ArrowUp {...sortIconProps} />
              ) : (
                <ArrowDown {...sortIconProps} />
              )
            ) : (
              <Flex w='3' mx='1' />
            )}
            <Text
              fontSize='sm'
              color='gray.700'
              cursor='pointer'
              onClick={onToggleSort}
              fontWeight={!isSorted ? 'normal' : 'medium'}
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
                    opacity={isOpen ? 1 : 0}
                    className='filter-icon-button'
                    icon={
                      <FilterLines
                        boxSize='3'
                        color={isOpen ? 'gray.700' : 'gray.400'}
                      />
                    }
                  />
                </PopoverTrigger>
                <PopoverContent>
                  <PopoverBody>{renderFilter?.(header.column)}</PopoverBody>
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
