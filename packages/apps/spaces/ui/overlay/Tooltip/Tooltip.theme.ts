import { defineStyleConfig } from '@chakra-ui/react';

export const Tooltip = defineStyleConfig({
  baseStyle: {
    py: '2',
    px: '3',
    boxShadow: 'lg',
    fontSize: 'md',
    color: 'gray.700',
    borderRadius: 'md',
    bg: 'white',
    border: '1px solid',
    borderColor: 'gray.100',
  },
  variants: {
    dark: {
      bg: 'gray.700',
      border: 'none',
      color: 'white'
    },
  },
});
