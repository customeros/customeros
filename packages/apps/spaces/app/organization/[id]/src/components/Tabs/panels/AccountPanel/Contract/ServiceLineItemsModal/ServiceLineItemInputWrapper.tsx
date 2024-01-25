'use client';
import React, { FC, ReactNode } from 'react';

import { Box } from '@ui/layout/Box';

export const ServiceLineItemInputWrapper: FC<{
  isDeleted: boolean;
  children: ReactNode;
  width: string | number;
}> = ({ children, isDeleted, width }) => {
  return (
    <Box fontSize='sm' w={width} pointerEvents={isDeleted ? 'none' : 'auto'}>
      {children}
    </Box>
  );
};
