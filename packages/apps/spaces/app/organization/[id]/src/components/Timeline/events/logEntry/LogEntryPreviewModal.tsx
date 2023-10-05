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
import { VStack } from '@ui/layout/Stack';
import { LogEntryWithAliases } from '@organization/src/components/Timeline/types';
import { User } from '@graphql/types';
import { Box } from '@ui/layout/Box';
import noteImg from 'public/images/note-img-preview.png';
import { LogEntryDatePicker } from './preview/LogEntryDatePicker';
import { Image } from '@ui/media/Image';
import { LogEntryExternalLink } from './preview/LogEntryExternalLink';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useSession } from 'next-auth/react';
import { PreviewTags } from './preview/tags/PreviewTags';
import { PreviewEditor } from './preview/PreviewEditor';
import { useGetTagsQuery } from '@organization/src/graphql/getTags.generated';
import { useLogEntryUpdateContext } from '@organization/src/components/Timeline/events/logEntry/context/LogEntryUpdateModalContext';
import { HtmlContentRenderer } from '@ui/presentation/HtmlContentRenderer/HtmlContentRenderer';
import React from 'react';

const getAuthor = (user: User) => {
  if (!user?.firstName && !user.lastName) {
    return 'Unknown';
  }

  return `${user.firstName} ${user.lastName}`.trim();
};

export const LogEntryPreviewModal: React.FC = () => {
  const { closeModal, modalContent } = useTimelineEventPreviewContext();
  const { data: session } = useSession();
  const event = modalContent as LogEntryWithAliases;
  const author = getAuthor(event?.logEntryCreatedBy);
  const authorEmail = event?.logEntryCreatedBy?.emails?.[0]?.email;
  const client = getGraphQLClient();
  const { data } = useGetTagsQuery(client);

  const isAuthor =
    event.logEntryCreatedBy?.emails?.findIndex(
      (e) => session?.user?.email === e.email,
    ) !== -1;
  const { formId } = useLogEntryUpdateContext();

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
              Log entry
            </Heading>
          </Flex>
          <Flex direction='row' justifyContent='flex-end' alignItems='center'>
            <Tooltip label='Copy link' placement='bottom'>
              <IconButton
                variant='ghost'
                aria-label='Copy link to this entry'
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
        <Box position='relative'>
          <Image
            src={noteImg}
            alt=''
            height={123}
            width={174}
            position='absolute'
            top={-2}
            right={-3}
          />
        </Box>
        <VStack gap={2} alignItems='flex-start'>
          <Flex direction='column'>
            <LogEntryDatePicker event={event} formId={formId} />
          </Flex>
          <Flex direction='column'>
            <Text fontSize='sm' fontWeight='semibold'>
              Author
            </Text>
            <Tooltip label={authorEmail} hasArrow>
              <Text fontSize='sm'>{author}</Text>
            </Tooltip>
          </Flex>

          <Flex direction='column' w='full'>
            <Text fontSize='sm' fontWeight='semibold'>
              Entry
            </Text>

            {!isAuthor && (
              <HtmlContentRenderer
                fontSize='sm'
                noOfLines={undefined}
                htmlContent={`${event?.content}`}
              />
            )}
            {isAuthor && (
              <PreviewEditor
                formId={formId}
                initialContent={`${event?.content}`}
                tags={data?.tags}
                onClose={closeModal}
              />
            )}
          </Flex>

          <PreviewTags
            isAuthor={isAuthor}
            tags={event?.tags}
            tagOptions={data?.tags}
            id={event.id}
          />

          {event?.externalLinks?.[0]?.externalUrl && (
            <LogEntryExternalLink externalLink={event?.externalLinks?.[0]} />
          )}
        </VStack>
      </CardBody>
    </>
  );
};
