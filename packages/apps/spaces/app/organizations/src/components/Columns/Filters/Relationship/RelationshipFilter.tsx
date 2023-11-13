'use client';
import { useEffect } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { Organization } from '@graphql/types';
import { Checkbox, CheckboxGroup } from '@ui/form/Checkbox';

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

  return (
    <CheckboxGroup size='md' value={filter.value?.map((v) => String(v))}>
      <FilterHeader
        isChecked={toggle.isActive}
        onToggle={toggle.handleChange}
        onDisplayChange={toggle.handleClick}
      />
      <VStack spacing={2} align='flex-start'>
        <Checkbox value='true' onChange={handleSelect(true)}>
          <Text fontSize='sm'>Customer</Text>
        </Checkbox>
        <Checkbox value='false' onChange={handleSelect(false)}>
          <Text fontSize='sm'>Prospect</Text>
        </Checkbox>
      </VStack>
    </CheckboxGroup>
  );
};
