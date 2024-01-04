import React from 'react';
import { useRouter } from 'next/navigation';

import { IMessage, PopoverNotificationCenter } from '@novu/notification-center';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { Badge } from '@ui/presentation/Badge';
import { DateTimeUtils } from '@spaces/utils/date';
import { Avatar, AvatarBadge } from '@ui/media/Avatar';
import { ArrowsRight } from '@ui/media/icons/ArrowsRight';
import { EmptyNotifications } from '@shared/components/Notifications/EmptyNotifications';
import { NotificationsHeader } from '@shared/components/Notifications/NotificationsHeader';

interface NotificationCenterProps {}

export const NotificationCenter: React.FC<NotificationCenterProps> = () => {
  const router = useRouter();

  function handlerOnNotificationClick(message: IMessage) {
    if (message?.cta?.data?.url) {
      router.push(message.cta.data.url);
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
      listItem={(message, _, onNotificationClick) => (
        <Flex
          px={4}
          mb={5}
          role='button'
          cursor='pointer'
          tabIndex={0}
          onClick={onNotificationClick}
        >
          <Avatar
            size='sm'
            name={'UN'}
            variant='roundedSquareSmall'
            src={undefined}
          >
            {!message.read && <AvatarBadge boxSize='10px' bg='#0BA5EC' />}
          </Avatar>
          <Flex direction='column' ml={3} gap={1}>
            <Text fontSize='sm' lineHeight='1' noOfLines={2}>
              {message?.content as string}
            </Text>
            <Text fontSize='xs' lineHeight='1' color='gray.500'>
              {DateTimeUtils.timeAgo(message?.createdAt as string, {
                includeMin: true,
                addSuffix: true,
              })}
            </Text>
          </Flex>
        </Flex>
      )}
    >
      {({ unseenCount }) => (
        <Button
          px='3'
          w='full'
          size='md'
          variant='ghost'
          fontSize='sm'
          textDecoration='none'
          fontWeight='regular'
          justifyContent='flex-start'
          borderRadius='md'
          color={'gray.500'}
          leftIcon={<ArrowsRight color='inherit' boxSize='5' />}
          _focus={{
            boxShadow: 'sidenavItemFocus',
          }}
        >
          <Flex justifyContent='space-between' flex={1}>
            <span>Up next</span>
            {!!unseenCount && (
              <Badge
                px={5}
                variant='outline'
                borderRadius='xl'
                boxShadow='none'
                border='1px solid'
                borderColor='gray.300'
                fontWeight='regular'
              >
                {unseenCount}
              </Badge>
            )}
          </Flex>
        </Button>
      )}
    </PopoverNotificationCenter>
  );
};
