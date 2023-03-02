import React from 'react';
import { Tooltip as PrimereactTooltip, TooltipProps } from 'primereact/tooltip';

export const Tooltip: React.FC<TooltipProps> = ({ children, ...rest }) => {
  return <PrimereactTooltip {...rest}>{children}</PrimereactTooltip>;
};
