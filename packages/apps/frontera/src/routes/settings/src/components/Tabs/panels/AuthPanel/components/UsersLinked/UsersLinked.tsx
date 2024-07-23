import { observer } from 'mobx-react-lite';
import { OauthToken } from '@store/Settings/OauthTokenStore.store';

import { Avatar } from '@ui/media/Avatar';
import { Plus } from '@ui/media/icons/Plus';
import { Google } from '@ui/media/logos/Google';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Microsoft } from '@ui/media/icons/Microsoft';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { LinkBroken01 } from '@ui/media/icons/LinkBroken01';
import { RefreshCcw01 } from '@ui/media/icons/RefreshCcw01';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

interface UsersLinkedProps {
  title: string;
  tokenType: string;
}

export const UsersLinked = observer(
  ({ title, tokenType }: UsersLinkedProps) => {
    const store = useStore();

    const tokens: OauthToken[] =
      store.settings.oauthToken.tokens?.filter(
        (token) => token.type === tokenType,
      ) ?? [];

    return (
      <div className='flex flex-col gap-2 mb-8'>
        <div className='flex justify-between items-center'>
          <p className='font-semibold text-sm'>{title}</p>
          <Menu>
            <MenuButton>
              <div className='flex items-center gap-2 text-primary-700 text-xs font-semibold'>
                <Plus className='size-3' />
                Link account
              </div>
            </MenuButton>
            <MenuList>
              <MenuItem
                className='text-sm'
                onClick={() =>
                  store.settings.oauthToken.enableSync(tokenType, 'google')
                }
              >
                <Google className='mr-2' />
                Google Workspace
              </MenuItem>
              <MenuItem
                className='text-sm'
                onClick={() =>
                  store.settings.oauthToken.enableSync(tokenType, 'azure-ad')
                }
              >
                <Microsoft className='mr-2' />
                Microsoft Outlook
              </MenuItem>
            </MenuList>
          </Menu>
        </div>

        {tokens.length === 0 && (
          <p className='text-gray-500 text-sm'>No accounts connected</p>
        )}

        {tokens.map((token, idx) => {
          return (
            <div
              key={`${token.email}_${idx}`}
              className='flex justify-between items-center group'
            >
              <div className='flex gap-2'>
                <Avatar
                  size='xs'
                  src={''}
                  name={token.email}
                  variant={'outlineCircle'}
                />
                <p className='text-sm'>{token.email}</p>
              </div>
              <div className='flex items-center'>
                <Button
                  className='opacity-0 group-hover:opacity-100'
                  leftIcon={<LinkBroken01 />}
                  colorScheme='gray'
                  variant='ghost'
                  size='xxs'
                  onClick={() =>
                    store.settings.oauthToken.disableSync(
                      token.email,
                      token.provider,
                    )
                  }
                >
                  Unlink
                </Button>
                {token.needsManualRefresh && (
                  <Tooltip
                    className='max-w-[320px]'
                    label={`Your conversations and meetings are no longer syncing because access to your ${
                      token.provider === 'azure-ad'
                        ? 'Microsoft Outlook'
                        : 'Google Workspace'
                    } account has expired`}
                  >
                    <Button
                      colorScheme='warning'
                      variant='ghost'
                      leftIcon={<RefreshCcw01 className='text-warning-500' />}
                      size='xxs'
                      onClick={() =>
                        store.settings.oauthToken.enableSync(
                          tokenType,
                          token.provider,
                        )
                      }
                    >
                      Re-allow
                    </Button>
                  </Tooltip>
                )}
              </div>
            </div>
          );
        })}
      </div>
    );
  },
);
