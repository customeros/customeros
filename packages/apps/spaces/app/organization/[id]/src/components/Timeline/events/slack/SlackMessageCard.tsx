import { PropsWithChildren } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Avatar } from '@ui/media/Avatar';
import { User01 } from '@ui/media/icons/User01';
import { Card, CardBody, CardProps } from '@ui/presentation/Card';
import { ViewInSlackButton } from '@organization/src/components/Timeline/events/slack/ViewInSlackButton';
// @ts-expect-error types not available
import { escapeForSlackWithMarkdown } from 'slack-to-html';

interface SlackMessageCardProps extends PropsWithChildren {
  name: string;
  sourceUrl?: string | null;
  profilePhotoUrl?: null | string;
  content: string;
  onClick?: () => void;
  date: string;
  w?: CardProps['w'];
  showDateOnHover?: boolean;
}

export const SlackMessageCard: React.FC<SlackMessageCardProps> = ({
  name,
  sourceUrl,
  profilePhotoUrl,
  content,
  onClick,
  children,
  date,
  w,
  showDateOnHover,
}) => {
  const displayContent: string = (() => {
    const sanitizeContent = content.replace(/\n/g, '<br/>');
    const slack = escapeForSlackWithMarkdown(sanitizeContent);
    const regex = /(@[\w]+)/g;
    return slack.replace(
      regex,
      (matched: string): string =>
        `<span class='slack-mention'>${matched.replace(/_/g, ' ')}</span>`,
    );
  })();

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
        cursor={onClick ? 'pointer' : 'unset'}
        boxShadow='xs'
        borderColor='gray.200'
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
                  <Text
                    color={showDateOnHover ? 'transparent' : 'gray.500'}
                    ml={2}
                    fontSize='xs'
                    className='slack-stub-date'
                  >
                    {date}
                  </Text>
                </Flex>

                <ViewInSlackButton url={sourceUrl} />
              </Flex>
              <Text
                className='slack-container'
                pointerEvents={showDateOnHover ? 'none' : 'initial'}
                noOfLines={showDateOnHover ? 4 : undefined}
                dangerouslySetInnerHTML={{ __html: displayContent }}
              />
              {children}
            </Flex>
          </Flex>
        </CardBody>
      </Card>
    </>
  );
};
