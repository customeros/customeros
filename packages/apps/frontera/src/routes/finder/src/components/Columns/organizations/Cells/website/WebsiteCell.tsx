import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { removeTrailingSlash } from '@utils/removeTrailingSlash';
import { getExternalUrl, getFormattedLink } from '@utils/getExternalLink';

interface WebsiteCellProps {
  organizationId: string;
}

export const WebsiteCell = observer(({ organizationId }: WebsiteCellProps) => {
  const store = useStore();
  const organization = store.organizations.getById(organizationId);

  const website = organization?.value?.website;

  const formattedLink = website && getFormattedLink(website);

  return (
    <div
      className='flex items-center cursor-pointer'
      onClick={(e) => {
        e.stopPropagation();
        window.open(getExternalUrl(website ?? '/'), '_blank', 'noopener');
      }}
    >
      <p className='text-gray-700  truncate'>
        {website?.length && formattedLink ? (
          removeTrailingSlash(formattedLink)
        ) : organization?.isEnriching ? (
          <span
            className='text-gray-400'
            data-test='organization-website-in-all-orgs-table'
          >
            Enriching...
          </span>
        ) : (
          <span
            className='text-gray-400'
            data-test='organization-website-in-all-orgs-table'
          >
            Not set
          </span>
        )}
      </p>
    </div>
  );
});
