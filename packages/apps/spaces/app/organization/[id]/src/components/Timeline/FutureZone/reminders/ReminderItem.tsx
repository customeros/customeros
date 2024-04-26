import { useForm } from 'react-inverted-form';
import { useRef, useState, FormEvent, useEffect, useCallback } from 'react';

import { produce } from 'immer';
import isEqual from 'lodash/isEqual';
import { useDidMount, useDebounceFn } from 'rooks';

import { cn } from '@ui/utils/cn';
import { Button } from '@ui/form/Button/Button';
import { useTimelineMeta } from '@organization/src/components/Timeline/state';
import { FormAutoresizeTextarea } from '@ui/form/Textarea/FormAutoresizeTextarea';

import { ReminderEditForm } from './types';
import { ReminderPostit, ReminderDueDatePicker } from '../../shared';

interface ReminderItem {
  index: number;
  currentOwner: string;
  data: ReminderEditForm;
  onDismiss: (id: string) => void;
  onChange: (value: ReminderEditForm) => void;
}

export const ReminderItem = ({
  data,
  index,
  onChange,
  onDismiss,
  currentOwner,
}: ReminderItem) => {
  const ref = useRef<HTMLTextAreaElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const formId = `reminder-edit-form-${data.id}`;
  const [timelineMeta, setTimelineMeta] = useTimelineMeta();
  const [debouncedOnChange] = useDebounceFn(
    (arg) => onChange(arg as ReminderEditForm),
    1000,
  );
  const { recentlyCreatedId, recentlyUpdatedId } = timelineMeta.reminders;
  const isMutating = data.id === 'TEMP';
  const [isFocused, setIsFocused] = useState(false);

  const onSubmit = useCallback(
    async (values: ReminderEditForm) =>
      !isEqual(values, data) ? onChange(values) : undefined,
    [onChange],
  );

  const { handleSubmit, setDefaultValues } = useForm<ReminderEditForm>({
    formId,
    defaultValues: data,
    onSubmit,
    stateReducer: (_, action, next) => {
      if (isMutating) return next;

      if (action.type === 'FIELD_CHANGE') {
        switch (action.payload.name) {
          case 'date': {
            onChange(next.values);
            break;
          }
          default: {
            debouncedOnChange(next.values);
            break;
          }
        }
      }

      return next;
    },
  });

  const updateReminder = () => {
    setIsFocused(false);
    handleSubmit({} as FormEvent<HTMLFormElement>);
  };

  useEffect(() => {
    setDefaultValues(data);
  }, [currentOwner, data.id]);

  useDidMount(() => {
    if (['TEMP', recentlyCreatedId, recentlyUpdatedId].includes(data.id)) {
      ref.current?.focus();

      if (data.id === recentlyCreatedId) {
        setTimelineMeta((prev) =>
          produce(prev, (draft) => {
            draft.reminders.recentlyCreatedId = '';
            draft.reminders.recentlyUpdatedId = '';
          }),
        );
      }
    }
  });

  useEffect(() => {
    if (
      data.id === recentlyUpdatedId ||
      data.id === 'TEMP' ||
      data.id === recentlyCreatedId
    ) {
      containerRef.current && containerRef.current.scrollIntoView();
    }
  }, [recentlyUpdatedId, data.id, index]);

  return (
    <ReminderPostit
      ref={containerRef}
      className={cn(
        data.id === recentlyUpdatedId ? 'shadow-ringWarning' : 'shadow-none',
      )}
      owner={data?.owner === currentOwner ? undefined : data?.owner}
      isFocused={isFocused}
      isMutating={isMutating}
      onClickOutside={() => {
        setTimelineMeta((prev) =>
          produce(prev, (draft) => {
            draft.reminders.recentlyUpdatedId = '';
          }),
        );
      }}
    >
      <FormAutoresizeTextarea
        className='px-2 pb-0 text-sm font-light font-sticky hover:border-transparent focus-within:border-transparent focus-within:hover:border-transparent'
        ref={ref}
        border
        readOnly={isMutating}
        name='content'
        formId={formId}
        onBlur={updateReminder}
        onFocus={() => setIsFocused(true)}
        cacheMeasurements
        maxRows={isFocused ? undefined : 3}
        placeholder='What should we remind you about?'
      />
      <div className='flex items-center px-4 w-full justify-between mb-2'>
        <ReminderDueDatePicker name='date' formId={formId} />

        <Button
          size='sm'
          variant='ghost'
          colorScheme='warning'
          className='text-[#B7791F] hover:bg-transparent hover:text-warning-900 focus:shadow-ringWarning'
          onClick={() => onDismiss(data.id)}
        >
          Dismiss
        </Button>
      </div>
    </ReminderPostit>
  );
};
