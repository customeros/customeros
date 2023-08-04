import { defineStyle, defineStyleConfig } from '@chakra-ui/react';

export const Button = defineStyleConfig({
  baseStyle: {
    borderRadius: '0.5rem',
  },
  variants: {
    // outline: ({ colorScheme }) => ({
    //   color: `${colorScheme}.700`,
    //   borderColor: `${colorScheme}.300`,
    //   _hover: {
    //     bg: `${colorScheme}.50`,
    //   },
    //   _focusVisible: {
    //     boxShadow: `0 0 0 4px var(--chakra-colors-${colorScheme}-100)`,
    //   },
    // }),
    outline: ({ colorScheme }) => ({
      color: `${colorScheme}.700`,
      border: `1px solid`,
      bg: `${colorScheme}.25`,
      borderColor: `${colorScheme}.200`,
      _hover: {
        background: 'primary.50',
        color: 'primary.700',
        borderColor: 'primary.200',
      },
      // _focus: {
      //   background: 'primary.50',
      //   color: 'primary.700',
      //   borderColor: 'primary.200',
      // },
      _focusVisible: {
        background: 'primary.50',
        color: 'primary.700',
        borderColor: 'primary.200',
        boxShadow: '0 0 0 4px var(--chakra-colors-primary-100)',
      },
    }),
  },
  defaultProps: {},
});
