import { useSearchParams } from 'react-router-dom';
import { memo, useRef, useMemo, useState, useEffect } from 'react';

import { produce } from 'immer';
import { FilterItem } from '@store/types';
import { subDays } from 'date-fns/subDays';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { Radio, RadioGroup } from '@ui/form/Radio/Radio';
import {
  ColumnViewType,
  LastTouchpointType,
  ComparisonOperator,
} from '@graphql/types';

import { allTime, touchpoints } from './util';
import { FilterHeader } from '../../../shared/Filters/abstract';

const defaultFilter: FilterItem = {
  property: ColumnViewType.OrganizationsLastTouchpoint,
  value: { after: '', types: [] },
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.Eq,
};

export const LastTouchpointFilter = observer(() => {
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');

  const store = useStore();
  const tableViewDef = store.tableViewDefs.getById(preset ?? '');
  const filter =
    tableViewDef?.getFilter(defaultFilter.property) ?? defaultFilter;

  const toggle = () => {
    tableViewDef?.toggleFilter(filter);
  };

  const [week, month, quarter] = useMemo(
    () =>
      [7, 30, 90].map((value) => {
        return subDays(new Date(), value).toISOString().split('T')[0];
      }),
    [],
  );

  const isAllSelected =
    filter.value.types.length === touchpoints.length &&
    filter.value.types.length > 0;

  const handleSelectAll = () => {
    if (isAllSelected) {
      tableViewDef?.setFilter({
        ...filter,
        active: true,
        value: {
          ...filter.value,
          types: [],
        },
      });
    } else {
      tableViewDef?.setFilter({
        ...filter,
        active: true,
        value: {
          ...filter.value,
          types: touchpoints.map(({ value }) => value),
        },
      });
    }
  };

  const handleSelect = (value: LastTouchpointType) => {
    const index = filter.value.types?.indexOf(value);

    if (index === -1) {
      tableViewDef?.setFilter({
        ...filter,
        active: true,
        value: {
          ...filter.value,
          types: [...filter.value.types, value],
        },
      });
    } else {
      tableViewDef?.setFilter({
        ...filter,
        active: true,
        value: {
          ...filter.value,
          types: filter.value.types.filter(
            (v: LastTouchpointType) => v !== value,
          ),
        },
      });
    }
  };

  const handleDateChange = (value: string) => {
    tableViewDef?.setFilter({
      ...filter,
      value: {
        ...filter.value,
        after: value,
      },
      active: filter.active || true,
    });
  };

  return (
    <>
      <FilterHeader
        onToggle={toggle}
        onDisplayChange={() => {}}
        isChecked={filter.active ?? false}
      />

      <RadioGroup
        value={filter.value.after}
        name='last-touchpoint-before'
        onValueChange={handleDateChange}
        className='border-b pb-2 border-gray-200'
      >
        <div className='flex flex-col gap-2 items-start'>
          <Radio value={week}>
            <span className='text-sm'>Last 7 days</span>
          </Radio>
          <Radio value={month}>
            <span className='text-sm'>Last 30 days</span>
          </Radio>
          <Radio value={quarter}>
            <span className='text-sm'>Last 90 days</span>
          </Radio>
          <Radio value={allTime}>
            <span className='text-sm'>All time</span>
          </Radio>
        </div>
      </RadioGroup>

      <div className='flex flex-col space-y-2 items-start mt-2 px-[4px] mx-[-4px] relative overflow-x-hidden overflow-y-auto'>
        <Checkboxes
          onCheck={handleSelect}
          value={filter.value.types}
          onCheckAll={handleSelectAll}
          isAllSelected={isAllSelected}
        />
      </div>
    </>
  );
});

interface CheckboxOptionsProps {
  value: string[];
  isDisabled?: boolean;
  isAllSelected: boolean;
  onCheckAll: () => void;
  onCheck: (value: LastTouchpointType) => void;
}

const makeState = (values: string[]) =>
  values.reduce((acc, curr) => ({ ...acc, [curr]: true }), {});

const allCheckedState = touchpoints.reduce(
  (acc, { value }) => ({ ...acc, [value]: true }),
  {},
);
const allUnchecked = touchpoints.reduce(
  (acc, { value }) => ({ ...acc, [value]: false }),
  {},
);

const Checkboxes = memo(
  ({
    value = [],
    onCheck,
    onCheckAll,
    isAllSelected,
  }: CheckboxOptionsProps) => {
    const timeoutRef = useRef<NodeJS.Timeout>();
    const [_isAllChecked, setIsAllChecked] = useState(() => isAllSelected);
    const [checked, setChecked] = useState<Record<string, boolean>>(
      makeState(value),
    );

    const handleCheck = (v: string) => {
      setChecked((prev) =>
        produce(prev, (draft) => {
          draft[v] = !draft[v];
        }),
      );

      timeoutRef.current && clearTimeout(timeoutRef.current);
      timeoutRef.current = setTimeout(
        () => onCheck(v as LastTouchpointType),
        250,
      );
    };

    const handleCheckAll = () => {
      setIsAllChecked((prev) => !prev);
      setChecked(_isAllChecked ? allUnchecked : allCheckedState);

      timeoutRef.current && clearTimeout(timeoutRef.current);
      timeoutRef.current = setTimeout(onCheckAll, 250);
    };

    useEffect(() => {
      setIsAllChecked(Object.values(checked).every((v) => v));
    }, [checked]);

    return (
      <>
        <div className='sticky top-0 w-full z-10 bg-white border-b border-gray-200 pb-2'>
          <Checkbox
            size='sm'
            onChange={handleCheckAll}
            isChecked={_isAllChecked}
          >
            <span className='text-sm'>
              {_isAllChecked ? 'Deselect all' : 'Select all'}
            </span>
          </Checkbox>
        </div>
        {touchpoints.map(({ label, value }) => (
          <Checkbox
            size='sm'
            key={value}
            iconSize='md'
            isChecked={checked[value]}
            onChange={() => handleCheck(value)}
            className='rounded-sm border border-gray-200'
          >
            <span className='text-sm line-clamp-1'>{label}</span>
          </Checkbox>
        ))}
      </>
    );
  },
);
