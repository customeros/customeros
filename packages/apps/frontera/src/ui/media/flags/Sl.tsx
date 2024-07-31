import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Sl = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g fillRule='evenodd'>
      <path fill='#0000cd' d='M0 320.3h640V480H0z' />
      <path fill='#fff' d='M0 160.7h640v159.6H0z' />
      <path fill='#00cd00' d='M0 0h640v160.7H0z' />
    </g>
  </svg>
);
