import React from 'react';
import { useController } from 'react-hook-form';
import classNames from 'classnames';
import styles from './input.module.scss';

interface ControlledTextAreaProps {
  control: any;
  name: string;
  placeholder: string;
  label: string;
  required?: boolean;
  inputSize?: 'xxxs' | 'xxs' | 'xs' | 'sm' | 'md' | 'lg';
}

export const ControlledTextArea: React.FC<ControlledTextAreaProps> = ({
  control,
  name,
  label,
  placeholder,
  required = false,
  inputSize = 'xxs',
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
      <label className={styles.label}>{label}</label>
      <textarea
        {...field}
        className={classNames(styles.input, {
          [styles.error]: invalid && isTouched,
          [styles?.[inputSize]]: inputSize,
        })}
        placeholder={placeholder}
      />
    </>
  );
};
