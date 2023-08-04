'use client';

import {Card, CardBody, CardHeader} from "@ui/layout/Card";
import {Button} from "@ui/form/Button";
import React from "react";
import {signIn, signOut} from "next-auth/react";
import {useJune} from "@spaces/hooks/useJune";

const placeholders = {
  valueProposition: `Value proposition (A company's value prop is its raison d'être, its sweet spot, its jam. It's the special sauce that makes customers come back for more. It's the secret behind "Shut up and take my money!")`,
};

export const OAuthPanel = () => {
    const analytics = useJune();

    const handleSyncGoogleClick = async () => {
        const res = await signIn(
            'google',
            { callbackUrl: '/tenant'},
            {prompt: "login", scope: "openid email profile https://www.googleapis.com/auth/gmail.readonly https://www.googleapis.com/auth/gmail.send https://www.googleapis.com/auth/gmail.compose https://www.googleapis.com/auth/calendar.events"});
    };

    const handleSignOutClick = () => {
        analytics?.reset();
        signOut();
    };

  return (
      <>
          <Card>
            <CardHeader> <b>OAuth Settings</b> </CardHeader>
            <CardBody>

              <Button onClick={handleSyncGoogleClick}> Sync Google Account</Button>

              <br/>
              <br/>
              <Button onClick={handleSignOutClick}>Sign Out</Button>
            </CardBody>
          </Card>
      </>
  );
};
