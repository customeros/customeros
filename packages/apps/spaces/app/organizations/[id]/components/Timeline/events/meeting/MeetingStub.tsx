import { convert } from 'html-to-text';

import { Meeting } from '@graphql/types';
import { Flex } from '@ui/layout/Flex';
import { Center } from '@ui/layout/Center';
import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { Icons } from '@ui/media/Icon';
import { Card, CardBody } from '@ui/presentation/Card';

import { getParticipants, getParticipantName } from '../utils';
import { useTimelineEventPreviewContext } from '../../preview/TimelineEventsPreviewContext/TimelineEventPreviewContext';
import { MeetingIcon } from './icons';

interface MeetingStubProps {
  data: Meeting;
}

export const MeetingStub = ({ data }: MeetingStubProps) => {
  const owner = getParticipantName(data.createdBy[0]);
  const firstParticipant = getParticipantName(data.attendedBy?.[0]);
  const [participants, remaining] = getParticipants(data);
  const { openModal } = useTimelineEventPreviewContext();

  const note = convert(data?.note?.[0]?.content ?? '', {
    preserveNewlines: true,
  });
  const agenda = convert(data?.agenda ?? '', { preserveNewlines: true });

  return (
    <Card
      variant='outline'
      size='md'
      maxWidth={549}
      position='unset'
      cursor='pointer'
      boxShadow='xs'
      borderColor='gray.200'
      borderRadius='lg'
      onClick={() => openModal(data)}
      _hover={{ boxShadow: 'md' }}
      transition='all 0.2s ease-out'
    >
      <CardBody px='3' py='2'>
        <Flex w='full' justify='space-between' position='relative' gap='3'>
          <VStack spacing='0' alignItems='flex-start'>
            <Text
              fontSize='sm'
              fontWeight='semibold'
              color='gray.700'
              noOfLines={1}
            >
              {data?.name ?? '(No title)'}
            </Text>
            <Flex>
              <Text
                fontSize='sm'
                color='gray.700'
                noOfLines={note || agenda ? 1 : 3}
                maxW='463px'
              >
                {owner || firstParticipant}{' '}
                <Text as='span' color='gray.500'>
                  met
                </Text>{' '}
                {participants}
              </Text>
              {remaining && (
                <Text
                  ml='1'
                  fontSize='sm'
                  as='span'
                  color='gray.500'
                  whiteSpace='nowrap'
                >
                  {` + ${remaining}`}
                </Text>
              )}
            </Flex>

            {(note || agenda) && (
              <Flex align='flex-start'>
                {note && (
                  <Icons.File2 boxSize='3' mt='1' mr='1' color='gray.500' />
                )}
                <Text fontSize='sm' color='gray.500' noOfLines={2}>
                  {note || agenda}
                </Text>
              </Flex>
            )}
          </VStack>

          <Center minW='12' h='10' fontSize='xxx-large'>
            <MeetingIcon />
            <Text
              position='absolute'
              fontSize='xl'
              fontWeight='semibold'
              mt='4px'
              color='gray.700'
            >
              {new Date(data?.startedAt).getDate()}
            </Text>
          </Center>
        </Flex>
      </CardBody>
    </Card>
  );
};
