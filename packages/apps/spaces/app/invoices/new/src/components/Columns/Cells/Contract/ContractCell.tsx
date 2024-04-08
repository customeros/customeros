import Link from 'next/link';

import { Eye } from '@ui/media/icons/Eye';

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
      <div className='flex items-center gap-1'>
        <Link
          href={`/organization/${organizationId}?tab=account`}
          className='font-medium line-clamp-1 text-gray-700 no-underline hover:no-underline hover:text-gray-900 transition-colors contract-cell peer'
        >
          {value}
        </Link>
        <Eye
          boxSize='4'
          className='opacity-0 peer-hover:opacity-100 transition-opacity text-gray-400'
        />
      </div>
    </div>
  );
};
