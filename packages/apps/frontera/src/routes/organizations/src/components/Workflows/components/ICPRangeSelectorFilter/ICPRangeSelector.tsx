import { useCallback } from 'react';

import { toJS } from 'mobx';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Input, InputProps } from '@ui/form/Input';
import {
  WorkflowType,
  ColumnViewType,
  ComparisonOperator,
} from '@shared/types/__generated__/graphql.types';

interface RangeSelectorProps extends Omit<InputProps, 'onChange'> {
  filter: string;
  years?: boolean;
  property: string;
  placeholder: string;
}

const formatNumberWithCommas = (value: string | number | undefined): string => {
  if (value === undefined || value === '' || value === 0) return '';
  const numString = value?.toString();

  return numString?.replace(/\B(?=(\d{3})+(?!\d))/g, ',');
};

export const ICPRangeSelector = observer(
  ({
    filter,
    placeholder,
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

    const minValue = years
      ? workFlow?.getFilter(`${property}`)?.value[0]
        ? new Date().getFullYear() -
          workFlow?.getFilter(`${property}`)?.value[0]
        : ''
      : formatNumberWithCommas(workFlow?.getFilter(`${property}`)?.value[0]);

    const maxValue = years
      ? workFlow?.getFilter(`${property}`)?.value[1]
        ? new Date().getFullYear() -
          workFlow?.getFilter(`${property}`)?.value[1]
        : ''
      : formatNumberWithCommas(workFlow?.getFilter(`${property}`)?.value[1]);

    const typeOfFilter =
      filter === 'between'
        ? ComparisonOperator.Between
        : filter === 'less than'
        ? ComparisonOperator.Lt
        : ComparisonOperator.Gt;

    const handleChage = useCallback(
      (e: React.ChangeEvent<HTMLInputElement>, isMaxTrue: boolean) => {
        const value = Number(e.target.value);
        workFlow?.update((v) => {
          v.live = false;

          return v;
        });
        if (value !== undefined && !isMaxTrue) {
          workFlow?.setFilter({
            property: property,
            value:
              typeOfFilter === ComparisonOperator.Between
                ? [value, workFlow?.getFilter(`${property}`)?.value[1]]
                : [value],
            operation: typeOfFilter,
          });
        }

        if (value !== undefined && isMaxTrue) {
          workFlow?.setFilter({
            property: property,
            value: [workFlow?.getFilter(`${property}`)?.value[0], value],
            operation: typeOfFilter,
          });
        }

        if (isMaxTrue && value === undefined) {
          workFlow?.removeFilter(property);
        }

        if (value === undefined && !isMaxTrue) {
          workFlow?.removeFilter(
            ColumnViewType.OrganizationsLinkedinFollowerCount,
          );
        }
      },
      [workFlow, property, filter],
    );

    const handleYearsChange = useCallback(
      (e: React.ChangeEvent<HTMLInputElement>, isMaxTrue: boolean) => {
        workFlow?.update((v) => {
          v.live = false;

          return v;
        });
        const value = Number(e.target.value);
        if (value !== undefined && !isMaxTrue) {
          workFlow?.setFilter({
            property: ColumnViewType.OrganizationsYearFounded,
            value:
              typeOfFilter === ComparisonOperator.Between
                ? [
                    new Date().getFullYear() - Number(value),
                    workFlow?.getFilter(`${property}`)?.value[1] || '',
                  ]
                : [new Date().getFullYear() - Number(value)],
            operation: typeOfFilter,
          });
        }

        if (value !== undefined && isMaxTrue) {
          workFlow?.setFilter({
            property: ColumnViewType.OrganizationsYearFounded,
            value: [
              workFlow?.getFilter(`${property}`)?.value[0] || '',
              new Date().getFullYear() - Number(value),
            ],
            operation: typeOfFilter,
          });
        }

        if (isMaxTrue && value === undefined) {
          workFlow?.removeFilter(ColumnViewType.OrganizationsYearFounded);
        }
        if (value === undefined && !isMaxTrue) {
          workFlow?.removeFilter(ColumnViewType.OrganizationsYearFounded);
        }
      },
      [workFlow, property, filter],
    );

    return (
      <div className='flex-1 flex items-center'>
        <Input
          variant='unstyled'
          type='text'
          value={minValue || ''}
          placeholder={filter === 'between' ? 'Min' : `${placeholder}`}
          style={{
            width: filter !== 'between' && !years ? '100%' : '50px',
          }}
          onChange={(e) => {
            years ? handleYearsChange(e, false) : handleChage(e, false);
          }}
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
          value={maxValue}
          onChange={(e) =>
            years ? handleYearsChange(e, true) : handleChage(e, true)
          }
          {...rest}
        />
        {filter === 'between' && years && <span>yrs</span>}
      </div>
    );
  },
);
