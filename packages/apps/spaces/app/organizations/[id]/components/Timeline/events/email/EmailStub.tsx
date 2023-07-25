'use client';
import React, { FC } from 'react';
import Image from 'next/image';
import { Card, CardBody, CardFooter } from '@ui/presentation/Card';
import { Text } from '@ui/typography/Text';
import { VStack } from '@chakra-ui/react';
import { convert } from 'html-to-text';
import { getEmailParticipantsName } from '@spaces/utils/getParticipantsName';
import { EmailParticipant, InteractionEvent } from '@graphql/types';
import { useTimelineEventPreviewContext } from '@organization/components/Timeline/preview/TimelineEventsPreviewContext/TimelineEventPreviewContext';
import { getEmailParticipantsByType } from '@organization/components/Timeline/events/email/utils';

export const EmailStub: FC<{ email: InteractionEvent }> = ({ email }) => {
  const { openModal } = useTimelineEventPreviewContext();
  const text = convert(email?.content || '', {
    preserveNewlines: true,
  });
  const { to, cc } = getEmailParticipantsByType(email?.sentTo || []);

  return (
    <>
      <Card
        variant='outline'
        size='md'
        fontSize='14px'
        background='gray.50'
        flexDirection='row'
        maxWidth={549}
        position='unset'
        aspectRatio='9/2'
        onClick={() => openModal(email)}
      >
        <CardBody
          pt={5}
          pb={5}
          pl={5}
          pr={0}
          overflow={'hidden'}
          flexDirection='row'
        >
          <VStack align='flex-start' spacing={0}>
            <Text as='p' noOfLines={1}>
              <Text as={'span'} fontWeight={500}>
                {getEmailParticipantsName(email?.sentBy as unknown as EmailParticipant[] || [])}
              </Text>{' '}
              <Text as={'span'} color='#6C757D'>
                emailed
              </Text>{' '}
              <Text as={'span'} fontWeight={500} marginRight={2}>
                {getEmailParticipantsName(to)}
              </Text>{' '}
              {!!cc.length && (
                <>
                  <Text as={'span'} color='#6C757D'>
                    CC:
                  </Text>{' '}
                  <Text as={'span'}>{getEmailParticipantsName(cc)}</Text>
                </>
              )}
            </Text>

            <Text fontWeight={500} noOfLines={1}>
              {email.interactionSession?.name}
            </Text>

            <Text noOfLines={2} wordBreak='break-word'>
              {text}
            </Text>
          </VStack>
        </CardBody>
        <CardFooter pt={5} pb={5} pr={5} pl={0}>
          <div>
            <Image
              src={'/backgrounds/organization/poststamp1.webp'}
              alt='Email'
              width={54}
              height={70}
            />
          </div>
        </CardFooter>
      </Card>
    </>
  );
};
