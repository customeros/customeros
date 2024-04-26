import { useRef } from 'react';
import { useParams } from 'next/navigation';
import { SelectInstance } from 'react-select';

import { produce } from 'immer';
import { useQueryClient } from '@tanstack/react-query';
import { useConnections } from '@integration-app/react';

import { Select } from '@ui/form/Select';
import { SelectOption } from '@ui/utils/types';
import { Button } from '@ui/form/Button/Button';
import { Link01 } from '@ui/media/icons/Link01';
import { Unthread } from '@ui/media/logos/Unthread';
import { toastError } from '@ui/presentation/Toast';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick';
import { useGetIssuesQuery } from '@organization/src/graphql/getIssues.generated';
import { useOrganizationQuery } from '@organization/src/graphql/organization.generated';
import { useSlackChannelsQuery } from '@organization/src/graphql/slackChannels.generated';
import { useUpdateOrganizationMutation } from '@shared/graphql/updateOrganization.generated';
import { useInfiniteGetTimelineQuery } from '@organization/src/graphql/getTimeline.generated';
import { useTimelineMeta } from '@organization/src/components/Timeline/state/TimelineMeta.atom';

interface ChannelLinkSelectProps {
  from: Date;
}

export const ChannelLinkSelect = ({ from }: ChannelLinkSelectProps) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const id = useParams()?.id as string;
  const ref = useRef(null);
  const selectRef = useRef<SelectInstance>(null);
  const { open: isOpen, onClose, onOpen } = useDisclosure();
  const { items } = useConnections();
  const [timelineMeta] = useTimelineMeta();

  const { data: organization, isPending: organizationIsPending } =
    useOrganizationQuery(client, { id });
  const { data, isPending } = useSlackChannelsQuery(client, {
    pagination: { page: 0, limit: 1000 },
  });

  const isLoading = organizationIsPending || isPending;
  const organizationQueryKey = useOrganizationQuery.getKey({ id });
  const issuesQueryKey = useGetIssuesQuery.getKey({
    organizationId: id,
    from,
    size: 50,
  });
  const timelineQueryKey = useInfiniteGetTimelineQuery.getKey({
    ...timelineMeta.getTimelineVariables,
  });

  const updateOrganization = useUpdateOrganizationMutation(client, {
    onMutate: ({ input }) => {
      queryClient.cancelQueries({ queryKey: organizationQueryKey });

      const previousEntry = useOrganizationQuery.mutateCacheEntry(queryClient, {
        id,
      })((cache) => {
        return produce(cache, (draft) => {
          if (!draft.organization) return;
          draft.organization['slackChannelId'] = input.slackChannelId;
        });
      });

      return { previousEntry };
    },
    onError: (_, __, context) => {
      if (context?.previousEntry) {
        queryClient.setQueryData([organizationQueryKey], context.previousEntry);
      }
      toastError(
        `We couldn't update the slack channel.`,
        'update-slack-channel-error',
      );
    },
    onSettled: () => {
      onClose();
      setTimeout(() => {
        queryClient.invalidateQueries({ queryKey: organizationQueryKey });
        queryClient.invalidateQueries({ queryKey: issuesQueryKey });
        queryClient.invalidateQueries({ queryKey: timelineQueryKey });
      }, 1000);
    },
  });

  const handleChange = (value: SelectOption) => {
    updateOrganization.mutate({
      input: {
        id,
        slackChannelId: value?.value ?? '',
        patch: true,
      },
    });
  };

  const options: SelectOption[] =
    isPending || !data?.slack_Channels?.content
      ? []
      : data?.slack_Channels?.content?.map((el) => ({
          label: el?.channelName || el?.channelId,
          value: el?.channelId,
        }));

  const selectedChannelId = organization?.organization?.slackChannelId;
  const value = options.find((el) => el.value === selectedChannelId);

  const hasUnthreadIntegration = items
    .map((i) => i.integration?.key)
    .some((i) => ['unthread'].includes(i ?? ''));

  useOutsideClick({ ref, handler: onClose });

  if (!hasUnthreadIntegration) return null;

  if (!isOpen) {
    if (!value) {
      return (
        <Button
          size='sm'
          variant='ghost'
          onClick={() => {
            onOpen();
            setTimeout(
              () => selectRef.current && selectRef.current?.focus(),
              0,
            );
          }}
          isLoading={isLoading}
          leftIcon={<Link01 color='gray.500' />}
        >
          Link Unthread Slack channel
        </Button>
      );
    }

    return (
      <Tooltip label={`Unlink ${value.label}`} hasArrow>
        <Button
          size='sm'
          variant='outline'
          onClick={() => {
            onOpen();
            setTimeout(
              () => selectRef.current && selectRef.current?.focus(),
              0,
            );
          }}
          className='rounded-full'
          isLoading={isLoading}
          leftIcon={<Unthread />}
        >
          Unthread issues linked
        </Button>
      </Tooltip>
    );
  }

  return (
    <div className='w-[210px]' ref={ref}>
      <Select
        size='sm'
        isClearable
        ref={selectRef}
        value={value}
        options={options}
        onChange={handleChange}
        onBlur={onClose}
        noOptionsMessage={() => 'No channel found'}
        openMenuOnClick={!value}
        placeholder='Slack channel'
        isLoading={updateOrganization.isPending}
        leftElement={<Link01 className='text-gray-500 mr-2' />}
      />
    </div>
  );
};
