import { useParams } from 'next/navigation';
import { FormEvent, useEffect } from 'react';
import { useForm } from 'react-inverted-form';

import { produce } from 'immer';
import { useDebounceFn } from 'rooks';
import { useQueryClient } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { VStack } from '@ui/layout/Stack';
import { useModKey } from '@shared/hooks/useModKey';
import { toastError } from '@ui/presentation/Toast';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { useUpdateReminderMutation } from '@organization/src/graphql/updateReminder.generated';
import {
  RemindersQuery,
  useRemindersQuery,
} from '@organization/src/graphql/reminders.generated';

import { ReminderPostit, ReminderDueDatePicker } from '../../shared';

type ReminderEditForm = {
  id: string;
  date: string;
  owner: string;
  content: string;
};

export const Reminders = () => {
  const organizationId = useParams()?.id as string;
  const client = getGraphQLClient();
  const queryClient = useQueryClient();

  const { data, isPending } = useRemindersQuery(client, { organizationId });
  const remindersQueryKey = useRemindersQuery.getKey({ organizationId });

  const updateReminder = useUpdateReminderMutation(client, {
    onMutate: (values) => {
      queryClient.cancelQueries({ queryKey: remindersQueryKey });

      const previousEntries = useRemindersQuery.mutateCacheEntry(queryClient, {
        organizationId,
      })((cache) =>
        produce(cache, (draft) => {
          if (!draft) return;

          const foundReminder = draft.remindersForOrganization.find(
            (r) => r.metadata.id === values.input.id,
          );

          if (!foundReminder) return;
          foundReminder.content = values.input.content ?? '';
          foundReminder.dismissed = values.input.dismissed ?? false;
        }),
      );

      return { previousEntries };
    },
    onError: (_, __, context) => {
      if (context?.previousEntries) {
        queryClient.setQueryData(remindersQueryKey, context.previousEntries);
      }
      toastError(`We couldn't update the reminder`, 'update-reminder-error');
    },
    onSettled: () => {
      setTimeout(() => {
        queryClient.invalidateQueries({ queryKey: remindersQueryKey });
      }, 500);
    },
  });

  const onChange = (values: ReminderEditForm) => {
    updateReminder.mutate({
      input: {
        id: values.id,
        dueDate: values.date,
        content: values.content,
      },
    });
  };

  const onDismiss = (id: string) => {
    updateReminder.mutate({
      input: {
        id: id,
        dismissed: true,
      },
    });
  };

  const { data: globalCacheData } = useGlobalCacheQuery(client);

  const user = globalCacheData?.global_Cache?.user;
  const currentOwner = [user?.firstName, user?.lastName]
    .filter(Boolean)
    .join(' ');

  if (isPending) return <p>Loading...</p>;

  return (
    <VStack align='flex-start'>
      {data?.remindersForOrganization
        ?.filter((r) => !r.dismissed)
        ?.map((r) => (
          <ReminderItem
            key={r.metadata.id}
            currentOwner={currentOwner}
            onChange={onChange}
            onDismiss={onDismiss}
            data={mapReminderToForm(r)}
          />
        ))}
    </VStack>
  );
};

interface ReminderItem {
  currentOwner: string;
  data: ReminderEditForm;
  onDismiss: (id: string) => void;
  onChange: (value: ReminderEditForm) => void;
}

const ReminderItem = ({
  data,
  onChange,
  onDismiss,
  currentOwner,
}: ReminderItem) => {
  const formId = `reminder-edit-form-${data.id}`;
  const [debouncedOnChange] = useDebounceFn(
    (arg) => onChange(arg as ReminderEditForm),
    1000,
  );

  const makeContentStr = (content: string, owner: string) => {
    const strippedContent = content.replace(`for ${owner}: `, '');

    return currentOwner !== owner
      ? `for ${owner}: ${strippedContent}`
      : strippedContent;
  };

  const { handleSubmit, setDefaultValues } = useForm<ReminderEditForm>({
    formId,
    defaultValues: data,
    onSubmit: async (values) => {
      onChange(values);
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
        debouncedOnChange(next.values);
      }

      return next;
    },
  });

  const updateReminder = () => {
    handleSubmit({} as FormEvent<HTMLFormElement>);
  };

  useModKey('Enter', updateReminder);

  useEffect(() => {
    setDefaultValues({
      ...data,
      content: makeContentStr(data.content, data.owner),
    });
  }, [currentOwner]);

  return (
    <ReminderPostit>
      <FormAutoresizeTextarea
        px='4'
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

function mapReminderToForm(
  reminder?: RemindersQuery['remindersForOrganization'][number],
): ReminderEditForm {
  return {
    id: reminder?.metadata.id ?? '',
    date: reminder?.dueDate,
    content: reminder?.content ?? '',
    owner:
      reminder?.owner?.name ||
      [reminder?.owner?.firstName, reminder?.owner?.lastName]
        .filter(Boolean)
        .join(' '),
  };
}
