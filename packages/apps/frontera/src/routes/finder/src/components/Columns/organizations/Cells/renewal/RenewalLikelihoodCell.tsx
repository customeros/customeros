import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { useStore } from '@shared/hooks/useStore';

import { getLikelihoodColor, getRenewalLikelihoodLabel } from './utils';

interface RenewalLikelihoodCellProps {
  id: string;
}

export const RenewalLikelihoodCell = observer(
  ({ id }: RenewalLikelihoodCellProps) => {
    const store = useStore();
    const organization = store.organizations.getById(id);
    const value = organization?.value?.renewalSummaryRenewalLikelihood;

    const colors = value ? getLikelihoodColor(value) : 'text-gray-400';

    if (!organization?.value) return null;

    const canUpdate = organization.value.contracts?.length;

    return (
      <div
        className={cn(
          'flex gap-1 items-center cursor-default group/likelihood',
          {
            'cursor-pointer': canUpdate,
          },
        )}
        onClick={() => {
          if (canUpdate) {
            store.ui.commandMenu.setType('UpdateHealthStatus');
            store.ui.commandMenu.setOpen(true);
          }
        }}
      >
        <span
          className={cn(colors)}
          data-test='organization-health-in-all-orgs-table'
        >
          {value ? getRenewalLikelihoodLabel(value) : 'No set'}
        </span>
      </div>
    );
  },
);
