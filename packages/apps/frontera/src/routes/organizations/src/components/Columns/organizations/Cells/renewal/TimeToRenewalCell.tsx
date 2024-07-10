import { DateTimeUtils } from '@utils/date.ts';

interface TimeToRenewalCellProps {
  nextRenewalDate: string;
}

export const TimeToRenewalCell = ({
  nextRenewalDate,
}: TimeToRenewalCellProps) => {
  if (!nextRenewalDate)
    return <span className='text-sm text-gray-400' data-test='organization-next-renewal-in-all-orgs-table'>Unknown</span>;
  const [value, unit] = DateTimeUtils.getDifferenceFromNow(nextRenewalDate);
  const isNegative = value && value < 0;

  return (
    <span className='text-sm text-gray-700'>
      {value ? Math.abs(value) : ''} {unit} {isNegative ? 'ago' : ''}
    </span>
  );
};
