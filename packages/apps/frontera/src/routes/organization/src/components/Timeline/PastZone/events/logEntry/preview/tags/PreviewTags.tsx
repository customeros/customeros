import { useForm } from 'react-inverted-form';
import React, { useRef, useEffect } from 'react';

import { useQueryClient } from '@tanstack/react-query';

import { Tag } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useResetLogEntryTagsMutation } from '@organization/graphql/resetLogEntryTags.generated';
import { TagsSelect } from '@organization/components/Timeline/FutureZone/TimelineActions/logger/components/TagSelect';
import {
  LogEntryTagsDto,
  LogEntryTagsFormDtoI,
} from '@organization/components/Timeline/PastZone/events/logEntry/preview/tags/LogEntryTagsDto';

export const PreviewTags: React.FC<{
  id: string;
  isAuthor: boolean;
  tags?: Array<Tag>;
  tagOptions?: Array<{ label: string; value: string }>;
}> = ({ isAuthor, tags = [], id, tagOptions }) => {
  const logEntryStartedAtValues = new LogEntryTagsDto({ tags });
  const formId = 'preview-modal-log-entry-tag-update';
  const queryClient = useQueryClient();
  const client = getGraphQLClient();
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const updateLogEntryTags = useResetLogEntryTagsMutation(client, {
    onSuccess: () => {
      timeoutRef.current = setTimeout(
        () =>
          queryClient.invalidateQueries({ queryKey: ['GetTimeline.infinite'] }),
        500,
      );
    },
  });

  useForm<LogEntryTagsFormDtoI>({
    formId,
    defaultValues: logEntryStartedAtValues,

    stateReducer: (_state, action, next) => {
      if (action.type === 'FIELD_CHANGE') {
        updateLogEntryTags.mutate({
          id: id,
          input: [
            ...LogEntryTagsDto.toPayload({
              tags: action.payload.value,
            }).tags,
          ],
        });
      }

      return next;
    },
  });

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  return (
    <>
      {!isAuthor && (
        <p className='text-sm font-medium'>
          {tags.map(({ name }) => `#${name}`).join(' ')}
        </p>
      )}

      {isAuthor && (
        <p className='text-sm font-medium leading-1'>
          <TagsSelect name='tags' formId={formId} tags={tagOptions} />
        </p>
      )}
    </>
  );
};
