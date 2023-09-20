import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Card, CardBody } from '@ui/presentation/Card';

import { useTimelineEventPreviewContext } from '../../preview/TimelineEventsPreviewContext/TimelineEventPreviewContext';
import React from 'react';
import Image from 'next/image';
import noteIcon from 'public/images/event-ill-log-stub.png';
import { Box } from '@spaces/ui/layout/Box';
import { Mail01 } from '@ui/media/icons/Mail01';
import { Calendar } from '@ui/media/icons/Calendar';
import { Phone } from '@ui/media/icons/Phone';
import { MessageTextSquare01 } from '@ui/media/icons/MessageTextSquare01';
import { convert } from 'html-to-text';
import { LogEntryWithAliases } from '@organization/components/Timeline/types';

interface LogEntryStubProps {
  data: LogEntryWithAliases;
}

function addHashToTagContent(htmlInput: string): string {
  const regex =
    /(<span[^>]*data-mention-atom-name="tag"[^>]*>)(.*?)(<\/span>)/g;
  return htmlInput.replace(
    regex,
    (match: string, p1: string, p2: string, p3: string) => {
      return p1 + '#' + p2 + p3;
    },
  );
}

const getLogEntryIcon = (type: string) => {
  switch (type) {
    case 'email':
      return <Mail01 color='gray.500' boxSize={3} />;
    case 'meeting':
      return <Calendar color='gray.500' boxSize={3} />;
    case 'voicemail':
    case 'call':
      return <Phone color='gray.500' boxSize={3} />;
    case 'text-message':
      return <MessageTextSquare01 color='gray.500' boxSize={3} />;

    default:
      return null;
  }
};

export const LogEntryStub = ({ data }: LogEntryStubProps) => {
  const { openModal } = useTimelineEventPreviewContext();
  const fullName = `${data.logEntryCreatedBy?.firstName} ${data.logEntryCreatedBy?.lastName}`;
  const text = convert(addHashToTagContent(data?.content || ''), {
    preserveNewlines: true,
    selectors: [
      {
        selector: 'a',
        options: { hideLinkHrefIfSameAsText: true, ignoreHref: true },
      },
    ],
  });
  const logEntryIcon = (() => {
    const icon = getLogEntryIcon(data.tags?.[0]?.name);

    if (!icon) return null;
    return (
      <Flex
        zIndex={1}
        position='relative'
        bg='white'
        border='1px solid'
        borderColor='gray.200'
        borderRadius='md'
        p={2}
        right='-3px'
        top='3px'
      >
        {icon}
      </Flex>
    );
  })();

  return (
    <Card
      variant='outline'
      size='md'
      maxWidth={549}
      cursor='pointer'
      boxShadow='xs'
      borderColor='gray.200'
      borderRadius='lg'
      onClick={() => openModal(data)}
      _hover={{ boxShadow: 'md' }}
      transition='all 0.2s ease-out'
    >
      <CardBody px='3' py='2'>
        <Flex w='full' justify='space-between' position='relative'>
          <Text w={390} noOfLines={4} color='gray.700' fontSize='sm'>
            <Text as='span'>{fullName}</Text>
            <Text as='span' color='gray.500' mx={1}>
              wrote
            </Text>
            <Text as='span'>{text}</Text>
          </Text>

          <Box h={94}>
            <Box position='absolute' top={-2} right={-3}>
              <Image src={noteIcon} alt='' height={94} width={124} />
            </Box>

            {logEntryIcon && logEntryIcon}
          </Box>
        </Flex>
      </CardBody>
    </Card>
  );
};
