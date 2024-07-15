import { useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

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
    </>
  );
});
