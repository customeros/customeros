import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Bubbles = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 24 24'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <circle cx='6.5' cy='7.5' r='2.5' stroke='currentColor' strokeWidth='2' />
    <circle
      cx='16.5'
      cy='5.5'
      r='1.5'
      fill='currentColor'
      stroke='currentColor'
      strokeWidth='2'
    />
    <circle cx='14.5' cy='15.5' r='4.5' stroke='currentColor' strokeWidth='2' />
  </svg>
);
