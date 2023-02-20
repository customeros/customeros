import React, { ChangeEventHandler } from 'react';
import { DebounceInput } from 'react-debounce-input';
import styles from './input.module.scss';
interface DebouncedInputProps {
  value?: string;
  onChange: ChangeEventHandler<HTMLInputElement>;
  placeholder?: string;
  minLength?: number;
}

export const DebouncedInput = ({
  onChange,
  placeholder = '',
  minLength = 3,
}: DebouncedInputProps) => {
  return (
    <>
      <DebounceInput
        className={styles.input}
        minLength={minLength}
        debounceTimeout={300}
        onChange={onChange}
        placeholder={placeholder}
      />
    </>
  );
};
