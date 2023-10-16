import { FC, PropsWithChildren } from 'react';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Avatar } from '@ui/media/Avatar';
import { User01 } from '@ui/media/icons/User01';
import { Card, CardBody } from '@ui/presentation/Card';
import { ViewInExternalAppButton } from '@ui/form/Button';
import { DateTimeUtils } from '@spaces/utils/date';
import { Zendesk } from '@ui/media/logos/Zendesk';

interface IssueCommentCardProps extends PropsWithChildren {
  name: string;
  sourceUrl?: string | null;
  profilePhotoUrl?: null | string;
  content: string;
  date: string;
  isCustomer?: boolean;

  showDateOnHover?: boolean;
}

export const IssueCommentCard: FC<IssueCommentCardProps> = ({
  name,
  sourceUrl,
  profilePhotoUrl,
  content,
  date,
  isCustomer,
}) => {
  return (
    <>
      <Card
        ml={isCustomer ? 0 : 6}
        variant='outline'
        size='md'
        fontSize='14px'
        background='white'
        flexDirection='row'
        width='calc(100% - 24px)'
        position='unset'
        cursor='unset'
        boxShadow='xs'
        borderColor='gray.200'
      >
        <CardBody p={3} overflow={'hidden'}>
          <Flex gap={3} flex={1}>
            <Avatar
              name={name}
              size='md'
              icon={<User01 color='gray.500' height='1.8rem' />}
              border={
                profilePhotoUrl
                  ? 'none'
                  : '1px solid var(--chakra-colors-primary-200)'
              }
              src={profilePhotoUrl || undefined}
            />
            <Flex direction='column' flex={1} position='relative'>
              <Flex justifyContent='space-between' flex={1}>
                <Flex>
                  <Text color='gray.700' fontWeight={600}>
                    {name}
                  </Text>
                  <Text color='gray.500' ml={2} fontSize='xs'>
                    {DateTimeUtils.formatTime(date)}
                  </Text>
                </Flex>

                <ViewInExternalAppButton
                  icon={<Zendesk boxSize={4} />}
                  url={sourceUrl}
                />
              </Flex>
              <Text fontSize='sm'>{content}</Text>
            </Flex>
          </Flex>
        </CardBody>
      </Card>
    </>
  );
};
