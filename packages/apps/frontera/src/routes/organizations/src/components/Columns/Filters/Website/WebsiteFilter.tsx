import { useState, useEffect, RefObject, useCallback } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { Organization } from '@graphql/types';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';

import { useWebsiteFilter, WebsiteFilterSelector } from './WebsiteFilter.atom';
import { FilterHeader, useFilterToggle, DebouncedSearchInput } from '../shared';

interface WebsiteFilterProps {
  initialFocusRef: RefObject<HTMLInputElement>;
  onFilterValueChange?: Column<Organization>['setFilterValue'];
}

export const WebsiteFilter = ({
  initialFocusRef,
  onFilterValueChange,
}: WebsiteFilterProps) => {
  const [filter, setFilter] = useWebsiteFilter();
  const [displayValue, setDisplayValue] = useState(() => filter.value);
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

  const handleDisplayChange = useCallback(
    (value: string) => {
      setDisplayValue(value.trim());
      !filter.showEmpty && toggle.setIsActive(!!value.trim());
    },
    [setDisplayValue, toggle.setIsActive, filter.showEmpty],
  );

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

  const handleShowEmpty = useCallback(
    (isChecked: boolean) => {
      setFilter((prev) => {
        const next = produce(prev, (draft) => {
          draft.showEmpty = isChecked;
        });

        if (!next.value) {
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
        className='mt-2'
        size='md'
        onChange={(isChecked) => handleShowEmpty(isChecked as boolean)}
        isChecked={filter.showEmpty}
        labelProps={{ className: 'text-sm mt-2' }}
      >
        Unknown
      </Checkbox>
    </>
  );
};
