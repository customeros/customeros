import { useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { GoogleToken } from '@store/Settings/Google.store';
import {
  useConnections,
  useIntegrations,
  useIntegrationApp,
} from '@integration-app/react';

import { Button } from '@ui/form/Button/Button';
import { Google } from '@ui/media/logos/Google';
import { Switch } from '@ui/form/Switch/Switch';
import { Trash01 } from '@ui/media/icons/Trash01';
import { useStore } from '@shared/hooks/useStore';
import { Outlook } from '@ui/media/logos/Outlook';
import { toastError } from '@ui/presentation/Toast';
import { Spinner } from '@ui/feedback/Spinner/Spinner';
import { PlusSquare } from '@ui/media/icons/PlusSquare';

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

  const renderTokenArea = (tokenType: string) => {
    const tokenLabel = tokenType.charAt(0).toUpperCase() + tokenType.slice(1);

    const tokens: GoogleToken[] =
      store.settings.google.tokens?.filter(
        (token) => token.type === tokenType,
      ) ?? [];

    return (
      <>
        <div className='mb-2'>
          <div className='flex items-center'>
            <Button
              className='size-[30px] p-0 border-1'
              onClick={() => store.settings.google.enableSync(tokenType)}
            >
              <PlusSquare className='size-4' />
            </Button>
            {tokenLabel}
          </div>

          {store.settings.google.isLoading && (
            <Spinner
              label='Google Loading'
              className='text-white fill-success-500 size-5 ml-2'
            />
          )}
          {tokens && (
            <div className='grid grid-cols-1 gap-3 mt-1 ml-3'>
              {tokens.map((token: GoogleToken, i: number) => (
                <div
                  key={token.email + '_' + i}
                  className='grid grid-cols-[200px_minmax(100px,_1fr)] gap-2 items-center'
                >
                  <div className='flex text-sm font-semibold'>
                    {token.email}
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
                      <Trash01 className='size-4' />
                    </Button>

                    {token.needsManualRefresh && (
                      <Button
                        className='font-semibold rounded-lg py-1 px-3 text-sm items-center'
                        variant='outline'
                        colorScheme='gray'
                        size='xs'
                        onClick={() =>
                          store.settings.google.enableSync(tokenType)
                        }
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
      </>
    );
  };

  return (
    <>
      <div className='bg-gray-25 rounded-2xl flex-col flex relative max-w-[50%] '>
        <div className='px-6 pb-2'>
          <div className='flex gap-1 items-center mb-2 pt-5 '>
            <Google className='size-6' />
            <h1 className='text-gray-700 text-lg '>Google Workspace</h1>
          </div>
          <div className='w-full border-b border-gray-100' />
        </div>

        <div className='p-6 pr-0 pt-0 '>
          {renderTokenArea('PERSONAL')}
          {renderTokenArea('WORKSPACE')}
          {renderTokenArea('OUTBOUND')}
        </div>
      </div>

      <div className='bg-gray-25 rounded-2xl flex-col flex relative max-w-[50%] '>
        <div className='px-6 pb-2'>
          <div className='flex gap-1 items-center mb-2 pt-5 '>
            <Outlook className='size-6' />
            <h1 className='text-gray-700 text-lg'>Microsoft 365</h1>
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

      {/*<div className='bg-gray-25 rounded-2xl flex-col mt-4 flex relative max-w-[50%] '>*/}
      {/*  <div className='px-6 pb-2'>*/}
      {/*    <div className='flex items-center gap-1 mb-2'>*/}
      {/*      <Slack className='size-6' />*/}
      {/*      <h1 className='text-gray-700 text-lg'>Slack</h1>*/}
      {/*    </div>*/}
      {/*    <div className='w-full border-b border-gray-100' />*/}
      {/*  </div>*/}

      {/*  <div className='p-6 pr-0 pt-0'>*/}
      {/*    <p className='line-clamp-2 mt-2 mb-3'>*/}
      {/*      Enable Slack Integration to get access to your Slack workspace*/}
      {/*    </p>*/}

      {/*    <div className='flex space-x-4 items-center'>*/}
      {/*      <div className='flex alig-middle space-x-1'>*/}
      {/*        <Slack className='size-6' />*/}
      {/*        <label className='mb-0'>Sync Slack</label>*/}
      {/*      </div>*/}
      {/*      {store.settings.slack.isLoading && (*/}
      {/*        <Spinner*/}
      {/*          label='Slack Loading'*/}
      {/*          className='text-white fill-success-500 size-5 ml-2'*/}
      {/*        />*/}
      {/*      )}*/}
      {/*      {!store.settings.slack.isLoading && (*/}
      {/*        <Switch*/}
      {/*          isChecked={store.settings.slack.enabled}*/}
      {/*          colorScheme='success'*/}
      {/*          onChange={(isChecked) => handleSlackToggle(isChecked)}*/}
      {/*        />*/}
      {/*      )}*/}
      {/*    </div>*/}
      {/*  </div>*/}
      {/*</div>*/}
    </>
  );
});
