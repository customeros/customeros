import { observer } from 'mobx-react-lite';
import { OauthToken } from '@store/Settings/OauthTokenStore.store';

import { Avatar } from '@ui/media/Avatar';
import { Plus } from '@ui/media/icons/Plus';
import { Google } from '@ui/media/logos/Google';
import { Button } from '@ui/form/Button/Button';
import { Link01 } from '@ui/media/icons/Link01';
import { useStore } from '@shared/hooks/useStore';
import { Microsoft } from '@ui/media/icons/Microsoft';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { LinkBroken01 } from '@ui/media/icons/LinkBroken01';
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
      <div className='flex flex-col gap-2'>
        <div className='flex justify-between items-center'>
          <p className='font-semibold'>{title}</p>
          {title !== 'Team' && (
            <Menu>
              <MenuButton>
                <div className='flex items-center gap-2 text-primary-700 text-[12px] font-semibold'>
                  <Plus className='size-3' />
                  Link account
                </div>
              </MenuButton>
              <MenuList>
                <MenuItem
                  onClick={() =>
                    store.settings.oauthToken.enableSync(tokenType, 'google')
                  }
                >
                  <Google className='mr-2' />
                  Google Workspace
                </MenuItem>
                <MenuItem
                  onClick={() =>
                    store.settings.oauthToken.enableSync(tokenType, 'azure-ad')
                  }
                >
                  <Microsoft className='mr-2' />
                  Microsoft Outlook
                </MenuItem>
              </MenuList>
            </Menu>
          )}
        </div>

        {tokens.length === 0 && (
          <p className='text-gray-500'>No accounts connected</p>
        )}

        {tokens.map((token, idx) => {
          return (
            <div
              key={`${token.email}_${idx}`}
              className='flex justify-between hover:bg-gray-50 items-center group'
            >
              <div className='flex gap-2'>
                <Avatar
                  size='xs'
                  src={''}
                  name={token.email}
                  variant={'outlineCircle'}
                />
                <p>{token.email}</p>
              </div>
              {title !== 'Team' && (
                <div className='flex items-center'>
                  <Button
                    className='opacity-0 group-hover:opacity-100'
                    leftIcon={<LinkBroken01 />}
                    colorScheme='gray'
                    variant='ghost'
                    size='xs'
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
                      label={`Your conversations and meetings are no longer syncing because access to your ${
                        token.provider === 'azure-ad'
                          ? 'Microsoft Outlook'
                          : 'Google Workspace'
                      } account has expired`}
                    >
                      <Button
                        colorScheme='warning'
                        variant='ghost'
                        leftIcon={<Link01 className='text-warning-500' />}
                        size='xs'
                        onClick={() =>
                          store.settings.oauthToken.enableSync(
                            tokenType,
                            token.provider,
                          )
                        }
                      >
                        Re-link account
                      </Button>
                    </Tooltip>
                  )}
                </div>
              )}
            </div>
          );
        })}
      </div>
    );
  },
);
