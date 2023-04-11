import React, { FC, InputHTMLAttributes } from 'react';
import styles from './checkbox.module.scss';
import { ChangeHandler } from 'react-hook-form';

interface Props extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  type: 'checkbox' | 'radio' | '';
  onChange: ChangeHandler;
}

export const Checkbox: FC<Props> = ({ label, type, ...rest }) => {
  return (
    <label className={styles.label}>
      <input
        type={type}
        className={styles.input}
        {...rest}
        onKeyDown={(e) => e.key === 'Enter' && rest.onChange(e)}
      />
      <span>{label}</span>
    </label>
  );
};
