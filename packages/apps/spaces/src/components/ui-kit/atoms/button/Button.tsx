import React, {
  ButtonHTMLAttributes,
  EventHandler,
  FC,
  ReactNode,
} from 'react';
import styles from './button.module.scss';

interface Props extends ButtonHTMLAttributes<HTMLButtonElement> {
  icon?: ReactNode;
  onClick: EventHandler<never>;
  ariaLabel?: string;
  children?: React.ReactNode;
  mode?:
    | 'default'
    | 'primary'
    | 'secondary'
    | 'danger'
    | 'link'
    | 'dangerLink'
    | 'text';
}

export const Button: FC<Props> = ({
  icon,
  onClick,
  children,
  mode = 'default',
  ...rest
}) => {
  return (
    <button
      {...rest}
      onClick={onClick}
      className={`${styles.button} ${styles[mode]}`}
    >
      <>
        {icon && icon}
        {children}
      </>
    </button>
  );
};
