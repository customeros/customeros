import React from 'react';
import { useRouter } from 'next/navigation';

import { IMessage, PopoverNotificationCenter } from '@novu/notification-center';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { Tooltip } from '@ui/overlay/Tooltip';
import { Badge } from '@ui/presentation/Badge';
import { DateTimeUtils } from '@spaces/utils/date';
import { Avatar, AvatarBadge } from '@ui/media/Avatar';
import { ArrowsRight } from '@ui/media/icons/ArrowsRight';
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
        const content: false | string[] =
          typeof message.content === 'string' &&
          message.content.split('owner of ');

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
              opacity={message.read ? 0.5 : 1}
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
                size='sm'
                name={'UN'}
                variant='roundedSquareSmall'
                src={undefined}
              >
                {!message.seen && <AvatarBadge boxSize='10px' bg='#0BA5EC' />}
              </Avatar>
              <Flex direction='column' ml={3} gap={1}>
                <Text fontSize='sm' lineHeight='1' noOfLines={2}>
                  {content && content[0]}
                  <Text as='span' fontWeight='medium'>
                    {content && content[1]}
                  </Text>
                </Text>
                <Text fontSize='xs' lineHeight='1' color='gray.500'>
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
          <Flex justifyContent='space-between' flex={1} alignItems='center'>
            <span>Up next</span>
            {!!unseenCount && (
              <Badge
                w={5}
                h={5}
                display='flex'
                alignItems='center'
                justifyContent='center'
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
