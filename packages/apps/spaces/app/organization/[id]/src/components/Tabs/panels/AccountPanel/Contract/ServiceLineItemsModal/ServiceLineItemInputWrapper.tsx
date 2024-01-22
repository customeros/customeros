'use client';
import React, { FC, ReactNode } from 'react';

import { Box } from '@ui/layout/Box';

export const ServiceLineItemInputWrapper: FC<{
  isDeleted: boolean;
  children: ReactNode;
  width: string | number;
}> = ({ children, isDeleted, width }) => {
  return (
    <Box
      fontSize='sm'
      w={width}
      sx={{
        '&': {
          position: 'relative',
        },
        '&:after': isDeleted
          ? {
              content: '""',
              position: 'absolute',
              left: '0',
              width: '100%',
              zIndex: '1',
              height: '1px',
              bg: 'gray.700',
              top: '50%',
              animation: 'line-animated 0.25s forwards',
            }
          : {},
      }}
    >
      {children}
    </Box>
  );
};
