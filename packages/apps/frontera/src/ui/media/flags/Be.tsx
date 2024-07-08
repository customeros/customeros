import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Be = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g fillRule='evenodd' strokeWidth='1pt'>
      <path fill='#000001' d='M0 0h213.3v480H0z' />
      <path fill='#ffd90c' d='M213.3 0h213.4v480H213.3z' />
      <path fill='#f31830' d='M426.7 0H640v480H426.7z' />
    </g>
  </svg>
);
