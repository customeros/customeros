import React from 'react';
import Image from 'next/image';
import { useField } from 'react-inverted-form';

import { Button } from '@ui/form/Button';
import { Box, Flex } from '@chakra-ui/react';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { FloatingReferenceSuggestions } from '@ui/form/RichTextEditor/FloatingReferenceSuggestions';
import { RichTextEditor } from '@ui/form/RichTextEditor/RichTextEditor';
import { useGetTagsQuery } from '@organization/src/graphql/getTags.generated';
import { useTimelineActionLogEntryContext } from '@organization/src/components/Timeline/TimelineActions/context/TimelineActionLogEntryContext';

import { TagsSelect } from './TagSelect';
import noteIcon from 'public/images/event-ill-log.png';
import { useGetMentionOptionsQuery } from '@organization/src/graphql/getMentionOptions.generated';
import { useParams } from 'next/navigation';
import { getMentionOptionLabel } from '@organization/src/components/Timeline/events/utils';
import { Contact } from '@graphql/types';
import { KeymapperCreate } from '@ui/form/RichTextEditor/components/keyboardShortcuts/KeymapperCreate';
import { KeymapperClose } from '@ui/form/RichTextEditor/components/keyboardShortcuts/KeymapperClose';
import { useTimelineActionContext } from '@organization/src/components/Timeline/TimelineActions/context/TimelineActionContext';

export const Logger = () => {
  const id = useParams()?.id as string;
  const { onCreateLogEntry, remirrorProps, isSaving, checkCanExitSafely } =
    useTimelineActionLogEntryContext();

  const client = getGraphQLClient();
  const { getInputProps } = useField(
    'content',
    'organization-create-log-entry',
  );
  const { value } = getInputProps();
  const { data } = useGetTagsQuery(client);
  const { data: mentionData } = useGetMentionOptionsQuery(client, {
    id,
  });
  const { showEditor } = useTimelineActionContext();

  const handleClose = () => {
    const canClose = checkCanExitSafely();

    if (canClose) {
      showEditor(null);
    }
  };
  const isLogEmpty = !value?.length || value === `<p style=""></p>`;

  const mentionOptions = (mentionData?.organization?.contacts?.content ?? [])
    .map((e) => ({ label: getMentionOptionLabel(e as Contact), id: e.id }))
    .filter((e) => Boolean(e.label)) as { label: string; id: string }[];

  return (
    <Flex
      flexDirection='column'
      position='relative'
      className='customeros-logger'
    >
      <Box position='absolute' top={-4} right={-6}>
        <Image src={noteIcon} alt='' height={123} width={174} />
      </Box>

      <RichTextEditor
        {...remirrorProps}
        placeholder='Log a conversation you had with a customer'
        formId='organization-create-log-entry'
        name='content'
        showToolbar={false}
      >
        <FloatingReferenceSuggestions
          tags={data?.tags?.map((e: { label: string; value: string }) => ({
            label: e.label,
            id: e.value,
          }))}
          mentionOptions={mentionOptions}
        />
        <KeymapperCreate onCreate={onCreateLogEntry} />
        <KeymapperClose onClose={handleClose} />
      </RichTextEditor>
      <Flex justifyContent='space-between' zIndex={8} fontSize='md'>
        <TagsSelect
          formId='organization-create-log-entry'
          name='tags'
          tags={data?.tags}
        />
        <Button
          variant='outline'
          colorScheme='gray'
          fontWeight={600}
          borderRadius='lg'
          pt={1}
          pb={1}
          pl={3}
          pr={3}
          size='sm'
          fontSize='sm'
          isDisabled={isSaving || isLogEmpty}
          isLoading={isSaving}
          loadingText='Sending'
          onClick={() => onCreateLogEntry()}
        >
          Log
        </Button>
      </Flex>
    </Flex>
  );
};
