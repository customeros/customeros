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
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { Subdomain } from '@organization/components/Tabs/panels/AboutPanel/components/Subdomain.tsx';

export const Domains = observer(() => {
  const store = useStore();
  const id = useParams()?.id as string;

  const [isExpanded, setIsExpanded] = useState(false);
  const [showMenu, setShowMenu] = useState(false);

  const organization = store.organizations.getById(id);

  if (!organization || !organization?.value) return null;

  const domains = organization.value.domainsDetails.filter(
    (e) => !e.primaryDomain,
  );
  const primaryDomain =
    organization.value.domainsDetails.find((e) => e.primaryDomain) ||
    domains[0];

  return (
    <div className='flex flex-col mt-1 '>
      <div className='flex items-center w-full gap-2  group'>
        <a
          target='_blank'
          rel='noreferrer noopener'
          href={getExternalUrl(primaryDomain?.domain || '/')}
          className='w-fit cursor-pointer text-sm flex items-center no-underline hover:no-underline text-gray-700'
        >
          <Globe06 className='mr-2 text-gray-500' />
          {primaryDomain?.domain}
        </a>
        <div
          className={cn('flex opacity-0 group-hover:opacity-100 gap-1', {
            '!opacity-100': showMenu,
          })}
        >
          <div>
            {domains?.length > 0 && (
              <IconButton
                size='xxs'
                variant='ghost'
                onClick={() => setIsExpanded(!isExpanded)}
                aria-label={isExpanded ? 'Collapse' : 'Expand'}
                icon={isExpanded ? <ChevronCollapse /> : <ChevronExpand />}
              />
            )}
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
              <MenuItem
                onClick={() => {
                  store.ui.commandMenu.setOpen(true);
                  store.ui.commandMenu.setType('RemoveDomain');
                  store.ui.commandMenu.setContext({
                    ...store.ui.commandMenu.context,
                    ids: [id],
                    meta: {
                      domain: primaryDomain?.domain,
                      isPrimary: true,
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
      </div>

      {isExpanded && (
        <div className='ml-6'>
          {domains.map(({ domain }) => (
            <Subdomain key={domain} domain={domain} />
          ))}
        </div>
      )}
    </div>
  );
});
