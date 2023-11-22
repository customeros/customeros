import { useField } from 'react-inverted-form';
import React, { PropsWithChildren } from 'react';

import { CheckboxProps } from '@chakra-ui/react';

import { CustomCheckbox } from '@ui/form/Checkbox/CustomCheckbox';

interface FormCheckboxProps extends Omit<CheckboxProps, 'onChange'> {
  name: string;
  formId: string;
}

export const FormCheckbox = ({
  name,
  formId,
  children,
  ...rest
}: PropsWithChildren<FormCheckboxProps>) => {
  const { getInputProps } = useField(name, formId);
  const { value: isChecked, onChange } = getInputProps();

  return (
    <CustomCheckbox
      isChecked={isChecked}
      onChange={onChange}
      top='0'
      size='md'
      zIndex='10'
      {...rest}
    >
      {children}
    </CustomCheckbox>
  );
};
