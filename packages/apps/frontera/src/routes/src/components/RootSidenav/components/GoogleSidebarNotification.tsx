import { observer } from 'mobx-react-lite';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { AlertCircle } from '@ui/media/icons/AlertCircle';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { RefreshCcw01 } from '@ui/media/icons/RefreshCcw01';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog/ConfirmDeleteDialog';

export const GoogleSidebarNotification = observer(() => {
  const store = useStore();

  const infoModal = useDisclosure();

  const requestAccess = () => {
    store.settings.google.enableSync();
  };

  return (
    <>
      {store.globalCache.value?.isGoogleTokenExpired && (
        <>
          <ConfirmDeleteDialog
            icon={<RefreshCcw01 />}
            isOpen={infoModal.open}
            onConfirm={requestAccess}
            onClose={infoModal.onClose}
            confirmButtonLabel='Re-allow'
            label='Re-allow access to Google'
            isLoading={store.settings.google.isLoading}
            body={
              <>
                <div className='text-sm'>
                  Access to Google typically expires after a week or when you
                  change your password.
                </div>
                <div className='text-sm mt-2'>
                  To resume syncing your conversations and meetings, you need to
                  re-allow access to Google.
                </div>
              </>
            }
          />

          <div className='bg-warning-25 border border-warning-200 w-full py-5 px-4 rounded-lg'>
            <FeaturedIcon size='md' colorScheme='warning' className='ml-[2px]'>
              <AlertCircle />
            </FeaturedIcon>

            <div className='text-warning-700 mt-2 text-sm'>
              Your conversations and meetings are no longer syncing because
              access to your Google account has expired
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
                isLoading={store.settings.google.isLoading}
                className='text-sm font-semibold hover:text-warning-900 text-warning-700'
              >
                Re-allow
              </Button>
            </div>
          </div>
        </>
      )}
    </>
  );
});
