import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { FilterItem } from '@store/types.ts';

import { useStore } from '@shared/hooks/useStore';
import { Checkbox } from '@ui/form/Checkbox/Checkbox.tsx';
import {
  ColumnViewType,
  ComparisonOperator,
  OpportunityRenewalLikelihood,
} from '@graphql/types';

import { FilterHeader } from '../abstract';

const defaultFilter: FilterItem = {
  property: ColumnViewType.OrganizationsRenewalLikelihood,
  value: [],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.In,
};

export const RenewalLikelihoodFilter = observer(
  ({ property }: { property?: ColumnViewType }) => {
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

    const handleSelect = (value: OpportunityRenewalLikelihood) => () => {
      tableViewDef?.setFilter({
        ...filter,
        value: filter.value.includes(value)
          ? filter.value.filter(
              (v: OpportunityRenewalLikelihood) => v !== value,
            )
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
            onChange={handleSelect(OpportunityRenewalLikelihood.HighRenewal)}
            isChecked={filter.value.includes(
              OpportunityRenewalLikelihood.HighRenewal,
            )}
          >
            <span className='text-sm'>High</span>
          </Checkbox>
          <Checkbox
            onChange={handleSelect(OpportunityRenewalLikelihood.MediumRenewal)}
            isChecked={filter.value.includes(
              OpportunityRenewalLikelihood.MediumRenewal,
            )}
          >
            <span className='text-sm'>Medium</span>
          </Checkbox>
          <Checkbox
            onChange={handleSelect(OpportunityRenewalLikelihood.LowRenewal)}
            isChecked={filter.value.includes(
              OpportunityRenewalLikelihood.LowRenewal,
            )}
          >
            <span className='text-sm'>Low</span>
          </Checkbox>
          <Checkbox
            onChange={handleSelect(OpportunityRenewalLikelihood.ZeroRenewal)}
            isChecked={filter.value.includes(
              OpportunityRenewalLikelihood.ZeroRenewal,
            )}
          >
            <span className='text-sm'>Zero</span>
          </Checkbox>
        </div>
      </>
    );
  },
);
