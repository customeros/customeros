import { ReactElement } from 'react';

import { cn } from '@ui/utils/cn.ts';
import { IconButton } from '@ui/form/IconButton';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';

interface ToolbarButtonProps {
  label: string;
  isActive?: boolean;
  icon: ReactElement;
  onClick: () => void;
}

export const ToolbarButton = ({
  onClick,
  isActive,
  icon,
  label,
}: ToolbarButtonProps) => {
  return (
    <Tooltip
      hasArrow
      label={label}
      side='bottom'
      align='center'
      className='customeros-remirror-tooltip'
    >
      <IconButton
        icon={icon}
        size={'xs'}
        variant='ghost'
        onClick={onClick}
        aria-label={label}
        className={cn('bg-transparent', {
          'text-gray-400': !isActive,
        })}
      />
    </Tooltip>
  );
};
