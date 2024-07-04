import { useSearchParams } from 'react-router-dom';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import {
  ColumnViewType,
  ComparisonOperator,
  OpportunityRenewalLikelihood,
} from '@graphql/types';

import { FilterHeader } from '../../../shared/Filters/abstract';

const defaultFilter: FilterItem = {
  property: ColumnViewType.OrganizationsRenewalLikelihood,
  value: [],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.In,
};

export const RenewalLikelihoodFilter = observer(() => {
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');

  const store = useStore();
  const tableViewDef = store.tableViewDefs.getById(preset ?? '');
  const filter =
    tableViewDef?.getFilter(defaultFilter.property) ?? defaultFilter;

  const toggle = () => {
    tableViewDef?.toggleFilter(filter);
  };

  const handleSelect = (value: OpportunityRenewalLikelihood) => () => {
    tableViewDef?.setFilter({
      ...filter,
      value: filter.value.includes(value)
        ? filter.value.filter((v: OpportunityRenewalLikelihood) => v !== value)
        : [...filter.value, value],
      active: true,
    });
  };

  return (
    <>
      <FilterHeader
        onToggle={toggle}
        onDisplayChange={() => {}}
        isChecked={filter.active ?? false}
      />
      <div className='flex flex-col space-y-2 items-start'>
        <Checkbox
          isChecked={filter.value.includes(
            OpportunityRenewalLikelihood.HighRenewal,
          )}
          onChange={handleSelect(OpportunityRenewalLikelihood.HighRenewal)}
        >
          <span className='text-sm'>High</span>
        </Checkbox>
        <Checkbox
          isChecked={filter.value.includes(
            OpportunityRenewalLikelihood.MediumRenewal,
          )}
          onChange={handleSelect(OpportunityRenewalLikelihood.MediumRenewal)}
        >
          <span className='text-sm'>Medium</span>
        </Checkbox>
        <Checkbox
          isChecked={filter.value.includes(
            OpportunityRenewalLikelihood.LowRenewal,
          )}
          onChange={handleSelect(OpportunityRenewalLikelihood.LowRenewal)}
        >
          <span className='text-sm'>Low</span>
        </Checkbox>
        <Checkbox
          isChecked={filter.value.includes(
            OpportunityRenewalLikelihood.ZeroRenewal,
          )}
          onChange={handleSelect(OpportunityRenewalLikelihood.ZeroRenewal)}
        >
          <span className='text-sm'>Zero</span>
        </Checkbox>
      </div>
    </>
  );
});
