import { useRef, useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { $getRoot, LexicalEditor } from 'lexical';
import noteIcon from '@assets/images/event-ill-log.png';

import { Button } from '@ui/form/Button/Button';
import { Editor } from '@ui/form/Editor/Editor';
import { useStore } from '@shared/hooks/useStore';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';
import { useTimelineActionLogEntryContext } from '@organization/components/Timeline/FutureZone/TimelineActions/context/TimelineActionLogEntryContext';

interface LoggerProps {
  hide: () => void;
}

export const Logger = observer(({ hide }: LoggerProps) => {
  const store = useStore();

  const [value, setValue] = useState<string>('');
  const editorRef = useRef<LexicalEditor | null>(null);
  const [hashags, setHashtags] = useState<string[]>([]);
  const [hashtagSearch, setHashtagSearch] = useState<string | null>(null);
  const [mentionsSearch, setMentionsSearch] = useState<string | null>(null);

  const { onCreateLogEntry } = useTimelineActionLogEntryContext();

  const handleChange = (html: string) => {
    setValue(html);

    if (html === '<p><br></p>') {
      store.ui.clearDirtyEditor();
    } else {
      store.ui.setDirtyEditor('log-entry');
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
    editorRef?.current?.update(() => {
      const root = $getRoot();

      root.clear();
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
    store.ui.clearDirtyEditor();
    hide();
  };

  const hashtags = store.tags
    .toArray()
    .map((t) => ({ value: t.value.id, label: t.value.name }))
    .filter((t) =>
      hashtagSearch
        ? t.label.toLowerCase().includes(hashtagSearch?.toLowerCase())
        : true,
    );

  const mentions = store.users
    .toArray()
    .map(
      ({ value: { name, lastName, firstName } }) =>
        name || [firstName, lastName].filter(Boolean).join(' '),
    )
    .filter((m) =>
      mentionsSearch
        ? m.toLowerCase().includes(mentionsSearch?.toLowerCase())
        : true,
    )
    .filter(Boolean) as string[];

  useEffect(() => {
    return () => {
      store.ui.clearDirtyEditor();
    };
  }, []);

  return (
    <div className='customeros-logger flex flex-col min-h-[123px] relative'>
      <div className='absolute top-[-16px] right-[-24px]'>
        <img alt='' width={174} height={135} src={noteIcon} />
      </div>

      <div className='z-2 w-full h-full'>
        <Editor
          className='mb-10'
          onChange={handleChange}
          mentionsOptions={mentions}
          hashtagsOptions={hashtags}
          namespace='LogEntryCreator'
          dataTest='timeline-log-editor'
          onHashtagSearch={setHashtagSearch}
          onMentionsSearch={setMentionsSearch}
          placeholder='Log a conversation you had with a customer'
          onHashtagsChange={(hashtags) =>
            setHashtags(hashtags.map((h) => h.label))
          }
        ></Editor>
        <Button
          size='xs'
          variant='outline'
          onClick={handleSave}
          className='absolute bottom-0 right-0'
          dataTest='timeline-log-confirmation-button'
        >
          Log
        </Button>
      </div>
      <ConfirmDeleteDialog
        colorScheme='primary'
        onConfirm={handleSave}
        onClose={handleDiscard}
        label='Log this log entry?'
        confirmButtonLabel='Log it'
        cancelButtonLabel='Discard'
        isOpen={store.ui.activeConfirmation === 'log-entry'}
        description='You have typed an unlogged entry. Do you want to log it to the timeline, or discard it?'
      />
    </div>
  );
});
