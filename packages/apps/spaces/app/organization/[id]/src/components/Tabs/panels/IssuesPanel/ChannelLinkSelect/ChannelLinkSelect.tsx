import { useRef } from 'react';
import { useParams } from 'next/navigation';

import { produce } from 'immer';
import { useQueryClient } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Select } from '@ui/form/SyncSelect';
import { Tooltip } from '@ui/overlay/Tooltip';
import { Slack } from '@ui/media/logos/Slack';
import { Link01 } from '@ui/media/icons/Link01';
import { toastError } from '@ui/presentation/Toast';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { SelectOption, useDisclosure, useOutsideClick } from '@ui/utils';
import { useGetIssuesQuery } from '@organization/src/graphql/getIssues.generated';
import { useOrganizationQuery } from '@organization/src/graphql/organization.generated';
import { useSlackChannelsQuery } from '@organization/src/graphql/slackChannels.generated';
import { useUpdateOrganizationMutation } from '@shared/graphql/updateOrganization.generated';

interface ChannelLinkSelectProps {
  from: Date;
}

export const ChannelLinkSelect = ({ from }: ChannelLinkSelectProps) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const id = useParams()?.id as string;
  const ref = useRef(null);
  const { isOpen, onClose, onOpen } = useDisclosure();

  const { data: organization, isPending: organizationIsPending } =
    useOrganizationQuery(client, { id });
  const { data, isPending } = useSlackChannelsQuery(client, {
    pagination: { page: 0, limit: 1000 },
  });

  const isLoading = organizationIsPending || isPending;
  const queryKey = useOrganizationQuery.getKey({ id });
  const updateOrganization = useUpdateOrganizationMutation(client, {
    onMutate: ({ input }) => {
      queryClient.cancelQueries({ queryKey });

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
        queryClient.setQueryData([queryKey], context.previousEntry);
      }
      toastError(
        `We couldn't update the slack channel.`,
        'update-slack-channel-error',
      );
    },
    onSettled: () => {
      setTimeout(() => {
        queryClient.invalidateQueries({ queryKey });
        queryClient.invalidateQueries({
          queryKey: useGetIssuesQuery.getKey({
            organizationId: id,
            from,
            size: 50,
          }),
        });
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

  useOutsideClick({ ref, handler: onClose });

  if (!isOpen) {
    if (!value) {
      return (
        <Button
          size='sm'
          variant='ghost'
          onClick={onOpen}
          color='gray.500'
          fontWeight='normal'
          isLoading={isLoading}
          leftIcon={<Link01 color='gray.500' />}
        >
          Link Slack channel
        </Button>
      );
    }

    return (
      <Tooltip label={`Edit channel ${value.label}`} hasArrow>
        <Button
          size='sm'
          variant='outline'
          onClick={onOpen}
          color='gray.500'
          fontWeight='normal'
          borderRadius='full'
          isLoading={isLoading}
          leftIcon={<Slack />}
        >
          Channel linked
        </Button>
      </Tooltip>
    );
  }

  return (
    <Flex w='210px' ref={ref}>
      <Select
        size='sm'
        isClearable
        options={options}
        value={value}
        onChange={handleChange}
        openMenuOnClick={!value}
        placeholder='Paste Slack Channel ID'
        isLoading={updateOrganization.isPending}
        leftElement={<Link01 color='gray.500' mr='2' />}
      />
    </Flex>
  );
};
