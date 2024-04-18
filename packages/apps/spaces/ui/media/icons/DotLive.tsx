import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const DotLive = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 24 24'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <circle cx='12' cy='12' r='8.25' stroke='currentColor' strokeWidth='0.5' />
    <circle cx='12' cy='12' r='4' fill='currentColor' />
  </svg>
);
