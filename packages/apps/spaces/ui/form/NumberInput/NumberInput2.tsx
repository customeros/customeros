import { Input, InputProps } from '../Input/Input2';

interface NumberInputProps extends InputProps {}

export const NumberInput = ({ ...props }: NumberInputProps) => {
  return <Input {...props} type='number' />;
};
