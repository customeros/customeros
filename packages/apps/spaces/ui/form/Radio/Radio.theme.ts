import { radioAnatomy } from '@chakra-ui/anatomy';
import { createMultiStyleConfigHelpers } from '@chakra-ui/react';

const { definePartsStyle, defineMultiStyleConfig } =
  createMultiStyleConfigHelpers(radioAnatomy.keys);

const baseStyle = definePartsStyle(({ colorScheme }) => ({
  control: {
    borderColor: 'gray.300',
    border: '1px solid',
    color: 'primary.600',
    _hover: {
      borderColor: 'primary.600',
      bg: 'primary.100',
    },
    _focusVisible: {
      boxShadow: '0px 0px 0px 4px #F4EBFF',
      borderColor: 'primary.300',
    },
    _disabled: {
      bg: 'gray.100',
      borderColor: 'gray.300',
    },
    _checked: {
      bg: 'primary.50',
      '&:before': {
        content: '""',
        display: 'block',
        width: 2,
        height: 2,
        bg: 'primary.600',
        borderRadius: 'full',
      },
      _hover: {
        bg: 'primary.100',
      },
      _disabled: {
        borderColor: 'gray.300',
        bg: 'gray.100',
        '&:before': {
          content: '""',
          bg: 'gray.300',
        },
      },
    },
  },
}));

export const radioTheme = defineMultiStyleConfig({ baseStyle });
