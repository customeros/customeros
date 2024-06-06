import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { useStore } from '@shared/hooks/useStore';
import { ContractBillingCycle } from '@graphql/types';

const billingCycleLabels: Record<ContractBillingCycle, string> = {
  ANNUAL_BILLING: 'Annually',
  MONTHLY_BILLING: 'Monthly',
  QUARTERLY_BILLING: 'Quarterly',
  CUSTOM_BILLING: 'Custom',
  NONE: 'None',
};

export const BillingCycleCell = observer(
  ({ contractId }: { contractId: string }) => {
    const store = useStore();
    const billingCycle =
      store.contracts?.value?.get(contractId)?.value?.billingDetails
        ?.billingCycle;

    return (
      <span className={cn(billingCycle ? 'text-gray-700' : 'text-gray-500')}>
        {billingCycle ? billingCycleLabels[billingCycle] : 'Unknown'}
      </span>
    );
  },
);
