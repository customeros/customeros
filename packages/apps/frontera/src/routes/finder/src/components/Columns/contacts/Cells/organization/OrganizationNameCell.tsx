import { useRef } from 'react';
import { Link } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';

import { Edit03 } from '@ui/media/icons/Edit03';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';

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
      <div className='flex items-center gap-2 group/orgName'>
        <span className='inline truncate'>
          {org.length ? (
            <Link
              to={href}
              ref={linkRef}
              className='inline text-gray-700 no-underline hover:no-underline font-normal cursor-pointer'
            >
              <span
                onClick={() => {}}
                className='inline text-gray-700 no-underline hover:no-underline font-normal cursor-pointer'
              >
                {contactStore?.value.primaryOrganizationName}
              </span>
            </Link>
          ) : (
            <span className='text-gray-400'>None</span>
          )}
        </span>

        <IconButton
          size='xxs'
          variant='ghost'
          icon={<Edit03 />}
          className='opacity-0 group-hover/orgName:opacity-100 mt-[3px]'
          aria-label={`navigate-to-${contactStore?.value.primaryOrganizationName}`}
          onClick={() => {
            store.ui.commandMenu.setType('EditLatestOrgActive');
            store.ui.commandMenu.setOpen(true);
          }}
        />
      </div>
    );
  },
);

function getHref(id: string, lastPositionParams: string | undefined) {
  return `/organization/${id}?${lastPositionParams || 'tab=people'}`;
}
