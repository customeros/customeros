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
  variants: {},
  defaultProps: {
    size: 'lg',
    colorScheme: 'primary',
  },
});
