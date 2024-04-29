import { Input, InputProps } from '../Input';

interface NumberInputProps extends InputProps {}

export const NumberInput = ({ ...props }: NumberInputProps) => {
  return <Input {...props} type='number' />;
};
