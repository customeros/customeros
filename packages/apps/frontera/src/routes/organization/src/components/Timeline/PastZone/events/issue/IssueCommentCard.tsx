import { FC, PropsWithChildren } from 'react';

import { match } from 'ts-pattern';
import { escapeForSlackWithMarkdown } from 'slack-to-html';

import { cn } from '@ui/utils/cn';
import { DateTimeUtils } from '@utils/date';
import { User01 } from '@ui/media/icons/User01';
import { Avatar } from '@ui/media/Avatar/Avatar';
import { ExternalSystemType } from '@graphql/types';
import { Card, CardContent } from '@ui/presentation/Card/Card';
import { HtmlContentRenderer } from '@ui/presentation/HtmlContentRenderer';

interface IssueCommentCardProps extends PropsWithChildren {
  name: string;
  date: string;
  type?: string;
  content: string;
  isPrivate?: boolean;
  isCustomer?: boolean;
  showDateOnHover?: boolean;
  profilePhotoUrl?: null | string;
}

export const IssueCommentCard: FC<IssueCommentCardProps> = ({
  name,
  date,
  type,
  content,
  isPrivate,
  isCustomer,
  profilePhotoUrl,
}) => {
  return (
    <>
      <Card
        className={cn(
          isPrivate ? 'bg-transparent' : 'bg-white shadow-xs',
          isCustomer ? 'ml-0' : 'ml-6',
          'text-[14px] flex flex-row w-[calc(100%-24px)]',
        )}
      >
        <CardContent className='p-3 overflow-hidden'>
          <div className='flex gap-3 flex-1'>
            <Avatar
              name={name}
              size='md'
              icon={<User01 className='text-primary-500 size-5' />}
              className={cn(profilePhotoUrl ? '' : 'border border-primary-200')}
              src={profilePhotoUrl || undefined}
            />
            <div className='flex flex-col flex-1 relative'>
              <div className='flex justify-between flex-1'>
                <div className='flex items-baseline'>
                  <span className='text-gray-700 font-semibold'>{name}</span>
                  <span className='text-gray-500 ml-2 text-xs'>
                    {DateTimeUtils.formatTime(date)}
                  </span>
                </div>
              </div>
              {match(type)
                .with(ExternalSystemType.Slack, () => (
                  <span
                    className='slack-container'
                    dangerouslySetInnerHTML={{
                      __html: escapeForSlackWithMarkdown(content),
                    }}
                  />
                ))
                .otherwise(() => (
                  <HtmlContentRenderer htmlContent={content} />
                ))}
            </div>
          </div>
        </CardContent>
      </Card>
    </>
  );
};
