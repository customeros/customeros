import React, { ChangeEventHandler, CSSProperties, ReactNode } from 'react';
import { DebounceInput, DebounceInputProps } from 'react-debounce-input';
import styles from './input.module.scss';
import classNames from 'classnames';
import { InlineLoader } from '@ui/presentation/inline-loader';

interface DebouncedInputProps
  extends Partial<
    Omit<DebounceInputProps<HTMLInputElement, HTMLInputElement>, 'children'>
  > {
  onChange: ChangeEventHandler<HTMLInputElement>;
  onClick?: ChangeEventHandler<HTMLInputElement>;
  placeholder?: string;
  minLength?: number;
  debounceTimeout?: number;
  inputSize?: 'xxxs' | 'xxs' | 'xs' | 'sm' | 'md' | 'lg';
  children?: ReactNode;
  inlineMode?: boolean;
  saving?: boolean;
  className?: string;
  customStyles?: CSSProperties | undefined;
  inputRef?: any;
}

export const DebouncedInput = ({
  onChange,
  onClick,
  placeholder = '',
  minLength = 3,
  children,
  inputSize = 'md',
  debounceTimeout = 300,
  inlineMode,
  className = '',
  customStyles,
  inputRef,
  saving,
  ...rest
}: DebouncedInputProps) => {
  return (
    <div
      className={classNames(styles.wrapper, {
        [styles.inlineMode]: inlineMode,
      })}
    >
      <DebounceInput
        {...rest}
        size={rest?.value?.length || placeholder?.length}
        inputRef={inputRef}
        className={classNames(styles.input, {
          [styles?.[inputSize]]: inputSize,
          [styles.xxxs]: inlineMode,
          [`${className}`]: className !== '',
        })}
        minLength={minLength}
        debounceTimeout={debounceTimeout}
        onChange={onChange}
        onClick={onClick}
        placeholder={placeholder}
        spellCheck={false}
        autoCorrect={'off'}
        autoCapitalize={'off'}
        style={customStyles}
      />

      {children && <span className={styles.icon}>{children}</span>}
      {saving && <InlineLoader />}
    </div>
  );
};
