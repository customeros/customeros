import { useField } from 'react-inverted-form';

import { InputProps } from '@ui/form/Input';
import { UrlInput } from './UrlInput';

interface FormUrlInputProps extends InputProps {
  name: string;
  formId: string;
}

export const FormUrlInput = ({ name, formId, ...props }: FormUrlInputProps) => {
  const { getInputProps } = useField(name, formId);

  return <UrlInput {...getInputProps()} {...props} />;
};
