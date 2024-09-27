import { getCommittedPeriodLabel } from '@shared/util/committedPeriodLabel';

interface PeriodCellProps {
  className?: string;
  committedPeriodInMonths?: string | number;
}

export const PeriodCell = ({ committedPeriodInMonths }: PeriodCellProps) => {
  if (!committedPeriodInMonths) {
    return <p className='text-gray-400'> {committedPeriodInMonths}</p>;
  }

  return (
    <div className='flex items-center'>
      {getCommittedPeriodLabel(committedPeriodInMonths)}
    </div>
  );
};
