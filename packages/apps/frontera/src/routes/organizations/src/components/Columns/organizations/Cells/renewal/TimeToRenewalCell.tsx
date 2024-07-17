import { DateTimeUtils } from '@utils/date';

interface TimeToRenewalCellProps {
  nextRenewalDate: string;
}

export const TimeToRenewalCell = ({
  nextRenewalDate,
}: TimeToRenewalCellProps) => {
  if (!nextRenewalDate)
    return (
      <span
        className='text-sm text-gray-400'
        data-test='organization-next-renewal-in-all-orgs-table'
      >
        No contract
      </span>
    );

  return (
    <span className='text-sm text-gray-700'>
      {DateTimeUtils.format(
        nextRenewalDate,
        DateTimeUtils.dateWithAbreviatedMonth,
      )}
    </span>
  );
};
