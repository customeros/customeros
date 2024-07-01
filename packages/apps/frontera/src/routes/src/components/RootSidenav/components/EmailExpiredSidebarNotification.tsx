import React from 'react';
import { useNavigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { XClose } from '@ui/media/icons/XClose';
import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { AlertCircle } from '@ui/media/icons/AlertCircle';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { RefreshCcw01 } from '@ui/media/icons/RefreshCcw01';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog/ConfirmDeleteDialog';

export const EmailExpiredSidebarNotification = observer(() => {
  const store = useStore();
  const navigate = useNavigate();

  const [isOpen, setIsOpen] = React.useState(true);
  const infoModal = useDisclosure();

  const requestAccess = () => {
    navigate(`/settings`);
  };

  return (
    <>
      {store.globalCache.value?.inactiveEmailTokens &&
        store.globalCache.value?.inactiveEmailTokens.length > 0 && (
          <>
            <ConfirmDeleteDialog
              icon={<RefreshCcw01 />}
              isOpen={infoModal.open}
              onConfirm={requestAccess}
              onClose={infoModal.onClose}
              confirmButtonLabel='Re-allow'
              hideCloseButton={false}
              label='Re-allow access to emails'
              isLoading={store.settings.oauthToken.isLoading}
              body={
                <>
                  <div className='text-sm'>
                    Access to Google and Microsoft typically expires after a
                    week or when you change your password.
                  </div>
                  <div className='text-sm mt-2'>
                    To resume syncing your conversations and meetings, you need
                    to re-allow access to the expired emails.
                  </div>
                </>
              }
            />

            {isOpen && (
              <div className='bg-warning-25 border border-warning-200 w-full py-5 px-4 rounded-lg'>
                <div className={'flex flex-row justify-between items-center'}>
                  <FeaturedIcon
                    size='md'
                    colorScheme='warning'
                    className='ml-[2px]'
                  >
                    <AlertCircle />
                  </FeaturedIcon>

                  <IconButton
                    variant='ghost'
                    colorScheme='gray'
                    icon={<XClose />}
                    aria-label='Close dialog'
                    onClick={() => {
                      setIsOpen(false);
                    }}
                  />
                </div>

                <div className='text-warning-700 mt-2 text-sm'>
                  Your conversations and meetings are no longer syncing because
                  access to some of your email accounts has expired
                </div>
                <div className='flex mt-2 justify-center'>
                  <Button
                    colorScheme='gray'
                    variant='ghost'
                    className='font-semibold hover:gray-700 text-sm'
                    onClick={() => {
                      infoModal.onOpen();
                    }}
                  >
                    Learn more
                  </Button>
                  <Button
                    variant='ghost'
                    colorScheme='warning'
                    onClick={requestAccess}
                    isLoading={store.settings.oauthToken.isLoading}
                    className='text-sm font-semibold hover:text-warning-900 text-warning-700'
                  >
                    Re-allow
                  </Button>
                </div>
              </div>
            )}
          </>
        )}
    </>
  );
});
