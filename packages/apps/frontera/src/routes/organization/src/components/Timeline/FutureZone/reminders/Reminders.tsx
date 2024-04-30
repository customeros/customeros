import { useEffect } from 'react';
import { useParams } from 'react-router-dom';

import { produce } from 'immer';
import { useQueryClient } from '@tanstack/react-query';

import { toastError } from '@ui/presentation/Toast';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useTimelineMeta } from '@organization/components/Timeline/state';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { useUpdateReminderMutation } from '@organization/graphql/updateReminder.generated';
import {
  RemindersQuery,
  useRemindersQuery,
} from '@organization/graphql/reminders.generated';

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
          foundReminder.dueDate = values.input.dueDate ?? '';
          foundReminder.metadata.lastUpdated = new Date().toISOString();

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

  const remindersLength = data?.remindersForOrganization?.length ?? 0;
  const user = globalCacheData?.global_Cache?.user;
  const currentOwner = [user?.firstName, user?.lastName]
    .filter(Boolean)
    .join(' ');

  useEffect(() => {
    setTimelineMeta((prev) => ({
      ...prev,
      remindersCount: remindersLength,
    }));
  }, [remindersLength]);

  if (isPending) return null;

  return (
    <div className='flex flex-col items-start gap-[0.5rem]'>
      {data?.remindersForOrganization
        ?.filter((r) => !r.dismissed)
        .sort((a, b) => {
          const diff =
            new Date(a?.dueDate).valueOf() - new Date(b?.dueDate).valueOf();

          if (diff === 0)
            return (
              new Date(a.metadata.lastUpdated).valueOf() -
              new Date(b.metadata.lastUpdated).valueOf()
            );

          return diff;
        })
        .map((r, i) => (
          <ReminderItem
            index={i}
            key={r.metadata.id}
            currentOwner={currentOwner}
            onChange={onChange}
            onDismiss={onDismiss}
            data={mapReminderToForm(r)}
          />
        ))}
    </div>
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
