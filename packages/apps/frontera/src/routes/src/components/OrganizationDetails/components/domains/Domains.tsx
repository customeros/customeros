import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { Domain } from '@graphql/types';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Button } from '@ui/form/Button/Button.tsx';
import { XCircle } from '@ui/media/icons/XCircle.tsx';
import { Globe01 } from '@ui/media/icons/Globe01.tsx';
import { getExternalUrl } from '@utils/getExternalLink.ts';
import { PlusCircle } from '@ui/media/icons/PlusCircle.tsx';
import { DotsVertical } from '@ui/media/icons/DotsVertical.tsx';
import { ChevronExpand } from '@ui/media/icons/ChevronExpand.tsx';
import { ChevronCollapse } from '@ui/media/icons/ChevronCollapse.tsx';
import {
  Menu,
  MenuItem,
  MenuList,
  MenuButton,
} from '@ui/overlay/Menu/Menu.tsx';

import { Subdomain } from './Subdomain.tsx';

export const Domains = observer(({ id }: { id: string }) => {
  const store = useStore();

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
        <div className='flex items-center w-full gap-2 group h-full'>
          <Globe01 className='text-gray-500' />
          <span
            tabIndex={0}
            role={'button'}
            className='text-sm text-gray-400'
            data-test='org-about-domain-empty'
            onClick={() => {
              store.ui.commandMenu.setOpen(true);
              store.ui.commandMenu.setType('AddNewDomain');
              store.ui.commandMenu.setContext({
                ids: [id],
                entity: 'Organization',
              });
            }}
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
          <div
            className={cn('flex items-center w-full gap-2 group', {
              'ml-6': index !== 0,
            })}
          >
            {index === 0 && <Globe01 className='text-gray-500' />}

            <a
              target='_blank'
              rel='noreferrer noopener'
              href={getExternalUrl(domainGroup.primaryDomain.domain || '/')}
              className='w-fit cursor-pointer text-sm flex items-center no-underline hover:no-underline text-gray-700'
            >
              <span
                className='hover:underline'
                data-test='org-about-domain-filled'
              >
                {domainGroup.primaryDomain.domain}
              </span>
            </a>
            <div
              className={cn(
                'flex opacity-0 group-hover:opacity-100 gap-1 items-center',
                {
                  '!opacity-100': showMenus[domainGroup.primaryDomain.domain],
                },
              )}
            >
              <div>
                {domainGroup.subdomains?.length > 0 && (
                  <Button
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
                    leftIcon={
                      expandedDomains[domainGroup.primaryDomain.domain] ? (
                        <ChevronCollapse />
                      ) : (
                        <ChevronExpand />
                      )
                    }
                  >
                    {domainGroup.subdomains?.length > 0 && (
                      <span className='text-xs font-medium mr-1'>
                        +{domainGroup.subdomains.length}
                      </span>
                    )}
                  </Button>
                )}
              </div>

              <Menu
                onOpenChange={(isOpen) =>
                  handleMenuOpen(domainGroup.primaryDomain.domain, isOpen)
                }
              >
                <MenuButton asChild>
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
