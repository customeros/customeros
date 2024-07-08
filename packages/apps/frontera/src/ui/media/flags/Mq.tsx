import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Mq = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#231f1e' d='M0 0h640v480H0z' />
    <path fill='#00a650' d='M0 0h640v240H0z' />
    <path fill='#ef1923' d='m0 0 320 240L0 480z' />
  </svg>
);
