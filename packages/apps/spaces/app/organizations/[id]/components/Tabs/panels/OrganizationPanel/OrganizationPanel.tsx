'use client';
import { FC, PropsWithChildren, ReactNode } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Box } from '@ui/layout/Box';
import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';

interface OrganizationPanelProps extends PropsWithChildren {
  title: string;
  bgImage?: string;
  actionItem?: ReactNode;
}
export const OrganizationPanel: FC<OrganizationPanelProps> = ({
  bgImage,
  title,
  actionItem,
  children,
}) => {
  return (
    <Box
      p={0}
      flex={1}
      as={Flex}
      flexDirection='column'
      height='100%'
      backgroundImage={bgImage ? bgImage : ''}
      backgroundRepeat='no-repeat'
      backgroundSize='contain'
    >
      <Flex justify='space-between' pt='4' pb='4' px='6'>
        <Text fontSize='lg' color='gray.700' fontWeight='semibold'>
          {title}
        </Text>
        {actionItem && actionItem}
      </Flex>

      <VStack
        spacing='2'
        w='full'
        h='100%'
        justify='stretch'
        overflowY='auto'
        px='6'
        pb={8}
      >
        {children}
      </VStack>
    </Box>
  );
};
