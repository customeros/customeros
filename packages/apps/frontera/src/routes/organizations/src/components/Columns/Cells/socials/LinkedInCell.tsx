import { useState } from 'react';

import { IconButton } from '@ui/form/IconButton/IconButton';
import { Social } from '@shared/types/__generated__/graphql.types';
import { LinkExternal02 } from '@ui/media/icons/LinkExternal02.tsx';
import {
  getExternalUrl,
  getFormattedLink,
} from '@spaces/utils/getExternalLink';

interface SocialsCellProps {
  socials?: Social[] | null;
}

export const LinkedInCell = ({ socials }: SocialsCellProps) => {
  const [isHovered, setIsHovered] = useState(false);

  if (!socials?.length) return <p className='text-gray-400'>Unknown</p>;
  const linkedIn = socials.find((social) => social?.url?.includes('linkedin'));
  if (!linkedIn?.url) {
    return <p className='text-gray-400'>Unknown</p>;
  }

  const formattedLink = getFormattedLink(linkedIn.url).replace(
    /^linkedin\.com\//,
    '',
  );

  return (
    <div
      className='flex items-center'
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      <p className='text-gray-700 cursor-default truncate'>{formattedLink}</p>
      {isHovered && (
        <IconButton
          className='ml-1 rounded-[5px]'
          variant='ghost'
          size='xs'
          onClick={() =>
            window.open(
              getExternalUrl(linkedIn.url ?? '/'),
              '_blank',
              'noopener',
            )
          }
          aria-label='organization website'
          icon={<LinkExternal02 className='text-gray-500' />}
        />
      )}
    </div>
  );
};
