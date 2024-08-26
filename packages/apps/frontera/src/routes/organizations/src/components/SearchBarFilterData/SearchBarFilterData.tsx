import { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

import { isEqual } from 'lodash';
import { observer } from 'mobx-react-lite';
import { FilterItem } from '@store/types.ts';

import { Switch } from '@ui/form/Switch';
import { useStore } from '@shared/hooks/useStore';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { useTablePlaceholder } from '@organizations/hooks/useTablePlaceholder.tsx';
import { filterOptions } from '@organizations/components/SearchBarFilterData/utils.ts';
import {
  Menu,
  MenuItem,
  MenuList,
  MenuButton,
} from '@ui/overlay/Menu/Menu.tsx';
import { useTableColumnOptionsMap } from '@organizations/hooks/useTableColumnOptionsMap.tsx';

export const SearchBarFilterData = observer(() => {
  const store = useStore();
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');
  const tableViewName = store.tableViewDefs.getById(preset || '')?.value.name;
  const [filters, setFilters] = useState<
    Map<string, FilterItem & { active: boolean }>
  >(new Map());

  const [isOpen, setIsOpen] = useState(false);
  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');
  const { multi: multiResultPlaceholder, single: singleResultPlaceholder } =
    useTablePlaceholder(tableViewName);

  const [optionsMap] = useTableColumnOptionsMap(tableViewDef?.value?.tableType);
  const filterOptionMap = { ...filterOptions, ...optionsMap };
  const appliedFilters = tableViewDef
    ?.getFilters()
    ?.AND?.filter(
      ({ filter }: { filter: FilterItem & { active: boolean } }) =>
        filter.active,
    );

  useEffect(() => {
    if (!isOpen) {
      const newFiltersMap: Map<string, FilterItem & { active: boolean }> =
        new Map(
          appliedFilters?.map(
            (item: { filter: FilterItem & { active: boolean } }) => [
              item.filter.property,
              item.filter,
            ],
          ) || [],
        );

      if (
        !isEqual(
          Array.from(filters.entries()),
          Array.from(newFiltersMap.entries()),
        )
      ) {
        setFilters(newFiltersMap);
      }
    }
  }, [appliedFilters, isOpen]);

  const isEmailFilterActive = () => {
    return (
      (filters.get('EMAIL_VERIFICATION')?.active ||
        filters.get('CONTACTS_EMAILS')?.active) ??
      false
    );
  };
  const totalResults = store.ui.searchCount;

  const tableName =
    totalResults === 1 ? singleResultPlaceholder : multiResultPlaceholder;

  const handleApplyChanges = () => {
    filters.forEach((filter, property) => {
      if (!filter.active) {
        tableViewDef?.removeFilter(property);
      }
    });
  };

  const handleChange = (property: string, active: boolean) => {
    setFilters((prev) => {
      const newFilters = new Map(prev);

      if (newFilters.has(property)) {
        newFilters.set(property, { ...newFilters.get(property)!, active });
      }

      return newFilters;
    });
  };

  const renderEmailFilter = () => {
    if (filters.has('EMAIL_VERIFICATION') || filters.has('CONTACTS_EMAILS')) {
      return (
        <MenuItem
          key='EmailFilter'
          className='flex justify-between font-normal capitalize mb-1'
          onClick={(e) => {
            e.stopPropagation();
            e.preventDefault();
          }}
        >
          Email
          <div className='ml-2 flex items-center'>
            <Switch
              size='sm'
              isChecked={isEmailFilterActive()}
              onChange={(newState) => {
                handleChange('EMAIL_VERIFICATION', newState);
                handleChange('CONTACTS_EMAILS', newState);
              }}
            />
          </div>
        </MenuItem>
      );
    }

    return null;
  };

  const renderFilterItem = (property: string) => {
    if (property === 'EMAIL_VERIFICATION' || property === 'CONTACTS_EMAILS') {
      return null; // These will be handled by renderEmailFilter
    }

    return (
      <MenuItem
        key={property}
        className='flex justify-between font-normal capitalize mb-1'
        onClick={(e) => {
          e.stopPropagation();
          e.preventDefault();
        }}
      >
        {filterOptionMap[property]}
        <div className='ml-2 flex items-center'>
          <Switch
            size='sm'
            isChecked={filters.get(property)?.active ?? false}
            onChange={(newState) => {
              handleChange(property, newState);
            }}
          />
        </div>
      </MenuItem>
    );
  };

  return (
    <div className='flex flex-row items-center gap-1'>
      <SearchSm className='size-5' />
      <div
        data-test={`search-${tableName}`}
        className={'font-medium flex items-center gap-1 break-keep w-max '}
      >
        {totalResults}{' '}
        {appliedFilters?.length ? (
          <Menu
            onOpenChange={(open) => {
              setIsOpen(open);

              if (!open) {
                handleApplyChanges();
              }
            }}
          >
            <MenuButton className='min-h-[38px] outline-none focus:outline-none underline text-gray-500'>
              filtered
            </MenuButton>
            <MenuList side='bottom' align='start' className='min-w-12'>
              <p className='font-medium mx-2 mb-2 min-w-[210px]'>
                <span className='capitalize mr-1'>{tableName}</span>
                filtered by:{' '}
              </p>
              {renderEmailFilter()}
              {appliedFilters?.map(
                ({ filter: { property } }: { filter: FilterItem }) =>
                  renderFilterItem(property),
              )}
            </MenuList>
          </Menu>
        ) : (
          ''
        )}{' '}
        {tableName}:
      </div>
    </div>
  );
});
