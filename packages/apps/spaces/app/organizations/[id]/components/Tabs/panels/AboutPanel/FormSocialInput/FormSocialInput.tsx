import { useState, useRef } from 'react';
import { useField } from 'react-inverted-form';

import {
  InputGroup,
  InputGroupProps,
  InputLeftElement,
} from '@ui/form/InputGroup';
import { Input } from '@ui/form/Input';

import { SocialIcon } from './SocialIcons';
import { SocialInput } from './SocialInput';

interface FormSocialInputProps extends InputGroupProps {
  name: string;
  formId: string;
  leftElement?: React.ReactNode;
}

export const FormSocialInput = ({
  name,
  formId,
  leftElement,
  ...rest
}: FormSocialInputProps) => {
  const { getInputProps } = useField(name, formId);
  const { value: values, onChange, onBlur } = getInputProps();

  const newInputRef = useRef<HTMLInputElement>(null);
  const [newValue, setNewValue] = useState('');

  const handleChange =
    (index: number) => (e: React.ChangeEvent<HTMLInputElement>) => {
      const next = [...values];
      next[index] = e.target.value;
      onChange(next);
    };

  const handleBlur =
    (index: number) => (e: React.FocusEvent<HTMLInputElement>) => {
      if (!e.target.value) {
        const next = [...values];
        next.splice(index, 1);
        onBlur?.(next);
      } else {
        onBlur?.(values);
      }
    };

  const handleRemoveKeyDown =
    (index: number) => (e: React.KeyboardEvent<HTMLInputElement>) => {
      if (e.key === 'Backspace' && !values[index]) {
        const next = [...values];
        next.splice(index, 1);
        onBlur?.(next);
        newInputRef.current?.focus();
      }
    };

  const handleAddKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      if (newValue) {
        onBlur([...values, newValue]);
        setNewValue('');
      }
    }
  };

  const handleAdd = (e: React.ChangeEvent<HTMLInputElement>) => {
    setNewValue(e.target.value);
  };

  const handleAddBlur = () => {
    if (newValue) {
      onBlur?.([...values, newValue]);
      setNewValue('');
    }
  };

  return (
    <>
      {(values as string[])?.map((v, index) => (
        <SocialInput
          key={index}
          value={v}
          index={index}
          leftElement={leftElement}
          onBlur={handleBlur(index)}
          onChange={handleChange(index)}
          onKeyDown={handleRemoveKeyDown(index)}
        />
      ))}

      <InputGroup {...rest}>
        {leftElement && (
          <InputLeftElement>
            <SocialIcon url={newValue}>{leftElement}</SocialIcon>
          </InputLeftElement>
        )}
        <Input
          value={newValue}
          ref={newInputRef}
          onChange={handleAdd}
          onBlur={handleAddBlur}
          onKeyDown={handleAddKeyDown}
          {...rest}
        />
      </InputGroup>
    </>
  );
};
