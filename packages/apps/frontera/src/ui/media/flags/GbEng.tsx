import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const GbEng = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#fff' d='M0 0h640v480H0z' />
    <path fill='#ce1124' d='M281.6 0h76.8v480h-76.8z' />
    <path fill='#ce1124' d='M0 201.6h640v76.8H0z' />
  </svg>
);
