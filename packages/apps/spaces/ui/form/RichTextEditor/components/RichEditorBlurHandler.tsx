import { FC } from 'react';
import { useField } from 'react-inverted-form';

import { useEditorEvent } from '@remirror/react';

export const RichEditorBlurHandler: FC<{
  name: string;
  formId: string;
}> = ({ formId, name }) => {
  const { getInputProps } = useField(name, formId);
  const { value, onBlur } = getInputProps();
  useEditorEvent('blur', () => {
    onBlur(value);
  });

  return <div />;
};
