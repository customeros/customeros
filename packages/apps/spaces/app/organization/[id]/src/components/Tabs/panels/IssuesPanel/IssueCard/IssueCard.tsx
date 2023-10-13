'use client';
import React, { useRef } from 'react';
import { Flex } from '@ui/layout/Flex';
import { Avatar } from '@ui/media/Avatar';
import { Text } from '@ui/typography/Text';
import { Card, CardHeader } from '@ui/presentation/Card';
import { User01 } from '@ui/media/icons/User01';
import { Heading } from '@ui/typography/Heading';
import { Tag, TagLabel } from '@ui/presentation/Tag';
import { DateTimeUtils } from '@spaces/utils/date';
import { getExternalUrl } from '@spaces/utils/getExternalLink';
import { toastError } from '@ui/presentation/Toast';
import { Issue } from '@graphql/types';
// import { getContactDisplayName } from '@spaces/utils/getContactName';
// import { useContactOrUserDisplayName } from '@shared/hooks/useContactOrUserDisplayData';

// TODO uncomment commented out code as soon as COS-464 is merged
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
  // const getDisplayName = useContactOrUserDisplayName();
  // const requestorName = getDisplayName(issue.requestedBy);
  const statusColorScheme = (() => getStatusColor(issue.status))();

  // const getLastCreatedNote = (notes: Array<Note>) => {
  //   const sortedNotes = notes.sort((a, b) => b.createdAt - a.createdAt);
  //   return sortedNotes[0].createdAt;
  // };

  const handleOpenInExternalApp = () => {
    if (issue?.externalLinks?.[0]?.externalUrl) {
      // replacing this https://gasposhelp.zendesk.com/api/v2/tickets/24.json -> https://gasposhelp.zendesk.com/agent/tickets/24
      const replacedUrl = issue?.externalLinks?.[0]?.externalUrl
        .replace('api/v2', 'agent')
        .replace('.json', '');

      window.open(getExternalUrl(replacedUrl), '_blank', 'noreferrer noopener');
      return;
    }
    toastError(
      'This issue is not connected to external source',
      `${issue.id}-tab-panel-open-in-external-app-error`,
    );
  };

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
      onClick={handleOpenInExternalApp}
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
            name={''}
            variant='outlined'
            src={undefined} // todo
            border={'1px solid var(--chakra-colors-primary-200)'}
            icon={<User01 color='primary.700' height='1.8rem' />}
          />

          <Flex direction='column' flex={1}>
            <Heading size='sm' fontSize='sm'>
              {issue?.subject ?? '[No subject]'}
            </Heading>

            <Text fontSize='sm'>
              Requested{' '}
              {DateTimeUtils.timeAgo(issue?.createdAt, { addSuffix: true })}
              {/* by <Text as='span' fontWeight='medium' color='gray.700' mx={1}>*/}
              {/*  {requestorName}*/}
              {/*</Text>*/}
            </Text>

            {/*{!!issue?.notes?.length && (*/}
            {/*  <Text fontSize='sm' color='gray.500'>*/}
            {/*    Last response was{' '}*/}
            {/*    {DateTimeUtils.timeAgo(getLastCreatedNote(issue.notes), {*/}
            {/*      addSuffix: true,*/}
            {/*    })}*/}
            {/*  </Text>*/}
            {/*)}*/}
          </Flex>

          {['closed', 'solved'].includes(issue.status.toLowerCase()) && (
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
              right={2}
            >
              <TagLabel>{issue.status}</TagLabel>
            </Tag>
          )}
        </Flex>
      </CardHeader>
    </Card>
  );
};
