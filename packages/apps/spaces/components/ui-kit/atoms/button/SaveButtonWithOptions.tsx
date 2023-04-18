import React, { FC, ReactNode, ButtonHTMLAttributes, useRef } from 'react';
import { Menu } from 'primereact/menu';
import { Send } from '../icons';
import { IconButton } from '../icon-button';

interface Props extends ButtonHTMLAttributes<HTMLButtonElement> {
  icon?: ReactNode;
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

export const SaveButtonWithOptions: FC<any> = ({
  icon,
  onClick,
  children,
  mode = 'default',
  label,
  ...rest
}) => {
  const menu = useRef(null);

  return (
    <>
      <Menu model={rest.items} popup ref={menu} />
      <IconButton
        mode='primary'
        size='xxxxs'
        icon={<Send style={{ transform: 'scale(0.9)' }} />}
        // @ts-expect-error fixme
        onClick={(e) => menu?.current?.toggle(e)}
        style={{
          borderRadius: '2px',
        }}
      />
    </>
  );
};
