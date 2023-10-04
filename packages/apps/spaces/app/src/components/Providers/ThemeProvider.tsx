'use client';
import { ChakraProvider } from '@chakra-ui/react';
import { CacheProvider } from '@chakra-ui/next-js';

import { theme } from '@ui/theme/theme';

export const ThemeProvider = ({ children }: { children: React.ReactNode }) => {
  return (
    <CacheProvider>
      <ChakraProvider theme={theme}>{children}</ChakraProvider>
    </CacheProvider>
  );
};
