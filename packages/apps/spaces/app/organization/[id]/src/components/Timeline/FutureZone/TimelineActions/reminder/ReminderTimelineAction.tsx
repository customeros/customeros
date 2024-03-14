import { FormEvent } from 'react';
import { useParams } from 'next/navigation';
import { useForm } from 'react-inverted-form';

import { produce } from 'immer';
import { InfiniteData, useQueryClient } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { useModKey } from '@shared/hooks/useModKey';
import { toastError } from '@ui/presentation/Toast';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { useTimelineMeta } from '@organization/src/components/Timeline/state';
import { useCreateReminderMutation } from '@organization/src/graphql/createReminder.generated';
import {
  RemindersQuery,
  useRemindersQuery,
} from '@organization/src/graphql/reminders.generated';
import { useTimelineRefContext } from '@organization/src/components/Timeline/context/TimelineRefContext';
import {
  GetTimelineQuery,
  useInfiniteGetTimelineQuery,
} from '@organization/src/graphql/getTimeline.generated';
import { useTimelineActionContext } from '@organization/src/components/Timeline/FutureZone/TimelineActions/context/TimelineActionContext';

import { ReminderPostit, ReminderDueDatePicker } from '../../../shared';

type ReminderForm = {
  date: string;
  content: string;
};

export const ReminderTimelineAction = () => {
  const organizationId = useParams()?.id as string;
  const client = getGraphQLClient();
  const { openedEditor, closeEditor } = useTimelineActionContext();
  const queryClient = useQueryClient();
  const { virtuosoRef } = useTimelineRefContext();
  const [timelineMeta] = useTimelineMeta();
  const { data: globalCacheData } = useGlobalCacheQuery(client);

  const remindersQueryKey = useRemindersQuery.getKey({ organizationId });
  const createReminder = useCreateReminderMutation(client, {
    onMutate: (values) => {
      queryClient.cancelQueries({ queryKey: remindersQueryKey });

      const previousEntries = useRemindersQuery.mutateCacheEntry(queryClient, {
        organizationId,
      })((cache) =>
        produce(cache, (draft) => {
          if (!draft) return;

          draft.remindersForOrganization.push({
            metadata: {
              id: 'TEMP',
            },
            dueDate: values.input.dueDate,
            content: values.input.content,
            owner: {
              id: globalCacheData?.global_Cache?.user?.id ?? '',
              firstName: globalCacheData?.global_Cache?.user?.firstName ?? '',
              lastName: globalCacheData?.global_Cache?.user?.lastName ?? '',
            },
            dismissed: false,
          });
        }),
      );

      return { previousEntries };
    },
    onError: (_, __, context) => {
      if (context?.previousEntries) {
        queryClient.setQueryData(remindersQueryKey, context.previousEntries);
      }
      toastError(`We couldn't create the reminder`, 'create-reminder-error');
    },
    onSettled: () => {
      setTimeout(() => {
        queryClient.invalidateQueries({ queryKey: remindersQueryKey });
      }, 500);
    },
  });

  const isOpen = openedEditor === 'reminder';

  const { handleSubmit, reset } = useForm<ReminderForm>({
    formId: 'reminder-form',
    defaultValues: {
      content: '',
      date: new Date().toISOString(),
    },
    onSubmit: async (values) => {
      const remindersCount =
        queryClient.getQueryData<RemindersQuery>(remindersQueryKey)
          ?.remindersForOrganization.length ?? 0;

      createReminder.mutate({
        input: {
          content: values.content,
          dueDate: values.date,
          organizationId,
          userId: globalCacheData?.global_Cache?.user?.id ?? '',
        },
      });

      const timelineData = queryClient.getQueryData<
        InfiniteData<GetTimelineQuery>
      >(useInfiniteGetTimelineQuery.getKey(timelineMeta.getTimelineVariables));

      const timelineItemsLength = Object.values(
        timelineData?.pages ?? [],
      ).reduce(
        (acc, curr) => curr.organization?.timelineEventsTotalCount + acc,
        0,
      );

      setTimeout(() => {
        virtuosoRef?.current?.scrollToIndex(
          timelineItemsLength + remindersCount + 1,
        );
        closeEditor();
        reset();
      }, 0);
    },
  });

  useModKey('Enter', () => {
    handleSubmit({} as FormEvent<HTMLFormElement>);
  });

  if (!isOpen) {
    return null;
  }

  return (
    <ReminderPostit
      _hover={{
        '& #sticky-body > #sticky-footer > button': {
          visibility: 'visible',
        },
      }}
    >
      <FormAutoresizeTextarea
        autoFocus
        px='4'
        fontFamily='sticky'
        fontSize='sm'
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
      <Flex
        mb='2'
        px='4'
        w='full'
        align='center'
        id='sticky-footer'
        justify='space-between'
      >
        <ReminderDueDatePicker name='date' formId='reminder-form' />
        <Button
          size='sm'
          variant='ghost'
          visibility='hidden'
          colorScheme='yellow'
          onClick={closeEditor}
        >
          Dismiss
        </Button>
      </Flex>
    </ReminderPostit>
  );
};
