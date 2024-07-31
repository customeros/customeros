import { FC } from 'react';
import { useField } from 'react-inverted-form';

import { Input } from '../Input/Input';
import { FormInputProps } from '../Input/FormInput';
import {
  InputGroup,
  LeftElement,
  RightElement,
} from '../InputGroup/InputGroup';

interface FormNumberInputGroupProps extends FormInputProps {
  name: string;
  formId: string;
  label?: string;
  isLabelVisible?: boolean;
  leftElement?: React.ReactNode;
  rightElement?: React.ReactNode;
  labelProps?: React.LabelHTMLAttributes<HTMLLabelElement>;
}

export const FormNumberInputGroup: FC<FormNumberInputGroupProps> = ({
  name,
  formId,
  leftElement,
  rightElement,
  label,
  labelProps,
  ...rest
}) => {
  const { getInputProps } = useField(name, formId);

  return (
    <div>
      <label {...labelProps}>{label}</label>

      <InputGroup>
        {leftElement && <LeftElement>{leftElement}</LeftElement>}

        <Input
          {...rest}
          {...getInputProps()}
          type='number'
          autoComplete='off'
          placeholder={rest?.placeholder || ''}
          className='w-full hover:border-transparent focus:hover:border-transparent focus:border-transparent'
        />

        {rightElement && <RightElement>{rightElement}</RightElement>}
      </InputGroup>
    </div>
  );
};
