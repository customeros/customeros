import { signIn, useSession } from 'next-auth/react';

import { FeaturedIcon } from '@ui/media/Icon';
import { AlertCircle } from '@ui/media/icons/AlertCircle';
import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { useDisclosure } from '@ui/utils';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';
import { RefreshCcw01 } from '@ui/media/icons/RefreshCcw01';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import {
  GetOAuthUserSettings,
  OAuthUserSettingsInterface,
} from 'services/settings/settingsService';

export const GoogleSidebarNotification = () => {
  const client = getGraphQLClient();
  const { data: globalCacheQuery } = useGlobalCacheQuery(client);
  const { data: session } = useSession();

  const infoModal = useDisclosure();

  const requestAccess = () => {
    // @ts-expect-error look into it
    GetOAuthUserSettings(session.user.playerIdentityId).then(
      async (res: OAuthUserSettingsInterface) => {
        const scopes = ['openid', 'email', 'profile'];
        if (res.gmailSyncEnabled) {
          scopes.push(
            'https://www.googleapis.com/auth/gmail.readonly',
            'https://www.googleapis.com/auth/gmail.send',
          );
        }
        if (res.googleCalendarSyncEnabled) {
          scopes.push('https://www.googleapis.com/auth/calendar');
        }

        await signIn(
          'google',
          { callbackUrl: '/' },
          {
            prompt: 'login',
            scope: scopes.join(' '),
          },
        );
      },
    );
  };

  return (
    <>
      {globalCacheQuery?.global_Cache?.gmailOauthTokenNeedsManualRefresh && (
        <>
          <ConfirmDeleteDialog
            colorScheme={'purple'}
            icon={<RefreshCcw01 />}
            label='Re-allow access to Google'
            confirmButtonLabel='Re-allow'
            body={
              <>
                <Box fontSize={'14px'}>
                  Access to Google typically expires after a week or when you
                  change your password.
                </Box>
                <Box fontSize={'14px'} mt={2}>
                  To resume syncing your conversations and meetings, you need to
                  re-allow access to Google.
                </Box>
              </>
            }
            isOpen={infoModal.isOpen}
            onClose={infoModal.onClose}
            onConfirm={requestAccess}
          />

          <Box
            bg={'warning.25'}
            border={'solid 1px'}
            borderColor={'warning.200'}
            w={'full'}
            padding={'20px 16px'}
          >
            <FeaturedIcon size='md' minW='10' colorScheme={'orange'}>
              <AlertCircle />
            </FeaturedIcon>

            <Box color={'warning.700'} mt={2}>
              Your conversations and meetings are no longer syncing because
              access to your Google account has expired
            </Box>
            <Flex mt={2}>
              <Button
                variant={'ghost'}
                fontSize={'14px'}
                onClick={() => {
                  infoModal.onOpen();
                }}
              >
                Learn more
              </Button>
              <Button
                variant={'ghost'}
                fontSize={'14px'}
                color={'warning.700'}
                onClick={requestAccess}
              >
                Re-allow
              </Button>
            </Flex>
          </Box>
        </>
      )}
    </>
  );
};
