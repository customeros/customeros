'use client';
import React, { FC } from 'react';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { CardBody, CardHeader, CardFooter, Card } from '@ui/layout/Card';
import { IssueBgPattern } from '@ui/media/logos/IssueBgPattern';
import { Tag, TagLabel } from '@ui/presentation/Tag';
import { CustomTicketTearStyle } from './styles';
import { IssueWithAliases } from '@organization/src/components/Timeline/types';
import { useTimelineEventPreviewMethodsContext } from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';
import { MarkdownContentRenderer } from '@ui/presentation/MarkdownContentRenderer/MarkdownContentRenderer';
function getStatusColor(status: string) {
  if (status === 'solved' || status === 'closed') {
    return 'gray';
  }
  return 'blue';
}

export const IssueStub: FC<{ data: IssueWithAliases }> = ({ data }) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();
  const statusColorScheme = getStatusColor(data.issueStatus);

  return (
    <Card
      variant='outline'
      size='md'
      fontSize='14px'
      background='white'
      flexDirection='row'
      position='unset'
      maxW={502}
      cursor='pointer'
      boxShadow='none'
      border='1px solid'
      borderColor='gray.200'
      transition='all 0.2s ease-out'
      _hover={{
        boxShadow: 'md',
        '& .footer:after': {
          borderBottomColor: '#e5e3e3',
          transition: 'all 3s ease-out',
          background: 'linear-gradient(182deg, #dcd9d9, #ececec)',
        },
      }}
      onClick={() => openModal(data.id)}
    >
      <Flex boxShadow='xs' pr={2} p={3} direction='column' flex={1}>
        <CardHeader fontWeight='semibold' p={0} noOfLines={1}>
          {data?.subject ?? '[No subject]'}
        </CardHeader>
        <CardBody p={0} maxW='calc(476px - 77px)'>
          <Text color='gray.500' noOfLines={3} fontSize='sm'>
            {data?.description ? (
              <MarkdownContentRenderer
                markdownContent={data?.description}
                showAsInlineText
              />
            ) : (
              '[No description]'
            )}
          </Text>
        </CardBody>
      </Flex>
      <CardFooter
        p={0}
        className='footer'
        position='relative'
        h='108px'
        display='flex'
        flexDirection='column'
        justifyContent='center'
        minW='72px'
        borderLeft='1px dashed'
        borderColor='gray.200'
        boxShadow='xs'
        sx={CustomTicketTearStyle}
      >
        <Flex
          direction='column'
          alignItems='center'
          justifyContent='center'
          overflow='hidden'
          h='103px'
          minW='66px'
          position='relative'
          borderRadius='md'
        >
          {!!data?.externalLinks?.length && (
            <Text mb={1} zIndex={1} fontWeight='semibold' color='gray.500'>
              {data?.externalLinks[0]?.externalId}
            </Text>
          )}

          <Tag
            zIndex={1}
            size='sm'
            variant='outline'
            colorScheme='blue'
            border='1px solid'
            background='white'
            borderColor={`${[statusColorScheme]}.200`}
            backgroundColor={`${[statusColorScheme]}.50`}
            color={`${[statusColorScheme]}.700`}
            boxShadow='none'
            fontWeight='normal'
            minHeight={6}
            width='min-content'
            cursor='pointer'
          >
            <TagLabel>
              {['solved', 'closed'].includes(data.issueStatus)
                ? 'Closed'
                : 'Open'}
            </TagLabel>
          </Tag>
          <IssueBgPattern position='absolute' width='120%' height='100%' />
        </Flex>
      </CardFooter>
    </Card>
  );
};
