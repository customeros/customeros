import { useState, useEffect, RefObject, useCallback } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { Organization } from '@graphql/types';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';

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
  const [displayValue, setDisplayValue] = useState<string>(() => filter.value);
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
          const nextValue = value.trim();

          draft.value = nextValue;

          if (!draft.showEmpty) {
            draft.isActive = !!nextValue;
          }
        });

        return next;
      });
    },
    [setFilter],
  );

  const handleDisplayChange = useCallback(
    (value: string) => {
      setDisplayValue(value.trim());
      !filter.showEmpty && toggle.setIsActive(!!value.trim());
    },
    [setDisplayValue, toggle.setIsActive, filter.showEmpty],
  );

  const handleShowEmpty = useCallback(
    (isChecked: boolean) => {
      setFilter((prev) => {
        const next = produce(prev, (draft) => {
          draft.showEmpty = isChecked;
        });

        if (!next.value.trim()) {
          toggle.setIsActive(isChecked);
        }

        return next;
      });
    },
    [setFilter, setDisplayValue, toggle.setIsActive],
  );

  useEffect(() => {
    onFilterValueChange?.(filterValue.isActive ? filterValue : undefined);
  }, [filterValue.value, filterValue.isActive, filterValue.showEmpty]);

  return (
    <>
      <FilterHeader
        isChecked={toggle.isActive}
        onToggle={toggle.handleChange}
        onDisplayChange={toggle.handleClick}
      />

      <DebouncedSearchInput
        value={displayValue}
        ref={initialFocusRef}
        onChange={handleChange}
        onDisplayChange={handleDisplayChange}
      />

      <Checkbox
        size='md'
        className='mt-2'
        isChecked={filter.showEmpty}
        onChange={(isChecked) => handleShowEmpty(isChecked as boolean)}
      >
        <p className='text-sm'>Unnamed</p>
      </Checkbox>
    </>
  );
};
