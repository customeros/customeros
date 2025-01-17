import { Link } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { getExternalUrl, getFormattedLink } from '@utils/getExternalLink';

interface SocialsCellProps {
  organizationId: string;
}

export const OrganizationLinkedInCell = observer(
  ({ organizationId }: SocialsCellProps) => {
    const store = useStore();
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
      <div className='flex items-center cursor-pointer'>
        <Link to={url} target='_blank' className='flex items-center'>
          <p className='text-gray-700 truncate hover:underline'>
            {displayLink}
          </p>
        </Link>
      </div>
    );
  },
);
