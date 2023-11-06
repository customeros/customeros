export type { ImageProps } from '@chakra-ui/next-js';
import NextImage, { ImageProps as NextImageProps } from 'next/image';

import { chakra, ChakraComponent } from '@chakra-ui/react';
// export { Image } from '@chakra-ui/next-js';

export const Image: ChakraComponent<'img', NextImageProps> = chakra(NextImage, {
  shouldForwardProp: (prop) =>
    [
      'src',
      'alt',
      'width',
      'height',
      'fill',
      'loader',
      'quality',
      'priority',
      'loading',
      'placeholder',
      'blurDataURL',
      'unoptimized',
      'onLoadingComplete',
    ].includes(prop),
});
