import React from 'react';
import { useNavigate } from 'react-router-dom';

import {
  IMessage,
  ListItem,
  PopoverNotificationCenter,
} from '@novu/notification-center';

import { cn } from '@ui/utils/cn';
import { DateTimeUtils } from '@utils/date';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Avatar, AvatarBadge } from '@ui/media/Avatar/Avatar';
import { CountButton } from '@shared/components/Notifications/CountButton';
import { EmptyNotifications } from '@shared/components/Notifications/EmptyNotifications';
import { NotificationsHeader } from '@shared/components/Notifications/NotificationsHeader';

import './override.css';

interface NotificationCenterProps {}

export const NotificationCenter: React.FC<NotificationCenterProps> = () => {
  const navigate = useNavigate();

  function handlerOnNotificationClick(message: IMessage) {
    if (message?.cta?.data?.url) {
      navigate(message?.cta?.data?.url as string);
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
          loaderColor: '#9E77ED',
          popover: {
            arrowColor: 'transparent',
          },
          layout: {
            borderRadius: '16px',
          },
        },
      }}
      onNotificationClick={handlerOnNotificationClick}
      listItem={CustomListItem}
    >
      {({ unseenCount }) => <CountButton unseenCount={unseenCount} />}
    </PopoverNotificationCenter>
  );
};

const CustomListItem: ListItem = (message, _, onNotificationClick) => {
  const parsedMessage = new DOMParser()?.parseFromString(
    message?.content as string,
    'text/html',
  )?.documentElement?.textContent;

  const content: false | string[] =
    typeof parsedMessage === 'string' && parsedMessage?.split('owner of ');

  const cursorClass = cn(
    message.payload.isArchived ? 'cursor-default' : 'cursor-pointer',
  );

  return (
    <Tooltip
      hasArrow
      label={message.content as string}
      side='bottom'
      align='center'
    >
      <div className='flex ml-6 mr-4 gap-2 mt-2 mb-3 items-start'>
        <div
          className={cn('flex', cursorClass)}
          role='button'
          tabIndex={message.payload.isArchived ? -1 : 0}
          onClick={message.payload.isArchived ? undefined : onNotificationClick}
        >
          <Avatar
            size='sm'
            name={'UN'}
            variant='roundedSquareSmall'
            src={undefined}
            className={cn(message.read ? 'opacity-50' : 'opacity-100')}
            badge={
              !message.seen ? <AvatarBadge className='bg-[#0BA5EC]' /> : <> </>
            }
          />
          <div className='flex flex-col text-gray-700'>
            <p
              className={cn(
                'text-sm leading-4 truncate text-inherit',
                cursorClass,
              )}
            >
              {content && `${content[0]} owner of `}
              <span className='font-medium text-inherit'>
                {content &&
                  (content[1]?.trim()?.length ? content[1] : 'Unnamed')}
              </span>
            </p>
            <p
              className={cn(
                'text-xs leading-4 cursor-default',
                message.read ? 'text-gray-400' : 'text-gray-500',
              )}
            >
              {DateTimeUtils.timeAgo(message?.createdAt as string, {
                includeMin: true,
                addSuffix: true,
              })}
            </p>
          </div>
        </div>
      </div>
    </Tooltip>
  );
};
