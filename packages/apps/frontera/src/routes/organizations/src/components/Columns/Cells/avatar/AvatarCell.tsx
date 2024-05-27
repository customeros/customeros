import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { useLocalStorage } from 'usehooks-ts';

import { Image } from '@ui/media/Image/Image';
import { Avatar } from '@ui/media/Avatar/Avatar';
import {
  Popover,
  PopoverTrigger,
  PopoverContent,
} from '@ui/overlay/Popover/Popover';

interface AvatarCellProps {
  id: string;
  name: string;
  icon?: string | null;
  logo?: string | null;
  description?: string;
}

export const AvatarCell = ({
  name,
  id,
  icon,
  logo,
  description,
}: AvatarCellProps) => {
  const [isOpen, setIsOpen] = useState(false);
  const navigate = useNavigate();
  const [tabs] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-player-last-position`, { root: 'organization' });

  const src = icon || logo;
  const lastPositionParams = tabs[id];
  const href = getHref(id, lastPositionParams);
  const fullName = name || 'Unnamed';

  return (
    <div className='items-center'>
      <Popover open={isOpen} onOpenChange={setIsOpen}>
        <PopoverTrigger>
          <Avatar
            className='text-primary-700 cursor-pointer focus:outline-none'
            textSize='sm'
            variant='outlineSquare'
            size='xs'
            src={src || undefined}
            name={fullName}
            onMouseEnter={() => setIsOpen(true)}
            onMouseLeave={() => setIsOpen(false)}
            onClick={() => {
              navigate(href);
            }}
          />
        </PopoverTrigger>

        <PopoverContent className='w-[264px]'>
          {(logo || icon) && (
            <Image
              src={logo || icon || undefined}
              className='h-[36px] w-fit mb-1'
            />
          )}
          <p className='text-md font-semibold'>{fullName}</p>
          <p className='text-xs'>{description}</p>
        </PopoverContent>
      </Popover>
    </div>
  );
};

function getHref(id: string, lastPositionParams: string | undefined) {
  return `/organization/${id}?${lastPositionParams || 'tab=about'}`;
}
