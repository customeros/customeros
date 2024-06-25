import { memo, useState } from 'react';
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
  variant?:
    | 'outlineSquare'
    | 'circle'
    | 'shadowed'
    | 'roundedSquareSmall'
    | 'roundedSquare'
    | 'roundedSquareShadowed'
    | 'outline'
    | 'outlineSquareSmall'
    | 'outlineCircle'
    | null
    | undefined;
}

export const AvatarCell = memo(
  ({
    name,
    id,
    icon,
    logo,
    description,
    variant = 'outlineSquare',
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
      <div className='items-center ml-[1px]'>
        <Popover open={isOpen} onOpenChange={setIsOpen}>
          <PopoverTrigger>
            <Avatar
              className='text-gray-700 cursor-pointer focus:outline-none'
              textSize='xs'
              variant={variant}
              tabIndex={-1}
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

          <PopoverContent
            className='w-[264px]'
            onCloseAutoFocus={(e) => e.preventDefault()}
          >
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
  },
  (prevProps, nextProps) => {
    return (
      prevProps.icon === nextProps.icon && prevProps.logo === nextProps.logo
    );
  },
);

function getHref(id: string, lastPositionParams: string | undefined) {
  return `/organization/${id}?${lastPositionParams || 'tab=about'}`;
}
