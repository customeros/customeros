'use client';

import {Card, CardBody, CardHeader} from "@ui/layout/Card";
import {Text} from "@ui/typography/Text";
import React from "react";

export const BillingInfoPanel = () => {
  return (
    <>
        <Card>
            <CardHeader> <b>Billing Info</b> </CardHeader>
            <CardBody>
                <Text> You have <b>324353</b> Contacts.</Text>
                <br/>
                <Text>Your next invoice is for $0.</Text>
            </CardBody>
        </Card>
    </>
  );
};
