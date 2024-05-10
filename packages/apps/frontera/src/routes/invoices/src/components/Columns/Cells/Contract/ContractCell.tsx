import { Link } from 'react-router-dom';

export const ContractCell = ({
  value,
  organizationId,
  organizationName,
}: {
  value: string;
  organizationId: string;
  organizationName: string;
}) => {
  return (
    <div>
      <p className='text-xs text-gray-500'>{organizationName}</p>
      <Link
        to={`/organization/${organizationId}?tab=account`}
        className='font-medium line-clamp-1 text-gray-700 no-underline hover:no-underline hover:text-gray-900 transition-colors'
      >
        {value}
      </Link>
    </div>
  );
};
