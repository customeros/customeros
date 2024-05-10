import { cn } from '@ui/utils/cn';
import { ContractBillingCycle } from '@graphql/types';

const billingCycleLabels: Record<ContractBillingCycle, string> = {
  ANNUAL_BILLING: 'Annually',
  MONTHLY_BILLING: 'Monthly',
  QUARTERLY_BILLING: 'Quarterly',
  CUSTOM_BILLING: 'Custom',
  NONE: 'None',
};

export const BillingCycleCell = ({
  value,
}: {
  value: ContractBillingCycle;
}) => {
  return (
    <span className={cn(value ? 'text-gray-700' : 'text-gray-500')}>
      {value ? billingCycleLabels[value] : 'Unknown'}
    </span>
  );
};
