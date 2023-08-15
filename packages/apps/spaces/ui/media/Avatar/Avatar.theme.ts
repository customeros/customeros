import { createMultiStyleConfigHelpers } from '@chakra-ui/styled-system';

const helpers = createMultiStyleConfigHelpers([
  'badge',
  'excessLabel',
  'container',
  'label',
]);

export const Avatar = helpers.defineMultiStyleConfig({
  baseStyle: () => ({
    container: {
      bg: 'primary.100',
      border: '1px solid transparent',
      // using & selector to work around borderColor bug
      '&': {
        borderColor: 'primary.200',
      },
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
