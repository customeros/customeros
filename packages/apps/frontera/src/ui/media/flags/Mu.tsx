import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Mu = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g fillRule='evenodd'>
      <path fill='#00a04d' d='M0 360h640v120H0z' />
      <path fill='#151f6d' d='M0 120h640v120H0z' />
      <path fill='#ee2737' d='M0 0h640v120H0z' />
      <path fill='#ffcd00' d='M0 240h640v120H0z' />
    </g>
  </svg>
);
