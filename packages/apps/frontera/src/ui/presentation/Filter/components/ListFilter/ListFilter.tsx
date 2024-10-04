import { useRef, useState, useEffect } from 'react';

import { Input } from '@ui/form/Input';
import { flags } from '@ui/media/flags';
import { Avatar } from '@ui/media/Avatar';
import { Check } from '@ui/media/icons/Check';
import { Button } from '@ui/form/Button/Button';
import { User01 } from '@ui/media/icons/User01';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { ComparisonOperator } from '@shared/types/__generated__/graphql.types';

import { handleOperatorName, handlePropertyPlural } from '../../utils/utils';

interface ListFilterProps {
  filterName: string;
  operatorName: string;
  filterValue: string[];
  onMultiSelectChange: (ids: string[]) => void;
  options: { id: string; label: string; avatar?: string }[];
}

export const ListFilter = ({
  onMultiSelectChange,
  options,
  filterValue,
  filterName,
  operatorName,
}: ListFilterProps) => {
  const [selectedIds, setSelectedIds] = useState<string[]>(filterValue || []);
  const [search, setSearch] = useState('');
  const [isOpen, setIsOpen] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (!filterValue) {
      if (filterName) {
        setTimeout(() => {
          setIsOpen(true);
        }, 100);
      }
    }
  }, [filterName]);

  const handleItemClick = (id: string) => {
    let newSelectedIds;

    if (filterName === 'Ownership type') {
      newSelectedIds = [id];
    } else {
      newSelectedIds = selectedIds?.includes(id)
        ? selectedIds.filter((selectedId) => selectedId !== id)
        : [...selectedIds, id];
    }

    setSelectedIds(newSelectedIds);
    onMultiSelectChange(newSelectedIds);
  };

  const filterValueLabels = options
    .filter((option) => selectedIds?.includes(option.id))
    .map((option) => option.label);

  useEffect(() => {
    setTimeout(() => {
      if (isOpen && inputRef.current) {
        inputRef.current.focus();
      }
    }, 100);
  }, [isOpen]);

  if (
    operatorName === ComparisonOperator.IsEmpty ||
    operatorName === ComparisonOperator.IsNotEmpty
  )
    return null;

  return (
    <Menu open={isOpen} onOpenChange={(open) => setIsOpen(open)}>
      <MenuButton asChild>
        <Button
          size='xs'
          colorScheme='grayModern'
          className='border-l-0 rounded-none text-gray-700 bg-white font-normal'
        >
          {filterValueLabels.length === 1
            ? filterValueLabels?.[0]
            : filterValueLabels.length > 1
            ? `${filterValueLabels.length} ${handlePropertyPlural(
                filterName,
                selectedIds,
              )}`
            : '...'}
        </Button>
      </MenuButton>
      <MenuList
        align='start'
        side='bottom'
        className='max-h-[400px] overflow-auto'
      >
        <Input
          size='xs'
          ref={inputRef}
          variant='unstyled'
          className='px-2.5'
          onChange={(e) => setSearch(e.target.value)}
          placeholder={`${filterName} ${handleOperatorName(
            operatorName as ComparisonOperator,
            'list',
          )}`}
        />
        {options
          .filter((o) => o.label?.toLowerCase().includes(search?.toLowerCase()))
          .map((option) => {
            const flag = flags[option.id];

            return (
              <MenuItem
                key={option.id}
                onKeyDown={(e) => {
                  e.stopPropagation();
                }}
                onClick={(e) => {
                  e.preventDefault();
                  handleItemClick(option.id);
                }}
              >
                <div className='flex items-center justify-between w-full'>
                  <div className='flex items-center gap-2'>
                    {(option.avatar || filterName === 'Owner') && (
                      <Avatar
                        size='xxs'
                        textSize='xxs'
                        variant='outlineCircle'
                        src={option.avatar ?? ''}
                        name={option.label ?? 'Unnamed'}
                        icon={<User01 className='text-gray-500 size-3' />}
                      />
                    )}
                    {filterName === 'Country' && <span>{flag}</span>}
                    <span>{option.label}</span>
                  </div>

                  {filterValueLabels.includes(option.label) && (
                    <Check className='text-primary-600' />
                  )}
                </div>
              </MenuItem>
            );
          })}
      </MenuList>
    </Menu>
  );
};
