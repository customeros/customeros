import { FormEvent, useEffect } from 'react';
import { useForm } from 'react-inverted-form';

import { produce } from 'immer';
import { useQuery, useQueryClient } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { VStack } from '@ui/layout/Stack';
import { useModKey } from '@shared/hooks/useModKey';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';

// import { useTimelineMeta } from '@organization/src/components/Timeline/state';
// import {
//   GetTimelineQuery,
//   useInfiniteGetTimelineQuery,
// } from '@organization/src/graphql/getTimeline.generated';

import { ReminderPostit, ReminderDueDatePicker } from '../../shared';

type Reminder = {
  id: string;
  date: string;
  owner: string;
  content: string;
  isDismissed: boolean;
};

const mockData: Reminder[] = [
  {
    id: '1',
    date: '2021-10-01:12:30:00Z',
    content: 'Reminder 1',
    owner: 'customerostest',
    isDismissed: false,
  },
  {
    id: '2',
    date: '2021-10-02',
    content: 'Reminder 2',
    owner: 'Gigel',
    isDismissed: false,
  },
  {
    id: '3',
    date: '2021-10-03',
    content: 'Reminder 3',
    owner: 'Frone',
    isDismissed: false,
  },
];

export const Reminders = () => {
  const client = getGraphQLClient();
  const { data, isPending } = useQuery<Reminder[]>({
    queryKey: ['reminders'],
    queryFn: async () => {
      return new Promise((resolve) => resolve(mockData));
    },
  });

  // const [timelineMeta] = useTimelineMeta();

  const { data: globalCacheData } = useGlobalCacheQuery(client);

  const user = globalCacheData?.global_Cache?.user;
  const currentOwner = [user?.firstName, user?.lastName]
    .filter(Boolean)
    .join(' ');

  if (isPending) return <p>Loading...</p>;

  return (
    <VStack align='flex-start'>
      {data
        ?.filter((r) => !r.isDismissed)
        ?.map((r) => (
          <ReminderItem key={r.id} data={r} currentOwner={currentOwner} />
        ))}
    </VStack>
  );
};

interface ReminderItem {
  data: Reminder;
  currentOwner: string;
}

const ReminderItem = ({ data, currentOwner }: ReminderItem) => {
  const queryClient = useQueryClient();
  const formId = `reminder-edit-form-${data.id}`;

  const makeContentStr = (content: string, owner: string) => {
    const strippedContent = content.replace(`for ${owner}: `, '');

    return currentOwner !== owner
      ? `for ${owner}: ${strippedContent}`
      : strippedContent;
  };

  const { handleSubmit, setDefaultValues } = useForm<Reminder>({
    formId,
    defaultValues: data,
    onSubmit: async (values) => {
      queryClient.setQueryData<Reminder[]>(['reminders'], (cache) => {
        return produce(cache, (draft) => {
          if (!draft) return;

          const foundReminder = draft.find((r) => r.id === values.id);
          if (!foundReminder) return;

          foundReminder.content = values.content;
        });
      });
    },
    stateReducer: (_, action, next) => {
      if (action.type === 'FIELD_CHANGE' && action.payload.name === 'content') {
        return {
          ...next,
          values: {
            ...next.values,
            content: makeContentStr(next.values.content, next.values.owner),
          },
        };
      }

      return next;
    },
  });

  const updateReminder = () => {
    handleSubmit({} as FormEvent<HTMLFormElement>);
  };

  const dismissReminder = () => {
    queryClient.setQueryData<Reminder[]>(['reminders'], (cache) => {
      return produce(cache, (draft) => {
        if (!draft) return;

        const foundIdx = draft.findIndex((r) => r.id === data.id);

        draft.splice(foundIdx, 1);
      });
    });
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
          onClick={dismissReminder}
        >
          Dismiss
        </Button>
      </Flex>
    </ReminderPostit>
  );
};
