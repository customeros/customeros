import React, { EventHandler, FC, ReactNode } from 'react';
import styles from './icon-button.module.scss';

interface Props {
  icon?: ReactNode;
  onClick: EventHandler<any>;
  ariaLabel?: string;
  children?: HTMLCollection | undefined;
  mode?: 'default' | 'primary' | 'secondary';
  disabled?: boolean;
  className?: string;
  title?: string;
  style?: any;
  type?: string;
}

export const IconButton: FC<Props> = ({
  icon,
  onClick,
  children,
  mode = 'default',
  ...rest
}) => {
  return (
    <div
      {...rest}
      onClick={onClick}
      role='button'
      title={rest.ariaLabel}
      tabIndex={0}
      style={rest?.style}
      className={`${styles.button} ${styles[mode]} ${rest.className}`}
    >
      <>
        {icon && icon}
        {children}
      </>
    </div>
  );
};
