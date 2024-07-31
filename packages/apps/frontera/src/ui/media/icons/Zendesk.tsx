import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Zendesk = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 32 32'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#07363D' d='M17 5L17 20.7143L30 5H17Z' />
    <path
      fill='#07363D'
      d='M8.50003 12C12.0562 12 15 8.41579 15 5.00002H2.00002C2.00002 8.41579 4.94386 12 8.50003 12Z'
    />
    <path
      fill='#07363D'
      d='M17 27C17 23.5842 19.9439 20 23.5 20C27.0562 20 30 23.5842 30 27H17Z'
    />
    <path fill='#07363D' d='M15 27V11.2857L2 27H15Z' />
  </svg>
);
