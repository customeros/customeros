import { useEffect } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { Organization, OrganizationRelationship } from '@graphql/types';
import { relationshipOptions } from '@organizations/components/Columns/Cells/relationship/util.ts';

import { FilterHeader, useFilterToggle } from '../shared/FilterHeader';
import {
  useRelationshipFilter,
  RelationshipFilterSelector,
} from './RelationshipFilter.atom';

interface RelationshipFilterProps {
  onFilterValueChange?: Column<Organization>['setFilterValue'];
}

export const RelationshipFilter = ({
  onFilterValueChange,
}: RelationshipFilterProps) => {
  const [filter, setFilter] = useRelationshipFilter();
  const filterValue = useRecoilValue(RelationshipFilterSelector);

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

  const handleSelect = (value: OrganizationRelationship) => () => {
    setFilter((prev) => {
      const next = produce(prev, (draft) => {
        draft.isActive = true;

        if (draft.value.includes(value)) {
          draft.value = draft.value.filter((item) => item !== value);
        } else {
          draft.value.push(value);
        }
      });

      toggle.setIsActive(next.isActive);

      return next;
    });
  };

  useEffect(() => {
    onFilterValueChange?.(filterValue.isActive ? filterValue.value : undefined);
  }, [filterValue.value.length, filterValue.isActive]);

  return (
    <>
      <FilterHeader
        isChecked={toggle.isActive}
        onToggle={toggle.handleChange}
        onDisplayChange={toggle.handleClick}
      />
      <div className='flex flex-col gap-2 items-start'>
        {relationshipOptions.map((option) => (
          <Checkbox
            key={option.value.toString()}
            isChecked={filterValue.value.includes(option.value)}
            onChange={handleSelect(option.value)}
          >
            <p className='text-sm'>{option.label}</p>
          </Checkbox>
        ))}
      </div>
    </>
  );
};
