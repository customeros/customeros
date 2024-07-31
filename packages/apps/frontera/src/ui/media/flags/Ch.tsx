import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Ch = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g strokeWidth='1pt' fillRule='evenodd'>
      <path fill='red' d='M0 0h640v480H0z' />
      <g fill='#fff'>
        <path d='M170 195h300v90H170z' />
        <path d='M275 90h90v300h-90z' />
      </g>
    </g>
  </svg>
);
