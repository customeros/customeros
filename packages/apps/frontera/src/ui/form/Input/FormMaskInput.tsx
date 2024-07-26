import { useIMask } from 'react-imask';
import React, { useEffect } from 'react';
import { useField } from 'react-inverted-form';

import { Input } from '@ui/form/Input/Input';

import { FormInputProps } from './FormInput';

interface FormMaskInputProps extends FormInputProps {
  name: string;
  label: string;
  formId: string;
  labelProps: React.LabelHTMLAttributes<HTMLLabelElement>;
  options: {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    opts: any;
    onAccept?: () => void;
    onComplete?: () => void;
  };
}

/**
 * @deprecated Use `<MaskedInput />` instead
 */
export const FormMaskInput = ({
  labelProps,
  label,
  formId,
  name,
  options: { opts, onAccept, onComplete },
  ...props
}: FormMaskInputProps) => {
  const { ref, setUnmaskedValue } = useIMask(opts, {
    onAccept: onAccept,
    onComplete: onComplete,
  });
  const { getInputProps } = useField(name, formId);
  const { value, onChange } = getInputProps();

  useEffect(() => {
    if (value) {
      setUnmaskedValue(value);
    }
  }, [value]);

  return (
    <div>
      <label {...labelProps}>{label}</label>
      {/* @ts-expect-error-ignore-now*/}
      <Input ref={ref} onChange={onChange} autoComplete='off' {...props} />
    </div>
  );
};
