import { useRef, useState, useEffect } from 'react';

import { produce } from 'immer';
// import { useDidMount } from 'rooks';
import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Reminder } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { useTimelineMeta } from '@organization/components/Timeline/state';
import { AutoresizeTextarea } from '@ui/form/Textarea/AutoresizeTextarea';

// import { ReminderEditForm } from './types';
import { ReminderPostit, ReminderDueDatePicker } from '../../shared';

interface ReminderItem {
  id: string;
  index: number;
  currentOwner: string;
}

export const ReminderItem = observer(
  ({ id, index, currentOwner }: ReminderItem) => {
    const store = useStore();

    const reminder = store.reminders.value.get(id);
    const owner =
      reminder?.value.owner?.name ||
      [reminder?.value.owner?.firstName, reminder?.value.owner?.lastName]
        .filter(Boolean)
        .join(' ');

    const ref = useRef<HTMLTextAreaElement>(null);
    const containerRef = useRef<HTMLDivElement>(null);
    const [timelineMeta, setTimelineMeta] = useTimelineMeta();
    const { recentlyCreatedId, recentlyUpdatedId } = timelineMeta.reminders;
    // const isMutating = data.id === 'TEMP';
    const [isFocused, setIsFocused] = useState(false);

    const handleContentChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
      reminder?.update((value) => {
        const property = e.target.name as keyof Reminder;

        value[property] = e.target.value;

        return value;
      });
    };

    const handleDueDateChange = (dateStr: string) => {
      reminder?.update((value) => {
        value.dueDate = dateStr;

        return value;
      });
    };

    const handleDismiss = () => {
      reminder?.update((value) => {
        value.dismissed = true;

        return value;
      });
    };

    // useDidMount(() => {
    //   if (['TEMP', recentlyCreatedId, recentlyUpdatedId].includes(data.id)) {
    //     ref.current?.focus();

    //     if (data.id === recentlyCreatedId) {
    //       setTimelineMeta((prev) =>
    //         produce(prev, (draft) => {
    //           draft.reminders.recentlyCreatedId = '';
    //           draft.reminders.recentlyUpdatedId = '';
    //         }),
    //       );
    //     }
    //   }
    // });

    useEffect(() => {
      if (
        id === recentlyUpdatedId ||
        id === 'TEMP' ||
        id === recentlyCreatedId
      ) {
        containerRef.current && containerRef.current.scrollIntoView();
      }
    }, [recentlyUpdatedId, id, index]);

    return (
      <ReminderPostit
        ref={containerRef}
        isFocused={isFocused}
        owner={owner === currentOwner ? undefined : owner}
        className={cn(
          id === recentlyUpdatedId ? 'shadow-ringWarning' : 'shadow-none',
        )}
        // isMutating={isMutating}
        onClickOutside={() => {
          setTimelineMeta((prev) =>
            produce(prev, (draft) => {
              draft.reminders.recentlyUpdatedId = '';
            }),
          );
        }}
      >
        <AutoresizeTextarea
          border
          ref={ref}
          name='content'
          cacheMeasurements
          // readOnly={isMutating}
          onChange={handleContentChange}
          onBlur={() => setIsFocused(false)}
          onFocus={() => setIsFocused(true)}
          maxRows={isFocused ? undefined : 3}
          data-test='timeline-reminder-editor'
          value={reminder?.value.content ?? ''}
          onKeyDown={(e) => e.stopPropagation()}
          placeholder='What should we remind you about?'
          className='px-2 pb-0 text-sm font-light font-sticky hover:border-transparent focus-within:border-transparent focus-within:hover:border-transparent'
        />
        <div className='flex items-center px-4 w-full justify-between mb-2'>
          <ReminderDueDatePicker
            onChange={handleDueDateChange}
            value={reminder?.value.dueDate}
          />

          <Button
            size='sm'
            variant='ghost'
            colorScheme='warning'
            onClick={handleDismiss}
            dataTest='timeline-reminder-dismiss'
            className='text-[#B7791F] hover:bg-transparent hover:text-warning-900 focus:shadow-ringWarning'
          >
            Dismiss
          </Button>
        </div>
      </ReminderPostit>
    );
  },
);
