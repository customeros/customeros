import { useSearchParams } from 'react-router-dom';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import {
  ColumnViewType,
  ComparisonOperator,
  OrganizationRelationship,
} from '@graphql/types';
import { relationshipOptions } from '@organizations/components/Columns/organizations/Cells/relationship/util.ts';

import { FilterHeader } from '../../../shared/Filters/abstract/FilterHeader';

const defaultFilter: FilterItem = {
  property: ColumnViewType.OrganizationsRelationship,
  value: [],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.In,
};

export const RelationshipFilter = observer(() => {
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');

  const store = useStore();
  const tableViewDef = store.tableViewDefs.getById(preset ?? '');
  const filter =
    tableViewDef?.getFilter(defaultFilter.property) ?? defaultFilter;

  const toggle = () => {
    tableViewDef?.toggleFilter(filter);
  };

  const handleSelect = (value: OrganizationRelationship) => () => {
    const newValue = filter.value.includes(value)
      ? filter.value.filter((v: OrganizationRelationship) => v !== value)
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
        {relationshipOptions.map((option) => (
          <Checkbox
            key={option.value.toString()}
            onChange={handleSelect(option.value)}
            isChecked={filter.value.includes(option.value)}
          >
            <p className='text-sm'>{option.label}</p>
          </Checkbox>
        ))}
      </div>
    </>
  );
});
