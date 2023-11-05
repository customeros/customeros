'use client';
import { Heading } from '@ui/typography/Heading';
import { Card, CardBody, CardHeader } from '@ui/presentation/Card';

export const MainSection = ({ children }: { children?: React.ReactNode }) => {
  return (
    <Card
      flex='3'
      h='100%'
      bg='#FCFCFC'
      borderRadius='2xl'
      flexDirection='column'
      overflow='hidden'
      boxShadow='none'
      position='relative'
      background='gray.25'
      minWidth={609}
      padding={0}
    >
      <CardHeader px={6} pb={2}>
        <Heading as='h1' fontSize='lg' color='gray.700'>
          Timeline
        </Heading>
      </CardHeader>
      <CardBody pr={0} pt={0} p={0} position='unset'>
        {children}
      </CardBody>
    </Card>
  );
};
