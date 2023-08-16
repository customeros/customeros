'use client';
import React, { FC, useState } from 'react';
import { IconButton } from '@ui/form/IconButton';
import ExternalLink from '@spaces/atoms/icons/ExternalLink';
import Slack from '@spaces/atoms/icons/Slack';

export const ViewInSlackButton: FC<{ url: string }> = ({ url }) => {
  const [hovered, setHovered] = useState(false);

  return (
    <IconButton
      aria-label='View in slack'
      size='xs'
      variant='outline'
      onMouseEnter={() => setHovered(true)}
      onMouseLeave={() => setHovered(false)}
      icon={hovered ? <ExternalLink height={16} /> : <Slack height={16} />}
    />
  );
};
