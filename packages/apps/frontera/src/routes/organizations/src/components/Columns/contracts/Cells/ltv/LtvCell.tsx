import { formatCurrency } from '@utils/getFormattedCurrencyNumber.ts';

interface LtvCellProps {
  ltv: number;
  currency: string;
}

export const LtvCell = ({ currency, ltv }: LtvCellProps) => {
  if (!ltv) {
    return <p className='text-gray-400'>Unknown</p>;
  }

  return (
    <div className='flex items-center'>
      {formatCurrency(ltv, 2, currency || 'USD')}
    </div>
  );
};
