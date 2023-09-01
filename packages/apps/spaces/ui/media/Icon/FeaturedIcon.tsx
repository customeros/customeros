import { isValidElement, cloneElement, ReactNode } from 'react';

import { useToken } from '@ui/utils';
import { ButtonProps } from '@ui/form/Button';
import { Flex, FlexProps } from '@ui/layout/Flex';

import { IconProps } from './Icons';

type FeaturedIconSize = 'sm' | 'md' | 'lg' | 'xl';

interface FeaturedIconProps extends Omit<FlexProps, 'children'> {
  size?: FeaturedIconSize;
  colorScheme?: ButtonProps['colorScheme'];
  children: ReactNode | ((props: IconProps) => ReactNode);
}

const boxProps: Record<FeaturedIconSize, FlexProps> = {
  sm: {
    w: '8',
    h: '8',
  },
  md: {
    w: '10',
    h: '10',
  },
  lg: {
    w: '12',
    h: '12',
  },
  xl: {
    w: '14',
    h: '14',
  },
};

const getIconProps = (
  size: FeaturedIconSize,
  color50: string,
  color100: string,
  color600: string,
) => {
  const props = {
    sm: {
      boxSize: '4',
      boxShadow: `0 0 0 6px ${color100}, 0 0 0 10px ${color50}`,
    },
    md: {
      boxSize: '5',
      boxShadow: `0 0 0 7px ${color100}, 0 0 0 13px ${color50}`,
    },
    lg: {
      boxSize: '6',
      boxShadow: `0 0 0 8px ${color100}, 0 0 0 16px ${color50}`,
    },
    xl: {
      boxSize: '7',
      boxShadow: `0 0 0 9px ${color100}, 0 0 0 19px ${color50}`,
    },
  };

  return {
    ...props[size],
    borderRadius: 'full',
    bg: color100,
    color: color600,
  };
};

export const FeaturedIcon = ({
  size = 'md',
  children,
  colorScheme = 'gray',
  ...props
}: FeaturedIconProps) => {
  const [color50, color100, color500] = useToken('colors', [
    `${colorScheme}.50`,
    `${colorScheme}.100`,
    `${colorScheme}.500`,
  ]);
  const iconProps = getIconProps(size, color50, color100, color500);

  const Icon = isValidElement(children)
    ? // eslint-disable-next-line @typescript-eslint/no-explicit-any
      cloneElement<any>(children, { ...iconProps, overflow: 'visible' })
    : typeof children === 'function'
    ? children(iconProps)
    : children;

  return (
    <Flex justify='center' align='center' {...boxProps[size]} {...props}>
      {Icon}
    </Flex>
  );
};
