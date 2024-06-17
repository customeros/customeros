import { useRef, useMemo } from 'react';

import { DateTimeUtils } from '@utils/date';
import { User01 } from '@ui/media/icons/User01';
import { Issue, Contact } from '@graphql/types';
import { Avatar } from '@ui/media/Avatar/Avatar';
import { Tag, TagLabel } from '@ui/presentation/Tag/Tag';
import { Card, CardHeader } from '@ui/presentation/Card/Card';
import { getParticipant, getParticipantName } from '@organization/hooks/utils';
import { useTimelineEventPreviewMethodsContext } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

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
    [issue.submittedBy],
  );
  const reportedBy = useMemo(
    () => (issue.reportedBy ? getParticipantName(issue.reportedBy) : undefined),

    [issue.reportedBy],
  );
  const profilePhoto = useMemo(
    () =>
      issue.reportedBy
        ? getParticipant(issue.reportedBy)
        : issue?.submittedBy
        ? getParticipant(issue.submittedBy)
        : undefined,

    [issue.reportedBy, issue.submittedBy],
  );

  const participantName = useMemo(
    () => (
      <span className='inline'>
        {reportedBy ? `Reported` : `Submitted`}

        {(reportedBy || submittedBy) && <span className='mx-1'>by</span>}
        <span className='font-bold mr-1'>
          {reportedBy ? reportedBy : submittedBy}
        </span>
      </span>
    ),
    [reportedBy, submittedBy],
  );

  const titleWidth = useMemo(() => {
    if (isStatusClosed) {
      return 'auto';
    }

    return issue?.status === 'pending' ? 250 : 260;
  }, [isStatusClosed, issue?.status]);

  const displayStatus = issue.status.replaceAll('_', ' ').replaceAll('-', ' ');

  return (
    <Card
      key={issue.id}
      className='w-full shadow-xs cursor-pointer rounded-lg border border-gray-200 bg-white hover:shadow-md p-3 max-w-[400px]'
      ref={cardRef}
      onClick={() => openModal(issue.id)}
    >
      <CardHeader>
        <div className='flex flex-1 gap-2 items-start flex-wrap relative'>
          <Avatar
            size='md'
            name={submittedBy ?? reportedBy}
            variant='circle'
            className='border border-primary-200 text-primary-700'
            src={
              (profilePhoto as unknown as Contact)?.profilePhotoUrl ?? undefined
            }
            icon={<User01 color='primary.700' height='1.8rem' />}
          />

          <div className='flex flex-col flex-1 line-clamp-1 ml-2'>
            <h2
              className='text-sm  font-semibold'
              style={{ maxWidth: titleWidth }}
            >
              {issue?.subject ?? '[No subject]'}
            </h2>

            <span className='text-sm mt-1 mb-[2px] leading-3 relative'>
              {participantName}
              {DateTimeUtils.timeAgo(issue?.createdAt, { addSuffix: true })}
            </span>

            {!!issue?.updatedAt && (
              <span className='text-sm text-gray-500 leading-3'>
                Last response was{' '}
                {DateTimeUtils.timeAgo(issue.updatedAt, {
                  addSuffix: true,
                })}
              </span>
            )}
          </div>

          {!isStatusClosed && (
            <Tag size='md' variant='outline' colorScheme={statusColorScheme}>
              <TagLabel className='capitalize'>{displayStatus}</TagLabel>
            </Tag>
          )}
        </div>
      </CardHeader>
    </Card>
  );
};
