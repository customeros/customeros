import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Cr = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g fillRule='evenodd' strokeWidth='1pt'>
      <path fill='#0000b4' d='M0 0h640v480H0z' />
      <path fill='#fff' d='M0 75.4h640v322.3H0z' />
      <path fill='#d90000' d='M0 157.7h640v157.7H0z' />
    </g>
  </svg>
);
