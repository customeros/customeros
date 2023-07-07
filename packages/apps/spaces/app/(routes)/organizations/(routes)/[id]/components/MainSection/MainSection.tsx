'use client';

import { Flex } from '@ui/layout/Flex';

export const MainSection = ({ children }: { children?: React.ReactNode }) => {
  return (
    <Flex
      flex='3'
      h='calc(100vh - 2rem)'
      bg='#FCFCFC'
      borderRadius='2xl'
      shadow='base'
    >
      {children}
    </Flex>
  );
};
