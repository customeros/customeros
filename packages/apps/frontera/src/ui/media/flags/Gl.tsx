import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Gl = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#fff' d='M0 0h640v480H0z' />
    <path
      fill='#d00c33'
      d='M0 240h640v240H0zm80 0a160 160 0 1 0 320 0 160 160 0 0 0-320 0'
    />
  </svg>
);
