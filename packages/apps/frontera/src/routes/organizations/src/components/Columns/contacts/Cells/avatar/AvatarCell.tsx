import { memo } from 'react';
import { useNavigate } from 'react-router-dom';

import { useLocalStorage } from 'usehooks-ts';

import { Avatar } from '@ui/media/Avatar/Avatar';

interface AvatarCellProps {
  id: string;
  name: string;
  icon?: string | null;
  logo?: string | null;
}

export const AvatarCell = memo(
  ({ name, id, icon, logo }: AvatarCellProps) => {
    const navigate = useNavigate();
    const [tabs] = useLocalStorage<{
      [key: string]: string;
    }>(`customeros-player-last-position`, { root: 'organization' });

    const src = icon || logo;
    const lastPositionParams = tabs[id];
    const href = getHref(id, 'tab=people' || lastPositionParams);
    const fullName = name || 'Unnamed';

    return (
      <div className='items-center ml-[1px]'>
        <Avatar
          className='text-gray-700 cursor-pointer focus:outline-none'
          textSize='xs'
          variant='outlineCircle'
          tabIndex={-1}
          size='xs'
          src={src || undefined}
          name={fullName}
          onClick={() => {
            navigate(href);
          }}
        />
      </div>
    );
  },
  (prevProps, nextProps) => {
    return (
      prevProps.icon === nextProps.icon && prevProps.logo === nextProps.logo
    );
  },
);

function getHref(id: string, lastPositionParams: string | undefined) {
  return `/organization/${id}?${lastPositionParams || 'tab=people'}`;
}