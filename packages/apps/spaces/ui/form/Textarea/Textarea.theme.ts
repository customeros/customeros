import { defineStyleConfig } from '@chakra-ui/react';

export const Textarea = defineStyleConfig({
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
