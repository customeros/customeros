import { defineStyleConfig } from '@chakra-ui/react';

export const Button = defineStyleConfig({
  baseStyle: {
    borderRadius: '0.5rem',
  },
  variants: {
    outline: ({ colorScheme }) => {
      if (colorScheme === 'gray') {
        return {
          bg: 'white',
          color: `${colorScheme}.500`,
          border: `1px solid`,
          borderColor: `${colorScheme}.300`,
          _hover: {
            bg: `${colorScheme}.50`,
            color: `${colorScheme}.700`,
            borderColor: `${colorScheme}.300`,
          },
          _focusVisible: {
            bg: `${colorScheme}.50`,
            color: `${colorScheme}.700`,
            borderColor: `${colorScheme}.300`,
            boxShadow: `0 0 0 4px var(--chakra-colors-${colorScheme}-50)`,
          },
          _active: {
            bg: `${colorScheme}.50`,
            color: `${colorScheme}.700`,
            borderColor: `${colorScheme}.300`,
          },
        };
      }

      return {
        bg: `${colorScheme}.50`,
        color: `${colorScheme}.700`,
        border: `1px solid`,
        borderColor: `${colorScheme}.200`,
        _hover: {
          bg: `${colorScheme}.100`,
          color: `${colorScheme}.700`,
          borderColor: `${colorScheme}.200`,
        },
        _focusVisible: {
          bg: `${colorScheme}.100`,
          color: `${colorScheme}.700`,
          borderColor: `${colorScheme}.300`,
          boxShadow: `0 0 0 4px var(--chakra-colors-${colorScheme}-100)`,
        },
        _active: {
          bg: `${colorScheme}.100`,
          color: `${colorScheme}.700`,
          borderColor: `${colorScheme}.200`,
        },
      };
    },
    link: () => {
      return {
        color: 'gray.500',
        fontWeight: 'normal',
        _hover: {
          color: 'primary.700',
        },
        _focusVisible: {
          color: 'primary.700',
        },
        _active: {
          color: 'primary.700',
        },
      };
    },
  },
  defaultProps: {},
});
