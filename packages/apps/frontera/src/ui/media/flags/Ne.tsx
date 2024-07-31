import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Ne = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#0db02b' d='M0 0h640v480H0z' />
    <path fill='#fff' d='M0 0h640v320H0z' />
    <path fill='#e05206' d='M0 0h640v160H0z' />
    <circle r='68' cx='320' cy='240' fill='#e05206' />
  </svg>
);
