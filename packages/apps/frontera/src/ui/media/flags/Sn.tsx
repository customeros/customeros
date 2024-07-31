import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Sn = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g strokeWidth='1pt' fillRule='evenodd'>
      <path fill='#0b7226' d='M0 0h213.3v480H0z' />
      <path fill='#ff0' d='M213.3 0h213.3v480H213.3z' />
      <path fill='#bc0000' d='M426.6 0H640v480H426.6z' />
    </g>
    <path
      fill='#0b7226'
      d='M342 218.8h71.8l-56.6 43.6 20.7 69.3-56.6-43.6-56.6 41.6 20.7-67.3-56.6-43.6h69.8l22.7-71.3z'
    />
  </svg>
);
