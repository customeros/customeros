import { useParams } from 'next/navigation';

import { produce } from 'immer';
import { useQueryClient } from '@tanstack/react-query';

import { VStack } from '@ui/layout/Stack';
import { toastError } from '@ui/presentation/Toast';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { useTimelineMeta } from '@organization/src/components/Timeline/state';
import { useUpdateReminderMutation } from '@organization/src/graphql/updateReminder.generated';
import {
  RemindersQuery,
  useRemindersQuery,
} from '@organization/src/graphql/reminders.generated';

import { ReminderEditForm } from './types';
import { ReminderItem } from './ReminderItem';

export const Reminders = () => {
  const organizationId = useParams()?.id as string;
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const [_, setTimelineMeta] = useTimelineMeta();

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

          setTimelineMeta((prev) =>
            produce(prev, (draft) => {
              draft.reminders.recentlyUpdatedId = values.input.id ?? '';
            }),
          );
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

  if (isPending) return null;

  return (
    <VStack align='flex-start'>
      {data?.remindersForOrganization
        ?.filter((r) => !r.dismissed)
        .sort(
          (a, b) =>
            new Date(a?.dueDate).valueOf() - new Date(b?.dueDate).valueOf(),
        )
        .map((r) => (
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
