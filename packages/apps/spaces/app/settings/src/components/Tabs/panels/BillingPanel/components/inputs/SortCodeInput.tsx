import { useIMask } from 'react-imask';
import React, { useEffect } from 'react';
import { useField } from 'react-inverted-form';

import { twMerge } from 'tailwind-merge';

import { Input } from '@ui/form/Input/Input2';
import { FormInputProps } from '@ui/form/Input/FormInput';

const opts = {
  mask: '00-00-00',
  definitions: {
    '0': /[0-9]/,
  },
};
export const SortCodeInput = ({
  isLabelVisible,
  labelProps,
  label,
  formId,
  name,
  className,
  ...props
}: FormInputProps) => {
  const { ref, setUnmaskedValue } = useIMask<HTMLInputElement>(
    opts /* { onAccept, onComplete } */,
  );
  const { getInputProps } = useField(name, formId);
  const { value, onChange } = getInputProps();

  useEffect(() => {
    if (value) {
      setUnmaskedValue(value);
    }
  }, [value]);

  return (
    <div className={twMerge(className)}>
      {isLabelVisible ? (
        <label {...labelProps}>{label}</label>
      ) : (
        <span>
          <label>{label}</label>
        </span>
      )}

      <Input ref={ref} {...props} onChange={onChange} autoComplete='off' />
    </div>
  );
};
