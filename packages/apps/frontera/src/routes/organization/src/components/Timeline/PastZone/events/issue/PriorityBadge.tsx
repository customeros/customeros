import React, { ReactNode } from 'react';

import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
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
    <AlertSquare role='presentation' className='block text-red-600 size-5' />
  ),
};

export const PriorityBadge: React.FC<PriorityBadgeProps> = ({ priority }) => {
  return (
    <Tooltip className='capitalize' label={`${priority} priority`}>
      <div aria-label={priority} className='flex items-end'>
        {colorMap[priority]}
      </div>
    </Tooltip>
  );
};
