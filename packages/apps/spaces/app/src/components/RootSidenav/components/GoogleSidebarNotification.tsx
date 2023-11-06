import { signIn, useSession } from 'next-auth/react';
import {
  GetGoogleSettings,
  OAuthUserSettingsInterface,
} from 'services/settings/settingsService';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { useDisclosure } from '@ui/utils';
import { FeaturedIcon } from '@ui/media/Icon';
import { AlertCircle } from '@ui/media/icons/AlertCircle';
import { RefreshCcw01 } from '@ui/media/icons/RefreshCcw01';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';

export const GoogleSidebarNotification = () => {
  const client = getGraphQLClient();
  const { data: globalCacheQuery } = useGlobalCacheQuery(client);
  const { data: session } = useSession();

  const infoModal = useDisclosure();

  const requestAccess = () => {
    if (!session) return;
    GetGoogleSettings(session.user.playerIdentityId).then(
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
      {globalCacheQuery?.global_Cache?.isGoogleTokenExpired && (
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
            borderRadius='8px'
          >
            <FeaturedIcon size='md' minW='10' colorScheme={'orange'}>
              <AlertCircle />
            </FeaturedIcon>

            <Box color={'warning.700'} mt={2} fontSize='sm'>
              Your conversations and meetings are no longer syncing because
              access to your Google account has expired
            </Box>
            <Flex mt={2} justify='space-between'>
              <Button
                colorScheme='warning'
                variant={'link'}
                fontSize={'14px'}
                color='gray.500'
                fontWeight='semibold'
                _hover={{
                  color: 'gray.700',
                }}
                onClick={() => {
                  infoModal.onOpen();
                }}
              >
                Learn more
              </Button>
              <Button
                colorScheme='warning'
                variant={'link'}
                fontSize={'14px'}
                fontWeight='semibold'
                _hover={{
                  color: 'warning.900',
                }}
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
