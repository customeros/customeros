import { createMultiStyleConfigHelpers } from '@chakra-ui/styled-system';

const helpers = createMultiStyleConfigHelpers(['field', 'addon']);

export const Input = helpers.defineMultiStyleConfig({
  baseStyle: {
    field: {
      color: 'gray.700',
      _placeholder: {
        color: 'gray.400',
      },
    },
  },
  variants: {
    flushed: {
      field: {
        borderColor: 'transparent',
        _focus: {
          borderColor: 'primary.500',
          boxShadow: 'unset',
        },
      },
    },
  },
  defaultProps: {
    variant: 'flushed',
  },
});
