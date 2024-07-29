import { useState } from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { FilterItem } from '@store/types.ts';
import { useDeepCompareEffect } from 'rooks';

import { Switch } from '@ui/form/Switch';
import { useStore } from '@shared/hooks/useStore';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { useTablePlaceholder } from '@organizations/hooks/useTablePlaceholder.tsx';
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
  const [filters, setFilters] = useState<Array<{ filter: FilterItem }>>([]);
  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');
  const { multi: multiResultPlaceholder, single: singleResultPlaceholder } =
    useTablePlaceholder(tableViewName);

  const [optionsMap] = useTableColumnOptionsMap(tableViewDef?.value?.tableType);

  const appliedFilters = tableViewDef
    ?.getFilters()
    ?.AND?.filter(
      ({ filter }: { filter: FilterItem & { active: boolean } }) =>
        filter.active,
    );

  useDeepCompareEffect(() => {
    setFilters(appliedFilters);
  }, [appliedFilters]);

  const totalResults = store.ui.searchCount;

  const tableName =
    totalResults === 1 ? singleResultPlaceholder : multiResultPlaceholder;

  const handleApplyChanges = () => {
    filters.forEach(({ filter }) => {
      if (filter.active) return;

      tableViewDef?.removeFilter(filter.property);
    });
  };
  const handleChange = (property: string, active: boolean) => {
    const filter = filters.find((item) => item?.filter.property === property);
    if (filter) {
      setFilters((prev) => {
        return prev.map((item) => {
          if (item?.filter.property === property) {
            return {
              ...item,
              filter: {
                ...item.filter,
                active,
              },
            };
          }

          return item;
        });
      });
    }
  };

  return (
    <div className='flex flex-row items-center gap-1'>
      <SearchSm className='size-5' />
      <div
        className={
          'font-medium flex items-center gap-1 break-keep w-max mb-[2px]'
        }
        data-test={`search-${tableName}`}
      >
        {totalResults}{' '}
        {appliedFilters?.length ? (
          <Menu
            onOpenChange={(open) => {
              if (!open) {
                handleApplyChanges();
              }
            }}
          >
            <MenuButton className='min-h-[40px] outline-none focus:outline-none underline text-gray-500'>
              filtered
            </MenuButton>
            <MenuList side='bottom' align='start' className='min-w-12'>
              <p className='font-medium mx-2 mb-2 min-w-[210px]'>
                <span className='capitalize mr-1'>{tableName}</span>
                filtered by:{' '}
              </p>
              {appliedFilters?.map(({ filter }: { filter: FilterItem }) => (
                <MenuItem
                  key={filter.property}
                  className='flex justify-between font-normal capitalize mb-1 '
                  onClick={(e) => {
                    e.stopPropagation();
                    e.preventDefault();
                  }}
                >
                  {optionsMap[filter.property]}

                  <div className='ml-2 flex items-center'>
                    <Switch
                      size='sm'
                      isChecked={
                        filters.find(
                          (e) => e?.filter.property === filter.property,
                        )?.filter.active
                      }
                      onChange={(e) => {
                        handleChange(filter.property, e);
                      }}
                    />
                  </div>
                </MenuItem>
              ))}
            </MenuList>
          </Menu>
        ) : (
          ''
        )}{' '}
        {tableName}
      </div>
    </div>
  );
});
