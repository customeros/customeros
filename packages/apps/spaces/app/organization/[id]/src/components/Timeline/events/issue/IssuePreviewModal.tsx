import { CardBody } from '@ui/presentation/Card';
import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import {
  useTimelineEventPreviewMethodsContext,
  useTimelineEventPreviewStateContext,
} from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';
import React from 'react';
import { Tag, TagLabel } from '@ui/presentation/Tag';
// import { IssueCommentCard } from '@organization/src/components/Timeline/events/issue/IssueCommentCard';
import { DateTimeUtils } from '@spaces/utils/date';
import {
  Priority,
  PriorityBadge,
} from '@organization/src/components/Timeline/events/issue/PriorityBadge';
import { Divider, HStack } from '@chakra-ui/react';
import { getExternalUrl } from '@spaces/utils/getExternalLink';
import { toastError } from '@ui/presentation/Toast';
import { MarkdownContentRenderer } from '@ui/presentation/MarkdownContentRenderer/MarkdownContentRenderer';
import { TimelineEventPreviewHeader } from '@organization/src/components/Timeline/preview/header/TimelineEventPreviewHeader';

function getStatusColor(status: string) {
  if (['closed', 'solved'].includes(status?.toLowerCase())) {
    return 'gray';
  }
  return 'blue';
}

export const IssuePreviewModal: React.FC = () => {
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();
  const issue = modalContent as any;
  const statusColorScheme = getStatusColor(issue.issueStatus);

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
      `${issue.id}-stub-open-in-external-app-error`,
    );
  };

  return (
    <>
      <TimelineEventPreviewHeader
        name={issue.subject ?? ''}
        onClose={closeModal}
        copyLabel='Copy link to this issue'
      />

      <CardBody
        mt={0}
        maxHeight='calc(100vh - 9rem)'
        p={6}
        pt={0}
        overflow='auto'
      >
        <HStack gap={2} mb={2} position='relative'>
          {issue?.priority && (
            <PriorityBadge
              priority={issue.priority.toLowerCase() as Priority}
            />
          )}
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
            <TagLabel textTransform='capitalize'>{issue.issueStatus}</TagLabel>
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
            cursor='pointer'
            onClick={handleOpenInExternalApp}
          >
            <TagLabel>#{issue?.externalLinks?.[0]?.externalId}</TagLabel>
          </Tag>
        </HStack>
        <Text fontSize='sm' mb={2}>
          <MarkdownContentRenderer markdownContent={issue?.description ?? ''} />
        </Text>

        {issue?.tags?.length && (
          <Text color='gray.500' fontSize='sm' mb={6}>
            {issue.tags.map((e: { name: string }) => e.name).join(' • ')}
          </Text>
        )}

        <Flex mb={2} alignItems='center'>
          <Text fontSize='sm' whiteSpace='nowrap'>
            Issue requested on
          </Text>
          {/*<Text mx={1} fontSize='sm' whiteSpace='nowrap'>*/}
          {/*  {issue?.requestedBy}*/}
          {/*</Text>*/}
          <Text
            color='gray.400'
            fontSize='sm'
            whiteSpace='nowrap'
            ml={1}
            mr={2}
          >
            {DateTimeUtils.format(issue?.createdAt, DateTimeUtils.dateWithHour)}
          </Text>
          <Divider orientation='horizontal' borderBottomColor='gray.200' />
        </Flex>
        {/* todo uncomment when data is available to query*/}
        {/*<VStack*/}
        {/*  gap={2}*/}
        {/*  w='full'*/}
        {/*  justifyContent='flex-start'*/}
        {/*  alignItems='flex-start'*/}
        {/*>*/}
        {/*  {Object.entries(xyz)?.map((values) => (*/}
        {/*    <React.Fragment key={values[0]}>*/}
        {/*      <Flex mb={2} alignItems='center' w='full'>*/}
        {/*        <Text color='gray.400' fontSize='sm' whiteSpace='nowrap' mr={2}>*/}
        {/*          {DateTimeUtils.format(values[0], DateTimeUtils.dateWithHour)}*/}
        {/*        </Text>*/}
        {/*        <Divider orientation='horizontal' />*/}
        {/*      </Flex>*/}
        {/*      {values[1]?.map((e) => (*/}
        {/*        <IssueCommentCard*/}
        {/*          key={e.id}*/}
        {/*          name={e.createdBy?.name}*/}
        {/*          content={e.content}*/}
        {/*          date={e.createdAt}*/}
        {/*        />*/}
        {/*      ))}*/}
        {/*    </React.Fragment>*/}
        {/*  ))}*/}
        {/*</VStack>*/}

        {['solved', 'closed'].includes(issue.issueStatus?.toLowerCase) && (
          <Text>
            Issue closed
            <Text color='gray.400' as='span' ml={1}>
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
