import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Ga = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g fillRule='evenodd'>
      <path fill='#ffe700' d='M640 480H0V0h640z' />
      <path fill='#36a100' d='M640 160H0V0h640z' />
      <path fill='#006dbc' d='M640 480H0V320h640z' />
    </g>
  </svg>
);
