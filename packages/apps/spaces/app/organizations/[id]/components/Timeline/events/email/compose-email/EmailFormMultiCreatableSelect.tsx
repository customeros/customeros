import React, { forwardRef } from 'react';

import { SelectInstance } from '@ui/form/SyncSelect/Select';
import { useField } from 'react-inverted-form';
import { AsyncCreatableProps } from '@ui/form/SyncSelect';
import { emailRegex } from '@organization/components/Timeline/events/email/utils';
import { MultiCreatableSelect } from '@ui/form/MultiCreatableSelect';

interface FormSelectProps extends AsyncCreatableProps<any, any, any> {
  name: string;
  formId: string;
  customStyles?: any;
  withTooltip?: boolean;
  Option?: any;
}

export const EmailFormMultiCreatableSelect = forwardRef<
  SelectInstance,
  FormSelectProps
>(({ name, formId, ...rest }, ref) => {
  const { getInputProps } = useField(name, formId);
  const { id, onChange, onBlur, value } = getInputProps();
  const handleBlur = (stringVal: string) => {
    if (stringVal && emailRegex.test(stringVal)) {
      onBlur([...value, { label: stringVal, value: stringVal }]);
      return;
    }
    onBlur(value);
  };

  return (
    <MultiCreatableSelect
      ref={ref}
      id={id}
      formId={formId}
      name={name}
      value={value}
      onBlur={(e) => handleBlur(e.target.value)}
      onChange={onChange}
      {...rest}
    />
  );
});
