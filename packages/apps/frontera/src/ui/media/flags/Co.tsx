import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Co = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g fillRule='evenodd' strokeWidth='1pt'>
      <path fill='#ffe800' d='M0 0h640v480H0z' />
      <path fill='#00148e' d='M0 240h640v240H0z' />
      <path fill='#da0010' d='M0 360h640v120H0z' />
    </g>
  </svg>
);
