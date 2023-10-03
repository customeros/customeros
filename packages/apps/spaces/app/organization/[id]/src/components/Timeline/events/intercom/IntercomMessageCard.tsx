import React, { PropsWithChildren } from 'react';
import { Card, CardBody, CardProps } from '@ui/presentation/Card';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Avatar } from '@ui/media/Avatar';
import User from '@spaces/atoms/icons/User';

import { ViewInIntercomButton } from './ViewInIntercomButton';
import { HtmlContentRenderer } from '@ui/presentation/HtmlContentRenderer/HtmlContentRenderer';

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

              <HtmlContentRenderer
                pointerEvents={showDateOnHover ? 'none' : 'initial'}
                noOfLines={showDateOnHover ? 4 : undefined}
                htmlContent={content}
              />
              {children}
            </Flex>
          </Flex>
        </CardBody>
      </Card>
    </>
  );
};
