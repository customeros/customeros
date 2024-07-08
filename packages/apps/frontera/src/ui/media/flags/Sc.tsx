import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Sc = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#fff' d='M0 0h640v480H0Z' />
    <path fill='#d92223' d='M0 480V0h640v160z' />
    <path fill='#fcd955' d='M0 480V0h426.7z' />
    <path fill='#003d88' d='M0 480V0h213.3z' />
    <path fill='#007a39' d='m0 480 640-160v160z' />
  </svg>
);
