import { useState } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { XCircle } from '@ui/media/icons/XCircle';
import { Globe06 } from '@ui/media/icons/Globe06';
import { PlusCircle } from '@ui/media/icons/PlusCircle';
import { getExternalUrl } from '@utils/getExternalLink';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { ChevronExpand } from '@ui/media/icons/ChevronExpand';
import { ChevronCollapse } from '@ui/media/icons/ChevronCollapse';
import { CornerDownRightDot } from '@ui/media/icons/CornerDownRightDot';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

export const Domains = observer(() => {
  const store = useStore();
  const id = useParams()?.id as string;

  const [isExpanded, setIsExpanded] = useState(false);
  const [showMenu, setShowMenu] = useState(false);

  const organization = store.organizations.getById(id);

  if (!organization || !organization?.value) return null;

  return (
    <div className='flex flex-col mt-1 group'>
      <div className='flex items-center justify-between w-full  '>
        <a
          target='_blank'
          rel='noreferrer noopener'
          href={getExternalUrl(organization.value.domains[0] ?? '/')}
          className='w-full cursor-pointer text-sm flex items-center no-underline hover:no-underline text-gray-700'
        >
          <Globe06 className='mr-2' />
          {organization.value.domains[0]}
        </a>
        <div
          className={cn('flex opacity-0 group-hover:opacity-100 gap-1', {
            '!opacity-100': showMenu,
          })}
        >
          <div>
            <IconButton
              size='xxs'
              variant='ghost'
              onClick={() => setIsExpanded(!isExpanded)}
              aria-label={isExpanded ? 'Collapse' : 'Expand'}
              icon={isExpanded ? <ChevronCollapse /> : <ChevronExpand />}
            />
          </div>

          <Menu onOpenChange={(data) => setShowMenu(data)}>
            <MenuButton>
              <IconButton
                size='xxs'
                variant='ghost'
                aria-label={'Collapse'}
                icon={<DotsVertical />}
              />
            </MenuButton>

            <MenuList className='min-w-[100px]'>
              <MenuItem
                onClick={() => {
                  store.ui.commandMenu.setOpen(true);
                  store.ui.commandMenu.setType('AddNewDomain');
                  store.ui.commandMenu.setContext({
                    ids: [id],
                    entity: 'Organization',
                  });
                }}
              >
                <PlusCircle />
                Add domain
              </MenuItem>
              <MenuItem onClick={() => {}}>
                <XCircle />
                Remove domain
              </MenuItem>
            </MenuList>
          </Menu>
        </div>
      </div>

      {isExpanded && (
        <div>
          {organization.value.domains.slice(1).map((domain) => (
            <a
              target='_blank'
              rel='noreferrer noopener'
              href={getExternalUrl(domain ?? '/')}
              className='text-sm cursor-pointer flex items-center no-underline hover:no-underline text-gray-700'
            >
              <CornerDownRightDot className='mr-2 size-3' />
              {domain}
            </a>
          ))}
        </div>
      )}
    </div>
  );
});
