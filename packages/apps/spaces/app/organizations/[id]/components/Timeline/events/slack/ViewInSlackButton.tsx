'use client';
import React, { FC, useState } from 'react';
import { IconButton } from '@ui/form/IconButton';
import ExternalLink from '@spaces/atoms/icons/ExternalLink';
import Slack from '@spaces/atoms/icons/Slack';
import { Tooltip } from '@ui/overlay/Tooltip';

export const ViewInSlackButton: FC<{ url: string }> = ({ url }) => {
  const [hovered, setHovered] = useState(false);

  return (
    <Tooltip label='View in Slack' variant='dark' hasArrow>
      <IconButton
        aria-label='View in slack'
        size='xs'
        variant={hovered ? 'ghost' : 'outline'}
        onMouseEnter={() => setHovered(true)}
        onMouseLeave={() => setHovered(false)}
        icon={
          hovered ? (
            <ExternalLink height={16} color='var(--chakra-colors-gray-500)' />
          ) : (
            <Slack height={16} />
          )
        }
      />
    </Tooltip>
  );
};
