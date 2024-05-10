import { useParams } from 'react-router-dom';
import { useField } from 'react-inverted-form';

import noteIcon from '@assets/images/event-ill-log.png';

import { Contact } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { getMentionOptionLabel } from '@organization/hooks/utils';
import { RichTextEditor } from '@ui/form/RichTextEditor/RichTextEditor';
import { useGetTagsQuery } from '@organization/graphql/getTags.generated';
import { useGetMentionOptionsQuery } from '@organization/graphql/getMentionOptions.generated';
import { FloatingReferenceSuggestions } from '@ui/form/RichTextEditor/FloatingReferenceSuggestions';
import { KeymapperClose } from '@ui/form/RichTextEditor/components/keyboardShortcuts/KeymapperClose';
import { KeymapperCreate } from '@ui/form/RichTextEditor/components/keyboardShortcuts/KeymapperCreate';
import { useTimelineActionContext } from '@organization/components/Timeline/FutureZone/TimelineActions/context/TimelineActionContext';
import { useTimelineActionLogEntryContext } from '@organization/components/Timeline/FutureZone/TimelineActions/context/TimelineActionLogEntryContext';

import { TagsSelect } from './TagSelect';

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
    .filter((e) => Boolean(e.label)) as { id: string; label: string }[];

  return (
    <div className='customeros-logger flex flex-col min-h-[123px] relative'>
      <div className='absolute top-[-16px] right-[-24px]'>
        <img src={noteIcon} alt='' height={135} width={174} />
      </div>

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
      <div className='flex justify-between text-base'>
        <TagsSelect
          formId='organization-create-log-entry'
          name='tags'
          tags={data?.tags}
        />
        <Button
          className='font-semibold rounded-lg py-1 px-3 text-sm items-center'
          variant='outline'
          colorScheme='gray'
          size='xs'
          isDisabled={isSaving || isLogEmpty}
          isLoading={isSaving}
          loadingText='Sending'
          onClick={() => onCreateLogEntry()}
        >
          Log
        </Button>
      </div>
    </div>
  );
};
