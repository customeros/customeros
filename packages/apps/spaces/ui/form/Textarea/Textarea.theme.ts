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
        _focusVisible: {
          borderColor: 'teal.500',
          boxShadow: 'unset',
        },
      },
    },
  },
  defaultProps: {
    variant: 'flushed',
  },
});
