import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Fi = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#fff' d='M0 0h640v480H0z' />
    <path fill='#002f6c' d='M0 174.5h640v131H0z' />
    <path fill='#002f6c' d='M175.5 0h130.9v480h-131z' />
  </svg>
);
