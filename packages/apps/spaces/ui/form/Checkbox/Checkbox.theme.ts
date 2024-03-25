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
      _invalid: {
        borderColor: 'gray.400',
      },
      _focus: {
        boxShadow: 'ringPrimary',
        borderColor: `${colorScheme}.300`,
        backgroundColor: 'white',
        _invalid: {
          boxShadow: 'ringWarning',
          backgroundColor: 'warning.100',
          borderColor: 'warning.300',
        },
      },
      _hover: {
        backgroundColor: `${colorScheme}.100`,
        borderColor: `${colorScheme}.600`,
        _invalid: {
          borderColor: 'warning.300',
          backgroundColor: 'warning.100',
        },
      },

      _checked: {
        backgroundColor: `${colorScheme}.50`,
        borderColor: `${colorScheme}.600`,
        _invalid: {
          backgroundColor: 'warning.50',
          borderColor: 'warning.600',
          _before: {
            backgroundColor: 'warning.500',
          },
          _disabled: {
            boxShadow: 'unset',
            borderColor: 'gray.200',
            '& > div': {
              '& > *': {
                color: 'warning.500',
              },
            },
          },
          '& > div': {
            '& > *': {
              color: 'warning.500',
            },
          },
        },
        _hover: {
          backgroundColor: isInvalid ? 'warning.100' : `${colorScheme}.100`,
          borderColor: isInvalid ? 'warning.600' : `${colorScheme}.600`,
          _invalid: {
            backgroundColor: 'warning.50',
            borderColor: 'warning.500',
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
          backgroundColor: 'warning.100',
          boxShadow: 'ringWarning',
          _before: {
            backgroundColor: 'warning.500',
          },
          _disabled: {
            borderColor: 'gray.200',
            boxShadow: 'unset',
            '& > div': {
              '& > *': {
                color: 'warning.500',
              },
            },
          },
          '& > div': {
            '& > *': {
              color: 'warning.500',
            },
          },
        },
        _hover: {
          backgroundColor: `${colorScheme}.100`,
          borderColor: `${colorScheme}.600`,
          _invalid: {
            backgroundColor: 'warning.100',
            borderColor: 'warning.500',
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
          borderColor: 'warning.200',
          _hover: {
            boxShadow: 'unset',
            borderColor: 'warning.200',
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
