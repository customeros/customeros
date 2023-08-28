import React, { ReactElement, useRef } from 'react';
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
  const ref = useRef(null);
  return (
    <div ref={ref}>
      <Tooltip
        label={label}
        className='customeros-remirror-tooltip'
        hasArrow
        placement='bottom'
        portalProps={{ containerRef: ref }}
      >
        <IconButton
          className='customeros-remirror-button'
          bg='transparent'
          variant='ghost'
          aria-label={label}
          onClick={onClick}
          isActive={isActive}
          icon={icon}
        />
      </Tooltip>
    </div>
  );
};
