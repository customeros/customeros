import { useMemo } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import {
  cooldownPeriods,
  cooldownPeriodsMap,
  AddSlackChannelUsecase,
} from '@domain/usecases/agents/capabilities/add-slack-channel.usecase.ts';

import { Icon } from '@ui/media/Icon';
import { Switch } from '@ui/form/Switch';
import { Combobox } from '@ui/form/Combobox';
import { Slack } from '@ui/media/logos/Slack';
import { Button } from '@ui/form/Button/Button';
import { Tag, TagLabel } from '@ui/presentation/Tag';
import { Popover, PopoverTrigger, PopoverContent } from '@ui/overlay/Popover';

import { DisconnectSlackMenu } from './DisconnectSlackMenu.tsx';

export const SendSlackNotificationCapability = observer(() => {
  const { id } = useParams<{ id: string }>();

  const usecase = useMemo(() => new AddSlackChannelUsecase(id!), [id]);

  if (!usecase.isSlackEnabled) {
    return (
      <div>
        <div className='flex items-center justify-between mb-4'>
          <h2 className='text-sm font-medium '>
            Send Slack notification{` `}
            <span className='text-grayModern-500'>(optional)</span>
          </h2>
          <Switch
            size='sm'
            checked={usecase.slackNotificationActive}
            onChange={usecase.toggleSendSlackNotificationActive}
          />
        </div>

        {usecase.capabilityErrors && (
          <div className='bg-error-50 text-error-700 px-2 py-1 rounded-[4px] mb-4'>
            <Icon stroke='none' className='mr-2' name='dot-single' />
            <span className='text-sm'>{usecase.capabilityErrors}</span>
          </div>
        )}

        <div className='flex flex-col gap-4 w-full'>
          <div>
            <div className='flex justify-between flex-1'>
              <p className='text-sm font-medium'>Slack notifications</p>
            </div>
            <p className='text-sm mt-1'>
              We’ll notify you on Slack when a visitor is identified on your
              preferred websites
            </p>
          </div>
          <Button
            className='w-full'
            colorScheme='primary'
            onClick={usecase.enableSlack}
          >
            Connect a channel
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div>
      <div className='flex items-center justify-between mb-4'>
        <h2 className='text-sm font-medium '>
          Send Slack notification{` `}
          <span className='text-grayModern-500'>(optional)</span>
        </h2>
        <Switch
          size='sm'
          checked={usecase.slackNotificationActive}
          onChange={usecase.toggleSendSlackNotificationActive}
        />
      </div>
      <div className='flex flex-col gap-4 w-full'>
        <div>
          <div className='flex justify-between flex-1'>
            <p className='text-sm font-medium'>Slack notifications</p>
            <DisconnectSlackMenu />
          </div>
          <p className='text-sm mt-1'>
            We’ll notify you on Slack when a visitor is identified on your
            preferred websites
          </p>
        </div>

        <div className='w-full flex-1'>
          <Popover open={usecase.isOpen} onOpenChange={usecase.togglePopover}>
            <PopoverTrigger className={'flex items-center w-full'}>
              <div className='flex w-full items-center bg-white border border-grayModern-200 px-2 py-1 rounded-md min-h-8'>
                <Slack className='mr-3 text-gray-500' />
                {usecase.selectedChannel ? (
                  <div className='text-sm'>{usecase.selectedChannelName}</div>
                ) : (
                  <span className='text-gray-400 text-sm'>Slack channel</span>
                )}
              </div>
            </PopoverTrigger>
            <PopoverContent align='start' className='min-w-[264px] w-[375px]'>
              <Combobox
                options={usecase.options}
                value={usecase.selectedOption}
                placeholder='Search a channel'
                inputValue={usecase.inputValue}
                onInputChange={usecase.setInputValue}
                onChange={(newValue) => {
                  usecase.selectChannel(newValue.value);
                }}
                noOptionsMessage={({ inputValue }) => (
                  <div className='text-gray-700 px-3 py-1 mt-0.5 rounded-md bg-grayModern-100 gap-1 flex items-center'>
                    <span>{`No results matching "${inputValue}"`}</span>
                  </div>
                )}
              />
            </PopoverContent>
          </Popover>
        </div>

        <div>
          <p className='text-sm font-medium'>Notification cooldown period</p>
          <p className='text-sm mt-1'>
            How long should we wait before sending another notification after
            someone visited your website?
          </p>
        </div>

        <div className='flex items-center gap-1'>
          {cooldownPeriods.map((period) => (
            <Tag
              key={period}
              variant='subtle'
              onClick={() => usecase.setCooldownPeriod(period)}
              className='cursor-pointer hover:bg-primary-100 transition-colors'
              colorScheme={
                usecase.cooldownPeriod === period ? 'primary' : 'gray'
              }
            >
              <TagLabel>{cooldownPeriodsMap[period]}</TagLabel>
            </Tag>
          ))}
        </div>
      </div>
    </div>
  );
});
