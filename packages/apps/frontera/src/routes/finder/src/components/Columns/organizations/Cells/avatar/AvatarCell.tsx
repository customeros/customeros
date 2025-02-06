import { memo, useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { useLocalStorage } from 'usehooks-ts';

import { cn } from '@ui/utils/cn';
import { Image } from '@ui/media/Image/Image';
import { Building06 } from '@ui/media/icons/Building06';

interface AvatarCellProps {
  id: string;
  name: string;
  icon?: string | null;
  logo?: string | null;
  isEnriching: boolean;
}

export const AvatarCell = memo(
  ({ name, id, icon, logo, isEnriching }: AvatarCellProps) => {
    const navigate = useNavigate();

    const [tabs] = useLocalStorage<{
      [key: string]: string;
    }>(`customeros-player-last-position`, { root: 'finder' });

    const src = icon || logo;
    const fullName = name || 'Unnamed';
    const [status, setStatus] = useState('loading');

    const handleNavigate = () => {
      const lastPositionParams = tabs[id];
      const href = getHref(id, lastPositionParams);

      navigate(href);
    };

    return (
      <div className='items-center ml-[1px]'>
        <div
          onClick={handleNavigate}
          className={cn(
            'w-6 h-6 flex items-center justify-center rounded border border-gray-200 cursor-pointer focus:outline-none',
            {
              'animate-pulse': isEnriching,
            },
          )}
        >
          {src && status !== 'error' && (
            <Image
              src={src}
              alt={fullName}
              loading='lazy'
              decoding='async'
              onError={() => setStatus('error')}
              onLoad={() => setStatus('loaded')}
              className={cn('w-full h-full object-contain', {
                'opacity-0 size-0': status === 'loading',
                'opacity-100': status === 'loaded',
              })}
            />
          )}

          {(!src || status === 'error' || status === 'loading') && (
            <Building06 className='w-4 h-4 text-gray-700' />
          )}
        </div>
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
