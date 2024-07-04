import { useRef } from 'react';
import { Link } from 'react-router-dom';

import { useLocalStorage } from 'usehooks-ts';

import { TableCellTooltip } from '@ui/presentation/Table';

interface OrganizationCellProps {
  id: string;
  name: string;
}

export const OrganizationCell = ({ id, name }: OrganizationCellProps) => {
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
        <Link
          className='inline text-gray-700 no-underline hover:no-underline font-normal'
          ref={linkRef}
          to={href}
        >
          {fullName}
        </Link>
      </span>
    </TableCellTooltip>
  );
};

function getHref(id: string, lastPositionParams: string | undefined) {
  return `/organization/${id}?${lastPositionParams || 'tab=people'}`;
}
