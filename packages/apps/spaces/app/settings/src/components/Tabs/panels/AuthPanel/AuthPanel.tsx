'use client';
import { useRouter, useSearchParams } from 'next/navigation';
import React, { useState, useEffect, ChangeEvent } from 'react';

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

import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { Switch } from '@ui/form/Switch';
import { Text } from '@ui/typography/Text';
import { Spinner } from '@ui/feedback/Spinner';
import { Heading } from '@ui/typography/Heading';
import { FormLabel } from '@ui/form/FormElement';
import { Outlook } from '@ui/media/logos/Outlook';
import { HStack, VStack } from '@ui/layout/Stack';
import { Divider } from '@ui/presentation/Divider';
import { Card, CardBody, CardHeader } from '@ui/layout/Card';
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

  const handleSyncGoogleToggle = async (event: ChangeEvent) => {
    setGoogleSettingsLoading(true);
    const scopes = [
      'openid',
      'email',
      'profile',
      'https://www.googleapis.com/auth/gmail.readonly',
      'https://www.googleapis.com/auth/gmail.send',
      'https://www.googleapis.com/auth/calendar.readonly',
    ];

    if ((event.target as HTMLInputElement).checked) {
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

  const handleSlackToggle = async (event: ChangeEvent) => {
    setSlackSettingsLoading(true);

    if ((event.target as HTMLInputElement).checked) {
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
      <Card
        bg='#FCFCFC'
        borderRadius='2xl'
        flexDirection='column'
        boxShadow='none'
        position='relative'
        background='gray.25'
        maxW='50%'
      >
        <CardHeader px={6} pb={2}>
          <Flex gap='1' align='center' mb='2'>
            <Icons.GOOGLE boxSize='6' />
            <Heading as='h1' fontSize='lg' color='gray.700'>
              Google OAuth
            </Heading>
          </Flex>
          <Divider></Divider>
        </CardHeader>

        <CardBody padding={6} pr={0} pt={0} position='unset'>
          <Text noOfLines={2} mt={2} mb={3}>
            Enable OAuth Integration to get access to your google workspace
            emails and calendar events
          </Text>
          <Flex direction={'column'} gap={2} width={'250px'}>
            <HStack>
              <VStack alignItems={'start'}>
                <Flex gap='1' align='center'>
                  <Icons.GMAIL boxSize='6' />
                  <FormLabel mb='0'>Sync Google Mail</FormLabel>
                </Flex>
                <Flex gap='1' align='center'>
                  <Icons.GOOGLE_CALENDAR boxSize='6' />
                  <FormLabel mb='0'>Sync Google Calendar</FormLabel>
                </Flex>
              </VStack>

              {googleSettingsLoading && <Spinner size='sm' color='green.500' />}
              {!googleSettingsLoading && (
                <Switch
                  isChecked={googleSettings.gmailSyncEnabled}
                  colorScheme='green'
                  onChange={(event) => handleSyncGoogleToggle(event)}
                />
              )}
            </HStack>
          </Flex>
        </CardBody>
      </Card>

      <Card
        bg='#FCFCFC'
        borderRadius='2xl'
        flexDirection='column'
        boxShadow='none'
        position='relative'
        background='gray.25'
        maxW='50%'
      >
        <CardHeader px={6} pb={2}>
          <Flex gap='1' align='center' mb='2'>
            <Outlook boxSize='6' />
            <Heading as='h1' fontSize='lg' color='gray.700'>
              Microsoft Outlook
            </Heading>
          </Flex>
          <Divider></Divider>
        </CardHeader>

        <CardBody padding={6} pr={0} pt={0} position='unset'>
          <Text noOfLines={2} mt={2} mb={3}>
            Enable OAuth Integration to get access to your microsoft outlook
            emails
          </Text>
          <Flex direction={'column'} gap={2} width={'250px'}>
            <HStack>
              <VStack alignItems={'start'}>
                <Flex gap='1' align='center'>
                  <Outlook boxSize='6' />
                  <FormLabel mb='0'>Sync Microsoft Outlook</FormLabel>
                </Flex>
              </VStack>

              {loading ? (
                <Spinner size='sm' color='green.500' />
              ) : (
                <Switch
                  colorScheme='green'
                  onChange={handleOutlookToggle}
                  defaultChecked={!!outlookConnection}
                />
              )}
            </HStack>
          </Flex>
        </CardBody>
      </Card>

      <Card
        bg='#FCFCFC'
        borderRadius='2xl'
        flexDirection='column'
        boxShadow='none'
        position='relative'
        background='gray.25'
        mt={4}
      >
        <CardHeader px={6} pb={2}>
          <Flex gap='1' align='center' mb='2'>
            <Icons.Slack boxSize='6' />
            <Heading as='h1' fontSize='lg' color='gray.700'>
              Slack
            </Heading>
          </Flex>
          <Divider></Divider>
        </CardHeader>

        <CardBody padding={6} pr={0} pt={0} position='unset'>
          <Text noOfLines={2} mt={2} mb={3}>
            Enable Slack Integration to get access to your Slack workspace
          </Text>
          <Flex direction={'column'} gap={2} width={'250px'}>
            <HStack>
              <VStack alignItems={'start'}>
                <Flex gap='1' align='center'>
                  <Icons.Slack boxSize='6' />
                  <FormLabel mb='0'>Sync Slack</FormLabel>
                </Flex>
              </VStack>

              {slackSettingsLoading && <Spinner size='sm' color='green.500' />}
              {!slackSettingsLoading && (
                <Switch
                  isChecked={slackSettings.slackEnabled}
                  colorScheme='green'
                  onChange={(event) => handleSlackToggle(event)}
                />
              )}
            </HStack>
          </Flex>
        </CardBody>
      </Card>
    </>
  );
};
