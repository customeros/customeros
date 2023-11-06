import { PropsWithChildren } from 'react';

// @ts-expect-error types not available
import { escapeForSlackWithMarkdown } from 'slack-to-html';

import { Flex } from '@ui/layout/Flex';
import { Avatar } from '@ui/media/Avatar';
import { Text } from '@ui/typography/Text';
import { Slack } from '@ui/media/logos/Slack';
import { User01 } from '@ui/media/icons/User01';
import { ViewInExternalAppButton } from '@ui/form/Button';
import { Card, CardBody, CardProps } from '@ui/presentation/Card';

interface SlackMessageCardProps extends PropsWithChildren {
  name: string;
  date: string;
  content: string;
  w?: CardProps['w'];
  onClick?: () => void;
  ml?: CardProps['ml'];
  sourceUrl?: string | null;
  showDateOnHover?: boolean;
  profilePhotoUrl?: null | string;
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
  ml,
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
        ml={ml}
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
                <ViewInExternalAppButton
                  icon={<Slack height={16} />}
                  url={sourceUrl}
                />
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
