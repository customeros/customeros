import React, { forwardRef } from 'react';
import { useField } from 'react-inverted-form';

import { AsyncCreatableProps } from '@ui/form/SyncSelect';
import { SelectInstance } from '@ui/form/SyncSelect/Select';
import { MultiCreatableSelect } from '@ui/form/MultiCreatableSelect';
import { emailRegex } from '@organization/src/components/Timeline/events/email/utils';

interface FormSelectProps extends AsyncCreatableProps<any, any, any> {
  name: string;
  Option?: any;
  formId: string;
  customStyles?: any;
  withTooltip?: boolean;
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
