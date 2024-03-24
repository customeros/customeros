'use client';
import React, { useState, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

import axios from 'axios';
import { signIn, useSession } from 'next-auth/react';
import { useQueryClient } from '@tanstack/react-query';
import { RevokeAccess } from 'services/admin/userAdminService';
import {
  useConnections,
  useIntegrations,
  useIntegrationApp,
} from '@integration-app/react';
import {
  GetSlackSettings,
  GetGoogleSettings,
  SlackSettingsInterface,
  OAuthUserSettingsInterface,
} from 'services/settings/settingsService';

import { Icons } from '@ui/media/Icon';
import { Spinner } from '@ui/feedback/Spinner';
import { Switch } from '@ui/form/Switch/Switch2';
import { FormLabel } from '@ui/form/FormElement';
import { Outlook } from '@ui/media/logos/Outlook';
import { toastError, toastSuccess } from '@ui/presentation/Toast';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';

export const AuthPanel = () => {
  const iApp = useIntegrationApp();
  const { items: iIntegrations } = useIntegrations();
  const { items: iConnections, refresh, loading } = useConnections();
  const router = useRouter();
  const { data: session } = useSession();
  const queryClient = useQueryClient();
  const queryParams = useSearchParams();

  const outlookConnection = iConnections.find(
    (o) => o?.integration?.key === 'microsoft-outlook',
  );

  const handleOutlookToggle = async () => {
    const outlookIntegration = iIntegrations.find(
      (o) => o.key === 'microsoft-outlook',
    );

    if (!outlookIntegration) {
      toastError(
        'Microsoft Outlook integration not available',
        'get-intergration-data',
      );

      return;
    }

    try {
      await iApp
        .integration(outlookIntegration.key)
        .open({ showPoweredBy: false });
      await refresh();
    } catch (err) {
      toastError('Integration failed', 'get-intergration-data');
    }
  };

  useEffect(() => {
    if (
      queryParams &&
      queryParams.has('redirect_slack') &&
      queryParams.has('code')
    ) {
      setSlackSettingsLoading(true);

      axios
        .post(`/ua/slack/oauth/callback?code=${queryParams.get('code')}`)
        .then(({ data }) => {
          GetSlackSettings().then((res: SlackSettingsInterface) => {
            setSlackSettings(res);
            setSlackSettingsLoading(false);
          });
          router.push('/settings?tab=auth');
        })
        .catch((reason) => {
          router.push('/settings?tab=auth');
        });
    } else {
      setSlackSettingsLoading(true);
      GetSlackSettings().then((res: SlackSettingsInterface) => {
        setSlackSettings(res);
        setSlackSettingsLoading(false);
      });
    }
  }, [queryParams]);

  const [googleSettingsLoading, setGoogleSettingsLoading] = useState(true);
  const [googleSettings, setGoogleSettings] =
    useState<OAuthUserSettingsInterface>({
      gmailSyncEnabled: false,
      googleCalendarSyncEnabled: false,
    });

  const [slackSettingsLoading, setSlackSettingsLoading] = useState(true);
  const [slackSettings, setSlackSettings] = useState<SlackSettingsInterface>({
    slackEnabled: false,
  });

  useEffect(() => {
    if (session) {
      setGoogleSettingsLoading(true);
      GetGoogleSettings(session.user.playerIdentityId).then(
        (res: OAuthUserSettingsInterface) => {
          setGoogleSettings(res);
          setGoogleSettingsLoading(false);
        },
      );
    }
  }, [session]);

  const handleSyncGoogleToggle = async (isChecked: boolean) => {
    setGoogleSettingsLoading(true);
    const scopes = [
      'openid',
      'email',
      'profile',
      'https://www.googleapis.com/auth/gmail.readonly',
      'https://www.googleapis.com/auth/gmail.send',
      'https://www.googleapis.com/auth/calendar.readonly',
    ];

    if (isChecked) {
      const _ = await signIn(
        'google',
        { callbackUrl: '/settings?tab=oauth' },
        {
          prompt: 'consent',
          scope: scopes.join(' '),
        },
      );
    } else {
      RevokeAccess('google', {
        // @ts-expect-error look into it
        providerAccountId: session.user.playerIdentityId,
      })
        .then((data) => {
          // @ts-expect-error look into it
          GetGoogleSettings(session.user.playerIdentityId)
            .then((res: OAuthUserSettingsInterface) => {
              setGoogleSettings(res);
              setGoogleSettingsLoading(false);
              queryClient.invalidateQueries({
                queryKey: useGlobalCacheQuery.getKey(),
              });
            })
            .catch(() => {
              setGoogleSettingsLoading(false);
              toastError(
                'There was a problem on our side and we cannot load settings data at the moment, we are doing our best to solve it! ',
                'revoke-google-access',
              );
            });
          toastSuccess(
            'We have successfully revoked the access to your google account!',
            'revoke-google-access',
          );
        })
        .catch(() => {
          setGoogleSettingsLoading(false);
          toastError(
            'There was a problem on our side and we cannot load settings data at the moment, we are doing our best to solve it! ',
            'revoke-google-access',
          );
        });
    }
  };

  const handleSlackToggle = async (isChecked: boolean) => {
    setSlackSettingsLoading(true);

    if (isChecked) {
      axios
        .get(`/ua/slack/requestAccess`)
        .then(({ data }) => {
          location.href = data.url;
        })
        .catch((reason) => {
          toastError(
            'There was a problem on our side and we cannot load settings data at the moment, we are doing our best to solve it! ',
            'request-access-slack-access',
          );
          setSlackSettingsLoading(false);
        });
    } else {
      RevokeAccess('slack')
        .then((data) => {
          GetSlackSettings()
            .then((res: SlackSettingsInterface) => {
              setSlackSettings(res);
              setSlackSettingsLoading(false);
            })
            .catch(() => {
              setSlackSettingsLoading(false);
              toastError(
                'There was a problem on our side and we cannot load settings data at the moment, we are doing our best to solve it! ',
                'revoke-slack-access',
              );
            });
          toastSuccess(
            `We can't access your Slack workspace anymore`,
            'revoke-slack-access',
          );
        })
        .catch(() => {
          setSlackSettingsLoading(false);
          toastError(
            'There was a problem on our side and we cannot load settings data at the moment, we are doing our best to solve it! ',
            'revoke-slack-access',
          );
        });
    }
  };

  return (
    <>
      <div className='bg-gray-25 rounded-2xl flex-col flex relative max-w-[50%] '>
        <div className='px-6 pb-2'>
          <div className='flex gap-1 items-center mb-2 pt-5 '>
            <Icons.GOOGLE boxSize='6' />
            <h1 className='text-gray-700 text-lg '>Google OAuth</h1>
          </div>
          <div className='w-full border-b border-gray-100' />
        </div>

        <div className='p-6 pr-0 pt-0 '>
          <p className='line-clamp-2 mt-2 mb-3'>
            Enable OAuth Integration to get access to your google workspace
            emails and calendar events
          </p>

          <div className='flex flex-col gap-2 w-[250px]'>
            <div className='flex gap-2 items-center'>
              <div className='flex flex-col items-start gap-4'>
                <div className='flex gap-1 items-center'>
                  <Icons.GMAIL boxSize='6' />
                  <FormLabel mb='0'>Sync Google Mail</FormLabel>
                </div>

                <div className='flex gap-1 items-center'>
                  <Icons.GOOGLE_CALENDAR boxSize='6' />
                  <FormLabel mb='0'>Sync Google Calendar</FormLabel>
                </div>
              </div>

              {googleSettingsLoading && <Spinner size='sm' color='green.500' />}
              {!googleSettingsLoading && (
                <Switch
                  isChecked={googleSettings.gmailSyncEnabled}
                  onChange={(value) => handleSyncGoogleToggle(value)}
                  colorScheme='success'
                />
              )}
            </div>
          </div>
        </div>
      </div>

      <div className='bg-gray-25 rounded-2xl flex-col flex relative max-w-[50%] '>
        <div className='px-6 pb-2'>
          <div className='flex gap-1 items-center mb-2 pt-5 '>
            <Outlook boxSize='6' />
            <h1 className='text-gray-700 text-lg'>Microsoft Outlook</h1>
          </div>
          <div className='w-full border-b border-gray-100' />
        </div>

        <div className='p-6 pr-0 pt-0 '>
          <p className='line-clamp-2 mt-2 mb-3'>
            Enable OAuth Integration to get access to your microsoft outlook
            emails
          </p>

          <div className='flex space-x-4 items-center'>
            <div className='flex alig-middle space-x-1'>
              <Outlook boxSize='6' />
              <FormLabel mb='0'>Sync Microsoft Outlook</FormLabel>
            </div>
            {loading ? (
              <Spinner size='sm' color='green.500' />
            ) : (
              <Switch
                colorScheme='success'
                onChange={handleOutlookToggle}
                isChecked={!!outlookConnection}
              />
            )}
          </div>
        </div>
      </div>

      <div className='bg-gray-25 rounded-2xl flex-col mt-4 flex relative max-w-[50%] '>
        <div className='px-6 pb-2'>
          <div className='flex items-center gap-1 mb-2'>
            <Icons.Slack boxSize='6' />
            <h1 className='text-gray-700 text-lg'>Slack</h1>
          </div>
          <div className='w-full border-b border-gray-100' />
        </div>

        <div className='p-6 pr-0 pt-0'>
          <p className='line-clamp-2 mt-2 mb-3'>
            Enable Slack Integration to get access to your Slack workspace
          </p>

          <div className='flex space-x-4 items-center'>
            <div className='flex alig-middle space-x-1'>
              <Icons.Slack boxSize='6' />
              <FormLabel mb='0'>Sync Slack</FormLabel>
            </div>
            {slackSettingsLoading && <Spinner size='sm' color='green.500' />}
            {!slackSettingsLoading && (
              <Switch
                isChecked={slackSettings.slackEnabled}
                colorScheme='success'
                onChange={(isChecked) => handleSlackToggle(isChecked)}
              />
            )}
          </div>
        </div>
      </div>
    </>
  );
};
