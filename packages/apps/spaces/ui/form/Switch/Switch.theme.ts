import { switchAnatomy } from '@chakra-ui/anatomy';
import { createMultiStyleConfigHelpers } from '@chakra-ui/react';

const { definePartsStyle, defineMultiStyleConfig } =
  createMultiStyleConfigHelpers(switchAnatomy.keys);

const baseStyle = definePartsStyle(({ colorScheme }) => ({
  track: {
    bg: `${colorScheme === 'primary' ? 'primary.600' : `${colorScheme}.500`}`,
  },
}));

export const switchTheme = defineMultiStyleConfig({ baseStyle });
