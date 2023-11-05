'use client';
import Image from 'next/image';
import React, { FC } from 'react';

import { convert } from 'html-to-text';

import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { EmailParticipant } from '@graphql/types';
import { Card, CardBody, CardFooter } from '@ui/presentation/Card';
import { getEmailParticipantsName } from '@spaces/utils/getParticipantsName';
import { getEmailParticipantsByType } from '@organization/src/components/Timeline/events/email/utils';
import { useTimelineEventPreviewMethodsContext } from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

import { InteractionEventWithDate } from '../../types';

export const EmailStub: FC<{ email: InteractionEventWithDate }> = ({
  email,
}) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();
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
  const isSendByTenant = (email?.sentBy?.[0] as EmailParticipant)
    ?.emailParticipant?.users?.length;

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
        onClick={() => openModal(email.id)}
        _hover={{ boxShadow: 'md' }}
        transition='all 0.2s ease-out'
        ml={isSendByTenant ? 6 : 0}
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
