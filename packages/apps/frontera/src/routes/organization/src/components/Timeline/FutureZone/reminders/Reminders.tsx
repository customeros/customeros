import { useEffect } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useTimelineMeta } from '@organization/components/Timeline/state';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { useRemindersQuery } from '@organization/graphql/reminders.generated';

import { ReminderItem } from './ReminderItem';

export const Reminders = observer(() => {
  const store = useStore();
  const organizationId = useParams()?.id as string;
  const reminders =
    store.reminders.valueByOrganization
      .get(organizationId)
      ?.map((r) => r.value) ?? [];

  const client = getGraphQLClient();
  const [_, setTimelineMeta] = useTimelineMeta();

  const { data, isPending } = useRemindersQuery(client, { organizationId });
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

  useEffect(() => {
    store.reminders.bootstrapByOrganization(organizationId);
  }, []);

  if (isPending) return null;

  return (
    <div className='flex flex-col items-start gap-[0.5rem]'>
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
