'use client';

import {Card, CardBody, CardHeader} from "@ui/layout/Card";
import {Button} from "@ui/form/Button";
import React from "react";
import {signIn, signOut} from "next-auth/react";
import {useJune} from "@spaces/hooks/useJune";
import {Divider} from "@ui/presentation/Divider";
import {Text} from "@ui/typography/Text";
import {Heading} from "@ui/typography/Heading";

const placeholders = {
  valueProposition: `Value proposition (A company's value prop is its raison d'être, its sweet spot, its jam. It's the special sauce that makes customers come back for more. It's the secret behind "Shut up and take my money!")`,
};

export const AuthPanel = () => {
    const analytics = useJune();

    const handleSyncGoogleClick = async () => {
        const res = await signIn(
            'google',
            { callbackUrl: '/settings?tab=oauth'},
            {prompt: "login", scope: "openid email profile https://www.googleapis.com/auth/gmail.readonly https://www.googleapis.com/auth/gmail.send https://www.googleapis.com/auth/gmail.compose https://www.googleapis.com/auth/calendar.events"});
    };

    const handleSignOutClick = () => {
        analytics?.reset();
        signOut();
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
                  <Heading as='h1' fontSize='lg' color='gray.700'>
                      Google OAuth
                      <br/>
                      <Divider></Divider>
                  </Heading>
              </CardHeader>

              <CardBody padding={6} pr={0} pt={0} position='unset'>
                  <br/>
                  <Text>Enable OAuth Integration to get access to your google workspace emails and calendar events</Text>
                  <br/>
                  <Button onClick={handleSyncGoogleClick}> Sync Google Account</Button>
                  <br/>
                  <br/>
                  <Button onClick={handleSignOutClick}>Sign Out</Button>
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
