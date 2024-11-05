import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const RadioButton = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 16 16'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <rect
      x='0.665'
      y='0.665'
      rx='7.335'
      width='14.67'
      height='14.67'
      stroke-width='1'
      stroke='currentColor'
    />
    <circle r='3' cx='8' cy='8' fill='currentColor' />
  </svg>
);
