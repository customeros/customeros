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

const getBillingCycleLabel = (cycleInMonths: number) => {
  switch (cycleInMonths) {
    case 0:
      return billingCycleLabels.NONE;
    case 1:
      return billingCycleLabels.MONTHLY_BILLING;
    case 3:
      return billingCycleLabels.QUARTERLY_BILLING;
    case 12:
      return billingCycleLabels.ANNUAL_BILLING;
    default:
      return billingCycleLabels.CUSTOM_BILLING;
  }
};

export const BillingCycleCell = observer(({ id }: { id: string }) => {
  const store = useStore();

  const billingCycle = store.invoices?.toArray().find((i) => i.id === id)?.value
    .contract.billingDetails?.billingCycleInMonths;

  return (
    <div className={cn(billingCycle ? 'text-gray-700' : 'text-gray-500')}>
      {billingCycle ? getBillingCycleLabel(billingCycle) : 'Unknown'}
    </div>
  );
});
