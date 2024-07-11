import { useState, useEffect, useCallback } from 'react';

import { toJS } from 'mobx';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Input, InputProps } from '@ui/form/Input';
import { WorkflowType } from '@shared/types/__generated__/graphql.types';

interface RangeSelectorProps extends Omit<InputProps, 'onChange'> {
  filter: string;
  years?: boolean;
  property: string;
  placeholder: string;
  onChange: (values: [number | string, (number | string)?]) => void;
}

const formatNumberWithCommas = (value: string | number | undefined): string => {
  if (value === undefined || value === '') return '';
  const numString = value.toString();

  return numString.replace(/\B(?=(\d{3})+(?!\d))/g, ',');
};

export const RangeSelector = observer(
  ({
    filter,
    placeholder,
    onChange,
    property,
    years = false,
    ...rest
  }: RangeSelectorProps) => {
    const store = useStore();
    const getWorkFlow = store.workFlows
      .toArray()
      .filter((wf) =>
        toJS(wf.value.type === WorkflowType.IdealCustomerProfile),
      );
    const getWorkFlowId = getWorkFlow.map((wf) => wf.value.id);

    const workFlow = store.workFlows.getByType(getWorkFlowId[0]);

    const [minValue, setMinValue] = useState<number | string | undefined>(
      years
        ? new Date().getFullYear() - workFlow?.getFilter(property)?.value[0]
        : workFlow?.getFilter(property)?.value[0] || undefined,
    );
    const [maxValue, setMaxValue] = useState<number | string | undefined>(
      years
        ? new Date().getFullYear() - workFlow?.getFilter(property)?.value[1]
        : workFlow?.getFilter(property)?.value[1] || undefined,
    );
    const [hasInputChanged, setHasInputChanged] = useState(false);

    useEffect(() => {
      if (filter === 'between') {
        onChange([minValue as number, maxValue as number]);
        if (hasInputChanged) {
          workFlow?.update((wf) => {
            wf.live = false;

            return wf;
          });
        }
      } else {
        onChange([minValue as number]);
        if (hasInputChanged) {
          workFlow?.update((wf) => {
            wf.live = false;

            return wf;
          });
        }
      }
    }, [minValue, maxValue, filter, hasInputChanged]);

    const handleMinChange = useCallback(
      (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = e.target.value.replace(/,/g, '');
        setMinValue(value ? Number(value) : undefined);
        setHasInputChanged(true);
      },
      [workFlow?.value.condition],
    );

    const handleMaxChange = useCallback(
      (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = e.target.value.replace(/,/g, '');
        setMaxValue(value ? Number(value) : undefined);
        setHasInputChanged(true);
      },
      [workFlow?.value.condition],
    );

    return (
      <div className='flex-1 flex items-center'>
        <Input
          variant='unstyled'
          type='text'
          value={years ? minValue : formatNumberWithCommas(minValue)}
          placeholder={filter === 'between' ? 'Min' : `${placeholder}`}
          style={{
            width: filter !== 'between' && !years ? '100%' : '50px',
          }}
          onChange={handleMinChange}
          {...rest}
        />
        {years && (
          <>
            <span>yrs</span>
            <span
              className='mx-4 '
              style={{
                display: filter === 'between' ? 'block' : 'none',
              }}
            >
              -{' '}
            </span>
          </>
        )}
        {!years && (
          <span
            className='mr-[30px]'
            style={{
              display: filter === 'between' ? 'block' : 'none',
            }}
          >
            -{' '}
          </span>
        )}
        <Input
          style={{
            display: filter === 'between' ? 'block' : 'none',
          }}
          variant='unstyled'
          type='text'
          placeholder='Max'
          className='w-[50px]'
          value={years ? maxValue || '' : formatNumberWithCommas(maxValue)}
          onChange={handleMaxChange}
          {...rest}
        />
        {filter === 'between' && years && <span>yrs</span>}
      </div>
    );
  },
);
