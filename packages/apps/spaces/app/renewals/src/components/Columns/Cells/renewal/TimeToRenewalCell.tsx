import { getDifferenceFromNow } from '@shared/util/date';

interface TimeToRenewalCellProps {
  nextRenewalDate: string;
}

export const TimeToRenewalCell = ({
  nextRenewalDate,
}: TimeToRenewalCellProps) => {
  if (!nextRenewalDate) return <p className='text-sm text-gray-400'>Unknown</p>;
  const [value, unit] = getDifferenceFromNow(nextRenewalDate);

  return (
    <div className='text-sm text-gray-700'>
      {value} {unit}
    </div>
  );
};
