import { FC, Fragment } from 'react';
import { Link } from 'react-router-dom';

import { match } from 'ts-pattern';

import { DateTimeUtils } from '@utils/date';
import { Zendesk } from '@ui/media/logos/Zendesk';
import { getName } from '@utils/getParticipantsName';
import { getExternalUrl } from '@utils/getExternalLink';
import { Tag, TagLabel } from '@ui/presentation/Tag/Tag';
import { Divider } from '@ui/presentation/Divider/Divider';
import { Comment, InteractionEvent } from '@graphql/types';
import { CardFooter, CardContent } from '@ui/presentation/Card/Card';
import { IssueWithAliases } from '@organization/components/Timeline/types';
import { MarkdownContentRenderer } from '@ui/presentation/MarkdownContentRenderer/MarkdownContentRenderer';
import { IssueCommentCard } from '@organization/components/Timeline/PastZone/events/issue/IssueCommentCard';
import {
  Priority,
  PriorityBadge,
} from '@organization/components/Timeline/PastZone/events/issue/PriorityBadge';
import { TimelineEventPreviewHeader } from '@organization/components/Timeline/shared/TimelineEventPreview/header/TimelineEventPreviewHeader';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

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
        parse='slack'
      />

      <CardContent className='mt-0 p-6 pt-0 max-h-[calc(100vh-195px)] overflow-auto'>
        <div className='flex gap-2 mb-2 sticky top-0 bg-white pb-2 z-10'>
          {issue?.priority && (
            <PriorityBadge
              priority={issue.priority.toLowerCase() as Priority}
            />
          )}
          <Tag
            size='md'
            variant='outline'
            className='min-h-6 font-normal'
            colorScheme={statusColorScheme}
          >
            <TagLabel className='capitalize'>
              {issue.issueStatus.replaceAll('_', ' ').replaceAll('-', ' ')}
            </TagLabel>
          </Tag>
          <Tag
            size='md'
            className='bg-white border-gray-200 text-gray-500 font-normal min-h-6'
            variant='outline'
            colorScheme='gray'
          >
            <TagLabel>#{issue?.externalLinks?.[0]?.externalId}</TagLabel>
          </Tag>
        </div>
        <MarkdownContentRenderer
          className='text-sm mb-2'
          markdownContent={issue?.description ?? ''}
        />

        {issue?.tags?.length && (
          <span className='text-gray-500 text-sm mb-6'>
            {issue.tags.map((t) => t?.name).join(' • ')}
          </span>
        )}

        <div className='flex items-center mb-2'>
          <div className='flex items-baseline'>
            <span className='text-sm whitespace-nowrap'>
              {`Issue ${submittedBy ? 'submitted' : 'reported'} by`}
            </span>
            <span className='mx-1 text-sm whitespace-nowrap font-medium'>
              {submittedBy || reportedBy}
            </span>
            <span className='text-gray-400 text-sm whitespace-nowrap ml-1 mr-2'>
              {DateTimeUtils.format(
                issue?.createdAt,
                DateTimeUtils.dateWithHour,
              )}
            </span>
          </div>
          <Divider />
        </div>

        <div className='flex flex-col justify-start items-start w-full gap-2'>
          {Object.entries(commentsByDay)?.map(([date, comments]) => (
            <Fragment key={date}>
              {!DateTimeUtils.isSameDay(issue?.createdAt, date) && (
                <div className='flex items-center w-full'>
                  <span className='text-sm whitespace-nowrap text-gray-400 mr-2'>
                    {DateTimeUtils.format(date, DateTimeUtils.date)}
                  </span>
                  <Divider />
                </div>
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
                    isCustomer={isCustomer}
                    content={c.content ?? ''}
                    type={issue?.externalLinks?.[0]?.type}
                    isPrivate={c.__typename === 'Comment'}
                  />
                );
              })}
            </Fragment>
          ))}
        </div>

        {['solved', 'closed'].includes(issue.issueStatus?.toLowerCase()) && (
          <div className='flex items-center mt-2'>
            <div className='flex items-baseline'>
              <span className='text-sm whitespace-nowrap'>Issue closed</span>

              <span className='text-gray-400 text-sm whitespace-nowrap ml-2 mr-2'>
                {DateTimeUtils.format(
                  issue?.updatedAt,
                  DateTimeUtils.dateWithHour,
                )}
              </span>
            </div>
            <Divider />
          </div>
        )}
      </CardContent>

      {externalUrl && (
        <CardFooter className='p-6 pt-0 pb-5'>
          <div className='flex pt-4 align-middle'>
            <Link
              className='text-primary-700 inline-flex text-sm items-center'
              to={getExternalUrl(externalUrl)}
              target='_blank'
            >
              <Zendesk className='mr-2' />
              View in Zendesk
            </Link>
          </div>
        </CardFooter>
      )}
    </>
  );
};
