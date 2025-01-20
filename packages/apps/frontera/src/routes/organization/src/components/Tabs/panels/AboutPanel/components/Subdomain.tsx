import { useState } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { XCircle } from '@ui/media/icons/XCircle';
import { getExternalUrl } from '@utils/getExternalLink';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { CornerDownRightDot } from '@ui/media/icons/CornerDownRightDot';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

export const Subdomain = observer(({ domain }: { domain: string }) => {
  const store = useStore();
  const id = useParams()?.id as string;

  const [showMenu, setShowMenu] = useState(false);

  const organization = store.organizations.getById(id);

  if (!organization || !organization?.value) return null;

  return (
    <div className='text-sm cursor-pointer flex items-center no-underline hover:no-underline text-gray-700 group'>
      <a
        target='_blank'
        rel='noreferrer noopener'
        href={getExternalUrl(domain ?? '/')}
        className=' no-underline hover:no-underline'
      >
        <CornerDownRightDot className='mr-2 size-3 text-grayModern-400' />
        <span className='hover:underline'>{domain}</span>
      </a>
      <Menu onOpenChange={(data) => setShowMenu(data)}>
        <MenuButton>
          <IconButton
            size='xxs'
            variant='ghost'
            aria-label={'Collapse'}
            icon={<DotsVertical />}
            className={cn('ml-2 opacity-0 group-hover:opacity-100', {
              '!opacity-100': showMenu,
            })}
          />
        </MenuButton>

        <MenuList className='min-w-[100px]'>
          <MenuItem
            onClick={() => {
              store.ui.commandMenu.setOpen(true);
              store.ui.commandMenu.setType('RemoveDomain');
              store.ui.commandMenu.setContext({
                ...store.ui.commandMenu.context,
                ids: [id],
                meta: {
                  domain,
                },
              });
            }}
          >
            <XCircle />
            Remove domain
          </MenuItem>
        </MenuList>
      </Menu>
    </div>
  );
});
