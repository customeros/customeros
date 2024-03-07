import React, { ReactNode } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Tooltip } from '@ui/overlay/Tooltip';
import { AlertSquare } from '@ui/media/icons/AlertSquare';
import { PriorityLow } from '@ui/media/icons/PriorityLow';
import { PriorityHigh } from '@ui/media/icons/PriorityHigh';
import { PriorityMedium } from '@ui/media/icons/PriorityMedium';
import { DotsHorizontal } from '@ui/media/icons/DotsHorizontal';

export type Priority = 'low' | 'medium' | 'high' | 'normal' | 'urgent';

interface PriorityBadgeProps {
  priority: Priority;
}

const colorMap: Record<Priority, ReactNode> = {
  normal: <DotsHorizontal />,
  low: <PriorityLow />,
  medium: <PriorityMedium />,
  high: <PriorityHigh />,
  urgent: (
    <AlertSquare
      display='block'
      color='red.600'
      role='presentation'
      boxSize='5'
    />
  ),
};

export const PriorityBadge: React.FC<PriorityBadgeProps> = ({ priority }) => {
  return (
    <Tooltip label={`${priority} priority`} textTransform='capitalize'>
      <Flex alignItems='flex-end' role='presentation' aria-label={priority}>
        {colorMap[priority]}
      </Flex>
    </Tooltip>
  );
};
