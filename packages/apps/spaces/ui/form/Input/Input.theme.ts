import { createMultiStyleConfigHelpers } from '@chakra-ui/styled-system';

const helpers = createMultiStyleConfigHelpers(['field', 'addon']);

export const Input = helpers.defineMultiStyleConfig({
  variants: {
    flushed: {
      field: {
        _focus: {
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
