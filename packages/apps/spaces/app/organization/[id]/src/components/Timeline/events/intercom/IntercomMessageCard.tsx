import React, { PropsWithChildren } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Avatar } from '@ui/media/Avatar';
import { Text } from '@ui/typography/Text';
import { User01 } from '@ui/media/icons/User01';
import Intercom from '@ui/media/icons/Intercom';
import { ViewInExternalAppButton } from '@ui/form/Button';
import { Card, CardBody, CardProps } from '@ui/presentation/Card';
import { HtmlContentRenderer } from '@ui/presentation/HtmlContentRenderer/HtmlContentRenderer';

interface IntercomMessageCardProps extends PropsWithChildren {
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

export const IntercomMessageCard: React.FC<IntercomMessageCardProps> = ({
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
  return (
    <>
      <Card
        variant='outline'
        size='md'
        ml={ml}
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
              icon={<User01 color='gray.500' height='1.8rem' />}
              border={
                profilePhotoUrl
                  ? 'none'
                  : '1px solid var(--chakra-colors-primary-200)'
              }
              src={profilePhotoUrl || undefined}
            />
            <Flex
              direction='column'
              flex={1}
              position='relative'
              maxWidth={showDateOnHover ? 470 : 408}
            >
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

                <ViewInExternalAppButton
                  url={sourceUrl}
                  icon={
                    <Flex alignItems='center' justifyContent='center'>
                      <Intercom height={10} />
                    </Flex>
                  }
                />
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
