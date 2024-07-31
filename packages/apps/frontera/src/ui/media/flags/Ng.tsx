import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Ng = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g strokeWidth='1pt' fillRule='evenodd'>
      <path fill='#fff' d='M0 0h640v480H0z' />
      <path fill='#008753' d='M426.6 0H640v480H426.6zM0 0h213.3v480H0z' />
    </g>
  </svg>
);
