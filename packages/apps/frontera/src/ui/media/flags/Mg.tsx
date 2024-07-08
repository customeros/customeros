import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Mg = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g fillRule='evenodd' strokeWidth='1pt'>
      <path fill='#fc3d32' d='M213.3 0H640v240H213.3z' />
      <path fill='#007e3a' d='M213.3 240H640v240H213.3z' />
      <path fill='#fff' d='M0 0h213.3v480H0z' />
    </g>
  </svg>
);
