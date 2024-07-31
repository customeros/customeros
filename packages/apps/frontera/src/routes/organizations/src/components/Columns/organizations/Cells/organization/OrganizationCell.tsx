import { useRef } from 'react';
import { Link } from 'react-router-dom';

import { useLocalStorage } from 'usehooks-ts';

import { TableCellTooltip } from '@ui/presentation/Table';

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
  const [tabs] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-player-last-position`, { root: 'organization' });
  const linkRef = useRef<HTMLAnchorElement>(null);

  const lastPositionParams = tabs[id];
  const href = getHref(id, lastPositionParams);
  const fullName = name || 'Unnamed';

  return (
    <TableCellTooltip
      hasArrow
      align='start'
      side='bottom'
      label={fullName}
      targetRef={linkRef}
    >
      <span className='inline'>
        {isSubsidiary && (
          <span className='text-xs text-gray-500'>
            {parentOrganizationName}
          </span>
        )}
        <Link
          to={href}
          ref={linkRef}
          data-test='organization-name-in-all-orgs-table'
          className='inline text-gray-700 no-underline hover:no-underline font-medium'
        >
          {fullName}
        </Link>
      </span>
    </TableCellTooltip>
  );
};

function getHref(id: string, lastPositionParams: string | undefined) {
  return `/organization/${id}?${lastPositionParams || 'tab=about'}`;
}
