import { useMemo, useState, useEffect, RefObject } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import difference from 'lodash/difference';
import intersection from 'lodash/intersection';
import { Column } from '@tanstack/react-table';

import { Organization } from '@graphql/types';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { Tumbleweed } from '@ui/media/icons/Tumbleweed';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGetUsersQuery } from '@shared/graphql/getUsers.generated';

import { useOwnerFilter, OwnerFilterSelector } from './OwnerFilter.atom';
import { FilterHeader, useFilterToggle, DebouncedSearchInput } from '../shared';

interface OwnerFilterProps {
  initialFocusRef: RefObject<HTMLInputElement>;
  onFilterValueChange?: Column<Organization>['setFilterValue'];
}

export const OwnerFilter = ({
  initialFocusRef,
  onFilterValueChange,
}: OwnerFilterProps) => {
  const client = getGraphQLClient();
  const [filter, setFilter] = useOwnerFilter();
  const [searchValue, setSearchValue] = useState('');
  const filterValue = useRecoilValue(OwnerFilterSelector);

  const toggle = useFilterToggle({
    defaultValue: filter.isActive,
    onToggle: (setIsActive) => {
      setFilter((prev) => {
        const next = produce(prev, (draft) => {
          draft.isActive = !draft.isActive;
        });

        setIsActive(next.isActive);

        return next;
      });
    },
  });

  const { data } = useGetUsersQuery(client, {
    pagination: { limit: 100, page: 1 },
  });

  const users = useMemo(() => {
    const items = [
      { value: '__EMPTY__', label: 'Unknown' },
      ...(data?.users.content.map(({ id, name, firstName, lastName }) => ({
        value: id,
        label: name ? name : [firstName, lastName].filter(Boolean).join(' '),
      })) ?? []),
    ];

    if (!searchValue) return items;

    return items.filter(({ label }) =>
      label?.toLowerCase().includes(searchValue.toLowerCase()),
    );
  }, [data?.users.content, searchValue]);

  const userIds = users.map(({ value }) => value);
  const isAllSelected =
    intersection(filter.value, userIds).length === users.length &&
    users.length > 0;

  const handleSelectAll = () => {
    setFilter((prev) => {
      const next = produce(prev, (draft) => {
        draft.isActive = true;

        if (isAllSelected) {
          draft.value = draft.value.filter((item) => !userIds.includes(item));

          if (draft.value.length === 0) {
            draft.isActive = false;
          }

          return;
        }

        if (searchValue) {
          draft.value = [...userIds, ...difference(draft.value, userIds)];

          return;
        }

        draft.value = userIds;
      });

      toggle.setIsActive(next.isActive);

      return next;
    });
  };

  const handleSelect = (value: string) => () => {
    setFilter((prev) => {
      const next = produce(prev, (draft) => {
        draft.isActive = true;

        if (draft.value.includes(value)) {
          draft.value = draft.value.filter((item) => item !== value);

          if (draft.value.length === 0) {
            draft.isActive = false;
          }
        } else {
          draft.value.push(value);
        }
      });

      toggle.setIsActive(next.isActive);

      return next;
    });
  };

  useEffect(() => {
    onFilterValueChange?.(filterValue.isActive ? filterValue : undefined);
  }, [filterValue.value.length, filterValue.isActive, filterValue.showEmpty]);

  return (
    <>
      <FilterHeader
        isChecked={toggle.isActive}
        onToggle={toggle.handleChange}
        onDisplayChange={toggle.handleClick}
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

        {users.length > 0 ? (
          users.map(({ value, label }) => (
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
            <Tumbleweed className='mr-10 size-8 text-gray-400 self-end' />
            <span className='text-sm text-gray-500'>
              Empty here in <b>No Resultsville</b>
            </span>
          </div>
        )}
      </div>
    </>
  );
};
