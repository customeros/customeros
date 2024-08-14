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
          size='xs'
          variant='ghost'
          leftIcon={<Link01 color='gray.500' />}
          onClick={() => {
            onOpen();
            setTimeout(
              () => selectRef.current && selectRef.current?.focus(),
              0,
            );
          }}
        >
          Link Unthread Slack channel
        </Button>
      );
    }

    return (
      <Tooltip hasArrow label={`Unlink ${value.label}`}>
        <Button
          size='sm'
          variant='outline'
          leftIcon={<Unthread />}
          className='rounded-full'
          onClick={() => {
            onOpen();
            setTimeout(
              () => selectRef.current && selectRef.current?.focus(),
              0,
            );
          }}
        >
          Unthread issues linked
        </Button>
      </Tooltip>
    );
  }

  return (
    <div ref={ref} className='w-[210px]'>
      <Select
        size='sm'
        isClearable
        value={value}
        ref={selectRef}
        onBlur={onClose}
        options={options}
        onChange={handleChange}
        openMenuOnClick={!value}
        placeholder='Slack channel'
        noOptionsMessage={() => 'No channel found'}
        leftElement={<Link01 className='text-gray-500 mr-2' />}
      />
    </div>
  );
});
