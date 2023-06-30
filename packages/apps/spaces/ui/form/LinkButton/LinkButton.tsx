import { ReactNode, PropsWithChildren } from 'react';
import classNames from 'classnames';
import Link, { LinkProps } from 'next/link';

import styles from './LinkButton.module.scss';

interface LinkButtonProps extends LinkProps, PropsWithChildren {
  icon?: ReactNode;
  ariaLabel?: string;
  className?: string;
  colorScheme?:
    | 'default'
    | 'primary'
    | 'secondary'
    | 'danger'
    | 'link'
    | 'dangerLink'
    | 'text';
}

export const LinkButton = ({
  icon,
  onClick,
  children,
  colorScheme = 'default',
  ...rest
}: LinkButtonProps) => {
  return (
    <Link
      {...rest}
      onClick={onClick}
      className={classNames(styles.button, styles[colorScheme], rest.className)}
    >
      <>
        {icon && icon}
        {children}
      </>
    </Link>
  );
};
