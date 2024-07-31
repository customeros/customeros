import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Bl = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#fff' d='M0 0h640v480H0z' />
    <path fill='#000091' d='M0 0h213.3v480H0z' />
    <path fill='#e1000f' d='M426.7 0H640v480H426.7z' />
  </svg>
);
