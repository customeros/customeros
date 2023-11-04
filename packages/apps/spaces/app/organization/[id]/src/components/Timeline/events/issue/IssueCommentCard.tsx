import { FC, PropsWithChildren } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Avatar } from '@ui/media/Avatar';
import { Text } from '@ui/typography/Text';
import { User01 } from '@ui/media/icons/User01';
import { DateTimeUtils } from '@spaces/utils/date';
import { Card, CardBody } from '@ui/presentation/Card';
import { HtmlContentRenderer } from '@ui/presentation/HtmlContentRenderer';

interface IssueCommentCardProps extends PropsWithChildren {
  name: string;
  date: string;
  content: string;
  isPrivate?: boolean;
  iscustomer?: boolean;
  showDateOnHover?: boolean;
  profilePhotoUrl?: null | string;
}

export const IssueCommentCard: FC<IssueCommentCardProps> = ({
  name,
  date,
  content,
  isPrivate,
  iscustomer,
  profilePhotoUrl,
}) => {
  return (
    <>
      <Card
        ml={iscustomer ? 0 : 6}
        variant='outline'
        size='md'
        fontSize='14px'
        background='white'
        flexDirection='row'
        width='calc(100% - 24px)'
        position='unset'
        cursor='unset'
        boxShadow={isPrivate ? 'none' : 'xs'}
        borderColor={isPrivate ? 'transparent' : 'gray.200'}
        bg={isPrivate ? 'transparent' : 'white'}
      >
        <CardBody p={3} overflow={'hidden'}>
          <Flex gap={3} flex={1}>
            <Avatar
              name={name}
              size='md'
              icon={<User01 color='primary.500' boxSize='5' />}
              border={
                profilePhotoUrl
                  ? 'none'
                  : '1px solid var(--chakra-colors-primary-200)'
              }
              src={profilePhotoUrl || undefined}
            />
            <Flex direction='column' flex={1} position='relative'>
              <Flex justifyContent='space-between' flex={1}>
                <Flex align='baseline'>
                  <Text color='gray.700' fontWeight={600}>
                    {name}
                  </Text>
                  <Text color='gray.500' ml={2} fontSize='xs'>
                    {DateTimeUtils.formatTime(date)}
                  </Text>
                </Flex>
              </Flex>
              <HtmlContentRenderer htmlContent={content} />
            </Flex>
          </Flex>
        </CardBody>
      </Card>
    </>
  );
};
