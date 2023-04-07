import React, {
  ButtonHTMLAttributes,
  FC,
  ReactEventHandler,
  ReactNode,
} from 'react';
import styles from './icon-button.module.scss';
import classNames from 'classnames';

interface Props extends ButtonHTMLAttributes<HTMLButtonElement> {
  icon?: ReactNode;
  onClick: ReactEventHandler;
  isSquare?: boolean;
  mode?:
    | 'default'
    | 'primary'
    | 'secondary'
    | 'accent'
    | 'text'
    | 'danger'
    | 'dangerLink';
  size?: 'xxxxs' | 'xxxs' | 'xxs' | 'xs' | 'sm' | 'md' | 'lg';
}

export const IconButton: FC<Props> = ({
  icon,
  onClick,
  mode = 'default',
  size = 'xxs',
  isSquare = false,
  ...rest
}) => {
  return (
    <button
      {...rest}
      onClick={onClick}
      role={rest?.role || 'button'}
      title={rest?.title}
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
