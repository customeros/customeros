import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Cg = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='cg-a'>
        <path fillOpacity='.7' d='M-79.5 32h640v480h-640z' />
      </clipPath>
    </defs>
    <g
      strokeWidth='1pt'
      fillRule='evenodd'
      clipPath='url(#cg-a)'
      transform='translate(79.5 -32)'
    >
      <path fill='#ff0' d='M-119.5 32h720v480h-720z' />
      <path fill='#00ca00' d='M-119.5 32v480l480-480z' />
      <path fill='red' d='M120.5 512h480V32z' />
    </g>
  </svg>
);
