'use client';
import React, { useRef, useMemo } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Avatar } from '@ui/media/Avatar';
import { Text } from '@ui/typography/Text';
import { User01 } from '@ui/media/icons/User01';
import { Issue, Contact } from '@graphql/types';
import { Heading } from '@ui/typography/Heading';
import { DateTimeUtils } from '@spaces/utils/date';
import { Tag, TagLabel } from '@ui/presentation/Tag';
import { Card, CardHeader } from '@ui/presentation/Card';
import {
  getParticipant,
  getParticipantName,
} from '@organization/src/hooks/utils';
import { useTimelineEventPreviewMethodsContext } from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

interface IssueCardProps {
  issue: Issue;
}
function getStatusColor(status: string) {
  if (['closed', 'solved'].includes(status.toLowerCase())) {
    return 'gray';
  }

  return 'blue';
}
export const IssueCard = ({ issue }: IssueCardProps) => {
  const cardRef = useRef<HTMLDivElement>(null);
  const { openModal } = useTimelineEventPreviewMethodsContext();
  const statusColorScheme = (() => getStatusColor(issue.status))();
  const isStatusClosed = useMemo(
    () => ['closed', 'solved'].includes(issue.status.toLowerCase()),
    [issue.status],
  );

  const submittedBy = useMemo(
    () =>
      issue.submittedBy ? getParticipantName(issue.submittedBy) : undefined,
    [issue.id],
  );
  const reportedBy = useMemo(
    () => (issue.reportedBy ? getParticipantName(issue.reportedBy) : undefined),

    [issue.id],
  );
  const profilePhoto = useMemo(
    () =>
      issue.reportedBy
        ? getParticipant(issue.reportedBy)
        : issue?.submittedBy
        ? getParticipant(issue.submittedBy)
        : undefined,

    [issue.id],
  );

  const participantName = useMemo(
    () => (
      <Text display='inline'>
        {reportedBy ? `Reported` : `Submitted`}

        {(reportedBy || submittedBy) && (
          <Text as='span' mx={1}>
            by
          </Text>
        )}
        <Text fontWeight='bold' as='span' mr={1}>
          {issue?.reportedBy ? reportedBy : submittedBy}
        </Text>
      </Text>
    ),
    [reportedBy, submittedBy],
  );

  const titleWidth = useMemo(() => {
    if (isStatusClosed) {
      return 'auto';
    }

    return issue?.status === 'pending' ? 250 : 260;
  }, [isStatusClosed, issue?.status]);

  return (
    <Card
      key={issue.id}
      w='full'
      ref={cardRef}
      boxShadow={'xs'}
      size='sm'
      cursor='pointer'
      borderRadius='lg'
      border='1px solid'
      borderColor='gray.200'
      onClick={() => openModal(issue.id)}
      _hover={{
        boxShadow: 'md',
        '& > div > #confirm-button': {
          opacity: '1',
          pointerEvents: 'auto',
        },
      }}
      transition='all 0.2s ease-out'
    >
      <CardHeader>
        <Flex flex='1' gap='4' alignItems='flex-start' flexWrap='wrap'>
          <Avatar
            size='md'
            name={submittedBy ?? reportedBy}
            variant='outlined'
            src={
              (profilePhoto as unknown as Contact)?.profilePhotoUrl ?? undefined
            }
            border={'1px solid var(--chakra-colors-primary-200)'}
            icon={<User01 color='primary.700' height='1.8rem' />}
          />

          <Flex direction='column' flex={1}>
            <Heading
              mt={1}
              size='sm'
              fontSize='sm'
              noOfLines={1}
              maxW={titleWidth}
            >
              {issue?.subject ?? '[No subject]'}
            </Heading>

            <Text fontSize='sm' mt={1} mb='2px' lineHeight={1}>
              {participantName}
              {DateTimeUtils.timeAgo(issue?.createdAt, { addSuffix: true })}
            </Text>

            {!!issue?.updatedAt && (
              <Text fontSize='sm' color='gray.500' lineHeight={1}>
                Last response was{' '}
                {DateTimeUtils.timeAgo(issue.updatedAt, {
                  addSuffix: true,
                })}
              </Text>
            )}
          </Flex>

          {!isStatusClosed && (
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
              position='absolute'
              right={3}
            >
              <TagLabel textTransform='capitalize'>{issue.status}</TagLabel>
            </Tag>
          )}
        </Flex>
      </CardHeader>
    </Card>
  );
};
