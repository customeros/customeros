import { useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { Switch } from '@ui/form/Switch';
import { Spinner } from '@ui/feedback/Spinner';
import { useStore } from '@shared/hooks/useStore';

import { UsersLinked, LinkedInSettings } from './components';

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
      <div className='bg-gray-25 flex-col flex relative max-w-[550px] px-6 pb-4 pt-2 '>
        <div className='flex gap-4 flex-col'>
          <div className='flex flex-col'>
            <h1
              data-test='settings-accounts-header'
              className='text-gray-700 font-semibold '
            >
              Accounts
            </h1>
          </div>
        </div>
      </div>

      <div className='flex flex-col max-w-[550px] px-6 '>
        <div className='flex items-center gap-1'>
          <h2 className='text-gray-700 text-sm font-medium'>Email</h2>
          <div className='w-full border-b border-gray-100 mx-2' />
        </div>
        <p className='line-clamp-2 mt-2 mb-4 text-sm'>
          Get all your customer contacts, conversations and meetings in one
          place by importing them from Google workspace or Microsoft Outlook.
        </p>

        <UsersLinked title='Team' tokenType='WORKSPACE' />
        <UsersLinked title='Outbound' tokenType='OUTBOUND' />
      </div>

      <article className='flex-col flex relative max-w-[550px] '>
        <div className='px-6 flex items-center w-full'>
          <div className='flex items-center gap-1'>
            <h2 className='text-gray-700 text-sm font-medium'>Slack</h2>
          </div>
          <div className='w-full border-b border-gray-100 mx-2' />
          {store.settings.slack.isLoading && (
            <Spinner
              label='Slack Loading'
              className='text-white fill-success-500 size-5 ml-2'
            />
          )}
          {!store.settings.slack.isLoading && (
            <div className='flex items-center'>
              <Switch
                size='sm'
                colorScheme='primary'
                isChecked={store.settings.slack.enabled}
                onChange={(isChecked) => handleSlackToggle(isChecked)}
              />
            </div>
          )}
        </div>

        <div className='p-6 pr-0 pt-0'>
          <p className='line-clamp-2 mt-2 mb-3 text-sm'>
            Sync Slack messages to your organizations’ timeline and direct
            website visitors to specified channels
          </p>
        </div>
      </article>

      <LinkedInSettings />
    </>
  );
});
