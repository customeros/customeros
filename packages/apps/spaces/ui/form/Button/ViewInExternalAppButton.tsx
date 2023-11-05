import React, { FC, useState, ReactElement } from 'react';

import { Tooltip } from '@ui/overlay/Tooltip';
import { Link03 } from '@ui/media/icons/Link03';
import { IconButton } from '@ui/form/IconButton';
import { getExternalUrl } from '@spaces/utils/getExternalLink';

export const ViewInExternalAppButton: FC<{
  icon: ReactElement;
  url?: string | null;
}> = ({ url, icon }) => {
  const [hovered, setHovered] = useState(false);

  return (
    <Tooltip label={url ? 'View in Slack' : ''} variant='dark' hasArrow>
      <IconButton
        aria-label='View in slack'
        size='xs'
        position='absolute'
        right={0}
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
        icon={hovered ? <Link03 height={16} color='gray.500' /> : icon}
      />
    </Tooltip>
  );
};
