import { defineStyleConfig } from '@chakra-ui/react';

export const Textarea = defineStyleConfig({
  baseStyle: {
    _placeholder: {
      color: 'gray.400',
    },
  },
  variants: {
    flushed: {
      field: {
        borderColor: 'transparent',
        _focusVisible: {
          boxShadow: 'unset',
        },
      },
    },
  },
  defaultProps: {
    variant: 'flushed',
    colorScheme: 'primary',
  },
});
