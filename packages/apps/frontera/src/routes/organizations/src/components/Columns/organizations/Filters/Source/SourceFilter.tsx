import { useSearchParams } from 'react-router-dom';

import uniqBy from 'lodash/uniqBy';
import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';

import { FilterHeader } from '../../../shared/Filters/abstract/FilterHeader';

const defaultFilter: FilterItem = {
  property: ColumnViewType.OrganizationsLeadSource,
  value: [],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.In,
};

export const SourceFilter = observer(() => {
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');

  const store = useStore();
  const tableViewDef = store.tableViewDefs.getById(preset ?? '');
  const filter =
    tableViewDef?.getFilter(defaultFilter.property) ?? defaultFilter;
  const options = uniqBy(store.organizations.toArray(), 'value.leadSource')
    .map((v) => v.value.leadSource)
    .filter(Boolean);

  const toggle = () => {
    tableViewDef?.toggleFilter(filter);
  };

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

  return (
    <>
      <FilterHeader
        onToggle={toggle}
        onDisplayChange={() => {}}
        isChecked={filter.active ?? false}
      />
      <div className='flex flex-col gap-2 items-start'>
        {options.map((option) => (
          <Checkbox
            key={option}
            isChecked={filter.value.includes(option)}
            onChange={handleSelect(option ?? '')}
          >
            <p className='text-sm'>{option}</p>
          </Checkbox>
        ))}
      </div>
    </>
  );
});
