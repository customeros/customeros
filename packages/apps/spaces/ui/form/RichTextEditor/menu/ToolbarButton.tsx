import React, { ReactElement } from 'react';
import { IconButton } from '@ui/form/IconButton';
import { Tooltip } from '@ui/overlay/Tooltip';

interface ToolbarButtonProps {
  onClick: () => void;
  isActive?: boolean;
  icon: ReactElement;
  label: string;
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
