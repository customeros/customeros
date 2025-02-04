import { useState } from 'react';

import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';

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
    const [previewCard, setPreviewCard] = useLocalStorage('previewCard', false);
    const [status, setStatus] = useState('loading');

    const isEnriching = contactStore?.isEnriching;

    return (
      <div className='items-center ml-[1px]'>
        <div
          onClick={() => {
            if (previewCard === true && store.ui.focusRow === id) {
              setPreviewCard(false);
            } else {
              store.ui.setFocusRow(id);
              setPreviewCard(true);
            }
          }}
          className={cn(
            'w-6 h-6 flex items-center justify-center rounded border border-gray-200 cursor-pointer focus:outline-none',
            {
              'animate-pulse': isEnriching,
              'cursor-default': !canNavigate,
            },
          )}
        >
          {src && status !== 'error' && (
            <>
              <Image
                src={src}
                alt={fullName}
                loading='lazy'
                decoding='async'
                fetchPriority='low'
                onError={() => setStatus('error')}
                onLoad={() => setStatus('loaded')}
                className={cn('w-full h-full object-contain', {
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
