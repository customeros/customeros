import { RefObject } from 'react';
import { useSearchParams } from 'react-router-dom';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Input, ResizableInput } from '@ui/form/Input';
import { Radio, RadioGroup } from '@ui/form/Radio/Radio';
import { FilterHeader } from '@shared/components/Filters';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';

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
        active: true,
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
          disabled={!filter.active}
          onValueChange={handleChange}
        >
          <div className='gap-2 flex flex-col items-start'>
            <RadioGroup
              value={filter.operation}
              onValueChange={(newType) =>
                handleOperatorChange(newType as ComparisonOperator)
              }
            >
              <Radio value={ComparisonOperator.Lt}>
                <label className='text-sm'>Less than</label>
              </Radio>
              <Radio value={ComparisonOperator.Gt}>
                <label className='text-sm'>More than</label>
              </Radio>
              <Radio value={ComparisonOperator.Between}>
                <label className='text-sm'>Between</label>
              </Radio>
            </RadioGroup>
          </div>

          <div>
            {(filter.operation === ComparisonOperator.Lt ||
              filter.operation === ComparisonOperator.Gt) && (
              <div>
                <label className='font-semibold text-sm capitalize flex flex-col'>
                  {label}

                  {suffix ? (
                    <div>
                      <ResizableInput
                        step={1}
                        size='xs'
                        type='number'
                        name='contacts-count'
                        value={filter.value ?? '0'}
                        onFocus={(e) => e.target.select()}
                        className='text-gray-700 font-normal min-h-3'
                        onChange={(e) => handleChange(e.target.value)}
                        placeholder={
                          filter.operation === ComparisonOperator.Lt
                            ? 'Max'
                            : `Min`
                        }
                        defaultValue={
                          filter.operation === ComparisonOperator.Lt
                            ? 'Max'
                            : `Min`
                        }
                      />
                      <span className='font-normal ml-1 lowercase'>
                        {filter.value === '1' ? suffix : `${suffix}s`}
                      </span>
                    </div>
                  ) : (
                    <Input
                      step={1}
                      size='xs'
                      type='number'
                      name='contacts-count'
                      defaultValue={filter.value ?? ''}
                      onFocus={(e) => e.target.select()}
                      placeholder={`Number of ${label}`}
                      className='text-gray-700 font-normal'
                      onChange={(e) => handleChange(e.target.value)}
                    />
                  )}
                </label>
              </div>
            )}

            {filter.operation === ComparisonOperator.Between && (
              <div className='flex justify-between gap-2'>
                <label className='font-semibold text-sm flex flex-col w-[50%] gap-1'>
                  Min {label}
                  {suffix ? (
                    <div>
                      <ResizableInput
                        step={1}
                        size='xs'
                        name='name'
                        placeholder={`Min`}
                        defaultValue={'Min'}
                        value={filter.value[0] ?? ''}
                        onFocus={(e) => e.target.select()}
                        className='text-gray-700 font-normal min-h-3 '
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
                      step={1}
                      size='xs'
                      name='name'
                      placeholder={`Min`}
                      value={filter.value[0] ?? ''}
                      onFocus={(e) => e.target.select()}
                      className='text-gray-700 font-normal'
                      onChange={(e) =>
                        handleChange([e.target.value, filter.value?.[1]])
                      }
                    />
                  )}
                </label>
                <label className='font-semibold text-sm flex flex-col w-[50%] gap-1'>
                  Max {label}
                  {suffix ? (
                    <div>
                      <ResizableInput
                        step={1}
                        size='xs'
                        name='name'
                        placeholder='Max'
                        defaultValue={'Max'}
                        value={filter.value[1] ?? ''}
                        onFocus={(e) => e.target.select()}
                        className='text-gray-700 font-normal min-h-3 '
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
                      step={1}
                      size='xs'
                      name='name'
                      placeholder='Max'
                      value={filter.value[1] ?? ''}
                      onFocus={(e) => e.target.select()}
                      className='text-gray-700 font-normal'
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
