import { observer } from 'mobx-react-lite';

import { Icon } from '@ui/media/Icon';
import { Image } from '@ui/media/Image/Image';
import { Spinner } from '@ui/feedback/Spinner';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import logoCustomerOs from '@shared/assets/customer-os-small.png';

export const LogoSection = observer(() => {
  const store = useStore();
  const isPlatformOwner = store?.globalCache?.value?.isPlatformOwner ?? false;

  return (
    <div className='flex justify-between'>
      <div className='py-2 pr-3 pl-[18px] flex items-center gap-1.5'>
        <Image
          width={20}
          height={20}
          alt='CustomerOS'
          fallbackSrc={logoCustomerOs}
          className='logo-image rounded'
          src={store.settings.tenant.value?.workspaceLogo}
        />
        <span className='font-semibold  text-start w-[fit-content] overflow-hidden text-ellipsis whitespace-nowrap'>
          {store.settings.tenant.value?.workspaceName || 'CustomerOS'}
        </span>

        {store.common?.impersonateAccounts?.length > 1 && isPlatformOwner && (
          <Tooltip label='Switch workspaces (G then W)'>
            <IconButton
              size='xxs'
              variant='ghost'
              aria-label='Switch workspace'
              dataTest={'switch-workspace-button'}
              icon={<Icon name={'arrow-switch-horizontal-02'} />}
              onClick={() => {
                store.ui.commandMenu.setType('SwitchWorkspace');
                store.ui.commandMenu.setOpen(true);
              }}
            />
          </Tooltip>
        )}
      </div>

      {(store.isSyncing || store.isBootstrapping) &&
        store.windowManager.networkStatus === 'online' && (
          <Tooltip
            side='bottom'
            delayDuration={0}
            label='Syncing the latest changes'
          >
            <div className='flex items-center'>
              <Spinner
                size='sm'
                label='Syncing'
                className='text-gray-300 fill-gray-700 mr-3'
              />
            </div>
          </Tooltip>
        )}

      {store.windowManager.networkStatus === 'offline' && (
        <Tooltip
          side='bottom'
          delayDuration={0}
          label={
            <>
              <p className='font-medium'>You’re offline</p>
              <p className='max-w-40'>
                Your changes are saving locally but may conflict with others
                when you reconnect
              </p>
            </>
          }
        >
          <div className='flex items-center'>
            <Icon name='cloud-off' className='mr-3 text-gray-500' />
          </div>
        </Tooltip>
      )}
    </div>
  );
});
