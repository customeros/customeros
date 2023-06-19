import React, {
  ButtonHTMLAttributes,
  FC,
  MouseEventHandler,
  ReactNode,
} from 'react';
import styles from './icon-button.module.scss';
import classNames from 'classnames';

interface Props extends ButtonHTMLAttributes<HTMLButtonElement> {
  icon?: ReactNode;
  onClick?: MouseEventHandler<HTMLButtonElement>;
  isSquare?: boolean;
  mode?:
    | 'default'
    | 'primary'
    | 'secondary'
    | 'accent'
    | 'text'
    | 'danger'
    | 'success'
    | 'subtle'
    | 'dangerLink';
  size?: 'xxxxs' | 'xxxs' | 'xxs' | 'xs' | 'sm' | 'md' | 'lg';
  label: string;
}

export const IconButton: FC<Props> = ({
  icon,
  onClick,
  mode = 'default',
  size = 'xxs',
  isSquare = false,
  label,
  ...rest
}) => {
  return (
    <button
      {...rest}
      onClick={onClick}
      role={rest?.role || 'button'}
      aria-label={label}
      style={rest?.style}
      className={classNames(
        styles.button,
        styles[mode],
        styles[size],
        rest.className,
        {
          [styles.square]: isSquare,
        },
      )}
    >
      {icon && icon}
    </button>
  );
};
