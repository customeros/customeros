import { createMultiStyleConfigHelpers } from '@chakra-ui/styled-system';
import { avatarAnatomy } from '@chakra-ui/anatomy';

const { definePartsStyle } = createMultiStyleConfigHelpers(avatarAnatomy.keys);
const helpers = createMultiStyleConfigHelpers([
  'badge',
  'excessLabel',
  'container',
  'label',
]);

const roundedSquare = definePartsStyle({
  badge: {
    bg: 'gray.500',
    border: '2px solid',
  },
  container: {
    borderRadius: 'md',
    boxShadow: 'none',
    bg: 'primary.50',
    '& img': {
      borderRadius: 'md',
    },
  },
  excessLabel: {
    bg: 'gray.800',
    color: 'white',
    borderRadius: 'xl',
    border: '2px solid',

    // let's also provide dark mode alternatives
    _dark: {
      bg: 'gray.400',
      color: 'gray.900',
    },
  },
});

export const Avatar = helpers.defineMultiStyleConfig({
  baseStyle: () => ({
    container: {
      bg: 'primary.100',
      boxShadow: 'avatarRing',
    },
    label: {
      color: 'primary.700',
      fontSize: 'lg',
      fontWeight: 'bold',
    },
  }),
  sizes: {
    lg: {
      container: {
        w: '12',
        h: '12',
      },
    },
  },
  variants: {
    roundedSquare,
    roundedSquareSmall: {
      ...roundedSquare,
      container: {
        ...roundedSquare.container,
        textDecoration: 'none !important',
        borderRadius: 'sm',
        '& img': {
          borderRadius: 'sm',
        },
      },
      _hover: {
        textDecoration: 'none',
      },
      _focusVisible: {
        textDecoration: 'none',
      },
      label: {
        textDecoration: 'none !important',
        fontSize: 'sm',
      },
    },
    shadowed: {
      container: {
        boxShadow: 'avatarRing',
      },
    },
  },
  defaultProps: {
    size: 'lg',
    colorScheme: 'primary',
  },
});
