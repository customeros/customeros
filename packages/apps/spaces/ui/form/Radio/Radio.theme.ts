import { radioAnatomy } from '@chakra-ui/anatomy';
import { createMultiStyleConfigHelpers } from '@chakra-ui/react';

const { definePartsStyle, defineMultiStyleConfig } =
  createMultiStyleConfigHelpers(radioAnatomy.keys);

const baseStyle = definePartsStyle(({ colorScheme }) => ({
  control: {
    borderColor: 'gray.300',
    bg: `${colorScheme === 'primary' ? 'primary.600' : `${colorScheme}.500`}`,
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
  },
}));

export const radioTheme = defineMultiStyleConfig({ baseStyle });
