import React, { PropsWithChildren } from 'react';

import { Box, BoxProps } from '@ui/layout/Box';

export const SvgIcon: React.FC<PropsWithChildren & BoxProps> = ({
  children,
  ...props
}) => {
  return <Box {...props}>{children}</Box>;
};
