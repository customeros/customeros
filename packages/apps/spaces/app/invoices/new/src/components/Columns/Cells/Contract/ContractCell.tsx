import Link from 'next/link';

export const ContractCell = ({
  value,
  organizationId,
}: {
  value: string;
  organizationId: string;
}) => {
  return (
    <Link
      href={`/organization/${organizationId}?tab=account`}
      className='font-medium line-clamp-1 text-gray-700 no-underline hover:no-underline hover:text-gray-900 transition-colors'
    >{`${value}'s contract`}</Link>
  );
};
