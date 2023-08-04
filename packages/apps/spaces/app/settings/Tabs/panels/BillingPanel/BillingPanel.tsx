'use client';

import {Card, CardBody, CardHeader} from "@ui/layout/Card";
import {Text} from "@ui/typography/Text";
import React from "react";
import {Heading} from "@ui/typography/Heading";

export const BillingInfoPanel = () => {
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
            minWidth={609}
        >
            <CardHeader px={6} pb={2}>
                <Heading as='h1' fontSize='lg' color='gray.700'>
                    <b>Billing details</b>
                </Heading>
            </CardHeader>
            <CardBody>
                <Text>Other Authentication methods coming soon</Text>
            </CardBody>
        </Card>
    </>
  );
};
