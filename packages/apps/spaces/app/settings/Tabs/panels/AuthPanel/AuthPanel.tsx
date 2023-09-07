'use client';

import { Card, CardBody, CardHeader } from '@ui/layout/Card';
import React, { ChangeEvent, useEffect, useState } from 'react';
import { signIn } from 'next-auth/react';
import { Divider } from '@ui/presentation/Divider';
import { Text } from '@ui/typography/Text';
import { Heading } from '@ui/typography/Heading';
import { Switch } from '@ui/form/Switch';
import { Flex, FormLabel } from '@chakra-ui/react';
import { Icons } from '@spaces/ui/media/Icon';
import {
  GetOAuthUserSettings,
  OAuthUserSettingsInterface,
} from '../../../../../services/settings/settingsService';
import { useSession } from 'next-auth/react';
import { GetServerSidePropsContext } from 'next';
import { getServerSession } from 'next-auth/next';
import { authOptions } from '../../../../../pages/api/auth/[...nextauth]';

export const AuthPanel = () => {
  const { data: session } = useSession();
  const [oAuthSettings, setOAuthSettings] =
    useState<OAuthUserSettingsInterface>({
      gmailSyncEnabled: false,
      googleCalendarSyncEnabled: false,
    });

  useEffect(() => {
    if (session) {
      // @ts-expect-error look into it
      GetOAuthUserSettings(session.user.playerIdentityId).then(
        (res: OAuthUserSettingsInterface) => {
          setOAuthSettings(res);
        },
      );
    }
  }, [session]);

  const handleSyncGoogleMailClick = async (event: ChangeEvent) => {
    if ((event.target as HTMLInputElement).checked) {
      const res = await signIn(
        'google',
        { callbackUrl: '/settings?tab=oauth' },
        {
          prompt: 'login',
          scope:
            'openid email profile https://www.googleapis.com/auth/gmail.readonly https://www.googleapis.com/auth/gmail.send',
        },
      );
    } else {
      window.location.replace(
        process.env.CUSTOMER_OS_GOOGLE_MANAGE_ACCESS_URL as string,
      );
    }
  };

  const handleSyncGoogleCalendarClick = async (event: ChangeEvent) => {
    if ((event.target as HTMLInputElement).checked) {
      const res = await signIn(
        'google',
        { callbackUrl: '/settings?tab=oauth' },
        {
          prompt: 'login',
          scope:
            'openid email profile https://www.googleapis.com/auth/calendar.events',
        },
      );
    } else {
      window.location.replace(
        process.env.CUSTOMER_OS_GOOGLE_MANAGE_ACCESS_URL as string,
      );
    }
  };

  return (
    <>
      <Card
        flex='3'
        h='calc(100vh - 1rem)'
        bg='#FCFCFC'
        borderRadius='2xl'
        flexDirection='column'
        boxShadow='none'
        position='relative'
        background='gray.25'
        minWidth={700}
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
          <br />
          <Text>
            Enable OAuth Integration to get access to your google workspace
            emails and calendar events
          </Text>
          <br />
          <Flex direction={'column'} gap={2} width={'250px'}>
            <Flex justifyContent={'space-between'}>
              <Flex gap='1' align='center'>
                <Icons.GMAIL boxSize='6' />
                <FormLabel htmlFor={'changeEmailSyncSwitchButton'} mb='0'>
                  Sync Google Mail
                </FormLabel>
              </Flex>
              <Switch
                id={'changeGmailSyncSwitchButton'}
                isChecked={oAuthSettings.gmailSyncEnabled}
                colorScheme='green'
                onChange={(event) => handleSyncGoogleMailClick(event)}
              ></Switch>
            </Flex>
            <Flex justifyContent={'space-between'}>
              <Flex gap='1' align='center'>
                <Icons.GOOGLE_CALENDAR boxSize='6' />
                <FormLabel
                  htmlFor={'changeGoogleCalendarSyncSwitchButton'}
                  mb='0'
                >
                  Sync Google Calendar
                </FormLabel>
              </Flex>
              <Switch
                id={'changeGoogleCalendarSyncSwitchButton'}
                isChecked={oAuthSettings.googleCalendarSyncEnabled}
                colorScheme='green'
                onChange={(event) => handleSyncGoogleCalendarClick(event)}
              ></Switch>
            </Flex>
          </Flex>
        </CardBody>
      </Card>
      <Divider borderWidth={'2px'}></Divider>

      <Card
        flex='3'
        h='calc(100vh - 1rem)'
        bg='#FCFCFC'
        borderRadius='2xl'
        flexDirection='column'
        boxShadow='none'
        position='relative'
        background='gray.25'
        minWidth={609}
      >
        <CardHeader px={6} pb={2}>
          <Heading as='h1' fontSize='lg' color='gray.700'>
            <b>Other Auth</b>
          </Heading>
        </CardHeader>
        <CardBody>
          <Text>Other Authentication methods coming soon</Text>
        </CardBody>
      </Card>
    </>
  );
};

export async function getServerSideProps(context: GetServerSidePropsContext) {
  const session = await getServerSession(context.req, context.res, authOptions);
  return { props: { session } };
}
