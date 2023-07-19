'use client';

import { Flex } from '@ui/layout/Flex';

export const MainSection = ({ children }: { children?: React.ReactNode }) => {
  return (
    <Flex
      flex='3'
      h='calc(100vh - 1rem)'
      bg='white'
      borderRadius='2xl'
      border='1px solid'
      borderColor='gray.200'
    >
      {children}
    </Flex>
  );
};
