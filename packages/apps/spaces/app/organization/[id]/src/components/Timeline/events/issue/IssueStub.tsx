'use client';
import React, { FC } from 'react';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { CardBody, CardHeader, CardFooter, Card } from '@ui/layout/Card';
import { IssueBgPattern } from '@ui/media/logos/IssueBgPattern';
import { Tag, TagLabel } from '@ui/presentation/Tag';
import { CustomTicketTearStyle } from './styles';
import { IssueWithAliases } from '@organization/src/components/Timeline/types';
import { getExternalUrl } from '@spaces/utils/getExternalLink';
import { toastError } from '@ui/presentation/Toast';
function getStatusColor(
  status: string | 'New' | 'Open' | 'Pending' | 'On hold' | 'Solved',
) {
  if (status === 'Solved') {
    return 'gray';
  }
  return 'blue';
}

export const IssueStub: FC<{ data: IssueWithAliases }> = ({ data }) => {
  // const { openModal } = useTimelineEventPreviewContext(); // todo uncomment when modal is ready
  const statusColorScheme = (() => getStatusColor(data.issueStatus))();
  const handleOpenInExternalApp = () => {
    if (data?.externalLinks?.[0]?.externalUrl) {
      window.open(
        getExternalUrl(data.externalLinks[0].externalUrl),
        '_blank',
        'noreferrer noopener',
      );
      return;
    }
    toastError(
      'This issue is not connected to external source',
      `${data.id}-stub-open-in-external-app-error`,
    );
  };

  return (
    <Card
      variant='outline'
      size='md'
      fontSize='14px'
      background='white'
      flexDirection='row'
      position='unset'
      maxW={476}
      cursor='default' // todo change to pointer when modal is ready
      boxShadow='none'
      border='1px solid'
      borderColor='gray.200'
      onClick={handleOpenInExternalApp} // todo remove when COS-464 is merged
      // onClick={() => openModal(data)}
      // TODO uncomment when modal is ready
      // _hover={{
      //   '&:hover .slack-stub-date': {
      //     color: 'gray.500',
      //   },
      // }}
    >
      <Flex boxShadow='xs' pr={2} direction='column' flex={1}>
        <CardHeader fontWeight='semibold' p={2} pb={0} pr={0} noOfLines={1}>
          {data?.subject ?? '[No subject]'}
        </CardHeader>
        <CardBody p={2} pt={0} pr={0}>
          <Text color='gray.500' noOfLines={3}>
            {data?.description ?? '[No description]'}
          </Text>
        </CardBody>
      </Flex>
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
          {!!data?.externalLinks?.length && (
            <Text mb={2} zIndex={1} fontWeight='semibold' color='gray.500'>
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
          >
            <TagLabel>{data.status === 'Solved' ? 'Closed' : 'Open'}</TagLabel>
          </Tag>
          <IssueBgPattern position='absolute' width='120%' height='100%' />
        </Flex>
      </CardFooter>
    </Card>
  );
};
