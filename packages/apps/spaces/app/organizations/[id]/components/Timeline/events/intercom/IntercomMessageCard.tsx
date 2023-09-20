import React, { PropsWithChildren } from 'react';
import parse, {
  Element,
  domToReact,
  HTMLReactParserOptions,
} from 'html-react-parser';
import linkifyHtml from 'linkify-html';

import { Card, CardBody, CardProps } from '@ui/presentation/Card';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Avatar } from '@ui/media/Avatar';
import User from '@spaces/atoms/icons/User';

import { ViewInIntercomButton } from './ViewInIntercomButton';
import { ImageAttachment } from './ImageAttachment';

interface IntercomMessageCardProps extends PropsWithChildren {
  name: string;
  sourceUrl?: string | null;
  profilePhotoUrl?: null | string;
  content: string;
  onClick?: () => void;
  date: string;
  w?: CardProps['w'];
  showDateOnHover?: boolean;
}

export const IntercomMessageCard: React.FC<IntercomMessageCardProps> = ({
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
    return linkifyHtml(content, {
      defaultProtocol: 'https',
      rel: 'noopener noreferrer',
    });
  })();

  const parseOptions: HTMLReactParserOptions = {
    replace: (domNode) => {
      if (domNode instanceof Element) {
        switch (domNode.name) {
          case 'td': {
            return (
              <Flex
                flexDir='column'
                noOfLines={showDateOnHover ? 4 : undefined}
              >
                {domToReact(domNode.children)}
              </Flex>
            );
          }
          case 'img': {
            return <ImageAttachment {...domNode.attribs} />;
          }
          default:
            return;
        }
      }
    },
  };

  const parsedDisplayContent = parse(displayContent, parseOptions);

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
          '&:hover .intercom-stub-date': {
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
              icon={
                <User color={'var(--chakra-colors-gray-500)'} height='1.8rem' />
              }
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
                  <Text
                    color={showDateOnHover ? 'transparent' : 'gray.500'}
                    ml='2'
                    fontSize='xs'
                    className='intercom-stub-date'
                  >
                    {date}
                  </Text>
                </Flex>

                <ViewInIntercomButton url={sourceUrl} />
              </Flex>
              <Flex
                flexDir='column'
                pointerEvents={showDateOnHover ? 'none' : 'initial'}
                noOfLines={showDateOnHover ? 4 : undefined}
                sx={{
                  '& ol, ul': {
                    pl: '5',
                  },
                  '& pre': {
                    whiteSpace: 'normal',
                    fontSize: '12px',
                    color: 'gray.700',
                    border: '1px solid',
                    borderColor: 'gray.300',
                    borderRadius: '4',
                    p: '2',
                    py: '1',
                    my: '2',
                  },
                }}
              >
                {parsedDisplayContent}
              </Flex>
              {children}
            </Flex>
          </Flex>
        </CardBody>
      </Card>
    </>
  );
};
