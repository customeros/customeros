'use client';
import { useEffect, RefObject, useCallback } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { Organization } from '@graphql/types';

import { FilterHeader, useFilterToggle, DebouncedSearchInput } from '../shared';
import {
  useOrganizationFilter,
  OrganizationFilterSelector,
} from './OrganizationFilter.atom';

interface OrganizationFilterProps {
  initialFocusRef: RefObject<HTMLInputElement>;
  onFilterValueChange?: Column<Organization>['setFilterValue'];
}

export const OrganizationFilter = ({
  initialFocusRef,
  onFilterValueChange,
}: OrganizationFilterProps) => {
  const [filter, setFilter] = useOrganizationFilter();
  const filterValue = useRecoilValue(OrganizationFilterSelector);

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
    onFilterValueChange?.(filterValue.isActive ? filterValue.value : undefined);
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
