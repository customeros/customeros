interface PeriodCellProps {
  className?: string;
  committedPeriodInMonths?: string | number;
}

export function getCommittedPeriodLabel(months: string | number) {
  if (`${months}` === '1') {
    return 'Month';
  }

  if (`${months}` === '3') {
    return 'Quarter';
  }

  if (`${months}` === '12') {
    return 'Annual';
  }

  return `${months}-month`;
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
