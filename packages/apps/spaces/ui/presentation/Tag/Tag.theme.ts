import { defineStyleConfig } from '@chakra-ui/react';

export const TagTheme = defineStyleConfig({
  baseStyle: {},
  variants: {
    outline: ({ colorScheme }) => {
      return {
        container: {
          bg: `${colorScheme}.50`,
          color: `${colorScheme}.700`,
          border: `1px solid`,
          borderColor: `${colorScheme}.200`,
          boxShadow: 'none',
          fontSize: 'sm',
          fontWeight: 'normal',
        },
      };
    },
  },
  defaultProps: {},
});
