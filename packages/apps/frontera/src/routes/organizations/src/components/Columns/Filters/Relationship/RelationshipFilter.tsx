import { useEffect } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { Organization } from '@graphql/types';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';

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

  const handleSelect = (value: boolean) => () => {
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

  //need to be checked if is ok
  return (
    <>
      <FilterHeader
        isChecked={toggle.isActive}
        onToggle={toggle.handleChange}
        onDisplayChange={toggle.handleClick}
      />
      <div className='flex flex-col gap-2 items-start'>
        <Checkbox isChecked={true} onChange={handleSelect(true)}>
          <p className='text-sm'>Customer</p>
        </Checkbox>
        <Checkbox isChecked={false} onChange={handleSelect(false)}>
          <p className='text-sm'>Prospect</p>
        </Checkbox>
      </div>
    </>
  );
};
