import React from 'react';
import { useRouter } from 'next/navigation';

import { IMessage, PopoverNotificationCenter } from '@novu/notification-center';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Tooltip } from '@ui/overlay/Tooltip';
import { DateTimeUtils } from '@spaces/utils/date';
import { Avatar, AvatarBadge } from '@ui/media/Avatar';
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
            <Flex
              px={4}
              mb={5}
              role='button'
              cursor={message.payload.isArchived ? 'default' : 'pointer'}
              tabIndex={message.payload.isArchived ? -1 : 0}
              onClick={
                message.payload.isArchived ? undefined : onNotificationClick
              }
            >
              <Avatar
                opacity={message.read ? 0.5 : 1}
                size='sm'
                name={'UN'}
                variant='roundedSquareSmall'
                src={undefined}
              >
                {!message.seen && <AvatarBadge boxSize='10px' bg='#0BA5EC' />}
              </Avatar>
              <Flex
                direction='column'
                ml={3}
                gap={1}
                color={message.read ? 'gray.400' : 'gray.700'}
              >
                <Text
                  fontSize='sm'
                  lineHeight='1'
                  noOfLines={2}
                  color='inherit'
                >
                  {content && `${content[0]} owner of `}
                  <Text as='span' fontWeight='medium' color='inherit'>
                    {content &&
                      (content[1]?.trim()?.length ? content[1] : 'Unnamed')}
                  </Text>
                </Text>
                <Text
                  fontSize='xs'
                  lineHeight='1'
                  color={message.read ? 'gray.400' : 'gray.500'}
                >
                  {DateTimeUtils.timeAgo(message?.createdAt as string, {
                    includeMin: true,
                    addSuffix: true,
                  })}
                </Text>
              </Flex>
            </Flex>
          </Tooltip>
        );
      }}
    >
      {({ unseenCount }) => <CountButton unseenCount={unseenCount} />}
    </PopoverNotificationCenter>
  );
};
