import React, { ChangeEventHandler, EventHandler, SyntheticEvent } from 'react';
import styles from './input.module.scss';
import classNames from 'classnames';
interface InputProps extends Partial<HTMLInputElement> {
  onChange: ChangeEventHandler<HTMLInputElement>;
  onClick?: EventHandler<SyntheticEvent>;
  icon?: React.ReactNode;
  label: string;
  error?: string;
}

export const Input = ({
  onChange,
  placeholder = '',
  icon,
  label,
  error,
}: InputProps) => {
  return (
    <>
      <label className={styles.wrapper}>
        <div> {label} </div>
        <input
          className={classNames(styles.input, {
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
