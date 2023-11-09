import { FC, Fragment } from 'react';

import { match } from 'ts-pattern';

import { Flex } from '@ui/layout/Flex';
import { Link } from '@ui/navigation/Link';
import { Text } from '@ui/typography/Text';
import { Zendesk } from '@ui/media/logos/Zendesk';
import { VStack, HStack } from '@ui/layout/Stack';
import { Divider } from '@ui/presentation/Divider';
import { DateTimeUtils } from '@spaces/utils/date';
import { Tag, TagLabel } from '@ui/presentation/Tag';
import { Comment, InteractionEvent } from '@graphql/types';
import { getName } from '@spaces/utils/getParticipantsName';
import { CardBody, CardFooter } from '@ui/presentation/Card';
import { getExternalUrl } from '@spaces/utils/getExternalLink';
import { IssueWithAliases } from '@organization/src/components/Timeline/types';
import { IssueCommentCard } from '@organization/src/components/Timeline/events/issue/IssueCommentCard';
import { MarkdownContentRenderer } from '@ui/presentation/MarkdownContentRenderer/MarkdownContentRenderer';
import {
  Priority,
  PriorityBadge,
} from '@organization/src/components/Timeline/events/issue/PriorityBadge';
import { TimelineEventPreviewHeader } from '@organization/src/components/Timeline/preview/header/TimelineEventPreviewHeader';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

function getStatusColor(status: string) {
  if (['closed', 'solved'].includes(status?.toLowerCase())) {
    return 'gray';
  }

  return 'blue';
}

export const IssuePreviewModal: FC = () => {
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();
  const issue = modalContent as IssueWithAliases;
  const statusColorScheme = getStatusColor(issue.issueStatus);

  const reportedBy = match(issue.reportedBy)
    .with({ __typename: 'ContactParticipant' }, ({ contactParticipant }) =>
      getName(contactParticipant),
    )
    .with(
      { __typename: 'OrganizationParticipant' },
      ({ organizationParticipant }) => getName(organizationParticipant),
    )
    .with({ __typename: 'UserParticipant' }, ({ userParticipant }) =>
      getName(userParticipant),
    )
    .otherwise(() => '');

  const submittedBy = match(issue.submittedBy)
    .with({ __typename: 'ContactParticipant' }, ({ contactParticipant }) =>
      getName(contactParticipant),
    )
    .with(
      { __typename: 'OrganizationParticipant' },
      ({ organizationParticipant }) => getName(organizationParticipant),
    )
    .with({ __typename: 'UserParticipant' }, ({ userParticipant }) =>
      getName(userParticipant),
    )
    .otherwise(() => '');

  const commentsByDay = [...issue.comments, ...issue.interactionEvents]
    .sort((a, b) =>
      new Date(a.createdAt).valueOf() > new Date(b.createdAt).valueOf()
        ? 1
        : -1,
    )
    .reduce((acc, curr) => {
      const day = curr.createdAt.split('T')[0];

      if (acc[day]) {
        acc[day].push(curr);
      } else {
        acc[day] = [curr];
      }

      return acc;
    }, {} as Record<string, (Comment | InteractionEvent)[]>);

  const externalUrl = (() => {
    const url = issue?.externalLinks?.[0]?.externalUrl;
    if (!url) return null;

    return url.replace('api/v2', 'agent').replace('.json', '');
  })();

  return (
    <>
      <TimelineEventPreviewHeader
        name={issue.subject ?? ''}
        onClose={closeModal}
        copyLabel='Copy link'
      />

      <CardBody
        mt={0}
        maxH='calc(100vh - 4rem - 56px - 51px - 16px - 8px);'
        p={6}
        pt={0}
        overflow='auto'
      >
        <HStack
          gap={2}
          mb={2}
          position='sticky'
          top='0'
          bg='white'
          pb='2'
          zIndex='10'
        >
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
          >
            <TagLabel>#{issue?.externalLinks?.[0]?.externalId}</TagLabel>
          </Tag>
        </HStack>
        <Text fontSize='sm' mb={2}>
          <MarkdownContentRenderer markdownContent={issue?.description ?? ''} />
        </Text>

        {issue?.tags?.length && (
          <Text color='gray.500' fontSize='sm' mb={6}>
            {issue.tags.map((t) => t?.name).join(' • ')}
          </Text>
        )}

        <Flex align='center' mb='2'>
          <Flex alignItems='baseline'>
            <Text fontSize='sm' whiteSpace='nowrap'>
              {`Issue ${submittedBy ? 'submitted' : 'reported'} by`}
            </Text>
            <Text mx={1} fontSize='sm' whiteSpace='nowrap' fontWeight='medium'>
              {submittedBy || reportedBy}
            </Text>
            <Text
              color='gray.400'
              fontSize='sm'
              whiteSpace='nowrap'
              ml={1}
              mr={2}
            >
              {DateTimeUtils.format(
                issue?.createdAt,
                DateTimeUtils.dateWithHour,
              )}
            </Text>
          </Flex>
          <Divider
            orientation='horizontal'
            borderBottomColor='gray.200'
            h='full'
          />
        </Flex>

        <VStack
          gap={2}
          w='full'
          justifyContent='flex-start'
          alignItems='flex-start'
        >
          {Object.entries(commentsByDay)?.map(([date, comments]) => (
            <Fragment key={date}>
              {!DateTimeUtils.isSameDay(issue?.createdAt, date) && (
                <Flex alignItems='center' w='full'>
                  <Text
                    color='gray.400'
                    fontSize='sm'
                    whiteSpace='nowrap'
                    mr={2}
                  >
                    {DateTimeUtils.format(date, DateTimeUtils.date)}
                  </Text>
                  <Divider orientation='horizontal' />
                </Flex>
              )}

              {comments?.map((c) => {
                const name = match(c)
                  .with({ __typename: 'Comment' }, (c) => {
                    return !c.createdBy ? '' : getName(c.createdBy);
                  })
                  .with({ __typename: 'InteractionEvent' }, (c) => {
                    return match(c.sentBy?.[0])
                      .with(
                        { __typename: 'UserParticipant' },
                        ({ userParticipant }) => {
                          return getName(userParticipant);
                        },
                      )
                      .otherwise(() => '');
                  })
                  .otherwise(() => '');

                const isCustomer = match(c)
                  .with({ __typename: 'InteractionEvent' }, (e) =>
                    match(e.sentBy?.[0])
                      .with(
                        { __typename: 'OrganizationParticipant' },
                        () => true,
                      )
                      .otherwise(() => false),
                  )
                  .otherwise(() => false);

                return (
                  <IssueCommentCard
                    key={c.id}
                    name={name}
                    date={c.createdAt}
                    content={c.content ?? ''}
                    iscustomer={isCustomer}
                    isPrivate={c.__typename === 'Comment'}
                  />
                );
              })}
            </Fragment>
          ))}
        </VStack>

        {['solved', 'closed'].includes(issue.issueStatus?.toLowerCase()) && (
          <Flex align='center' mt='2'>
            <Flex alignItems='baseline'>
              <Text fontSize='sm' whiteSpace='nowrap'>
                Issue closed
              </Text>

              <Text
                color='gray.400'
                fontSize='sm'
                whiteSpace='nowrap'
                ml={2}
                mr={2}
              >
                {DateTimeUtils.format(
                  issue?.updatedAt,
                  DateTimeUtils.dateWithHour,
                )}
              </Text>
            </Flex>
            <Divider
              orientation='horizontal'
              borderBottomColor='gray.200'
              h='full'
            />
          </Flex>
        )}
      </CardBody>

      {externalUrl && (
        <CardFooter p='6' pt='0' pb='5'>
          <Flex pt='4' align='center'>
            <Link
              display='inline-flex'
              href={getExternalUrl(externalUrl)}
              fontSize='sm'
              color='primary.700'
              target='_blank'
              alignItems='center'
            >
              <Zendesk boxSize='4' mr='2' />
              View in Zendesk
            </Link>
          </Flex>
        </CardFooter>
      )}
    </>
  );
};
