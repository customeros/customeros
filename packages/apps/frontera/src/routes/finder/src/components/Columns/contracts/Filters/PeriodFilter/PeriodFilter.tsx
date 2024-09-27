import { useState, RefObject } from 'react';
import { useSearchParams } from 'react-router-dom';

import uniqBy from 'lodash/uniqBy';
import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';
import difference from 'lodash/difference';
import intersection from 'lodash/intersection';

import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { Tumbleweed } from '@ui/media/icons/Tumbleweed';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';
import { getCommittedPeriodLabel } from '@shared/util/committedPeriodLabel.ts';

import { DebouncedSearchInput } from '../../../shared/Filters/abstract';
import { FilterHeader } from '../../../shared/Filters/abstract/FilterHeader';

interface IndustryFilterProps {
  initialFocusRef: RefObject<HTMLInputElement>;
}

const defaultFilter: FilterItem = {
  property: ColumnViewType.ContractsPeriod,
  value: [],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.In,
};

export const PeriodFilter = observer(
  ({ initialFocusRef }: IndustryFilterProps) => {
    const [searchParams] = useSearchParams();
    const preset = searchParams.get('preset');

    const [searchValue, setSearchValue] = useState('');

    const store = useStore();
    const tableViewDef = store.tableViewDefs.getById(preset ?? '');
    const filter =
      tableViewDef?.getFilter(defaultFilter.property) ?? defaultFilter;
    const options = uniqBy(
      store.contracts.toArray(),
      'value.committedPeriodInMonths',
    )
      .map((v) => v.value.committedPeriodInMonths)
      .filter(Boolean)
      .sort((a, b) => a - b)
      .filter((v) => {
        if (searchValue) {
          return getCommittedPeriodLabel(v)
            ?.toLowerCase()
            ?.includes(searchValue?.toLowerCase() as string);
        }

        return true;
      }) as number[];

    const toggle = () => {
      tableViewDef?.toggleFilter(filter);
    };

    const isAllSelected =
      intersection(filter.value, options).length === options.length &&
      options.length > 0;

    const handleSelect = (value: number) => () => {
      const newValue = filter.value.includes(value)
        ? filter.value.filter((v: number) => v !== value)
        : [...filter.value, value];

      tableViewDef?.setFilter({
        ...filter,
        value: newValue,
        active: newValue.length > 0,
      });
    };

    const handleSelectAll = () => {
      let nextValue: number[] = [];

      if (isAllSelected) {
        tableViewDef?.setFilter({
          ...filter,
          value: difference(filter.value, options),
          active: false,
        });

        return;
      }

      if (searchValue) {
        nextValue = [...options, ...difference(filter.value, options)];
      } else {
        nextValue = options;
      }

      tableViewDef?.setFilter({
        ...filter,
        value: nextValue,
        active: nextValue.length > 0,
      });
    };

    return (
      <>
        <FilterHeader
          onToggle={toggle}
          onDisplayChange={() => {}}
          isChecked={filter.active ?? false}
        />
        <DebouncedSearchInput
          value={searchValue}
          ref={initialFocusRef}
          onDisplayChange={(v) => setSearchValue(v)}
          onChange={(v) => {
            setSearchValue(v);

            if ((v.length && !filter.active) || (!v.length && filter.active)) {
              toggle();
            }
          }}
        />

        <div className='flex flex-col w-full h-[13rem] items-start gap-2 mt-2 px-1 mx-[-4px] overflow-x-hidden overflow-y-auto'>
          {options.length > 1 && (
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
                      (searchValue && options.length > 2
                        ? ` ${options.length}`
                        : '')}
                </span>
              </Checkbox>
            </div>
          )}

          {options.length > 0 ? (
            options.map((option) => (
              <Checkbox
                key={option}
                onChange={handleSelect(option)}
                isChecked={filter.value.includes(option)}
              >
                <span className='text-sm line-clamp-1'>
                  {getCommittedPeriodLabel(option)}
                </span>
              </Checkbox>
            ))
          ) : (
            <div className='flex w-full justify-center items-center flex-col'>
              <Tumbleweed className='size-8 text-gray-400' />
              <span className='text-center text-sm text-gray-500'>
                Empty here in <b>No Resultsville</b>
              </span>
            </div>
          )}
        </div>
      </>
    );
  },
);
