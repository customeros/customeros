import { createMultiStyleConfigHelpers } from '@chakra-ui/styled-system';

const helpers = createMultiStyleConfigHelpers([
  'icon',
  'label',
  'container',
  'control',
]);

export const Checkbox = helpers.defineMultiStyleConfig({
  baseStyle: {
    control: {
      border: '1px solid',
      borderColor: 'gray.300',
      borderRadius: '4px',
      transition: 'all 0.3s ease',
      _focus: {
        boxShadow: 'ringPrimary',
        borderColor: 'primary.500',
        backgroundColor: 'primary.100',
        _invalid: {
          boxShadow: 'inputInvalid',
          backgroundColor: 'red.100',
        },
      },
      _hover: {
        borderColor: 'primary.300',
        boxShadow: 'ringPrimary',
        backgroundColor: 'primary.100',
        _invalid: {
          borderColor: 'red.300',
          boxShadow: 'inputInvalid',
          backgroundColor: 'red.100',
        },
      },
      _invalid: {
        boxShadow: 'inputInvalid',
      },
      _checked: {
        backgroundColor: 'primary.100',
        _invalid: {
          backgroundColor: 'red.100',
          boxShadow: 'inputInvalid',
          _before: {
            backgroundColor: 'red.500',
          },
          _disabled: {
            boxShadow: 'unset',
            borderColor: 'gray.200',
            '& > div': {
              '& > *': {
                color: 'red.500',
              },
            },
          },
          '& > div': {
            '& > *': {
              color: 'red.500',
            },
          },
        },
        _hover: {
          backgroundColor: 'primary.100',
          borderColor: 'primary.500',
          _invalid: {
            backgroundColor: 'red.100',
            borderColor: 'red.500',
          },
        },
        _before: {
          backgroundColor: 'primary.500',
        },
      },
      _indeterminate: {
        backgroundColor: 'primary.100',
        boxShadow: 'ringPrimary',
        _invalid: {
          backgroundColor: 'red.100',
          boxShadow: 'inputInvalid',
          _before: {
            backgroundColor: 'red.500',
          },
          _disabled: {
            borderColor: 'gray.200',
            boxShadow: 'unset',
            '& > div': {
              '& > *': {
                color: 'red.500',
              },
            },
          },
          '& > div': {
            '& > *': {
              color: 'red.500',
            },
          },
        },
        _hover: {
          backgroundColor: 'primary.100',
          borderColor: 'primary.500',
          _invalid: {
            backgroundColor: 'red.100',
            borderColor: 'red.500',
          },
        },
        _disabled: {
          _hover: {
            borderColor: 'gray.200',
          },
        },
        _before: {
          backgroundColor: 'primary.500',
        },
      },
      _disabled: {
        boxShadow: 'unset',
        borderColor: 'gray.200',
        _focus: {
          boxShadow: 'unset',
        },
        _invalid: {
          boxShadow: 'unset',
          borderColor: 'red.200',
          _hover: {
            boxShadow: 'unset',
            borderColor: 'red.200',
            backgroundColor: 'gray.100',
          },
        },
        _checked: {
          _before: {
            backgroundColor: 'gray.400',
          },
          '& > div': {
            '& > *': {
              color: 'gray.400',
            },
          },
        },
        _indeterminate: {
          _before: {
            backgroundColor: 'gray.400',
          },
          '& > div': {
            '& > *': {
              color: 'gray.400',
            },
          },
        },
        _hover: {
          borderColor: 'gray.200',
          boxShadow: 'unset',
        },
      },
    },
    icon: {
      color: 'primary.500',
    },
  },
  sizes: {
    sm: {
      control: {
        borderRadius: '4px',
      },
    },
    md: {
      control: {
        borderRadius: '4px',
      },
    },
    lg: {
      control: {
        borderRadius: '6px',
      },
    },
  },
  defaultProps: {
    size: 'lg',
    colorScheme: 'primary',
  },
});
