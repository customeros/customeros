import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Microsoft = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 32 32'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <rect x='17' y='17' width='10' height='10' fill='#FEBA08' />
    <rect x='5' y='17' width='10' height='10' fill='#05A6F0' />
    <rect y='5' x='17' width='10' height='10' fill='#80BC06' />
    <rect x='5' y='5' width='10' height='10' fill='#F25325' />
  </svg>
);
