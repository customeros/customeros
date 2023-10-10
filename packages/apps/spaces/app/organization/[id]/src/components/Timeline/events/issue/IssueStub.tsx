'use client';
import React, { FC } from 'react';
import { useTimelineEventPreviewContext } from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';
import { Flex } from '@ui/layout/Flex';
import { Box, Card, CardFooter, CardHeader, Text } from '@chakra-ui/react';
import { CardBody } from '@chakra-ui/card';
import { IssueBgPattern } from '@ui/media/logos/IssueBgPattern';
import { Tag, TagLabel } from '@ui/presentation/Tag';
import { DataSource } from '@graphql/types';
import { CustomTicketTearStyle } from './styles';
function getStatusColor(
  status: 'New' | 'Open' | 'Pending' | 'On hold' | 'Solved',
) {
  if (status === 'Solved') {
    return 'gray';
  }
  return 'blue';
}

// TODO generate issue types and use then instead of any
export const IssueStub: FC<{ issueEvent: any }> = ({ issueEvent }) => {
  const { openModal } = useTimelineEventPreviewContext();
  const statusColorScheme = (() =>
    issueEvent.appSource === DataSource.ZendeskSupport
      ? getStatusColor(issueEvent.status)
      : 'gray')();
  return (
    <Card
      variant='outline'
      size='md'
      fontSize='14px'
      background='white'
      flexDirection='row'
      position='unset'
      maxW={476}
      cursor='pointer'
      boxShadow='none'
      border='1px solid'
      borderColor='gray.200'
      onClick={() => openModal(issueEvent)}
      _hover={{
        '&:hover .slack-stub-date': {
          color: 'gray.500',
        },
      }}
    >
      <Box boxShadow='xs' pr={2}>
        <CardHeader fontWeight='semibold' p={2} pb={0} pr={0} noOfLines={1}>
          {issueEvent?.subject ?? '[No subject]'}
        </CardHeader>
        <CardBody p={2} pt={0} pr={0}>
          <Text color='gray.500' noOfLines={3}>
            {issueEvent?.description ?? '[No description]'}
          </Text>
        </CardBody>
      </Box>
      <CardFooter
        p={0}
        position='relative'
        h='100px'
        display='flex'
        flexDirection='column'
        justifyContent='center'
        minW='66px'
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
          h='100px'
          minW='66px'
          position='relative'
          borderRadius='md'
        >
          <Text mb={2} zIndex={1} fontWeight='semibold' color='gray.500'>
            {issueEvent?.issueNumber}
          </Text>
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
          >
            <TagLabel>
              {issueEvent.status === 'Solved' ? 'Closed' : 'Open'}
            </TagLabel>
          </Tag>
          <IssueBgPattern position='absolute' width='120%' height='100%' />
        </Flex>
      </CardFooter>
    </Card>
  );
};
