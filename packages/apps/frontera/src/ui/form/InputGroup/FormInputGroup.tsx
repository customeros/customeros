import { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import { Input } from '../Input/Input';
import { FormInputProps } from '../Input/FormInput';
import { InputGroup, LeftElement, RightElement } from './InputGroup';

interface FormInputGroupProps extends FormInputProps {
  name: string;
  formId: string;
  label?: string;
  autoFocus?: boolean;
  isLabelVisible?: boolean;
  leftElement?: React.ReactNode;
  rightElement?: React.ReactNode;
  labelProps?: React.HTMLAttributes<HTMLLabelElement>;
}

export const FormInputGroup = forwardRef((props: FormInputGroupProps, ref) => {
  const {
    name,
    formId,
    label,
    isLabelVisible,
    leftElement,
    rightElement,
    labelProps,
    autoFocus,
    value: _,
    ...rest
  } = props;
  const { getInputProps } = useField(name, formId);

  return (
    <div>
      <label {...labelProps}>{label}</label>

      <InputGroup {...rest}>
        {leftElement && <LeftElement>{leftElement}</LeftElement>}
        <Input
          ref={ref as React.Ref<HTMLInputElement>}
          {...getInputProps()}
          autoComplete='off'
          {...rest}
          variant='group'
          autoFocus={autoFocus}
        />
        {rightElement && <RightElement>{rightElement}</RightElement>}
      </InputGroup>
    </div>
  );
});
