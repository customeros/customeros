'use client';
import React from 'react';
import { Flex } from '@ui/layout/Flex';

export const TabsContainer = ({ children }: { children?: React.ReactNode }) => {
  return (
    <Flex
      maxW='50%'
      minW='600px'
      h='100%'
      bg='gray.25'
      flexDir='column'
      border='1px solid'
      borderColor='gray.200'
      borderRadius='2xl'
      overflow='hidden'
    >
      {children}
    </Flex>
  );
};
