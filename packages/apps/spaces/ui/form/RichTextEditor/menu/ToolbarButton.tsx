import React, { ReactElement } from 'react';

import { Tooltip } from '@ui/overlay/Tooltip';
import { IconButton } from '@ui/form/IconButton';

interface ToolbarButtonProps {
  label: string;
  isActive?: boolean;
  icon: ReactElement;
  onClick: () => void;
}
export const ToolbarButton: React.FC<ToolbarButtonProps> = ({
  onClick,
  isActive,
  icon,
  label,
}) => {
  return (
    <Tooltip
      label={label}
      className='customeros-remirror-tooltip'
      hasArrow
      placement='bottom'
    >
      <IconButton
        fontSize={2}
        className='customeros-remirror-button'
        bg='transparent'
        variant='ghost'
        aria-label={label}
        onClick={onClick}
        isActive={isActive}
        icon={icon}
      />
    </Tooltip>
  );
};
