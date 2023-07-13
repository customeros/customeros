'use client';

import { Flex } from '@ui/layout/Flex';

export default function OrganizationLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <Flex flexDir='row' gap='4'>
      {children}
    </Flex>
  );
}
