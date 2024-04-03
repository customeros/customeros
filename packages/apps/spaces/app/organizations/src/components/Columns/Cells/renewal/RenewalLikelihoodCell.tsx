import { cn } from '@ui/utils/cn';
import { OpportunityRenewalLikelihood } from '@graphql/types';

import { getLikelihoodColor, getRenewalLikelihoodLabel } from './utils';

interface RenewalLikelihoodCellProps {
  value?: OpportunityRenewalLikelihood | null;
}

export const RenewalLikelihoodCell = ({
  value,
}: RenewalLikelihoodCellProps) => {
  const colors = value ? getLikelihoodColor(value) : 'text-gray-400';

  return (
    <div className='w-full' key={Math.random()}>
      <span className={cn('cursor-default', colors)}>
        {value ? getRenewalLikelihoodLabel(value) : 'Unknown'}
      </span>
    </div>
  );
};
