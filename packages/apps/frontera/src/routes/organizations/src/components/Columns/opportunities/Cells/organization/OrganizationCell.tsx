import { useRef } from 'react';
import { useNavigate } from 'react-router-dom';

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
  const navigate = useNavigate();

  const linkRef = useRef<HTMLParagraphElement>(null);
  const fullName = name || 'Unnamed';

  const handleNavigate = () => {
    const lastPositionParams = tabs[id];
    const href = getHref(id, lastPositionParams);

    if (!href) return;

    navigate(href);
  };

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
        <p
          role='button'
          ref={linkRef}
          onClick={handleNavigate}
          data-test='organization-name-in-all-orgs-table'
          className='overflow-ellipsis overflow-hidden font-medium no-underline hover:no-underline cursor-pointer'
        >
          {fullName}
        </p>
      </span>
    </TableCellTooltip>
  );
};

function getHref(id: string, lastPositionParams: string | undefined) {
  return `/organization/${id}?${lastPositionParams || 'tab=about'}`;
}
