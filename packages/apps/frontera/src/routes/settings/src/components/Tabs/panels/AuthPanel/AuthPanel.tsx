import { useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { OwnerInput } from '@settings/components/Tabs/panels/AuthPanel/OwnerInput';
import {
  useConnections,
  useIntegrations,
  useIntegrationApp,
} from '@integration-app/react';

import { User } from '@graphql/types';
import { Gmail } from '@ui/media/icons/Gmail';
import { Slack } from '@ui/media/icons/Slack';
import { Button } from '@ui/form/Button/Button';
import { Google } from '@ui/media/logos/Google';
import { Switch } from '@ui/form/Switch/Switch';
import { useStore } from '@shared/hooks/useStore';
import { Outlook } from '@ui/media/logos/Outlook';
import { toastError } from '@ui/presentation/Toast';
import { Spinner } from '@ui/feedback/Spinner/Spinner';

export const AuthPanel = observer(() => {
  const iApp = useIntegrationApp();
  const { items: iIntegrations } = useIntegrations();
  const { items: iConnections, refresh, loading } = useConnections();
  const store = useStore();
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

          <div className='flex flex-col gap-2 w-[250px]'>
            <div className='flex flex-col'>
              <Button
                className='font-semibold rounded-lg py-1 px-3 text-sm items-center'
                variant='outline'
                colorScheme='gray'
                size='xs'
                onClick={() => store.settings.google.enableSync()}
              >
                <Gmail className='size-6' />
                Add gmail account
              </Button>

              {store.settings.google.isLoading && (
                <Spinner
                  label='Google Loading'
                  className='text-white fill-success-500 size-5 ml-2'
                />
              )}
              {!store.settings.google.isLoading && (
                <div className='grid grid-cols-1 gap-3 mt-3'>
                  {store.settings.google.tokens.map((token, i) => (
                    <div
                      key={token.email + '_' + i}
                      className='grid grid-cols-[200px_200px_minmax(100px,_1fr)] gap-2 items-center'
                    >
                      <div className='flex text-sm font-semibold'>
                        {token.email}
                      </div>

                      <div className='flex'>
                        <OwnerInput
                          id={token.userId}
                          owner={{ id: token.userId } as User}
                          onSelect={(owner) => {
                            store.settings.google.updateUser(
                              token.email,
                              owner.value,
                            );
                          }}
                        />
                      </div>

                      <div className='flex flex-row gap-3'>
                        <Button
                          className='font-semibold rounded-lg py-1 px-3 text-sm items-center'
                          variant='outline'
                          colorScheme='gray'
                          size='xs'
                          onClick={() =>
                            store.settings.google.disableSync(token.email)
                          }
                        >
                          Remove
                        </Button>

                        {token.needsManualRefresh && (
                          <Button
                            className='font-semibold rounded-lg py-1 px-3 text-sm items-center'
                            variant='outline'
                            colorScheme='gray'
                            size='xs'
                            onClick={() => store.settings.google.enableSync()}
                          >
                            Re-allow
                          </Button>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
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
