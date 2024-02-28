'use client';
import { useParams } from 'next/navigation';

import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { HStack } from '@ui/layout/Stack';
import { Heading } from '@ui/typography/Heading';
import { useChannel } from '@shared/hooks/useChannel';
import { UserHexagon } from '@shared/components/UserHexagon';
import { Card, CardBody, CardHeader } from '@ui/presentation/Card';

export const MainSection = ({ children }: { children?: React.ReactNode }) => {
  const organizationId = useParams()?.id as string;
  const { presentUsers, username } = useChannel(
    `organization:${organizationId}`,
  );
  const isPresenceEnabled = useFeatureIsOn('presence');

  return (
    <Card
      flex='3'
      h='100%'
      bg='#FCFCFC'
      borderRadius='unset'
      flexDirection='column'
      overflow='hidden'
      boxShadow='none'
      position='relative'
      background='gray.25'
      minWidth={609}
      padding={0}
    >
      <CardHeader
        px={6}
        pb={2}
        display='flex'
        alignItems='center'
        flexDirection={'row'}
        justifyContent='space-between'
      >
        <Heading as='h1' fontSize='lg' color='gray.700'>
          Timeline
        </Heading>
        {isPresenceEnabled && (
          <HStack>
            {presentUsers.map(([user, color]) => (
              <UserHexagon
                key={user}
                name={user}
                color={color}
                isCurrent={user === username}
              />
            ))}
          </HStack>
        )}
      </CardHeader>
      <CardBody pr={0} pt={0} p={0} position='unset'>
        {children}
      </CardBody>
    </Card>
  );
};
