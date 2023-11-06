'use client';

import { ButtonProps } from '@ui/form/Button';
import { Flex, FlexProps } from '@ui/layout/Flex';

interface DotProps extends FlexProps {
  colorScheme?: ButtonProps['colorScheme'];
}

export const Dot = ({ colorScheme = 'gray', ...props }: DotProps) => {
  return (
    <Flex
      w='10px'
      h='10px'
      borderRadius='full'
      bg={`${colorScheme}.500`}
      {...props}
    />
  );
};
