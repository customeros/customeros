import { escapeForSlackWithMarkdown } from 'slack-to-html';

import { cn } from '@ui/utils/cn';
import { Slack } from '@ui/media/logos/Slack';
import { User01 } from '@ui/media/icons/User01';
import { Avatar } from '@ui/media/Avatar/Avatar';
import { ViewInExternalAppButton } from '@ui/form/Button';
import { Card, CardContent } from '@ui/presentation/Card/Card';

interface SlackMessageCardProps {
  name: string;
  date: string;
  content: string;
  className?: string;
  onClick?: () => void;
  sourceUrl?: string | null;
  showDateOnHover?: boolean;
  children?: React.ReactNode;
  profilePhotoUrl?: null | string;
}

export const SlackMessageCard = ({
  name,
  sourceUrl,
  profilePhotoUrl,
  content,
  onClick,
  className,
  children,
  date,
  showDateOnHover,
}: SlackMessageCardProps) => {
  const displayContent: string = (() => {
    const sanitizeContent = content.replace(/\n/g, '<br/>');
    const slack = escapeForSlackWithMarkdown(sanitizeContent);
    const regex = /(?<=^|\s)@(\w+)/g;

    return slack.replace(
      regex,
      (matched: string): string =>
        `<span class='slack-mention'>${matched.replace(/_/g, ' ')}</span>`,
    );
  })();

  return (
    <>
      <Card
        onClick={() => onClick?.()}
        className={cn(
          className,
          onClick ? 'cursor-pointer' : '',
          'max-w-[549px] text-sm bg-white flex shadow-xs border border-gray-200 [slack-stub-date]:hover:text-gray-500 hover:shadow-md transition-all duration-200 ease-out',
        )}
      >
        <CardContent className='p-3 overflow-hidden w-full'>
          <div className='flex flex-1 gap-3'>
            <Avatar
              size='md'
              name={name}
              variant='roundedSquare'
              src={profilePhotoUrl || undefined}
              icon={<User01 className='text-gray-500 size-7' />}
              className={cn(profilePhotoUrl ? '' : 'border border-gray-200')}
            />
            <div className='flex flex-col flex-1 relative'>
              <div className='flex justify-between flex-1'>
                <div className='flex items-center'>
                  <p className='text-gray-700 font-semibold'>{name}</p>
                  <p
                    className={cn(
                      showDateOnHover ? 'transparent' : 'text-gray-500',
                      'ml-2 text-xs slack-stub-date',
                    )}
                  >
                    {date}
                  </p>
                </div>
                <ViewInExternalAppButton url={sourceUrl} icon={<Slack />} />
              </div>
              <p
                dangerouslySetInnerHTML={{ __html: displayContent }}
                className={cn(
                  showDateOnHover
                    ? 'pointer-events-none line-clamp-4'
                    : 'pointer-events-auto',
                  'slack-container',
                )}
              />
              {children}
            </div>
          </div>
        </CardContent>
      </Card>
    </>
  );
};
