import React, {
  ButtonHTMLAttributes,
  EventHandler,
  FC,
  ReactNode,
} from 'react';
import styles from './icon-button.module.scss';

interface Props extends ButtonHTMLAttributes<HTMLButtonElement> {
  icon?: ReactNode;
  onClick: EventHandler<never>;
  mode?: 'default' | 'primary' | 'secondary' | 'text';
  size?: 'xs' | 'sm' | 'md' | 'lg';
}

export const IconButton: FC<Props> = ({
  icon,
  onClick,
  mode = 'default',
  ...rest
}) => {
  return (
    <button
      {...rest}
      onClick={onClick}
      role={rest?.role || 'button'}
      title={rest?.title}
      tabIndex={0}
      style={rest?.style}
      className={`${styles.button} ${styles[mode]} ${rest.className}`}
    >
      {icon && icon}
    </button>
  );
};
