'use client';
import React from 'react';

import { Flex } from '@ui/layout/Flex';

export const TabsContainer = ({ children }: { children?: React.ReactNode }) => {
  return (
    <Flex
      minW='400px'
      h='100%'
      bg='gray.25'
      flexDir='column'
      borderRight='1px solid'
      borderColor='gray.200'
      overflow='hidden'
    >
      {children}
    </Flex>
  );
};
