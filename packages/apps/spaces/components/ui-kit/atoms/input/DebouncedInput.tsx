import React, { ChangeEventHandler } from 'react';
import { DebounceInput, DebounceInputProps } from 'react-debounce-input';
import styles from './input.module.scss';
import classNames from 'classnames';
interface DebouncedInputProps
  extends Partial<DebounceInputProps<HTMLInputElement, HTMLInputElement>> {
  onChange: ChangeEventHandler<HTMLInputElement>;
  placeholder?: string;
  minLength?: number;
  inputSize?: 'xxxs' | 'xxs' | 'xs' | 'sm' | 'md' | 'lg';
}

export const DebouncedInput = ({
  onChange,
  placeholder = '',
  minLength = 3,
  children,
  inputSize = 'md',
  ...rest
}: DebouncedInputProps) => {
  return (
    <div className={styles.wrapper}>
      <DebounceInput
        className={classNames(styles.input, {
          [styles?.[inputSize]]: inputSize,
        })}
        minLength={minLength}
        debounceTimeout={300}
        onChange={onChange}
        placeholder={placeholder}
        {...rest}
      />

      {children && <span className={styles.icon}>{children}</span>}
    </div>
  );
};
