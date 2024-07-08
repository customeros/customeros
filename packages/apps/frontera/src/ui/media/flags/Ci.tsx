import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Ci = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g fillRule='evenodd'>
      <path fill='#00cd00' d='M426.8 0H640v480H426.8z' />
      <path fill='#ff9a00' d='M0 0h212.9v480H0z' />
      <path fill='#fff' d='M212.9 0h214v480h-214z' />
    </g>
  </svg>
);
