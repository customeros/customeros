import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Id = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#e70011' d='M0 0h640v240H0Z' />
    <path fill='#fff' d='M0 240h640v240H0Z' />
  </svg>
);
