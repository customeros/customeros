import { cn } from '@ui/utils/cn';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';

interface RenewalForecastCellProps {
  amount?: number | null;
  potentialAmount?: number | null;
}

export const RenewalForecastCell = ({
  amount = null,
  potentialAmount = null,
}: RenewalForecastCellProps) => {
  const formattedAmount =
    amount !== null && amount >= 0 ? formatCurrency(amount, 0) : 'Unknown';
  const formattedPotentialAmount = formatCurrency(potentialAmount ?? 0, 0);

  const showPotentialAmount =
    amount !== null &&
    potentialAmount !== null &&
    formattedAmount !== formattedPotentialAmount;

  if (formattedAmount === 'Unknown')
    return <p className='text-gray-400'>Unknown</p>;

  return (
    <span className='flex flex-col justify-center'>
      <span
        className={cn('text-sm', amount ? 'text-gray-700' : 'text-gray-500')}
      >
        {formattedAmount}
      </span>
      {showPotentialAmount && (
        <span className='text-sm text-gray-500 line-through'>
          {formattedPotentialAmount}
        </span>
      )}
    </span>
  );
};
