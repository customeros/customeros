import React, { useEffect, useRef } from 'react';
import { RichTextEditor } from '@ui/form/RichTextEditor/RichTextEditor';
import { useRemirror } from '@remirror/react';
import { basicEditorExtensions } from '@ui/form/RichTextEditor/extensions';
import { Box, Flex } from '@chakra-ui/react';
import { Button } from '@ui/form/Button';
import { TagSuggestor } from './TagSuggestor';

import { TagsSelect } from './TagSelect';
import Image from 'next/image';
import noteIcon from 'public/images/event-ill-log.png';
import { useForm } from 'react-inverted-form';
import { LogEntryDto, LogEntryDtoI } from './LogEntry.dto';
import { useCreateLogEntryMutation } from '@organization/graphql/createLogEntry.generated';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useParams } from 'next/navigation';

import { useGetTagsQuery } from '@organization/graphql/getTags.generated';
import { invalidateAccountDetailsQuery } from '@organization/components/Tabs/panels/AccountPanel/utils';

export const Logger: React.FC<{ invalidateQuery: any }> = ({
  invalidateQuery,
}) => {
  const id = useParams()?.id as string;
  const client = getGraphQLClient();
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  const { data, isLoading } = useGetTagsQuery(client);
  const logEntryValues: LogEntryDtoI = new LogEntryDto();
  const { state, reset } = useForm<LogEntryDtoI>({
    formId: 'organization-create-log-entry',
    defaultValues: logEntryValues,

    stateReducer: (_, _a, next) => {
      return next;
    },
  });
  const remirrorProps = useRemirror({
    extensions: basicEditorExtensions,
  });
  const createLogEntryMutation = useCreateLogEntryMutation(client, {
    onSuccess: (data, variables, context) => {
      reset();
      timeoutRef.current = setTimeout(() => invalidateQuery(), 500);
    },
    onError: () => {
      console.log('🏷️ ----- : ERROR');
    },
  });

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  const onCreateLogEntry = () => {
    console.log('🏷️ ----- state.values.tags: ', state.values.tags);
    const logEntryPayload = LogEntryDto.toPayload({
      ...logEntryValues,
      tags: state.values.tags,
      content: state.values.content,
      contentType: state.values.contentType,
    });
    createLogEntryMutation.mutate({
      organizationId: id,

      logEntry: logEntryPayload,
    });
  };

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
          className='customeros-remirror-submit-button'
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
          isDisabled={createLogEntryMutation.isLoading}
          isLoading={createLogEntryMutation.isLoading}
          loadingText='Sending'
          onClick={onCreateLogEntry}
        >
          Log
        </Button>
      </Flex>
    </Flex>
  );
};
