import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Cz = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#fff' d='M0 0h640v240H0z' />
    <path fill='#d7141a' d='M0 240h640v240H0z' />
    <path fill='#11457e' d='M360 240 0 0v480z' />
  </svg>
);
