import { ContractStatus } from '@graphql/types';

interface StatusCellProps {
  className?: string;
  status?: ContractStatus | null;
}

export function renderStatusNode(type: ContractStatus | null | undefined) {
  switch (type) {
    case ContractStatus.Draft:
      return <>Draft</>;

    case ContractStatus.Live:
      return <>Live</>;

    case ContractStatus.Ended:
      return <>Ended</>;

    case ContractStatus.OutOfContract:
      return <>Out of contract</>;

    case ContractStatus.Scheduled:
      return <>Scheduled</>;

    default:
      return '';
  }
}

export const StatusCell = ({ status }: StatusCellProps) => {
  return <div className='flex items-center'>{renderStatusNode(status)}</div>;
};
