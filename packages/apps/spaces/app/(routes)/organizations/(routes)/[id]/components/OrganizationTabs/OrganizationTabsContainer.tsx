'use client';

import { Flex } from '@ui/layout/Flex';

export const OrganizationTabsContainer = ({
  children,
}: {
  children?: React.ReactNode;
}) => {
  return (
    <Flex
      w='full'
      h='100%'
      bg='white'
      shadow='base'
      flexDir='column'
      borderRadius='2xl'
    >
      {children}
    </Flex>
  );
};
