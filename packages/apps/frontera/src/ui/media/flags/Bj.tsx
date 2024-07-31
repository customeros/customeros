import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Bj = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='bj-a'>
        <path fill='gray' d='M67.6-154h666v666h-666z' />
      </clipPath>
    </defs>
    <g clipPath='url(#bj-a)' transform='matrix(.961 0 0 .7207 -65 111)'>
      <g strokeWidth='1pt' fillRule='evenodd'>
        <path fill='#319400' d='M0-154h333v666H0z' />
        <path fill='#ffd600' d='M333-154h666v333H333z' />
        <path fill='#de2110' d='M333 179h666v333H333z' />
      </g>
    </g>
  </svg>
);
