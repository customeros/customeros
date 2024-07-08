import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Ae = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#00732f' d='M0 0h640v160H0z' />
    <path fill='#fff' d='M0 160h640v160H0z' />
    <path fill='#000001' d='M0 320h640v160H0z' />
    <path fill='red' d='M0 0h220v480H0z' />
  </svg>
);
