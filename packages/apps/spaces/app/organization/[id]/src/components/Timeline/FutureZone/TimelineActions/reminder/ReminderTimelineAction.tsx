import { FormEvent } from 'react';
import { useForm } from 'react-inverted-form';

import { produce } from 'immer';
import { InfiniteData, useQueryClient } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { useModKey } from '@shared/hooks/useModKey';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
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
  const client = getGraphQLClient();
  const { openedEditor, closeEditor } = useTimelineActionContext();
  const queryClient = useQueryClient();
  const { virtuosoRef } = useTimelineRefContext();
  const [timelineMeta] = useTimelineMeta();
  const { data: globalCacheData } = useGlobalCacheQuery(client);

  const user = globalCacheData?.global_Cache?.user;
  const currentOwner = [user?.firstName, user?.lastName]
    .filter(Boolean)
    .join(' ');

  const isOpen = openedEditor === 'reminder';

  const { handleSubmit, reset } = useForm<ReminderForm>({
    formId: 'reminder-form',
    defaultValues: {
      content: '',
    },
    onSubmit: async (values) => {
      let remindersCount = 0;

      queryClient.setQueryData<
        { id: string; date: string; owner: string; content: string }[]
      >(['reminders'], (cache) => {
        return produce(cache, (draft) => {
          if (!draft) return;

          draft.unshift({
            id: Math.random().toString(),
            date: new Date().toISOString(),
            content: values.content,
            owner: currentOwner,
          });

          remindersCount = draft.length;
        });
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
          timelineItemsLength + remindersCount,
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
        <Text>24 Mar • 09:09</Text>
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
