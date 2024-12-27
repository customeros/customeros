import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { getExternalUrl, getFormattedLink } from '@utils/getExternalLink';

interface SocialsCellProps {
  organizationId: string;
}

export const OrganizationLinkedInCell = observer(
  ({ organizationId }: SocialsCellProps) => {
    const store = useStore();
    const [metaKey, setMetaKey] = useState(false);
    const organization = store.organizations.getById(organizationId);

    const linkedIn = organization?.value?.socialMedia.find((social) =>
      social.url.includes('linkedin'),
    );

    if (organization?.isEnriching && !linkedIn) {
      return <span className='text-gray-400'>Enriching...</span>;
    }

    if (!linkedIn) {
      return <span className='text-gray-400'>Not set</span>;
    }

    const link = linkedIn.url;
    const alias = linkedIn.alias;
    const formattedLink = getFormattedLink(link).replace(
      /^linkedin\.com\/(?:in\/|company\/)?/,
      '/',
    );

    const displayLink = alias ? `/${alias}` : formattedLink;
    const url = formattedLink
      ? link.includes('linkedin')
        ? getExternalUrl(`https://linkedin.com/company${displayLink}`)
        : getExternalUrl(link)
      : '';

    return (
      <div
        className='flex items-center cursor-pointer'
        onKeyUp={() => metaKey && setMetaKey(false)}
        onKeyDown={(e) => e.metaKey && setMetaKey(true)}
        onClick={(e) => {
          if (e.metaKey) {
            e.stopPropagation();
            window.open(url, '_blank', 'noopener');

            return;
          }
          store.ui.commandMenu.setType('EditCompanyLinkedin');
          store.ui.commandMenu.setOpen(true);
        }}
      >
        <p className='text-gray-700 truncate'>{displayLink}</p>
      </div>
    );
  },
);
