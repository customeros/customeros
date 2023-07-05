import { ButtonHTMLAttributes, ReactNode } from 'react';
import classNames from 'classnames';

import styles from './Button.module.scss';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  icon?: ReactNode;
  ariaLabel?: string;
  colorScheme?:
    | 'default'
    | 'primary'
    | 'secondary'
    | 'danger'
    | 'link'
    | 'dangerLink'
    | 'text';
}

export const Button = ({
  icon,
  onClick,
  children,
  colorScheme = 'default',
  ...rest
}: ButtonProps) => {
  return (
    <button
      {...rest}
      onClick={onClick}
      className={classNames(styles.button, styles[colorScheme], rest.className)}
    >
      {icon && icon}
      {children}
    </button>
  );
};
