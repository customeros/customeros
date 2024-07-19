import { useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { Switch } from '@ui/form/Switch';
import { Slack } from '@ui/media/icons/Slack';
import { Spinner } from '@ui/feedback/Spinner';
import { useStore } from '@shared/hooks/useStore';

import { UsersLinked } from './components';

export const AuthPanel = observer(() => {
  const store = useStore();
  const [queryParams] = useSearchParams();

  useEffect(() => {
    if (
      queryParams &&
      queryParams.has('redirect_slack') &&
      queryParams.has('code')
    ) {
      store.settings.slack.oauthCallback(queryParams.get('code') as string);
    }
  }, [queryParams]);

  const handleSlackToggle = async (isChecked: boolean) => {
    if (isChecked) {
      store.settings.slack.enableSync();
    } else {
      store.settings.slack.disableSync();
    }
  };

  return (
    <>
      <div className='bg-gray-25 flex-col flex relative max-w-[550px] px-6 py-4 '>
        <div className='flex gap-4 flex-col'>
          <div className='flex flex-col'>
            <h1 className='text-gray-700 text-lg font-semibold '>Accounts</h1>
            <p>
              Get all your customer contacts, conversations and meetings in one
              place by importing them from Google workspace or Microsoft
              Outlook.
            </p>
          </div>
          <UsersLinked title='Team' tokenType='WORKSPACE' />
          <UsersLinked title='Outbound' tokenType='OUTBOUND' />
        </div>
      </div>

      <div className='bg-gray-25 rounded-2xl flex-col mt-4 flex relative max-w-[50%] '>
        <div className='px-6 pb-2'>
          <div className='flex items-center gap-1 mb-2'>
            <Slack className='size-6' />
            <h1 className='text-gray-700 text-lg'>Slack</h1>
          </div>
          <div className='w-full border-b border-gray-100' />
        </div>

        <div className='p-6 pr-0 pt-0'>
          <p className='line-clamp-2 mt-2 mb-3'>
            Enable Slack Integration to get access to your Slack workspace
          </p>

          <div className='flex space-x-4 items-center'>
            <div className='flex alig-middle space-x-1'>
              <Slack className='size-6' />
              <label className='mb-0'>Sync Slack</label>
            </div>
            {store.settings.slack.isLoading && (
              <Spinner
                label='Slack Loading'
                className='text-white fill-success-500 size-5 ml-2'
              />
            )}
            {!store.settings.slack.isLoading && (
              <Switch
                isChecked={store.settings.slack.enabled}
                colorScheme='success'
                onChange={(isChecked) => handleSlackToggle(isChecked)}
              />
            )}
          </div>
        </div>
      </div>
    </>
  );
});
