import React, { ChangeEventHandler, ReactNode } from 'react';
import { DebounceInput, DebounceInputProps } from 'react-debounce-input';
import styles from './input.module.scss';
import classNames from 'classnames';
interface DebouncedInputProps
  extends Partial<
    Omit<DebounceInputProps<HTMLInputElement, HTMLInputElement>, 'children'>
  > {
  onChange: ChangeEventHandler<HTMLInputElement>;
  placeholder?: string;
  minLength?: number;
  debounceTimeout?: number;
  inputSize?: 'xxxs' | 'xxs' | 'xs' | 'sm' | 'md' | 'lg';
  children?: ReactNode;
}

export const DebouncedInput = ({
  onChange,
  placeholder = '',
  minLength = 3,
  children,
  inputSize = 'md',
  debounceTimeout = 300,
  ...rest
}: DebouncedInputProps) => {
  return (
    <div
      className={classNames(styles.wrapper, {
        //@ts-expect-error fixme
        [styles?.[rest.className]]: rest?.className,
      })}
    >
      <DebounceInput
        {...rest}
        className={classNames(styles.input, {
          [styles?.[inputSize]]: inputSize,
        })}
        minLength={minLength}
        debounceTimeout={debounceTimeout}
        onChange={onChange}
        placeholder={placeholder}
      />

      {children && <span className={styles.icon}>{children}</span>}
    </div>
  );
};
