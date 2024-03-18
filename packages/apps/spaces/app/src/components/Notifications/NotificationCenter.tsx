import React from 'react';
import { useRouter } from 'next/navigation';

import { IMessage, PopoverNotificationCenter } from '@novu/notification-center';

import { cn } from '@ui/utils/cn';
import { Tooltip } from '@ui/overlay/Tooltip/';
import { DateTimeUtils } from '@spaces/utils/date';
// import { Avatar, AvatarBadge } from '@ui/media/Avatar';
import { Avatar, AvatarBadge } from '@ui/media/Avatar/Avatar';
import { CountButton } from '@shared/components/Notifications/CountButton';
import { EmptyNotifications } from '@shared/components/Notifications/EmptyNotifications';
import { NotificationsHeader } from '@shared/components/Notifications/NotificationsHeader';

interface NotificationCenterProps {}

export const NotificationCenter: React.FC<NotificationCenterProps> = () => {
  const router = useRouter();
  const isProduction = process.env.NEXT_PUBLIC_PRODUCTION === 'true';

  if (isProduction) {
    // todo remove after feature is released to production
    return null;
  }

  function handlerOnNotificationClick(message: IMessage) {
    if (message?.cta?.data?.url) {
      router.push(message?.cta?.data?.url as string);
    }
  }

  return (
    <PopoverNotificationCenter
      colorScheme='light'
      position='right-end'
      showUserPreferences={false}
      emptyState={<EmptyNotifications />}
      header={() => <NotificationsHeader />}
      footer={() => <div />}
      theme={{
        light: {
          loaderColor: 'primary.500',
          popover: {
            arrowColor: 'transparent',
          },
        },
      }}
      onNotificationClick={handlerOnNotificationClick}
      listItem={(message, _, onNotificationClick) => {
        const parsedMessage = new DOMParser()?.parseFromString(
          message?.content as string,
          'text/html',
        )?.documentElement?.textContent;
        const content: false | string[] =
          typeof parsedMessage === 'string' &&
          parsedMessage?.split('owner of ');

        return (
          <Tooltip
            hasArrow
            placement='bottom-start'
            label={
              message.payload.isArchived
                ? 'This organization has been archived'
                : ''
            }
          >
            <div
              className={cn(
                message.payload.isArchived
                  ? 'cursor-default'
                  : 'cursor-pointer',
                'flex px-4 mb-5',
              )}
              role='button'
              tabIndex={message.payload.isArchived ? -1 : 0}
              onClick={
                message.payload.isArchived ? undefined : onNotificationClick
              }
              style={{
                cursor: message.payload.isArchived ? 'default' : 'pointer',
              }}
            >
              <Avatar
                size='sm'
                name={'UN'}
                variant='roundedSquareSmall'
                src={undefined}
                className={cn(message.read ? 'opacity-5' : 'opacity-10')}
                badge={
                  !message.seen ? (
                    <AvatarBadge className='bg-[#0BA5EC]' />
                  ) : (
                    <> </>
                  )
                }
              />
            </div>
            <div className='flex flex-col ml-3 gap-1 text-gray-700'>
              <p className='text-sm leading-4 truncate text-inherit'>
                {content && `${content[0]} owner of `}
                <span className='font-medium text-inherit'>
                  {content &&
                    (content[1]?.trim()?.length ? content[1] : 'Unnamed')}
                </span>
              </p>
              <p
                className={cn(
                  message.read ? 'text-gray-400' : 'text-gray-500',
                  'text-xs leading-4',
                )}
              >
                {DateTimeUtils.timeAgo(message?.createdAt as string, {
                  includeMin: true,
                  addSuffix: true,
                })}
              </p>
            </div>
          </Tooltip>
        );
      }}
    >
      {({ unseenCount }) => <CountButton unseenCount={unseenCount} />}
    </PopoverNotificationCenter>
  );
};
