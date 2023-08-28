import { convert } from 'html-to-text';

import { Meeting } from '@graphql/types';
import { Flex } from '@ui/layout/Flex';
import { Center } from '@ui/layout/Center';
import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { Icons } from '@ui/media/Icon';
import { Card, CardBody } from '@ui/presentation/Card';

import MeetingIcon from './meetingIcon.svg';

interface MeetingStubProps {
  data: Meeting;
}

export const MeetingStub = ({ data }: MeetingStubProps) => {
  const owner = (() => {
    if (data?.createdBy?.[0]?.__typename === 'ContactParticipant') {
      const participant = data?.createdBy?.[0]?.contactParticipant;
      return participant?.firstName ?? participant.emails?.[0]?.email;
    }
    return '';
  })();

  const firstParticipant = (() => {
    // use 1st participant as owner if there's no owner
    if (data?.attendedBy?.[0]?.__typename === 'ContactParticipant') {
      const participant = data?.attendedBy?.[0]?.contactParticipant;
      return participant?.firstName ?? participant.emails?.[0]?.email;
    }
    return '';
  })();

  const note = convert(data?.note?.[0]?.html);
  const agenda = data?.agenda ?? '';

  const [participants, remaining] = (() => {
    const count = data?.attendedBy?.length;

    if (data?.attendedBy?.length) {
      const fullArray = data?.attendedBy
        ?.map((participant) => {
          if (participant?.__typename === 'ContactParticipant') {
            return (
              participant?.contactParticipant?.firstName ??
              participant?.contactParticipant?.emails?.[0]?.email
            );
          }
          return '';
        })
        .filter(Boolean);

      if (!note || !agenda)
        return [fullArray.join(count > 2 ? ', ' : ' and '), ''];

      return fullArray
        .filter((v, i) => {
          if (!owner && i === 0 && v) return false; // skip 1st participant if there's no owner
          return v && i < (!owner ? 3 : 2);
        })
        .join(count > 2 ? ', ' : ' and ')
        .concat(count > 2 ? ` + ${count - (!owner ? 3 : 2)}` : '')
        .split(' + ');
    }

    return [];
  })();

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
    >
      <CardBody p='3'>
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

          <Center minW='12' h='10'>
            <MeetingIcon />
            <Text
              position='absolute'
              fontSize='xl'
              fontWeight='semibold'
              mt='4px'
              color='gray.700'
            >
              {new Date(data?.createdAt).getDate()}
            </Text>
          </Center>
        </Flex>
      </CardBody>
    </Card>
  );
};
