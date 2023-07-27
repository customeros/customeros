import { defineStyleConfig } from '@chakra-ui/react';

export const Button = defineStyleConfig({
  baseStyle: {
    borderRadius: '0.5rem',
  },
  variants: {
    outline: ({ colorScheme }) => ({
      color: `${colorScheme}.700`,
      borderColor: `${colorScheme}.300`,
      _hover: {
        bg: `${colorScheme}.50`,
      },
      _focusVisible: {
        boxShadow: `0 0 0 4px var(--chakra-colors-${colorScheme}-100)`,
      },
    }),
  },
  defaultProps: {},
});
