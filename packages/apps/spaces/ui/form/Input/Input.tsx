import { forwardRef } from 'react';

import {
  Input as ChakraInput,
  InputProps as ChakraInputProps,
} from '@chakra-ui/react';
export { FormLabel, FormControl } from '@chakra-ui/react';
export type { InputProps } from '@chakra-ui/react';

interface Props extends ChakraInputProps {
  disabled?: boolean;
  required?: boolean;
  readOnly?: boolean;
}
export const Input = forwardRef(({ ...props }: Props, ref) => {
  return <ChakraInput ref={ref} {...props} data-1p-ignore />;
});
