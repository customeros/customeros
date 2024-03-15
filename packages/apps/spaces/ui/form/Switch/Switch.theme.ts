import { switchAnatomy } from '@chakra-ui/anatomy';
import { createMultiStyleConfigHelpers } from '@chakra-ui/react';

const { definePartsStyle, defineMultiStyleConfig } =
  createMultiStyleConfigHelpers(switchAnatomy.keys);

const baseStyle = definePartsStyle(({ colorScheme }) => ({
  track: {
    _checked: {
      bg: `${colorScheme === 'primary' ? 'primary.600' : `${colorScheme}.500`}`,
      _invalid: { bg: 'warning.500' },
    },
  },
  thumb: {
    _invalid: {
      position: 'relative',
      '&:after': {
        content: '"!"',
        width: '100%',
        background: 'gray.200',
        color: 'warning.500',
        position: 'absolute',
        top: '6px',
        right: '-5px',
        fontSize: 'xs',
        fontWeight: 'bold',
      },
    },
  },
}));

export const switchTheme = defineMultiStyleConfig({ baseStyle });
