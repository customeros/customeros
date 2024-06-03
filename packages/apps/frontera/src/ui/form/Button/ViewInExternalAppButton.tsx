import { FC, useState, ReactElement } from 'react';

import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { getExternalUrl } from '@utils/getExternalLink';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { LinkExternal02 } from '@ui/media/icons/LinkExternal02';

export const ViewInExternalAppButton: FC<{
  icon: ReactElement;
  url?: string | null;
}> = ({ url, icon }) => {
  const [hovered, setHovered] = useState(false);

  return (
    <Tooltip label={url ? 'View in Slack' : ''}>
      <IconButton
        aria-label='View in slack'
        size='xxs'
        colorScheme='gray'
        className='absolute right-0'
        isDisabled={!url}
        variant={hovered ? 'ghost' : 'outline'}
        onMouseEnter={() => setHovered(true)}
        onMouseLeave={() => setHovered(false)}
        onClick={(e) => {
          e.preventDefault();
          e.stopPropagation();
          if (url) {
            window.open(getExternalUrl(url), '_blank', 'noopener');
          }
        }}
        icon={hovered ? <LinkExternal02 className='text-gray-500' /> : icon}
      />
    </Tooltip>
  );
};
