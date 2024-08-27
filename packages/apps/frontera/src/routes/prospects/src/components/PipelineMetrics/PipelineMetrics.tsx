import { Currency } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';

interface PipelineMetricsProps {
  count: number;
  totalArr: number;
  currency: Currency;
  totalWeightedArr: number;
}

export const PipelineMetrics = ({
  currency,
  count = 0,
  totalArr = 0,
  totalWeightedArr = 0,
}: PipelineMetricsProps) => {
  const store = useStore();
  const numberOfColumns =
    store.settings.tenant.value?.opportunityStages.length ?? 0;

  return (
    <div
      style={{ minWidth: `${numberOfColumns * 150}px` }}
      className='
       sticky top-[42px] z-20 bg-white '
    >
      <div className='px-3 py-2 mx-4 mt-4 mb-4 bg-gray-100 flex justify-center gap-4 rounded-[4px] '>
        <span className=''>
          <span className='font-semibold'>{count}</span>{' '}
          <span className='text-gray-500'>opportunities</span>
        </span>
        <p className='font-semibold'>•</p>
        <span>
          <span className='font-semibold'>
            {formatCurrency(totalArr, 2, currency)}{' '}
          </span>
          <span className='text-gray-500 text-medium'>ARR estimate</span>
        </span>
        <p className='font-semibold'>•</p>
        <span className=''>
          <span className='font-semibold'>
            {formatCurrency(totalWeightedArr, 2, currency)}{' '}
          </span>
          <span className='text-gray-500 text-medium'>
            Weighted ARR estimate
          </span>
        </span>
      </div>
    </div>
  );
};
