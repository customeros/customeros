import { useState } from 'react';

import { IconButton } from '@ui/form/IconButton/IconButton';
import { LinkExternal02 } from '@ui/media/icons/LinkExternal02';
import { getExternalUrl, getFormattedLink } from '@utils/getExternalLink';

interface WebsiteCellProps {
  website?: string | null;
}

export const WebsiteCell = ({ website }: WebsiteCellProps) => {
  const [isHovered, setIsHovered] = useState(false);

  if (!website) return <p className='text-gray-400'>Unknown</p>;

  return (
    <div
      className='flex items-center'
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      <p className='text-gray-700 cursor-default truncate'>
        {getFormattedLink(website)}
      </p>
      {isHovered && (
        <IconButton
          className='ml-1 rounded-[5px]'
          variant='ghost'
          size='xs'
          onClick={() =>
            window.open(getExternalUrl(website ?? '/'), '_blank', 'noopener')
          }
          aria-label='organization website'
          icon={<LinkExternal02 className='text-gray-500' />}
        />
      )}
    </div>
  );
};
