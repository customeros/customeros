import { useRef } from 'react';
import { Link } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';

import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { LinkExternal02 } from '@ui/media/icons/LinkExternal02';

interface OrganizationNameCellProps {
  org: string;
  orgId: string;
  contactId: string;
}
export const OrganizationNameCell = observer(
  ({ org, contactId, orgId }: OrganizationNameCellProps) => {
    const store = useStore();

    const [tabs] = useLocalStorage<{
      [key: string]: string;
    }>(`customeros-player-last-position`, { root: 'organization' });
    const linkRef = useRef<HTMLAnchorElement>(null);

    const lastPositionParams = tabs[orgId];
    const href = getHref(orgId, lastPositionParams);
    const contactStore = store.contacts.value.get(contactId);

    const isEnriching = contactStore?.isEnriching;

    if (!org?.length && isEnriching) {
      return (
        <p className='text-gray-400'>
          {isEnriching ? 'Enriching...' : 'Not set'}
        </p>
      );
    }

    return (
      <div className='flex items-center gap-2 group'>
        <span className='inline truncate'>
          {org.length ? (
            <span
              className='inline text-gray-700 no-underline hover:no-underline font-normal cursor-pointer'
              onClick={() => {
                store.ui.commandMenu.setType('EditLatestOrgActive');
                store.ui.commandMenu.setOpen(true);
              }}
            >
              {contactStore?.value.primaryOrganizationName}
            </span>
          ) : (
            <span className='text-gray-400'>None</span>
          )}
        </span>

        {contactStore?.value.primaryOrganizationName && (
          <Link to={href} ref={linkRef}>
            <IconButton
              size='xxs'
              variant='ghost'
              icon={<LinkExternal02 />}
              className='opacity-0 group-hover:opacity-100'
              aria-label={`navigate-to-${contactStore?.value.primaryOrganizationName}`}
            />
          </Link>
        )}
      </div>
    );
  },
);

function getHref(id: string, lastPositionParams: string | undefined) {
  return `/organization/${id}?${lastPositionParams || 'tab=people'}`;
}
