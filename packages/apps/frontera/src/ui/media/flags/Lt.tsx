import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Lt = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g strokeWidth='1pt' fillRule='evenodd' transform='scale(.64143 .96773)'>
      <rect
        rx='0'
        ry='0'
        width='1063'
        height='708.7'
        fill='#006a44'
        transform='scale(.93865 .69686)'
      />
      <rect
        rx='0'
        ry='0'
        y='475.6'
        width='1063'
        height='236.2'
        fill='#c1272d'
        transform='scale(.93865 .69686)'
      />
      <path fill='#fdb913' d='M0 0h997.8v164.6H0z' />
    </g>
  </svg>
);
