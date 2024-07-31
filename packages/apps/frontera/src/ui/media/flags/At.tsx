import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const At = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#fff' d='M0 160h640v160H0z' />
    <path fill='#c8102e' d='M0 0h640v160H0zm0 320h640v160H0z' />
  </svg>
);
