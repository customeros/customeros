import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Mv = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#d21034' d='M0 0h640v480H0z' />
    <path fill='#007e3a' d='M120 120h400v240H120z' />
    <circle cx='350' cy='240' r='80' fill='#fff' />
    <circle cx='380' cy='240' r='80' fill='#007e3a' />
  </svg>
);
