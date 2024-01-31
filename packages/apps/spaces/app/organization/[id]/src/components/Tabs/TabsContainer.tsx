'use client';

import { AnimatePresence } from 'framer-motion';

import { Flex } from '@ui/layout/Flex';

export const TabsContainer = ({ children }: { children?: React.ReactNode }) => {
  return (
    <AnimatePresence mode='wait'>
      <Flex
        w='full'
        h='100%'
        bg='gray.25'
        flexDir='column'
        borderRight='1px solid'
        borderColor='gray.200'
        overflow='hidden'
      >
        {children}
      </Flex>
    </AnimatePresence>
  );
};
