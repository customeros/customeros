import { getDifferenceFromNow } from '@shared/util/date';

interface TimeToRenewalCellProps {
  nextRenewalDate: string;
}

export const TimeToRenewalCell = ({
  nextRenewalDate,
}: TimeToRenewalCellProps) => {
  if (!nextRenewalDate)
    return <span className='text-sm text-gray-400'>Unknown</span>;
  const [value, unit] = getDifferenceFromNow(nextRenewalDate);
  const isNegative = `${value}`.includes('-');

  return (
    <span className='text-sm text-gray-700'>
      {`${value}`.replace('-', '')} {unit} {isNegative ? 'ago' : ''}
    </span>
  );
};
