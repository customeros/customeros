import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Lc = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g fillRule='evenodd'>
      <path fill='#65cfff' d='M0 0h640v480H0z' />
      <path fill='#fff' d='m318.9 42 162.7 395.3-322.6.9z' />
      <path fill='#000001' d='m319 96.5 140.8 340-279 .8z' />
      <path fill='#ffce00' d='m318.9 240.1 162.7 197.6-322.6.5z' />
    </g>
  </svg>
);
