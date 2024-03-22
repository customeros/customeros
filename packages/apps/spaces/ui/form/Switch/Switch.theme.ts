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
      _checked: {
        position: 'relative',
        '&:after': {
          content: '"!"',
          width: '1ch',
          background: 'gray.200',
          color: 'warning.500',
          top: '6px',
          position: 'absolute',
          right: '1px',
          fontSize: 'xs',
          fontWeight: 'bold',
        },
      },
    },
  },
}));

export const switchTheme = defineMultiStyleConfig({ baseStyle });
