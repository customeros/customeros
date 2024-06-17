import { useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { OauthToken } from '@store/Settings/OauthTokenStore.store';

import { Button } from '@ui/form/Button/Button';
import { Google } from '@ui/media/logos/Google';
import { Trash01 } from '@ui/media/icons/Trash01';
import { useStore } from '@shared/hooks/useStore';
import { Microsoft } from '@ui/media/icons/Microsoft';
import { Spinner } from '@ui/feedback/Spinner/Spinner';
import { PlusSquare } from '@ui/media/icons/PlusSquare';

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

  const renderTokenArea = (tokenType: string) => {
    const tokenLabel = tokenType.charAt(0).toUpperCase() + tokenType.slice(1);

    const tokens: OauthToken[] =
      store.settings.oauthToken.tokens?.filter(
        (token) => token.type === tokenType,
      ) ?? [];

    return (
      <>
        <div className='mb-2'>
          {tokenLabel}

          <div className='flex items-center pt-2 pb-2'>
            <Button
              className='p-0 mr-4 border-0'
              variant={'ghost'}
              onClick={() =>
                store.settings.oauthToken.enableSync(tokenType, 'google')
              }
            >
              <PlusSquare className='size-4' /> Google
            </Button>
            <Button
              className='p-0 border-0'
              variant={'ghost'}
              onClick={() =>
                store.settings.oauthToken.enableSync(tokenType, 'azure-ad')
              }
            >
              <PlusSquare className='size-4' /> Microsoft 365
            </Button>
          </div>

          <div className='w-full border-b border-gray-100' />

          {store.settings.oauthToken.isLoading && (
            <Spinner
              label='Google Loading'
              className='text-white fill-success-500 size-5 ml-2'
            />
          )}
          {tokens && (
            <div className='grid grid-cols-1 gap-3 mt-1'>
              {tokens.map((token: OauthToken, i: number) => (
                <div
                  key={token.email + '_' + i}
                  className='grid grid-cols-[200px_minmax(100px,_1fr)] gap-2 items-center'
                >
                  <div className='flex text-sm font-semibold'>
                    {token.provider === 'google' ? (
                      <Google className='size-5 mr-2' />
                    ) : (
                      <Microsoft className='size-5 mr-2' />
                    )}

                    {token.email}
                  </div>

                  <div className='flex flex-row gap-3'>
                    <Button
                      className='font-semibold rounded-lg py-1 px-3 text-sm items-center'
                      variant='outline'
                      colorScheme='gray'
                      size='xs'
                      onClick={() =>
                        store.settings.oauthToken.disableSync(
                          token.email,
                          token.provider,
                        )
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
                          store.settings.oauthToken.enableSync(
                            tokenType,
                            token.provider,
                          )
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
