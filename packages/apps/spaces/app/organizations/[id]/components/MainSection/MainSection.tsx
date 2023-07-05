'use client';

import { CardHeader, Card, CardBody } from '@ui/layout/Card';
import { Divider, Heading } from '@chakra-ui/react';

export const MainSection = ({ children }: { children?: React.ReactNode }) => {
  return (
    <Card
      flex='3'
      h='calc(100vh - 2rem)'
      bg='#FCFCFC'
      borderRadius='2xl'
      shadow='base'
      flexDirection='column'
      position='relative'
      maxWidth={700}
      minWidth={609}
    >
      <CardHeader pr={6} pl={6} pt={3} pb={2}>
        <Heading as='h1' fontSize='2xl'>
          Timeline
        </Heading>
      </CardHeader>
      <Divider color='#EAECF0' />
      <CardBody padding={6} pr={0} pt={0} position='unset'>
        {children}
      </CardBody>
    </Card>
  );
};
