import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Image } from '@ui/media/Image/Image';
import { User02 } from '@ui/media/icons/User02';
import { useStore } from '@shared/hooks/useStore';

interface AvatarCellProps {
  id: string;
  name: string;
  icon?: string | null;
  logo?: string | null;
  canNavigate?: boolean;
}

export const AvatarCell = observer(
  ({ name, id, icon, logo, canNavigate }: AvatarCellProps) => {
    const store = useStore();
    const src = icon || logo;
    const fullName = name || 'Unnamed';
    const contactStore = store.contacts.value.get(id);
    const [status, setStatus] = useState('loading');

    const isEnriching = contactStore?.isEnriching;

    return (
      <div className='items-center ml-[1px]'>
        <div
          className={cn(
            'w-6 h-6 flex items-center justify-center  border border-gray-200 cursor-pointer focus:outline-none rounded-full',
            {
              'animate-pulse': isEnriching,
              'cursor-default': !canNavigate,
            },
          )}
          onClick={() => {
            if (store.ui.showPreviewCard && store.ui.focusRow === id) {
              store.ui.setShowPreviewCard(false);
            } else {
              store.ui.setFocusRow(id);
              store.ui.setShowPreviewCard(true);
            }
          }}
        >
          {src && status !== 'error' && (
            <>
              <Image
                src={src}
                alt={fullName}
                loading='lazy'
                decoding='async'
                onError={() => setStatus('error')}
                onLoad={() => setStatus('loaded')}
                className={cn('w-full h-full object-contain rounded-full', {
                  'opacity-0 size-0': status === 'loading',
                  'opacity-100': status === 'loaded',
                })}
              />
            </>
          )}
          {(!src || status === 'error' || status === 'loading') && (
            <User02 className='w-4 h-4 text-gray-700' />
          )}
        </div>
      </div>
    );
  },
);
