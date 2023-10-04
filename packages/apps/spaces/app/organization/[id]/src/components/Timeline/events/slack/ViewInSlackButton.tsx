'use client';
import React, { FC, useState } from 'react';
import { IconButton } from '@ui/form/IconButton';
import { Link03 } from '@ui/media/icons/Link03';
import Slack from '@spaces/atoms/icons/Slack';
import { Tooltip } from '@ui/overlay/Tooltip';
import { getExternalUrl } from '@spaces/utils/getExternalLink';

export const ViewInSlackButton: FC<{ url?: string | null }> = ({ url }) => {
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
        icon={
          hovered ? (
            <Link03 height={16} color='gray.500' />
          ) : (
            <Slack height={16} />
          )
        }
      />
    </Tooltip>
  );
};
