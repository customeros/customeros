'use client';
import { useState, ReactNode, useEffect } from 'react';

import { Flex } from '@ui/layout/Flex';
import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { Box, BoxProps } from '@ui/layout/Box';

interface OrganizationPanelProps extends BoxProps {
  title: string;
  bgImage?: string;
  withFade?: boolean;
  actionItem?: ReactNode;
  shouldBlockPanelScroll?: boolean; // fix for https://linear.app/customer-os/issue/COS-619/scrollbar-overlaps-the-renewal-modals-in-safari
}
export const OrganizationPanel = ({
  bgImage,
  title,
  actionItem,
  children,
  withFade = false,
  shouldBlockPanelScroll = false,
  ...props
}: OrganizationPanelProps) => {
  const [isMounted, setIsMounted] = useState(!withFade);

  useEffect(() => {
    if (!withFade) return;
    setIsMounted(true);
  }, []);

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
      {...props}
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
        overflowY={shouldBlockPanelScroll ? 'hidden' : 'auto'}
        px='6'
        pb={8}
        opacity={isMounted ? 1 : 0}
        transition='opacity 0.3s ease-in-out'
      >
        {children}
      </VStack>
    </Box>
  );
};
