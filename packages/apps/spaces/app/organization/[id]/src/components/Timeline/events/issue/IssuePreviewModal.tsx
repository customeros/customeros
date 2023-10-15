import { CardHeader, CardBody } from '@ui/presentation/Card';
import { Heading } from '@ui/typography/Heading';
import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { Tooltip } from '@ui/presentation/Tooltip';
import { IconButton } from '@ui/form/IconButton';
import { useTimelineEventPreviewContext } from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';
import { Link03 } from '@ui/media/icons/Link03';
import { XClose } from '@ui/media/icons/XClose';
import copy from 'copy-to-clipboard';
import React from 'react';
import { Tag, TagLabel } from '@ui/presentation/Tag';
import { IssueCommentCard } from '@organization/src/components/Timeline/events/issue/IssueCommentCard';
import { DateTimeUtils } from '@spaces/utils/date';
import { PriorityBadge } from '@organization/src/components/Timeline/events/issue/PriorityBadge';
import { Divider, HStack, VStack } from '@chakra-ui/react';

function getStatusColor(status: string) {
  if (['closed', 'solved'].includes(status.toLowerCase())) {
    return 'gray';
  }
  return 'blue';
}
const mentionedByNotes = [
  {
    id: '1',
    content: 'Note 1',
    createdAt: '2023-05-01T00:00:00Z',
    contentType: 'text',
    createdBy: {
      id: 'user1',
      name: 'User One',
      firstName: 'User',
      lastName: 'One',
    },
  },
  {
    id: '2',
    content: 'Note 2',
    createdAt: '2023-05-01T02:00:00Z',
    contentType: 'text',
    createdBy: {
      id: 'user2',
      name: 'User Two',
      firstName: 'User',
      lastName: 'Two',
    },
  },
  {
    id: '3',
    content: 'Note 3',
    createdAt: '2023-05-02T00:00:00Z',
    contentType: 'text',
    createdBy: {
      id: 'user3',
      name: 'User Three',
      firstName: 'User',
      lastName: 'Three',
    },
  },
];

function restructureData(data) {
  return data.reduce((acc, note) => {
    const date = new Date(note.createdAt).toDateString();
    if (!acc[date]) {
      acc[date] = [];
    }
    acc[date].push(note);
    return acc;
  }, {});
}

const xyz = restructureData(mentionedByNotes);
export const IssuePreviewModal: React.FC = () => {
  const { closeModal, modalContent } = useTimelineEventPreviewContext();
  const issue = modalContent as any;
  const statusColorScheme = getStatusColor(issue.issueStatus);
  console.log('🏷️ -----  issue: ', issue);
  return (
    <>
      <CardHeader
        py='4'
        px='6'
        pb='1'
        position='sticky'
        top={0}
        borderRadius='xl'
      >
        <Flex
          direction='row'
          justifyContent='space-between'
          alignItems='center'
        >
          <Flex alignItems='center'>
            <Heading size='sm' fontSize='lg'>
              {issue?.subject ?? 'Issue'}
            </Heading>
          </Flex>
          <Flex direction='row' justifyContent='flex-end' alignItems='center'>
            <Tooltip label='Copy link' placement='bottom'>
              <IconButton
                variant='ghost'
                aria-label='Copy link to this issue'
                color='gray.500'
                fontSize='sm'
                size='sm'
                mr={1}
                icon={<Link03 color='gray.500' boxSize='4' />}
                onClick={() => copy(window.location.href)}
              />
            </Tooltip>
            <Tooltip label='Close' aria-label='close' placement='bottom'>
              <IconButton
                variant='ghost'
                aria-label='Close preview'
                color='gray.500'
                fontSize='sm'
                size='sm'
                icon={<XClose color='gray.500' boxSize='5' />}
                onClick={closeModal}
              />
            </Tooltip>
          </Flex>
        </Flex>
      </CardHeader>
      <CardBody
        mt={0}
        maxHeight='calc(100vh - 9rem)'
        p={6}
        pt={0}
        overflow='auto'
      >
        <HStack gap={2} mb={2}>
          <PriorityBadge priority={issue?.priority?.toLowerCase()} />

          <Tag
            size='sm'
            variant='outline'
            colorScheme='blue'
            border='1px solid'
            background='white'
            borderColor={`${[statusColorScheme]}.200`}
            backgroundColor={`${[statusColorScheme]}.50`}
            color={`${[statusColorScheme]}.700`}
            boxShadow='none'
            fontWeight='normal'
            minHeight={6}
          >
            <TagLabel>{issue.issueStatus}</TagLabel>
          </Tag>
          <Tag
            size='sm'
            variant='outline'
            colorScheme='blue'
            border='1px solid'
            background='white'
            borderColor={`gray.200`}
            backgroundColor={`white`}
            color={`gray.500`}
            boxShadow='none'
            fontWeight='normal'
            minHeight={6}
          >
            <TagLabel>{issue?.externalLinks?.[0]?.externalId}</TagLabel>
          </Tag>
        </HStack>

        <Text fontSize='sm' mb={2}>
          {issue?.description}
        </Text>

        {issue?.tags?.length && (
          <Text color='gray.500' fontSize='sm' mb={6}>
            {issue.tags.map((e: any) => e.name).join(' • ')}
          </Text>
        )}

        <Flex mb={2} alignItems='center'>
          <Text fontSize='sm' whiteSpace='nowrap'>
            Issue requested by
          </Text>
          <Text mx={1} fontSize='sm' whiteSpace='nowrap'>
            {issue?.requestedBy}
          </Text>
          <Text color='gray.400' fontSize='sm' whiteSpace='nowrap' mr={2}>
            {DateTimeUtils.format(issue?.createdAt, DateTimeUtils.dateWithHour)}
          </Text>
          <Divider orientation='horizontal' />
        </Flex>

        <VStack
          gap={2}
          w='full'
          justifyContent='flex-start'
          alignItems='flex-start'
        >
          {Object.entries(xyz)?.map((values) => (
            <React.Fragment key={values[0]}>
              <Flex mb={2} alignItems='center' w='full'>
                <Text color='gray.400' fontSize='sm' whiteSpace='nowrap' mr={2}>
                  {DateTimeUtils.format(values[0], DateTimeUtils.dateWithHour)}
                </Text>
                <Divider orientation='horizontal' />
              </Flex>

              {values[1]?.map((e) => (
                <IssueCommentCard
                  key={e.id}
                  name={e.createdBy?.name}
                  content={e.content}
                  date={e.createdAt}
                />
              ))}
            </React.Fragment>
          ))}
        </VStack>

        {['solved', 'closed'].includes(issue.issueStatus?.toLowerCase) && (
          <Text>
            Issue closed
            {/*by(todo user data)*/}
            <Text color='gray.400' as='span'>
              {DateTimeUtils.format(
                issue.updatedAt,
                DateTimeUtils.dateWithHour,
              )}
            </Text>
          </Text>
        )}
      </CardBody>
    </>
  );
};
