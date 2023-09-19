import React, { useEffect, useRef } from 'react';
import { RichTextEditor } from '@ui/form/RichTextEditor/RichTextEditor';
import { Box, Flex } from '@chakra-ui/react';
import { Button } from '@ui/form/Button';
import { TagSuggestor } from './TagSuggestor';

import { TagsSelect } from './TagSelect';
import Image from 'next/image';
import noteIcon from '../../../../../../../../public/images/event-ill-log.png';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import { useGetTagsQuery } from '@organization/graphql/getTags.generated';
import { useTimelineActionLogEntryContext } from '@organization/components/Timeline/TimelineActions/TimelineActionsContext/TimelineActionLogEntryContext';

export const Logger = () => {
  const { onCreateLogEntry, remirrorProps, isSaving } =
    useTimelineActionLogEntryContext();
  const client = getGraphQLClient();
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const { data } = useGetTagsQuery(client);
  //
  // const createLogEntryMutation = useCreateLogEntryMutation(client, {
  //   onSuccess: () => {
  //     // reset();
  //     timeoutRef.current = setTimeout(() => invalidateQuery(), 500);
  //   },
  // });

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  return (
    <Flex
      flexDirection='column'
      position='relative'
      className='customeros-logger'
    >
      <Box position='absolute' top={-6} right={-6}>
        <Image src={noteIcon} alt='' height={123} width={174} />
      </Box>

      <RichTextEditor
        {...remirrorProps}
        placeholder='Log conversation you had with a customer'
        formId='organization-create-log-entry'
        name='content'
        showToolbar={false}
      >
        <TagSuggestor
          tags={data?.tags?.map((e: { label: string; value: string }) => ({
            label: e.label,
            id: e.value,
          }))}
        />
      </RichTextEditor>
      <Flex justifyContent='space-between' zIndex={3}>
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
          isDisabled={isSaving}
          isLoading={isSaving}
          loadingText='Sending'
          onClick={onCreateLogEntry}
        >
          Log
        </Button>
      </Flex>
    </Flex>
  );
};
