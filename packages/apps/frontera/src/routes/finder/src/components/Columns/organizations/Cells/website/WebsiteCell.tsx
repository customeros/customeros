import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { removeTrailingSlash } from '@utils/removeTrailingSlash.ts';
import { getExternalUrl, getFormattedLink } from '@utils/getExternalLink';

interface WebsiteCellProps {
  organizationId: string;
}

export const WebsiteCell = observer(({ organizationId }: WebsiteCellProps) => {
  const store = useStore();
  const [metaKey, setMetaKey] = useState(false);
  const organization = store.organizations.getById(organizationId);
  const enrichedOrganizations = organization?.value?.enrichDetails;

  const enrichingStatus =
    !enrichedOrganizations?.enrichedAt &&
    enrichedOrganizations?.requestedAt &&
    !enrichedOrganizations?.failedAt;

  const website = organization?.value?.website;

  const formattedLink = website && getFormattedLink(website);

  return (
    <div
      className='flex items-center cursor-pointer'
      onKeyUp={() => metaKey && setMetaKey(false)}
      onKeyDown={(e) => {
        if (e.metaKey) {
          setMetaKey(true);
        }
      }}
      onClick={(e) => {
        if (e.metaKey) {
          e.stopPropagation();
          window.open(getExternalUrl(website ?? '/'), '_blank', 'noopener');

          return;
        }
        store.ui.commandMenu.setType('RenameOrganizationProperty');
        store.ui.commandMenu.setContext({
          ...store.ui.commandMenu.context,
          property: 'website',
        });
        store.ui.commandMenu.setOpen(true);
      }}
    >
      <p className='text-gray-700  truncate'>
        {website?.length && formattedLink ? (
          removeTrailingSlash(formattedLink)
        ) : enrichingStatus ? (
          <span className='text-gray-400'>Enriching...</span>
        ) : (
          <span className='text-gray-400'>Not set</span>
        )}
      </p>
    </div>
  );
});
