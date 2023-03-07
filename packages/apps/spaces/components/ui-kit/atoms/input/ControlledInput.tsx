import React from 'react';
import { useController } from 'react-hook-form';
import classNames from 'classnames';
import styles from './input.module.scss';

interface ControlledInputProps {
  control: any;
  name: string;
  placeholder: string;
  label: string;
  hideLabel?: boolean;
  id?: string;
  required?: boolean;
  inputSize?: 'xxxs' | 'xxs' | 'xs' | 'sm' | 'md' | 'lg';
}

export const ControlledInput: React.FC<ControlledInputProps> = ({
  control,
  name,
  label,
  placeholder,
  required = false,
  inputSize = 'xxs',
  hideLabel = false,
  id,
}) => {
  const {
    field,
    fieldState: { invalid, isTouched, isDirty },
    formState: { touchedFields, dirtyFields },
  } = useController({
    name,
    control,
    rules: { required },
  });

  return (
    <>
      <label
        htmlFor={id}
        className={classNames(styles.label, {
          [styles.labelHidden]: hideLabel,
        })}
      >
        {label}
      </label>
      <input
        {...field}
        id={id}
        className={classNames(styles.input, {
          [styles.error]: invalid && isTouched,
          [styles?.[inputSize]]: inputSize,
        })}
        placeholder={placeholder}
      />
    </>
  );
};
