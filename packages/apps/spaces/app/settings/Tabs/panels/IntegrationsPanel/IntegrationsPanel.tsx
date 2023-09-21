'use client';
import { useIntegrationApp } from '@integration-app/react';

import { Button } from '@ui/form/Button';
import { Heading } from '@ui/typography/Heading';
import { Card, CardBody, CardHeader } from '@ui/layout/Card';

export const IntegrationsPanel = () => {
  const iApp = useIntegrationApp();

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
        maxHeight='calc(100vh - 1rem)'
      >
        <CardHeader px={6} pb={2}>
          <Heading as='h1' fontSize='lg' color='gray.700'>
            <b>Data Integrations</b>
          </Heading>
        </CardHeader>
        <CardBody overflow='auto'>
          <Button onClick={() => iApp.open()}>Choose Integrations</Button>
        </CardBody>
      </Card>
    </>
  );
};
