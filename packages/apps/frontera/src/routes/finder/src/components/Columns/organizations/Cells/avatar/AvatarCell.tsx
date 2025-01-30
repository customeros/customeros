import { memo, useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { useLocalStorage } from 'usehooks-ts';

import { cn } from '@ui/utils/cn';
import { Image } from '@ui/media/Image/Image';
import { Building06 } from '@ui/media/icons/Building06';
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
  isEnriching: boolean;
}

export const AvatarCell = memo(
  ({ name, id, icon, logo, description, isEnriching }: AvatarCellProps) => {
    const [isOpen, setIsOpen] = useState(false);
    const navigate = useNavigate();

    const [tabs] = useLocalStorage<{
      [key: string]: string;
    }>(`customeros-player-last-position`, { root: 'organization' });

    const src = icon || logo;
    const fullName = name || 'Unnamed';

    const handleNavigate = () => {
      const lastPositionParams = tabs[id];
      const href = getHref(id, lastPositionParams);

      navigate(href);
    };

    return (
      <div className='items-center ml-[1px]'>
        <Popover open={isOpen} onOpenChange={setIsOpen}>
          <PopoverTrigger>
            <div
              onClick={handleNavigate}
              onMouseEnter={() => setIsOpen(true)}
              onMouseLeave={() => setIsOpen(false)}
              className={cn(
                'w-6 h-6 flex items-center justify-center rounded border border-gray-200 cursor-pointer focus:outline-none',
                {
                  'animate-pulse': isEnriching,
                },
              )}
            >
              {src ? (
                <Image
                  src={src}
                  alt={fullName}
                  loading='lazy'
                  decoding='async'
                  fetchPriority='low'
                  className='w-full h-full object-contain'
                />
              ) : (
                <Building06 className='w-4 h-4 text-gray-700' />
              )}
            </div>
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
            <p className='text-xs line-clamp-[8]'>{description}</p>
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
