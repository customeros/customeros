import { useRef } from 'react';
import { useParams } from 'react-router-dom';
import { SelectInstance } from 'react-select';

import { observer } from 'mobx-react-lite';
import { useQueryClient } from '@tanstack/react-query';
import { useConnections } from '@integration-app/react';

import { Select } from '@ui/form/Select';
import { SelectOption } from '@ui/utils/types';
import { Button } from '@ui/form/Button/Button';
import { Link01 } from '@ui/media/icons/Link01';
import { useStore } from '@shared/hooks/useStore';
import { Unthread } from '@ui/media/logos/Unthread';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick';
import { useInfiniteGetTimelineQuery } from '@organization/graphql/getTimeline.generated';
import { useTimelineMeta } from '@organization/components/Timeline/state/TimelineMeta.atom';

export const ChannelLinkSelect = observer(() => {
  const store = useStore();
  const id = useParams()?.id as string;
  const slackChannels = store.settings.slack.channels;
  const organization = store.organizations.value.get(id);

  const queryClient = useQueryClient();

  const ref = useRef(null);
  const selectRef = useRef<SelectInstance>(null);
  const { open: isOpen, onClose, onOpen } = useDisclosure();
  const { items } = useConnections();
  const [timelineMeta] = useTimelineMeta();

  const timelineQueryKey = useInfiniteGetTimelineQuery.getKey({
    ...timelineMeta.getTimelineVariables,
  });

  const handleChange = (value: SelectOption) => {
    organization?.update((val) => {
      val.slackChannelId = value?.value ?? '';

      return val;
    });

    onClose();
    setTimeout(() => {
      queryClient.invalidateQueries({ queryKey: timelineQueryKey });
      store.timelineEvents.invalidateTimeline(id);
    }, 1000);
  };

  const options: SelectOption[] = slackChannels.map((el) => ({
    label: el?.channelName || el?.channelId,
    value: el?.channelId,
  }));

  const selectedChannelId = organization?.value?.slackChannelId;
  const value = options.find((el) => el.value === selectedChannelId);

  const hasUnthreadIntegration = items
    .map((i) => i.integration?.key)
    .some((i) => ['unthread'].includes(i ?? ''));

  useOutsideClick({ ref, handler: onClose });

  if (store.demoMode || !hasUnthreadIntegration) return null;

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
        leftElement={<Link01 className='text-gray-500 mr-2' />}
      />
    </div>
  );
});
