import { useState, RefObject } from 'react';
import { useSearchParams } from 'react-router-dom';

import { FilterItem } from '@store/types';
import difference from 'lodash/difference';
import { observer } from 'mobx-react-lite';
import intersection from 'lodash/intersection';

import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { Tumbleweed } from '@ui/media/icons/Tumbleweed';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';

import { FilterHeader, DebouncedSearchInput } from '../shared';

interface OwnerFilterProps {
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

export const OwnerFilter = observer(({ initialFocusRef }: OwnerFilterProps) => {
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');

  const store = useStore();
  const tableViewDef = store.tableViewDefs.getById(preset ?? '');
  const filter =
    tableViewDef?.getFilter(defaultFilter.property) ?? defaultFilter;

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

  const options = [
    { value: '__EMPTY__', label: 'Unknown' },
    ...(users
      .map((u) => ({
        value: u.id,
        label: u.name,
      }))
      .filter((o) => o.label) ?? []),
  ].filter((v) => {
    return searchValue ? v.value !== '__EMPTY__' : true;
  });

  const userIds = options.map(({ value }) => value);
  const isAllSelected =
    intersection(filter.value, userIds).length === users.length + 1 &&
    users.length > 0;

  const handleSelectAll = () => {
    let nextValue: string[] = [];

    if (isAllSelected) {
      tableViewDef?.setFilter({
        ...filter,
        value: [],
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
        onChange={(v) => setSearchValue(v)}
        onDisplayChange={(v) => setSearchValue(v)}
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
                    (searchValue && users.length > 2 ? ` ${users.length}` : '')}
              </span>
            </Checkbox>
          </div>
        )}

        {options.length > 0 ? (
          options.map(({ value, label }) => (
            <Checkbox
              key={value}
              isChecked={filter.value.includes(value)}
              onChange={handleSelect(value)}
            >
              <span className='text-sm line-clamp-1'>{label}</span>
            </Checkbox>
          ))
        ) : (
          <div className='flex w-full justify-center items-center flex-col'>
            <Tumbleweed className='mr-10 size-8 text-gray-400 self-end' />
            <span className='text-sm text-gray-500'>
              Empty here in <b>No Resultsville</b>
            </span>
          </div>
        )}
      </div>
    </>
  );
});
