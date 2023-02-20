import React, { ChangeEventHandler } from 'react';
import { DebounceInput } from 'react-debounce-input';
import styles from './input.module.scss';
interface DebouncedInputProps {
  value?: string;
  onChange: ChangeEventHandler<HTMLInputElement>;
  placeholder?: string;
  minLength?: number;
  children?: React.ReactNode;
}

export const DebouncedInput = ({
  onChange,
  placeholder = '',
  minLength = 3,
  children,
}: DebouncedInputProps) => {
  return (
    <div className={styles.wrapper}>
      <DebounceInput
        className={styles.input}
        minLength={minLength}
        debounceTimeout={300}
        onChange={onChange}
        placeholder={placeholder}
      />

      {children && <span className={styles.icon}>{children}</span>}
    </div>
  );
};
