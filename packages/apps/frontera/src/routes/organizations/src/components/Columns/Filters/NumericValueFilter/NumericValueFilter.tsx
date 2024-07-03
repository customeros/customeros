import { RefObject } from 'react';
import { useSearchParams } from 'react-router-dom';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Input, ResizableInput } from '@ui/form/Input';
import { Radio, RadioGroup } from '@ui/form/Radio/Radio';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';

import { FilterHeader } from '../shared';

const defaultFilter: FilterItem = {
  property: ColumnViewType.OrganizationsContactCount,
  value: 0,
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.Lte,
};

interface NumericValueFilterProps {
  label: string;
  suffix?: string;
  property?: string;

  initialFocusRef?: RefObject<HTMLInputElement>;
}
export const NumericValueFilter = observer(
  ({ label, property, suffix }: NumericValueFilterProps) => {
    const [searchParams] = useSearchParams();
    const preset = searchParams.get('preset');

    const store = useStore();
    const tableViewDef = store.tableViewDefs.getById(preset ?? '');
    const filter = tableViewDef?.getFilter(
      property || defaultFilter.property,
    ) ?? { ...defaultFilter, property: property || defaultFilter.property };

    const toggle = () => {
      tableViewDef?.toggleFilter(filter);
    };

    const handleChange = (value: string | string[]) => {
      tableViewDef?.setFilter({
        ...filter,
        value,
        active: true,
      });
    };
    const handleOperatorChange = (operation: ComparisonOperator) => {
      tableViewDef?.setFilter({
        ...filter,
        operation,
        value:
          operation === ComparisonOperator.Between
            ? [0, filter.value]
            : Array.isArray(filter.value)
            ? filter.value[1]
            : filter.value,
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
          name='timeToRenewal'
          value={filter.value}
          onValueChange={handleChange}
          disabled={!filter.active}
        >
          <div className='gap-2 flex flex-col items-start'>
            <RadioGroup
              value={filter.operation}
              onValueChange={(newType) =>
                handleOperatorChange(newType as ComparisonOperator)
              }
            >
              <Radio value={ComparisonOperator.Lte}>
                <label className='text-sm'>Less than</label>
              </Radio>
              <Radio value={ComparisonOperator.Gte}>
                <label className='text-sm'>More than</label>
              </Radio>
              <Radio value={ComparisonOperator.Between}>
                <label className='text-sm'>Between</label>
              </Radio>
            </RadioGroup>
          </div>

          <div>
            {(filter.operation === ComparisonOperator.Lte ||
              filter.operation === ComparisonOperator.Gte) && (
              <div>
                <label className='font-semibold text-sm capitalize flex flex-col'>
                  {label}

                  {suffix ? (
                    <div>
                      <ResizableInput
                        className='text-gray-700 font-normal min-h-3'
                        name='contacts-count'
                        type='number'
                        size='xs'
                        step={1}
                        onFocus={(e) => e.target.select()}
                        placeholder={`${label}`}
                        value={filter.value ?? '0'}
                        defaultValue={'0'}
                        onChange={(e) => handleChange(e.target.value)}
                      />
                      <span className='font-normal ml-1 lowercase'>
                        {filter.value === '1' ? suffix : `${suffix}s`}
                      </span>
                    </div>
                  ) : (
                    <Input
                      className='text-gray-700 font-normal'
                      name='contacts-count'
                      type='number'
                      size='xs'
                      step={1}
                      onFocus={(e) => e.target.select()}
                      placeholder={`${label}`}
                      defaultValue={filter.value ?? ''}
                      onChange={(e) => handleChange(e.target.value)}
                    />
                  )}
                </label>
              </div>
            )}

            {filter.operation === ComparisonOperator.Between && (
              <div className='flex min-w-[280px] justify-between'>
                <label className='font-semibold text-sm flex flex-col w-[50%]'>
                  Min {label}
                  {suffix ? (
                    <div>
                      <ResizableInput
                        className='text-gray-700 font-normal min-h-3 '
                        name='name'
                        size='xs'
                        step={1}
                        onFocus={(e) => e.target.select()}
                        placeholder={`min ${label}`}
                        value={filter.value[0] ?? ''}
                        defaultValue={'0'}
                        onChange={(e) =>
                          handleChange([e.target.value, filter.value?.[1]])
                        }
                      />
                      <span className='font-normal ml-1'>
                        {filter.value?.[0] === '1' ? suffix : `${suffix}s`}
                      </span>
                    </div>
                  ) : (
                    <Input
                      className='text-gray-700 font-normal'
                      name='name'
                      size='xs'
                      step={1}
                      onFocus={(e) => e.target.select()}
                      placeholder={`min ${label}`}
                      value={filter.value[0] ?? ''}
                      onChange={(e) =>
                        handleChange([e.target.value, filter.value?.[1]])
                      }
                    />
                  )}
                </label>
                <label className='font-semibold text-sm flex flex-col w-[50%]'>
                  Max {label}
                  {suffix ? (
                    <div>
                      <ResizableInput
                        className='text-gray-700 font-normal min-h-3 '
                        name='name'
                        size='xs'
                        step={1}
                        defaultValue={'0'}
                        onFocus={(e) => e.target.select()}
                        placeholder={`min ${label}`}
                        value={filter.value[1] ?? ''}
                        onChange={(e) =>
                          handleChange([filter.value?.[0], e.target.value])
                        }
                      />
                      <span className='font-normal ml-1'>
                        {filter.value?.[1] === '1' ? suffix : `${suffix}s`}
                      </span>
                    </div>
                  ) : (
                    <Input
                      className='text-gray-700 font-normal'
                      name='name'
                      size='xs'
                      step={1}
                      onFocus={(e) => e.target.select()}
                      placeholder={`min ${label}`}
                      value={filter.value[1] ?? ''}
                      onChange={(e) =>
                        handleChange([filter.value?.[0], e.target.value])
                      }
                    />
                  )}
                </label>
              </div>
            )}
          </div>
        </RadioGroup>
      </>
    );
  },
);
