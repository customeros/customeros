import React, { ChangeEventHandler, EventHandler, SyntheticEvent } from 'react';
import styles from './input.module.scss';
import classNames from 'classnames';
interface InputProps extends Partial<HTMLInputElement> {
  onChange: ChangeEventHandler<HTMLInputElement>;
  onClick?: EventHandler<SyntheticEvent>;
  icon?: React.ReactNode;
  label: string;
  error?: string;
  inputSize?: 'xxxs' | 'xxs' | 'xs' | 'sm' | 'md' | 'lg';
}

export const Input = ({
  onChange,
  placeholder = '',
  icon,
  label,
  error,
  inputSize = 'xxs',
}: InputProps) => {
  return (
    <>
      <label className={styles.wrapper}>
        <div> {label} </div>
        <input
          className={classNames(styles.input, styles[inputSize], {
            [styles.error]: !!error,
          })}
          onChange={onChange}
          placeholder={placeholder}
        />

        {icon && <span className={styles.icon}>{icon}</span>}
      </label>
      {error}
    </>
  );
};
