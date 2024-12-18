import { useEffect } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { useTimelineMeta } from '@organization/components/Timeline/state';

import { ReminderItem } from './ReminderItem';

export const Reminders = observer(() => {
  const store = useStore();
  const organizationId = useParams()?.id as string;
  const reminders =
    store.reminders.valueByOrganization
      .get(organizationId)
      ?.map((r) => r.value) ?? [];

  const [_, setTimelineMeta] = useTimelineMeta();

  const remindersLength = reminders?.length ?? 0;
  const user = store.globalCache.value?.user;
  const currentOwner = [user?.firstName, user?.lastName]
    .filter(Boolean)
    .join(' ');

  useEffect(() => {
    setTimelineMeta((prev) => ({
      ...prev,
      remindersCount: remindersLength,
    }));
  }, [remindersLength]);

  useEffect(() => {
    store.reminders.bootstrapByOrganization(organizationId);
  }, []);

  if (store?.reminders?.isLoading) return null;

  return (
    <div
      data-test={'timeline-reminder-list'}
      className='flex flex-col items-start gap-[0.5rem]'
    >
      {reminders
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
            id={r.metadata.id}
            key={r.metadata.id}
            currentOwner={currentOwner}
          />
        ))}
    </div>
  );
});
