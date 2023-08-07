'use client';

import { Flex } from '@ui/layout/Flex';

export const SideSection = ({ children }: { children?: React.ReactNode }) => {
  return (
    <Flex flex='1' h='100%' minW='28rem'>
      {children}
    </Flex>
  );
};
