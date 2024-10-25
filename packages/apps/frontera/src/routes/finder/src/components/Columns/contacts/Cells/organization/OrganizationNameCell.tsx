import { Link } from 'react-router-dom';
import { useRef, useState } from 'react';

import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';

import { cn } from '@ui/utils/cn';
import { Combobox } from '@ui/form/Combobox';
import { Edit03 } from '@ui/media/icons/Edit03';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
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
          <Link
            to={href}
            ref={linkRef}
            className='inline text-gray-700 no-underline hover:no-underline font-normal'
          >
            {org}
          </Link>
        </span>
        <Popover>
          <PopoverTrigger asChild>
            <IconButton
              size='xxs'
              variant='ghost'
              icon={<Edit03 />}
              aria-label='edit-organization'
              className={cn('opacity-0 ml-2', isHoverd && 'opacity-100')}
            />
          </PopoverTrigger>
          <PopoverContent align='end' side='bottom' className='w-[200px]'>
            <Combobox
              options={options}
              closeMenuOnSelect={true}
              onChange={(value) => {
                contactStore?.linkOrganization(value.value);
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
