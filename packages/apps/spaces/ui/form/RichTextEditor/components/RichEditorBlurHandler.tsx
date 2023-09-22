import { FC } from 'react';
import { useEditorEvent } from '@remirror/react';
import { useField } from 'react-inverted-form';

export const RichEditorBlurHandler: FC<{
  formId: string;
  name: string;
}> = ({ formId, name }) => {
  const { getInputProps } = useField(name, formId);
  const { value, onBlur } = getInputProps();
  useEditorEvent('blur', () => {
    onBlur(value);
  });

  return <div />;
};
