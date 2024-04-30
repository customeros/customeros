// import { signIn, useSession } from 'next-auth/react';
// import {
//   GetGoogleSettings,
//   OAuthUserSettingsInterface,
// } from 'src/services/settings/settingsService';

import { Button } from '@ui/form/Button/Button';
import { AlertCircle } from '@ui/media/icons/AlertCircle';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { RefreshCcw01 } from '@ui/media/icons/RefreshCcw01';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog/ConfirmDeleteDialog';

export const GoogleSidebarNotification = () => {
  const client = getGraphQLClient();
  const { data: globalCacheQuery } = useGlobalCacheQuery(client);
  // const { data: session } = useSession();

  const infoModal = useDisclosure();

  //TODO:SAME HERE NEED TO BE UPDATE THE SESION

  // const requestAccess = () => {
  //   if (!session) return;
  //   GetGoogleSettings(session.user.playerIdentityId).then(
  //     async (res: OAuthUserSettingsInterface) => {
  //       const scopes = [
  //         'openid',
  //         'email',
  //         'profile',
  //         'https://www.googleapis.com/auth/gmail.readonly',
  //         'https://www.googleapis.com/auth/gmail.send',
  //         'https://www.googleapis.com/auth/calendar.readonly',
  //       ];

  //       await signIn(
  //         'google',
  //         { callbackUrl: '/' },
  //         {
  //           prompt: 'consent',
  //           scope: scopes.join(' '),
  //         },
  //       );
  //     },
  //   );
  // };

  return (
    <>
      {globalCacheQuery?.global_Cache?.isGoogleTokenExpired && (
        <>
          <ConfirmDeleteDialog
            icon={<RefreshCcw01 />}
            label='Re-allow access to Google'
            confirmButtonLabel='Re-allow'
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
            isOpen={infoModal.open}
            onClose={infoModal.onClose}
            // onConfirm={requestAccess}
            onConfirm={() => {}}
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
                colorScheme='warning'
                variant='ghost'
                className='text-sm font-semibold hover:text-warning-900 text-warning-700'
                // onClick={requestAccess}
              >
                Re-allow
              </Button>
            </div>
          </div>
        </>
      )}
    </>
  );
};
