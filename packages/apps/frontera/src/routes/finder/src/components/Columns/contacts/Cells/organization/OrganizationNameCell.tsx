import { Link } from 'react-router-dom';
import { useRef, useState } from 'react';

import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';

import { cn } from '@ui/utils/cn';
import { Combobox } from '@ui/form/Combobox';
import { Edit03 } from '@ui/media/icons/Edit03';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Organization } from '@shared/types/__generated__/graphql.types';
import {
  Popover,
  PopoverTrigger,
  PopoverContent,
} from '@ui/overlay/Popover/Popover';

interface OrganizationNameCellProps {
  org: string;
  orgId: string;
  contactId: string;
}
export const OrganizationNameCell = observer(
  ({ org, contactId, orgId }: OrganizationNameCellProps) => {
    const store = useStore();
    const [isOpen, setIsOpen] = useState(false);
    const [isHoverd, setIsHoverd] = useState(false);

    const [tabs] = useLocalStorage<{
      [key: string]: string;
    }>(`customeros-player-last-position`, { root: 'organization' });
    const linkRef = useRef<HTMLAnchorElement>(null);

    const lastPositionParams = tabs[orgId];
    const href = getHref(orgId, lastPositionParams);

    const organizations = store.organizations.toArray();

    const options = organizations.map((org) => ({
      label: org.value.name,
      value: org.value.metadata.id,
    }));

    const contactStore = store.contacts.value.get(contactId);

    return (
      <div
        className='flex items-center'
        onMouseEnter={() => setIsHoverd(true)}
        onMouseLeave={() => setIsHoverd(false)}
      >
        <span className='inline truncate'>
          {org.length ? (
            <Link
              to={href}
              ref={linkRef}
              className='inline text-gray-700 no-underline hover:no-underline font-normal'
            >
              {
                contactStore?.value.latestOrganizationWithJobRole?.organization
                  .name
              }
            </Link>
          ) : (
            <span className='text-gray-400'>None</span>
          )}
        </span>
        <Popover open={isOpen} onOpenChange={(value) => setIsOpen(value)}>
          <PopoverTrigger asChild>
            <IconButton
              size='xxs'
              variant='ghost'
              icon={<Edit03 />}
              aria-label='edit-organization'
              onClick={() => setIsOpen(true)}
              className={cn('opacity-0 ml-2', isHoverd && 'opacity-100')}
            />
          </PopoverTrigger>
          <PopoverContent align='end' side='bottom' className='w-[200px]'>
            <Combobox
              options={options}
              onChange={(value) => {
                contactStore?.value.organizations.content.map((org) => {
                  if (org.metadata.id === value.value) return;
                  contactStore.value.organizations.content.push({
                    metadata: {
                      id: value.value,
                      // eslint-disable-next-line @typescript-eslint/no-explicit-any
                    } as any,
                    id: value.value,
                    name: value.label,
                  } as Organization);
                });
                contactStore?.commit();

                if (contactStore?.value.latestOrganizationWithJobRole) {
                  contactStore.value.latestOrganizationWithJobRole.organization =
                    {
                      metadata: {
                        id: value.value,
                        // eslint-disable-next-line @typescript-eslint/no-explicit-any
                      } as any,
                      id: value.value,
                      name: value.label,
                    } as Organization;
                }
                setIsOpen(false);
                contactStore?.commit({ syncOnly: true });
              }}
            />
          </PopoverContent>
        </Popover>
      </div>
    );
  },
);

function getHref(id: string, lastPositionParams: string | undefined) {
  return `/organization/${id}?${lastPositionParams || 'tab=people'}`;
}
