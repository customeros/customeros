'use client';

import {Card, CardBody, CardHeader} from '@ui/layout/Card';
import React, {ChangeEvent, useEffect, useState} from 'react';
import {signIn, useSession} from 'next-auth/react';
import {Divider} from '@ui/presentation/Divider';
import {Text} from '@ui/typography/Text';
import {Heading} from '@ui/typography/Heading';
import {Switch} from '@ui/form/Switch';
import {Flex, FormLabel, HStack, Spinner, VStack} from '@chakra-ui/react';
import {Icons} from '@ui/media/Icon';
import {
    GetOAuthUserSettings,
    OAuthUserSettingsInterface,
} from '../../../../../../../services/settings/settingsService';
import {GetServerSidePropsContext} from 'next';
import {getServerSession} from 'next-auth/next';
import {authOptions} from '../../../../../../../pages/api/auth/[...nextauth]';
import {RevokeAccess} from '../../../../../../../services/admin/userAdminService';
import {toastError, toastSuccess} from '@ui/presentation/Toast';

export const AuthPanel = () => {
    const {data: session} = useSession();

    const [googleToggleLoading, setGoogleToggleLoading] = useState(true);
    const [oAuthSettings, setOAuthSettings] =
        useState<OAuthUserSettingsInterface>({
            gmailSyncEnabled: false,
            googleCalendarSyncEnabled: false,
        });

    useEffect(() => {
        if (session) {
            setGoogleToggleLoading(true);
            // @ts-expect-error look into it
            GetOAuthUserSettings(session.user.playerIdentityId).then(
                (res: OAuthUserSettingsInterface) => {
                    setOAuthSettings(res);
                    setGoogleToggleLoading(false);
                },
            );
        }
    }, [session]);

    const handleSyncGoogleToggle = async (event: ChangeEvent) => {
        setGoogleToggleLoading(true);
        const scopes = [
            'openid',
            'email',
            'profile',
            'https://www.googleapis.com/auth/gmail.readonly',
            'https://www.googleapis.com/auth/gmail.send',
            'https://www.googleapis.com/auth/calendar.readonly',
        ];

        if ((event.target as HTMLInputElement).checked) {
            const res = await signIn(
                'google',
                {callbackUrl: '/settings?tab=oauth'},
                {
                    prompt: 'login',
                    scope: scopes.join(' '),
                },
            );
        } else {
            RevokeAccess({
                // @ts-expect-error look into it
                providerAccountId: session.user.playerIdentityId,
                provider: 'google',
            })
                .then((data: any) => {
                    // @ts-expect-error look into it
                    GetOAuthUserSettings(session.user.playerIdentityId).then(
                        (res: OAuthUserSettingsInterface) => {
                            setOAuthSettings(res);
                            setGoogleToggleLoading(false);
                        },
                    ).catch(() => {
                        setGoogleToggleLoading(false);
                        toastError(
                            'There was a problem on our side and we cannot load settings data at the moment, we are doing our best to solve it! ',
                            'revoke-google-access',
                        );
                    });
                    toastSuccess(
                        'We have successfully revoked access to your google account!',
                        'revoke-google-access',
                    );
                })
                .catch(() => {
                    setGoogleToggleLoading(false);
                    toastError(
                        'There was a problem on our side and we cannot load settings data at the moment, we are doing our best to solve it! ',
                        'revoke-google-access',
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
            >
                <CardHeader px={6} pb={2}>
                    <Flex gap='1' align='center' mb='2'>
                        <Icons.GOOGLE boxSize='6'/>
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
                                    <Icons.GMAIL boxSize='6'/>
                                    <FormLabel htmlFor={'changeEmailSyncSwitchButton'} mb='0'>
                                        Sync Google Mail
                                    </FormLabel>
                                </Flex>
                                <Flex gap='1' align='center'>
                                    <Icons.GOOGLE_CALENDAR boxSize='6'/>
                                    <FormLabel
                                        htmlFor={'changeGoogleCalendarSyncSwitchButton'}
                                        mb='0'
                                    >
                                        Sync Google Calendar
                                    </FormLabel>
                                </Flex>
                            </VStack>

                            {
                                googleToggleLoading &&
                                <Spinner size='sm' color='green.500' />
                            }
                            {!googleToggleLoading && (
                                <Switch
                                    id={'changeGmailSyncSwitchButton'}
                                    isChecked={oAuthSettings.gmailSyncEnabled}
                                    colorScheme='green'
                                    onChange={(event) => handleSyncGoogleToggle(event)}
                                />
                            )}
                        </HStack>
                    </Flex>
                </CardBody>
            </Card>
            <Divider/>

            <Card
                bg='#FCFCFC'
                borderRadius='2xl'
                flexDirection='column'
                boxShadow='none'
                position='relative'
                background='gray.25'
                w='full'
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
    return {props: {session}};
}
