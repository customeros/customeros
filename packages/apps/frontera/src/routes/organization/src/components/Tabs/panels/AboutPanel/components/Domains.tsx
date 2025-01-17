import { useState } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Domain } from '@graphql/types';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { XCircle } from '@ui/media/icons/XCircle';
import { Globe01 } from '@ui/media/icons/Globe01';
import { PlusCircle } from '@ui/media/icons/PlusCircle';
import { getExternalUrl } from '@utils/getExternalLink';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { ChevronExpand } from '@ui/media/icons/ChevronExpand';
import { ChevronCollapse } from '@ui/media/icons/ChevronCollapse';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

import { Subdomain } from './Subdomain';

export const Domains = observer(() => {
  const store = useStore();
  const id = useParams()?.id as string;

  // Track expanded state for each domain group
  const [expandedDomains, setExpandedDomains] = useState<
    Record<string, boolean>
  >({});
  const [showMenus, setShowMenus] = useState<Record<string, boolean>>({});

  const organization = store.organizations.getById(id);

  if (!organization || !organization?.value) return null;

  const formattedData = formatDomains(organization.value.domainsDetails);

  const toggleExpanded = (domainKey: string) => {
    setExpandedDomains((prev) => ({
      ...prev,
      [domainKey]: !prev[domainKey],
    }));
  };

  const handleMenuOpen = (domainKey: string, isOpen: boolean) => {
    setShowMenus((prev) => ({
      ...prev,
      [domainKey]: isOpen,
    }));
  };

  if (!formattedData.length) {
    return (
      <div className='flex flex-col mt-1 h-6'>
        <div
          tabIndex={0}
          role={'button'}
          className='flex items-center w-full gap-2 group h-full'
          onClick={() => {
            store.ui.commandMenu.setOpen(true);
            store.ui.commandMenu.setType('AddNewDomain');
            store.ui.commandMenu.setContext({
              ids: [id],
              entity: 'Organization',
            });
          }}
        >
          <Globe01 className='text-gray-500' />
          <span
            className='text-sm text-gray-400'
            data-test='org-about-domain-empty'
          >
            Add domain
          </span>
        </div>
      </div>
    );
  }

  return (
    <div className='flex flex-col mt-1'>
      {formattedData.map((domainGroup, index) => (
        <div className='flex flex-col' key={domainGroup.primaryDomain.domain}>
          <div className='flex items-center w-full gap-2 group'>
            <a
              target='_blank'
              rel='noreferrer noopener'
              href={getExternalUrl(domainGroup.primaryDomain.domain || '/')}
              className='w-fit cursor-pointer text-sm flex items-center no-underline hover:no-underline text-gray-700'
            >
              {index === 0 && <Globe01 className='mr-2 text-gray-500' />}
              <span
                data-test='org-about-domain-filled'
                className={cn({
                  'ml-6': index !== 0,
                })}
              >
                {domainGroup.primaryDomain.domain}
              </span>
            </a>
            <div
              className={cn('flex opacity-0 group-hover:opacity-100 gap-1', {
                '!opacity-100': showMenus[domainGroup.primaryDomain.domain],
              })}
            >
              <div>
                {domainGroup.subdomains?.length > 0 && (
                  <IconButton
                    size='xxs'
                    variant='ghost'
                    onClick={() =>
                      toggleExpanded(domainGroup.primaryDomain.domain)
                    }
                    aria-label={
                      expandedDomains[domainGroup.primaryDomain.domain]
                        ? 'Collapse'
                        : 'Expand'
                    }
                    icon={
                      expandedDomains[domainGroup.primaryDomain.domain] ? (
                        <ChevronCollapse />
                      ) : (
                        <ChevronExpand />
                      )
                    }
                  />
                )}
              </div>
              <Menu
                onOpenChange={(isOpen) =>
                  handleMenuOpen(domainGroup.primaryDomain.domain, isOpen)
                }
              >
                <MenuButton>
                  <IconButton
                    size='xxs'
                    variant='ghost'
                    aria-label={'Collapse'}
                    icon={<DotsVertical />}
                  />
                </MenuButton>

                <MenuList className='min-w-[100px]'>
                  {index === 0 && (
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
                  )}

                  <MenuItem
                    onClick={() => {
                      store.ui.commandMenu.setOpen(true);
                      store.ui.commandMenu.setType('RemoveDomain');
                      store.ui.commandMenu.setContext({
                        ...store.ui.commandMenu.context,
                        ids: [id],
                        meta: {
                          domain: domainGroup.primaryDomain.domain,
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

          {expandedDomains[domainGroup.primaryDomain.domain] && (
            <div className='ml-6'>
              {domainGroup.subdomains.map(({ domain }) => (
                <Subdomain key={domain} domain={domain} />
              ))}
            </div>
          )}
        </div>
      ))}
    </div>
  );
});
interface DomainGroup {
  subdomains: Domain[];
  primaryDomain: Domain;
}

function formatDomains(domainsDetails: Domain[]): DomainGroup[] {
  const primaryDomainsMap = new Map();
  const result: DomainGroup[] = [];

  for (const domain of domainsDetails) {
    if (domain.primary) {
      const resultObj = {
        primaryDomain: domain,
        subdomains: [],
      };

      primaryDomainsMap.set(domain.domain, result.length);
      result.push(resultObj);
    }
  }

  for (const domain of domainsDetails) {
    if (!domain.primary) {
      const primaryIndex = primaryDomainsMap.get(domain.primaryDomain);

      if (primaryIndex !== undefined && !!domain) {
        result[primaryIndex].subdomains.push(domain);
      }
    }
  }

  return result;
}
