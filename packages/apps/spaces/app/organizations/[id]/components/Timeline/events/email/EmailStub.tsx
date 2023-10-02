'use client';
import React, { FC } from 'react';
import { Card, CardBody, CardFooter } from '@ui/presentation/Card';
import { Text } from '@ui/typography/Text';
import { VStack } from '@ui/layout/Stack';
import { convert } from 'html-to-text';
import { getEmailParticipantsName } from '@spaces/utils/getParticipantsName';
import { EmailParticipant } from '@graphql/types';
import { useTimelineEventPreviewContext } from '@organization/components/Timeline/preview/context/TimelineEventPreviewContext';
import { getEmailParticipantsByType } from '@organization/components/Timeline/events/email/utils';
import Image from 'next/image';
import { InteractionEventWithDate } from '../../types';

export const EmailStub: FC<{ email: InteractionEventWithDate }> = ({
  email,
}) => {
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
        boxShadow='xs'
        fontSize='14px'
        border='1px solid'
        borderColor='gray.200'
        background='white'
        flexDirection='row'
        maxWidth={549}
        position='unset'
        cursor='pointer'
        borderRadius='lg'
        onClick={() => openModal(email)}
        _hover={{ boxShadow: 'md' }}
        transition='all 0.2s ease-out'
      >
        <CardBody px='3' py='2' pr='0' overflow={'hidden'} flexDirection='row'>
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
        <CardFooter py='2' pr='3' pl='3' ml='1' display='block'>
          <Image
            src={'/backgrounds/organization/post-stamp.webp'}
            alt='Email'
            width={48}
            height={70}
          />
        </CardFooter>
      </Card>
    </>
  );
};
