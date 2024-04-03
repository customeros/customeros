import Link from 'next/link';

interface OrganizationCellProps {
  id: string;
  name: string;
  isSubsidiary: boolean;
  parentOrganizationName: string;
}

export const OrganizationCell = ({
  id,
  name,
  isSubsidiary,
  parentOrganizationName,
}: OrganizationCellProps) => {
  const href = `/organization/${id}?tab=account`;
  const fullName = name || 'Unnamed';

  return (
    <div className='flex flex-col line-clamp-1 '>
      {isSubsidiary && (
        <span className='text-xs text-gray-500'>{parentOrganizationName}</span>
      )}
      <Link
        className='text-gray-700 font-semibold hover:no-underline no-underline'
        href={href}
      >
        {fullName}
      </Link>
    </div>
  );
};
