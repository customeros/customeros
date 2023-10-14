import { CardHeader, CardBody } from '@ui/presentation/Card';
import { Heading } from '@ui/typography/Heading';
import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { Tooltip } from '@ui/presentation/Tooltip';
import { IconButton } from '@ui/form/IconButton';
import { useTimelineEventPreviewContext } from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';
import { Link03 } from '@ui/media/icons/Link03';
import { XClose } from '@ui/media/icons/XClose';
import copy from 'copy-to-clipboard';
import { User } from '@graphql/types';
import { Box } from '@ui/layout/Box';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useSession } from 'next-auth/react';
import { useGetTagsQuery } from '@organization/src/graphql/getTags.generated';
import React from 'react';
import { Tag, TagLabel } from '@ui/presentation/Tag';

const getAuthor = (user: User) => {
  if (!user?.firstName && !user?.lastName) {
    return 'Unknown';
  }

  return `${user.firstName} ${user.lastName}`.trim();
};

function getStatusColor(status: string) {
  if (['closed', 'solved'].includes(status.toLowerCase())) {
    return 'gray';
  }
  return 'blue';
}

export const IssuePreviewModal: React.FC = () => {
  const { closeModal, modalContent } = useTimelineEventPreviewContext();
  const { data: session } = useSession();
  const issue = modalContent as any;
  const client = getGraphQLClient();
  const { data } = useGetTagsQuery(client);
  const statusColorScheme = getStatusColor(issue.issueStatus);

  return (
    <>
      <CardHeader
        py='4'
        px='6'
        pb='1'
        position='sticky'
        top={0}
        borderRadius='xl'
      >
        <Flex
          direction='row'
          justifyContent='space-between'
          alignItems='center'
        >
          <Flex alignItems='center'>
            <Heading size='sm' fontSize='lg'>
              {issue?.subject ?? 'Issue'}
            </Heading>
          </Flex>
          <Flex direction='row' justifyContent='flex-end' alignItems='center'>
            <Tooltip label='Copy link' placement='bottom'>
              <IconButton
                variant='ghost'
                aria-label='Copy link to this issue'
                color='gray.500'
                fontSize='sm'
                size='sm'
                mr={1}
                icon={<Link03 color='gray.500' boxSize='4' />}
                onClick={() => copy(window.location.href)}
              />
            </Tooltip>
            <Tooltip label='Close' aria-label='close' placement='bottom'>
              <IconButton
                variant='ghost'
                aria-label='Close preview'
                color='gray.500'
                fontSize='sm'
                size='sm'
                icon={<XClose color='gray.500' boxSize='5' />}
                onClick={closeModal}
              />
            </Tooltip>
          </Flex>
        </Flex>
      </CardHeader>
      <CardBody
        mt={0}
        maxHeight='calc(100vh - 9rem)'
        p={6}
        pt={0}
        overflow='auto'
      >
        <Flex>
          <Box>todo priority</Box>
          <Tag
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
          >
            <TagLabel>{issue.status}</TagLabel>
          </Tag>
          <Tag
            size='sm'
            variant='outline'
            colorScheme='blue'
            border='1px solid'
            background='white'
            borderColor={`gray.200`}
            backgroundColor={`white`}
            color={`gray.500`}
            boxShadow='none'
            fontWeight='normal'
            minHeight={6}
          >
            <TagLabel>{issue?.externalLinks?.[0]?.externalId}</TagLabel>
          </Tag>
        </Flex>
        <Text fontSize='sm'>{issue?.description}</Text>

        {issue?.tags?.length && (
          <Flex>
            {issue.tags.map((tag: { id: string; name: string }) => (
              <Text key={`issue-tag-list-${tag.id}`} as='span' color='gray.500'>
                {tag.name}
              </Text>
            ))}
          </Flex>
        )}

        <Text>Issue requested by (todo user data) // todo date</Text>
        {/* todo */}

        {['solved', 'closed'].includes(issue.issueStatus?.toLowerCase) && (
          <Text>Issue closed by (todo user data) // todo date</Text>
        )}
      </CardBody>
    </>
  );
};
