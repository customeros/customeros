import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Lv = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g fillRule='evenodd'>
      <path fill='#fff' d='M0 0h640v480H0z' />
      <path fill='#981e32' d='M0 0h640v192H0zm0 288h640v192H0z' />
    </g>
  </svg>
);
