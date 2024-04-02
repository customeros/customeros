'use client';
import { useEffect } from 'react';

import { produce } from 'immer';
import { useRecoilValue } from 'recoil';
import { Column } from '@tanstack/react-table';

import { Checkbox } from '@ui/form/Checkbox/Checkbox2';
import { OpportunityRenewalLikelihood } from '@graphql/types';

import { FilterHeader, useFilterToggle } from '../shared';
import {
  useRenewalLikelihoodFilter,
  RenewalLikelihoodFilterSelector,
} from './RenewalLikelihoodFilter.atom';

interface RenewalLikelihoodFilterProps<T> {
  column: Column<T>;
}

export const RenewalLikelihoodFilter = <T,>({
  column,
}: RenewalLikelihoodFilterProps<T>) => {
  const [filter, setFilter] = useRenewalLikelihoodFilter();
  const filterValue = useRecoilValue(RenewalLikelihoodFilterSelector);

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

  const handleSelect = (value: OpportunityRenewalLikelihood) => () => {
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
    column.setFilterValue(filterValue.isActive ? filterValue.value : undefined);
  }, [filterValue.value.length, filterValue.isActive]);

  return (
    <>
      <FilterHeader
        isChecked={toggle.isActive}
        onToggle={toggle.handleChange}
        onDisplayChange={toggle.handleClick}
      />
      <div className='flex flex-col space-y-2 items-start'>
        <Checkbox
          isChecked={filter.value.includes(
            OpportunityRenewalLikelihood.HighRenewal,
          )}
          onChange={() =>
            handleSelect(OpportunityRenewalLikelihood.HighRenewal)
          }
        >
          <span className='text-sm'>High</span>
        </Checkbox>
        <Checkbox
          isChecked={filter.value.includes(
            OpportunityRenewalLikelihood.MediumRenewal,
          )}
          onChange={() =>
            handleSelect(OpportunityRenewalLikelihood.MediumRenewal)
          }
        >
          <span className='text-sm'>Medium</span>
        </Checkbox>
        <Checkbox
          isChecked={filter.value.includes(
            OpportunityRenewalLikelihood.LowRenewal,
          )}
          onChange={() => handleSelect(OpportunityRenewalLikelihood.LowRenewal)}
        >
          <span className='text-sm'>Low</span>
        </Checkbox>
        <Checkbox
          isChecked={filter.value.includes(
            OpportunityRenewalLikelihood.ZeroRenewal,
          )}
          onChange={() =>
            handleSelect(OpportunityRenewalLikelihood.ZeroRenewal)
          }
        >
          <span className='text-sm'>Zero</span>
        </Checkbox>
      </div>
    </>
  );
};
