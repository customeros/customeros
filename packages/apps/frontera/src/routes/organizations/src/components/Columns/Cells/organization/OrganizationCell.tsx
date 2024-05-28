import { Link } from 'react-router-dom';
import { useRef, useState, useEffect } from 'react';

import { useLocalStorage } from 'usehooks-ts';

import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';

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
  const [isOverflowing, setIsOverflowing] = useState(false);

  const lastPositionParams = tabs[id];
  const href = getHref(id, lastPositionParams);
  const fullName = name || 'Unnamed';

  useEffect(() => {
    const element = linkRef.current;
    if (element) {
      const isOverflow = element.scrollWidth > element.clientWidth;
      setIsOverflowing(isOverflow);
    }
  }, [linkRef]);

  return (
    <Tooltip
      hasArrow
      align='start'
      side='bottom'
      label={isOverflowing ? fullName : ''}
    >
      <div className='flex flex-col line-clamp-1'>
        {isSubsidiary && (
          <span className='text-xs text-gray-500'>
            {parentOrganizationName}
          </span>
        )}
        <Link
          className='line-clamp-1 font-semibold text-gray-700 no-underline hover:no-underline'
          ref={linkRef}
          to={href}
        >
          {fullName}
        </Link>
      </div>
    </Tooltip>
  );
};

function getHref(id: string, lastPositionParams: string | undefined) {
  return `/organization/${id}?${lastPositionParams || 'tab=about'}`;
}
