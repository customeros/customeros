import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Ml = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g fillRule='evenodd'>
      <path fill='red' d='M425.8 0H640v480H425.7z' />
      <path fill='#009a00' d='M0 0h212.9v480H0z' />
      <path fill='#ff0' d='M212.9 0h214v480h-214z' />
    </g>
  </svg>
);
