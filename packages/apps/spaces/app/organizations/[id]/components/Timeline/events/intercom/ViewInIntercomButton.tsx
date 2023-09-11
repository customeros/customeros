'use client';
import React, { FC, useState } from 'react';
import { IconButton } from '@ui/form/IconButton';
import ExternalLink from '@spaces/atoms/icons/ExternalLink';
import { Tooltip } from '@ui/overlay/Tooltip';
import { getExternalUrl } from '@spaces/utils/getExternalLink';
import Intercom from '@ui/media/icons/Intercom';
import { Box } from '@ui/layout/Box';

export const ViewInIntercomButton: FC<{ url?: string | null }> = ({ url }) => {
  const [hovered, setHovered] = useState(false);

  return (
    <Tooltip label={url ? 'View in Intercom' : ''} variant='dark' hasArrow>
      <IconButton
        aria-label='View in Intercom'
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
            <ExternalLink height={16} color='var(--chakra-colors-gray-500)' />
          ) : (
            <Box boxSize={4}>
              <Intercom />
            </Box>
          )
        }
      />
    </Tooltip>
  );
};
