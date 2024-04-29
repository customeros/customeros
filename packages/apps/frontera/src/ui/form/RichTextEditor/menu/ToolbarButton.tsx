import React, { ReactElement } from 'react';

import { IconButton } from '@ui/form/IconButton';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';

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
      side='bottom'
      align='center'
    >
      <IconButton
        className='customeros-remirror-button bg-transparent'
        variant='ghost'
        aria-label={label}
        onClick={onClick}
        icon={icon}
        isDisabled={!isActive}
      />
    </Tooltip>
  );
};
