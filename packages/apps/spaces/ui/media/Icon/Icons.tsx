import React from 'react';

import { IconProps, createIcon } from '@chakra-ui/react';

import { customIcons } from './customIcons';

export const CustomIcons = { ...customIcons };
export type IconNames = keyof typeof CustomIcons;
export type IconsRecord = Record<IconNames, React.FC<IconProps>>;
export type { IconProps };

export const Icons: IconsRecord = Object.entries(CustomIcons).reduce(
  (prev, [name, icon]) => {
    // Custom Icons with multiple paths
    if (Array.isArray(icon)) {
      return {
        ...prev,
        [name]: createIcon({
          defaultProps: {
            fill: 'none',
            boxSize: '4',
          },
          path: icon,
        }),
      };
    }

    // Custom Icons with single path
    return {
      ...prev,
      [name]: createIcon({
        defaultProps: {
          boxSize: '4',
          fill: 'none',
          strokeWidth: '2',
          stroke: 'currentColor',
        },
        path: icon,
      }),
    };
  },
  {} as IconsRecord,
);
