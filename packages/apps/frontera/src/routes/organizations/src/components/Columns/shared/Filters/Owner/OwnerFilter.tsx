import { useState, RefObject } from 'react';
import { useSearchParams } from 'react-router-dom';

import uniqBy from 'lodash/uniqBy';
import difference from 'lodash/difference';
import { observer } from 'mobx-react-lite';
import { FilterItem } from '@store/types.ts';
import intersection from 'lodash/intersection';

import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox.tsx';
import { Tumbleweed } from '@ui/media/icons/Tumbleweed.tsx';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';

import { FilterHeader, DebouncedSearchInput } from '../abstract';

interface OwnerFilterProps {
  property?: ColumnViewType;
  initialFocusRef: RefObject<HTMLInputElement>;
}

const defaultFilter: FilterItem = {
  property: ColumnViewType.OrganizationsOwner,
  value: [],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.In,
};

export const OwnerFilter = observer(
  ({ initialFocusRef, property }: OwnerFilterProps) => {
    const [searchParams] = useSearchParams();
    const preset = searchParams.get('preset');

    const store = useStore();
    const tableViewDef = store.tableViewDefs.getById(preset ?? '');
    const filter = tableViewDef?.getFilter(
      property || defaultFilter.property,
    ) ?? { ...defaultFilter, property: property || defaultFilter.property };

    const toggle = () => {
      tableViewDef?.toggleFilter(filter);
    };

    const [searchValue, setSearchValue] = useState('');

    const users = store.users.toComputedArray((arr) => {
      if (searchValue) {
        return arr.filter((user) =>
          user.name.toLowerCase().includes(searchValue.toLowerCase()),
        );
      }

      return arr;
    });

    const options = uniqBy(
      [
        { value: '__EMPTY__', label: 'Unknown' },
        ...(users
          .map((u) => ({
            value: u.id,
            label: u.name,
          }))
          .filter((o) => o.label) ?? []),
      ],
      'label',
    ).filter((v) => {
      return searchValue ? v.value !== '__EMPTY__' : true;
    });

    const userIds = options.map(({ value }) => value);
    const isAllSelected =
      intersection(filter.value, userIds).length === userIds.length &&
      userIds.length > 0;

    const handleSelectAll = () => {
      let nextValue: string[] = [];

      if (isAllSelected) {
        tableViewDef?.setFilter({
          ...filter,
          value: difference(filter.value, userIds),
          active: false,
        });

        return;
      }

      if (searchValue) {
        nextValue = [...userIds, ...difference(filter.value, userIds)];
      } else {
        nextValue = userIds;
      }

      tableViewDef?.setFilter({
        ...filter,
        value: nextValue,
        active: nextValue.length > 0,
      });
    };

    const handleSelect = (value: string) => () => {
      const nextValue = filter.value.includes(value)
        ? filter.value.filter((item: string) => item !== value)
        : [...filter.value, value];

      tableViewDef?.setFilter({
        ...filter,
        value: nextValue,
        active: nextValue.length > 0,
      });
    };

    return (
      <>
        <FilterHeader
          onToggle={toggle}
          onDisplayChange={() => {}}
          isChecked={filter.active ?? false}
        />

        <DebouncedSearchInput
          value={searchValue}
          ref={initialFocusRef}
          onDisplayChange={(v) => setSearchValue(v)}
          onChange={(v) => {
            setSearchValue(v);

            if ((v.length && !filter.active) || (!v.length && filter.active)) {
              toggle();
            }
          }}
        />

        <div className='flex flex-col w-full h-[13rem] items-start gap-2 mt-2 px-1 mx-[-4px] overflow-x-hidden overflow-y-auto'>
          {users.length > 1 && (
            <div className='sticky top-0 w-full z-10 bg-white gap-2 flex flex-col pb-2 border-b border-gray-200'>
              <Checkbox
                className='top-0 z-10'
                isChecked={isAllSelected}
                onChange={handleSelectAll}
              >
                <span className='text-sm'>
                  {isAllSelected
                    ? 'Deselect all'
                    : 'Select all' +
                      (searchValue && users.length > 2
                        ? ` ${users.length}`
                        : '')}
                </span>
              </Checkbox>
            </div>
          )}

          {options.length > 0 ? (
            options.map(({ value, label }) => (
              <Checkbox
                key={value}
                onChange={handleSelect(value)}
                isChecked={filter.value.includes(value)}
              >
                <span className='text-sm line-clamp-1'>{label}</span>
              </Checkbox>
            ))
          ) : (
            <div className='flex w-full justify-center items-center flex-col'>
              <Tumbleweed className='size-8 text-gray-400' />
              <span className='text-center text-sm text-gray-500'>
                Empty here in <b>No Resultsville</b>
              </span>
            </div>
          )}
        </div>
      </>
    );
  },
);
