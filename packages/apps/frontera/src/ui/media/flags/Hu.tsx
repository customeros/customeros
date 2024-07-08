import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Hu = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g fillRule='evenodd'>
      <path fill='#fff' d='M640 480H0V0h640z' />
      <path fill='#388d00' d='M640 480H0V320h640z' />
      <path fill='#d43516' d='M640 160.1H0V.1h640z' />
    </g>
  </svg>
);
