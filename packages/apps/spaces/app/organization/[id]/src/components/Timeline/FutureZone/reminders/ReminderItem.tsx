import { useForm } from 'react-inverted-form';
import { useRef, FormEvent, useEffect } from 'react';

import { produce } from 'immer';
import { useDidMount, useDebounceFn } from 'rooks';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { useTimelineMeta } from '@organization/src/components/Timeline/state';

import { ReminderEditForm } from './types';
import { ReminderPostit, ReminderDueDatePicker } from '../../shared';

interface ReminderItem {
  currentOwner: string;
  data: ReminderEditForm;
  onDismiss: (id: string) => void;
  onChange: (value: ReminderEditForm) => void;
}

export const ReminderItem = ({
  data,
  onChange,
  onDismiss,
  currentOwner,
}: ReminderItem) => {
  const ref = useRef<HTMLTextAreaElement>(null);
  const formId = `reminder-edit-form-${data.id}`;
  const [timelineMeta, setTimelineMeta] = useTimelineMeta();
  const [debouncedOnChange] = useDebounceFn(
    (arg) => onChange(arg as ReminderEditForm),
    1000,
  );
  const { recentlyCreatedId, recentlyUpdatedId } = timelineMeta.reminders;

  const stripContent = (content: string, owner: string) => {
    const targetString = `for ${owner}: `;

    if (!content.startsWith(targetString)) return content;

    return content.replace(targetString, '');
  };
  const makeContentStr = (content: string, owner: string) => {
    const strippedContent = stripContent(content, owner);

    return currentOwner !== owner
      ? `for ${owner}: ${strippedContent}`
      : strippedContent;
  };

  const { handleSubmit, setDefaultValues } = useForm<ReminderEditForm>({
    formId,
    defaultValues: data,
    onSubmit: async (values) => {
      onChange({
        ...values,
        content: stripContent(values.content, values.owner),
      });
    },
    stateReducer: (_, action, next) => {
      if (action.type === 'FIELD_CHANGE') {
        if (action.payload.name === 'content') {
          return {
            ...next,
            values: {
              ...next.values,
              content: makeContentStr(next.values.content, next.values.owner),
            },
          };
        }
        debouncedOnChange({
          ...next.values,
          content: stripContent(next.values.content, next.values.owner),
        });
      }

      return next;
    },
  });

  const updateReminder = () => {
    handleSubmit({} as FormEvent<HTMLFormElement>);
  };

  useEffect(() => {
    setDefaultValues({
      ...data,
      content: makeContentStr(data.content, data.owner),
    });
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

  return (
    <ReminderPostit
      boxShadow={data.id === recentlyUpdatedId ? 'ringPrimary' : 'unset'}
      onClickOutside={() => {
        setTimelineMeta((prev) =>
          produce(prev, (draft) => {
            draft.reminders.recentlyUpdatedId = '';
          }),
        );
      }}
    >
      <FormAutoresizeTextarea
        px='4'
        ref={ref}
        fontFamily='sticky'
        fontSize='sm'
        name='content'
        formId={formId}
        onBlur={updateReminder}
        placeholder='Type your reminder here'
        borderBottom='unset'
        _hover={{
          borderBottom: 'unset',
        }}
        _focus={{
          borderBottom: 'unset',
        }}
      />
      <Flex align='center' px='4' w='full' justify='space-between' mb='2'>
        <ReminderDueDatePicker name='date' formId={formId} />

        <Button
          size='sm'
          variant='ghost'
          colorScheme='yellow'
          onClick={() => onDismiss(data.id)}
        >
          Dismiss
        </Button>
      </Flex>
    </ReminderPostit>
  );
};
