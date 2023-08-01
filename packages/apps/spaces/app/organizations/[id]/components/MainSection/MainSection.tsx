'use client';
import { CardHeader, Card, CardBody } from '@ui/presentation/Card';
import { Heading } from '@ui/typography/Heading';

export const MainSection = ({ children }: { children?: React.ReactNode }) => {
  return (
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
          Timeline
        </Heading>
      </CardHeader>
      <CardBody padding={6} pr={0} pt={0} position='unset'>
        {children}
      </CardBody>
    </Card>
  );
};
