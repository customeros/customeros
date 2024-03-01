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
        _hover: {
          borderColor: 'gray.300',
        },
        _focus: {
          borderColor: 'primary.500',
          boxShadow: 'unset',
        },
        _invalid: {
          boxShadow: 'unset',
          borderColor: 'error.500',
        },
      },
    },
  },
  defaultProps: {
    variant: 'flushed',
  },
});
