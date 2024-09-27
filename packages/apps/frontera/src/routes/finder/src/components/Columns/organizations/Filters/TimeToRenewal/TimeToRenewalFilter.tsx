import { useSearchParams } from 'react-router-dom';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';

import { DateTimeUtils } from '@utils/date';
import { useStore } from '@shared/hooks/useStore';
import { Calendar } from '@ui/media/icons/Calendar';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';
import { DatePickerUnderline } from '@ui/form/DatePicker/DatePickerUnderline';

import { FilterHeader } from '../../../shared/Filters/abstract';

const defaultFilter: FilterItem = {
  property: ColumnViewType.OrganizationsRenewalDate,
  value: '',
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.Lt,
};

export const TimeToRenewalFilter = observer(() => {
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');

  const store = useStore();
  const tableViewDef = store.tableViewDefs.getById(preset ?? '');
  const filter =
    tableViewDef?.getFilter(defaultFilter.property) ?? defaultFilter;

  const toggle = () => {
    tableViewDef?.toggleFilter({
      ...filter,
    });

    if (filter.active) {
      tableViewDef?.removeFilter(filter.property);
    }
  };

  const handleChange = (date: Date, isMax: boolean) => {
    const parsedDate = DateTimeUtils.format(
      String(date),
      DateTimeUtils.iso8601,
    );

    if (!isMax) {
      tableViewDef?.setFilter({
        ...filter,
        value: [parsedDate, filter.value[1]],
        operation: ComparisonOperator.Gt,
      });
    } else {
      tableViewDef?.setFilter({
        ...filter,
        value: [filter.value[0], parsedDate],
        operation: ComparisonOperator.Between,
      });
    }
  };

  return (
    <>
      <FilterHeader
        onToggle={toggle}
        onDisplayChange={() => {}}
        isChecked={filter.active ?? false}
      />
      <div className='flex justify-between'>
        <div className='flex flex-col'>
          <label className='font-semibold text-sm'>From</label>
          <div className='flex items-center'>
            <Calendar className='mr-1 text-gray-500' />
            <DatePickerUnderline
              size='sm'
              value={filter.value[0]}
              onChange={(value) => {
                if (value) handleChange(value, false);
              }}
            />
          </div>
        </div>
        <div className='flex flex-col'>
          <label className='font-semibold text-sm'>To</label>
          <div className='flex items-center'>
            <Calendar className='mr-1 text-gray-500' />
            <DatePickerUnderline
              size='sm'
              value={filter.value[1]}
              minDate={new Date(filter.value[0])}
              onChange={(value) => {
                if (value) handleChange(value, true);
              }}
            />
          </div>
        </div>
      </div>
    </>
  );
});
