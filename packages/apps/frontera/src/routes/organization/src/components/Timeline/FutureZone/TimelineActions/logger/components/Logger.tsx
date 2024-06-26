import { useState } from 'react';
import { useParams } from 'react-router-dom';

import { $getRoot } from 'lexical';
import { observer } from 'mobx-react-lite';
import noteIcon from '@assets/images/event-ill-log.png';
import { useLexicalComposerContext } from '@lexical/react/LexicalComposerContext';

import { Contact } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { Editor } from '@ui/form/Editor/Editor';
import { useStore } from '@shared/hooks/useStore';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { getMentionOptionLabel } from '@organization/hooks/utils';
import { MessageChatSquare } from '@ui/media/icons/MessageChatSquare';
import { useGetTagsQuery } from '@organization/graphql/getTags.generated';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';
import { useGetMentionOptionsQuery } from '@organization/graphql/getMentionOptions.generated';
import { useTimelineActionLogEntryContext } from '@organization/components/Timeline/FutureZone/TimelineActions/context/TimelineActionLogEntryContext';

interface LoggerProps {
  hide: () => void;
}

export const Logger = observer(({ hide }: LoggerProps) => {
  const store = useStore();
  const id = useParams()?.id as string;

  const [hashags, setHashtags] = useState<string[]>([]);
  const [value, setValue] = useState<string>('');
  const [hashtagSearch, setHashtagSearch] = useState<string | null>(null);
  const { onCreateLogEntry } = useTimelineActionLogEntryContext();

  const client = getGraphQLClient();
  const { data } = useGetTagsQuery(client);
  const { data: mentionData } = useGetMentionOptionsQuery(client, {
    id,
  });

  const handleChange = (html: string) => {
    setValue(html);
    if (html === '<p><br></p>') {
      store.ui.setDirtyEditor('log-entry');
    } else {
      store.ui.clearDirtyEditor();
    }
  };

  const handleSave = () => {
    // remove this code in order to switch the store logic on.
    onCreateLogEntry({
      payload: {
        content: value,
        tags: hashags.map((t) => ({ name: t })),
        contentType: 'text/html',
      },
      onSuccess: () => {
        store.ui.clearConfirmAction();
        hide();
      },
    });

    // uncomment the bellow code in order to switch the store logic on.
    // store.timelineEvents.logEntries.create(id, {
    //   content: value,
    //   tags: hashags,
    // });

    // const canClose = !value || value === '<p><br></p>';

    // if (canClose) {
    //   showEditor(null);
    // }
  };

  const handleDiscard = () => {
    store.ui.clearConfirmAction();
    hide();
  };

  const mentionOptions = (mentionData?.organization?.contacts?.content ?? [])
    .map((e) => getMentionOptionLabel(e as Contact))
    .filter(Boolean) as string[];
  const hashtagsOptions =
    data?.tags.filter((t) => t.label.includes(hashtagSearch ?? '')) ?? [];

  return (
    <div className='customeros-logger flex flex-col min-h-[123px] relative'>
      <div className='absolute top-[-16px] right-[-24px]'>
        <img src={noteIcon} alt='' height={135} width={174} />
      </div>

      <div className='z-2 w-full h-full'>
        <Editor
          className='mb-10'
          onChange={handleChange}
          mentionsOptions={mentionOptions}
          hashtagsOptions={hashtagsOptions}
          onHashtagSearch={setHashtagSearch}
          onHashtagCreate={(hashtag) => console.info(hashtag)}
          placeholder='Log a conversation you had with a customer'
          onHashtagsChange={(hashtags) =>
            setHashtags(hashtags.map((h) => h.label))
          }
        >
          <SaveButton onClick={handleSave} />
        </Editor>
      </div>
      <ConfirmDeleteDialog
        colorScheme='primary'
        onConfirm={handleSave}
        onClose={handleDiscard}
        label='Log this log entry?'
        confirmButtonLabel='Log it'
        cancelButtonLabel='Discard'
        isOpen={store.ui.activeConfirmation === 'log-entry'}
        icon={<MessageChatSquare className='text-primary-700' />}
        description='You have typed an unlogged entry. Do you want to log it to the timeline, or discard it?'
      />
    </div>
  );
});

const SaveButton = ({ onClick }: { onClick?: () => void }) => {
  const [editor] = useLexicalComposerContext();

  return (
    <Button
      size='xs'
      onClick={() => {
        onClick?.();
        editor.update(() => {
          const root = $getRoot();
          root.clear();
        });
      }}
      variant='outline'
      className='absolute bottom-0 right-0'
    >
      Log
    </Button>
  );
};
