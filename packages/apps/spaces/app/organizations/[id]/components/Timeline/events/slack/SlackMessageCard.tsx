import React, { PropsWithChildren } from 'react';
import { Card, CardBody } from '@ui/presentation/Card';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Avatar } from '@ui/media/Avatar';
import User from '@spaces/atoms/icons/User';
import { ViewInSlackButton } from '@organization/components/Timeline/events/slack/ViewInSlackButton';

interface SlackMessageCardProps extends PropsWithChildren {
  name: string;
  profilePhotoUrl?: null | string;
  content: string;
  onClick?: () => void;
  date: string;
  w?: any;
  showDateOnHover?: boolean;
}

export const SlackMessageCard: React.FC<SlackMessageCardProps> = ({
  name,
  profilePhotoUrl,
  content,
  onClick,
  children,
  date,
  w,
  showDateOnHover,
}) => {
  return (
    <>
      <Card
        variant='outline'
        size='md'
        fontSize='14px'
        background='white'
        flexDirection='row'
        maxWidth={w || 549}
        position='unset'
        cursor='pointer'
        boxShadow='xs'
        borderColor='gray.100'
        onClick={() => onClick?.()}
        _hover={{
          '&:hover .slack-stub-date': {
            color: 'gray.500',
          },
        }}
      >
        <CardBody p={3} overflow={'hidden'}>
          <Flex gap={3} flex={1}>
            <Avatar
              name={name}
              variant='roundedSquare'
              size='lg'
              icon={
                <User color={'var(--chakra-colors-gray-500)'} height='1.8rem' />
              }
              border={
                profilePhotoUrl
                  ? 'none'
                  : '2px solid var(--chakra-colors-primary-200)'
              }
              src={profilePhotoUrl || undefined}
            />
            <Flex direction='column' flex={1}>
              <Flex justifyContent='space-between' flex={1}>
                <Flex>
                  <Text color='gray.700' fontWeight={600}>
                    {name}
                  </Text>
                  <Text
                    color={showDateOnHover ? 'transparent' : 'gray.500'}
                    ml={2}
                    className='slack-stub-date'
                  >
                    {date}
                  </Text>
                </Flex>

                <ViewInSlackButton url='' />
              </Flex>

              <Text>{content}</Text>
              {children}
            </Flex>
          </Flex>
        </CardBody>
      </Card>
    </>
  );
};
