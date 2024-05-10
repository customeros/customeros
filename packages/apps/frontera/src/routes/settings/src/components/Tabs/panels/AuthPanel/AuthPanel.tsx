import { useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import {
  useConnections,
  useIntegrations,
  useIntegrationApp,
} from '@integration-app/react';

import { Gmail } from '@ui/media/icons/Gmail';
import { Slack } from '@ui/media/icons/Slack';
import { Google } from '@ui/media/logos/Google';
import { Switch } from '@ui/form/Switch/Switch';
import { useStore } from '@shared/hooks/useStore';
import { Outlook } from '@ui/media/logos/Outlook';
import { toastError } from '@ui/presentation/Toast';
import { GCalendar } from '@ui/media/icons/GCalendar';
import { Spinner } from '@ui/feedback/Spinner/Spinner';

export const AuthPanel = observer(() => {
  const iApp = useIntegrationApp();
  const { items: iIntegrations } = useIntegrations();
  const { items: iConnections, refresh, loading } = useConnections();
  const { sessionStore, settingsStore } = useStore();
  const [queryParams] = useSearchParams();

  const outlookConnection = iConnections.find(
    (o) => o?.integration?.key === 'microsoft-outlook',
  );

  const handleOutlookToggle = async () => {
    const outlookIntegration = iIntegrations.find(
      (o) => o.key === 'microsoft-outlook',
    );

    if (!outlookIntegration) {
      toastError(
        'Microsoft Outlook integration not available',
        'get-intergration-data',
      );

      return;
    }

    try {
      await iApp
        .integration(outlookIntegration.key)
        .open({ showPoweredBy: false });
      await refresh();
    } catch (err) {
      toastError('Integration failed', 'get-intergration-data');
    }
  };

  useEffect(() => {
    if (
      queryParams &&
      queryParams.has('redirect_slack') &&
      queryParams.has('code')
    ) {
      settingsStore.slack.oauthCallback(queryParams.get('code') as string);
    }
  }, [queryParams]);

  const handleSyncGoogleToggle = (isChecked: boolean) => {
    if (isChecked) {
      settingsStore.google.enableSync();
    } else {
      settingsStore.google.disableSync();
    }
  };

  const handleSlackToggle = async (isChecked: boolean) => {
    if (isChecked) {
      settingsStore.slack.enableSync();
    } else {
      settingsStore.slack.disableSync();
    }
  };

  return (
    <>
      <div className='bg-gray-25 rounded-2xl flex-col flex relative max-w-[50%] '>
        <div className='px-6 pb-2'>
          <div className='flex gap-1 items-center mb-2 pt-5 '>
            <Google className='size-6' />
            <h1 className='text-gray-700 text-lg '>Google OAuth</h1>
          </div>
          <div className='w-full border-b border-gray-100' />
        </div>

        <div className='p-6 pr-0 pt-0 '>
          <p className='line-clamp-2 mt-2 mb-3'>
            Enable OAuth Integration to get access to your google workspace
            emails and calendar events
          </p>

          <button onClick={sessionStore.authenticate}>Click me</button>

          <div className='flex flex-col gap-2 w-[250px]'>
            <div className='flex gap-2 items-center'>
              <div className='flex flex-col items-start gap-4'>
                <div className='flex gap-1 items-center'>
                  <Gmail className='size-6' />
                  <label className='mb-0'>Sync Google Mail</label>
                </div>

                <div className='flex gap-1 items-center'>
                  <GCalendar className='size-6' />
                  <label className='mb-0'>Sync Google Calendar</label>
                </div>
              </div>

              {settingsStore.google.isLoading && (
                <Spinner
                  label='Google Loading'
                  className='text-white fill-success-500 size-5 ml-2'
                />
              )}
              {!settingsStore.google.isLoading && (
                <Switch
                  isChecked={settingsStore.google.gmailEnabled}
                  onChange={(value) => handleSyncGoogleToggle(value)}
                  colorScheme='success'
                />
              )}
            </div>
          </div>
        </div>
      </div>

      <div className='bg-gray-25 rounded-2xl flex-col flex relative max-w-[50%] '>
        <div className='px-6 pb-2'>
          <div className='flex gap-1 items-center mb-2 pt-5 '>
            <Outlook className='size-6' />
            <h1 className='text-gray-700 text-lg'>Microsoft Outlook</h1>
          </div>
          <div className='w-full border-b border-gray-100' />
        </div>

        <div className='p-6 pr-0 pt-0 '>
          <p className='line-clamp-2 mt-2 mb-3'>
            Enable OAuth Integration to get access to your microsoft outlook
            emails
          </p>

          <div className='flex space-x-4 items-center'>
            <div className='flex alig-middle space-x-1'>
              <Outlook className='size-6' />
              <label className='mb-0'>Sync Microsoft Outlook</label>
            </div>
            {loading ? (
              <Spinner
                label='Outlook Loading'
                className='text-white fill-success-500 size-5 ml-2'
              />
            ) : (
              <Switch
                colorScheme='success'
                onChange={handleOutlookToggle}
                isChecked={!!outlookConnection}
              />
            )}
          </div>
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
            {settingsStore.slack.isLoading && (
              <Spinner
                label='Slack Loading'
                className='text-white fill-success-500 size-5 ml-2'
              />
            )}
            {!settingsStore.slack.isLoading && (
              <Switch
                isChecked={settingsStore.slack.enabled}
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
