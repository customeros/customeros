import { FormEvent } from 'react';
import { useForm } from 'react-inverted-form';

import { produce } from 'immer';
import { InfiniteData, useQueryClient } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { useModKey } from '@shared/hooks/useModKey';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { useTimelineMeta } from '@organization/src/components/Timeline/state';
import { useTimelineRefContext } from '@organization/src/components/Timeline/context/TimelineRefContext';
import {
  GetTimelineQuery,
  useInfiniteGetTimelineQuery,
} from '@organization/src/graphql/getTimeline.generated';
import { useTimelineActionContext } from '@organization/src/components/Timeline/FutureZone/TimelineActions/context/TimelineActionContext';

import { ReminderPostit } from '../../../shared/ReminderPostit';

type ReminderForm = {
  content: string;
};

export const ReminderTimelineAction = () => {
  const { openedEditor } = useTimelineActionContext();
  const queryClient = useQueryClient();
  const { virtuosoRef } = useTimelineRefContext();
  const [timelineMeta] = useTimelineMeta();

  const isOpen = openedEditor === 'reminder';

  const { handleSubmit } = useForm<ReminderForm>({
    formId: 'reminder-form',
    defaultValues: {
      content: '',
    },
    onSubmit: async (values) => {
      let remindersCount = 0;

      queryClient.setQueryData<{ id: string; date: string; content: string }[]>(
        ['reminders'],
        (cache) => {
          return produce(cache, (draft) => {
            if (!draft) return;

            draft.unshift({
              id: Math.random().toString(),
              date: new Date().toISOString(),
              content: values.content,
            });

            remindersCount = draft.length;
          });
        },
      );

      const timelineData = queryClient.getQueryData<
        InfiniteData<GetTimelineQuery>
      >(useInfiniteGetTimelineQuery.getKey(timelineMeta.getTimelineVariables));

      const timelineItemsLength = Object.values(
        timelineData?.pages ?? [],
      ).reduce(
        (acc, curr) => curr.organization?.timelineEventsTotalCount + acc,
        0,
      );

      setTimeout(
        () =>
          virtuosoRef?.current?.scrollToIndex(
            timelineItemsLength + remindersCount,
          ),
        0,
      );
    },
  });

  useModKey('Enter', () => {
    handleSubmit({} as FormEvent<HTMLFormElement>);
  });

  if (!isOpen) {
    return null;
  }

  return (
    <ReminderPostit>
      <FormAutoresizeTextarea
        autoFocus
        px='4'
        name='content'
        formId='reminder-form'
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
        <Text>24 Mar • 09:09</Text>
        <Button variant='ghost' colorScheme='yellow' size='sm'>
          Delete
        </Button>
      </Flex>
    </ReminderPostit>
  );
};
