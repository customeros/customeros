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

import { DebouncedSearchInput } from '../shared';
import { FilterHeader } from '../shared/FilterHeader';

interface IndustryFilterProps {
  initialFocusRef: RefObject<HTMLInputElement>;
}

const defaultFilter: FilterItem = {
  property: ColumnViewType.OrganizationsIndustry,
  value: [],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.In,
};

export const IndustryFilter = observer(
  ({ initialFocusRef }: IndustryFilterProps) => {
    const [searchParams] = useSearchParams();
    const preset = searchParams.get('preset');

    const [searchValue, setSearchValue] = useState('');

    const store = useStore();
    const tableViewDef = store.tableViewDefs.getById(preset ?? '');
    const filter =
      tableViewDef?.getFilter(defaultFilter.property) ?? defaultFilter;
    const options = uniqBy(store.organizations.toArray(), 'value.industry')
      .map((v) => v.value.industry)
      .filter(Boolean)
      .filter((v) => {
        if (searchValue) {
          return v?.includes(searchValue);
        }

        return true;
      }) as string[];

    const toggle = () => {
      tableViewDef?.toggleFilter(filter);
    };

    const isAllSelected =
      intersection(filter.value, options).length === options.length &&
      options.length > 0;

    const handleSelect = (value: string) => () => {
      const newValue = filter.value.includes(value)
        ? filter.value.filter((v: string) => v !== value)
        : [...filter.value, value];

      tableViewDef?.setFilter({
        ...filter,
        value: newValue,
        active: newValue.length > 0,
      });
    };

    const handleSelectAll = () => {
      let nextValue: string[] = [];

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
          onChange={(v) => setSearchValue(v)}
          onDisplayChange={(v) => setSearchValue(v)}
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
                isChecked={filter.value.includes(option)}
                onChange={handleSelect(option ?? '')}
              >
                <span className='text-sm line-clamp-1'>{option}</span>
              </Checkbox>
            ))
          ) : (
            <div className='flex w-full justify-center items-center flex-col'>
              <Tumbleweed className='mr-10 size-8 text-gray-400 self-end' />
              <span className='text-sm text-gray-500'>
                Empty here in <b>No Resultsville</b>
              </span>
            </div>
          )}
        </div>
      </>
    );
  },
);
