'use client';
import { useEffect, RefObject, useCallback } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { Organization } from '@graphql/types';

import { useWebsiteFilter, WebsiteFilterSelector } from './WebsiteFilter.atom';
import { FilterHeader, useFilterToggle, DebouncedSearchInput } from '../shared';

interface WebsiteFilterProps {
  column: Column<Organization>;
  initialFocusRef: RefObject<HTMLInputElement>;
}

export const WebsiteFilter = ({
  column,
  initialFocusRef,
}: WebsiteFilterProps) => {
  const [filter, setFilter] = useWebsiteFilter();
  const filterValue = useRecoilValue(WebsiteFilterSelector);

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

  const handleChange = useCallback(
    (value: string) => {
      setFilter((prev) => {
        const next = produce(prev, (draft) => {
          draft.value = value;
          if (!value) {
            draft.isActive = false;
          } else {
            draft.isActive = true;
          }
        });

        toggle.setIsActive(next.isActive);

        return next;
      });
    },
    [setFilter, toggle.setIsActive],
  );

  useEffect(() => {
    column.setFilterValue(filterValue.isActive ? filterValue.value : undefined);
  }, [filterValue.value, filterValue.isActive]);

  return (
    <>
      <FilterHeader
        isChecked={toggle.isActive}
        onToggle={toggle.handleChange}
        onDisplayChange={toggle.handleClick}
      />

      <DebouncedSearchInput
        ref={initialFocusRef}
        value={filter.value}
        onChange={handleChange}
      />
    </>
  );
};
