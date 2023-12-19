'use client';
import { useEffect } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { SelectOption } from '@ui/utils/types';
import { OnboardingStatus } from '@graphql/types';
import { Checkbox, CheckboxGroup } from '@ui/form/Checkbox';

import { FilterHeader, useFilterToggle } from '../shared/FilterHeader';
import {
  useOnboardingFilter,
  OnboardingFilterSelector,
} from './OnboardingFilter.atom';

type ComputedOnboardingStatus = OnboardingStatus | 'NOT_ONBOARDING';

const options: SelectOption<ComputedOnboardingStatus>[] = [
  { label: 'Not started', value: OnboardingStatus.NotStarted },
  { label: 'Late', value: OnboardingStatus.Late },
  { label: 'Stuck', value: OnboardingStatus.Stuck },
  { label: 'On track', value: OnboardingStatus.OnTrack },
  { label: 'Done', value: OnboardingStatus.Done },
  { label: 'Not onboarding', value: 'NOT_ONBOARDING' },
];
const notOnboardingOptions = [
  OnboardingStatus.NotApplicable,
  OnboardingStatus.Successful,
];

interface RelationshipFilterProps<T> {
  column: Column<T>;
}

export const OnboardingFilter = <T,>({
  column,
}: RelationshipFilterProps<T>) => {
  const [filter, setFilter] = useOnboardingFilter();
  const filterValue = useRecoilValue(OnboardingFilterSelector);

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

  const handleSelect = (value: ComputedOnboardingStatus) => () => {
    setFilter((prev) => {
      const next = produce(prev, (draft) => {
        const isComputed = value === 'NOT_ONBOARDING';
        draft.isActive = true;

        if (
          isComputed
            ? draft.value.some((v) => notOnboardingOptions.includes(v))
            : draft.value.includes(value)
        ) {
          draft.value = draft.value.filter((item) =>
            isComputed ? !notOnboardingOptions.includes(item) : item !== value,
          );
        } else {
          if (isComputed) {
            draft.value.push(...notOnboardingOptions);
          } else {
            draft.value.push(value);
          }
        }
      });

      toggle.setIsActive(next.isActive);

      return next;
    });
  };

  const handleSelectAll = () => {
    setFilter((prev) => {
      const next = produce(prev, (draft) => {
        draft.isActive = true;

        if (draft.value.length === options.length + 1) {
          draft.value = [];
        } else {
          draft.value = options
            .map((option) => option.value)
            .filter((v) => v !== 'NOT_ONBOARDING')
            .concat(...notOnboardingOptions) as OnboardingStatus[];
        }
      });

      toggle.setIsActive(next.isActive);

      return next;
    });
  };

  useEffect(() => {
    column.setFilterValue?.(
      filterValue.isActive ? filterValue.value : undefined,
    );
  }, [filterValue.value.length, filterValue.isActive]);

  const computedValues = computeFilterValues(filterValue.value);
  const isAllChecked = filterValue.value.length === options.length + 1;

  return (
    <CheckboxGroup size='md' value={computedValues}>
      <FilterHeader
        isChecked={toggle.isActive}
        onToggle={toggle.handleChange}
        onDisplayChange={toggle.handleClick}
      />
      <VStack spacing={2} align='flex-start'>
        <Checkbox isChecked={isAllChecked} onChange={handleSelectAll}>
          <Text fontSize='sm'>
            {isAllChecked ? 'Deselect all' : 'Select all'}
          </Text>
        </Checkbox>
        {options.map((option) => (
          <Checkbox
            key={option.label}
            value={option.value}
            onChange={handleSelect(option.value)}
          >
            <Text fontSize='sm'>{option.label}</Text>
          </Checkbox>
        ))}
      </VStack>
    </CheckboxGroup>
  );
};

function computeFilterValues(values: OnboardingStatus[]) {
  const outputSet = new Set(
    values.map((v) =>
      notOnboardingOptions.includes(v) ? 'NOT_ONBOARDING' : v,
    ),
  );

  return Array.from(outputSet);
}
