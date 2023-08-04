'use client';
import React, { FC } from 'react';
import { Card, CardBody, CardFooter } from '@ui/presentation/Card';
import { Text } from '@ui/typography/Text';
import { VStack } from '@ui/layout/Stack';
import { convert } from 'html-to-text';
import { getEmailParticipantsName } from '@spaces/utils/getParticipantsName';
import { EmailParticipant, InteractionEvent } from '@graphql/types';
import { useTimelineEventPreviewContext } from '@organization/components/Timeline/preview/TimelineEventsPreviewContext/TimelineEventPreviewContext';
import { getEmailParticipantsByType } from '@organization/components/Timeline/events/email/utils';
import Image from 'next/image';

export const EmailStub: FC<{ email: InteractionEvent }> = ({ email }) => {
  const { openModal } = useTimelineEventPreviewContext();
  const text = convert(email?.content || '', {
    preserveNewlines: true,
    selectors: [
      {
        selector: 'a',
        options: { hideLinkHrefIfSameAsText: true, ignoreHref: true },
      },
    ],
  });

  const { to, cc } = getEmailParticipantsByType(email?.sentTo || []);

  return (
    <>
      <Card
        variant='outline'
        size='md'
        fontSize='14px'
        background='white'
        flexDirection='row'
        maxWidth={549}
        position='unset'
        aspectRatio='9/2'
        cursor='pointer'
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
                {getEmailParticipantsName(
                  ([email?.sentBy?.[0]] as unknown as EmailParticipant[]) || [],
                )}
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

            <Text fontWeight='semibold' noOfLines={1}>
              {email.interactionSession?.name}
            </Text>

            <Text noOfLines={2} wordBreak='break-word'>
              {text}
            </Text>
          </VStack>
        </CardBody>
        <CardFooter pt={5} pb={5} pr={5} pl={0} ml={1}>
          <div>
            <Image
              src={'/backgrounds/organization/post-stamp.webp'}
              alt='Email'
              width={54}
              height={70}
              style={{ filter: 'box-shadow(0px 0.5px 1px #D8D8D8)' }}
            />
          </div>
        </CardFooter>
      </Card>
    </>
  );
};
