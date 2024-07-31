import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Dk = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#c8102e' d='M0 0h640.1v480H0z' />
    <path fill='#fff' d='M205.7 0h68.6v480h-68.6z' />
    <path fill='#fff' d='M0 205.7h640.1v68.6H0z' />
  </svg>
);
