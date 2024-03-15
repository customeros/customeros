import { createMultiStyleConfigHelpers } from '@chakra-ui/styled-system';

const helpers = createMultiStyleConfigHelpers([
  'icon',
  'label',
  'container',
  'control',
]);

export const Checkbox = helpers.defineMultiStyleConfig({
  baseStyle: ({ colorScheme, isInvalid, isDisabled }) => ({
    control: {
      border: '1px solid',
      borderColor: 'gray.300',
      borderRadius: '4px',
      transition: 'all 0.3s ease',
      pointerEvents: isDisabled ? 'none' : 'auto',
      opacity: isDisabled ? 0.4 : 1,
      _focus: {
        boxShadow: 'ringPrimary',
        borderColor: `${colorScheme}.300`,
        backgroundColor: 'white',
        _invalid: {
          boxShadow: 'inputInvalid',
          backgroundColor: 'red.100',
        },
      },
      _hover: {
        backgroundColor: isInvalid ? 'warning.100' : `${colorScheme}.100`,
        borderColor: isInvalid ? 'warning.600' : `${colorScheme}.600`,
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
        backgroundColor: isInvalid ? 'warning.50' : `${colorScheme}.50`,
        borderColor: isInvalid ? 'warning.600' : `${colorScheme}.600`,
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
          backgroundColor: isInvalid ? 'warning.100' : `${colorScheme}.100`,
          borderColor: isInvalid ? 'warning.600' : `${colorScheme}.600`,
          _invalid: {
            backgroundColor: 'red.100',
            borderColor: 'red.500',
          },
        },
        _before: {
          backgroundColor: `${colorScheme}.600`,
        },
      },
      _indeterminate: {
        backgroundColor: `${colorScheme}.100`,
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
          backgroundColor: isInvalid ? 'warning.100' : `${colorScheme}.100`,
          borderColor: isInvalid ? 'warning.600' : `${colorScheme}.600`,
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
          backgroundColor: `${colorScheme}.600`,
        },
      },
      _disabled: {
        boxShadow: 'unset',
        borderColor: 'gray.300',
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
          borderColor: 'gray.300',
          boxShadow: 'unset',
        },
      },
    },
    label: {
      pointerEvents: isDisabled ? 'none' : 'auto',
      opacity: isDisabled ? 0.4 : 1,
      _disabled: {
        color: 'gray.500',
      },
    },
    icon: {
      color: isInvalid ? 'warning.600' : `${colorScheme}.600`, // workaround, _invalid styles are not picked up
    },
  }),
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
