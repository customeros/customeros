import { useField } from 'react-inverted-form';

import { CurrencyInput, CurrencyInputProps } from './CurrencyInput';

interface FormCurrencyInputProps
  extends Omit<CurrencyInputProps, 'value' | 'onChange'> {
  name: string;
  formId: string;
}

export const FormCurrencyInput = ({
  name,
  formId,
  ...rest
}: FormCurrencyInputProps) => {
  const { getInputProps } = useField(name, formId);
  const { value, onChange, onBlur } = getInputProps();

  return (
    <CurrencyInput
      value={value}
      onChange={onChange}
      onBlur={() => onBlur(value)}
      {...rest}
    />
  );
};
