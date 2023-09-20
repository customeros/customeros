import React from 'react';
import { CardHeader, CardBody } from '@ui/presentation/Card';
import { Heading } from '@ui/typography/Heading';
import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { Tooltip } from '@ui/presentation/Tooltip';
import { IconButton } from '@ui/form/IconButton';
import { useTimelineEventPreviewContext } from '../../preview/TimelineEventsPreviewContext/TimelineEventPreviewContext';
import CopyLink from '@spaces/atoms/icons/CopyLink';
import Times from '@spaces/atoms/icons/Times';
import copy from 'copy-to-clipboard';
import { VStack } from '@ui/layout/Stack';
import { LogEntryWithAliases } from '@organization/components/Timeline/types';
import { User } from '@graphql/types';
import { Box } from '@ui/layout/Box';
import Image from 'next/image';
import noteIcon from 'public/images/event-ill-log-preview.png';
import { LogEntryDatePicker } from './LogEntryDatePicker';

const getAuthor = (user: User) => {
  if (!user?.firstName && !user.lastName) {
    return 'Unknown';
  }

  return `${user.firstName} ${user.lastName}`.trim();
};

export const LogEntryPreviewModal: React.FC = () => {
  const { closeModal, modalContent } = useTimelineEventPreviewContext();
  const event = modalContent as LogEntryWithAliases;
  const author = getAuthor(event?.logEntryCreatedBy);

  return (
    <>
      <CardHeader pb={1} position='sticky' top={0} borderRadius='xl'>
        <Flex
          direction='row'
          justifyContent='space-between'
          alignItems='center'
        >
          <Flex mb={2} alignItems='center'>
            <Heading size='sm' fontSize='lg'>
              Log entry
            </Heading>
          </Flex>
          <Flex direction='row' justifyContent='flex-end' alignItems='center'>
            <Tooltip label='Copy link to this entry' placement='bottom'>
              <IconButton
                variant='ghost'
                aria-label='Copy link to this entry'
                color='gray.500'
                size='sm'
                mr={1}
                icon={<CopyLink color='gray.500' height='18px' />}
                onClick={() => copy(window.location.href)}
              />
            </Tooltip>
            <Tooltip label='Close' aria-label='close' placement='bottom'>
              <IconButton
                variant='ghost'
                aria-label='Close preview'
                color='gray.500'
                size='sm'
                icon={<Times color='gray.500' height='24px' />}
                onClick={closeModal}
              />
            </Tooltip>
          </Flex>
        </Flex>
      </CardHeader>
      <CardBody mt={0} maxHeight='50%' pb={6}>
        <VStack gap={2} alignItems='flex-start' position='relative'>
          <Box position='absolute' top={-2} right={-3}>
            <Image src={noteIcon} alt='' height={123} width={174} />
          </Box>
          <Flex direction='column'>
            <LogEntryDatePicker event={event} />
          </Flex>
          <Flex direction='column'>
            <Text size='sm' fontWeight='semibold'>
              Author
            </Text>
            <Text size='sm'>{author}</Text>
          </Flex>

          <Flex direction='column'>
            <Text size='sm' fontWeight='semibold'>
              Entry
            </Text>
            <Text
              className='slack-container'
              dangerouslySetInnerHTML={{ __html: `${event?.content}` }}
            />
          </Flex>

          <Text size='sm' fontWeight='medium'>
            {event.tags.map(({ name }) => `#${name}`).join(' ')}
          </Text>
        </VStack>
      </CardBody>
    </>
  );
};
