import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Bubbles = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 24 24'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <circle r='2.5' cx='6.5' cy='7.5' strokeWidth='2' stroke='currentColor' />
    <circle
      r='1.5'
      cy='5.5'
      cx='16.5'
      strokeWidth='2'
      fill='currentColor'
      stroke='currentColor'
    />
    <circle r='4.5' cx='14.5' cy='15.5' strokeWidth='2' stroke='currentColor' />
  </svg>
);
